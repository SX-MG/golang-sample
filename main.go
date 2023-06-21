package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"post-manage/model"
	"post-manage/routers"
	"post-manage/test"
	"post-manage/utils"
	"post-manage/view/formVerify"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"
	"github.com/sirupsen/logrus"
)

// 系统启动后的初始化函数，自动调用
func init() {
	s := ` ___                                                   _         _                       _              
	/  _>  ___   ___  ___ ._ _ _ ._ _ _  ___ ._ _   ___  _| |._ _ _ <_>._ _   ___ _ _  ___ _| |_ ___ ._ _ _ 
	| <_/\\/ . \\ / | '/ . \\| ' ' || ' ' |/ . \\| ' | <_> |/ . || ' ' || || ' | <_-<| | |<_-<  | | / ._>| ' ' |
	'____/\\___/ \\_|_.\\___/|_|_|_||_|_|_|\\___/|_|_| <___|\\___||_|_|_||_||_|_| /__/'_. |/__/  |_| \\___.|_|_|_|
																				 <___'                      
	`
	fmt.Println(s)                               //输出banner
	utils.InitSysConfigInfo()                    //从配置文件中读取系统配置
	formVerify.InitFormMeesages()                //初始化表单验证信息
	utils.InitDBConnPoolInfo()                   //初始化数据库连接池
	logrus.SetLevel(logrus.TraceLevel)           //设置日志级别，默认最高级别
	logrus.SetFormatter(&logrus.JSONFormatter{}) //日志格式设置 json，方便后期采集
	utils.AutoCareatDBTableStruct()              //自动创建表结构
}

// 整个系统的入口函数
func main() {
	beforeSysInit()
	Task()
	s := gin.Default()
	s.Use(utils.InitWebSession()) //初始化web session
	s.Use(utils.GlobalPermissionCheck)
	routers.CommonRoutersInit(s) //应用全局权限判断的实现
	routers.AdminUserRoutersInit(s)
	routers.LogRoutersInit(s)
	routers.RoleRoutersInit(s)                                                           //路由的初始化 user管理的路由
	s.Run(":" + fmt.Sprintf("%v", utils.GetSysSettingByKey(utils.ServerPortSettingKey))) //在特定端口启动web服务器
	test.DoConnOracle()
}

// 初始化默认数据
func beforeSysInit() {
	args := os.Args
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.ToLower(arg) == "sysinit" {
			fmt.Println("首次使用，系统初始化...")
			db := utils.GetDBInfo()
			db.Delete(&model.AdminUserModel{})
			dp := utils.GetSysSettingByKey(utils.DefaultAdminUserPassword)

			initUser := model.AdminUserModel{
				UserName:   "admin",
				RealName:   "管理员",
				CreateDate: time.Now().Unix(),
				Password:   fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", dp)))),
				LoginCount: 0,
			}
			db.Create(&initUser)
			fmt.Println("默认用户初始化成功")
		}
	}
}

func Task() {
	cron := cronexpr.MustParse("*/10 * * * * * *") //用cron库生成一个cronexpr.Expression对象
	next := cron.Next(time.Now())                  //计算下次触发时间的时间对象
	for {
		now := time.Now()                        //每次循环计算获取当前时间
		if next.Before(now) || next.Equal(now) { //下次触发时间与当前时间进行对比，等于或者时间已到 则进行任务触发
			doPost()
			files, e := filepath.Glob("d:\\*.xlsx")
			if e != nil {
				fmt.Println("删除文件前发生异常,", e.Error())
			} else {
				for _, f := range files {
					if err := os.Remove(f); err != nil {
						fmt.Println("删除文件发生异常,", err.Error())
					}
				}
			}
			err := utils.ExcelUtil{}.CreateNewExcelFile("d:\\" + fmt.Sprintf("%v", time.Now().Unix()) + ".xlsx")
			if err != nil {
				fmt.Println("excel文件创建失败,", err.Error())
			}
			next = cron.Next(now) //重新计算下次任务时间的时间对象
		}
		select {
		case <-time.NewTicker(time.Second).C: //每秒扫描一遍 循环频率设定
		}
	}
}

func doPost() {
	data := url.Values{"start": {"0"}, "offset": {"xxxx"}}
	body := strings.NewReader(data.Encode())
	clt := http.Client{}
	resp, err := clt.Post("http://localhost:8088/common/doUserLogout", "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		respBody := string(content)
		fmt.Println("respBody:", respBody)
	}
}
