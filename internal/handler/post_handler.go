package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aziz-shoko/goblog/internal/service"
)

type CreatePostRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type CreatePostResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type PostHandler struct {
	Service *service.PostServiceRepository
}

func NewPostHandler(service *service.PostServiceRepository) *PostHandler {
	return &PostHandler{
		Service: service,
	}
}

// CreatePost handles POST /posts requests
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Call service
	post, err := h.Service.CreatePost(req.Name, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build response
	response := CreatePostResponse{
		ID:        post.ID,
		Name:      post.Name,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetPostByID handles the getting post part
func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/post/")

	// Call service
	post, err := h.Service.Store.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Build response
	response := CreatePostResponse{
		ID:        post.ID,
		Name:      post.Name,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetPostsAll returns all the posts
func (h *PostHandler) GetPostsAll(w http.ResponseWriter, r *http.Request) {
	// call service
	posts, err := h.Service.Store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Build Response
	response := []CreatePostResponse{}
	for _, post := range posts {
		response = append(response, CreatePostResponse{
			ID:        post.ID,
			Name:      post.Name,
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteAllPosts handler
func (h *PostHandler) DeleteAllPosts(w http.ResponseWriter, r *http.Request) {
	// call service 
	err := h.Service.Store.DeleteAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}