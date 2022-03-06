package usercase

import (
	"github.com/beecool-cocktail/application-backend/domain"
)

func matchedByUserOfUpdateUserInfo(mockUserCache *domain.UserCache, mockUser *domain.User) bool {

	return mockUserCache.Name == mockUser.Name
}