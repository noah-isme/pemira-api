package adminuser

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"pemira-api/internal/http/response"
	"pemira-api/internal/shared"
	"pemira-api/internal/shared/constants"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List admin users: GET /admin/users
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := ListFilter{
		Search: strings.TrimSpace(q.Get("search")),
		Role:   strings.ToUpper(strings.TrimSpace(q.Get("role"))),
	}
	if a := q.Get("active"); a != "" {
		val := strings.ToLower(a) == "true"
		filter.Active = &val
	}

	page := parseIntDefault(q.Get("page"), 1)
	limit := parseIntDefault(q.Get("limit"), 50)

	items, meta, err := h.svc.List(r.Context(), filter, page, limit)
	if err != nil {
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal mengambil daftar admin")
		return
	}

	response.Success(w, http.StatusOK, map[string]interface{}{
		"items":       items,
		"page":        meta.CurrentPage,
		"limit":       meta.PerPage,
		"total_items": meta.Total,
		"total_pages": meta.TotalPages,
	})
}

// Create admin user: POST /admin/users
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid")
		return
	}

	user, err := h.svc.Create(r.Context(), req)
	if err != nil {
		switch err {
		case shared.ErrBadRequest:
			response.BadRequest(w, "VALIDATION_ERROR", "Data tidak valid atau role tidak diperbolehkan")
			return
		case shared.ErrDuplicateEntry:
			response.Conflict(w, "DUPLICATE", "Username atau email sudah digunakan")
			return
		default:
			response.InternalServerError(w, "INTERNAL_ERROR", "Gagal membuat admin")
			return
		}
	}

	response.Success(w, http.StatusCreated, user)
}

// Detail: GET /admin/users/{userID}
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, chi.URLParam(r, "userID"))
	if !ok {
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if err == shared.ErrNotFound {
			response.NotFound(w, "NOT_FOUND", "User tidak ditemukan")
			return
		}
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal mengambil user")
		return
	}

	response.Success(w, http.StatusOK, user)
}

// Update: PATCH /admin/users/{userID}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, chi.URLParam(r, "userID"))
	if !ok {
		return
	}

	var req UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid")
		return
	}

	// Normalize role string if provided
	if req.Role != nil {
		rv := constants.Role(strings.ToUpper(string(*req.Role)))
		req.Role = &rv
	}

	user, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		switch err {
		case shared.ErrBadRequest:
			response.BadRequest(w, "VALIDATION_ERROR", "Role tidak valid")
			return
		case shared.ErrDuplicateEntry:
			response.Conflict(w, "DUPLICATE", "Username atau email sudah digunakan")
			return
		case shared.ErrNotFound:
			response.NotFound(w, "NOT_FOUND", "User tidak ditemukan")
			return
		default:
			response.InternalServerError(w, "INTERNAL_ERROR", "Gagal memperbarui user")
			return
		}
	}

	response.Success(w, http.StatusOK, user)
}

// Reset password: POST /admin/users/{userID}/reset-password
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, chi.URLParam(r, "userID"))
	if !ok {
		return
	}

	var req ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid")
		return
	}

	if err := h.svc.ResetPassword(r.Context(), id, req.NewPassword); err != nil {
		switch err {
		case shared.ErrBadRequest:
			response.BadRequest(w, "VALIDATION_ERROR", "Password baru tidak valid")
			return
		case shared.ErrNotFound:
			response.NotFound(w, "NOT_FOUND", "User tidak ditemukan")
			return
		default:
			response.InternalServerError(w, "INTERNAL_ERROR", "Gagal reset password")
			return
		}
	}

	response.Success(w, http.StatusOK, map[string]bool{"success": true})
}

// Activate/Deactivate
func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	h.setActive(w, r, true)
}

func (h *Handler) Deactivate(w http.ResponseWriter, r *http.Request) {
	h.setActive(w, r, false)
}

func (h *Handler) setActive(w http.ResponseWriter, r *http.Request, active bool) {
	id, ok := parseID(w, chi.URLParam(r, "userID"))
	if !ok {
		return
	}
	user, err := h.svc.SetActive(r.Context(), id, active)
	if err != nil {
		if err == shared.ErrNotFound {
			response.NotFound(w, "NOT_FOUND", "User tidak ditemukan")
			return
		}
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal memperbarui status user")
		return
	}
	response.Success(w, http.StatusOK, user)
}

// Delete: DELETE /admin/users/{userID}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, chi.URLParam(r, "userID"))
	if !ok {
		return
	}
	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		if err == shared.ErrNotFound {
			response.NotFound(w, "NOT_FOUND", "User tidak ditemukan")
			return
		}
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal menghapus user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, raw string) (int64, bool) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(w, "VALIDATION_ERROR", "ID tidak valid")
		return 0, false
	}
	return id, true
}

func parseIntDefault(raw string, def int) int {
	if raw == "" {
		return def
	}
	if v, err := strconv.Atoi(raw); err == nil && v > 0 {
		return v
	}
	return def
}
