package utils

/*
数据库连接池的工具
*/
import (
	"fmt"
	"os"
	"post-manage/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

// 获取mysql连接的对象
func getMysqlConnectDB() *gorm.DB {
	//配置MySQL连接参数
	username := GetSysSettingByKey(DBUsernameSettingKey) //账号
	password := GetSysSettingByKey(DBPasswordSettingKey) //密码
	host := GetSysSettingByKey(DBHostSettingKey)         //数据库地址，可以是Ip或者域名
	port := GetSysSettingByKey(DBPortSettingKey)         //数据库端口
	Dbname := GetSysSettingByKey(DBNameSettingKey)       //数据库名
	timeout := GetSysSettingByKey(DBTimeoutSettingKey)   //连接超时，10秒

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, fmt.Sprintf("%v", password), host, port, Dbname, timeout)
	logrus.Debugln(dsn)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	} else {
		logrus.Warnln("数据库连接成功")
		return db
	}
}

// 给系统调用的数据库连接池初始化函数
func InitDBConnPoolInfo() {
	_db = getMysqlConnectDB()
	if _db == nil {
		os.Exit(1)
	} else {
		sqlDB, _ := _db.DB()
		//设置数据库连接池参数
		sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
		sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	}
}

// 获取数据库连接对象
func GetDBInfo() *gorm.DB {
	if _db == nil {
		logrus.Errorln("没有可用的数据库连接")
		panic("没有可用的数据库连接")
	}
	d := GetSysSettingByKey(ServerDebugModeSettingKey)
	if d == "open" {
		return _db.Debug()
	}
	return _db
}

// 初始化数据表结构 自动建表
func AutoCareatDBTableStruct() {
	db := GetDBInfo()
	if !db.Migrator().HasTable(model.AdminUserModel{}) {
		db.AutoMigrate(model.AdminUserModel{})
	}
	if !db.Migrator().HasTable(model.LogModel{}) {
		db.AutoMigrate(model.LogModel{})
	}
	if !db.Migrator().HasTable(model.RoleModel{}) {
		db.AutoMigrate(model.RoleModel{})
	}
	if !db.Migrator().HasTable(model.RoleUserModel{}) {
		db.AutoMigrate(model.RoleUserModel{})
	}
}
