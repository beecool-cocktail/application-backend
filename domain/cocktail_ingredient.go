package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type CocktailIngredient struct {
	ID               int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID       int64     `gorm:"type:bigint(64) NOT NULL;index:idx_cocktail_id"`
	IngredientName   string    `gorm:"type:varchar(16) NOT NULL DEFAULT ''; comment: 成分名稱"`
	IngredientAmount float32   `gorm:"type:float NOT NULL DEFAULT 0; comment:成分數量"`
	IngredientUnit   string    `gorm:"type:varchar(16) NOT NULL DEFAULT ''; comment: 成分單位"`
	CreatedDate      time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type CocktailIngredientMySQLRepository interface {
	StoreTx(ctx context.Context, tx *gorm.DB, c *CocktailIngredient) error
	QueryByCocktailId(ctx context.Context, id int64) ([]CocktailIngredient, error)
}

type CocktailIngredientUsecase interface {
	Store(ctx context.Context, c *CocktailIngredient) error
}