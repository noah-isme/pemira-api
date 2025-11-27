package electionvoter

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"pemira-api/internal/auth"
	"pemira-api/internal/http/response"
	"pemira-api/internal/shared"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// AdminLookup handles GET /admin/elections/{electionID}/voters/lookup?nim=...
func (h *Handler) AdminLookup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	electionID, ok := parseID(w, chi.URLParam(r, "electionID"))
	if !ok {
		return
	}
	nim := strings.TrimSpace(r.URL.Query().Get("nim"))
	if nim == "" {
		response.BadRequest(w, "VALIDATION_ERROR", "parameter nim wajib diisi")
		return
	}

	res, err := h.svc.LookupByNIM(ctx, electionID, nim)
	if err != nil {
		if err == shared.ErrNotFound {
			response.NotFound(w, "NOT_FOUND", "Pemilih dengan NIM tersebut tidak ditemukan")
			return
		}
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal mengambil data pemilih")
		return
	}

	response.Success(w, http.StatusOK, res)
}

// AdminUpsert handles POST /admin/elections/{electionID}/voters
func (h *Handler) AdminUpsert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	electionID, ok := parseID(w, chi.URLParam(r, "electionID"))
	if !ok {
		return
	}

	var req UpsertAndEnrollInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid")
		return
	}

	res, err := h.svc.UpsertAndEnroll(ctx, electionID, req)
	if err != nil {
		switch err {
		case shared.ErrBadRequest:
			response.BadRequest(w, "VALIDATION_ERROR", "Data wajib diisi atau tidak valid")
			return
		case shared.ErrDuplicateEntry:
			response.Conflict(w, "DUPLICATE", "NIM sudah terdaftar di pemilu ini")
			return
		default:
			response.InternalServerError(w, "INTERNAL_ERROR", "Gagal menyimpan data pemilih")
			return
		}
	}

	response.Success(w, http.StatusOK, res)
}

// AdminList handles GET /admin/elections/{electionID}/voters
func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	electionID, ok := parseID(w, chi.URLParam(r, "electionID"))
	if !ok {
		return
	}

	q := r.URL.Query()
	filter := ListFilter{
		Search:           strings.TrimSpace(q.Get("search")),
		VoterType:        q.Get("voter_type"),
		Status:           q.Get("status"),
		VotingMethod:     q.Get("voting_method"),
		FacultyCode:      q.Get("faculty_code"),
		StudyProgramCode: q.Get("study_program_code"),
	}

	if cy := q.Get("cohort_year"); cy != "" {
		if v, err := strconv.Atoi(cy); err == nil {
			filter.CohortYear = &v
		}
	}
	if tps := q.Get("tps_id"); tps != "" {
		if v, err := strconv.ParseInt(tps, 10, 64); err == nil {
			filter.TPSID = &v
		}
	}

	filter, err := ValidateFilter(filter)
	if err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Filter tidak valid")
		return
	}

	page := parseIntDefault(q.Get("page"), 1)
	limit := parseIntDefault(q.Get("limit"), 50)

	items, meta, err := h.svc.List(ctx, electionID, filter, page, limit)
	if err != nil {
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal mengambil daftar pemilih")
		return
	}

	resp := map[string]interface{}{
		"items":       items,
		"page":        meta.CurrentPage,
		"limit":       meta.PerPage,
		"total_items": meta.Total,
		"total_pages": meta.TotalPages,
	}
	response.Success(w, http.StatusOK, resp)
}

// AdminPatch handles PATCH /admin/elections/{electionID}/voters/{enrollmentID}
func (h *Handler) AdminPatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	electionID, ok := parseID(w, chi.URLParam(r, "electionID"))
	if !ok {
		return
	}
	enrollmentID, ok := parseID(w, chi.URLParam(r, "voterID"))
	if !ok {
		return
	}

	var req UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid")
		return
	}

	ev, err := h.svc.UpdateEnrollment(ctx, electionID, enrollmentID, req)
	if err != nil {
		switch err {
		case shared.ErrBadRequest:
			response.BadRequest(w, "VALIDATION_ERROR", "status / voting_method tidak valid")
			return
		case shared.ErrNotFound:
			response.NotFound(w, "NOT_FOUND", "Data pemilih tidak ditemukan")
			return
		default:
			response.InternalServerError(w, "INTERNAL_ERROR", "Gagal memperbarui data")
			return
		}
	}

	response.Success(w, http.StatusOK, ev)
}

// VoterSelfRegister handles POST /voters/me/elections/{electionID}/register
func (h *Handler) VoterSelfRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := auth.FromContext(ctx)
	if !ok || authUser.VoterID == nil {
		response.Forbidden(w, "FORBIDDEN", "Akses tidak diizinkan")
		return
	}

	electionID, ok := parseID(w, chi.URLParam(r, "electionID"))
	if !ok {
		return
	}

	var req SelfRegisterInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "VALIDATION_ERROR", "Body tidak valid")
		return
	}

	ev, err := h.svc.SelfRegister(ctx, electionID, *authUser.VoterID, req)
	if err != nil {
		switch err {
		case shared.ErrDuplicateEntry:
			response.Conflict(w, "DUPLICATE", "NIM sudah terdaftar di pemilu ini")
			return
		case shared.ErrBadRequest:
			response.BadRequest(w, "VALIDATION_ERROR", "Metode voting tidak valid")
			return
		case shared.ErrNotFound:
			response.NotFound(w, "NOT_FOUND", "Voter tidak ditemukan")
			return
		default:
			response.InternalServerError(w, "INTERNAL_ERROR", "Gagal mendaftarkan pemilih")
			return
		}
	}

	response.Success(w, http.StatusOK, ev)
}

// VoterStatus handles GET /voters/me/elections/{electionID}/status
func (h *Handler) VoterStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := auth.FromContext(ctx)
	if !ok || authUser.VoterID == nil {
		response.Forbidden(w, "FORBIDDEN", "Akses tidak diizinkan")
		return
	}

	electionID, ok := parseID(w, chi.URLParam(r, "electionID"))
	if !ok {
		return
	}

	ev, err := h.svc.GetStatus(ctx, electionID, *authUser.VoterID)
	if err != nil {
		if err == shared.ErrNotFound {
			response.NotFound(w, "NOT_FOUND", "Belum terdaftar di pemilu ini")
			return
		}
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal mengambil status pemilih")
		return
	}

	response.Success(w, http.StatusOK, ev)
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
