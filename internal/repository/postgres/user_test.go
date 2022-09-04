package postgres

import (
	"errors"
	"log"
	"testing"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserPostgres_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewUserPostgres(db)

	type args struct {
		entity.User
	}

	type mockBehavior func(args args)

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
				entity.User{
					Id:       1,
					Username: "test_username",
				},
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id", "username"}).
					AddRow(args.Id, args.Username)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(args.Id).
					WillReturnRows(rows)
			},
			want: entity.User{
				Id:       1,
				Username: "test_username",
			},
		},
		{
			name: "Fail Select",
			args: args{
				entity.User{
					Id:       1,
					Username: "test_username",
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(args.Id).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			user, err := r.GetUserById(testCase.args.Id)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, user)
		})
	}
}

func TestUserPostgres_FindUserByName(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewUserPostgres(db)

	type args struct {
		filter      string
		currentUser int
		user        entity.User
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         []entity.User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				filter:      "user",
				currentUser: 1,
				user: entity.User{
					Id:       2,
					Username: "test_user2",
				},
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id", "username"}).
					AddRow(args.user.Id, args.user.Username)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs("%"+args.filter+"%", args.currentUser).
					WillReturnRows(rows)
			},
			want: []entity.User{
				{
					Id:       2,
					Username: "test_user2",
				},
			},
		},
		{
			name: "OK Empty",
			args: args{
				filter:      "user",
				currentUser: 1,
				user: entity.User{
					Id:       2,
					Username: "test_user2",
				},
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id", "username"})
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs("%"+args.filter+"%", args.currentUser).
					WillReturnRows(rows)
			},
			want: []entity.User{},
		},
		{
			name: "Fail Select",
			args: args{
				filter:      "user",
				currentUser: 1,
				user: entity.User{
					Id:       2,
					Username: "test_user2",
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs("%"+args.filter+"%", args.currentUser).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			users, err := r.FindUserByName(testCase.args.filter, testCase.args.currentUser)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(testCase.want), len(users))

			for i := 0; i < len(testCase.want); i++ {
				assert.Equal(t, testCase.want[i], users[i])
			}
		})
	}
}
