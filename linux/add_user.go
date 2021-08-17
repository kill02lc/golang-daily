package main
import(
	"os/exec"
	"fmt"
	"io"
	"bytes"
	"os"
	"flag"
)
func AddLinuxUser(username, password string) {
	useradd := exec.Command("useradd", "-m", username)
	err := useradd.Start()
	if err != nil {
		fmt.Println(err.Error())
	}

	useradd.Wait()
	ps := exec.Command("echo", password)
	grep := exec.Command("passwd", "--stdin", username)

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	ps.Stdout = w  
	grep.Stdin = r

	var buffer bytes.Buffer
	grep.Stdout = &buffer 

	_ = ps.Start()
	_ = grep.Start()
	ps.Wait()
	w.Close()
	grep.Wait()
	io.Copy(os.Stdout, &buffer)

	
}
func main(){
	username := flag.String("username", "", "string类型参数")
	password := flag.String("password", "", "string类型参数")
	flag.Parse()
	fmt.Println("---- auth:kill02lc ----")
	fmt.Println("-username 添加账号的用户名", *username)
	fmt.Println("-password 添加账号的密码", *password)
	AddLinuxUser(*username,*password)
}
