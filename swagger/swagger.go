package swagger

import (
	"bytes"
	"github.com/beecool-cocktail/application-backend/viewmodels"
)


// swagger:parameters popularCocktailListRequest
type swaggerPopularCocktailListRequest struct {
	// in: body
	Body viewmodels.GetPopularCocktailListRequest
}

// swagger:response popularCocktailListResponse
type swaggerPopularCocktailListResponse struct {
	// in: body
	Body viewmodels.GetPopularCocktailListResponse
}


// swagger:parameters googleAuthenticateRequest
type swaggerGoogleAuthenticateRequest struct {
	// in: body
	Body viewmodels.GoogleAuthenticateRequest
}

// swagger:response googleAuthenticateResponse
type swaggerGoogleAuthenticateResponse struct {
	// in: body
	Body viewmodels.GoogleAuthenticateResponse
}

// swagger:parameters logoutRequest
type swaggerLogoutRequest struct {
	// in: body
	Body viewmodels.LogoutRequest
}

// swagger:response getUserInfoResponse
type swaggerGetUserInfoResponse struct {
	// in: body
	Body viewmodels.GetUserInfoResponse
}

// swagger:parameters updateUserInfoRequest
type swaggerUpdateUserInfoRequest struct {
	// in: body
	Body viewmodels.UpdateUserInfoRequest
}

// swagger:parameters updateUserPhoto
type swaggerUpdateUserPhoto struct {
	// in: formData
	//
	// swagger:file
	File *bytes.Buffer `json:"file"`
}

// swagger:response updateUserPhotoResponse
type swaggerUpdateUserPhotoResponse struct {
	// in: body
	Body viewmodels.UpdateUserPhotoResponse
}