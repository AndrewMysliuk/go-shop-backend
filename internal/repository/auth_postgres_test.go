package repository

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	a := NewAuthPostgres(db)

	type args struct {
		dataId    string
		item      domain.UserSignUp
		createdAt time.Time
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				dataId: "34c8d3e6-b8d7-43dc-847e-5764c4114856",
				item: domain.UserSignUp{
					Name:     "Test_User",
					Surname:  "Test_Surname",
					Email:    "test@gmail.com",
					Phone:    "+4456781234",
					Password: "_9Z9sL~i3H4Kb33jcKJ9-8rZ+&-uRk#",
					Role:     "ADMIN",
				},
				createdAt: time.Now(),
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.dataId)
				prep := mock.ExpectPrepare("INSERT INTO users")
				prep.ExpectQuery().
					WithArgs(args.dataId, args.item.Name, args.item.Surname, args.item.Email, args.item.Phone, args.item.Role, args.item.Password, args.createdAt).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},

		{
			name: "Empty Fields",
			args: args{
				dataId: "",
				item: domain.UserSignUp{
					Name:     "",
					Surname:  "",
					Email:    "",
					Phone:    "",
					Password: "",
					Role:     "",
				},
				createdAt: time.Now(),
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.dataId).RowError(0, errors.New("insert error"))
				prep := mock.ExpectPrepare("INSERT INTO users")
				prep.ExpectQuery().
					WithArgs(args.dataId, args.item.Name, args.item.Surname, args.item.Email, args.item.Phone, args.item.Role, args.item.Password, args.createdAt).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := a.CreateUser(testCase.args.item, testCase.args.dataId, testCase.args.createdAt)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.dataId, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	a := NewAuthPostgres(db)

	type args struct {
		email    string
		password string
	}

	testTable := []struct {
		name    string
		mock    func()
		args    args
		want    domain.User
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "surname", "email", "phone", "role", "password_hash", "created_at"}).
					AddRow("34c8d3e6-b8d7-43dc-847e-5764c4114856", "Test_Name", "Test_Surname", "test@gmail.com", "+4412345678", "ADMIN", "_9Z9sL~i3H4Kb33jcKJ9-8rZ+&-uRk#", time.Date(2022, 01, 12, 13, 8, 21, 32963, time.Local))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1 AND password_hash = $2")).
					WithArgs("test@gmail.com", "_9Z9sL~i3H4Kb33jcKJ9-8rZ+&-uRk#").
					WillReturnRows(rows)
			},
			args: args{
				email:    "test@gmail.com",
				password: "_9Z9sL~i3H4Kb33jcKJ9-8rZ+&-uRk#",
			},
			want: domain.User{
				Id:        "34c8d3e6-b8d7-43dc-847e-5764c4114856",
				Name:      "Test_Name",
				Surname:   "Test_Surname",
				Email:     "test@gmail.com",
				Phone:     "+4412345678",
				Role:      "ADMIN",
				Password:  "_9Z9sL~i3H4Kb33jcKJ9-8rZ+&-uRk#",
				CreatedAt: time.Date(2022, 01, 12, 13, 8, 21, 32963, time.Local),
			},
		},

		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "surname", "email", "phone", "role", "password_hash", "created_at"})

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1 AND password_hash = $2")).
					WithArgs("test@gmail.com", "_9Z9sL~i3H4Kb33jcKJ9-8rZ+&-uRk#").
					WillReturnRows(rows)
			},
			args: args{
				email:    "test@gmail.com",
				password: "_9Z9sL~i3H4Kb33jcKJ9-8rZ+&-uRk#",
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := a.GetUser(testCase.args.email, testCase.args.password)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
