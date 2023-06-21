package model

//模型定义 日志模型
type LogModel struct {
	Id         int    `json:"id" gorm:"column:f_id;primary_key;AUTO_INCREMENT" form:"id" uri:"id"`
	UserName   string `json:"userName" gorm:"column:f_userName" form:"userName" uri:"userName"`
	CreateDate int64  `json:"createDate" gorm:"column:f_createDate" form:"createDate" uri:"createDate"`
	Title      string `json:"title" gorm:"column:f_title" form:"title" uri:"title"`
	LogType    int    `json:"logType" gorm:"column:f_logType" form:"logType" uri:"logType"`
	Category   string `json:"category" gorm:"column:f_category" form:"category" uri:"category"`
	UserId     int    `json:"userId" gorm:"column:f_userId" form:"userId" uri:"userId"`
	Ip         string `json:"ip" gorm:"column:f_ip" form:"ip" uri:"ip"`
	ClientType string `json:"clientType" gorm:"column:f_clientType" form:"clientType" uri:"clientType"`
	Content    string `json:"content" gorm:"column:f_content" form:"content" uri:"content"`
}

func (u *LogModel) TableName() string {
	return "t_Log"
}
