package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/vincent-petithory/dataurl"
)

func matchedByUserOfUpdateUserInfo(mockUser *domain.User, requestData viewmodels.UpdateUserInfoRequest, userID int64) bool {

	return mockUser.ID == userID &&
		mockUser.Name == requestData.Name &&
		mockUser.IsCollectionPublic == requestData.IsCollectionPublic
}

func matchedByUserImageOfUpdateUserInfo(mockUserImage *domain.UserImage, dataUrl *dataurl.DataURL, userID int64) bool {

	return mockUserImage.ID == userID &&
		mockUserImage.Data == string(dataUrl.Data) &&
		mockUserImage.Type == dataUrl.MediaType.ContentType()
}