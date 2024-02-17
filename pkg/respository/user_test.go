package repository

import (
	models "cleancode/pkg/entity"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_GetUserDetails(t *testing.T) {
	tests := []struct {
		name    string
		args    int
		stub    func(mockSQL sqlmock.Sqlmock)
		want    models.UserDetail
		wantErr error
	}{
		{
			name: "success",
			args: 2,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectQuery := `SELECT u.username,u.email,u.phonenumber FROM users u WHERE u.id = ?`
				mockSQL.ExpectQuery(expectQuery).WillReturnRows(sqlmock.NewRows([]string{"username", "email", "phonenumber"}).AddRow("bibin", "bibin@gmail.com", "9087678564"))
			},
			want: models.UserDetail{
				UserName:    "bibin",
				Email:       "bibin@gmail.com",
				PhoneNumber: "9087678564",
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectQuery := `SELECT u.firstname,u.lastname,u.email,u.phone FROM users u WHERE u.id = ?`
				mockSQL.ExpectQuery(expectQuery).WillReturnError(errors.New("error"))
			},
			want:    models.UserDetail{},
			wantErr: errors.New("could not get user details"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.stub(mockSQL)
			u := NewUserRepository(gormDB)

			result, err := u.UserDetails(tt.args)

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
