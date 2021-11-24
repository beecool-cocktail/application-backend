package internal

import (
	_userHandlerHttpDelivery "github.com/beecool-cocktail/application-backend/api/user/delivery/http"
	_userRepo "github.com/beecool-cocktail/application-backend/api/user/repository/mysql"
	_userCache "github.com/beecool-cocktail/application-backend/api/user/repository/redis"
	_userUsecase "github.com/beecool-cocktail/application-backend/api/user/usercase"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/sirupsen/logrus"
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

	store := cookie.NewStore([]byte("secret"))
	s.HTTP.Use(sessions.Sessions("mysession", store))

	middlewareHandler := middleware.NewMiddlewareHandler(s.Redis)
	// CORSMiddleware for all handler
	s.HTTP.Use(middlewareHandler.CORSMiddleware())

	userMySQLRepo := _userRepo.NewMySQLUserRepository(s.DB)
	userRedisRepo := _userCache.NewRedisUserRepository(s.Redis)

	userUsecase := _userUsecase.NewUserUsecase(userMySQLRepo, userRedisRepo)

	_userHandlerHttpDelivery.NewUserHandler(s.HTTP, userUsecase, middlewareHandler)
}