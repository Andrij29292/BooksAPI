package database

import (
	"database/sql"
	"fmt"
	book "hometask/Book"

	_ "github.com/lib/pq"
)



type Database struct {
	b *sql.DB
}

func (d *Database) Start() (err error) {
	d.b, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=12345 dbname=postgres sslmode=disable")

	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	err = d.b.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil

}

func (d *Database) GetById(id int) (book book.Book, err error) {
	err = d.b.QueryRow("select * from books where id = $1", id).Scan(
		&book.ID, &book.Name, &book.Author, &book.PagesCount, &book.RegisteredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return book, fmt.Errorf("id %d is not found", id)
		}
		return book, fmt.Errorf("error with finding the book")
	}

	return book, nil
}

func (d *Database) GetAll() (books []book.Book, err error) {
	rows, err := d.b.Query("select * from books")
	if err != nil {
		return nil, fmt.Errorf("error querying all books: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var book book.Book
		if err = rows.Scan(&book.ID, &book.Name, &book.Author, &book.PagesCount, &book.RegisteredAt); err != nil {
			return nil, fmt.Errorf("error scanning book row: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over book rows: %w", err)
	}

	return books, nil
}

func (d *Database) DeleteById(id int) (err error) {
	_, err = d.b.Exec("delete from books where id = $1", id)
	return err
}

func (d *Database) Insert(b book.Book) (err error) {
	tx, err := d.b.Begin()
	if err != nil {
		return fmt.Errorf("transaction initialization error: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("insert into books (name, author, pagesCount) values ($1, $2, $3)",
		b.Name, b.Author, b.PagesCount,
	)
	if err != nil {
		return fmt.Errorf("inserting error: %w", err)
	}
	_, err = tx.Exec("insert into logs (entity, action) values ('book', 'created')")
	if err != nil {
		return fmt.Errorf("logging error: %w", err)
	}

	return tx.Commit()
}

func (d *Database) UpdateById(id int, newBook book.Book) (err error) {
	_, err = d.b.Exec("update books set name=$1, author=$2, pagesCount=$3 where id = $4",
		newBook.Name, newBook.Author, newBook.PagesCount, id,
	)
	return err
}

func (d *Database) End() error {
	if d.b != nil {
		return d.b.Close()
	}
	return nil
}
