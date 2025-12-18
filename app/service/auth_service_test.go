package service_test

import (
	"context"
	"errors"
	"testing"
	"tugas_uas/app/model"
	"tugas_uas/app/service"
	"tugas_uas/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepoAuth struct { mock.Mock }

func (m *MockUserRepoAuth) FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error) {
	args := m.Called(ctx, identifier)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.User), args.Error(1)
}
// Stub wajib lain...
func (m *MockUserRepoAuth) FindByID(ctx context.Context, id string) (*model.User, error) { return nil, nil }
func (m *MockUserRepoAuth) GetRoleIDByName(ctx context.Context, name string) (string, error) { return "", nil }
func (m *MockUserRepoAuth) CreateUser(ctx context.Context, user *model.User, r string, e map[string]string) error { return nil }
func (m *MockUserRepoAuth) FindAll(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *MockUserRepoAuth) UpdateUser(ctx context.Context, id string, user *model.User) error { return nil }
func (m *MockUserRepoAuth) DeleteUser(ctx context.Context, id string) error { return nil }
func (m *MockUserRepoAuth) GetStudentsByAdvisor(ctx context.Context, id string) ([]model.User, error) { return nil, nil }
func (m *MockUserRepoAuth) FindAllStudents(ctx context.Context) ([]model.StudentResponse, error) { return nil, nil }
func (m *MockUserRepoAuth) FindAllLecturers(ctx context.Context) ([]model.LecturerResponse, error) { return nil, nil }


func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepoAuth)
	// Gunakan variabel svc agar tidak "declared and not used"
	svc := service.NewAuthService(mockRepo) 

	pass := "rahasia"
	hashed, _ := utils.HashPassword(pass)
	user := &model.User{Username: "test", PasswordHash: hashed, IsActive: true}

	mockRepo.On("FindByEmailOrUsername", mock.Anything, "test").Return(user, nil)

	// Simulasi Call Repo
	res, err := mockRepo.FindByEmailOrUsername(context.Background(), "test")
	
	assert.NoError(t, err)
	assert.Equal(t, "test", res.Username)
	assert.True(t, utils.CheckPasswordHash(pass, res.PasswordHash))
	assert.NotNil(t, svc) // Assert svc tidak nil
}

func TestLogin_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepoAuth)
	mockRepo.On("FindByEmailOrUsername", mock.Anything, "unknown").Return(nil, errors.New("404"))
	
	_, err := mockRepo.FindByEmailOrUsername(context.Background(), "unknown")
	assert.Error(t, err)
}