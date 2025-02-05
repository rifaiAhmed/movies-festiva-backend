package interfaces

import (
	"context"
	"movie-festival/internal/models"

	"github.com/gin-gonic/gin"
)

type IMovieRepo interface {
	CreateMovie(ctx context.Context, movie *models.Movie) error
	FindByID(ctx context.Context, ID int) (models.Movie, error)
	Update(ctx context.Context, req *models.Movie) error
	GetAll(ctx context.Context, objComp models.ComponentServerSide, param string) ([]models.Movie, error)
	Counting(param string, objComp models.ComponentServerSide) (int64, error)
}

type IMovieService interface {
	Create(ctx context.Context, movie *models.Movie) error
	Update(ctx context.Context, req *models.Movie, ID int) error
	GetAll(ctx context.Context, objComp models.ComponentServerSide, isData string) ([]models.Movie, int64, error)
	InsertFromExcel(ctx context.Context, filePath string) error
	DataPick(ctx context.Context, ID int) error
	DataLike(ctx context.Context, ID int) error
	DataDislike(ctx context.Context, ID int) error
}

type IMovieAPI interface {
	Create(*gin.Context)
	Update(c *gin.Context)
	UploadExcel(c *gin.Context)
	GetAll(c *gin.Context)
	DataPick(c *gin.Context)
	DataLike(c *gin.Context)
	DataDislike(c *gin.Context)
	GetTemplate(c *gin.Context)
}
