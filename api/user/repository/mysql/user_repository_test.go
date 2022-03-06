package mysql

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_userMySQLRepository_QueryById(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID: 1,
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(mockUser.ID)

	query := "SELECT * FROM `users` WHERE id = ? LIMIT 1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)
		d := NewMySQLUserRepository(db)
		client, _ := d.QueryById(context.TODO(), 1)
		assert.Equal(t, mockUser, client)
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(2).WillReturnError(gorm.ErrRecordNotFound)
		d := NewMySQLUserRepository(db)
		_, err = d.QueryById(context.TODO(), 2)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func Test_userMySQLRepository_Store(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:       1,
		Account:  "test_account",
		Password: "pass",
	}

	const sqlInsert = "INSERT INTO `users` (`account`,`password`,`id`) VALUES (?,?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).
		WithArgs(mockUser.Account, mockUser.Password, mockUser.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	d := NewMySQLUserRepository(db)

	d.Store(context.TODO(), mockUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_userMySQLRepository_UpdateBasicInfo(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:                 123456,
		Name:               "Andy",
		IsCollectionPublic: true,
	}

	sqlUpdate := "UPDATE `users` SET `is_collection_public`=?,`name`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.IsCollectionPublic, mockUser.Name, mockUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateBasicInfo(context.TODO(), mockUser)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_userMySQLRepository_UpdateImage(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.UserImage{
		ID:          123456,
		Destination: "static/images/",
	}

	sqlUpdate := "UPDATE `users` SET `photo`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.Destination, mockUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateImage(context.TODO(), mockUser)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_userMySQLRepository_UpdateBasicInfoTx(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:                 123456,
		Name:               "Andy",
		IsCollectionPublic: true,
	}

	sqlUpdate := "UPDATE `users` SET `is_collection_public`=?,`name`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.IsCollectionPublic, mockUser.Name, mockUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateBasicInfoTx(context.TODO(), db, mockUser)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_userMySQLRepository_UpdateImageTx(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.UserImage{
		ID:          123456,
		Destination: "static/images/",
	}

	sqlUpdate := "UPDATE `users` SET `photo`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.Destination, mockUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateImageTx(context.TODO(), db, mockUser)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
