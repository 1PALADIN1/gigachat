package postgres

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestMessagePostgress_AddMessageToChat(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewMessagePostgress(db)

	type args struct {
		userId    int
		username  string
		chatId    int
		sendTime  time.Time
		message   string
		messageId int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         entity.Message
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId:    1,
				username:  "test_username",
				chatId:    1,
				messageId: 1,
				sendTime:  time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				message:   "Hello world!",
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(a.messageId)

				mock.ExpectQuery("INSERT INTO messages").
					WithArgs(a.message, a.sendTime, a.chatId, a.userId).
					WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"username"}).
					AddRow(a.username)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(a.userId).
					WillReturnRows(rows)
			},
			want: entity.Message{
				Id:       1,
				SendTime: time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				Text:     "Hello world!",
				UserId:   1,
				Username: "test_username",
				ChatId:   1,
			},
		},
		{
			name: "Fail Insert",
			args: args{
				userId:    1,
				username:  "test_username",
				chatId:    1,
				messageId: 1,
				sendTime:  time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				message:   "Hello world!",
			},
			mockBehavior: func(a args) {
				mock.ExpectQuery("INSERT INTO messages").
					WithArgs(a.message, a.sendTime, a.chatId, a.userId).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "Fail Select",
			args: args{
				userId:    1,
				username:  "test_username",
				chatId:    1,
				messageId: 1,
				sendTime:  time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				message:   "Hello world!",
			},
			mockBehavior: func(a args) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(a.messageId)

				mock.ExpectQuery("INSERT INTO messages").
					WithArgs(a.message, a.sendTime, a.chatId, a.userId).
					WillReturnRows(rows)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(a.userId).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			message, err := r.AddMessageToChat(testCase.args.userId, testCase.args.chatId, testCase.args.message, testCase.args.sendTime)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, message)
		})
	}
}

func TestMessagePostgress_GetAllMessages(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewMessagePostgress(db)

	type args struct {
		chatId    int
		sendTime  time.Time
		messageId int
		message   string
		userId    int
		username  string
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         []entity.Message
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				chatId:    1,
				sendTime:  time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				messageId: 1,
				message:   "Hello world!",
				userId:    1,
				username:  "test_username",
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id", "send_date_time", "message", "user_id", "chat_id", "username"}).
					AddRow(args.messageId, args.sendTime, args.message, args.userId, args.chatId, args.username)

				mock.ExpectQuery("SELECT (.+) FROM messages m INNER JOIN users u ON (.+) WHERE (.+)").
					WithArgs(args.chatId).
					WillReturnRows(rows)
			},
			want: []entity.Message{
				{
					Id:       1,
					SendTime: time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
					Text:     "Hello world!",
					UserId:   1,
					Username: "test_username",
					ChatId:   1,
				},
			},
		},
		{
			name: "Select Fail",
			args: args{
				chatId:    1,
				sendTime:  time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				messageId: 1,
				message:   "Hello world!",
				userId:    1,
				username:  "test_username",
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM messages m INNER JOIN users u ON (.+) WHERE (.+)").
					WithArgs(args.chatId).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			messages, err := r.GetAllMessages(testCase.args.chatId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(testCase.want), len(messages))

			for i := 0; i < len(testCase.want); i++ {
				assert.Equal(t, testCase.want[i], messages[i])
			}
		})
	}
}
