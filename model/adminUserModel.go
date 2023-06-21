package model

//模型定义 管理员模型
type AdminUserModel struct {
	Id            int    `json:"id" gorm:"column:f_id;primary_key;AUTO_INCREMENT" form:"id" uri:"id"`
	UserName      string `json:"userName" gorm:"column:f_userName" form:"userName" uri:"userName"`
	RealName      string `json:"realName" gorm:"column:f_realName" form:"realName" uri:"realName"`
	Password      string `json:"password" gorm:"column:f_passwd" form:"password" uri:"password"`
	CreateDate    int64  `json:"createDate" gorm:"column:f_createDate" form:"createDate" uri:"createDate"`
	LoginCount    int    `json:"loginCount" gorm:"column:f_loginCount" form:"loginCount" uri:"loginCount"`
	LastLoginTime int64  `json:"lastLoginTime" gorm:"column:f_lastLoginTime" form:"lastLoginTime" uri:"lastLoginTime"`
}

func (u *AdminUserModel) TableName() string {
	return "t_adminUser"
}
