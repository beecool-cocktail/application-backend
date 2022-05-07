package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type CocktailImage struct {
	ImageID      int64
	CocktailID   int64
	Data         string
	Name         string
	Type         string
	Destination  string
	IsCoverPhoto bool
	IsLowQuality bool
}

type APICocktail struct {
	CocktailID       int64
	UserID           int64
	UserName         string
	Title            string
	Description      string
	CoverPhoto       CocktailPhoto
	Photos           []CocktailPhoto
	LowQualityPhotos []CocktailPhoto
	Ingredients      []CocktailIngredient
	Steps            []CocktailStep
	IsCollected      bool
	CreatedDate      string
}

type CocktailRedis struct {
	CocktailID       int64 `structs:"cocktail_id"`
	CollectionCounts int   `structs:"collection_counts"`
}

type Cocktail struct {
	ID          int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID  int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_cocktail_id"`
	UserID      int64     `gorm:"type:bigint(64) NOT NULL;index:idx_user_id; comment: 作者id"`
	Title       string    `gorm:"type:varchar(16) NOT NULL;; comment: 調酒名稱"`
	Description string    `gorm:"type:varchar(512) NOT NULL; comment: 調酒介紹"`
	Category    int       `gorm:"type:tinyint(1) NOT NULL DEFAULT 0; comment: 類型 0=草稿, 1=正式"`
	CreatedDate time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type CocktailMySQLRepository interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationMySQLRepository) ([]Cocktail, int64, error)
	QueryByCocktailID(ctx context.Context, id int64) (Cocktail, error)
	QueryFormalByUserID(ctx context.Context, id int64) ([]Cocktail, error)
	StoreTx(ctx context.Context, tx *gorm.DB, c *Cocktail) error
	DeleteTx(ctx context.Context, tx *gorm.DB, id int64) error
	UpdateTx(ctx context.Context, tx *gorm.DB, c *Cocktail) (int64, error)
	UpdateCategoryTx(ctx context.Context, tx *gorm.DB, c *Cocktail) (int64, error)
}

type CocktailUsecase interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationUsecase, userID int64) ([]APICocktail, int64, error)
	QueryByCocktailID(ctx context.Context, cocktailID, userID int64) (APICocktail, error)
	QueryFormalByUserID(ctx context.Context, id int64) ([]APICocktail, error)
	QueryDraftByCocktailID(ctx context.Context, cocktailID, userID int64) (APICocktail, error)
	Store(ctx context.Context, c *Cocktail, cig []CocktailIngredient, cs []CocktailStep, ci []CocktailImage, userID int64) error
	Delete(ctx context.Context, cocktailID, userID int64) error
	Update(ctx context.Context, c *Cocktail, cig []CocktailIngredient, cs []CocktailStep, ci []CocktailImage, userID int64) error
	MakeDraftToFormal(ctx context.Context, cocktailID, userID int64) error
}

type CocktailFileRepository interface {
	SaveAsWebp(ctx context.Context, ci *CocktailImage) error
	SaveAsWebpInLQIP(ctx context.Context, ci *CocktailImage) error
	UpdateAsWebp(ctx context.Context, ci *CocktailImage) error
	UpdateAsWebpInLQIP(ctx context.Context, ci *CocktailImage) error
}
