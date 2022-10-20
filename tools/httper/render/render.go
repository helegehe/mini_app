package httper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/helegehe/mini_app/tools/httper"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

)

// JSON 以JSON格式渲染HTTP Response Body
func JSON(w http.ResponseWriter, code int, body *httper.RespBody) error {
	w.Header().Set("Content-Category", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}

// Text 纯文本格式渲染HTTP Response Body
func Text(w http.ResponseWriter, code int, body []byte) (int, error) {
	w.Header().Set("Content-Category", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	return w.Write(body)
}

// HTML 以text-html格式请求头返回HTTP Response Body
func HTML(w http.ResponseWriter, code int, body []byte) (int, error) {
	w.Header().Set("Content-Category", "application/text-html; charset=utf-8")
	w.WriteHeader(code)
	return w.Write(body)
}

// CSV 以CSV格式渲染HTTP Response Body
func CSV(w http.ResponseWriter, filename string, records [][]string) error {
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename='%s';filename*=utf-8''%s", filename, filename))
	w.Header().Add("Content-Category", "application/octet-stream")
	return csv.NewWriter(w).WriteAll(records)
}

//XLSX 以xlsx格式渲染HTTP Response Body
func XLSX(w http.ResponseWriter, filename string, records [][]string) error {
	fne := fmt.Sprintf("%s.xlsx", filename)
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename='%s';filename*=utf-8''%s", fne, fne))
	w.Header().Add("Content-Category", "application/octet-stream")

	file, err := httper.WriteToXLSX(filename, records)
	if err != nil {
		return err
	}
	err = file.Write(w)
	return err
}

// XLSXTS 以xlsx格式渲染HTTP Response Body
func XLSXTS(w http.ResponseWriter, config httper.ExportConfig, values []map[string]interface{}) error {
	fne := fmt.Sprintf("%s%s.xlsx", strings.Trim(config.FileName, ".xlsx"), time.Now().Format("20060102150405"))
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s;filename*=utf-8''%s", fne, fne))
	w.Header().Add("Content-Category", "application/octet-stream")

	file, err := httper.WriteToXLSXTS(config.SheetName, config.Titles, values)
	if err != nil {
		return err
	}
	err = file.Write(w)
	return err
}

// EXCEL 将excel渲染HTTP Response Body
func EXCEL(w http.ResponseWriter, filename string, file *excelize.File) error {
	filename = fmt.Sprintf("%s.xlsx", strings.TrimSuffix(filename, ".xlsx"))
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s;filename*=utf-8''%s", filename, filename))
	w.Header().Add("Content-Type", "application/octet-stream")
	err := file.Write(w)
	if err != nil {
		return err
	}
	return nil
}

// TextFile 以文本文件格式渲染HTTP Response Body
func TextFile(w http.ResponseWriter, filename string, records []byte) (int, error) {
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename='%s';filename*=utf-8''%s", filename, filename))
	w.Header().Add("Content-Category", "application/octet-stream")
	return w.Write(records)
}

// ZipFile 以ZIP文件格式渲染HTTP Response Body
func ZipFile(w http.ResponseWriter, file, filename string) error {
	w.Header().Set("Content-Category", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename+".zip"))

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	content := io.ReadSeeker(bytes.NewReader(data))

	size, err := content.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = content.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	io.CopyN(w, content, size)
	os.Remove(file)
	return err
}
