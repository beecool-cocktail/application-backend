package internal

import (
	_cocktailHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/cocktail/delivery/http"
	_cocktailMySQLRepo "github.com/beecool-cocktail/application-backend/api/cocktail/repository/mysql"
	_cocktailUsecase "github.com/beecool-cocktail/application-backend/api/cocktail/usecase"
	_socialAccountGoogleOAuth "github.com/beecool-cocktail/application-backend/api/social-account/repository/google-oauth2"
	_socialAccountMySQLRepo "github.com/beecool-cocktail/application-backend/api/social-account/repository/mysql"
	_socialAccountUsecase "github.com/beecool-cocktail/application-backend/api/social-account/usecase"
	_userHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/user/delivery/http"
	_userRepo "github.com/beecool-cocktail/application-backend/api/user/repository/mysql"
	_userCache "github.com/beecool-cocktail/application-backend/api/user/repository/redis"
	_userUsecase "github.com/beecool-cocktail/application-backend/api/user/usercase"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Init(cfgFile string) {

	appService, err := service.NewService(cfgFile)
	if err != nil {
		logrus.Panicf("start appService failed - %s", err)
	}

	initializeRoutes(appService)
	go util.StartUserIdGenerator()

	logrus.Fatal(appService.HTTP.Run(appService.Configure.HTTP.Address + ":" + appService.Configure.HTTP.Port))
}

func initializeRoutes(s *service.Service) {
	g := s.Configure.Others.GoogleOAuth2
	googleOAuthConfig := &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		RedirectURL:  g.RedirectURL,
		Scopes:       g.Scopes,
		Endpoint:     google.Endpoint,
	}

	// CORSMiddleware for all handler
	s.HTTP.Use(middleware.CORSMiddleware())

	userMySQLRepo := _userRepo.NewMySQLUserRepository(s.DB)
	socialAccountMySQLRepo := _socialAccountMySQLRepo.NewMySQLSocialAccountRepository(s.DB)
	cocktailMySQLRepo := _cocktailMySQLRepo.NewMySQLCocktailRepository(s.DB)

	userRedisRepo := _userCache.NewRedisUserRepository(s.Redis)

	socialAccountGoogleOAuthRepo := _socialAccountGoogleOAuth.NewGoogleOAuthSocialAccountRepository(googleOAuthConfig)


	userUsecase := _userUsecase.NewUserUsecase(userMySQLRepo, userRedisRepo)
	socialAccountUsecase := _socialAccountUsecase.NewSocialAccountUsecase(userMySQLRepo, socialAccountMySQLRepo, socialAccountGoogleOAuthRepo)
	cocktailUsecase := _cocktailUsecase.NewCocktailUsecase(cocktailMySQLRepo)

	_userHandlerHttpDelivery.NewUserHandler(s, userUsecase, socialAccountUsecase)
	_cocktailHandlerHttpDelivery.NewCocktailHandler(s, cocktailUsecase)
}