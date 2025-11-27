package adminuser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"pemira-api/internal/shared"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func NewPgRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) List(ctx context.Context, filter ListFilter, pag shared.PaginationParams) ([]User, int64, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.Search != "" {
		where = append(where, fmt.Sprintf("(username ILIKE $%d OR email ILIKE $%d OR full_name ILIKE $%d)", len(args)+1, len(args)+1, len(args)+1))
		args = append(args, "%"+filter.Search+"%")
	}
	if filter.Role != "" {
		where = append(where, fmt.Sprintf("role = $%d", len(args)+1))
		args = append(args, filter.Role)
	}
	if filter.Active != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *filter.Active)
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM user_accounts %s`, whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	limitPos := len(args) + 1
	offsetPos := len(args) + 2
	args = append(args, pag.Limit(), pag.Offset())
	listQuery := fmt.Sprintf(`
		SELECT id, username, email, full_name, role, voter_id, tps_id, lecturer_id, staff_id, is_active,
		       last_login_at, login_count, created_at, updated_at
		FROM user_accounts
		%s
		ORDER BY updated_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitPos, offsetPos)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var items []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return items, total, nil
}

func (r *pgRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, full_name, role, voter_id, tps_id, lecturer_id, staff_id, is_active,
		       last_login_at, login_count, created_at, updated_at
		FROM user_accounts
		WHERE id = $1
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, query, id)
	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *pgRepository) Create(ctx context.Context, in CreateInput, passwordHash string) (*User, error) {
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}

	query := `
		INSERT INTO user_accounts (username, email, full_name, password_hash, role, voter_id, tps_id, lecturer_id, staff_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, username, email, full_name, role, voter_id, tps_id, lecturer_id, staff_id, is_active,
		          last_login_at, login_count, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		in.Username,
		in.Email,
		in.FullName,
		passwordHash,
		in.Role,
		in.VoterID,
		in.TPSID,
		in.LecturerID,
		in.StaffID,
		isActive,
	)

	user, err := scanUser(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, shared.ErrDuplicateEntry
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (r *pgRepository) Update(ctx context.Context, id int64, in UpdateInput) (*User, error) {
	setParts := []string{}
	args := []interface{}{id}

	if in.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, *in.Email)
	}
	if in.FullName != nil {
		setParts = append(setParts, fmt.Sprintf("full_name = $%d", len(args)+1))
		args = append(args, *in.FullName)
	}
	if in.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role = $%d", len(args)+1))
		args = append(args, *in.Role)
	}
	if in.VoterID != nil {
		setParts = append(setParts, fmt.Sprintf("voter_id = $%d", len(args)+1))
		args = append(args, *in.VoterID)
	}
	if in.TPSID != nil {
		setParts = append(setParts, fmt.Sprintf("tps_id = $%d", len(args)+1))
		args = append(args, *in.TPSID)
	}
	if in.LecturerID != nil {
		setParts = append(setParts, fmt.Sprintf("lecturer_id = $%d", len(args)+1))
		args = append(args, *in.LecturerID)
	}
	if in.StaffID != nil {
		setParts = append(setParts, fmt.Sprintf("staff_id = $%d", len(args)+1))
		args = append(args, *in.StaffID)
	}
	if in.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *in.IsActive)
	}

	if len(setParts) == 0 {
		// nothing to update
		return r.GetByID(ctx, id)
	}
	setParts = append(setParts, "updated_at = NOW()")

	query := fmt.Sprintf(`
		UPDATE user_accounts
		SET %s
		WHERE id = $1
		RETURNING id, username, email, full_name, role, voter_id, tps_id, lecturer_id, staff_id, is_active,
		          last_login_at, login_count, created_at, updated_at
	`, strings.Join(setParts, ", "))

	row := r.db.QueryRow(ctx, query, args...)
	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, shared.ErrNotFound
		}
		if isUniqueViolation(err) {
			return nil, shared.ErrDuplicateEntry
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return user, nil
}

func (r *pgRepository) ResetPassword(ctx context.Context, id int64, passwordHash string) error {
	query := `UPDATE user_accounts SET password_hash = $2, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id, passwordHash)
	if err != nil {
		return fmt.Errorf("reset password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (r *pgRepository) SetActive(ctx context.Context, id int64, active bool) (*User, error) {
	query := `
		UPDATE user_accounts
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, email, full_name, role, voter_id, tps_id, lecturer_id, staff_id, is_active,
		          last_login_at, login_count, created_at, updated_at
	`
	row := r.db.QueryRow(ctx, query, id, active)
	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("set active: %w", err)
	}
	return user, nil
}

func (r *pgRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM user_accounts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func scanUser(scanner interface {
	Scan(dest ...interface{}) error
}) (*User, error) {
	var (
		voterID     sql.NullInt64
		tpsID       sql.NullInt64
		lecturerID  sql.NullInt64
		staffID     sql.NullInt64
		lastLoginAt sql.NullTime
		loginCount  sql.NullInt64
		user        User
	)

	err := scanner.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Role,
		&voterID,
		&tpsID,
		&lecturerID,
		&staffID,
		&user.IsActive,
		&lastLoginAt,
		&loginCount,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if voterID.Valid {
		user.VoterID = &voterID.Int64
	}
	if tpsID.Valid {
		user.TPSID = &tpsID.Int64
	}
	if lecturerID.Valid {
		user.LecturerID = &lecturerID.Int64
	}
	if staffID.Valid {
		user.StaffID = &staffID.Int64
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	if loginCount.Valid {
		val := int(loginCount.Int64)
		user.LoginCount = &val
	}

	return &user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}
