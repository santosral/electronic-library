package handler

import (
	"electronic-library/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type BookDetailHandler struct {
	BookRepo *repository.BookDetailRepository
}

func NewBookDetailHandler(repo *repository.BookDetailRepository) *BookDetailHandler {
	return &BookDetailHandler{BookRepo: repo}
}

func (h *BookDetailHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := NewErrorResponse("Invalid HTTP method", nil)
		json.NewEncoder(w).Encode(response)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := NewErrorResponse("Missing title query parameter", nil)
		json.NewEncoder(w).Encode(response)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "5"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "5"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 5
	}

	books, totalCount, err := h.BookRepo.SearchByTitle(r.Context(), title, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := NewErrorResponse("Error searching books", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	pagination := &Pagination{
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	w.WriteHeader(http.StatusOK)
	response := NewSuccessResponse(books, pagination)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := NewErrorResponse("Error encoding response:", err)
		json.NewEncoder(w).Encode(response)
	}
}
