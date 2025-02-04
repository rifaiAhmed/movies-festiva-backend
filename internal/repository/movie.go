package repository

import (
	"context"
	"movie-festival/internal/models"
	"strings"

	"gorm.io/gorm"
)

type MovieRepo struct {
	DB *gorm.DB
}

func (r *MovieRepo) CreateMovie(ctx context.Context, movie *models.Movie) error {
	return r.DB.Create(movie).Error
}

func (r *MovieRepo) FindByID(ctx context.Context, ID int) (models.Movie, error) {
	var (
		resp = models.Movie{}
	)
	if err := r.DB.Where("id = ?", ID).First(&resp).Error; err != nil {
		return resp, err
	}
	return resp, nil
}

func (r *MovieRepo) Update(ctx context.Context, req *models.Movie) error {
	return r.DB.Save(req).Error
}

func (r *MovieRepo) GetAll(ctx context.Context, objComp models.ComponentServerSide, param string) ([]models.Movie, error) {
	var (
		resp    []models.Movie
		limit   = objComp.Limit
		isOrder = objComp.SortBy + ` ` + objComp.SortType
		isWhere = isWhere(param, objComp)
		query   = r.DB.Table("movies")
	)
	if isWhere != "" {
		query.Where(isWhere)
	}
	if err := query.Order(isOrder).Limit(limit).Find(&resp).Error; err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *MovieRepo) Counting(param string, objComp models.ComponentServerSide) (int64, error) {
	var (
		count   int64
		isWhere = isWhere(param, objComp)
		query   = r.DB.Table("movies")
	)
	if isWhere != "" {
		query = query.Where(isWhere)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func isWhere(param string, objComp models.ComponentServerSide) string {
	var (
		isWhere = param
	)
	if objComp.Search != "" {
		if isWhere != "" {
			isWhere += ` and (LOWER(description) LIKE '%` + strings.ToLower(objComp.Search) + `%' OR LOWER(artist) LIKE '%` + strings.ToLower(objComp.Search) + `%' OR LOWER(genre) LIKE '%` + strings.ToLower(objComp.Search) + `%' OR LOWER(url) LIKE '%` + strings.ToLower(objComp.Search) + `%')`
		} else {
			isWhere = `(LOWER(description) LIKE '%` + strings.ToLower(objComp.Search) + `%' OR LOWER(artist) LIKE '%` + strings.ToLower(objComp.Search) + `%' OR LOWER(genre) LIKE '%` + strings.ToLower(objComp.Search) + `%' OR LOWER(url) LIKE '%` + strings.ToLower(objComp.Search) + `%')`
		}
	}

	return isWhere
}
