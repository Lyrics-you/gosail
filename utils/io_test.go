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
	return
}

func Test_io_SplitUserHost(t *testing.T) {
	str := "192.168.245.131"
	user, host := SplitUserHost(str)
	fmt.Println(user, host)
	t.Log()
	return
}

func Test_io_GetByte(t *testing.T) {
	result, _ := GetByte("C:\\Users\\Taragrade\\.ssh\\id_rsa.pub")
	fmt.Println(result)
	t.Log(result)
	return
}

func Test_io_GetString(t *testing.T) {
	result, _ := GetString("..\\examples\\echodate")
	fmt.Println(result)
	t.Log(result)
	return
}

func Test_io_GetJson(t *testing.T) {
	result, _ := GetJson("ssh.json")
	fmt.Println(result)
	t.Log(result)
	return
}
