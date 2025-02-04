package services

import (
	"context"
	"errors"
	"fmt"
	"movie-festival/internal/interfaces"
	"movie-festival/internal/models"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type MovieService struct {
	MovieRepo interfaces.IMovieRepo
}

func (s *MovieService) Create(ctx context.Context, movie *models.Movie) error {
	return s.MovieRepo.CreateMovie(ctx, movie)
}

func (s *MovieService) Update(ctx context.Context, req *models.Movie, ID int) error {
	obj, err := s.MovieRepo.FindByID(ctx, ID)
	if err != nil {
		return errors.New("failed find data by id")
	}
	req.ID = obj.ID
	req.LikeCount = obj.LikeCount
	req.DislikeCount = obj.DislikeCount
	req.CountSelected = obj.CountSelected
	req.CreatedAt = obj.CreatedAt
	req.UpdatedAt = time.Now()
	return s.MovieRepo.Update(ctx, req)
}

func (s *MovieService) GetAll(ctx context.Context, objComp models.ComponentServerSide, isData string) ([]models.Movie, int64, error) {
	var (
		obj   []models.Movie
		count int64
		param = ""
	)

	objChan := make(chan []models.Movie)
	countChan := make(chan int64)
	errChan := make(chan error, 2)
	if strings.ToLower(isData) == "selected" {
		param = "count_selected > 0"
	}
	if strings.ToLower(isData) == "like" {
		param = "like_count > 0"
	}
	go func() {
		movies, err := s.MovieRepo.GetAll(ctx, objComp, param)
		if err != nil {
			errChan <- err
			return
		}
		objChan <- movies
	}()

	go func() {
		total, err := s.MovieRepo.Counting(param, objComp)
		if err != nil {
			errChan <- err
			return
		}
		countChan <- total
	}()

	for i := 0; i < 2; i++ {
		select {
		case movies := <-objChan:
			obj = movies
		case total := <-countChan:
			count = total
		case err := <-errChan:
			return nil, 0, err
		}
	}

	return obj, count, nil
}

func (s *MovieService) InsertFromExcel(ctx context.Context, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("failed to read sheet: %v", err)
	}

	var movies []models.Movie

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 6 {
			continue
		}

		movie := models.Movie{
			Title:       row[0],
			Description: row[1],
			Durasi:      row[2],
			Url:         row[3],
			Genre:       row[4],
			Artist:      row[5],
		}

		movies = append(movies, movie)
		err = s.MovieRepo.CreateMovie(ctx, &movie)
	}

	if len(movies) == 0 {
		return fmt.Errorf("no valid movies found in the file")
	}

	return nil
}

func (s *MovieService) DataPick(ctx context.Context, ID int) error {
	var (
		obj models.Movie
	)
	obj, err := s.MovieRepo.FindByID(ctx, ID)
	if err != nil {
		return errors.New("failed find data by id")
	}
	obj.CountSelected += 1
	return s.MovieRepo.Update(ctx, &obj)
}

func (s *MovieService) DataLike(ctx context.Context, ID int) error {
	var (
		obj models.Movie
	)
	obj, err := s.MovieRepo.FindByID(ctx, ID)
	if err != nil {
		return errors.New("failed find data by id")
	}
	obj.LikeCount += 1
	return s.MovieRepo.Update(ctx, &obj)
}

func (s *MovieService) DataDislike(ctx context.Context, ID int) error {
	var (
		obj models.Movie
	)
	obj, err := s.MovieRepo.FindByID(ctx, ID)
	if err != nil {
		return errors.New("failed find data by id")
	}
	obj.DislikeCount += 1
	return s.MovieRepo.Update(ctx, &obj)
}
