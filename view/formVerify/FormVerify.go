package formVerify

//表单字段验证工具类
import (
	"fmt"
	"net/http"
	"os"
	"post-manage/utils"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

var config *viper.Viper

// 初始化消息定义
func InitFormMeesages() {
	logrus.Warnln("初始化表单验证信息")
	config = viper.New()
	config.AddConfigPath("./conf/")
	config.SetConfigName("message")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Errorln("找不到message配置文件..")
			os.Exit(1)
		} else {
			logrus.Errorln("配置文件出错..")
			os.Exit(1)
		}
	}
}

// controller中调用的表单字段验证函数
func DoFormFieldVerify(messageRootKey, formType string, entity interface{}, context *gin.Context) bool {
	r := reflect.TypeOf(entity)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	if r.NumField() == 0 {
		return true
	}
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		propName := field.Name
		fieldName := field.Tag.Get("form")
		if !(notnullCheck(messageRootKey, formType, entity, propName, fieldName, context)) {
			return false
		}
		if !(lengthCheck(messageRootKey, formType, entity, propName, fieldName, context)) {
			return false
		}
		if !(numCheck(messageRootKey, formType, entity, propName, fieldName, context)) {
			return false
		}
	}
	return true
}

func notnullCheck(messageRootKey, formType string, entity interface{}, propName, fieldName string, context *gin.Context) bool {
	nullableEnabKey := fmt.Sprintf("%v.%v.%v.nullable.enab", messageRootKey, formType, fieldName)
	nullableEnabVal := config.GetBool(nullableEnabKey)
	if nullableEnabVal {
		nullableFailedStatusKey := fmt.Sprintf("%v.%v.%v.nullable.failedStatus", messageRootKey, formType, fieldName)
		nullableFailedMsgKey := fmt.Sprintf("%v.%v.%v.nullable.failedMsg", messageRootKey, formType, fieldName)
		nullableFailedStatusVal := config.GetInt(nullableFailedStatusKey)
		nullableFailedMsgVal := config.GetString(nullableFailedMsgKey)
		if nullableFailedStatusVal == 0 {
			nullableFailedStatusVal = -101
		}
		if nullableFailedMsgVal == "" {
			nullableFailedMsgVal = fieldName + "不能为空"
		}
		//判断字段值是否满足要求
		fieldValue := reflect.ValueOf(entity)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}
		var fv interface{}
		rfv := ""

		switch fieldValue.FieldByName(propName).Kind() {
		case reflect.String:
			fv = fieldValue.FieldByName(propName).String()
			rfv = fv.(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fv = fieldValue.FieldByName(propName).Int()
			rfv = fmt.Sprintf("%v", fv)
		case reflect.Float32, reflect.Float64:
			fv = fieldValue.FieldByName(propName).Float()
			rfv = fmt.Sprintf("%v", fv)
		case reflect.Bool:
			fv = fieldValue.FieldByName(propName).Bool()
			rfv = fmt.Sprintf("%v", fv)
		default:
			rfv = ""
		}
		if rfv == "0" {
			rfv = ""
		}
		if rfv == "" {
			context.JSON(http.StatusOK, utils.GetCustomResponseMsg(nullableFailedStatusVal, 0, nullableFailedMsgVal))
			return false
		}
	}
	return true
}
func lengthCheck(messageRootKey, formType string, entity interface{}, propName, fieldName string, context *gin.Context) bool {
	lenEnabKey := fmt.Sprintf("%v.%v.%v.len.enab", messageRootKey, formType, fieldName)
	lenEnabVal := config.GetBool(lenEnabKey)
	if lenEnabVal {
		lenFailedStatusKey := fmt.Sprintf("%v.%v.%v.len.failedStatus", messageRootKey, formType, fieldName)
		lenFailedMsgKey := fmt.Sprintf("%v.%v.%v.len.failedMsg", messageRootKey, formType, fieldName)
		lenMaxMsgKey := fmt.Sprintf("%v.%v.%v.len.maxLength", messageRootKey, formType, fieldName)
		lenMinMsgKey := fmt.Sprintf("%v.%v.%v.len.minLength", messageRootKey, formType, fieldName)
		lenFailedStatusVal := config.GetInt(lenFailedStatusKey)
		lenFailedMsgVal := config.GetString(lenFailedMsgKey)
		lenMaxVal := config.GetInt(lenMaxMsgKey)
		lenMinVal := config.GetInt(lenMinMsgKey)
		if lenFailedStatusVal == 0 {
			lenFailedStatusVal = -101
		}
		if lenFailedMsgVal == "" {
			lenFailedMsgVal = fieldName + "长度应介于" + strconv.Itoa(lenMinVal) + "-" + strconv.Itoa(lenMaxVal) + "之间"
		}
		//判断字段值是否满足要求
		fieldValue := reflect.ValueOf(entity)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}
		var fv interface{}
		rfv := ""

		switch fieldValue.FieldByName(propName).Kind() {
		case reflect.String:
			fv = fieldValue.FieldByName(propName).String()
			rfv = fv.(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fv = fieldValue.FieldByName(propName).Int()
			rfv = fmt.Sprintf("%v", fv)
		case reflect.Float32, reflect.Float64:
			fv = fieldValue.FieldByName(propName).Float()
			rfv = fmt.Sprintf("%v", fv)
		case reflect.Bool:
			fv = fieldValue.FieldByName(propName).Bool()
			rfv = fmt.Sprintf("%v", fv)
		default:
			rfv = ""
		}
		if rfv == "0" {
			rfv = ""
		}
		if utf8.RuneCountInString(rfv) < lenMinVal || utf8.RuneCountInString(rfv) > lenMaxVal {
			context.JSON(http.StatusOK, utils.GetCustomResponseMsg(lenFailedStatusVal, 0, lenFailedMsgVal))
			return false
		}
	}
	return true
}
func numCheck(messageRootKey, formType string, entity interface{}, propName, fieldName string, context *gin.Context) bool {
	lenEnabKey := fmt.Sprintf("%v.%v.%v.num.enab", messageRootKey, formType, fieldName)
	lenEnabVal := config.GetBool(lenEnabKey)
	if lenEnabVal {
		numFailedStatusKey := fmt.Sprintf("%v.%v.%v.num.failedStatus", messageRootKey, formType, fieldName)
		numFailedMsgKey := fmt.Sprintf("%v.%v.%v.num.failedMsg", messageRootKey, formType, fieldName)
		numFailedStatusVal := config.GetInt(numFailedStatusKey)
		numFailedMsgVal := config.GetString(numFailedMsgKey)
		if numFailedStatusVal == 0 {
			numFailedStatusVal = -101
		}
		if numFailedMsgVal == "" {
			numFailedMsgVal = fieldName + "必须是数字"
		}
		//判断字段值是否满足要求
		fieldValue := reflect.ValueOf(entity)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}
		var fv interface{}
		rfv := ""

		switch fieldValue.FieldByName(propName).Kind() {
		case reflect.String:
			fv = fieldValue.FieldByName(propName).String()
			rfv = fv.(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fv = fieldValue.FieldByName(propName).Int()
			rfv = fmt.Sprintf("%v", fv)
		case reflect.Float32, reflect.Float64:
			fv = fieldValue.FieldByName(propName).Float()
			rfv = fmt.Sprintf("%v", fv)
		case reflect.Bool:
			fv = fieldValue.FieldByName(propName).Bool()
			rfv = fmt.Sprintf("%v", fv)
		default:
			rfv = ""
		}
		if rfv == "0" {
			rfv = ""
		}
		_, err := strconv.Atoi(rfv)
		if err != nil {
			context.JSON(http.StatusOK, utils.GetCustomResponseMsg(numFailedStatusVal, 0, numFailedMsgVal))
			return false
		}
	}
	return true
}
