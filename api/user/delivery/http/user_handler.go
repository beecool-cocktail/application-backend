package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

type UserHandler struct {
	Configure   *service.Configure
	UserUsecase domain.UserUsecase
}

func NewUserHandler(s *service.Service, userUsecase domain.UserUsecase) {
	handler := &UserHandler{
		Configure:   s.Configure,
		UserUsecase: userUsecase,
	}

	s.HTTP.GET("/api/google-login", handler.SocialLogin)
}

// swagger:operation GET /google-login login googleLogin
// ---
// summary: Login with google OAuth2
// description: todo
// security:
// - Bearer: []
// responses:
//  201:
//   headers:
//     Authorization:
//       type: string
//       description: jwt token
func (u *UserHandler) SocialLogin(c *gin.Context) {
	g := u.Configure.Others.GoogleOAuth2
	googleOAuth2Config := &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		RedirectURL:  g.RedirectURL,
		Scopes:       g.Scopes,
		Endpoint:     google.Endpoint,
	}

	url := googleOAuth2Config.AuthCodeURL("whispering-corner")

	c.Redirect(http.StatusTemporaryRedirect, url)
}