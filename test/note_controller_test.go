package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/milhamap/simple-note/app"
	"github.com/milhamap/simple-note/controller"
	"github.com/milhamap/simple-note/helper"
	"github.com/milhamap/simple-note/middleware"
	"github.com/milhamap/simple-note/model/domain"
	"github.com/milhamap/simple-note/repository"
	"github.com/milhamap/simple-note/service"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/simple_note_test")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	noteRepository := repository.NewNoteRepository()
	noteService := service.NewNoteService(noteRepository, db, validate)
	noteController := controller.NewNoteController(noteService)
	router := app.NewRouter(noteController)

	return middleware.NewAuthMiddleware(router)
}

func truncateNote(db *sql.DB) {
	db.Exec("TRUNCATE note")
}

func TestSuccessCreateNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title": "Belajar Golang", "content": "Belajar Golang dari dasar"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/notes", requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Belajar Golang", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, "Belajar Golang dari dasar", responseBody["data"].(map[string]interface{})["content"])
}

func TestFailCreateNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title":"", "content":""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/notes", requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestSuccessUpdateNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)

	tx, _ := db.Begin()
	noteRepository := repository.NewNoteRepository()
	note := noteRepository.Store(context.Background(), tx, domain.Note{
		Title:   "Belajar Golang",
		Content: "Belajar Golang dari dasar",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title": "Belajar Golang", "content": "Belajar Golang dari dasar dengan menggunakan goland"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/notes/"+strconv.Itoa(note.Id), requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, note.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Belajar Golang", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, "Belajar Golang dari dasar dengan menggunakan goland", responseBody["data"].(map[string]interface{})["content"])
}

func TestFailUpdateNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)

	tx, _ := db.Begin()
	noteRepository := repository.NewNoteRepository()
	note := noteRepository.Store(context.Background(), tx, domain.Note{
		Title:   "Belajar Golang",
		Content: "Belajar Golang dari dasar",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title":"", "content":""}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/notes/"+strconv.Itoa(note.Id), requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestSuccessGetNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)

	tx, _ := db.Begin()
	noteRepository := repository.NewNoteRepository()
	note := noteRepository.Store(context.Background(), tx, domain.Note{
		Title:   "Belajar Golang",
		Content: "Belajar Golang dari dasar",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/notes/"+strconv.Itoa(note.Id), nil)
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, note.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, note.Title, responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, note.Content, responseBody["data"].(map[string]interface{})["content"])
}

func TestFailGetNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/notes/404", nil)
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestSuccessDeleteNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)

	tx, _ := db.Begin()
	noteRepository := repository.NewNoteRepository()
	note := noteRepository.Store(context.Background(), tx, domain.Note{
		Title:   "Belajar Golang",
		Content: "Belajar Golang dari dasar",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/notes/"+strconv.Itoa(note.Id), nil)
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestFailDeleteNote(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/notes/404", nil)
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestSuccessListNotes(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)

	tx, _ := db.Begin()
	noteRepository := repository.NewNoteRepository()
	note1 := noteRepository.Store(context.Background(), tx, domain.Note{
		Title:   "Belajar Golang",
		Content: "Belajar Golang dari dasar",
	})
	note2 := noteRepository.Store(context.Background(), tx, domain.Note{
		Title:   "Belajar NodeJS",
		Content: "Belajar NodeJS dari dasar",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/notes", nil)
	request.Header.Set("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	var notes = responseBody["data"].([]interface{})

	noteResponse1 := notes[0].(map[string]interface{})
	noteResponse2 := notes[1].(map[string]interface{})

	assert.Equal(t, note1.Id, int(noteResponse1["id"].(float64)))
	assert.Equal(t, note1.Title, noteResponse1["title"])
	assert.Equal(t, note1.Content, noteResponse1["content"])

	assert.Equal(t, note2.Id, int(noteResponse2["id"].(float64)))
	assert.Equal(t, note2.Title, noteResponse2["title"])
	assert.Equal(t, note2.Content, noteResponse2["content"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	truncateNote(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/notes", nil)
	request.Header.Set("X-API-KEY", "SALAH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	responseBody := map[string]interface{}{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responseBody["status"])
}
