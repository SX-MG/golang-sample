package paramcheck

//辅助函数
import "github.com/gin-gonic/gin"

// 公共的、从query或post中获取参数值的方法
func GetParamValue(context gin.Context, fieldName string) string {
	v := context.Query(fieldName)
	if v == "" {
		v = context.PostForm(fieldName)
	}
	return v
}

type ParamCheckStruct struct {
	FieldName        string
	CheckFaildStatus int
	CheckFaildMsg    string
}

//生成用于参数检测的字段结构体
func MakeParamCheckInfo(fieldName string, checkFaildStatus int, checkFaildMsg string) ParamCheckStruct {
	o := ParamCheckStruct{}
	o.CheckFaildMsg = checkFaildMsg
	o.CheckFaildStatus = checkFaildStatus
	o.FieldName = fieldName
	return o
}
