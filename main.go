package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"lumbi.net/practice/todo-go/internal/controllers"
	"lumbi.net/practice/todo-go/internal/model"
	"lumbi.net/practice/todo-go/pkg/frontend"
	"lumbi.net/practice/todo-go/pkg/middleware"
	"lumbi.net/practice/todo-go/pkg/server"
)

func main() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Todo{})
	config := server.NewConfig().SetID("s01").SetListenAddress(":3000").SetName("todo-api")
	fs := http.FileServer(http.Dir("./static"))
	srv := server.NewServer(config, db)
	logger := middleware.NewColorLogger()
	srv.Use(middleware.LoggingMiddleware(logger))
	srv.Use(middleware.EnableCors)
	srv.Handle("GET /static/", http.StripPrefix("/static/", fs))
	srv.Route("POST /api/v2/todo", controllers.AddTodo(srv))
	srv.Route("GET /api/v2/todo", controllers.GetAllTodos(srv))
	srv.Route("GET /api/v2/todo/{id}", controllers.GetTodo(srv))
	srv.Route("GET /", func(w http.ResponseWriter, r *http.Request) {
		var todos []model.Todo
		res := srv.DB.Find(&todos)
		if res.Error != nil {
			server.ToError(w, "fallo al buscar todos", http.StatusNotFound)
			return
		}
		if res.RowsAffected == 0 {
			templ.Handler(frontend.Layout([]model.Todo{})).ServeHTTP(w, r)
			return
		} else {
			templ.Handler(frontend.Layout(todos)).ServeHTTP(w, r)
		}
	})
	srv.Route("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		var todos []model.Todo
		res := srv.DB.Find(&todos)
		if res.Error != nil {
			server.ToError(w, "fallo al buscar todos", http.StatusNotFound)
			return
		}
		if res.RowsAffected == 0 {
			server.ToError(w, "no se han encontrado todos", http.StatusNotFound)
			return
		} else {
			templ.Handler(frontend.TodoList(todos)).ServeHTTP(w, r)
		}
	})
	srv.Route("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			taskText := r.FormValue("todo")
			newTask := model.Todo{
				Title:     taskText,
				Completed: false,
			}
			srv.DB.Create(&newTask)
			templ.Handler(frontend.List([]model.Todo{newTask})).ServeHTTP(w, r)
		}
	})
	// Ruta para DELETE
	srv.Route("DELETE /todos/delete/", func(w http.ResponseWriter, r *http.Request) {

		idStr := strings.TrimPrefix(r.URL.Path, "/todos/delete/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			server.ToError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		res := srv.DB.Delete(&model.Todo{}, id)
		if res.Error != nil {
			server.ToError(w, "No se pudo borrar el todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Ruta para PATCH
	srv.Route("PATCH /todos/update/", func(w http.ResponseWriter, r *http.Request) {

		idStr := strings.TrimPrefix(r.URL.Path, "/todos/update/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			server.ToError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var todo model.Todo
		if err := srv.DB.First(&todo, id).Error; err != nil {
			server.ToError(w, "Todo no encontrado", http.StatusNotFound)
			return
		}

		todo.Completed = true
		srv.DB.Save(&todo)

		templ.Handler(frontend.TodoItem(todo)).ServeHTTP(w, r)
	})
	srv.Start()
}
