package viewmodels

type GoogleLoginRequest struct {
	RedirectPath      string `form:"redirect_path" json:"redirect_path"`
	CollectAfterLogin string `form:"collect_after_login" json:"collect_after_login"`
}

type GoogleAuthenticateRequest struct {
	//required: true
	Code  string `json:"code"`
	State string `json:"state"`
}

type GoogleAuthenticateResponse struct {
	//required: true
	Token             string `json:"token"`
	RedirectPath      string `json:"redirect_path"`
	CollectAfterLogin string `json:"collect_after_login"`
}

type LogoutRequest struct {
	//required: true
	UserID int64 `json:"user_id"`
}

type GetUserInfoResponse struct {
	//required: true
	UserID int64 `json:"user_id"`
	//required: true
	Name string `json:"user_name"`
	//required: true
	Email string `json:"email"`
	//required: true
	OriginAvatar string `json:"origin_avatar"`
	//required: true
	CropAvatar string `json:"crop_avatar"`
	//原圖長
	// required: true
	Height int `json:"height" binding:"required"`
	//原圖寬
	// required: true
	Width int `json:"width" binding:"required"`
	//座標 [左上XY, 右下XY]
	// required: true
	Coordinate []Coordinate `json:"coordinate" binding:"required,gte=2"`
	//旋轉角度
	// required: true
	Rotation float32 `json:"rotation" binding:"required"`
	//required: true
	NumberOfPost int64 `json:"number_of_post"`
	//required: true
	NumberOfCollection int64 `json:"number_of_collection"`
	// 是否公開收藏 false=不公開, true=公開
	//required: true
	IsCollectionPublic bool `json:"is_collection_public"`
}

type GetOtherUserInfoRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type GetOtherUserInfoResponse struct {
	//required: true
	UserID int64 `json:"user_id"`
	//required: true
	Name string `json:"user_name"`
	//required: true
	CropAvatar string `json:"crop_avatar"`
	//原圖長
	// required: true
	Height int `json:"height" binding:"required"`
	//原圖寬
	// required: true
	Width int `json:"width" binding:"required"`
	//座標 [左上XY, 右下XY]
	// required: true
	Coordinate []Coordinate `json:"coordinate" binding:"required,gte=2"`
	//旋轉角度
	// required: true
	Rotation float32 `json:"rotation" binding:"required"`
	//required: true
	NumberOfPost int64 `json:"number_of_post"`
	//required: true
	NumberOfCollection int64 `json:"number_of_collection"`
	// 是否公開收藏 false=不公開, true=公開
	//required: true
	IsCollectionPublic bool `json:"is_collection_public"`
}

type UpdateUserInfoRequest struct {
	// name for this user
	// example: Feen Lin
	Name *string `json:"name"`

	// public user collection post
	// example: false
	IsCollectionPublic *bool `json:"is_collection_public"`
}

type UpdateUserInfoResponse struct {
	//required: true
	Photo string `json:"photo"`
}

type CollectArticleRequest struct {
	// cocktail id
	// required: true
	// example: 123456
	ID int64 `json:"id" binding:"required"`
}

type GetUserFavoriteCocktailListRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type GetUserFavoriteCocktailListResponse struct {
	//required: true
	IsPublic bool `json:"is_public"`
	//required: true
	Total int64 `json:"total"`
	//required: true
	FavoriteCocktailList []FavoriteCocktail `json:"favorite_cocktail_list"`
}

type FavoriteCocktail struct {
	//required: true
	CocktailID int64 `json:"cocktail_id"`
	//required: true
	UserName string `json:"user_name"`
	//required: true
	Photo string `json:"photo"`
	//required: true
	Title string `json:"title"`
	//required: true
	IsCollected bool `json:"is_collected"`
	//required: true
	CollectedDate string `json:"collected_date"`
}

type DeleteFavoriteCocktailRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type DeleteFavoriteCocktailResponse struct {
	//required: true
	CommandID string `json:"command_id"`
}

type UpdateUserAvatarRequest struct {
	//原始照片
	OriginAvatar string `json:"origin_avatar"`

	//裁切後照片
	// required: true
	CropAvatar string `json:"crop_avatar" binding:"required"`

	//座標 [左上XY, 右下XY]
	// required: true
	Coordinate []Coordinate `json:"coordinate" binding:"required,gte=2"`

	//旋轉角度
	// required: true
	Rotation float32 `json:"rotation"`
}

type Coordinate struct {
	// required: true
	X float32 `json:"x"`
	// required: true
	Y float32 `json:"y"`
}

type UpdateUserAvatarResponse struct {
	//required: true
	Photo string `json:"photo"`
}
