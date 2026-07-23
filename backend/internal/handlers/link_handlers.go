package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go-links/internal/models"
	"go-links/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type LinkHandler struct {
	svc *service.LinkService
}

func NewLinkHandler(svc *service.LinkService) *LinkHandler {
	return &LinkHandler{svc: svc}
}

type errorResponse struct {
	Error     string `json:"error"`
	Field     string `json:"field,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message, RequestID: middleware.GetReqID(r.Context())})
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid JSON body")
		return
	}

	link, err := h.svc.CreateLink(req)
	if err != nil {
		var valErr *service.ValidationError
		switch {
		case errors.As(err, &valErr):
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error:     valErr.Message,
				Field:     valErr.Field,
				RequestID: middleware.GetReqID(r.Context()),
			})
		case errors.Is(err, service.ErrShortcut):
			writeError(w, r, http.StatusConflict, err.Error())
		default:
			writeError(w, r, http.StatusInternalServerError, "failed to create link")
		}
		return
	}

	writeJSON(w, http.StatusCreated, link)
}

func (h *LinkHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.svc.ListLinks()
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "failed to list links")
		return
	}
	writeJSON(w, http.StatusOK, links)
}

func (h *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid link id")
		return
	}

	link, err := h.svc.GetLink(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, "link not found")
			return
		}
		writeError(w, r, http.StatusInternalServerError, "failed to get link")
		return
	}

	writeJSON(w, http.StatusOK, link)
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid link id")
		return
	}

	if err := h.svc.DeleteLink(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, "link not found")
			return
		}
		writeError(w, r, http.StatusInternalServerError, "failed to delete link")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LinkHandler) RedirectShortcut(w http.ResponseWriter, r *http.Request) {
	shortcut := chi.URLParam(r, "shortcut")

	link, err := h.svc.ResolveShortcut(shortcut)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, "no link found for shortcut \""+shortcut+"\"")
			return
		}
		writeError(w, r, http.StatusInternalServerError, "failed to resolve shortcut")
		return
	}

	http.Redirect(w, r, link.Destination, http.StatusFound)
}

func parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.ParseInt(idStr, 10, 64)
}
