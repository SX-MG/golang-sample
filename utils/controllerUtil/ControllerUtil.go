package controllerUtil

//controller公共工具函数
import (
	"fmt"
	"post-manage/model"
	"post-manage/utils"
	"strings"

	"gorm.io/gorm"
)

type Paramet struct {
	Start       int    `json:"start" uri:"start" form:"start"`
	Limit       int    `json:"limit" uri:"limit" form:"limit"`
	OrderField  string `json:"orderField" uri:"orderField" form:"orderField"`
	OrderType   string `json:"orderType" uri:"orderType" form:"orderType"`
	QueryString string `json:"queryString" uri:"queryString" form:"queryString"`
	Ids         string `json:"ids" uri:"ids" form:"ids"`
}

// 分页封装
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 获取controller中基于queryString参数的where语句
func GetWhereString(fieldName string, oper string, queryString string) string {
	if queryString == "" {
		return ""
	} else {
		queryString = strings.ReplaceAll(queryString, "'", "")
		queryString = strings.ReplaceAll(queryString, "\"", "")
		queryString = strings.ReplaceAll(queryString, "`", "")
		switch oper {
		case "like":
			return fmt.Sprintf("%v like '%v'", fieldName, "%%"+queryString+"%%")
		case "eq":
			return fmt.Sprintf("%v = '%v'", fieldName, queryString)
		case "gt":
			return fmt.Sprintf("%v> '%v'", fieldName, queryString)
		case "gq":
			return fmt.Sprintf("%v >= '%v'", fieldName, queryString)
		case "lt":
			return fmt.Sprintf("%v < '%v'", fieldName, queryString)
		case "lq":
			return fmt.Sprintf("%v <= '%v'", fieldName, queryString)
		case "nq":
			return fmt.Sprintf("%v <> '%v'", fieldName, queryString)
		default:
			return fmt.Sprintf("%v like '%v'", fieldName, "%%"+queryString+"%%")
		}
	}
}

/*
*
排序判断 需要补充的逻辑：字段有效性的判断和排序方式的判断，现在主要靠前端传参
*/
func Order(orderField, orderType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if orderField == "" {
			return db
		}
		if orderType == "" {
			orderType = "ASC"
		}
		return db.Order(orderField + " " + orderType)
	}
}

func GetRowCount(model model.BaseModel, whereString string) int64 {
	db := utils.GetDBInfo()
	var total int64
	db.Table(model.TableName()).Select("f_id").Where(whereString).Count(&total)
	return total
}

func MakeWebList(model model.BaseModel, whereString string, rowList interface{}) WebList {
	rowTotal := GetRowCount(model, whereString)
	wl := WebList{
		RowTotal: int(rowTotal),
		RowList:  rowList,
	}
	return wl
}

type WebList struct {
	RowTotal int         `json:"rowTotal"`
	RowList  interface{} `json:"rowList"`
}
