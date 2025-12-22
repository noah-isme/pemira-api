package electionvoter

import (
	"context"

	"pemira-api/internal/shared"
)

type Repository interface {
	LookupByNIM(ctx context.Context, electionID int64, nim string) (*LookupResult, error)
	UpsertAndEnroll(ctx context.Context, electionID int64, in UpsertAndEnrollInput) (*UpsertAndEnrollResult, error)
	List(ctx context.Context, electionID int64, filter ListFilter, pag shared.PaginationParams) ([]ElectionVoter, int64, error)
	UpdateEnrollment(ctx context.Context, electionID int64, enrollmentID int64, in UpdateInput) (*ElectionVoter, error)
	SelfRegister(ctx context.Context, electionID int64, voterID int64, in SelfRegisterInput) (*ElectionVoter, error)
	GetStatus(ctx context.Context, electionID int64, voterID int64) (*ElectionVoter, error)
	BlacklistVoter(ctx context.Context, electionID, voterID int64, reason string) error
	UnblacklistVoter(ctx context.Context, electionID, voterID int64) error
}
