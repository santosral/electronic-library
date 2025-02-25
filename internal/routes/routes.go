package routes

import (
	"electronic-library/internal/handler"
	"electronic-library/pkg/middleware"
	"net/http"
)

func SetupBookDetailRoutes(mux *http.ServeMux, bookDetailHandler *handler.BookDetailHandler) {
	mux.HandleFunc("/book-details/search", middleware.Chain(bookDetailHandler.SearchBooks, middleware.Logging(), middleware.Method("GET"), middleware.Header()))
	mux.HandleFunc("/borrow", middleware.Chain(bookDetailHandler.BorrowBook, middleware.Logging(), middleware.Method("POST"), middleware.Header()))
	mux.HandleFunc("/return", middleware.Chain(bookDetailHandler.ReturnBook, middleware.Logging(), middleware.Method("POST"), middleware.Header()))
}

func SetupLoanDetailRoutes(mux *http.ServeMux, loanDetailHandler *handler.LoanDetailHandler) {
	mux.HandleFunc("/extend", middleware.Chain(loanDetailHandler.Extend, middleware.Logging(), middleware.Method("POST"), middleware.Header()))
}

func SetupRoutes(mux *http.ServeMux, bookDetailHandler *handler.BookDetailHandler, loanDetailHandler *handler.LoanDetailHandler) {
	SetupBookDetailRoutes(mux, bookDetailHandler)
	SetupLoanDetailRoutes(mux, loanDetailHandler)
}
