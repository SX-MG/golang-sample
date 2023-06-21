package utils

//日志工具类
import (
	"post-manage/model"
	"post-manage/utils/jwtUtils"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const LogType_warn int = 1
const LogType_debug int = 2
const LogType_error int = 3
const LogType_fail int = 4

func Warn(th interface{}, c *gin.Context, category, title, content string) {
	writeLog(LogType_warn, th, c, category, title, content)
}
func Debug(th interface{}, c *gin.Context, category, title, content string) {
	writeLog(LogType_debug, th, c, category, title, content)
}
func Error(th interface{}, c *gin.Context, category, title, content string) {
	writeLog(LogType_error, th, c, category, title, content)
}
func Fail(th interface{}, c *gin.Context, category, title, content string) {
	writeLog(LogType_fail, th, c, category, title, content)
}

func writeLog(logType int, th interface{}, c *gin.Context, category, title, content string) {
	log := model.LogModel{}
	log.Category = category
	log.Content = content
	log.Title = title
	log.LogType = logType
	log.CreateDate = time.Now().Unix()
	if th != nil {
		switch th.(type) {
		case struct{}:
			tp := reflect.TypeOf(th)
			if tp.Kind() == reflect.Ptr {
				tp = tp.Elem()
			}
			log.Content = log.Content + "|" + tp.Name()
		case int:
			log.Content = log.Content + "|" + strconv.Itoa(int(th.(int)))
		case string:
			log.Content = log.Content + "|" + th.(string)
		default:

		}
	}
	if c != nil {
		ip := c.ClientIP()
		log.Ip = ip
		accessToken := c.GetHeader("Access-Token")
		info, err := jwtUtils.ParseToken(accessToken)
		if err == nil {
			log.UserId = info.UserId
			log.UserName = info.UserName
		}
	}
	db := GetDBInfo()
	if !db.Migrator().HasTable(model.LogModel{}) {
		db.AutoMigrate(model.LogModel{})
	}
	db.Create(&log)
}
