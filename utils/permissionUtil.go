package utils

/*
系统全局的公共权限检测函数
*/
import (
	"net/http"
	"post-manage/utils/jwtUtils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 全局的访问权限控制实现
func GlobalPermissionCheck(c *gin.Context) {
	path := c.Request.URL.Path
	if path == "/common/doUserLogin" || path == "/common/getVeryCode" || path == "/common/checkVeryCode" {
		c.Next()
	} else {
		accessToken := c.Request.Header.Get("access-token")
		if accessToken == "" {
			c.JSON(http.StatusOK, GetDefaultNoPermissionResponseMsg())
			c.Abort()
		} else {
			myClaims := &jwtUtils.MyClaims{}
			myClaims, err := jwtUtils.ParseToken(accessToken)
			if err != nil {
				logrus.Warnln(err.Error())
				c.JSON(http.StatusOK, GetDefaultNoPermissionResponseMsg())
				c.Abort()
				return
			}
			if myClaims == nil {
				c.JSON(http.StatusOK, GetDefaultNoPermissionResponseMsg())
				c.Abort()
				return
			}
			c.Next()
		}
	}

}
