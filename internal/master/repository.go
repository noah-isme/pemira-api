package master

import "context"

type Repository interface {
	GetAllFaculties(ctx context.Context) ([]Faculty, error)
	GetAllStudyPrograms(ctx context.Context) ([]StudyProgram, error)
	GetStudyProgramsByFaculty(ctx context.Context, facultyID int64) ([]StudyProgram, error)
	GetAllLecturerUnits(ctx context.Context) ([]LecturerUnit, error)
	GetAllLecturerPositions(ctx context.Context) ([]LecturerPosition, error)
	GetLecturerPositionsByCategory(ctx context.Context, category string) ([]LecturerPosition, error)
	GetAllStaffUnits(ctx context.Context) ([]StaffUnit, error)
	GetAllStaffPositions(ctx context.Context) ([]StaffPosition, error)

	// Lookup methods for validation and ID retrieval
	GetFacultyByName(ctx context.Context, name string) (*Faculty, error)
	GetStudyProgramByName(ctx context.Context, facultyID int64, name string) (*StudyProgram, error)
	GetLecturerUnitByName(ctx context.Context, name string) (*LecturerUnit, error)
	GetLecturerPositionByName(ctx context.Context, name string) (*LecturerPosition, error)
	GetStaffUnitByName(ctx context.Context, name string) (*StaffUnit, error)
	GetStaffPositionByName(ctx context.Context, name string) (*StaffPosition, error)
}
