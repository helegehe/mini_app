package httper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var (
	author = "系统提示: "

	// DefaultWidth 默认Excel列宽
	DefaultWidth = 30.0
)

// 如果超过27*26列 暂不支持
func getCellLoc(row, col int) (cl string) {
	if col < 26 {
		cl = fmt.Sprintf("%c%d", col+65, row+1)
	}
	if col >= 26 && col < 27*26 {
		cl = fmt.Sprintf("%c%c%d", col/26+64, col%26+65, row+1)
	}
	return cl
}

func setSqref(sq, sqref string) string {
	if sq == "" {
		sq = sqref
	} else {
		sq = fmt.Sprintf("%s %s", sq, sqref)
	}
	return sq
}

//
func getCellLen(col int) (rl string) {
	var c string
	if col < 26 {
		c = fmt.Sprintf("%c", col+65)

	}
	if col >= 26 && col < 27*26 {
		c = fmt.Sprintf("%c%c", col/26+64, col%26+65)
	}
	for i := 1; i < 1000; i++ {
		rl = setSqref(rl, fmt.Sprintf("%s%d", c, i))
	}
	return rl
}

func numToCol(col int) (cl string) {
	if col < 26 {
		cl = fmt.Sprintf("%c", col+65)

	}
	if col >= 26 && col < 27*26 {
		cl = fmt.Sprintf("%c%c", col/26+64, col%26+65)
	}
	return cl
}

// `{"author":"Excelize: ","text":"This is a comment."}`
type comment struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

type cellStyle struct {
	Comment    string `json:"comment"`
	Validation string `json:"validation"`
}

// getStytle 根据第一行数据获取注解和数据校验数据
// 注解批注格式"SN(必填){comment:xxxxxx,validation:1,2,3}" validation中的数据以’,‘分割
func getStytle(data string) (*cellStyle, string) {
	start := strings.IndexByte(data, '{')
	end := strings.IndexByte(data, '}')
	//fmt.Printf("%s, %d, %d\n", data, start, end)
	if start == -1 {
		return nil, data
	}
	var cs cellStyle
	err := json.Unmarshal([]byte(data[start:end+1]), &cs)
	if err != nil {
		return nil, data[:start]
	}
	return &cs, data[:start]
}

// 获取列名
func getColName(col int) (c string) {
	if col < 26 {
		c = fmt.Sprintf("%c", col+65)

	}
	if col >= 26 && col < 27*26 {
		c = fmt.Sprintf("%c%c", col/26+64, col%26+65)
	}
	return c
}

// ExportConfig 导出配置，文件名称、sheet页名称、列名称
type ExportConfig struct {
	FileName  string       `json:"filename"`
	SheetName string       `json:"sheetName"`
	Titles    []ExcelTitle `json:"titles"`
}

// ExcelTitle 表格Title
type ExcelTitle struct {
	Key        string  `json:"key"`
	Name       string  `json:"name"`
	Comment    string  `json:"comment"`
	Validation string  `json:"validation"`
	Width      float64 `json:"width"`
}

// WriteToXLSXTS 将数据写入xlsx
func WriteToXLSXTS(sheetName string, titles []ExcelTitle, values []map[string]interface{}) (*excelize.File, error) {
	xlsx := excelize.NewFile()

	if sheetName == "" {
		sheetName = "Sheet1"
	}
	index := xlsx.NewSheet(sheetName)
	// write title
	for i := range titles {
		title := titles[i]
		rl := getCellLen(i)
		col := getCellLoc(0, i)
		xlsx.SetCellStr(sheetName, col, title.Name)
		if title.Comment != "" {
			_ = xlsx.AddComment(sheetName, col, fmt.Sprintf(`{"author":"%s","text":"%s"}`, author, title.Comment))
		}
		if title.Validation != "" {
			dvRange := excelize.NewDataValidation(true)
			dvRange.Sqref = rl
			_ = dvRange.SetDropList(strings.Split(title.Validation, ","))
			xlsx.AddDataValidation(sheetName, dvRange)
		}

		if title.Width <= 0 {
			title.Width = DefaultWidth
		}
		xlsx.SetColWidth(sheetName, getColName(i), getColName(i), title.Width)
	}

	// write data
	for j := range values {
		for i := range titles {
			title := titles[i]
			vv := values[j][title.Key]
			cl := getCellLoc(j+1, i)
			xlsx.SetCellValue(sheetName, cl, vv)
		}
	}
	xlsx.SetActiveSheet(index)

	return xlsx, nil
}

// writeToXLSX 将数据写入xlsx
func writeToXLSX(data map[string][][]string, title []string, name string) (*excelize.File, error) {
	xlsx := excelize.NewFile()
	for k := range data {
		index := xlsx.NewSheet(k)

		// 填标题
		if title != nil {
			for i := range title {
				cl := getCellLoc(0, i)
				rl := getCellLen(i)
				cs, da := getStytle(title[i])
				xlsx.SetCellStr(k, cl, da)
				// 配置注解
				if cs != nil && cs.Comment != "" && cl != "" {
					xlsx.AddComment(k, cl, fmt.Sprintf(`{"author":"%s","text":"%s"}`, author, cs.Comment))
				}
				// 设置数据校验
				if cs != nil && cs.Validation != "" {
					dvRange := excelize.NewDataValidation(true)
					dvRange.Sqref = rl
					dvRange.SetDropList(strings.Split(cs.Validation, ","))
					xlsx.AddDataValidation(k, dvRange)
				}
			}
		}
		// 填数据
		for i := range data[k] {
			// 针对每一行数据的头设置样式，比如批注，数据校验
			for j := range data[k][i] {
				cl := getCellLoc(i+1, j)
				xlsx.SetCellValue(k, cl, data[k][i][j])
			}
		}
		xlsx.SetActiveSheet(index)
		xlsx.SetColWidth(k, "A", numToCol(len(title)), 20)
	}

	// 保存到文件
	if name != "" {
		err := xlsx.SaveAs(name)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return xlsx, nil
}

const (
	// Device 导出设备表
	Device = "device"
	// DeviceScanned 导出设备表
	DeviceScanned = "device-scanned"
	// DeviceSetting 导出装机参数表
	DeviceSetting = "device-setting"
	// HardwareTemplate 导出硬件模版
	HardwareTemplate = "hardware-template"
	// SystemTemplate 导出系统模版
	SystemTemplate = "system-template"
	// ImageTemplate 导出镜像模版
	ImageTemplate = "Image-template"
	// DeviceOOBs 设备带外列表
	DeviceOOBs = "device_oob_list"
)

var (
	mTitle = map[string][]string{
		Device:        []string{`SN(必填)`, `业务负责人`, `负责人loginId{"comment":"导入时以此为准"}`, `品牌`, `型号`, `购买时间{"comment":"格式为2018-01-01"}`, `采购人`, `电源数量`, `分配状态{"comment":"unassigned:未分配,preassigned:预分配,assigned:已分配","validation":"unassigned,preassigned,assigned"}`, `带外账号`, `带外密码`, `设备等级{"comment":"1:低,2:中,3:高","validation":"1,2,3"}`, `备注`, `用途`, `规格`, `资产编号`, `业务IP{"comment":"格式为10.0.0.0,多IP用换行符分割"}`, `带外IP{"comment":"格式为10.0.0.0"}`, `操作系统`, `状态{"comment":"pre_deploy:等待部署,pre_online:等待上线,online:上线,offline:下线","validation":"pre_deploy,pre_online,online,offline"}`, `环境{"comment":"test:测试,development:开发,production:生产","validation":"test,development,production"}`, `硬件巡检状态{"comment":"只读无法导入"}`, `硬件巡检结果{"comment":"只读无法导入"}`, `管理人`, `管理人loginid{"comment":"导入时以此为准"}`, `硬件架构{"comment":"x86_64:x86_64架构,aarch64:aarch64架构,ppc64:ppc64架构","validation":"x86_64,aarch64,ppc64"}`, `数据中心`, `机房`, `机架(柜)`, `机位`, `标签`, `应用系统`, `SNMP端口`, `SNMP团体名`, `SNMP版本`, `SNMPv3用户名`, `SNMPv3认证协议`, `SNMPv3认证密码`, `SNMPv3加密协议`, `SNMPv3加密密码`},
		DeviceSetting: []string{`序列号(必填)`, `主机名`, `业务IP{"comment":"有3种格式I ip,mac  II ip,out/in, whichslot, whichport  III  ip,两个IP之间用换行符分割"}`, `硬件配置{"comment":"填写完整模版名称 比如Dell PowerEdge R730xd"}`, `PXE配置{"comment":"镜像安装不填此项"}`, `带外IP{"comment":"格式为10.0.0.0"}`, `镜像配置{"comment":"系统安装不填此项"}`, `操作系统用户名`, `操作系统密码`},
		DeviceOOBs:    []string{`序列号`, `厂商`, `型号`, `硬件架构`, `机房`, `位置`, `带外IP`, `带外用户`, `带外密码（密文）`, `带外密码修改时间`},
	}
)

// WriteToXLSX 将数据写入xlsx
func WriteToXLSX(which string, data [][]string) (*excelize.File, error) {
	return writeToXLSX(map[string][][]string{
		"Sheet1": data,
	}, mTitle[which], "")
}

// WriteToXLSXFile 将数据写入xlsx文件
func WriteToXLSXFile(which string, data [][]string) (*excelize.File, error) {
	return writeToXLSX(map[string][][]string{
		"Sheet1": data,
	}, mTitle[which], fmt.Sprintf("%s.xlsx", which))
}

// 读取excel
func ReadFromExcel(f string) (colums []map[string]interface{}, err error) {

	return
}
