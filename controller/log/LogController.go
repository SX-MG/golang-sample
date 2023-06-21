package log

/*
日志的controller定义
*/
import (
	"net/http"
	"post-manage/model"
	"post-manage/utils"
	"post-manage/utils/controllerUtil"
	"post-manage/view/formVerify"
	"strings"

	"github.com/gin-gonic/gin"
)

type LogController struct{}

// 列表
func (uc LogController) List(c *gin.Context) {
	paramet := &controllerUtil.Paramet{} //做form字段注入
	if err := c.BindQuery(&paramet); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("commonList", "common", paramet, c) {
		return
	}

	list := make([]model.LogModel, 1)                                                    //初始化list对象，实际上可以简化写一下
	db := utils.GetDBInfo()                                                              //获取数据库连接池中的数据库操作对象
	whereString := controllerUtil.GetWhereString("f_title", "like", paramet.QueryString) //生成where语句，可以再优化
	db.Where(whereString).Scopes(controllerUtil.Paginate(paramet.Start, paramet.Limit)).Scopes(controllerUtil.Order(paramet.OrderField, paramet.OrderType)).Find(&list)
	var md model.BaseModel = &model.LogModel{}
	c.JSON(http.StatusOK, utils.GetCustomResponseMsg(utils.DEFAULT_STATUS_CODE_OK, 0, utils.DEFAULT_SUCCESS_MSG, controllerUtil.MakeWebList(md, whereString, list)))

}

func (uc LogController) ClearAll(c *gin.Context) {
	paramet := &controllerUtil.Paramet{} //做form字段注入
	if err := c.BindQuery(&paramet); err != nil {
		c.JSON(http.StatusOK, utils.GetCustomResponseMsg(0, 0, err.Error()))
		return
	}
	//做字段参数的可用性检测，其他检测也可以做,如果有检测失败的，则直接返回不再继续
	if !formVerify.DoFormFieldVerify("Log", "deleteForm", paramet, c) {
		return
	}
	idList := strings.Split(paramet.Ids, ",")
	db := utils.GetDBInfo()
	db.Delete(&model.LogModel{}, idList)

	c.JSON(http.StatusOK, utils.GetDefaultSuccessResponseMsg())
}
