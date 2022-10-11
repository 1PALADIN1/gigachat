package postgres

import (
	"errors"
	"log"
	"testing"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type args struct {
		username     string
		passwordHash string
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				username:     "test_username",
				passwordHash: "test_password_hash",
			},
			id: 1,
			mockBehavior: func(args args, id int) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(args.username, args.passwordHash).
					WillReturnRows(rows)
			},
		},
		{
			name: "Insert Failure",
			args: args{
				username:     "",
				passwordHash: "test_password_hash",
			},
			id: 1,
			mockBehavior: func(args args, id int) {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(args.username, args.passwordHash).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			userId, err := r.CreateUser(entity.User{
				Username: testCase.args.username,
				Password: testCase.args.passwordHash,
			})

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.id, userId)
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type args struct {
		username     string
		passwordHash string
	}

	type mockBehavior func(args args, userId int)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         entity.User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				username:     "test_username",
				passwordHash: "test_password_hash",
			},
			mockBehavior: func(args args, userId int) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(userId)

				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(args.username, args.passwordHash).
					WillReturnRows(rows)
			},
			want: entity.User{
				Id: 1,
			},
		},
		{
			name: "Get Failure",
			args: args{
				username:     "test_username",
				passwordHash: "test_password_hash",
			},
			mockBehavior: func(args args, userId int) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(args.username, args.passwordHash).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.want.Id)

			user, err := r.GetUser(testCase.args.username, testCase.args.passwordHash)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.want.Id, user.Id)
		})
	}
}
