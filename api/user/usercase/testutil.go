package usercase

import "github.com/beecool-cocktail/application-backend/domain"

func matchedByUpdateToken(token *domain.UserCache, id int64) bool {
	return token.Id == id
}
