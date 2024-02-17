package usecase

import (
	models "cleancode/pkg/entity"
	mockRepository "cleancode/pkg/respository/mock"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	gomock "github.com/golang/mock/gomock"
)

func Test_AddAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock implementations for the repositories
	userRepo := mockRepository.NewMockUserRepository(ctrl)

	userUseCase := NewUserUseCase(userRepo)

	testData := map[string]struct {
		input   models.UserAddress
		stub    func(*mockRepository.MockUserRepository, models.UserAddress)
		wantErr error
	}{
		"success": {
			input: models.UserAddress{
				City:        "kochi",
				State:       "kerala",
				Street:      "maradu",
				Country:     "india",
				PhoneNumber: "8585823457",
				PostalCode:  "688541",
			},
			stub: func(userRepo *mockRepository.MockUserRepository, data models.UserAddress) {
				userRepo.EXPECT().AddAddress(1, data).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		"failure": {
			input: models.UserAddress{
				City:        "kochi",
				State:       "kerala",
				Street:      "maradu",
				Country:     "india",
				PhoneNumber: "8585823457",
				PostalCode:  "688541",
			},
			stub: func(userRepo *mockRepository.MockUserRepository, data models.UserAddress) {
				userRepo.EXPECT().AddAddress(1, data).Return(errors.New("could not add the address")).Times(1)
			},
			wantErr: errors.New("could not add the address"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(userRepo, test.input)
			err := userUseCase.AddAddress(1, test.input)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func Test_GetAllAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockRepository.NewMockUserRepository(ctrl)

	userUseCase := NewUserUseCase(userRepo)

	testData := map[string]struct {
		input   int
		stub    func(*mockRepository.MockUserRepository, int)
		want    []models.AddressInfoResponse
		wantErr error
	}{
		"success": {
			input: 0,
			stub: func(userRepo *mockRepository.MockUserRepository, data int) {
				userRepo.EXPECT().GetAllAddress(data).Times(1).Return([]models.AddressInfoResponse{}, nil)
			},
			want:    []models.AddressInfoResponse{},
			wantErr: nil,
		},
		"failed": {
			input: 1,
			stub: func(userRepo *mockRepository.MockUserRepository, data int) {
				userRepo.EXPECT().GetAllAddress(data).Times(1).Return([]models.AddressInfoResponse{}, errors.New("error"))
			},
			want:    []models.AddressInfoResponse{},
			wantErr: errors.New("error"),
		},
	}

	for _, test := range testData {
		test.stub(userRepo, test.input)
		result, err := userUseCase.GetAllAddress(test.input)
		assert.Equal(t, test.want, result)
		assert.Equal(t, test.wantErr, err)
	}
}
