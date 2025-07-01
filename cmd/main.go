package main

import (
	"log"
	"net/http"

	"github.com/aziz-shoko/goblog/internal/handler"
	"github.com/aziz-shoko/goblog/internal/service"
	"github.com/aziz-shoko/goblog/internal/store"
)

func main() {
	postStore := store.NewInMemoryStore()
	postService := service.NewPostService(postStore)
	postHandler := handler.NewPostHandler(postService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", handler.LoggingMiddleware(postHandler.CreatePost))
	mux.HandleFunc("GET /post/{id}", handler.LoggingMiddleware(postHandler.GetPostByID))
	mux.HandleFunc("GET /posts", handler.LoggingMiddleware(postHandler.GetPostsAll))
	mux.HandleFunc("DELETE /posts", handler.LoggingMiddleware(postHandler.DeleteAllPosts))

	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
