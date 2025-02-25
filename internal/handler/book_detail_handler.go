package handler

import (
	"electronic-library/internal/model"
	"electronic-library/internal/repository"
	"electronic-library/internal/request"
	"electronic-library/internal/service"
	"electronic-library/pkg"
	"electronic-library/pkg/response"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type BookDetailHandler struct {
	BookRepo      repository.BookDetailRepository
	BorrowService *service.BorrowService
	ReturnService *service.ReturnService
}

func NewBookDetailHandler(repo repository.BookDetailRepository, bs *service.BorrowService, rs *service.ReturnService) *BookDetailHandler {
	return &BookDetailHandler{
		BookRepo:      repo,
		BorrowService: bs,
		ReturnService: rs,
	}
}

func (h *BookDetailHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := pkg.GetLoggerFromContext(ctx).With(
		"handler", "SearchBooks",
		"method", r.Method,
		"path", r.URL.Path,
		"client_ip", r.RemoteAddr,
	)

	title := r.URL.Query().Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse("Missing title query parameter")
		json.NewEncoder(w).Encode(resp)
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

	books, err := h.BookRepo.SearchByTitle(r.Context(), title, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse("Internal server error")
		json.NewEncoder(w).Encode(resp)
		logger.ErrorContext(ctx, "Failed to search books")
		return
	}

	pagination := &response.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(books, pagination)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server error")
		logger.ErrorContext(ctx, "Failed to encode JSON", "handler_method", "SearchBooks")
	}
}

func (h *BookDetailHandler) BorrowBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := pkg.GetLoggerFromContext(ctx).With(
		"handler", "BorrowBook",
		"method", r.Method,
		"path", r.URL.Path,
		"client_ip", r.RemoteAddr,
	)

	var loanDetail *model.LoanDetail
	if err := json.NewDecoder(r.Body).Decode(&loanDetail); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewErrorResponse("Invalid JSON format"))
		return
	}

	var borrowBookRequest request.BorrowBook
	borrowBookRequest.BookID = loanDetail.BookID
	borrowBookRequest.NameOfBorrower = loanDetail.NameOfBorrower
	v := validator.New()
	if err := borrowBookRequest.Validate(v); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := response.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}

	loanDetail, err := h.BorrowService.Call(r.Context(), loanDetail)
	if err != nil {
		if err.Code == http.StatusInternalServerError {
			logger.ErrorContext(ctx, "Failed to borrow a book", "error", err)
		}

		w.WriteHeader(err.Code)
		resp := response.NewErrorResponse(err.Message)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp := response.NewErrorResponse(err.Error())
			json.NewEncoder(w).Encode(resp)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response.NewSuccessResponse(loanDetail, nil)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewErrorResponse("internal server error"))
		logger.ErrorContext(ctx, "Failed to encode JSON")
	}

}

func (h *BookDetailHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := pkg.GetLoggerFromContext(ctx).With(
		"method", r.Method,
		"path", r.URL.Path,
		"client_ip", r.RemoteAddr,
	)

	var loanDetail *model.LoanDetail
	if err := json.NewDecoder(r.Body).Decode(&loanDetail); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewErrorResponse("Invalid JSON format"))
		return
	}

	var returnBookRequest request.ReturnBook
	returnBookRequest.ID = loanDetail.ID
	v := validator.New()
	if err := returnBookRequest.Validate(v); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := response.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}

	loanDetail, err := h.ReturnService.Call(r.Context(), loanDetail)
	if err != nil {
		if err.Code == http.StatusInternalServerError {
			logger.ErrorContext(ctx, "Failed to return book", "error", err)
		}

		w.WriteHeader(err.Code)
		resp := response.NewErrorResponse(err.Message)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp := response.NewErrorResponse(err.Error())
			json.NewEncoder(w).Encode(resp)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response.NewSuccessResponse(loanDetail, nil)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(resp)
		logger.ErrorContext(ctx, "Failed to encode JSON", "error", err)
	}
}
