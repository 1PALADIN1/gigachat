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

func TestChatPostgres_GetChatIdByUsers(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewChatPostgres(db)

	type args struct {
		inputUserIds []int
		resultChatId int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		wantFound    bool
		wantChatId   int
	}{
		{
			name: "OK",
			args: args{
				inputUserIds: []int{1, 2},
				resultChatId: 3,
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"chat_id"}).
					AddRow(args.resultChatId)

				mock.ExpectQuery(`SELECT (.+) FROM users_chats out_t INNER JOIN (.+) ON (.+) WHERE (.+) GROUP BY (.+) HAVING (.+)`).
					WithArgs(args.inputUserIds[0], args.inputUserIds[1]).
					WillReturnRows(rows)
			},
			wantFound:  true,
			wantChatId: 3,
		},
		{
			name: "OK Not Found",
			args: args{
				inputUserIds: []int{1, 2},
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"chat_id"})

				mock.ExpectQuery(`SELECT (.+) FROM users_chats out_t INNER JOIN (.+) ON (.+) WHERE (.+) GROUP BY (.+) HAVING (.+)`).
					WithArgs(args.inputUserIds[0], args.inputUserIds[1]).
					WillReturnRows(rows)
			},
			wantFound: false,
		},
		{
			name: "Select Fail",
			args: args{
				inputUserIds: []int{1, 2},
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM users_chats out_t INNER JOIN (.+) ON (.+) WHERE (.+) GROUP BY (.+) HAVING (.+)`).
					WithArgs(args.inputUserIds[0], args.inputUserIds[1]).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "Not Enough Params",
			args: args{
				inputUserIds: []int{},
			},
			mockBehavior: func(args args) {},
			wantErr:      true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			chatId, found, err := r.GetChatIdByUsers(testCase.args.inputUserIds)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.wantFound, found)
			if !testCase.wantFound {
				return
			}

			assert.Equal(t, testCase.wantChatId, chatId)
		})
	}
}

func TestChatPostgres_CreateChat(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewChatPostgres(db)

	type args struct {
		chat         entity.Chat
		resultChatId int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		wantChatId   int
	}{
		{
			name: "OK",
			args: args{
				chat: entity.Chat{
					Title:       "test_chat",
					Description: "test_description",
					UserIds:     []int{1, 2},
				},
				resultChatId: 4,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := mock.NewRows([]string{"id"}).
					AddRow(args.resultChatId)

				mock.ExpectQuery("INSERT INTO chats").
					WithArgs(args.chat.Title, args.chat.Description).
					WillReturnRows(rows)

				for i := 0; i < len(args.chat.UserIds); i++ {
					mock.ExpectExec("INSERT INTO users_chats").
						WithArgs(args.chat.UserIds[i], args.resultChatId).
						WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
				}

				mock.ExpectCommit()
			},
			wantChatId: 4,
		},
		{
			name: "First Insert Fail",
			args: args{
				chat: entity.Chat{
					Title:       "test_chat",
					Description: "test_description",
					UserIds:     []int{1, 2},
				},
				resultChatId: 4,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO chats").
					WithArgs(args.chat.Title, args.chat.Description).
					WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Second Insert Fail",
			args: args{
				chat: entity.Chat{
					Title:       "test_chat",
					Description: "test_description",
					UserIds:     []int{1, 2},
				},
				resultChatId: 4,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := mock.NewRows([]string{"id"}).
					AddRow(args.resultChatId)

				mock.ExpectQuery("INSERT INTO chats").
					WithArgs(args.chat.Title, args.chat.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_chats").
					WithArgs(args.chat.UserIds[0], args.resultChatId).
					WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			chatId, err := r.CreateChat(testCase.args.chat)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.wantChatId, chatId)
		})
	}
}

func TestChatPostgres_GetAllChats(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewChatPostgres(db)

	type args struct {
		inputUserId int
		resultChat  entity.ChatResponse
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		want         []entity.ChatResponse
	}{
		{
			name: "OK",
			args: args{
				inputUserId: 1,
				resultChat: entity.ChatResponse{
					Chat: entity.Chat{
						Id:          1,
						Title:       "test_chat",
						Description: "test_description",
					},
					LastMessage:         "Hello world!",
					LastMessageUserId:   "1",
					LastMessageUsername: "test_username",
					LastMessageTime:     time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC).String(),
				},
			},
			mockBehavior: func(args args) {
				result := args.resultChat
				rows := mock.NewRows([]string{"id", "title", "description", "send_date_time", "message", "user_id", "username"}).
					AddRow(result.Chat.Id, result.Chat.Title, result.Chat.Description,
						result.LastMessageTime, result.LastMessage, result.LastMessageUserId, result.LastMessageUsername)

				mock.ExpectQuery("SELECT (.+) FROM users_chats us INNER JOIN chats (.+) LEFT JOIN (.+) INNER JOIN users (.+) WHERE (.+)").
					WithArgs(args.inputUserId).
					WillReturnRows(rows)
			},
			want: []entity.ChatResponse{
				{
					Chat: entity.Chat{
						Id:          1,
						Title:       "test_chat",
						Description: "test_description",
					},
					LastMessage:         "Hello world!",
					LastMessageUserId:   "1",
					LastMessageUsername: "test_username",
					LastMessageTime:     time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC).String(),
				},
			},
		},
		{
			name: "OK Empty Chats",
			args: args{
				inputUserId: 1,
				resultChat: entity.ChatResponse{
					Chat: entity.Chat{
						Id:          1,
						Title:       "test_chat",
						Description: "test_description",
					},
					LastMessage:         "Hello world!",
					LastMessageUserId:   "1",
					LastMessageUsername: "test_username",
					LastMessageTime:     time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC).String(),
				},
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id", "title", "description", "send_date_time", "message", "user_id", "username"})

				mock.ExpectQuery("SELECT (.+) FROM users_chats us INNER JOIN chats (.+) LEFT JOIN (.+) INNER JOIN users (.+) WHERE (.+)").
					WithArgs(args.inputUserId).
					WillReturnRows(rows)
			},
			want: []entity.ChatResponse{},
		},
		{
			name: "Select Fail",
			args: args{
				inputUserId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM users_chats us INNER JOIN chats (.+) LEFT JOIN (.+) INNER JOIN users (.+) WHERE (.+)").
					WithArgs(args.inputUserId).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			chats, err := r.GetAllChats(testCase.args.inputUserId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(testCase.want), len(chats))

			for i := 0; i < len(testCase.want); i++ {
				assert.Equal(t, testCase.want[i], chats[i])
			}
		})
	}
}

func TestChatPostgres_GetUserIdsByChatId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewChatPostgres(db)

	type args struct {
		inputChatId   int
		resultUserIds []int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		want         []int
	}{
		{
			name: "OK",
			args: args{
				inputChatId:   1,
				resultUserIds: []int{1, 2},
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id"})

				for _, userId := range args.resultUserIds {
					rows.AddRow(userId)
				}

				mock.ExpectQuery("SELECT (.+) FROM users_chats WHERE (.+)").
					WithArgs(args.inputChatId).
					WillReturnRows(rows)
			},
			want: []int{1, 2},
		},
		{
			name: "OK Empty",
			args: args{
				inputChatId: 1,
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id"})

				mock.ExpectQuery("SELECT (.+) FROM users_chats WHERE (.+)").
					WithArgs(args.inputChatId).
					WillReturnRows(rows)
			},
			want: []int{},
		},
		{
			name: "Select Fail",
			args: args{
				inputChatId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM users_chats WHERE (.+)").
					WithArgs(args.inputChatId).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			userIds, err := r.GetUserIdsByChatId(testCase.args.inputChatId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(testCase.want), len(userIds))

			for i := 0; i < len(testCase.want); i++ {
				assert.Equal(t, testCase.want[i], userIds[i])
			}
		})
	}
}
