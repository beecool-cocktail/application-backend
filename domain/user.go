package domain

import (
	"context"
	"time"
)

type User struct {
	ID                 int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	Account            string    `gorm:"type:varchar(20) NOT NULL;uniqueIndex:idx_account"`
	Password           string    `gorm:"type:varchar(100) NOT NULL DEFAULT ''"`
	Status             int       `gorm:"type:tinyint(1) NOT NULL DEFAULT 0"`
	Name               string    `gorm:"type:varchar(32) NOT NULL DEFAULT ''"`
	Email              string    `gorm:"type:varchar(64) NOT NULL DEFAULT ''"`
	Photo              string    `gorm:"type:varchar(128) NOT NULL"`
	NumberOfPost       int       `gorm:"type:int unsigned NOT NULL DEFAULT 0; comment: 貼文數"`
	NumberOfCollection int       `gorm:"type:int unsigned NOT NULL DEFAULT 0; comment: 收藏數"`
	IsCollectionPublic bool      `gorm:"type:tinyint(1) NOT NULL DEFAULT 0; comment: 公開收藏 0=不公開, 1=公開"`
	Remark             string    `gorm:"type:varchar(64) NOT NULL DEFAULT ''"`
	CreatedDate        time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type UserCache struct {
	Id           int64  `structs:"id"`
	Account      string `structs:"account"`
	Name         string `structs:"Name"`
	AccessToken  string `structs:"access_token"`
	RefreshToken string `structs:"refresh_token"`
	TokenExpire  string `structs:"token_expire"`
}

type UserMySQLRepository interface {
	Store(ctx context.Context, d *User) error
	QueryById(ctx context.Context, id int64) (*User, error)
}

type UserRedisRepository interface {
	Store(ctx context.Context, r *UserCache) error
	UpdateToken(ctx context.Context, r *UserCache) error
}

type UserUsecase interface {
	Logout(ctx context.Context, id int64) error
	QueryById(ctx context.Context, id int64) (*User, error)
}
