package adminuser

import (
	"context"
	"strings"

	"pemira-api/internal/auth"
	"pemira-api/internal/shared"
	"pemira-api/internal/shared/constants"
)

var allowedRoles = map[constants.Role]struct{}{
	constants.RoleAdmin:              {},
	constants.RoleSuperAdmin:         {},
	constants.RoleTPSOperator:        {},
	constants.RoleStudent:            {},
	constants.RoleLecturer:           {},
	constants.RoleStaff:              {},
	constants.Role("PANITIA"):        {},
	constants.Role("KETUA_TPS"):      {},
	constants.Role("OPERATOR_PANEL"): {},
	constants.Role("VIEWER"):         {},
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filter ListFilter, page, limit int) ([]User, shared.PaginationMeta, error) {
	pag := shared.NewPaginationParams(page, limit)
	items, total, err := s.repo.List(ctx, filter, pag)
	if err != nil {
		return nil, shared.PaginationMeta{}, err
	}
	meta := shared.NewPaginatedResponse(nil, pag, total).Meta
	return items, meta, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*User, error) {
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.TrimSpace(in.Email)
	in.FullName = strings.TrimSpace(in.FullName)
	in.Role = constants.Role(strings.ToUpper(string(in.Role)))

	if in.Username == "" || in.Password == "" || in.FullName == "" || in.Email == "" {
		return nil, shared.ErrBadRequest
	}
	if !roleAllowed(in.Role) {
		return nil, shared.ErrBadRequest
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, in, hash)
}

func (s *Service) Update(ctx context.Context, id int64, in UpdateInput) (*User, error) {
	if in.Role != nil {
		r := constants.Role(strings.ToUpper(string(*in.Role)))
		if !roleAllowed(r) {
			return nil, shared.ErrBadRequest
		}
		in.Role = &r
	}
	return s.repo.Update(ctx, id, in)
}

func (s *Service) ResetPassword(ctx context.Context, id int64, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" || len(newPassword) < 6 {
		return shared.ErrBadRequest
	}
	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repo.ResetPassword(ctx, id, hash)
}

func (s *Service) SetActive(ctx context.Context, id int64, active bool) (*User, error) {
	return s.repo.SetActive(ctx, id, active)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func roleAllowed(role constants.Role) bool {
	_, ok := allowedRoles[role]
	return ok
}
