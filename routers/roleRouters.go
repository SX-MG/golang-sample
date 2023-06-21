package routers

/*
路由分组和具体路由的定义
*/
import (
	"post-manage/controller/role"

	"github.com/gin-gonic/gin"
)

// 用户管理的路由定义
func RoleRoutersInit(s *gin.Engine) {
	userGroup := s.Group("/Role")
	{
		userGroup.GET("/list", role.RoleController{}.List)
		userGroup.POST("/add", role.RoleController{}.Add)
		userGroup.POST("/update", role.RoleController{}.Update)
		userGroup.POST("/delete", role.RoleController{}.Delete)
		userGroup.POST("/modifyRoleUserInfo", role.RoleController{}.ModifyRoleUserInfo)
	}
}
