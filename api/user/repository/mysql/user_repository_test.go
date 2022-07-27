package mysql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/testutil"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_userMySQLRepository_Store(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := domain.User{
		ID:       1,
		Account:  "user1",
		Password: "password",
	}

	const sqlInsert = "INSERT INTO `users` " +
		"(`account`,`password`,`id`) " +
		"VALUES (?,?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).
		WithArgs(mockUser.Account, mockUser.Password, mockUser.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	d := NewMySQLUserRepository(db)

	_ = d.Store(context.TODO(), &mockUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_userMySQLRepository_QueryById(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := domain.User{
		ID: 1,
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(mockUser.ID)

	query := "SELECT * FROM `users` WHERE id = ? LIMIT 1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockUser.ID).WillReturnRows(rows)
		d := NewMySQLUserRepository(db)
		report, _ := d.QueryById(context.TODO(), mockUser.ID)
		assert.Equal(t, mockUser, report)
	})
}

func Test_userMySQLRepository_UpdateBasicInfo(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:                 1,
		Name:               "user-name",
		Length:             10,
		Width:              20,
		CoordinateX1:       10,
		CoordinateY1:       20,
		CoordinateX2:       30,
		CoordinateY2:       40,
		IsCollectionPublic: true,
	}

	sqlUpdate := "UPDATE `users` SET `coordinate_x1`=?,`coordinate_x2`=?,`coordinate_y1`=?,`coordinate_y2`=?," +
		"`is_collection_public`=?,`length`=?,`name`=?,`width`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.CoordinateX1, mockUser.CoordinateX2, mockUser.CoordinateY1, mockUser.CoordinateY2,
				mockUser.IsCollectionPublic, mockUser.Length, mockUser.Name, mockUser.Width, mockUser.ID).
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

	mockUserImage := &domain.UserImage{
		ID:          1,
		Destination: "images/my_photo.png",
	}

	sqlUpdate := "UPDATE `users` SET `photo`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUserImage.Destination, mockUserImage.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateImage(context.TODO(), mockUserImage)
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
		ID:                 1,
		Name:               "user-name",
		Length:             10,
		Width:              20,
		CoordinateX1:       10,
		CoordinateY1:       20,
		CoordinateX2:       30,
		CoordinateY2:       40,
		IsCollectionPublic: true,
	}

	sqlUpdate := "UPDATE `users` SET `coordinate_x1`=?,`coordinate_x2`=?,`coordinate_y1`=?,`coordinate_y2`=?," +
		"`is_collection_public`=?,`length`=?,`name`=?,`width`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.CoordinateX1, mockUser.CoordinateX2, mockUser.CoordinateY1, mockUser.CoordinateY2,
				mockUser.IsCollectionPublic, mockUser.Length, mockUser.Name, mockUser.Width, mockUser.ID).
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

	mockUserImage := &domain.UserImage{
		ID:          1,
		Destination: "images/my_photo.png",
	}

	sqlUpdate := "UPDATE `users` SET `photo`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUserImage.Destination, mockUserImage.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateImageTx(context.TODO(), db, mockUserImage)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_userMySQLRepository_UpdateNumberOfPostTx(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:           1,
		NumberOfPost: 2,
	}

	sqlUpdate := "UPDATE `users` SET `number_of_post`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.NumberOfPost, mockUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateNumberOfPostTx(context.TODO(), db, mockUser)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_userMySQLRepository_UpdateNumberOfDraftTx(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:            1,
		NumberOfDraft: 2,
	}

	sqlUpdate := "UPDATE `users` SET `number_of_draft`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.NumberOfDraft, mockUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateNumberOfDraftTx(context.TODO(), db, mockUser)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_userMySQLRepository_UpdateNumberOfNumberOfCollectionTx(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:                 1,
		NumberOfCollection: 2,
	}

	sqlUpdate := "UPDATE `users` SET `number_of_collection`=? WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(mockUser.NumberOfCollection, mockUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		d := NewMySQLUserRepository(db)
		rowsAffected, _ := d.UpdateNumberOfNumberOfCollectionTx(context.TODO(), db, mockUser)
		assert.Equal(t, int64(1), rowsAffected)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
