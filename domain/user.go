package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                 int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	Account            string    `gorm:"type:varchar(20) NOT NULL;uniqueIndex:idx_account"`
	Password           string    `gorm:"type:varchar(100) NOT NULL DEFAULT ''"`
	Status             int       `gorm:"type:tinyint(1) NOT NULL DEFAULT 0"`
	Type               int       `gorm:"type:tinyint(1) NOT NULL DEFAULT 1; comment: 帳號類型 0=測試, 1=正式"`
	Name               string    `gorm:"type:varchar(32) NOT NULL DEFAULT ''"`
	Email              string    `gorm:"type:varchar(64) NOT NULL DEFAULT ''"`
	OriginAvatar       string    `gorm:"type:varchar(128) NOT NULL"`
	CropAvatar         string    `gorm:"type:varchar(128) NOT NULL"`
	Height             int       `gorm:"type:int NOT NULL DEFAULT 0; comment:照片長度"`
	Width              int       `gorm:"type:int NOT NULL DEFAULT 0; comment:照片寬度"`
	CoordinateX1       float32   `gorm:"type:float NOT NULL DEFAULT 0; comment:照片左上X座標"`
	CoordinateY1       float32   `gorm:"type:float NOT NULL DEFAULT 0; comment:照片左上Y座標"`
	CoordinateX2       float32   `gorm:"type:float NOT NULL DEFAULT 0; comment:照片左下X座標"`
	CoordinateY2       float32   `gorm:"type:float NOT NULL DEFAULT 0; comment:照片左下Y座標"`
	Rotation           float32   `gorm:"type:float NOT NULL DEFAULT 0; comment:照片選轉角度"`
	NumberOfPost       int       `gorm:"type:int unsigned NOT NULL DEFAULT 0; comment: 貼文數"`
	NumberOfCollection int       `gorm:"type:int unsigned NOT NULL DEFAULT 0; comment: 收藏數"`
	NumberOfDraft      int       `gorm:"type:int unsigned NOT NULL DEFAULT 0; comment: 草稿數"`
	IsCollectionPublic bool      `gorm:"type:tinyint(1) NOT NULL DEFAULT 0; comment: 公開收藏 0=不公開, 1=公開"`
	Remark             string    `gorm:"type:varchar(64) NOT NULL DEFAULT ''"`
	CreatedDate        time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type UserCache struct {
	Id           int64  `structs:"id"`
	Account      string `structs:"account"`
	Name         string `structs:"name"`
	AccessToken  string `structs:"access_token"`
	RefreshToken string `structs:"refresh_token"`
	TokenExpire  string `structs:"token_expire"`
}

type UserAvatar struct {
	UserID       int64
	OriginAvatar OriginAvatar
	CropAvatar   CropAvatar
}

type OriginAvatar struct {
	DataURL     string
	Type        string
	Destination string
}

type CropAvatar struct {
	DataURL     string
	Type        string
	Destination string
}

type UserMySQLRepository interface {
	Store(ctx context.Context, d *User) error
	QueryById(ctx context.Context, id int64) (User, error)
	UpdateUserOriginAvatarTx(ctx context.Context, tx *gorm.DB, ui *UserAvatar) (int64, error)
	UpdateUserCropAvatarTx(ctx context.Context, tx *gorm.DB, ui *UserAvatar) (int64, error)
	UpdateUserAvatarInfoTx(ctx context.Context, tx *gorm.DB, d *User) (int64, error)
	UpdateUserNameTx(ctx context.Context, tx *gorm.DB, d *User) (int64, error)
	UpdateUserCollectionStatus(ctx context.Context, d *User) (int64, error)
	UpdateNumberOfPostTx(ctx context.Context, tx *gorm.DB, d *User) (int64, error)
	UpdateNumberOfDraftTx(ctx context.Context, tx *gorm.DB, d *User) (int64, error)
	UpdateNumberOfNumberOfCollectionTx(ctx context.Context, tx *gorm.DB, d *User) (int64, error)
}

type UserRedisRepository interface {
	Store(ctx context.Context, r *UserCache) error
	UpdateToken(ctx context.Context, r *UserCache) error
	UpdateBasicInfo(ctx context.Context, r *UserCache) error
	QueryUserNameByID(ctx context.Context, id int64) (string, error)
}

type UserFileRepository interface {
	SaveOriginAvatarAsWebp(ctx context.Context, ui *OriginAvatar) (int, int, error)
	SaveCropAvatarAsWebp(ctx context.Context, ui *CropAvatar) error
}

type UserUsecase interface {
	Logout(ctx context.Context, id int64) error
	QueryById(ctx context.Context, id int64) (User, error)
	UpdateUserAvatar(ctx context.Context, d *User, ui *UserAvatar) error
	UpdateUserName(ctx context.Context, d *User) error
	UpdateUserCollectionStatus(ctx context.Context, d *User) error
	DeleteUserAvatar(ctx context.Context, userID int64) error
}
