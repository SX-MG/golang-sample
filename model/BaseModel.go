package model

//基础数据类型的定义
type BaseModel interface {
	TableName() string
}
