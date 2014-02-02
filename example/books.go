package main

import (
	"../../sleepy"
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var db, _ = sql.Open("sqlite3", "./development.sqlite3")
var dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

type Book struct {
	Id        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	BooksURL  string    `db:"books_url" json:"books_url"`
	Have      bool      `db:"have" json:"have"`
	Read      bool      `db:"read" json:"read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Books []Book

func (books *Books) NewModelInstance() interface{} {
	var book Book
	return &book
}

func (books *Books) PathLength() int {
	return len("/books/")
}

// GET /books/
func (books Books) Get(book interface{}, id string) (int, interface{}) {
	data, err := dbmap.Select(Book{}, "select * from books")
	if err != nil {
		fmt.Print(err)
	}

	return 200, data
}

// POST /books/
func (books Books) Post(book interface{}, id string) (int, interface{}) {
	err := dbmap.Insert(book)
	if err != nil {
		fmt.Println(err)
	}

	return 200, book
}

// DELETE /books/:id
func (books Books) Delete(book interface{}, id string) (int, interface{}) {
	delete_book, err := dbmap.Get(Book{}, id)
	if err != nil {
		fmt.Print(err)
	}

	dbmap.Delete(delete_book)

	return 200, delete_book
}

// PUT /books/:id
func (books Books) Put(book interface{}, id string) (int, interface{}) {
	dbmap.Update(book)

	return 200, book
}

func main() {
	books := new(Books)

	var api = sleepy.NewAPI()
	// the restful api is
	// GET		/books/
	// POST		/books/
	// DELETE	/books/:id
	// PUT		/books/:id
	api.AddResource(books, "/books/")
	api.Start(4000)
}
