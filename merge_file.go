package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//var filename_tong = "./test.txt"

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
func write_file(file_cont string, filename_tong string) bool {
	//	var str = "测试1\n测试2\n"
	var f *os.File
	var err1 error
	if checkFileIsExist(filename_tong) { //如果文件存在
		f, err1 = os.OpenFile(filename_tong, os.O_APPEND, 0666) //打开文件
		//	fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename_tong) //创建文件
		//fmt.Println("文件不存在")
	}
	defer f.Close()
	if err1 != nil {
		panic(err1)
	}
	w := bufio.NewWriter(f) //创建新的 Writer 对象
	w.WriteString(file_cont)
	//fmt.Printf("写入 %d 个字节n", n)
	w.Flush()
	return true
}
func main() {
	var path_c string
	var save_c string
	fmt.Println("请输入要文件整合的路径:")
	//当程序只是到fmt.Scanln(&name)程序会停止执行等待用户输入
	fmt.Scanln(&path_c)
	fmt.Println("请输入要文件保存名称:")
	fmt.Scanln(&save_c)
	save_c = path_c + "\\" + save_c
	fmt.Print(save_c)
	//pwd, _ := os.Getwd()
	//	fmt.Print(pwd)
	//var file_add []string
	//获取当前目录下的所有文件或目录信息
	//var temp1 string
	filepath.Walk(path_c, func(path string, info os.FileInfo, err error) error {
		//	fmt.Println(path)        //打印path信息

		if strings.Contains(info.Name(), ".") && !(strings.Contains(info.Name(), ".go")) && (info.Name() != "test.txt") { //筛选字符串有.
			//		fmt.Println(info.Name()) //打印文件
			fi, err := os.Open(info.Name())
			//	fmt.Println(info.Name())

			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			defer fi.Close()

			br := bufio.NewReader(fi)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}

				var w string = string(a)
				w = w + "\n"
				fmt.Print(w)
				write_file(w, save_c)
			}
			// file, err := os.Open(info.Name())
			// if err != nil {
			// 	panic(err)
			// }
			// defer file.Close()
			// content, err := ioutil.ReadAll(file)

		}

		return nil
	})

}
