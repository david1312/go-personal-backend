package response

import (
	"bytes"
	"context"
	"libra-internal/pkg/crashy"
	"libra-internal/pkg/log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Status    string      `json:"status"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type ResponseCallback struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Code      string `json:"code"`
	Error     string `json:"error"`
}

type ErrorResponseCallback struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (resp Response) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
func (er ErrorResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func Yay(rw http.ResponseWriter, r *http.Request, data interface{}, code int) {
	p := Response{
		Status:    "ok",
		RequestID: middleware.GetReqID(r.Context()),
		Data:      data,
	}

	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, "ok"))
	if err := render.Render(rw, r, p); err != nil {
		Nay(rw, r, crashy.New(err, crashy.ErrCodeDataWrite, "unexpected error occurred while processing your request"), http.StatusInternalServerError)
	}
}

func Nay(rw http.ResponseWriter, r *http.Request, err *crashy.Error, code int) {
	p := ErrorResponse{
		Status:    "error",
		RequestID: middleware.GetReqID(r.Context()),
		Code:      string(err.Code),
		Error:     err.Message,
	}

	log.Errorf("error: %v, request-id: %s", err.Unwrap(), middleware.GetReqID(r.Context()))

	render.Status(r, code)
	if err := render.Render(rw, r, p); err != nil {
		http.Error(rw, "unexpected error occurred while processing your request", http.StatusInternalServerError)
	}
}

func renderHttpError(rw http.ResponseWriter, r *http.Request, code int, err string) {
	p := ErrorResponse{
		Status:    "error",
		RequestID: middleware.GetReqID(r.Context()),
		Code:      err,
		Error:     err,
	}
	render.Status(r, code)
	if err := render.Render(rw, r, p); err != nil {
		http.Error(rw, "unexpected error occurred while processing your request", http.StatusInternalServerError)
	}
}

func ExpiredAccess(rw http.ResponseWriter, r *http.Request) {
	renderHttpError(rw, r, http.StatusUnauthorized, crashy.ErrCodeExpired)
}

func ForbiddenAccess(rw http.ResponseWriter, r *http.Request) {
	renderHttpError(rw, r, http.StatusForbidden, crashy.ErrCodeForbidden)
}

func UnauthorizedAccess(rw http.ResponseWriter, r *http.Request) {
	renderHttpError(rw, r, http.StatusUnauthorized, crashy.ErrCodeUnauthorized)
}

func (resp ResponseCallback) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
func (er ErrorResponseCallback) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func File(w http.ResponseWriter, r *http.Request, filename string, bfr *bytes.Buffer, code int) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.WriteHeader(code)

	w.Write(bfr.Bytes())
}
