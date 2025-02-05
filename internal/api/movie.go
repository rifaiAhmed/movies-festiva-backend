package api

import (
	"context"
	"movie-festival/constants"
	"movie-festival/helpers"
	"movie-festival/internal/interfaces"
	"movie-festival/internal/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type MovieAPI struct {
	MovieService interfaces.IMovieService
}

func (api *MovieAPI) Create(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.Movie
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse request: ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	err := api.MovieService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Error("failed to create wallet: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, req)
}

func (api *MovieAPI) Update(c *gin.Context) {
	var (
		log = helpers.Logger
		req *models.Movie
	)
	var inputID models.UriId
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		log.Error("failed to get id : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Error("failed to parse request: ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	err = api.MovieService.Update(c.Request.Context(), req, inputID.ID)
	if err != nil {
		log.Error("failed to create wallet: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, req)
}

func (api *MovieAPI) GetAll(c *gin.Context) {
	var (
		log             = helpers.Logger
		objComponent, _ = helpers.ComptServerSidePre(c)
		tipe            = c.Query("type")
	)
	if objComponent.Limit == 0 {
		objComponent.Limit = helpers.GetLimitData()
	}
	obj, total, err := api.MovieService.GetAll(c, objComponent, tipe)
	if err != nil {
		log.Error("failed to get data : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}
	if total == 0 {
		helpers.SendResponseHTTP(c, http.StatusOK, "data empty", nil)
		return
	}
	response := helpers.APIResponseView("Succesfully Get Data!", http.StatusOK, "Succesfully", total, objComponent.Limit, obj)
	response.Meta.CurrentPage = (int64(objComponent.Skip) / int64(objComponent.Limit)) + 1
	c.JSON(http.StatusOK, response)
}

func (api *MovieAPI) UploadExcel(c *gin.Context) {
	var (
		log = helpers.Logger
	)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	err = api.MovieService.InsertFromExcel(context.Background(), filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := os.Remove(filePath); err != nil {
		log.Println("Failed to delete file:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "File processed successfully"})
}

func (api *MovieAPI) DataPick(c *gin.Context) {
	var (
		log = helpers.Logger
	)
	var inputID models.UriId
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		log.Error("failed to get id : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	err = api.MovieService.DataPick(c, inputID.ID)
	if err != nil {
		log.Error("failed to Update data selected: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (api *MovieAPI) DataLike(c *gin.Context) {
	var (
		log = helpers.Logger
	)
	var inputID models.UriId
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		log.Error("failed to get id : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	err = api.MovieService.DataLike(c, inputID.ID)
	if err != nil {
		log.Error("failed to Update data liked: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (api *MovieAPI) DataDislike(c *gin.Context) {
	var (
		log = helpers.Logger
	)
	var inputID models.UriId
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		log.Error("failed to get id : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	err = api.MovieService.DataDislike(c, inputID.ID)
	if err != nil {
		log.Error("failed to Update data dislike: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (api *MovieAPI) GetTemplate(c *gin.Context) {
	var (
		pathFile = "./uploads/movies.xlsx"
		fileName = "movies.xlsx"
	)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(pathFile)
}
