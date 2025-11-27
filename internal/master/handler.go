package master

import (
	"net/http"
	"strconv"

	"pemira-api/internal/http/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetFacultyPrograms(w http.ResponseWriter, r *http.Request) {
	options, err := h.service.GetFacultyProgramOptions(r.Context())
	if err != nil {
		response.InternalServerError(w, "DATABASE_ERROR", "Failed to fetch faculty programs")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"faculties": options,
	})
}

func (h *Handler) GetFaculties(w http.ResponseWriter, r *http.Request) {
	faculties, err := h.service.GetAllFaculties(r.Context())
	if err != nil {
		response.InternalServerError(w, "DATABASE_ERROR", "Failed to fetch faculties")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": faculties,
	})
}

func (h *Handler) GetStudyPrograms(w http.ResponseWriter, r *http.Request) {
	facultyIDStr := r.URL.Query().Get("faculty_id")
	
	var programs []StudyProgram
	var err error

	if facultyIDStr != "" {
		facultyID, parseErr := strconv.ParseInt(facultyIDStr, 10, 64)
		if parseErr != nil {
			response.BadRequest(w, "VALIDATION_ERROR", "Invalid faculty_id")
			return
		}
		programs, err = h.service.GetStudyProgramsByFaculty(r.Context(), facultyID)
	} else {
		programs, err = h.service.GetAllStudyPrograms(r.Context())
	}

	if err != nil {
		response.InternalServerError(w, "DATABASE_ERROR", "Failed to fetch study programs")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": programs,
	})
}

func (h *Handler) GetLecturerUnits(w http.ResponseWriter, r *http.Request) {
	units, err := h.service.GetAllLecturerUnits(r.Context())
	if err != nil {
		response.InternalServerError(w, "DATABASE_ERROR", "Failed to fetch lecturer units")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": units,
	})
}

func (h *Handler) GetLecturerPositions(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	
	var positions []LecturerPosition
	var err error

	if category != "" {
		positions, err = h.service.GetLecturerPositionsByCategory(r.Context(), category)
	} else {
		positions, err = h.service.GetAllLecturerPositions(r.Context())
	}

	if err != nil {
		response.InternalServerError(w, "DATABASE_ERROR", "Failed to fetch lecturer positions")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": positions,
	})
}

func (h *Handler) GetStaffUnits(w http.ResponseWriter, r *http.Request) {
	units, err := h.service.GetAllStaffUnits(r.Context())
	if err != nil {
		response.InternalServerError(w, "DATABASE_ERROR", "Failed to fetch staff units")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": units,
	})
}

func (h *Handler) GetStaffPositions(w http.ResponseWriter, r *http.Request) {
	positions, err := h.service.GetAllStaffPositions(r.Context())
	if err != nil {
		response.InternalServerError(w, "DATABASE_ERROR", "Failed to fetch staff positions")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": positions,
	})
}
