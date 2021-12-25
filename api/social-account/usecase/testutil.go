package usecase

import "github.com/beecool-cocktail/application-backend/domain"

func matchedByUserRedis(mockUser *domain.UserCache, user *domain.User) bool {
	return mockUser.Id == user.ID &&
		mockUser.Name == user.Name &&
		mockUser.Account == user.Account
}