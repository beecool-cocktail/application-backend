package viewmodels

type GoogleAuthenticateRequest struct {
	Code string `json:"code"`
}

type GoogleAuthenticateResponse struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	UserID int64 `json:"user_id"`
}

type GetUserInfoResponse struct {
	UserID             int64  `json:"user_id"`
	Name               string `json:"user_name"`
	Email              string `json:"email"`
	Photo              string `json:"photo"`
	NumberOfPost       int    `json:"number_of_post"`
	NumberOfCollection int    `json:"number_of_collection"`
	// 是否公開收藏 false=不公開, true=公開
	IsCollectionPublic bool `json:"is_collection_public"`
}

type UpdateUserInfoRequest struct {
	// name for this user
	// required: true
	// example: Feen Lin
	Name               string `json:"name" binding:"required"`

	// public user collection post
	// required: true
	// example: false
	IsCollectionPublic bool   `json:"is_collection_public"`
}

type UpdateUserPhotoResponse struct {
	Photo string `json:"photo"`
}