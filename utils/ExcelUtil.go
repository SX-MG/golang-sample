package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelUtil struct{}

func (util ExcelUtil) CreateNewExcelFile(savePath string) error {
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet("Sheet2")
	// 设置单元格的值
	f.SetCellValue("Sheet2", "A1", "Hello world.")
	f.SetCellValue("Sheet1", "B1", 38.958)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs(savePath); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
