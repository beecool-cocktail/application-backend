package mysql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beecool-cocktail/application-backend/domain"
	social_account "github.com/beecool-cocktail/application-backend/enum/social-account"
	"github.com/beecool-cocktail/application-backend/testutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Test_socialAccountMySQLRepository_Store(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		Account: "account",
		Name:    "Andy",
		Email:   "abc123@gmail.com",
	}

	mockSocialAccount := &domain.SocialAccount{
		SocialID: "googleUUID",
		UserID:   1,
		Type:     social_account.ParseSocialAccountType(social_account.Google),
	}

	const sqlInsertSocialAccounts = "INSERT INTO `social_accounts` (`social_id`,`user_id`,`type`) VALUES (?,?,?)"
	const sqlInsertUsers = "INSERT INTO `users` (`account`,`name`,`email`) VALUES (?,?,?)"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlInsertUsers)).
			WithArgs(mockUser.Account, mockUser.Name, mockUser.Email).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(sqlInsertSocialAccounts)).
			WithArgs(mockSocialAccount.SocialID, mockSocialAccount.UserID, mockSocialAccount.Type).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		d := NewMySQLSocialAccountRepository(db)

		d.Store(context.TODO(), mockSocialAccount, mockUser)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Failed and Rollback", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlInsertUsers)).
			WithArgs(mockUser.Account, mockUser.Name, mockUser.Email).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		d := NewMySQLSocialAccountRepository(db)

		d.Store(context.TODO(), mockSocialAccount, mockUser)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_socialAccountMySQLRepository_QueryById(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockClient := &domain.SocialAccount{
		SocialID:      "googleUUID",
	}

	rows := sqlmock.NewRows([]string{"social_id"}).
		AddRow(mockClient.SocialID)

	query := "SELECT * FROM `social_accounts` WHERE social_id = ? LIMIT 1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("googleUUID").WillReturnRows(rows)
		d := NewMySQLSocialAccountRepository(db)
		client, _ := d.QueryById(context.TODO(), "googleUUID")
		assert.Equal(t, mockClient, client)
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("noGoogleUUID").WillReturnError(gorm.ErrRecordNotFound)
		d := NewMySQLSocialAccountRepository(db)
		_, err = d.QueryById(context.TODO(), "noGoogleUUID")
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
