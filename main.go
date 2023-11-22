package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rs/cors"

	"BookReview/models"
)

type BookInfo struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}

type CommentInfo struct {
	ISBN    string `json:"isbn"`
	Comment string `json:"comment"`
}

func main() {

	Mux := http.NewServeMux()

	Demo_CORS := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //all
		AllowedMethods:   []string{http.MethodPost, http.MethodGet},
		AllowedHeaders:   []string{"*"}, //all
		AllowCredentials: false,         //none
	})

	Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// データをJSON形式でフロントエンドに返す
		bookInfo := models.SelectDb()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bookInfo); err != nil {
			log.Fatal(err)
		}
	})

	Mux.HandleFunc("/book/post", func(w http.ResponseWriter, r *http.Request) {
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
		models.InsertBook(title, isbn)

	})

	Mux.HandleFunc("/comment/post", func(w http.ResponseWriter, r *http.Request) {
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
		models.InsertComment(isbn, comment)

	})

	Demo_Handler := Demo_CORS.Handler(Mux)

	// http.HandleFunc("/index/", indexHandler)
	log.Fatal(http.ListenAndServe(":3000", Demo_Handler))

}
