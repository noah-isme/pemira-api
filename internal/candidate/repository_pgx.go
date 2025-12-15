package candidate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	storage_go "github.com/supabase-community/storage-go"
)

// PgCandidateRepository implements CandidateRepository using pgxpool
type PgCandidateRepository struct {
	db *pgxpool.Pool
}

// NewPgCandidateRepository creates a new PostgreSQL candidate repository
func NewPgCandidateRepository(db *pgxpool.Pool) *PgCandidateRepository {
	return &PgCandidateRepository{db: db}
}

const qListCandidatesBase = `
SELECT
id,
election_id,
number,
name,
photo_url,
photo_media_id::text AS photo_media_id,
short_bio,
long_bio,
tagline,
faculty_name,
study_program_name,
cohort_year,
vision,
missions,
main_programs,
media,
social_links,
status,
created_at,
updated_at
FROM candidates
WHERE election_id = $1
`

// Compatibility query when candidates.photo_media_id column does not exist (pre-migration 016).
const qListCandidatesBaseNoPhotoMedia = `
SELECT
id,
election_id,
number,
name,
photo_url,
NULL::text AS photo_media_id,
short_bio,
long_bio,
tagline,
faculty_name,
study_program_name,
cohort_year,
vision,
missions,
main_programs,
media,
social_links,
status,
created_at,
updated_at
FROM candidates
WHERE election_id = $1
`

const qCountCandidatesBase = `
SELECT COUNT(*) FROM candidates WHERE election_id = $1
`

// ListByElection returns candidates for an election with filters and pagination
func (r *PgCandidateRepository) ListByElection(
	ctx context.Context,
	electionID int64,
	filter Filter,
) ([]Candidate, int64, error) {
	args := []any{electionID}
	where := ""

	// status filter
	if filter.Status != nil {
		args = append(args, *filter.Status)
		where += fmt.Sprintf(" AND status = $%d", len(args))
	}

	// simple search by name/tagline
	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		where += fmt.Sprintf(" AND (name ILIKE $%d OR tagline ILIKE $%d)", len(args), len(args))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	// total count
	countSQL := qCountCandidatesBase + where
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list query
	listSQL := qListCandidatesBase + where + `
ORDER BY number ASC
LIMIT $` + fmt.Sprint(len(args)+1) + `
OFFSET $` + fmt.Sprint(len(args)+2)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		if isUndefinedColumn(err, "photo_media_id") {
			listSQL = qListCandidatesBaseNoPhotoMedia + where + `
ORDER BY number ASC
LIMIT $` + fmt.Sprint(len(args)-1) + `
OFFSET $` + fmt.Sprint(len(args))
			// args already has limit, offset at the end
			rows, err = r.db.Query(ctx, listSQL, args...)
		}
		if err != nil {
			return nil, 0, err
		}
	}
	defer rows.Close()

	var candidates []Candidate
	for rows.Next() {
		c, err := scanCandidate(rows)
		if err != nil {
			return nil, 0, err
		}
		candidates = append(candidates, c)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return candidates, total, nil
}

const qGetCandidateByID = `
SELECT
id,
election_id,
number,
name,
photo_url,
photo_media_id::text AS photo_media_id,
short_bio,
long_bio,
tagline,
faculty_name,
study_program_name,
cohort_year,
vision,
missions,
main_programs,
media,
social_links,
status,
created_at,
updated_at
FROM candidates
WHERE election_id = $1 AND id = $2
`

const qGetCandidateByIDNoPhotoMedia = `
SELECT
id,
election_id,
number,
name,
photo_url,
NULL::text AS photo_media_id,
short_bio,
long_bio,
tagline,
faculty_name,
study_program_name,
cohort_year,
vision,
missions,
main_programs,
media,
social_links,
status,
created_at,
updated_at
FROM candidates
WHERE election_id = $1 AND id = $2
`

const qGetCandidateByCandidateID = `
SELECT
id,
election_id,
number,
name,
photo_url,
photo_media_id::text AS photo_media_id,
short_bio,
long_bio,
tagline,
faculty_name,
study_program_name,
cohort_year,
vision,
missions,
main_programs,
media,
social_links,
status,
created_at,
updated_at
FROM candidates
WHERE id = $1
`

const qGetCandidateByCandidateIDNoPhotoMedia = `
SELECT
id,
election_id,
number,
name,
photo_url,
NULL::text AS photo_media_id,
short_bio,
long_bio,
tagline,
faculty_name,
study_program_name,
cohort_year,
vision,
missions,
main_programs,
media,
social_links,
status,
created_at,
updated_at
FROM candidates
WHERE id = $1
`

// GetByID returns a single candidate by election and candidate ID
func (r *PgCandidateRepository) GetByID(
	ctx context.Context,
	electionID, candidateID int64,
) (*Candidate, error) {
	row := r.db.QueryRow(ctx, qGetCandidateByID, electionID, candidateID)
	c, err := scanCandidateRow(row)
	if err != nil && isUndefinedColumn(err, "photo_media_id") {
		row = r.db.QueryRow(ctx, qGetCandidateByIDNoPhotoMedia, electionID, candidateID)
		c, err = scanCandidateRow(row)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCandidateNotFound
		}
		return nil, err
	}

	return &c, nil
}

func (r *PgCandidateRepository) GetByCandidateID(
	ctx context.Context,
	candidateID int64,
) (*Candidate, error) {
	row := r.db.QueryRow(ctx, qGetCandidateByCandidateID, candidateID)
	c, err := scanCandidateRow(row)
	if err != nil && isUndefinedColumn(err, "photo_media_id") {
		row = r.db.QueryRow(ctx, qGetCandidateByCandidateIDNoPhotoMedia, candidateID)
		c, err = scanCandidateRow(row)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCandidateNotFound
		}
		return nil, err
	}
	return &c, nil
}

// scanCandidate scans a candidate from pgx.Rows
func scanCandidate(rows pgx.Rows) (Candidate, error) {
	var c Candidate
	var missionsRaw, mainProgramsRaw, mediaRaw, socialLinksRaw any
	var photoURL, shortBio, longBio, tagline, facultyName, studyProgramName, vision *string

	if err := rows.Scan(
		&c.ID,
		&c.ElectionID,
		&c.Number,
		&c.Name,
		&photoURL,
		&c.PhotoMediaID,
		&shortBio,
		&longBio,
		&tagline,
		&facultyName,
		&studyProgramName,
		&c.CohortYear,
		&vision,
		&missionsRaw,
		&mainProgramsRaw,
		&mediaRaw,
		&socialLinksRaw,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return Candidate{}, err
	}

	if err := scanJSON(missionsRaw, &c.Missions); err != nil {
		return Candidate{}, err
	}
	if err := scanJSON(mainProgramsRaw, &c.MainPrograms); err != nil {
		return Candidate{}, err
	}
	if err := scanJSON(mediaRaw, &c.Media); err != nil {
		return Candidate{}, err
	}
	if err := scanJSON(socialLinksRaw, &c.SocialLinks); err != nil {
		return Candidate{}, err
	}

	c.PhotoURL = derefString(photoURL)
	c.ShortBio = derefString(shortBio)
	c.LongBio = derefString(longBio)
	c.Tagline = derefString(tagline)
	c.FacultyName = derefString(facultyName)
	c.StudyProgramName = derefString(studyProgramName)
	c.Vision = derefString(vision)

	return c, nil
}

// scanCandidateRow scans a candidate from pgx.Row
func scanCandidateRow(row pgx.Row) (Candidate, error) {
	var c Candidate
	var missionsRaw, mainProgramsRaw, mediaRaw, socialLinksRaw any
	var photoURL, shortBio, longBio, tagline, facultyName, studyProgramName, vision *string

	err := row.Scan(
		&c.ID,
		&c.ElectionID,
		&c.Number,
		&c.Name,
		&photoURL,
		&c.PhotoMediaID,
		&shortBio,
		&longBio,
		&tagline,
		&facultyName,
		&studyProgramName,
		&c.CohortYear,
		&vision,
		&missionsRaw,
		&mainProgramsRaw,
		&mediaRaw,
		&socialLinksRaw,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return Candidate{}, err
	}

	if err := scanJSON(missionsRaw, &c.Missions); err != nil {
		return Candidate{}, err
	}
	if err := scanJSON(mainProgramsRaw, &c.MainPrograms); err != nil {
		return Candidate{}, err
	}
	if err := scanJSON(mediaRaw, &c.Media); err != nil {
		return Candidate{}, err
	}
	if err := scanJSON(socialLinksRaw, &c.SocialLinks); err != nil {
		return Candidate{}, err
	}

	c.PhotoURL = derefString(photoURL)
	c.ShortBio = derefString(shortBio)
	c.LongBio = derefString(longBio)
	c.Tagline = derefString(tagline)
	c.FacultyName = derefString(facultyName)
	c.StudyProgramName = derefString(studyProgramName)
	c.Vision = derefString(vision)

	return c, nil
}

// scanJSON scans JSONB data into a Go type
func scanJSON[T any](src any, dest *T) error {
	if src == nil {
		return nil
	}

	var b []byte

	switch v := src.(type) {
	case []byte:
		if len(v) == 0 {
			return nil
		}
		b = v
	case string:
		if v == "" {
			return nil
		}
		b = []byte(v)
	default:
		var err error
		b, err = json.Marshal(v)
		if err != nil {
			logJSONError(fmt.Errorf("marshal fallback failed for type %T: %w", src, err))
			return err
		}
	}

	err := json.Unmarshal(b, dest)
	if err != nil {
		trimmed := bytes.TrimSpace(b)
		destType := reflect.TypeOf(*dest)

		if destType != nil {
			if destType.Kind() == reflect.Struct && bytes.Equal(trimmed, []byte("[]")) {
				slog.Warn("scanJSON: expected object, got array; using zero value", "dest_type", destType.String())
				var zero T
				*dest = zero
				return nil
			}
			if destType.Kind() == reflect.Slice && bytes.Equal(trimmed, []byte("{}")) {
				slog.Warn("scanJSON: expected array, got object; using empty value", "dest_type", destType.String())
				var zero T
				*dest = zero
				return nil
			}
			if destType.Kind() == reflect.Slice && destType.Elem().Kind() == reflect.String && len(trimmed) > 0 && trimmed[0] == '"' {
				var single string
				if uerr := json.Unmarshal(trimmed, &single); uerr == nil {
					destValue := reflect.ValueOf(dest).Elem()
					slice := reflect.MakeSlice(destType, 1, 1)
					slice.Index(0).SetString(single)
					destValue.Set(slice)
					slog.Warn("scanJSON: expected string array, got string; wrapped into array", "dest_type", destType.String())
					return nil
				}
			}
		}

		logJSONError(fmt.Errorf("scanJSON unmarshal error: %w, data: %s", err, string(b)))
		return err
	}
	return nil
}

func logJSONError(err error) {
	slog.Error("candidate JSON scan failed", "err", err)
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func isUndefinedColumn(err error, column string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	if pgErr.Code != "42703" { // undefined_column
		return false
	}
	if column == "" {
		return true
	}
	if strings.EqualFold(pgErr.ColumnName, column) {
		return true
	}
	// ColumnName isn't always populated depending on where the error happened.
	return strings.Contains(pgErr.Message, fmt.Sprintf("column \"%s\"", column))
}

const qCreateCandidate = `
INSERT INTO candidates (
election_id, number, name, photo_url, short_bio, long_bio, tagline,
faculty_name, study_program_name, cohort_year, vision, missions,
main_programs, media, social_links, status
) VALUES (
$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
RETURNING id, election_id, number, name, photo_url, photo_media_id::text AS photo_media_id, short_bio, long_bio, tagline,
faculty_name, study_program_name, cohort_year, vision, missions, main_programs,
media, social_links, status, created_at, updated_at
`

// Create creates a new candidate
func (r *PgCandidateRepository) Create(ctx context.Context, candidate *Candidate) (*Candidate, error) {
	// Ensure non-nil slices and valid JSON
	if candidate.Missions == nil {
		candidate.Missions = []string{}
	}
	if candidate.MainPrograms == nil {
		candidate.MainPrograms = []MainProgram{}
	}
	if candidate.SocialLinks == nil {
		candidate.SocialLinks = []SocialLink{}
	}

	missionsJSON, err := json.Marshal(candidate.Missions)
	if err != nil {
		return nil, fmt.Errorf("marshal missions: %w", err)
	}
	mainProgramsJSON, err := json.Marshal(candidate.MainPrograms)
	if err != nil {
		return nil, fmt.Errorf("marshal main_programs: %w", err)
	}
	mediaJSON, err := json.Marshal(candidate.Media)
	if err != nil {
		return nil, fmt.Errorf("marshal media: %w", err)
	}
	socialLinksJSON, err := json.Marshal(candidate.SocialLinks)
	if err != nil {
		return nil, fmt.Errorf("marshal social_links: %w", err)
	}
	
	slog.Info("Creating candidate", 
		"missions", string(missionsJSON),
		"mainPrograms", string(mainProgramsJSON),
		"media", string(mediaJSON),
		"socialLinks", string(socialLinksJSON))

	row := r.db.QueryRow(ctx, qCreateCandidate,
		candidate.ElectionID,
		candidate.Number,
		candidate.Name,
		candidate.PhotoURL,
		candidate.ShortBio,
		candidate.LongBio,
		candidate.Tagline,
		candidate.FacultyName,
		candidate.StudyProgramName,
		candidate.CohortYear,
		candidate.Vision,
		string(missionsJSON),
		string(mainProgramsJSON),
		string(mediaJSON),
		string(socialLinksJSON),
		candidate.Status,
	)

	c, err := scanCandidateRow(row)
	if err != nil {
		slog.Error("Failed to scan candidate row after insert", "error", err)
		return nil, fmt.Errorf("scan candidate: %w", err)
	}

	return &c, nil
}

const qUpdateCandidate = `
UPDATE candidates SET
number = $3,
name = $4,
photo_url = $5,
short_bio = $6,
long_bio = $7,
tagline = $8,
faculty_name = $9,
study_program_name = $10,
cohort_year = $11,
vision = $12,
missions = $13,
main_programs = $14,
media = $15,
social_links = $16,
status = $17,
updated_at = NOW()
WHERE election_id = $1 AND id = $2
RETURNING id, election_id, number, name, photo_url, photo_media_id::text AS photo_media_id, short_bio, long_bio, tagline,
faculty_name, study_program_name, cohort_year, vision, missions, main_programs,
media, social_links, status, created_at, updated_at
`

// Update updates an existing candidate
func (r *PgCandidateRepository) Update(ctx context.Context, electionID, candidateID int64, candidate *Candidate) (*Candidate, error) {
	missionsJSON, _ := json.Marshal(candidate.Missions)
	mainProgramsJSON, _ := json.Marshal(candidate.MainPrograms)
	mediaJSON, _ := json.Marshal(candidate.Media)
	socialLinksJSON, _ := json.Marshal(candidate.SocialLinks)

	row := r.db.QueryRow(ctx, qUpdateCandidate,
		electionID,
		candidateID,
		candidate.Number,
		candidate.Name,
		candidate.PhotoURL,
		candidate.ShortBio,
		candidate.LongBio,
		candidate.Tagline,
		candidate.FacultyName,
		candidate.StudyProgramName,
		candidate.CohortYear,
		candidate.Vision,
		missionsJSON,
		mainProgramsJSON,
		mediaJSON,
		socialLinksJSON,
		candidate.Status,
	)

	c, err := scanCandidateRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCandidateNotFound
		}
		return nil, err
	}

	return &c, nil
}

const qDeleteCandidate = `
DELETE FROM candidates WHERE election_id = $1 AND id = $2
`

// Delete deletes a candidate
func (r *PgCandidateRepository) Delete(ctx context.Context, electionID, candidateID int64) error {
	result, err := r.db.Exec(ctx, qDeleteCandidate, electionID, candidateID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCandidateNotFound
	}

	return nil
}

const qUpdateStatus = `
UPDATE candidates SET status = $3, updated_at = NOW()
WHERE election_id = $1 AND id = $2
`

// UpdateStatus updates candidate status
func (r *PgCandidateRepository) UpdateStatus(ctx context.Context, electionID, candidateID int64, status CandidateStatus) error {
	result, err := r.db.Exec(ctx, qUpdateStatus, electionID, candidateID, status)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCandidateNotFound
	}

	return nil
}

const qCheckNumberExists = `
SELECT EXISTS(
SELECT 1 FROM candidates
WHERE election_id = $1 AND number = $2 AND ($3::bigint IS NULL OR id != $3)
)
`

// CheckNumberExists checks if candidate number is already taken in an election
func (r *PgCandidateRepository) CheckNumberExists(ctx context.Context, electionID int64, number int, excludeCandidateID *int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, qCheckNumberExists, electionID, number, excludeCandidateID).Scan(&exists)
	return exists, err
}

func scanCandidateMedia(row pgx.Row) (*CandidateMedia, error) {
	var media CandidateMedia
	err := row.Scan(
		&media.ID,
		&media.CandidateID,
		&media.Slot,
		&media.FileName,
		&media.ContentType,
		&media.SizeBytes,
		&media.Data,
		&media.CreatedAt,
		&media.CreatedByID,
	)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

// SaveProfileMedia stores/replace profile media and updates photo_media_id
func (r *PgCandidateRepository) SaveProfileMedia(
	ctx context.Context,
	candidateID int64,
	media CandidateMediaCreate,
) (*CandidateMedia, error) {
	// Upload to Supabase Storage
	storage, err := newSupabaseStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	// Generate path and upload
	ext := getExtension(media.ContentType)
	path := fmt.Sprintf("candidates/%d/profile_%d%s", candidateID, time.Now().Unix(), ext)
	bucket := getMediaBucket()

	publicURL, err := storage.Upload(ctx, bucket, path, media.Data, media.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to storage: %w", err)
	}

	// Update candidate photo_url in DB
	_, err = r.db.Exec(ctx, `
UPDATE candidates
SET photo_url = $1,
    updated_by_admin_id = $2,
    updated_at = NOW()
WHERE id = $3
`, publicURL, media.CreatedByID, candidateID)
	if err != nil {
		// Rollback: delete from Supabase
		_ = storage.Delete(ctx, bucket, path)
		return nil, err
	}

	return &CandidateMedia{
		ID:          media.ID,
		CandidateID: candidateID,
		Slot:        CandidateMediaSlotProfile,
		FileName:    media.FileName,
		ContentType: media.ContentType,
		SizeBytes:   media.SizeBytes,
		URL:         publicURL,
		CreatedAt:   time.Now(),
		CreatedByID: &media.CreatedByID,
	}, nil
}

// GetProfileMedia retrieves profile media; 404 if missing
func (r *PgCandidateRepository) GetProfileMedia(ctx context.Context, candidateID int64) (*CandidateMedia, error) {
	var photoURL *string
	err := r.db.QueryRow(ctx, `
SELECT photo_url FROM candidates WHERE id = $1
`, candidateID).Scan(&photoURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCandidateNotFound
		}
		return nil, err
	}
	if photoURL == nil || *photoURL == "" {
		return nil, ErrCandidateMediaNotFound
	}

	// Return media with URL from Supabase (photo_url field)
	return &CandidateMedia{
		CandidateID: candidateID,
		Slot:        CandidateMediaSlotProfile,
		URL:         *photoURL,
		ContentType: "image/jpeg",
	}, nil
}

// DeleteProfileMedia removes profile media and clears reference
func (r *PgCandidateRepository) DeleteProfileMedia(ctx context.Context, candidateID int64, adminID int64) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var mediaID *string
	err = tx.QueryRow(ctx, `
SELECT photo_media_id FROM candidates WHERE id = $1 FOR UPDATE
`, candidateID).Scan(&mediaID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrCandidateNotFound
		}
		return err
	}
	if mediaID == nil {
		return ErrCandidateMediaNotFound
	}

	if _, err := tx.Exec(ctx, `
UPDATE candidates
SET photo_media_id = NULL,
    updated_by_admin_id = $2,
    updated_at = NOW()
WHERE id = $1
`, candidateID, adminID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM candidate_media WHERE id = $1`, *mediaID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// AddMedia stores non-profile media
func (r *PgCandidateRepository) AddMedia(ctx context.Context, candidateID int64, media CandidateMediaCreate) (*CandidateMedia, error) {
	if media.Slot == CandidateMediaSlotProfile {
		return nil, ErrInvalidCandidateMediaSlot
	}

	var exists bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM candidates WHERE id = $1)`, candidateID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrCandidateNotFound
	}

	var createdAt time.Time
	err := r.db.QueryRow(ctx, `
INSERT INTO candidate_media (id, candidate_id, slot, file_name, content_type, size_bytes, data, created_by_admin_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING created_at
`, media.ID, candidateID, media.Slot, media.FileName, media.ContentType, media.SizeBytes, media.Data, media.CreatedByID).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	return &CandidateMedia{
		ID:          media.ID,
		CandidateID: candidateID,
		Slot:        media.Slot,
		FileName:    media.FileName,
		ContentType: media.ContentType,
		SizeBytes:   media.SizeBytes,
		Data:        media.Data,
		CreatedAt:   createdAt,
		CreatedByID: &media.CreatedByID,
	}, nil
}

// GetMedia fetches any media by id scoped to candidate
func (r *PgCandidateRepository) GetMedia(ctx context.Context, candidateID int64, mediaID string) (*CandidateMedia, error) {
	row := r.db.QueryRow(ctx, `
SELECT id, candidate_id, slot, file_name, content_type, size_bytes, data, created_at, created_by_admin_id
FROM candidate_media
WHERE candidate_id = $1 AND id = $2
`, candidateID, mediaID)
	media, err := scanCandidateMedia(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCandidateMediaNotFound
		}
		return nil, err
	}
	return media, nil
}

// DeleteMedia deletes a media asset and clears profile reference if needed
func (r *PgCandidateRepository) DeleteMedia(ctx context.Context, candidateID int64, mediaID string) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var slot CandidateMediaSlot
	err = tx.QueryRow(ctx, `
SELECT slot FROM candidate_media WHERE candidate_id = $1 AND id = $2 FOR UPDATE
`, candidateID, mediaID).Scan(&slot)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrCandidateMediaNotFound
		}
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM candidate_media WHERE id = $1`, mediaID); err != nil {
		return err
	}

	// Clear photo reference if the deleted media was profile
	if slot == CandidateMediaSlotProfile {
		if _, err := tx.Exec(ctx, `
UPDATE candidates SET photo_media_id = NULL, updated_at = NOW()
WHERE id = $1 AND photo_media_id = $2
`, candidateID, mediaID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// ListMediaMeta returns lightweight metadata for a candidate
func (r *PgCandidateRepository) ListMediaMeta(ctx context.Context, candidateID int64) ([]CandidateMediaMeta, error) {
	rows, err := r.db.Query(ctx, `
SELECT id, slot, file_name, content_type, size_bytes, created_at
FROM candidate_media
WHERE candidate_id = $1
ORDER BY created_at DESC
`, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CandidateMediaMeta
	for rows.Next() {
		var meta CandidateMediaMeta
		if err := rows.Scan(&meta.ID, &meta.Slot, &meta.Label, &meta.ContentType, &meta.SizeBytes, &meta.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, meta)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return items, nil
}

// Helper functions for Supabase
func newSupabaseStorage() (*supabaseStorage, error) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SECRET_KEY")
	if url == "" || key == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_SECRET_KEY required")
	}

	// Add apikey header for Supabase auth
	headers := map[string]string{
		"apikey": key,
	}
	client := storage_go.NewClient(url+"/storage/v1", key, headers)
	return &supabaseStorage{client: client, url: url}, nil
}

type supabaseStorage struct {
	client *storage_go.Client
	url    string
}

func (s *supabaseStorage) Upload(ctx context.Context, bucket, path string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.client.UploadFile(bucket, path, reader)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.url, bucket, path), nil
}

func (s *supabaseStorage) Delete(ctx context.Context, bucket, path string) error {
	_, err := s.client.RemoveFile(bucket, []string{path})
	return err
}

func getExtension(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}

func getMediaBucket() string {
	bucket := os.Getenv("SUPABASE_MEDIA_BUCKET")
	if bucket == "" {
		return "pemira"
	}
	return bucket
}

// GetActiveQRCode returns the active QR code for a candidate
func (r *PgCandidateRepository) GetActiveQRCode(ctx context.Context, candidateID int64) (*CandidateQRCode, error) {
	query := `
		SELECT id, election_id, candidate_id, version, qr_token, is_active
		FROM candidate_qr_codes
		WHERE candidate_id = $1 AND is_active = true
		LIMIT 1
	`

	var qr CandidateQRCode
	err := r.db.QueryRow(ctx, query, candidateID).Scan(
		&qr.ID,
		&qr.ElectionID,
		&qr.CandidateID,
		&qr.Version,
		&qr.QRToken,
		&qr.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No QR code found
		}
		return nil, err
	}

	return &qr, nil
}

// GetQRCodesByElection returns all active QR codes for candidates in an election
func (r *PgCandidateRepository) GetQRCodesByElection(ctx context.Context, electionID int64) (map[int64]*CandidateQRCode, error) {
	query := `
		SELECT id, election_id, candidate_id, version, qr_token, is_active
		FROM candidate_qr_codes
		WHERE election_id = $1 AND is_active = true
		ORDER BY candidate_id
	`

	rows, err := r.db.Query(ctx, query, electionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qrCodes := make(map[int64]*CandidateQRCode)
	for rows.Next() {
		var qr CandidateQRCode
		err := rows.Scan(
			&qr.ID,
			&qr.ElectionID,
			&qr.CandidateID,
			&qr.Version,
			&qr.QRToken,
			&qr.IsActive,
		)
		if err != nil {
			return nil, err
		}
		qrCodes[qr.CandidateID] = &qr
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return qrCodes, nil
}
