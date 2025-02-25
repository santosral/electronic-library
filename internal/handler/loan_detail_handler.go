package handler

import (
	"electronic-library/internal/model"
	"electronic-library/internal/request"
	"electronic-library/internal/service"
	"electronic-library/pkg"
	"electronic-library/pkg/response"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type LoanDetailHandler struct {
	ExtendService *service.ExtendService
}

func NewLoanDetailHandler(es *service.ExtendService) *LoanDetailHandler {
	return &LoanDetailHandler{ExtendService: es}
}

func (h *LoanDetailHandler) Extend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := pkg.GetLoggerFromContext(ctx).With(
		"handler", "Extend",
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

	var extendBookLoanDetail request.ExtendBookLoanDetail
	extendBookLoanDetail.ID = loanDetail.ID
	v := validator.New()
	if err := extendBookLoanDetail.Validate(v); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := response.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}

	ld, err := h.ExtendService.Call(r.Context(), loanDetail)
	if err != nil {
		if err.Code == http.StatusInternalServerError {
			logger.ErrorContext(ctx, "Failed to extend a book loan", "error", err)
		}

		w.WriteHeader(err.Code)
		resp := response.NewErrorResponse(err.Message)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			resp := response.NewErrorResponse(err.Error())
			json.NewEncoder(w).Encode(resp)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response.NewSuccessResponse(ld, nil)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewErrorResponse("Error encoding response:"))
	}
}
