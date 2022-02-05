package mysql

import (
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) domain.DBTransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Transaction(txFunc func(interface{}) error) (err error) {
	tx := r.db.Begin()
	if !errors.Is(tx.Error, nil) {
		return tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if !errors.Is(err, nil) {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = txFunc(tx)
	return err
}