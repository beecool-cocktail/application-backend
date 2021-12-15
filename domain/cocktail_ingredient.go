package domain

import "time"

type CocktailIngredient struct {
	ID               int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	IngredientID     int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_ingredient_id"`
	CocktailID       int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_cocktail_id"`
	IngredientName   string    `gorm:"type:varchar(16) NOT NULL DEFAULT ''; comment: 成分名稱"`
	IngredientAmount float32   `gorm:"type:float NOT NULL DEFAULT 0; comment:成分數量"`
	IngredientUnit   string    `gorm:"type:varchar(16) NOT NULL DEFAULT ''; comment: 成分單位"`
	CreatedDate      time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}
