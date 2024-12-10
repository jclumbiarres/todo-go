package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"lumbi.net/practice/todo-go/internal/model"
	"lumbi.net/practice/todo-go/pkg/server"
)

func AddTodo(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var insertData model.Todo
		err := json.NewDecoder(r.Body).Decode(&insertData)
		if err != nil {
			server.ToError(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(201)
		srv.DB.Create(&insertData)
	}
}

func GetAllTodos(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todos []model.Todo
		res := srv.DB.Find(&todos)
		if res.Error != nil {
			server.ToError(w, "fallo al buscar todos", http.StatusNotFound)
			return
		}
		if res.RowsAffected == 0 {
			server.ToError(w, "no se han encontrado todos", http.StatusNotFound)
			return
		}
		server.ToJSON(w, todos)
	}
}

func GetTodo(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonData model.Todo
		nan, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			server.ToError(w, "error al recuperar el todo", http.StatusBadRequest)
			return
		}
		salida := srv.DB.Where("id = ?", nan).First(&jsonData)
		if salida.RowsAffected == 0 {
			server.ToError(w, "no encontrado", http.StatusNotFound)
			return
		}
		// if reflect.ValueOf(jsonData).IsZero() {
		// 	http.Error(w, "no encontrado", http.StatusNotFound)
		// 	return
		// }
		server.ToJSON(w, jsonData)
	}
}

func CompleteTodo(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonData model.Todo
		nan, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			server.ToError(w, "error al recuperar el todo", http.StatusBadRequest)
			return
		}
		salida := srv.DB.Where("id = ?", nan).First(&jsonData)
		if salida.RowsAffected == 0 {
			server.ToError(w, "no encontrado", http.StatusNotFound)
			return
		}
		// if reflect.ValueOf(jsonData).IsZero() {
		// 	http.Error(w, "no encontrado", http.StatusNotFound)
		// 	return
		// }
		server.ToJSON(w, jsonData)
	}
}
