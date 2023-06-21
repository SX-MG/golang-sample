package test

import (
	"fmt"

	"github.com/cengsin/oracle"
	"gorm.io/gorm"
)

func DoConnOracle() {
	_, err := gorm.Open(oracle.Open("system/1qaz2wsx#EDC@localhost:1521/orcl"), &gorm.Config{})
	if err != nil {
		fmt.Println("连接oracle数据库异常：", err.Error())
	} else {
		fmt.Println("oracle数据库连接成功")

	}

	// do somethings
}
