package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type APIFavoriteCocktail struct {
	CocktailID  int64
	UserID      int64
	UserName    string
	Title       string
	CoverPhoto  string
	IsCollected bool
}

type FavoriteCocktail struct {
	ID          int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID  int64     `gorm:"type:bigint(64) NOT NULL;index:idx_favorite_cocktail,priority:2; comment: 調酒id"`
	UserID      int64     `gorm:"type:bigint(64) NOT NULL;index:idx_favorite_cocktail,priority:1; comment: 收藏者id"`
	CreatedDate time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type FavoriteCocktailMySQLRepository interface {
	StoreTx(ctx context.Context, tx *gorm.DB, c *FavoriteCocktail) error
	QueryByUserID(ctx context.Context, id int64, pagination PaginationMySQLRepository) ([]FavoriteCocktail, int64, error)
	QueryCountsByUserID(ctx context.Context, id int64) (int64, error)
	Delete(ctx context.Context, cocktailID, userID int64) error
	DeleteTx(ctx context.Context, tx *gorm.DB, cocktailID, userID int64) error
}

type FavoriteCocktailUsecase interface {
	Store(ctx context.Context, c *FavoriteCocktail) error
	QueryByUserID(ctx context.Context, targetUserID int64, pagination PaginationUsecase,
		needCollectedStatusUserID int64) ([]APIFavoriteCocktail, int64, error)
	QueryCountsByUserID(ctx context.Context, id int64) (int64, error)
	Delete(ctx context.Context, cocktailID, userID int64) (string, error)
}
