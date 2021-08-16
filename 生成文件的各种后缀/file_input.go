package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
func MakeFile(filepath string, ext_str string, file_name_olny_new string) {
	fmt.Println(filepath)
	fmt.Println(ext_str)

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	filenameWithSuffix := path.Base(filepath)
	fileSuffix := path.Ext(filenameWithSuffix)
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	var filenameOnly_new string
	fmt.Println(filenameOnly)
	if file_name_olny_new != "" {
		filenameWithSuffix_new := path.Base(file_name_olny_new)
		fileSuffix_new := path.Ext(filenameWithSuffix_new)
		filenameOnly_new = strings.TrimSuffix(filenameWithSuffix_new, fileSuffix_new)
	}
	err_dir := os.Mkdir(filenameOnly+"_ext", os.ModePerm)

	if err_dir != nil {
		fmt.Println(err_dir)
	}
	for _, v := range strings.SplitAfter(ext_str, ",") {
		trimStr := strings.Trim(v, ",")
		if file_name_olny_new != "" {
			err_w := ioutil.WriteFile(filenameOnly+"_ext"+"\\"+filenameOnly_new+"."+trimStr, content, 0644)
			fmt.Println(filenameOnly + "_ext" + "\\" + filenameOnly_new + "." + trimStr + " 生成成功!")
			if err_w != nil {
				panic(err_w)
			}
		} else {
			err_w := ioutil.WriteFile(filenameOnly+"_ext"+"\\"+filenameOnly+"."+trimStr, content, 0644)
			fmt.Println(filenameOnly + "_ext" + "\\" + filenameOnly + "." + trimStr + " 生成成功!")
			if err_w != nil {
				panic(err_w)
			}
		}

	}
}

//
func main() {
	add_suffix_file := flag.String("add_suffix_file", "", "string类型参数")
	add_suffix_file_all := flag.String("add_suffix_file_all", "", "string类型参数")
	add_suffix_str := flag.String("add_suffix_str", "", "string类型参数")
	flag.Parse()
	fmt.Println("---- auth:kill02lc ----")
	fmt.Println("-add_suffix_file 添加后缀的文件路径", *add_suffix_file)
	fmt.Println("-add_suffix_file_all 添加后缀的路径(添加底下所有文件的后缀)", *add_suffix_file_all)
	fmt.Println("-add_suffix_str 后缀字符串(参数格式:word,doc,docx,pdf,jpg,gif,rar,zip,ppt,ppt,7z)", *add_suffix_str)
	fmt.Println("\n")
	fmt.Println("使用样例:")
	fmt.Println(".\\file_input.exe -add_suffix_file yhk.txt -add_suffix_str word,doc,docx,pdf,jpg,gif,rar,zip,ppt,ppt,7z")
	fmt.Println(".\\file_input.exe -add_suffix_file E:\\temp\\ -add_suffix_str word,doc,docx,pdf,jpg,gif,rar,zip,ppt,ppt,7z")

	if (*add_suffix_file) != "" && (*add_suffix_str) != "" {
		if checkFileIsExist(*add_suffix_file) {
			filepath := *add_suffix_file
			ext_str := *add_suffix_str
			MakeFile(filepath, ext_str, "")
		} else {
			fmt.Println("文件不存在")
		}
	} else if (*add_suffix_file_all) != "" && (*add_suffix_str) != "" {
		filepath.Walk(*add_suffix_file_all, func(path string, info os.FileInfo, err error) error {

			if strings.Contains(info.Name(), ".") {
				MakeFile(*add_suffix_file_all+info.Name(), *add_suffix_str, info.Name())
			}
			return nil
		})

	} else {
		fmt.Println("\n")
		fmt.Println("missing parameters")
	}

}
