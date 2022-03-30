package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type CocktailStep struct {
	ID              int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID      int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_cocktail_id"`
	StepNumber      int       `gorm:"type:int(2) unsigned NOT NULL DEFAULT 1; comment: 步驟"`
	StepDescription string    `gorm:"type:varchar(64) NOT NULL DEFAULT ''; comment: 步驟介紹"`
	CreatedDate     time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type CocktailStepMySQLRepository interface {
	StoreTx(ctx context.Context, tx *gorm.DB, c *CocktailStep) error
	QueryByCocktailId(ctx context.Context, id int64) ([]CocktailStep, error)
	DeleteByCocktailIDTx(ctx context.Context, tx *gorm.DB, id int64) error
	UpdateTx(ctx context.Context, tx *gorm.DB, c *CocktailStep) (int64, error)
}

type CocktailStepUsecase interface {
	Store(ctx context.Context, c *CocktailStep) error
}
