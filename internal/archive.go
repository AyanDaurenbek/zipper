package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	allowedExtensions = []string{".pdf", ".jpeg"}
	archiveDir        = "./static/archives"
)

func DownloadAndZip(taskID string, links []string) (string, []string, error) {
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return "", nil, err
	}

	archivePath := filepath.Join(archiveDir, taskID+".zip")
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return "", nil, err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	var badLinks []string

	for i, link := range links {
		filename := filepath.Base(link)
		ext := strings.ToLower(filepath.Ext(filename))

		if !isAllowedExtension(ext) {
			badLinks = append(badLinks, link+" (неподдерживаемый тип)")
			continue
		}

		resp, err := http.Get(link)
		if err != nil || resp.StatusCode != 200 {
			badLinks = append(badLinks, link+" (недоступен)")
			continue
		}
		defer resp.Body.Close()

		fw, err := zipWriter.Create("file" + fmt.Sprint('1'+i) + ext)
		if err != nil {
			badLinks = append(badLinks, link+" (ошибка архивации)")
			continue
		}

		_, err = io.Copy(fw, resp.Body)
		if err != nil {
			badLinks = append(badLinks, link+" (ошибка копирования)")
			continue
		}
	}

	return "/static/archives/" + taskID + ".zip", badLinks, nil
}

func isAllowedExtension(ext string) bool {
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}
