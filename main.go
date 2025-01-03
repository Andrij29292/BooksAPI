package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	book "hometask/Book"
	database "hometask/DataBase"
)

var db = &database.Database{}

func ParseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.URL.Query().Get("id"))
}

func ParseBook(r *http.Request) (book.Book, error) {
	var b book.Book

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return book.Book{}, err
	}

	if err = json.Unmarshal(body, &b); err != nil {
		return book.Book{}, err
	}
	return b, nil

}

func Logging(r *http.Request) {
	log.Printf("[%s] - %s\n", r.Method, time.Now().Format("2006-01-02 15:04:05"))
}

func Book(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetBook(w, r)
	case http.MethodPut:
		PutBook(w, r)
	case http.MethodPost:
		PostBook(w, r)
	case http.MethodDelete:
		DeleteBook(w, r)
	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	Logging(r)
}

func Books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetBooks(w, r)
	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	Logging(r)
}

func main() {
	if err := db.Start(); err != nil {
		panic(err)
	}
	defer db.End()

	http.HandleFunc("/book", Book)
	http.HandleFunc("/books", Books)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	// Отримання параметра id із запиту
	id, err := ParseID(r)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid 'id' parameter. It must be a positive integer.", http.StatusBadRequest)
		return
	}

	book, err := db.GetById(id)
	if err != nil {
		http.Error(w, "Error retrieving the book from the database: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "Failed to encode response to JSON", http.StatusInternalServerError)
	}

}

func PutBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil || id < 0 {
		http.Error(w, "Invalid 'id' parameter. It must be a positive integer.", http.StatusBadRequest)
		return
	}

	b, err := ParseBook(r)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err = db.UpdateById(id, b); err != nil {
		http.Error(w, "Database cannot update book", http.StatusInternalServerError)
		return
	}
}

func PostBook(w http.ResponseWriter, r *http.Request) {
	b, err := ParseBook(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = db.Insert(b); err != nil {
		http.Error(w, "Database cannot post book", http.StatusInternalServerError)
	}

}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil || id < 0 {
		http.Error(w, "Invalid 'id' parameter. It must be a positive integer.", http.StatusBadRequest)
		return
	}

	err = db.DeleteById(id)
	if err != nil {
		http.Error(w, "Error delete the book from the database: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := db.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Failed to encode response to JSON", http.StatusInternalServerError)
	}
}
