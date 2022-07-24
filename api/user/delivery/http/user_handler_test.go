package http

import (
	"bytes"
	"encoding/json"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/domain/mocks"
	"github.com/beecool-cocktail/application-backend/testutil"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_GoogleAuthenticate(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockFavoriteCocktailUsecase := new(mocks.FavoriteCocktailUsecase)
	mockCocktailUsecase := new(mocks.CocktailUsecase)
	mockSocialAccountUsecase := new(mocks.SocialAccountUsecase)

	r := testutil.GetRouter()
	svc := testutil.GetService()
	logger := testutil.GetLogger()
	handler := UserHandler{
		Configure:               svc.Configure,
		Logger:                  logger,
		UserUsecase:             mockUserUsecase,
		FavoriteCocktailUsecase: mockFavoriteCocktailUsecase,
		CocktailUsecase:         mockCocktailUsecase,
		SocialAccountUsecase:    mockSocialAccountUsecase,
	}

	r.POST("/api/auth/google-authenticate", handler.GoogleAuthenticate)

	t.Run("Success", func(t *testing.T) {
		mockGoogleToken := &oauth2.Token{
			AccessToken:  "access_token",
			TokenType:    "token_type",
			RefreshToken: "refresh_token",
			Expiry:       time.Time{},
		}
		mockJWTToken := "jwt_token"

		requestData := viewmodels.GoogleAuthenticateRequest{
			Code: "code",
		}

		mockSocialAccountUsecase.
			On("Exchange",
				mock.Anything,
				requestData.Code,
			).Return(mockGoogleToken, nil).Once()

		mockSocialAccountUsecase.
			On("GetUserInfo",
				mock.Anything,
				mockGoogleToken,
			).Return(mockJWTToken, nil).Once()

		requestJsonString, _ := json.Marshal(requestData)

		var responseData viewmodels.ResponseData
		responseData.Data = viewmodels.GoogleAuthenticateResponse{
			Token: mockJWTToken,
		}
		responseData.ErrorCode = domain.CodeSuccess
		responseJsonString, _ := json.Marshal(responseData)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/auth/google-authenticate", bytes.NewBuffer(requestJsonString))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, string(responseJsonString), w.Body.String())

		mockSocialAccountUsecase.AssertExpectations(t)
	})
}

func TestNewUserHandler(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockFavoriteCocktailUsecase := new(mocks.FavoriteCocktailUsecase)
	mockCocktailUsecase := new(mocks.CocktailUsecase)
	mockSocialAccountUsecase := new(mocks.SocialAccountUsecase)

	r := testutil.GetRouter()
	svc := testutil.GetService()
	logger := testutil.GetLogger()
	handler := UserHandler{
		Configure:               svc.Configure,
		Logger:                  logger,
		UserUsecase:             mockUserUsecase,
		FavoriteCocktailUsecase: mockFavoriteCocktailUsecase,
		CocktailUsecase:         mockCocktailUsecase,
		SocialAccountUsecase:    mockSocialAccountUsecase,
	}

	r.POST("/api/auth/logout", handler.Logout)

	t.Run("Success", func(t *testing.T) {

		requestData := viewmodels.LogoutRequest{
			UserID: 123456,
		}

		mockUserUsecase.
			On("Logout",
				mock.Anything,
				requestData.UserID,
			).Return(nil).Once()

		requestJsonString, _ := json.Marshal(requestData)

		var responseData viewmodels.ResponseData
		responseData.Data = make([]interface{}, 0)
		responseData.ErrorCode = domain.CodeSuccess
		responseJsonString, _ := json.Marshal(responseData)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/auth/logout", bytes.NewBuffer(requestJsonString))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, string(responseJsonString), w.Body.String())

		mockUserUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_GetUserInfo(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockFavoriteCocktailUsecase := new(mocks.FavoriteCocktailUsecase)
	mockCocktailUsecase := new(mocks.CocktailUsecase)
	mockSocialAccountUsecase := new(mocks.SocialAccountUsecase)

	r := testutil.GetRouteWithcontext()
	svc := testutil.GetService()
	logger := testutil.GetLogger()
	handler := UserHandler{
		Configure:               svc.Configure,
		Logger:                  logger,
		UserUsecase:             mockUserUsecase,
		FavoriteCocktailUsecase: mockFavoriteCocktailUsecase,
		CocktailUsecase:         mockCocktailUsecase,
		SocialAccountUsecase:    mockSocialAccountUsecase,
	}

	r.POST("/api/user/info", handler.GetUserInfo)

	t.Run("Success", func(t *testing.T) {

		mockUser := domain.User{
			ID:                 1,
			Account:            "account",
			Password:           "password",
			Status:             0,
			Name:               "name",
			Email:              "email",
			Photo:              "photo",
			NumberOfPost:       1,
			NumberOfCollection: 2,
			NumberOfDraft:      3,
			IsCollectionPublic: true,
			Remark:             "good",
			CreatedDate:        time.Time{},
		}

		mockUserID := int64(1)
		mockNumberOfPost := int64(10)
		mockNumberOfCollection := int64(20)
		mockUserUsecase.
			On("QueryById",
				mock.Anything,
				mockUserID,
			).Return(mockUser, nil).Once()

		mockCocktailUsecase.
			On("QueryFormalCountsByUserID",
				mock.Anything,
				mockUser.ID,
			).Return(mockNumberOfPost, nil).Once()

		mockFavoriteCocktailUsecase.
			On("QueryCountsByUserID",
				mock.Anything,
				mockUser.ID,
			).Return(mockNumberOfCollection, nil).Once()

		var responseData viewmodels.ResponseData
		responseData.Data = viewmodels.GetUserInfoResponse{
			UserID:             mockUser.ID,
			Name:               mockUser.Name,
			Email:              mockUser.Email,
			Photo:              mockUser.Photo,
			NumberOfPost:       mockNumberOfPost,
			NumberOfCollection: mockNumberOfCollection,
			IsCollectionPublic: mockUser.IsCollectionPublic,
		}
		responseData.ErrorCode = domain.CodeSuccess
		responseJsonString, _ := json.Marshal(responseData)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/api/user/info", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, string(responseJsonString), w.Body.String())

		mockUserUsecase.AssertExpectations(t)
		mockCocktailUsecase.AssertExpectations(t)
		mockFavoriteCocktailUsecase.AssertExpectations(t)
	})
}
