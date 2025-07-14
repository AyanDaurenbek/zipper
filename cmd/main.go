package main

import (
	"log"
	"net/http"
	"zipper/internal"
)

func main() {
	manager := internal.NewTaskManager()

	http.HandleFunc("/", manager.HandleHome)
	http.HandleFunc("/task/create", manager.HandleCreateTask)
	http.HandleFunc("/task/", manager.HandleTaskPage)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
