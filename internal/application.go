package internal

import (
	_cocktailIngredientMySQLRepo "github.com/beecool-cocktail/application-backend/api/cockingredient/repository/mysql"
	_cocktailPhotoMySQLRepo "github.com/beecool-cocktail/application-backend/api/cockphoto/repository/mysql"
	_cocktailStepMySQLRepo "github.com/beecool-cocktail/application-backend/api/cockstep/repository/mysql"
	_cocktailHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/cocktail/delivery/http"
	_cocktailElasticSearchRepo "github.com/beecool-cocktail/application-backend/api/cocktail/repository/elastic"
	_cocktailFileRepo "github.com/beecool-cocktail/application-backend/api/cocktail/repository/file"
	_cocktailMySQLRepo "github.com/beecool-cocktail/application-backend/api/cocktail/repository/mysql"
	_cocktailRedisRepo "github.com/beecool-cocktail/application-backend/api/cocktail/repository/redis"
	_cocktailUsecase "github.com/beecool-cocktail/application-backend/api/cocktail/usecase"
	_commandHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/command/delivery/http"
	_commandRedis "github.com/beecool-cocktail/application-backend/api/command/repository/redis"
	_commandUsecase "github.com/beecool-cocktail/application-backend/api/command/usecase"
	_favoriteCocktailMySQLRepo "github.com/beecool-cocktail/application-backend/api/favoritecock/repository/mysql"
	_favoriteCocktailUsecase "github.com/beecool-cocktail/application-backend/api/favoritecock/usecase"
	_socialAccountGoogleOAuth "github.com/beecool-cocktail/application-backend/api/social-account/repository/google-oauth2"
	_socialAccountMySQLRepo "github.com/beecool-cocktail/application-backend/api/social-account/repository/mysql"
	_socialAccountRedisRepo "github.com/beecool-cocktail/application-backend/api/social-account/repository/redis"
	_socialAccountUsecase "github.com/beecool-cocktail/application-backend/api/social-account/usecase"
	_userHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/user/delivery/http"
	_userFileRepo "github.com/beecool-cocktail/application-backend/api/user/repository/file"
	_userRepo "github.com/beecool-cocktail/application-backend/api/user/repository/mysql"
	_userCache "github.com/beecool-cocktail/application-backend/api/user/repository/redis"
	_userUsecase "github.com/beecool-cocktail/application-backend/api/user/usercase"
	cmdPattern "github.com/beecool-cocktail/application-backend/command"
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

	if appService.Configure.HTTP.IsTLS {
		logrus.Fatal(appService.HTTP.RunTLS(appService.Configure.HTTP.Address+":"+appService.Configure.HTTP.Port,
			appService.Configure.HTTP.CertificateFile, appService.Configure.HTTP.KeyFile))
	} else {
		logrus.Fatal(appService.HTTP.Run(appService.Configure.HTTP.Address + ":" + appService.Configure.HTTP.Port))
	}
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

	s.HTTP.Static("/static", "/static/images")

	// Repository dependency injection
	transactionRepo := _transactionRepo.NewDBRepository(s.DB)
	userMySQLRepo := _userRepo.NewMySQLUserRepository(s.DB)
	socialAccountMySQLRepo := _socialAccountMySQLRepo.NewMySQLSocialAccountRepository(s.DB)
	cocktailMySQLRepo := _cocktailMySQLRepo.NewMySQLCocktailRepository(s.DB)
	cocktailIngredientMySQLRepo := _cocktailIngredientMySQLRepo.NewMySQLCocktailIngredientRepository(s.DB)
	cocktailStepMySQLRepo := _cocktailStepMySQLRepo.NewMySQLCocktailStepRepository(s.DB)
	cocktailPhotoMySQLRepo := _cocktailPhotoMySQLRepo.NewMySQLCocktailStepRepository(s.DB)
	favoriteCocktailMySQLRepo := _favoriteCocktailMySQLRepo.NewMySQLFavoriteCocktailRepository(s.DB)

	cocktailElasticSearchRepo := _cocktailElasticSearchRepo.NewElasticSearchCocktailRepository(s.Elastic)

	userRedisRepo := _userCache.NewRedisUserRepository(s.Redis)
	cocktailRedisRepo := _cocktailRedisRepo.NewRedisCocktailRepository(s.Redis)
	commandRedisRepo := _commandRedis.NewRedisCommandRepository(s.Redis)
	socialAccountRedisRepo := _socialAccountRedisRepo.NewRedisSocialAccountRepository(s.Redis)

	userFileRepo := _userFileRepo.NewFileUserRepository()
	cocktailFileMySQL := _cocktailFileRepo.NewFileUserRepository()

	socialAccountGoogleOAuthRepo := _socialAccountGoogleOAuth.NewGoogleOAuthSocialAccountRepository(googleOAuthConfig)

	//operator
	operatorHandler := cmdPattern.NewOperatorHandler()
	favoriteCocktailOperator := cmdPattern.NewFavoriteCocktailOperator(favoriteCocktailMySQLRepo, cocktailMySQLRepo, cocktailRedisRepo, transactionRepo)
	operatorHandler.SetOperator(cmdPattern.FavoriteCocktailDelete, favoriteCocktailOperator)

	// Usecase dependency injection
	userUsecase := _userUsecase.NewUserUsecase(s, userMySQLRepo, userRedisRepo, userFileRepo, transactionRepo)
	socialAccountUsecase := _socialAccountUsecase.NewSocialAccountUsecase(userMySQLRepo, userRedisRepo,
		socialAccountMySQLRepo, socialAccountRedisRepo, socialAccountGoogleOAuthRepo)
	favoriteCocktailUsecase := _favoriteCocktailUsecase.NewFavoriteCocktailUsecase(favoriteCocktailMySQLRepo,
		cocktailMySQLRepo, cocktailRedisRepo, cocktailPhotoMySQLRepo, userMySQLRepo, userRedisRepo, commandRedisRepo,
		transactionRepo)
	cocktailUsecase := _cocktailUsecase.NewCocktailUsecase(s, cocktailMySQLRepo, cocktailRedisRepo,
		cocktailElasticSearchRepo, cocktailFileMySQL, cocktailPhotoMySQLRepo, cocktailIngredientMySQLRepo,
		cocktailStepMySQLRepo, userMySQLRepo, favoriteCocktailMySQLRepo, transactionRepo)
	commandUsecase := _commandUsecase.NewCommandUsecase(s, commandRedisRepo, operatorHandler)

	// Delivery dependency injection
	_userHandlerHttpDelivery.NewUserHandler(s, userUsecase, socialAccountUsecase, cocktailUsecase, favoriteCocktailUsecase, *middlewareHandler)
	_cocktailHandlerHttpDelivery.NewCocktailHandler(s, cocktailUsecase, userUsecase, *middlewareHandler)
	_commandHandlerHttpDelivery.NewCommandHandler(s, commandUsecase, *middlewareHandler)
}
