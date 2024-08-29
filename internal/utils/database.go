package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Search(db *gorm.DB, query string, fields ...string) *gorm.DB {
	if query != "" {
		for i, field := range fields {
			if i == 0 {
				db = db.Where(field+" ILIKE ?", "%"+query+"%")
			} else {
				db = db.Or(field+" ILIKE ?", "%"+query+"%")
			}
		}
	}
	return db
}

func GetPaginationParams(c *gin.Context) (query string, page, pageSize int) {
	query = c.Query("q")
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return
}

func GetPaginationResponse(page, pageSize int, total int64) gin.H {
	return gin.H{
		"currentPage": page,
		"pageSize":    pageSize,
		"totalItems":  total,
		"totalPages":  int(math.Ceil(float64(total) / float64(pageSize))),
	}
}
