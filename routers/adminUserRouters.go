package routers

/*
路由分组和具体路由的定义
*/
import (
	"post-manage/controller/adminUser"

	"github.com/gin-gonic/gin"
)

// 用户管理的路由定义
func AdminUserRoutersInit(s *gin.Engine) {
	userGroup := s.Group("/adminUser")
	{
		userGroup.GET("/list", adminUser.AdminUserController{}.List)
		userGroup.POST("/add", adminUser.AdminUserController{}.Add)
		userGroup.POST("/update", adminUser.AdminUserController{}.Update)
		userGroup.GET("/delete", adminUser.AdminUserController{}.Delete)
		userGroup.POST("/resetPassword", adminUser.AdminUserController{}.ResetPassword)
		userGroup.POST("/updateCurrentUserPassword", adminUser.AdminUserController{}.UpdateCurrentUserPassword)
	}
}
