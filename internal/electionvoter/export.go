package electionvoter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"

	"pemira-api/internal/http/response"
)

// ExportToExcel handles GET /admin/elections/{electionID}/voters/export
func (h *Handler) ExportToExcel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	electionIDStr := chi.URLParam(r, "electionID")
	electionID, err := strconv.ParseInt(electionIDStr, 10, 64)
	if err != nil || electionID <= 0 {
		response.BadRequest(w, "INVALID_REQUEST", "electionID tidak valid.")
		return
	}

	// Get all voters (no pagination for export)
	filter := ListFilter{
		Search:           r.URL.Query().Get("search"),
		VoterType:        r.URL.Query().Get("voter_type"),
		Status:           r.URL.Query().Get("status"),
		VotingMethod:     r.URL.Query().Get("voting_method"),
		FacultyCode:      r.URL.Query().Get("faculty_code"),
		StudyProgramCode: r.URL.Query().Get("study_program_code"),
	}

	// Use large limit to get all data
	page := 1
	limit := 10000
	voters, _, err := h.svc.List(ctx, electionID, filter, page, limit)
	if err != nil {
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal mengambil data DPT.")
		return
	}

	// Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "DPT"
	f.SetSheetName("Sheet1", sheetName)

	// Set headers
	headers := []string{
		"No", "NIM", "Nama", "Tipe", "Fakultas", "Program Studi",
		"Angkatan", "Semester", "Status Akademik", "Email",
		"Status DPT", "Metode Voting", "Sudah Memilih", "Waktu Memilih",
		"Login Terakhir", "Blacklist", "Tanda Tangan",
	}

	// Style for header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 5, "B": 15, "C": 30, "D": 12, "E": 20, "F": 25,
		"G": 10, "H": 10, "I": 15, "J": 30,
		"K": 12, "L": 12, "M": 15, "N": 20,
		"O": 20, "P": 10, "Q": 25,
	}
	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// HTTP client for downloading images
	httpClient := &http.Client{Timeout: 5 * time.Second}

	// Fill data
	for i, v := range voters {
		row := i + 2

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), v.NIM)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), v.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), v.VoterType)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), derefStr(v.FacultyName))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), derefStr(v.StudyProgramName))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), derefInt(v.CohortYear))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), derefInt(v.Semester))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), derefStr(v.AcademicStatus))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), derefStr(v.Email))
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), v.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), v.VotingMethod)

		hasVoted := "Belum"
		if v.HasVoted != nil && *v.HasVoted {
			hasVoted = "Sudah"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), hasVoted)

		if v.VotedAt != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), v.VotedAt.In(time.FixedZone("WIB", 7*3600)).Format("02/01/2006 15:04"))
		}

		if v.LastLoginAt != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), v.LastLoginAt.In(time.FixedZone("WIB", 7*3600)).Format("02/01/2006 15:04"))
		}

		// Blacklist status
		blacklist := "Tidak"
		if v.IsBlacklisted {
			blacklist = "Ya"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), blacklist)

		// Embed signature image if URL exists
		if v.DigitalSignatureURL != nil && *v.DigitalSignatureURL != "" {
			imgData, ext := downloadImage(httpClient, *v.DigitalSignatureURL)
			if imgData != nil {
				// Set row height for image (approx 50 pixels = 37.5 points)
				f.SetRowHeight(sheetName, row, 50)

				cell := fmt.Sprintf("Q%d", row)
				if err := f.AddPictureFromBytes(sheetName, cell, &excelize.Picture{
					Extension: ext,
					File:      imgData,
					Format: &excelize.GraphicOptions{
						AutoFit: true,
						ScaleX:  0.3,
						ScaleY:  0.3,
					},
				}); err != nil {
					// Fallback to URL if image embedding fails
					f.SetCellValue(sheetName, cell, *v.DigitalSignatureURL)
				}
			} else {
				// Fallback to URL if download fails
				f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), *v.DigitalSignatureURL)
			}
		}
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("DPT_Election_%d_%s.xlsx", electionID, time.Now().Format("20060102_150405"))

	// Set response headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Write to response
	if err := f.Write(w); err != nil {
		response.InternalServerError(w, "INTERNAL_ERROR", "Gagal menulis file Excel.")
		return
	}
}

// downloadImage downloads image from URL and returns bytes and extension
func downloadImage(client *http.Client, url string) ([]byte, string) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ""
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ""
	}

	// Detect extension from content type or URL
	ext := ".png"
	contentType := resp.Header.Get("Content-Type")
	if contentType == "image/jpeg" || contentType == "image/jpg" {
		ext = ".jpg"
	} else if bytes.HasPrefix(data, []byte("\xff\xd8\xff")) {
		ext = ".jpg"
	}

	return data, ext
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
