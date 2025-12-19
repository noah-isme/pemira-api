package voting

import (
	"context"
	"fmt"
	
	"github.com/jackc/pgx/v5"
	"pemira-api/internal/candidate"
	"pemira-api/internal/shared"
)

type candidateRepository struct{}

func NewCandidateRepository() CandidateRepository {
	return &candidateRepository{}
}

func (r *candidateRepository) GetByIDWithTx(ctx context.Context, tx pgx.Tx, candidateID int64) (*candidate.Candidate, error) {
	query := `
		SELECT id, election_id, number, name, vision, photo_url, status, created_at, updated_at
		FROM candidates
		WHERE id = $1
	`
	
	var c candidate.Candidate
	
	var vision *string
	var status string
	
	err := tx.QueryRow(ctx, query, candidateID).Scan(
		&c.ID,
		&c.ElectionID,
		&c.Number,
		&c.Name,
		&vision,
		&c.PhotoURL,
		&status,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get candidate: %w", err)
	}
	
	return &c, nil
}
