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
