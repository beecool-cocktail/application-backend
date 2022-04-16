package swagger

import (
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
	Body struct {
		ErrorCode    string                                    `json:"error_code"`
		ErrorMessage string                                    `json:"error_message"`
		Data         viewmodels.GetPopularCocktailListResponse `json:"data"`
	}
}

// swagger:parameters getDraftCocktailListRequest
type swaggerGetDraftCocktailListRequest struct {
	// in: body
	Body viewmodels.GetDraftCocktailListRequest
}

// swagger:response getDraftCocktailListResponse
type swaggerGetDraftCocktailListResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                                  `json:"error_code"`
		ErrorMessage string                                  `json:"error_message"`
		Data         viewmodels.GetDraftCocktailListResponse `json:"data"`
	}
}

// swagger:parameters googleAuthenticateRequest
type swaggerGoogleAuthenticateRequest struct {
	// in: body
	Body viewmodels.GoogleAuthenticateRequest
}

// swagger:response googleAuthenticateResponse
type swaggerGoogleAuthenticateResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                                `json:"error_code"`
		ErrorMessage string                                `json:"error_message"`
		Data         viewmodels.GoogleAuthenticateResponse `json:"data"`
	}
}

// swagger:parameters logoutRequest
type swaggerLogoutRequest struct {
	// in: body
	Body viewmodels.LogoutRequest
}

// swagger:response getUserInfoResponse
type swaggerGetUserInfoResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                         `json:"error_code"`
		ErrorMessage string                         `json:"error_message"`
		Data         viewmodels.GetUserInfoResponse `json:"data"`
	}
}

// swagger:parameters updateUserInfoRequest
type swaggerUpdateUserInfoRequest struct {
	// in: body
	Body viewmodels.UpdateUserInfoRequest
}

// swagger:response updateUserPhotoResponse
type swaggerUpdateUserPhotoResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                            `json:"error_code"`
		ErrorMessage string                            `json:"error_message"`
		Data         viewmodels.UpdateUserInfoResponse `json:"data"`
	}
}

// swagger:parameters postArticleRequest
type swaggerPostArticleRequest struct {
	// in: body
	Body viewmodels.PostArticleRequest
}

// swagger:parameters postDraftArticleRequest
type swaggerPostDraftArticleRequest struct {
	// in: body
	Body viewmodels.PostDraftArticleRequest
}

// swagger:parameters getCocktailByIDRequest
type swaggerGetCocktailByIDRequest struct {
	// in: body
	Body viewmodels.GetCocktailByIDRequest
}

// swagger:response getCocktailByIDResponse
type swaggerGetCocktailByIDResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                             `json:"error_code"`
		ErrorMessage string                             `json:"error_message"`
		Data         viewmodels.GetCocktailByIDResponse `json:"data"`
	}
}

// swagger:parameters getCocktailDraftByIDRequest
type swaggerGetCocktailDraftByIDRequest struct {
	// in: body
	Body viewmodels.GetCocktailDraftByIDRequest
}

// swagger:response getCocktailDraftByIDResponse
type swaggerGetCocktailDraftByIDResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                                  `json:"error_code"`
		ErrorMessage string                                  `json:"error_message"`
		Data         viewmodels.GetCocktailDraftByIDResponse `json:"data"`
	}
}

// swagger:model updateDraftArticleRequest
type swaggerUpdateDraftArticleRequest struct {
	// in: body
	viewmodels.UpdateDraftArticleRequest
}

// swagger:parameters deleteDraftArticleRequest
type swaggerDeleteDraftArticleRequest struct {
	// in: body
	Body viewmodels.DeleteDraftArticleRequest
}

// swagger:parameters collectArticleRequest
type swaggerCollectArticleRequest struct {
	// in: body
	Body viewmodels.CollectArticleRequest
}

// swagger:response getUserFavoriteCocktailListResponse
type swaggerGetUserFavoriteCocktailListResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                                         `json:"error_code"`
		ErrorMessage string                                         `json:"error_message"`
		Data         viewmodels.GetUserFavoriteCocktailListResponse `json:"data"`
	}
}

// swagger:response getSelfCocktailListResponse
type swaggerGetSelfCocktailListResponse struct {
	// in: body
	Body struct {
		ErrorCode    string                                 `json:"error_code"`
		ErrorMessage string                                 `json:"error_message"`
		Data         viewmodels.GetSelfCocktailListResponse `json:"data"`
	}
}
