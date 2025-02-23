package handler

import (
	"electronic-library/internal/model"
	"electronic-library/internal/service"
	"encoding/json"
	"net/http"
)

type LoanDetailHandler struct {
	ExtendService service.ExtendService
}

func NewLoanDetailHandler(es service.ExtendService) *LoanDetailHandler {
	return &LoanDetailHandler{ExtendService: es}
}

func (h *LoanDetailHandler) Extend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(NewErrorResponse("Invalid HTTP method", nil))
		return
	}

	var loanDetail *model.LoanDetail
	if err := json.NewDecoder(r.Body).Decode(&loanDetail); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse("Invalid JSON format", err))
		return
	}

	loanDetail, err := h.ExtendService.Call(r.Context(), loanDetail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := NewErrorResponse("Error extending the loan", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(NewSuccessResponse(loanDetail, nil)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewErrorResponse("Error encoding response:", err))
	}
}
