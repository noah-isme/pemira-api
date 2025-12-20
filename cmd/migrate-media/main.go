// cmd/migrate-media/main.go
// Migration script to upload blob data from candidate_media to Supabase Storage
// and update candidates.photo_url with the public URL.
//
// Usage:
//   DATABASE_URL=postgres://... SUPABASE_URL=... SUPABASE_SECRET_KEY=... go run cmd/migrate-media/main.go
//
// Flags:
//   -dry-run    Preview what would be migrated without making changes
//   -limit N    Limit number of records to migrate

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	storage_go "github.com/supabase-community/storage-go"
)

type CandidateMedia struct {
	ID          string
	CandidateID int64
	Slot        string
	FileName    string
	ContentType string
	SizeBytes   int64
	Data        []byte
	StoragePath *string
}

type MigrationResult struct {
	CandidateID int64
	MediaID     string
	OldSize     int64
	NewURL      string
	Error       error
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Preview without making changes")
	limit := flag.Int("limit", 0, "Limit number of records (0 = all)")
	flag.Parse()

	// Load configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_SECRET_KEY are required")
	}

	bucket := os.Getenv("SUPABASE_MEDIA_BUCKET")
	if bucket == "" {
		bucket = "pemira"
	}

	ctx := context.Background()

	// Connect to database
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Supabase storage client
	storageClient := storage_go.NewClient(supabaseURL+"/storage/v1", supabaseKey, nil)

	// Query candidate_media with profile slot that has blob data
	query := `
		SELECT 
			cm.id, 
			cm.candidate_id, 
			cm.slot, 
			cm.file_name, 
			cm.content_type, 
			cm.size_bytes, 
			cm.data,
			cm.storage_path
		FROM candidate_media cm
		JOIN candidates c ON c.id = cm.candidate_id
		WHERE cm.slot = 'profile'
		  AND cm.data IS NOT NULL
		  AND LENGTH(cm.data) > 0
		  AND (c.photo_url IS NULL OR c.photo_url = '')
	`
	if *limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}

	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Fatalf("Failed to query candidate_media: %v", err)
	}
	defer rows.Close()

	var mediaList []CandidateMedia
	for rows.Next() {
		var m CandidateMedia
		if err := rows.Scan(
			&m.ID,
			&m.CandidateID,
			&m.Slot,
			&m.FileName,
			&m.ContentType,
			&m.SizeBytes,
			&m.Data,
			&m.StoragePath,
		); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		mediaList = append(mediaList, m)
	}

	log.Printf("Found %d profile photos to migrate", len(mediaList))

	if *dryRun {
		log.Println("=== DRY RUN MODE - No changes will be made ===")
		for _, m := range mediaList {
			log.Printf("Would migrate: candidate_id=%d, size=%d bytes, content_type=%s",
				m.CandidateID, m.SizeBytes, m.ContentType)
		}
		return
	}

	// Migrate each photo
	var results []MigrationResult
	for _, m := range mediaList {
		result := migratePhoto(ctx, db, storageClient, bucket, m)
		results = append(results, result)

		if result.Error != nil {
			log.Printf("❌ Failed candidate_id=%d: %v", m.CandidateID, result.Error)
		} else {
			log.Printf("✅ Migrated candidate_id=%d → %s", m.CandidateID, result.NewURL)
		}
	}

	// Summary
	var success, failed int
	for _, r := range results {
		if r.Error == nil {
			success++
		} else {
			failed++
		}
	}
	log.Printf("\n=== Migration Complete ===")
	log.Printf("Success: %d, Failed: %d, Total: %d", success, failed, len(results))
}

func migratePhoto(
	ctx context.Context,
	db *pgxpool.Pool,
	storage *storage_go.Client,
	bucket string,
	m CandidateMedia,
) MigrationResult {
	result := MigrationResult{
		CandidateID: m.CandidateID,
		MediaID:     m.ID,
		OldSize:     m.SizeBytes,
	}

	// Generate storage path
	ext := getExtension(m.ContentType)
	path := fmt.Sprintf("candidates/%d/profile_%d%s", m.CandidateID, time.Now().Unix(), ext)

	// Upload to Supabase
	reader := bytes.NewReader(m.Data)
	_, err := storage.UploadFile(bucket, path, reader)
	if err != nil {
		result.Error = fmt.Errorf("upload failed: %w", err)
		return result
	}

	// Build public URL
	supabaseURL := os.Getenv("SUPABASE_URL")
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucket, path)
	result.NewURL = publicURL

	// Update candidates.photo_url
	_, err = db.Exec(ctx, `
		UPDATE candidates 
		SET photo_url = $1, updated_at = NOW() 
		WHERE id = $2
	`, publicURL, m.CandidateID)
	if err != nil {
		result.Error = fmt.Errorf("update candidates failed: %w", err)
		return result
	}

	// Update candidate_media.storage_path (for reference)
	_, err = db.Exec(ctx, `
		UPDATE candidate_media 
		SET storage_path = $1 
		WHERE id = $2
	`, path, m.ID)
	if err != nil {
		log.Printf("Warning: failed to update storage_path for media %s: %v", m.ID, err)
		// Don't fail the whole operation for this
	}

	return result
}

func getExtension(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".bin"
	}
}
