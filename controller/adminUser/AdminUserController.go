package adminUser

/*
用户的controller定义
*/
import (
	"crypto/md5"
	"fmt"
	"net/http"
	"post-manage/model"
	"post-manage/utils"
	"post-manage/utils/controllerUtil"
	"post-manage/utils/jwtUtils"
	"post-manage/view/formVerify"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type AdminUserController struct{}

// 列表
func (uc AdminUserController) List(c *gin.Context) {
	paramet := &controllerUtil.Paramet{} //做form字段注入
	if err := c.BindQuery(&paramet); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("commonList", "common", paramet, c) {
		return
	}

	list := make([]model.AdminUserModel, 1)                                                 //初始化list对象，实际上可以简化写一下
	db := utils.GetDBInfo()                                                                 //获取数据库连接池中的数据库操作对象
	whereString := controllerUtil.GetWhereString("f_userName", "like", paramet.QueryString) //生成where语句，可以再优化
	db.Where(whereString).Scopes(controllerUtil.Paginate(paramet.Start, paramet.Limit)).Scopes(controllerUtil.Order(paramet.OrderField, paramet.OrderType)).Find(&list)
	var md model.BaseModel = &model.AdminUserModel{}
	c.JSON(http.StatusOK, utils.GetCustomResponseMsg(utils.DEFAULT_STATUS_CODE_OK, 0, utils.DEFAULT_SUCCESS_MSG, controllerUtil.MakeWebList(md, whereString, list)))

}
func (uc AdminUserController) Add(c *gin.Context) {
	m := model.AdminUserModel{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("adminUser", "addForm", m, c) {
		return
	}
	db := utils.GetDBInfo()
	var existCount int64 = 0
	db.Model(model.AdminUserModel{}).Where("f_userName=?", m.UserName).Count(&existCount)
	if existCount > 0 {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(-110, 0, "用户名不能重复"))
		return
	}
	m.CreateDate = time.Now().Unix()
	m.LoginCount = 0
	m.Password = fmt.Sprintf("%v", utils.GetSysSettingByKey(utils.DefaultAdminUserPassword))
	m.Password = fmt.Sprintf("%x", md5.Sum([]byte(m.Password)))
	db.Create(&m)
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg(m.Id))
}
func (uc AdminUserController) Update(c *gin.Context) {
	m := model.AdminUserModel{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("adminUser", "updateForm", m, c) {
		return
	}
	db := utils.GetDBInfo()
	oldModel := model.AdminUserModel{}
	db.Model(model.AdminUserModel{}).Where("f_id=?", m.Id).First(&oldModel)
	if oldModel.Id != m.Id {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, "没有找到原始数据,当前更新失败"))
		return
	}
	m.CreateDate = oldModel.CreateDate
	m.LastLoginTime = oldModel.LastLoginTime
	m.LoginCount = oldModel.LoginCount
	m.Password = oldModel.Password
	m.UserName = oldModel.UserName
	db.Model(model.AdminUserModel{}).Where("f_id=?", m.Id).Save(m)
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

func (uc AdminUserController) Delete(c *gin.Context) {
	paramet := &controllerUtil.Paramet{} //做form字段注入
	if err := c.BindQuery(&paramet); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("adminUser", "deleteForm", paramet, c) {
		return
	}
	idList := strings.Split(paramet.Ids, ",")
	db := utils.GetDBInfo()
	db.Delete(&model.AdminUserModel{}, idList)
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

func (uc AdminUserController) ResetPassword(c *gin.Context) {
	m := ResetPasswordParam{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("adminUser", "resetPasswordForm", m, c) {
		return
	}
	db := utils.GetDBInfo()
	oldModel := model.AdminUserModel{}
	db.Model(model.AdminUserModel{}).Where("f_id=?", m.Id).First(&oldModel)
	if oldModel.Id != m.Id {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, "没有找到原始数据,当前更新失败"))
		return
	}
	if oldModel.Password != m.OldPsd {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, "旧密码不正确"))
		return
	}
	oldModel.Password = utils.GetSysSettingByKey(utils.DefaultAdminUserPassword).(string)
	db.Model(model.AdminUserModel{}).Where("f_id=?", m.Id).Save(oldModel)
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

func (uc AdminUserController) UpdateCurrentUserPassword(c *gin.Context) {
	accessToken := c.Request.Header.Get("access-token")
	myClaims := &jwtUtils.MyClaims{}
	myClaims, err := jwtUtils.ParseToken(accessToken)
	if err != nil {
		logrus.Warnln(err.Error())
		c.JSON(http.StatusOK, utils.GetDefaultNoPermissionResponseMsg())
		c.Abort()
		return
	}
	if myClaims == nil {
		c.JSON(http.StatusOK, utils.GetDefaultNoPermissionResponseMsg())
		c.Abort()
		return
	}
	m := ResetPasswordParam{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("adminUser", "updatePasswordForm", m, c) {
		return
	}
	db := utils.GetDBInfo()
	oldModel := model.AdminUserModel{}
	db.Model(model.AdminUserModel{}).Where("f_id=?", myClaims.UserId).First(&oldModel)
	if oldModel.Password != m.OldPsd {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, "就密码错误"))
		return
	}

	oldModel.Password = m.NewPsd
	db.Model(model.AdminUserModel{}).Where("f_id=?", myClaims.UserId).Save(oldModel)
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

type ResetPasswordParam struct {
	OldPsd string `json:"oldPsd" form:"oldPsd" uri:"oldPsd"`
	NewPsd string `json:"newPsd" form:"newPsd" uri:"newPsd"`
	Id     int    `json:"id" form:"id" uri:"id"`
}
