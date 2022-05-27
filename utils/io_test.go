package utils

import (
	"fmt"
	"testing"
)

func Test_io_SplitString(t *testing.T) {
	str := "1,2,3"
	result := SplitString(str)
	fmt.Println(result)
	t.Log(result)

}

func Test_io_SplitUserHost(t *testing.T) {
	str := "192.168.245.131"
	user, host := SplitUserHost(str)
	fmt.Println(user, host)
	t.Log()
}

func Test_io_SplitUserHostPath(t *testing.T) {
	strs := []string{
		"root@192.168.245.131:/root/demo/",
		"root@192.168.245.131",
		"192.168.245.131:/root/demo/",
		"/root/demo/",
	}
	for _, str := range strs {
		user, host, path := SplitUserHostPath(str)
		fmt.Printf("%+q %+q %+q\n", user, host, path)
	}
	t.Log()
}

func Test_io_GetByte(t *testing.T) {
	result, _ := GetByte("C:\\Users\\Taragrade\\.ssh\\id_rsa.pub")
	fmt.Println(result)
	t.Log(result)

}

func Test_io_GetString(t *testing.T) {
	result, _ := GetString("..\\examples\\echodate")
	fmt.Println(result)
	t.Log(result)

}

func Test_io_GetJson(t *testing.T) {
	result, _ := GetJson("ssh.json")
	fmt.Println(result)
	t.Log(result)

}

func Test_info_GetAbsPath(t *testing.T) {
	strs := []string{
		"/root/demo/",
		"root@192.168.245.131",
		"192.168.245.131:/root/demo/",
		"/root/demo/",
	}
	for _, str := range strs {
		user, host, path := SplitUserHostPath(str)
		fmt.Printf("%+q %+q %+q\n", user, host, path)
	}
	t.Log()
}
