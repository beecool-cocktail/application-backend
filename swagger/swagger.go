package swagger

import "github.com/beecool-cocktail/application-backend/viewmodels"

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