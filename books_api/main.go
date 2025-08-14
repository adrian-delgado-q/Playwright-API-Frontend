package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Book model
type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" gorm:"not null"`
	Author string `json:"author" gorm:"not null"`
	ISBN   string `json:"isbn" gorm:"unique;not null"`
	Year   int    `json:"year"`
}

// Database instance
var db *gorm.DB

// Initialize database
func initDB() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "books.db"
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Book{})

	// Seed the database
	seedDatabase()
}

// Seed database with sample data
func seedDatabase() {
	var count int64
	db.Model(&Book{}).Count(&count)

	if count == 0 {
		books := []Book{
			{Title: "The Go Programming Language", Author: "Alan Donovan", ISBN: "9780134190440", Year: 2015},
			{Title: "Clean Code", Author: "Robert C. Martin", ISBN: "9780132350884", Year: 2008},
			{Title: "The Pragmatic Programmer", Author: "David Thomas", ISBN: "9780201616224", Year: 1999},
			{Title: "Design Patterns", Author: "Gang of Four", ISBN: "9780201633612", Year: 1994},
			{Title: "Refactoring", Author: "Martin Fowler", ISBN: "9780201485677", Year: 1999},
		}

		for _, book := range books {
			db.Create(&book)
		}
		fmt.Println("Database seeded with sample books")
	}
}

// API Handlers

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var books []Book
	db.Find(&books)
	json.NewEncoder(w).Encode(books)
}

// Get book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// Create new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" || book.ISBN == "" {
		http.Error(w, "Title, Author, and ISBN are required", http.StatusBadRequest)
		return
	}

	if err := db.Create(&book).Error; err != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update fields
	if updatedBook.Title != "" {
		book.Title = updatedBook.Title
	}
	if updatedBook.Author != "" {
		book.Author = updatedBook.Author
	}
	if updatedBook.ISBN != "" {
		book.ISBN = updatedBook.ISBN
	}
	if updatedBook.Year != 0 {
		book.Year = updatedBook.Year
	}

	db.Save(&book)
	json.NewEncoder(w).Encode(book)
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	if err := db.First(&book, id).Error; err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	db.Delete(&book)
	w.WriteHeader(http.StatusNoContent)
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize database
	initDB()

	// Setup routes
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/books", getBooks).Methods("GET")
	api.HandleFunc("/books", createBook).Methods("POST")
	api.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/books/{id}", getBook).Methods("GET")
	api.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	api.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	api.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	fmt.Println("Books API server starting on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
