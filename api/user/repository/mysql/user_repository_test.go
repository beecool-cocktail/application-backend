package mysql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beecool-cocktail/application-backend/testutil"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"

	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
)

func Test_userMySQLRepository_QueryById(t *testing.T) {
	db, mock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUser := &domain.User{
		ID:     1,
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
