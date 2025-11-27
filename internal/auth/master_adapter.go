package auth

import (
	"context"

	"pemira-api/internal/master"
)

type MasterRepositoryAdapter struct {
	repo master.Repository
}

func NewMasterRepositoryAdapter(repo master.Repository) *MasterRepositoryAdapter {
	return &MasterRepositoryAdapter{repo: repo}
}

func (a *MasterRepositoryAdapter) GetFacultyByName(ctx context.Context, name string) (*MasterFaculty, error) {
	faculty, err := a.repo.GetFacultyByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &MasterFaculty{
		ID:   faculty.ID,
		Name: faculty.Name,
	}, nil
}

func (a *MasterRepositoryAdapter) GetStudyProgramByName(ctx context.Context, facultyID int64, name string) (*MasterStudyProgram, error) {
	program, err := a.repo.GetStudyProgramByName(ctx, facultyID, name)
	if err != nil {
		return nil, err
	}
	return &MasterStudyProgram{
		ID:   program.ID,
		Name: program.Name,
	}, nil
}

func (a *MasterRepositoryAdapter) GetLecturerUnitByName(ctx context.Context, name string) (*MasterLecturerUnit, error) {
	unit, err := a.repo.GetLecturerUnitByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &MasterLecturerUnit{
		ID:   unit.ID,
		Name: unit.Name,
	}, nil
}

func (a *MasterRepositoryAdapter) GetLecturerPositionByName(ctx context.Context, name string) (*MasterLecturerPosition, error) {
	position, err := a.repo.GetLecturerPositionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &MasterLecturerPosition{
		ID:   position.ID,
		Name: position.Name,
	}, nil
}

func (a *MasterRepositoryAdapter) GetStaffUnitByName(ctx context.Context, name string) (*MasterStaffUnit, error) {
	unit, err := a.repo.GetStaffUnitByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &MasterStaffUnit{
		ID:   unit.ID,
		Name: unit.Name,
	}, nil
}

func (a *MasterRepositoryAdapter) GetStaffPositionByName(ctx context.Context, name string) (*MasterStaffPosition, error) {
	position, err := a.repo.GetStaffPositionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &MasterStaffPosition{
		ID:   position.ID,
		Name: position.Name,
	}, nil
}
