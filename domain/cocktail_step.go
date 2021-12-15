package domain

import "time"

type CocktailStep struct {
	ID              int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	StepID          int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_recipe_id"`
	CocktailID      int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_cocktail_id"`
	StepNumber      int       `gorm:"type:int(2) unsigned NOT NULL DEFAULT 1; comment: 步驟"`
	StepDescription string    `gorm:"type:varchar(64) NOT NULL DEFAULT ''; comment: 步驟介紹"`
	CreatedDate     time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}
