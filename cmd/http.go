package cmd

import (
	"log"
	"movie-festival/external"
	"movie-festival/helpers"
	"movie-festival/internal/api"
	"movie-festival/internal/interfaces"
	"movie-festival/internal/repository"
	"movie-festival/internal/services"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	d := dependencyInject()

	config := cors.DefaultConfig()
	allowOriginsEnv := helpers.GetEnv("ALLOW_ORIGINS", "")
	allowOrigins := strings.Split(allowOriginsEnv, ",")
	config.AllowOrigins = allowOrigins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r := gin.Default()
	r.Use(cors.New(config))
	r.GET("/health", d.HealthcheckAPI.HealthcheckHandlerHTTP)

	walletV1 := r.Group("/movie/v1")
	walletV1.GET("/", d.MiddlewareValidateToken, d.WalletAPI.GetAll)
	walletV1.POST("/", d.MiddlewareValidateToken, d.WalletAPI.Create)
	walletV1.PUT("/:id", d.MiddlewareValidateToken, d.WalletAPI.Update)
	walletV1.POST("/import", d.MiddlewareValidateToken, d.WalletAPI.UploadExcel)
	walletV1.PATCH("/pick/:id", d.MiddlewareValidateToken, d.WalletAPI.DataPick)
	walletV1.PATCH("/like/:id", d.MiddlewareValidateToken, d.WalletAPI.DataLike)
	walletV1.PATCH("/dislike/:id", d.MiddlewareValidateToken, d.WalletAPI.DataDislike)
	walletV1.GET("/template", d.WalletAPI.GetTemplate)

	err := r.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	HealthcheckAPI interfaces.IHealthcheckAPI
	WalletAPI      interfaces.IMovieAPI
	External       interfaces.IExternal
}

func dependencyInject() Dependency {
	healthcheckSvc := &services.Healthcheck{}
	healthcheckAPI := &api.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	movieRepo := &repository.MovieRepo{
		DB: helpers.DB,
	}

	movieSvc := &services.MovieService{
		MovieRepo: movieRepo,
	}
	movieAPI := &api.MovieAPI{
		MovieService: movieSvc,
	}

	external := &external.External{}

	return Dependency{
		HealthcheckAPI: healthcheckAPI,
		WalletAPI:      movieAPI,
		External:       external,
	}
}
