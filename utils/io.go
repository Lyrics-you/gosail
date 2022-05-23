package utils

import (
	"bufio"
	"encoding/json"
	"gosail/model"
	"io/ioutil"
	"os"
	"strings"
)

// var log = logger.Logger()

func SplitString(str string) (strList []string) {
	if str == "" {
		return
	}
	if strings.Contains(str, ",") {
		strList = strings.Split(str, ",")
	} else {
		strList = strings.Split(str, ";")
	}
	return
}

func SplitUserHost(str string) (user, host string) {
	if str == "" {
		return
	}
	if strings.Contains(str, "@") {
		user = strings.Split(str, "@")[0]
		host = strings.Split(str, "@")[1]
	} else {
		user = ""
		host = str
	}
	return
}

func GetByte(filepath string) ([]byte, error) {
	result, err := ioutil.ReadFile(filepath)
	if err != nil {
		// log.Errorf("read file %s error, %v", filepath, err)
		return nil, err
	}
	return result, nil
}

func GetString(filepath string) ([]string, error) {
	result := []string{}
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		// log.Errorf("read file %s error, %v", filepath, err)
		return nil, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	return result, nil
}

func GetJson(filePath string) (model.HostJson, error) {
	var result model.HostJson
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		// log.Errorf("read file %s error, %v", filePath, err)
		return result, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		// log.Errorf("read file %s error, %v", filePath, err)
		return result, err
	}
	return result, nil
}

func GetIpListFromString(ipString string) ([]string, error) {
	res := SplitString(ipString)
	var allIp []string
	if len(res) > 0 {
		for _, sip := range res {
			aip, err := ParseIp(sip)
			if err != nil {
				return aip, err
			}
			for _, ip := range aip {
				allIp = append(allIp, ip)
			}
		}
	}
	return allIp, nil
}

func GetIpListFromFile(filePath string) ([]string, error) {
	res, err := GetString(filePath)
	if err != nil {
		return nil, err
	}
	var allIp []string
	if len(res) > 0 {
		for _, sip := range res {

			aip, err := ParseIp(sip)
			if err != nil {
				return aip, err
			}
			for _, ip := range aip {
				allIp = append(allIp, ip)
			}
		}
	}
	return allIp, nil
}

func WriteIntoTxt(sshResult model.SSHResult, locate string) error {
	outputFile, outputError := os.OpenFile(locate+sshResult.Host+".txt", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		return outputError
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	//var outputString string

	outputString := sshResult.Result
	outputWriter.WriteString(outputString)
	outputWriter.Flush()
	return nil
}