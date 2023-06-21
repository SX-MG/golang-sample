package routers

/*
路由分组和具体路由的定义
*/
import (
	"post-manage/controller/log"

	"github.com/gin-gonic/gin"
)

// 用户管理的路由定义
func LogRoutersInit(s *gin.Engine) {
	userGroup := s.Group("/log")
	{
		userGroup.GET("/list", log.LogController{}.List)
		userGroup.POST("/clearAll", log.LogController{}.ClearAll)

	}
}
