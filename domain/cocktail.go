package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type CocktailImage struct {
	ID           int64
	Data         string
	Name         string
	Type         string
	Destination  string
	IsCoverPhoto bool
}

type APICocktail struct {
	CocktailID  int64
	UserID      int64
	Title       string
	Description string
	CoverPhoto  string
	Photos      []string
	Ingredients []CocktailIngredient
	Steps       []CocktailStep
	CreatedDate string
}

type Cocktail struct {
	ID          int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID  int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_cocktail_id"`
	UserID      int64     `gorm:"type:bigint(64) NOT NULL;index:idx_user_id; comment: 作者id"`
	Title       string    `gorm:"type:varchar(16) NOT NULL;; comment: 調酒名稱"`
	Description string    `gorm:"type:varchar(512) NOT NULL; comment: 調酒介紹"`
	Category    int      `gorm:"type:tinyint(1) NOT NULL DEFAULT 0; comment: 公開收藏 0=草稿, 1=正式"`
	CreatedDate time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type CocktailMySQLRepository interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationMySQLRepository) ([]Cocktail, int64, error)
	QueryByCocktailID(ctx context.Context, id int64) (Cocktail, error)
	StoreTx(ctx context.Context, tx *gorm.DB, c *Cocktail) error
}

type CocktailUsecase interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationUsecase) ([]APICocktail, int64, error)
	QueryByCocktailID(ctx context.Context, id int64) (APICocktail, error)
	QueryDraftByCocktailID(ctx context.Context, cocktailID, userID int64) (APICocktail, error)
	Store(ctx context.Context, c *Cocktail, cig []CocktailIngredient, cs []CocktailStep, ci []CocktailImage) error
}

type CocktailFileRepository interface {
	SaveAsWebp(ctx context.Context, ci *CocktailImage) error
}
