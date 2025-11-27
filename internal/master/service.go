package master

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetFacultyProgramOptions(ctx context.Context) ([]FacultyWithPrograms, error) {
	faculties, err := s.repo.GetAllFaculties(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get faculties: %w", err)
	}

	programs, err := s.repo.GetAllStudyPrograms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get study programs: %w", err)
	}

	programsByFaculty := make(map[int64][]string)
	for _, p := range programs {
		programName := fmt.Sprintf("%s %s", p.Level, p.Name)
		programsByFaculty[p.FacultyID] = append(programsByFaculty[p.FacultyID], programName)
	}

	result := make([]FacultyWithPrograms, 0, len(faculties))
	for _, f := range faculties {
		result = append(result, FacultyWithPrograms{
			Faculty:  f.Name,
			Programs: programsByFaculty[f.ID],
		})
	}

	return result, nil
}

func (s *Service) GetAllFaculties(ctx context.Context) ([]Faculty, error) {
	return s.repo.GetAllFaculties(ctx)
}

func (s *Service) GetAllStudyPrograms(ctx context.Context) ([]StudyProgram, error) {
	return s.repo.GetAllStudyPrograms(ctx)
}

func (s *Service) GetStudyProgramsByFaculty(ctx context.Context, facultyID int64) ([]StudyProgram, error) {
	return s.repo.GetStudyProgramsByFaculty(ctx, facultyID)
}

func (s *Service) GetAllLecturerUnits(ctx context.Context) ([]LecturerUnit, error) {
	return s.repo.GetAllLecturerUnits(ctx)
}

func (s *Service) GetAllLecturerPositions(ctx context.Context) ([]LecturerPosition, error) {
	return s.repo.GetAllLecturerPositions(ctx)
}

func (s *Service) GetLecturerPositionsByCategory(ctx context.Context, category string) ([]LecturerPosition, error) {
	return s.repo.GetLecturerPositionsByCategory(ctx, category)
}

func (s *Service) GetAllStaffUnits(ctx context.Context) ([]StaffUnit, error) {
	return s.repo.GetAllStaffUnits(ctx)
}

func (s *Service) GetAllStaffPositions(ctx context.Context) ([]StaffPosition, error) {
	return s.repo.GetAllStaffPositions(ctx)
}
