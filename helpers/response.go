package helpers

import (
	"movie-festival/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendResponseHTTP(c *gin.Context, code int, message string, data interface{}) {
	resp := Response{
		Message: message,
		Data:    data,
	}
	c.JSON(code, resp)
}

var limitData int = 25

func GetLimitData() int {
	return limitData
}

func ComptServerSidePre(c *gin.Context) (models.ComponentServerSide, error) {
	var objComponentServerSide models.ComponentServerSide

	limitInt, _ := strconv.Atoi(c.Query("limit"))
	skipInt, _ := strconv.Atoi(c.Query("skip"))

	condition := ""

	if limitInt != 0 {
		condition = "limit"
	}

	if c.Query("sort") != "" && c.Query("sortBy") != "" {
		condition = "sort"
	}

	if c.Query("search") != "" {
		condition = "search"
	}

	offset := skipInt

	objComponentServerSide = models.ComponentServerSide{
		Limit:     limitInt,
		Skip:      skipInt,
		SortType:  c.Query("sort_type"),
		SortBy:    strings.ToUpper(c.Query("sort_by")),
		Search:    strings.ToUpper(c.Query("search")),
		Offset:    offset,
		Condition: condition,
		From:      c.Query("from"),
		To:        c.Query("to"),
	}

	return objComponentServerSide, nil

}

func GetTotalPage(totalData int, limit int) int {
	bagi_awal := totalData / limit
	sisa_bagi := totalData % limit

	if sisa_bagi > 0 {
		bagi_awal = bagi_awal + 1
	}

	return bagi_awal
}

func APIResponseView(message string, code int, status string, count int64, limit int, data interface{}) ResponseView {
	meta := MetaResponse{
		Message:      message,
		Code:         code,
		Status:       status,
		ResponseTime: DateToStdNow(),
		TotalData:    count,
		TotalPage:    int64(GetTotalPage(int(count), limit)),
	}

	jsonResponse := ResponseView{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func DateToStdNow() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

type ResponseView struct {
	Meta MetaResponse `json:"meta"`
	Data interface{}  `json:"data"`
}

type MetaResponse struct {
	Message      string `json:"message"`
	Code         int    `json:"int"`
	Status       string `json:"status"`
	ResponseTime string `json:"response_time"`
	CurrentPage  int64  `json:"current_page"`
	TotalPage    int64  `json:"totalPages"`
	TotalData    int64  `json:"totalData"`
}
