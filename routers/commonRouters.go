package routers

/*
路由分组和具体路由的定义
*/
import (
	"post-manage/controller/common"

	"github.com/gin-gonic/gin"
)

// 用户管理的路由定义
func CommonRoutersInit(s *gin.Engine) {
	routerGroup := s.Group("/common")
	{
		routerGroup.POST("/doUserLogin", common.CommonController{}.DoUserLogin)
		routerGroup.POST("/doUserLogout", common.CommonController{}.DoUserLogout)
		routerGroup.GET("/getVeryCode", common.CommonController{}.DoMakeVeryCode)
		routerGroup.GET("/checkVeryCode", common.CommonController{}.DoCheckVeryCode)
	}
}
