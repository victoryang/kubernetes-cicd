package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 日志中输出 error 依赖 c.Error 函数，如果发生了 error，需要将其注册进去
// 在日志组件中将 username 和 info 字段打印，供接口 log 使用

// formatErr 错误流程返回错误给前端使用
func formatErr(err error) gin.H {
	return gin.H{"ErrMessage": fmt.Sprint(err)}
}

func HttpResponseWithForbidden(c *gin.Context, err error) {
	// 将 error 记录到 context 里面，供 log 使用
	if err != nil {
		defer func() {
			c.Error(err)
		}()

		c.AbortWithStatusJSON(http.StatusForbidden, formatErr(err))
	}
}

func HttpResponseWithBadRequest(c *gin.Context, err error) {
	// 将 error 记录到 context 里面，供 log 使用
	if err != nil {
		defer func() {
			c.Error(err)
		}()

		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
	}
}

func HttpResponseWithSuccess(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}