package utils

// 服务端反馈给客户端的消息体定义
type ResponseMsg struct {
	Status int         `json:"status"` // 状态
	Id     int         `json:"id"`     // 编辑等状态下用于返回ID
	Msg    string      `json:"msg"`    // 具体消息
	Extra  interface{} `json:"extra"`  //扩展内容，用于填充map、list等
}

// 生成responseMsg消息体内容
func MakeResponseMsg(status, id int, msg string, extra interface{}) ResponseMsg {
	rm := ResponseMsg{}
	rm.Status = status
	rm.Id = id
	rm.Msg = msg
	rm.Extra = extra
	return rm
}

const (
	DEFAULT_STATUS_CODE_OK           int    = 99
	DEFAULT_SUCCESS_MSG              string = "process success"
	DEFAULT_STATUS_CODE_FAILD        int    = -1
	DEFAULT_FAILD_MSG                string = "process faild"
	DEFAULT_STATUS_CODE_NOPERMISSION int    = -9
	DEFAULT_NOPERMISSION_MSG         string = "无权限"
)

// 获取默认的操作成功消息体返回值
func GetDefaultSuccessResponseMsg(args ...interface{}) ResponseMsg {
	id := 0
	msg := DEFAULT_SUCCESS_MSG
	if len(args) > 0 {
		switch args[0].(type) {
		case int:
			id = args[0].(int)
		default:
			id = 0
		}
	}
	if len(args) > 1 {
		switch args[1].(type) {
		case string:
			msg = args[1].(string)
		default:
			id = 0
		}

	}
	resp := MakeResponseMsg(DEFAULT_STATUS_CODE_OK, id, msg, nil)
	return resp
}

// 获取默认的操作失败消息体返回值
func GetDefaultFaildResponseMsg() ResponseMsg {
	resp := MakeResponseMsg(DEFAULT_STATUS_CODE_FAILD, 0, DEFAULT_FAILD_MSG, nil)
	return resp
}

// 获取默认的消息体无授权的返回值
func GetDefaultNoPermissionResponseMsg() ResponseMsg {
	resp := MakeResponseMsg(DEFAULT_STATUS_CODE_NOPERMISSION, 0, DEFAULT_NOPERMISSION_MSG, nil)
	return resp
}

// 获取自定义的消息体返回值
func GetCustomResponseMsg(status, id int, msg string, extra ...interface{}) ResponseMsg {
	var ext interface{}
	ext = nil
	if len(extra) > 0 {
		ext = extra[0]
	}
	resp := MakeResponseMsg(status, id, msg, ext)
	return resp
}
