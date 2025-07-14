package internal

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var errorTmpl = template.Must(template.ParseFiles(filepath.Join("templates", "error.html")))

type ErrorData struct {
	Code    int
	Title   string
	Message string
}

func RenderError(w http.ResponseWriter, code int, title, message string) {
	w.WriteHeader(code)
	_ = errorTmpl.Execute(w, ErrorData{
		Code:    code,
		Title:   title,
		Message: message,
	})
}
