package helper

import "github.com/gin-gonic/gin"

type Response struct {
	Status string `json:"status,omitempty"`
	Code   int32  `json:"code,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func ResponseOutput(c *gin.Context, status_code int32, status_message string, data interface{}) {
	resp := Response{
		Status: status_message,
		Code:   status_code,
		Data:   data,
	}
	c.JSON(int(status_code), resp)
}

func GenerateTotalPage(totalData, limit int64) (totalPage int64) {
	totalPage = totalData / limit
	modTotalPage := totalData % limit
	if modTotalPage > 0 {
		totalPage++
	}

	return totalPage
}

func GetOffsetAndLimit(page, pageSize int64) (offset, limit int64) {
	offset = (page - 1) * pageSize
	limit = pageSize

	return offset, limit
}
