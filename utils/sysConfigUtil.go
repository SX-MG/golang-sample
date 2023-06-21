package utils

//系统配置工具类
import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

/*
系统配置项管理工具类
*/
var config *viper.Viper
var allSettings = make(map[string]interface{})

// 系统配置的初始化函数
func InitSysConfigInfo() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infoln("Config file changed:", e.Name)
	})
	logrus.Infoln("读取系统配置信息...")
	config = viper.New()
	config.AddConfigPath("./conf/")
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Errorln("找不到配置文件..")
			os.Exit(1)
		} else {
			logrus.Errorln("配置文件出错..")
			os.Exit(1)
		}
	} else {
		allSettings = config.AllSettings()
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")

		fmt.Println("##########################################################")
		for key, val := range allSettings {
			printSettings(key, val)
		}
		fmt.Println("##########################################################")
	}
}

func printSettings(k string, v interface{}) {
	switch v.(type) {
	case map[string]interface{}:
		fmt.Printf("####%v\n", k)
		for sk, sv := range v.(map[string]interface{}) {
			printSettings(sk, sv)
		}
		fmt.Println("-----------------------------------------------------------")
	default:
		fmt.Printf("#####\t%10v\t\t\t\t%-4v\n", k, v)
	}
}

const (
	DBHostSettingKey          string = "database.host"
	DBPortSettingKey          string = "database.port"
	DBUsernameSettingKey      string = "database.username"
	DBPasswordSettingKey      string = "database.password"
	DBNameSettingKey          string = "database.dbname"
	DBTimeoutSettingKey       string = "database.timeout"
	LogLevelSettingKey        string = "log.level"
	ServerPortSettingKey      string = "server.port"
	ServerDebugModeSettingKey string = "server.debug"
	DefaultAdminUserPassword  string = "server.defaultPassword"
)

// 给外部文件调用的 获取系统配置的函数
func GetSysSettingByKey(key string) interface{} {
	// fmt.Println("=========", config)
	if config == nil {
		InitSysConfigInfo()
	}
	val := config.Get(key)
	if val == nil {
		logrus.Errorln("读取配置信息返回空,", key)
	}
	logrus.Debugln("read syscofig value:", val)
	return val
}
