package internal

import (
	_cocktailIngredientMySQLRepo "github.com/beecool-cocktail/application-backend/api/cockingredient/repository/mysql"
	_cocktailPhotoMySQLRepo "github.com/beecool-cocktail/application-backend/api/cockphoto/repository/mysql"
	_cocktailStepMySQLRepo "github.com/beecool-cocktail/application-backend/api/cockstep/repository/mysql"
	_cocktailHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/cocktail/delivery/http"
	_cocktailFileRepo "github.com/beecool-cocktail/application-backend/api/cocktail/repository/file"
	_cocktailMySQLRepo "github.com/beecool-cocktail/application-backend/api/cocktail/repository/mysql"
	_cocktailUsecase "github.com/beecool-cocktail/application-backend/api/cocktail/usecase"
	_socialAccountGoogleOAuth "github.com/beecool-cocktail/application-backend/api/social-account/repository/google-oauth2"
	_socialAccountMySQLRepo "github.com/beecool-cocktail/application-backend/api/social-account/repository/mysql"
	_socialAccountUsecase "github.com/beecool-cocktail/application-backend/api/social-account/usecase"
	_userHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/user/delivery/http"
	_userFileRepo "github.com/beecool-cocktail/application-backend/api/user/repository/file"
	_userRepo "github.com/beecool-cocktail/application-backend/api/user/repository/mysql"
	_userCache "github.com/beecool-cocktail/application-backend/api/user/repository/redis"
	_userUsecase "github.com/beecool-cocktail/application-backend/api/user/usercase"
	_transactionRepo "github.com/beecool-cocktail/application-backend/db/repository/mysql"
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

	middlewareHandler := middleware.NewMiddlewareHandler(s)
	// CORSMiddleware for all handler
	s.HTTP.Use(middlewareHandler.CORSMiddleware())

	// Repository dependency injection
	transactionRepo := _transactionRepo.NewDBRepository(s.DB)
	userMySQLRepo := _userRepo.NewMySQLUserRepository(s.DB)
	socialAccountMySQLRepo := _socialAccountMySQLRepo.NewMySQLSocialAccountRepository(s.DB)
	cocktailMySQLRepo := _cocktailMySQLRepo.NewMySQLCocktailRepository(s.DB)
	cocktailIngredientMySQLRepo := _cocktailIngredientMySQLRepo.NewMySQLCocktailIngredientRepository(s.DB)
	cocktailStepMySQLRepo := _cocktailStepMySQLRepo.NewMySQLCocktailStepRepository(s.DB)
	cocktailPhotoMySQLRepo := _cocktailPhotoMySQLRepo.NewMySQLCocktailStepRepository(s.DB)

	userRedisRepo := _userCache.NewRedisUserRepository(s.Redis)

	userFileRepo := _userFileRepo.NewFileUserRepository()
	cocktailFileMySQL := _cocktailFileRepo.NewFileUserRepository()

	socialAccountGoogleOAuthRepo := _socialAccountGoogleOAuth.NewGoogleOAuthSocialAccountRepository(googleOAuthConfig)

	// Usecase dependency injection
	userUsecase := _userUsecase.NewUserUsecase(s, userMySQLRepo, userRedisRepo, userFileRepo, transactionRepo)
	socialAccountUsecase := _socialAccountUsecase.NewSocialAccountUsecase(userMySQLRepo, userRedisRepo,
		socialAccountMySQLRepo, socialAccountGoogleOAuthRepo)
	cocktailUsecase := _cocktailUsecase.NewCocktailUsecase(s, cocktailMySQLRepo, cocktailFileMySQL,
		cocktailPhotoMySQLRepo, cocktailIngredientMySQLRepo, cocktailStepMySQLRepo, userMySQLRepo, transactionRepo)

	// Delivery dependency injection
	_userHandlerHttpDelivery.NewUserHandler(s, userUsecase, socialAccountUsecase, *middlewareHandler)
	_cocktailHandlerHttpDelivery.NewCocktailHandler(s, cocktailUsecase, *middlewareHandler)
}
