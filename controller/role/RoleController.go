package role

/*
角色的controller定义
*/
import (
	"net/http"
	"post-manage/model"
	"post-manage/utils"
	"post-manage/utils/controllerUtil"
	"post-manage/view/formVerify"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// 列表
func (uc RoleController) List(c *gin.Context) {
	paramet := &controllerUtil.Paramet{} //做form字段注入
	if err := c.BindQuery(&paramet); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("commonList", "common", paramet, c) {
		return
	}

	list := make([]model.RoleModel, 1)                                                      //初始化list对象，实际上可以简化写一下
	db := utils.GetDBInfo()                                                                 //获取数据库连接池中的数据库操作对象
	whereString := controllerUtil.GetWhereString("f_roleName", "like", paramet.QueryString) //生成where语句，可以再优化
	db.Where(whereString).Scopes(controllerUtil.Paginate(paramet.Start, paramet.Limit)).Scopes(controllerUtil.Order(paramet.OrderField, paramet.OrderType)).Find(&list)
	var md model.BaseModel = &model.RoleModel{}
	c.JSON(http.StatusOK, utils.GetCustomResponseMsg(utils.DEFAULT_STATUS_CODE_OK, 0, utils.DEFAULT_SUCCESS_MSG, controllerUtil.MakeWebList(md, whereString, list)))

}
func (uc RoleController) Add(c *gin.Context) {
	m := model.RoleModel{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("role", "addForm", m, c) {
		return
	}
	db := utils.GetDBInfo()
	var existCount int64 = 0
	db.Model(model.RoleModel{}).Where("f_roleName=?", m.RoleName).Count(&existCount)
	if existCount > 0 {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(-110, 0, "角色名称不能重复"))
		return
	}
	m.CreateDate = int(time.Now().Unix())
	db.Create(&m)
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}
func (uc RoleController) Update(c *gin.Context) {
	m := model.RoleModel{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("role", "updateForm", m, c) {
		return
	}
	db := utils.GetDBInfo()
	oldModel := model.RoleModel{}
	db.Model(model.RoleModel{}).Where("f_id=?", m.Id).First(&oldModel)
	if oldModel.Id != m.Id {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, "没有找到原始数据,当前更新失败"))
		return
	}
	m.CreateDate = oldModel.CreateDate
	db.Model(model.RoleModel{}).Where("f_id=?", m.Id).Save(m)
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

func (uc RoleController) Delete(c *gin.Context) {
	paramet := &controllerUtil.Paramet{} //做form字段注入
	if err := c.BindQuery(&paramet); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("Role", "deleteForm", paramet, c) {
		return
	}
	idList := strings.Split(paramet.Ids, ",")
	db := utils.GetDBInfo()
	db.Where("f_roleId in (?)", strings.Split(paramet.Ids, ",")).Delete(model.RoleUserModel{})
	db.Delete(&model.RoleModel{}, idList)
	// for i := 0; i < len(idList); i++ {
	// 	uid := idList[i]
	// 	roleId, err := strconv.Atoi(uid)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	db.Model(model.RoleModel{}).Delete("f_roleId=?", roleId)
	// }
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

// 修改角色所属用户信息
func (uc RoleController) ModifyRoleUserInfo(c *gin.Context) {
	m := RoleUserParam{} //做form字段注入
	if err := c.Bind(&m); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("role", "modifyRoleUserInfo", m, c) {
		return
	}
	db := utils.GetDBInfo()
	db.Where("f_roleId=?", m.RoleId).Delete(model.RoleUserModel{})
	userIds := strings.Split(m.UserIds, ",")
	for i := 0; i < len(userIds); i++ {
		uid := userIds[i]
		userId, err := strconv.Atoi(uid)
		if err != nil {
			continue
		}
		roleUser := model.RoleUserModel{}
		roleUser.RoleId = m.RoleId
		roleUser.UserId = userId
		db.Create(&roleUser)
	}
	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}

type RoleUserParam struct {
	RoleId  int    `json:"roleId" form:"roleId" uri:"roleId"`
	UserIds string `json:"userIds" form:"userIds" uri:"userIds"`
}
