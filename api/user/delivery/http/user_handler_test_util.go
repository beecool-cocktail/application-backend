package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/viewmodels"
)

func matchedByUpdateUserInfoTable(user *domain.User, requestData viewmodels.UpdateUserInfoRequest, id int64) bool {
	return user.ID == id &&
		user.Name == requestData.Name &&
		user.IsCollectionPublic == requestData.IsCollectionPublic
}