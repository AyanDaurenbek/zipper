package internal

import (
	"html/template"
	"net/http"
	"strings"
)

var (
	indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
	taskTmpl  = template.Must(template.ParseFiles("templates/task.html"))
)

func (m *TaskManager) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		RenderError(w, 404, "Страница не найдена", "Такой страницы не существует.")
		return
	}

	m.tasksMu.Lock()
	defer m.tasksMu.Unlock()

	if err := indexTmpl.Execute(w, m.tasks); err != nil {
		RenderError(w, 500, "Ошибка шаблона", "Не удалось отрендерить главную страницу.")
	}
}

func (m *TaskManager) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, 405, "Метод не поддерживается", "Ожидался POST-запрос.")
		return
	}

	task, err := m.CreateTask()
	if err != nil {
		RenderError(w, 503, "Сервер занят", "Достигнуто максимальное количество активных задач. Попробуйте позже.")
		return
	}

	http.Redirect(w, r, "/task/"+task.ID, http.StatusSeeOther)
}

func (m *TaskManager) HandleTaskPage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		RenderError(w, 404, "Неверный путь", "Страница задачи не найдена.")
		return
	}
	taskID := parts[2]

	task, ok := m.GetTask(taskID)
	if !ok {
		RenderError(w, 404, "Задача не найдена", "Указанная задача отсутствует.")
		return
	}

	if r.Method == http.MethodPost {
		link := r.FormValue("link")
		if link != "" {
			if err := m.AddLink(taskID, link); err != nil {
				RenderError(w, 400, "Ошибка добавления ссылки", err.Error())
				return
			}
			http.Redirect(w, r, "/task/"+taskID, http.StatusSeeOther)
			return
		}
	}

	if err := taskTmpl.Execute(w, task); err != nil {
		RenderError(w, 500, "Ошибка шаблона", "Не удалось отобразить задачу.")
	}
}
