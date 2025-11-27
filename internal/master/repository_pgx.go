package master

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRepository struct {
	db *pgxpool.Pool
}

func NewPgxRepository(db *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{db: db}
}

func (r *PgxRepository) GetAllFaculties(ctx context.Context) ([]Faculty, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM faculties ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faculties []Faculty
	for rows.Next() {
		var f Faculty
		if err := rows.Scan(&f.ID, &f.Code, &f.Name, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		faculties = append(faculties, f)
	}
	return faculties, rows.Err()
}

func (r *PgxRepository) GetAllStudyPrograms(ctx context.Context) ([]StudyProgram, error) {
	query := `SELECT id, faculty_id, code, name, level, created_at, updated_at FROM study_programs ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []StudyProgram
	for rows.Next() {
		var p StudyProgram
		if err := rows.Scan(&p.ID, &p.FacultyID, &p.Code, &p.Name, &p.Level, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, rows.Err()
}

func (r *PgxRepository) GetStudyProgramsByFaculty(ctx context.Context, facultyID int64) ([]StudyProgram, error) {
	query := `SELECT id, faculty_id, code, name, level, created_at, updated_at 
	          FROM study_programs WHERE faculty_id = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, facultyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []StudyProgram
	for rows.Next() {
		var p StudyProgram
		if err := rows.Scan(&p.ID, &p.FacultyID, &p.Code, &p.Name, &p.Level, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, rows.Err()
}

func (r *PgxRepository) GetAllLecturerUnits(ctx context.Context) ([]LecturerUnit, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM lecturer_units ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []LecturerUnit
	for rows.Next() {
		var u LecturerUnit
		if err := rows.Scan(&u.ID, &u.Code, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (r *PgxRepository) GetAllLecturerPositions(ctx context.Context) ([]LecturerPosition, error) {
	query := `SELECT id, category, code, name, created_at, updated_at FROM lecturer_positions ORDER BY category, name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []LecturerPosition
	for rows.Next() {
		var p LecturerPosition
		if err := rows.Scan(&p.ID, &p.Category, &p.Code, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, rows.Err()
}

func (r *PgxRepository) GetLecturerPositionsByCategory(ctx context.Context, category string) ([]LecturerPosition, error) {
	query := `SELECT id, category, code, name, created_at, updated_at 
	          FROM lecturer_positions WHERE category = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []LecturerPosition
	for rows.Next() {
		var p LecturerPosition
		if err := rows.Scan(&p.ID, &p.Category, &p.Code, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, rows.Err()
}

func (r *PgxRepository) GetAllStaffUnits(ctx context.Context) ([]StaffUnit, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM staff_units ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []StaffUnit
	for rows.Next() {
		var u StaffUnit
		if err := rows.Scan(&u.ID, &u.Code, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (r *PgxRepository) GetAllStaffPositions(ctx context.Context) ([]StaffPosition, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM staff_positions ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []StaffPosition
	for rows.Next() {
		var p StaffPosition
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, rows.Err()
}

func (r *PgxRepository) GetFacultyByName(ctx context.Context, name string) (*Faculty, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM faculties WHERE name = $1`
	var f Faculty
	err := r.db.QueryRow(ctx, query, name).Scan(&f.ID, &f.Code, &f.Name, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *PgxRepository) GetStudyProgramByName(ctx context.Context, facultyID int64, name string) (*StudyProgram, error) {
	query := `SELECT id, faculty_id, code, name, level, created_at, updated_at 
	          FROM study_programs WHERE faculty_id = $1 AND name = $2`
	var p StudyProgram
	err := r.db.QueryRow(ctx, query, facultyID, name).Scan(
		&p.ID, &p.FacultyID, &p.Code, &p.Name, &p.Level, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PgxRepository) GetLecturerUnitByName(ctx context.Context, name string) (*LecturerUnit, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM lecturer_units WHERE name = $1`
	var u LecturerUnit
	err := r.db.QueryRow(ctx, query, name).Scan(&u.ID, &u.Code, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PgxRepository) GetLecturerPositionByName(ctx context.Context, name string) (*LecturerPosition, error) {
	query := `SELECT id, category, code, name, created_at, updated_at 
	          FROM lecturer_positions WHERE name = $1`
	var p LecturerPosition
	err := r.db.QueryRow(ctx, query, name).Scan(
		&p.ID, &p.Category, &p.Code, &p.Name, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PgxRepository) GetStaffUnitByName(ctx context.Context, name string) (*StaffUnit, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM staff_units WHERE name = $1`
	var u StaffUnit
	err := r.db.QueryRow(ctx, query, name).Scan(&u.ID, &u.Code, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PgxRepository) GetStaffPositionByName(ctx context.Context, name string) (*StaffPosition, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM staff_positions WHERE name = $1`
	var p StaffPosition
	err := r.db.QueryRow(ctx, query, name).Scan(&p.ID, &p.Code, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
