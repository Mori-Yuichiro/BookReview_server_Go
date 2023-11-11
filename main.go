package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Book struct {
	// ID    int    `json:"id"`
	Title   string `json:"title"`
	ISBN    string `json:"isbn"`
	Comment string `json:"comment"`
}

type BookInfo struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}

type CommentInfo struct {
	ISBN    string `json:"isbn"`
	Comment string `json:"comment"`
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

func selectDb() []Book {
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

func insertBook(title, isbn string) {
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

func insertComment(isbn, comment string) {
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

func main() {
	bookInfo := selectDb()

	Demo_Mux := http.NewServeMux()

	Demo_CORS := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //all
		AllowedMethods:   []string{http.MethodPost, http.MethodGet},
		AllowedHeaders:   []string{"*"}, //all
		AllowCredentials: false,         //none
	})

	Demo_Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// データをJSON形式でフロントエンドに返す
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bookInfo); err != nil {
			log.Fatal(err)
		}
	})

	Demo_Mux.HandleFunc("/book/post", func(w http.ResponseWriter, r *http.Request) {
		len := r.ContentLength
		body := make([]byte, len) // Content-Length と同じサイズの byte 配列を用意
		r.Body.Read(body)         // byte 配列にリクエストボディを読み込む

		var book BookInfo
		err := json.Unmarshal(body, &book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		title := book.Title
		isbn := book.ISBN
		insertBook(title, isbn)

	})

	Demo_Mux.HandleFunc("/comment/post", func(w http.ResponseWriter, r *http.Request) {
		len := r.ContentLength
		body := make([]byte, len) // Content-Length と同じサイズの byte 配列を用意
		r.Body.Read(body)         // byte 配列にリクエストボディを読み込む

		var commentInfo CommentInfo
		err := json.Unmarshal(body, &commentInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		isbn := commentInfo.ISBN
		comment := commentInfo.Comment
		insertComment(isbn, comment)

	})

	Demo_Handler := Demo_CORS.Handler(Demo_Mux)

	// http.HandleFunc("/index/", indexHandler)
	log.Fatal(http.ListenAndServe(":3000", Demo_Handler))

}
