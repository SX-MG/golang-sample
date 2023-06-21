package model

//模型定义 角色与管理员的对照模型
type RoleUserModel struct {
	Id     int `json:"id" gorm:"column:f_id;primary_key;AUTO_INCREMENT" form:"id" uri:"id"`
	RoleId int `json:"roleId" gorm:"column:f_roleId" form:"roleId" uri:"roleId"`
	UserId int `json:"userId" gorm:"column:f_userId" form:"userId" uri:"userId"`
}

func (u *RoleUserModel) TableName() string {
	return "t_roleUser"
}
