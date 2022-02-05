package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/domain/mocks"
	"github.com/beecool-cocktail/application-backend/testutil"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestUserHandler_GoogleAuthenticate(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockSocialAccountUsecase := new(mocks.SocialAccountUsecase)
	mockOauth2Token := &oauth2.Token{
		AccessToken: "Oauth2 token",
	}
	mockJWTToken := "JWT token"

	requestData := viewmodels.GoogleAuthenticateRequest{
		Code: "google-oauth-code",
	}

	r := testutil.GetRouteWithcontext()
	logger := testutil.GetLogger()

	handler := UserHandler{
		Logger:               logger,
		UserUsecase:          mockUserUsecase,
		SocialAccountUsecase: mockSocialAccountUsecase,
	}

	r.POST("/api/google-authenticate", handler.GoogleAuthenticate)

	t.Run("Success", func(t *testing.T) {
		var responseData viewmodels.ResponseData
		mockSocialAccountUsecase.
			On("Exchange",
				mock.Anything,
				mock.MatchedBy(func(code string) bool {
					return code == requestData.Code
				})).
			Return(mockOauth2Token, nil).Once()
		mockSocialAccountUsecase.
			On("GetUserInfo",
				mock.Anything,
				mock.MatchedBy(func(token *oauth2.Token) bool {
					return token == mockOauth2Token
				})).
			Return(mockJWTToken, nil).Once()

		requestJsonString, _ := json.Marshal(requestData)
		responseData.ErrorCode = domain.CodeSuccess
		responseData.Data = viewmodels.GoogleAuthenticateResponse{
			Token: mockJWTToken,
		}
		responseJsonString, _ := json.Marshal(responseData)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/google-authenticate", bytes.NewBuffer(requestJsonString))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, string(responseJsonString), w.Body.String())
	})
}

func TestUserHandler_Logout(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockSocialAccountUsecase := new(mocks.SocialAccountUsecase)

	requestData := viewmodels.LogoutRequest{
		UserID: 1,
	}

	r := testutil.GetRouteWithcontext()
	logger := testutil.GetLogger()

	handler := UserHandler{
		Logger:               logger,
		UserUsecase:          mockUserUsecase,
		SocialAccountUsecase: mockSocialAccountUsecase,
	}

	r.POST("/api/user/logout", handler.Logout)

	t.Run("Success", func(t *testing.T) {
		var responseData viewmodels.ResponseData
		mockUserUsecase.
			On("Logout",
				mock.Anything,
				mock.MatchedBy(func(userID int64) bool {
					return userID == requestData.UserID
				})).
			Return(nil).Once()

		requestJsonString, _ := json.Marshal(requestData)
		responseData.ErrorCode = domain.CodeSuccess
		responseData.Data = make([]interface{}, 0)
		responseJsonString, _ := json.Marshal(responseData)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/user/logout", bytes.NewBuffer(requestJsonString))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, string(responseJsonString), w.Body.String())
	})
}

func TestUserHandler_GetUserInfo(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockSocialAccountUsecase := new(mocks.SocialAccountUsecase)
	mockUser := domain.User{
		ID:                 testutil.UserID,
		Name:               "Andy",
		Email:              "abc123@gmail.com",
		Photo:              "static/images/image01.png",
		NumberOfPost:       2,
		NumberOfCollection: 110,
		IsCollectionPublic: true,
	}

	r := testutil.GetRouteWithcontext()
	logger := testutil.GetLogger()

	handler := UserHandler{
		Logger:               logger,
		UserUsecase:          mockUserUsecase,
		SocialAccountUsecase: mockSocialAccountUsecase,
	}

	r.POST("/api/user/info", handler.GetUserInfo)

	t.Run("Success", func(t *testing.T) {
		var responseData viewmodels.ResponseData
		mockUserUsecase.
			On("QueryById",
				mock.Anything,
				mock.MatchedBy(func(userID int64) bool {
					return userID == testutil.UserID
				})).
			Return(&mockUser, nil).Once()

		responseData.ErrorCode = domain.CodeSuccess
		responseData.Data = viewmodels.GetUserInfoResponse{
			UserID:             mockUser.ID,
			Name:               mockUser.Name,
			Email:              mockUser.Email,
			Photo:              mockUser.Photo,
			NumberOfPost:       mockUser.NumberOfPost,
			NumberOfCollection: mockUser.NumberOfCollection,
			IsCollectionPublic: mockUser.IsCollectionPublic,
		}
		responseJsonString, _ := json.Marshal(responseData)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/user/info", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, string(responseJsonString), w.Body.String())
	})
}
