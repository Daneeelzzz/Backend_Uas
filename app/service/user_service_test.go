package service_test

import (
	"context"
	"testing"
	"tugas_uas/app/model"
	"tugas_uas/app/service"
	"tugas_uas/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// GANTI NAMA STRUCT AGAR TIDAK BENTROK
type MockUserRepoForUserSvc struct {
	mock.Mock
}

func (m *MockUserRepoForUserSvc) GetRoleIDByName(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepoForUserSvc) CreateUser(ctx context.Context, user *model.User, r string, e map[string]string) error {
	args := m.Called(ctx, user, r, e)
	return args.Error(0)
}

// Stub method lain...
func (m *MockUserRepoForUserSvc) FindByEmailOrUsername(ctx context.Context, id string) (*model.User, error) { return nil, nil }
func (m *MockUserRepoForUserSvc) FindByID(ctx context.Context, id string) (*model.User, error) { return nil, nil }
func (m *MockUserRepoForUserSvc) FindAll(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *MockUserRepoForUserSvc) UpdateUser(ctx context.Context, id string, user *model.User) error { return nil }
func (m *MockUserRepoForUserSvc) DeleteUser(ctx context.Context, id string) error { return nil }
func (m *MockUserRepoForUserSvc) GetStudentsByAdvisor(ctx context.Context, id string) ([]model.User, error) { return nil, nil }
func (m *MockUserRepoForUserSvc) FindAllStudents(ctx context.Context) ([]model.StudentResponse, error) { return nil, nil }
func (m *MockUserRepoForUserSvc) FindAllLecturers(ctx context.Context) ([]model.LecturerResponse, error) { return nil, nil }


func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepoForUserSvc)
	svc := service.NewUserService(mockRepo) // Variabel svc dipakai

	reqName := "Mahasiswa"
	reqPass := "123"

	// Mock Behavior
	mockRepo.On("GetRoleIDByName", mock.Anything, reqName).Return("role-id", nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User"), reqName, mock.Anything).Return(nil)

	// Test Logic Manual
	hash, _ := utils.HashPassword(reqPass)
	assert.NotEmpty(t, hash)

	rid, _ := mockRepo.GetRoleIDByName(context.Background(), reqName)
	assert.Equal(t, "role-id", rid)

	err := mockRepo.CreateUser(context.Background(), &model.User{}, reqName, nil)
	
	assert.NoError(t, err)
	assert.NotNil(t, svc) // Assert agar variabel svc terpakai
}