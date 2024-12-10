package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"lumbi.net/practice/todo-go/internal/model"
	"lumbi.net/practice/todo-go/pkg/server"
)

func setupTestServer(t *testing.T) *server.Server {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.Todo{})

	config := server.NewConfig().SetID("test").SetListenAddress(":0").SetName("test-server")
	return server.NewServer(config, db)
}

func TestGetAllTodos(t *testing.T) {
	srv := setupTestServer(t)

	// Crear algunos todos de prueba
	todos := []model.Todo{
		{Title: "Test 1", Completed: false},
		{Title: "Test 2", Completed: true},
	}
	for _, todo := range todos {
		srv.DB.Create(&todo)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v2/todo", nil)
	w := httptest.NewRecorder()

	handler := GetAllTodos(srv)
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.Todo
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestGetTodo(t *testing.T) {
	srv := setupTestServer(t)

	// Crear un todo de prueba
	todo := model.Todo{Title: "Test", Completed: false}
	srv.DB.Create(&todo)

	req := httptest.NewRequest(http.MethodGet, "/api/v2/todo/1", nil)
	w := httptest.NewRecorder()

	handler := GetTodo(srv)
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Todo
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, todo.Title, response.Title)
}

func TestAddTodo(t *testing.T) {
	srv := setupTestServer(t)

	newTodo := model.Todo{
		Title:     "Nuevo Todo",
		Completed: false,
	}

	body, _ := json.Marshal(newTodo)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/todo", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := AddTodo(srv)
	handler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Todo
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, newTodo.Title, response.Title)
}
