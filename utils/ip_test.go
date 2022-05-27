package utils

import (
	"fmt"
	"testing"
)

func Test_ip_ParseIp(t *testing.T) {
	str := "1.1.1.1/32"
	result, _ := ParseIp(str)
	fmt.Println(result)
	str = "1.1.1.1/32-1.1.1.2/32"
	result, _ = ParseIp(str)
	fmt.Println(result)
	t.Log(result)
}

func Test_ip_GetAvailableIP(t *testing.T) {
	str := "191.168.245.131/29"
	result, _ := GetAvailableIP(str)
	fmt.Println(result)
	t.Log(result)

}

func Test_ip_IPAddressToCIDR(t *testing.T) {
	str := "191.168.245.131/29"
	result := IPAddressToCIDR(str)
	fmt.Println(result)
	t.Log(result)
}
