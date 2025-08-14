package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	db.AutoMigrate(&Book{})
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/books", getBooks).Methods("GET")
	api.HandleFunc("/books/{id}", getBook).Methods("GET")
	api.HandleFunc("/books", createBook).Methods("POST")
	api.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	api.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	return r
}

func TestMain(m *testing.M) {
	setupTestDB()
	code := m.Run()
	os.Exit(code)
}

func clearDB() {
	db.Exec("DELETE FROM books")
	db.Exec("DELETE FROM sqlite_sequence WHERE name='books'")
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.Code)
	}
}

func TestGetBooksEmpty(t *testing.T) {
	clearDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/books", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.Code)
	}

	var books []Book
	json.Unmarshal(response.Body.Bytes(), &books)

	if len(books) != 0 {
		t.Errorf("Expected empty array, got %d books", len(books))
	}
}

func TestCreateBook(t *testing.T) {
	clearDB()
	router := setupRouter()

	book := Book{
		Title:  "Test Book",
		Author: "Test Author",
		ISBN:   "1234567890123",
		Year:   2023,
	}

	jsonData, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", response.Code)
	}

	var createdBook Book
	json.Unmarshal(response.Body.Bytes(), &createdBook)

	if createdBook.Title != book.Title {
		t.Errorf("Expected title %s, got %s", book.Title, createdBook.Title)
	}
	if createdBook.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func TestCreateBookInvalidData(t *testing.T) {
	router := setupRouter()

	book := Book{
		Title: "Test Book",
		// Missing Author and ISBN
	}

	jsonData, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", response.Code)
	}
}

func TestGetBookByID(t *testing.T) {
	clearDB()
	router := setupRouter()

	// First create a book
	book := Book{
		Title:  "Test Book",
		Author: "Test Author",
		ISBN:   "1234567890456",
		Year:   2023,
	}
	db.Create(&book)

	req, _ := http.NewRequest("GET", "/api/v1/books/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.Code)
	}

	var retrievedBook Book
	json.Unmarshal(response.Body.Bytes(), &retrievedBook)

	if retrievedBook.Title != book.Title {
		t.Errorf("Expected title %s, got %s", book.Title, retrievedBook.Title)
	}
}

func TestGetBookByIDNotFound(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/books/999", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", response.Code)
	}
}

func TestUpdateBook(t *testing.T) {
	clearDB()
	router := setupRouter()

	// First create a book
	book := Book{
		Title:  "Original Title",
		Author: "Original Author",
		ISBN:   "1234567890789",
		Year:   2023,
	}
	db.Create(&book)

	// Update the book
	updatedData := Book{
		Title: "Updated Title",
		Year:  2024,
	}

	jsonData, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest("PUT", "/api/v1/books/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.Code)
	}

	var updatedBook Book
	json.Unmarshal(response.Body.Bytes(), &updatedBook)

	if updatedBook.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", updatedBook.Title)
	}
	if updatedBook.Author != "Original Author" {
		t.Errorf("Expected author to remain 'Original Author', got %s", updatedBook.Author)
	}
}

func TestDeleteBook(t *testing.T) {
	clearDB()
	router := setupRouter()

	// First create a book
	book := Book{
		Title:  "Book to Delete",
		Author: "Author",
		ISBN:   "1234567890012",
		Year:   2023,
	}
	db.Create(&book)

	req, _ := http.NewRequest("DELETE", "/api/v1/books/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", response.Code)
	}

	// Verify book is deleted
	var deletedBook Book
	err := db.First(&deletedBook, 1).Error
	if err == nil {
		t.Error("Expected book to be deleted")
	}
}
