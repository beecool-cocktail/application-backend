package domain

import (
	"context"
	"golang.org/x/oauth2"
	"time"
)

type SocialAccount struct {
	ID          int64     `gorm:"type:bigint(64) NOT NULL auto_increment;primary_key"`
	SocialID    string    `gorm:"type:varchar(64) NOT NULL;uniqueIndex:idx_social_id"`
	UserID      int64     `gorm:"type:bigint(20) NOT NULL;uniqueIndex:idx_user_id"`
	Type        int       `gorm:"type:tinyint(1) NOT NULL DEFAULT 0"`
	CreatedDate time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_date"`
}

type GoogleUserInfo struct {
	Sub           string
	Email         string
	VerifiedEmail string
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
}

type State struct {
	RedirectPath      string `structs:"redirect_path"`
	CollectAfterLogin string `structs:"collect_after_login"`
}

type SocialAccountMySQLRepository interface {
	QueryById(ctx context.Context, id string) (*SocialAccount, error)
	Store(ctx context.Context, s *SocialAccount, u *User) (int64, error)
}

type SocialAccountRedisRepository interface {
	StoreState(ctx context.Context, key string, state State) error
	GetState(ctx context.Context, key string) (State, error)
	DeleteState(ctx context.Context, key string) error
}

type SocialAccountGoogleOAuthRepository interface {
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error)
}

type SocialAccountUsecase interface {
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (string, error)
	GenerateState(ctx context.Context, state State) (string, error)
	GetState(ctx context.Context, stateString string) (State, error)
}
