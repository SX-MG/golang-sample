package model

//模型定义 角色模型
type RoleModel struct {
	Id         int    `json:"id" gorm:"column:f_id;primary_key;AUTO_INCREMENT" form:"id" uri:"id"`
	RoleName   string `json:"roleName" gorm:"column:f_roleName" form:"roleName" uri:"roleName"`
	CreateDate int    `json:"createDate" gorm:"column:f_createDate" form:"createDate" uri:"createDate"`
	Order      int    `json:"order" gorm:"column:f_order" form:"order" uri:"order"`
}

func (u *RoleModel) TableName() string {
	return "t_role"
}
