package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageQuery 分页查询结果。
type PageQuery struct {
	Page     int
	PageSize int
	Offset   int
}

// ParsePageQuery 从 query 解析 page / page_size，带默认值与上限。
func ParsePageQuery(c *gin.Context) PageQuery {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PageQuery{Page: page, PageSize: pageSize, Offset: (page - 1) * pageSize}
}

// PageResult 分页返回结构。
type PageResult struct {
	List     any `json:"list"`
	Total    int64 `json:"total"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
}
