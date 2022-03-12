package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type CocktailPhoto struct {
	ID           int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID   int64     `gorm:"type:bigint(64) NOT NULL;index:idx_cocktail_id"`
	Photo        string    `gorm:"type:varchar(128) NOT NULL; comment: 調酒圖片"`
	IsCoverPhoto bool      `gorm:"type:tinyint(1) NOT NULL DEFAULT 0; comment: 是否為封面照 0=no, 1=yes"`
	CreatedDate  time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type CocktailPhotoMySQLRepository interface {
	StoreTx(ctx context.Context, tx *gorm.DB, c *CocktailPhoto) error
	QueryCoverPhotoByCocktailId(ctx context.Context, id int64) (string, error)
	QueryPhotosByCocktailId(ctx context.Context, id int64) ([]string, error)
}

type CocktailPhotoUsecase interface {
	Store(ctx context.Context, c *CocktailPhoto) error
}
