package domain

import (
	"context"
	"time"
)

type Cocktail struct {
	ID          int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID  int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_cocktail_id"`
	Title       string    `gorm:"type:varchar(16) NOT NULL;; comment: 調酒名稱"`
	Description string    `gorm:"type:varchar(512) NOT NULL; comment: 調酒介紹"`
	UserID      int64     `gorm:"type:bigint(64) NOT NULL;index:idx_user_id; comment: 作者id"`
	UserName    string    `gorm:"type:varchar(32) NOT NULL DEFAULT '';comment: 作者名稱"`
	CreatedDate time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type CocktailMySQLRepository interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationMySQLRepository) (*[]Cocktail, int64, error)
}

type CocktailUsecase interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationUsecase) (*[]Cocktail, int64, error)
}
