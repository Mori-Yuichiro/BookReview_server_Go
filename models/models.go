package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Book struct {
	// ID    int    `json:"id"`
	Title   string `json:"title"`
	ISBN    string `json:"isbn"`
	Comment string `json:"comments"`
}

func connectDb() *sql.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// エラーハンドリング
		fmt.Println(err)
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	DBName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	c := mysql.Config{
		DBName:    DBName,
		User:      user,
		Passwd:    password,
		Addr:      "localhost:" + port,
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}
	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func SelectDb() []Book {
	db := connectDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Println(err)
		return nil
	} else {
		fmt.Println("接続完了")
	}

	// SQLの実行
	rows, err := db.Query("SELECT title, books.isbn, comments FROM books inner join comments on comments.isbn = books.isbn")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var book Book
		// err := rows.Scan(&book.ID, &book.Title, &book.ISBN)
		err := rows.Scan(&book.Title, &book.ISBN, &book.Comment)
		if err != nil {
			fmt.Println(err)
		}
		books = append(books, book)
	}
	fmt.Println("*********************************")
	return books
}

func InsertBook(title, isbn string) {
	db := connectDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Println(err)
		return
	} else {
		fmt.Println("接続完了")
	}

	// SQLの実行
	in, err := db.Prepare("insert into books (title, isbn) select * from (select ? as title, ? as isbn) as tmp where not exists (select * from books where title=? and isbn=?)")
	if err != nil {
		fmt.Println(err)
	}
	res, err := in.Exec(title, isbn, title, isbn)
	if err != nil {
		log.Fatal(err)
	}

	// 結果の取得
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(lastInsertID)
}

func InsertComment(isbn, comment string) {
	db := connectDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Println(err)
		return
	} else {
		fmt.Println("接続完了")
	}

	// SQLの実行
	in, err := db.Prepare("insert into comments (isbn, comments) values (?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	res, err := in.Exec(isbn, comment)
	if err != nil {
		log.Fatal(err)
	}

	// 結果の取得
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(lastInsertID)
}
