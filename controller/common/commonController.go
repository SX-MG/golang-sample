package common

/*
公共通用的controller定义
*/
import (
	"crypto/md5"
	"net/http"
	"post-manage/model"
	"post-manage/utils"
	"post-manage/utils/jwtUtils"
	"post-manage/view/formVerify"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CommonController struct {
	UserName string `json:"userName" form:"userName" uri:"userName"`
	Password string `json:"password" form:"password" uri:"password"`
	VeryCode string `json:"veryCode" form:"veryCode" uri:"veryCode"`
}

// 用户登录
func (uc CommonController) DoUserLogin(c *gin.Context) {
	m := CommonController{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("userLogin", "common", m, c) {
		return
	}
	db := utils.GetDBInfo()
	model := model.AdminUserModel{}
	db.Where("f_userName=? and f_passwd=?", m.UserName, md5.Sum([]byte(m.Password))).First(&model)
	if model.Id == 0 {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(-2, 0, "没有该用户信息"))
		utils.Warn(CommonController{}, c, "用户登录", "登陆失败", "用户"+model.UserName+"登录失败")
	} else {
		signingString, err := jwtUtils.MakeAccessToken(model.Id, model.UserName)
		if err != nil {
			logrus.Panicln(err.Error)
		}
		c.Writer.Header().Set("access-token", signingString)
		c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
		utils.Warn(CommonController{}, c, "用户登录", "登陆成功", "用户"+model.UserName+"登录成功")
	}
}

// 登出
func (uc CommonController) DoUserLogout(c *gin.Context) {
	//先通过设置jwt token为空的方式，强制客户端退出登录
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

// 生成验证码
func (uc CommonController) DoMakeVeryCode(c *gin.Context) {
	utils.Captcha(c, 4)
}

// 检查验证码是否可用
func (uc CommonController) DoCheckVeryCode(c *gin.Context) {
	code := c.Query("veryCode")
	if code == "" {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(-99, 0, "缺少参数"))
		return
	}
	if utils.CaptchaVerify(c, code) {
		c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
	} else {
		c.JSON(http.StatusOK, utils.GetDefaultFaildResponseMsg())
	}
}
