package repository

import (
	"cleancode/pkg/entity"
	models "cleancode/pkg/entity"
	"errors"
	"reflect"
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
				expectQuery := `SELECT u.username,u.email,u.number FROM users u WHERE u.id = ?`
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
				expectQuery := `SELECT u.username,u.email,u.number FROM users u WHERE u.id = ?`
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

func Test_GetAllAddress(t *testing.T) {
	testCase := []struct {
		name    string
		args    int
		stub    func(mockSQL sqlmock.Sqlmock)
		want    []models.AddressInfoResponse
		wantErr error
	}{
		{
			name: "Success",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM addresses WHERE user_id = \$1`
				mockSQL.ExpectQuery(expectedQuery).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "street", "city", "state", "postal_code", "country"}).
						AddRow(1, "cheleri", "kannur", "kerala", "670604", "india"))
			},
			want: []entity.AddressInfoResponse{
				{
					ID:         1,
					Street:     "cheleri",
					City:       "kannur",
					State:      "kerala",
					PostalCode: "670604",
					Country:    "india",
				},
			},
			wantErr: nil,
		},
		{
			name: "failed",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM addresses WHERE user_id = \$1`
				mockSQL.ExpectQuery(expectedQuery).WithArgs(1).
					WillReturnError(errors.New("error"))
			},
			want:    []models.AddressInfoResponse{},
			wantErr: errors.New("error"),
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.stub(mockSQL)
			u := NewUserRepository(gormDB)
			result, err := u.GetAllAddress(tt.args)
			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserSignUp(t *testing.T) {
	type args struct {
		input models.UserSignUp
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(mockSQL sqlmock.Sqlmock)
		want       models.UserDetailsResponse
		wantErr    error
	}{
		{
			name: "success signup user",
			args: args{
				input: models.UserSignUp{
					Username: "bibin",
					Email:    "bibin@gmail.com",
					Password: "12345",
					Number:   "7565748990",
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users\(username, email, password, number\) VALUES\(\$1, \$2, \$3, \$4\) RETURNING id, username, email, number`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("bibin", "bibin@gmail.com", "12345", "7565748990").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "number"}).
						AddRow(1, "bibin", "bibin@gmail.com", "7565748990"))
			},

			want: models.UserDetailsResponse{
				Id:       1,
				Username: "bibin",
				Email:    "bibin@gmail.com",
				Number:   "7565748990",
			},
			wantErr: nil,
		},
		{
			name: "error signup user",
			args: args{
				input: models.UserSignUp{},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users\(username, email, password, number\) VALUES\(\$1, \$2, \$3, \$4\) RETURNING id, username, email, number`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("", "", "", "").
					WillReturnError(errors.New("email should be unique"))
			},

			want:    models.UserDetailsResponse{},
			wantErr: errors.New("email should be unique"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.beforeTest(mockSQL)
			u := NewUserRepository(gormDB)
			got, err := u.UserSignUp(tt.args.input)
			assert.Equal(t, tt.wantErr, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.UserSignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
