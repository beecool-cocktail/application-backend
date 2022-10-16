package domain

import (
	"context"
	"github.com/bsm/redislock"
	"gorm.io/gorm"
	"time"
)

const CocktailsIndex = "cocktails"
const CocktailsMapping = `
{
	"mappings": {
		"properties": {
			"cocktail_id": {
				"type": "keyword"
			},
			"user_id": {
				"type": "keyword"
			},
			"title": {
				"type": "text",
				"analyzer": "ik_max_word",
            	"search_analyzer": "ik_max_word"
			},
			"description": {
				"type": "text",
				"analyzer": "ik_max_word",
            	"search_analyzer": "ik_max_word"
			},
			"ingredients": {
				"type": "text",
				"analyzer": "ik_max_word",
            	"search_analyzer": "ik_max_word"
			},
			"steps": {
				"type": "text",
				"analyzer": "ik_max_word",
            	"search_analyzer": "ik_max_word"
			},
			"created_date": {
				"type": "date"
			}
		}
	}
}`

type CocktailElasticSearch struct {
	CocktailID  int64     `json:"cocktail_id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Steps       []string  `json:"steps"`
	Ingredients []string  `json:"ingredients"`
	CreatedDate time.Time `json:"created_date"`
}

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

type CocktailCollection struct {
	CocktailID       int64   `structs:"cocktail_id"`
	CollectionCounts float64 `structs:"collection_counts"`
}

type Cocktail struct {
	ID                 int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	CocktailID         int64     `gorm:"type:bigint(64) NOT NULL;uniqueIndex:idx_cocktail_id"`
	UserID             int64     `gorm:"type:bigint(64) NOT NULL;index:idx_user_id; comment: 作者id"`
	Title              string    `gorm:"type:varchar(30) NOT NULL;; comment: 調酒名稱"`
	Description        string    `gorm:"type:varchar(512) NOT NULL; comment: 調酒介紹"`
	Category           int       `gorm:"type:tinyint(1) NOT NULL DEFAULT 0; comment: 類型 0=草稿, 1=正式"`
	NumberOfCollection int       `gorm:"type:int unsigned NOT NULL DEFAULT 0; comment: 收藏數"`
	CreatedDate        time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type CocktailMySQLRepository interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationMySQLRepository) ([]Cocktail, int64, error)
	QueryByCocktailID(ctx context.Context, id int64) (Cocktail, error)
	QueryFormalByUserID(ctx context.Context, id int64) ([]Cocktail, error)
	QueryFormalCountsByUserID(ctx context.Context, id int64) (int64, error)
	StoreTx(ctx context.Context, tx *gorm.DB, c *Cocktail) error
	DeleteTx(ctx context.Context, tx *gorm.DB, id int64) error
	UpdateTx(ctx context.Context, tx *gorm.DB, c *Cocktail) (int64, error)
	UpdateCategoryTx(ctx context.Context, tx *gorm.DB, c *Cocktail) (int64, error)
	IncreaseNumberOfCollectionTx(ctx context.Context, tx *gorm.DB, cocktailID int64) (int64, error)
	DecreaseNumberOfCollectionTx(ctx context.Context, tx *gorm.DB, cocktailID int64) (int64, error)
}

type CocktailRedisRepository interface {
	GetCocktailCollectionNumberLock(ctx context.Context, key string, ttl, retryInterval time.Duration,
		retryTimes int) (*redislock.Lock, error)
	ReleaseCocktailCollectionNumberLock(ctx context.Context, lock *redislock.Lock) error
}

type CocktailFileRepository interface {
	SaveAsWebp(ctx context.Context, ci *CocktailImage) error
	SaveAsWebpInLQIP(ctx context.Context, ci *CocktailImage) error
	UpdateAsWebp(ctx context.Context, ci *CocktailImage) error
	UpdateAsWebpInLQIP(ctx context.Context, ci *CocktailImage) error
}

type CocktailElasticSearchRepository interface {
	Index(ctx context.Context, c *CocktailElasticSearch) error
	Search(ctx context.Context, text string, from, size int) ([]CocktailElasticSearch, int64, error)
	Update(ctx context.Context, c *CocktailElasticSearch) error
	Delete(ctx context.Context, id int64) error
}

type CocktailUsecase interface {
	GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination PaginationUsecase, needCollectedStatusUserID int64) ([]APICocktail, int64, error)
	Search(ctx context.Context, keyword string, from, size int, userID int64) ([]APICocktail, int64, error)
	QueryByCocktailID(ctx context.Context, cocktailID, needCollectedStatusUserID int64) (APICocktail, error)
	QueryFormalByUserID(ctx context.Context, targetUserID int64, needCollectedStatusUserID int64) ([]APICocktail, error)
	QueryDraftByCocktailID(ctx context.Context, cocktailID, userID int64) (APICocktail, error)
	QueryFormalCountsByUserID(ctx context.Context, id int64) (int64, error)
	Store(ctx context.Context, c *Cocktail, cig []CocktailIngredient, cs []CocktailStep, ci []CocktailImage, userID int64) error
	Delete(ctx context.Context, cocktailID, userID int64) error
	Update(ctx context.Context, c *Cocktail, cig []CocktailIngredient, cs []CocktailStep, ci []CocktailImage, userID int64) error
	MakeDraftToFormal(ctx context.Context, cocktailID, userID int64) error
}
