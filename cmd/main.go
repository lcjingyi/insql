package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jingyi/insql/pkg"
)

func main() {
	var (
		url     = flag.String("url", "", "目标URL (例如: http://example.com/page.php?id=1)")
		method  = flag.String("method", "union", "注入方法: union, error, boolean, time")
		timeout = flag.Int("timeout", 10, "请求超时时间(秒)")
		verbose = flag.Bool("verbose", false, "详细输出")
	)
	
	// 使用verbose变量避免未使用警告
	_ = verbose
	_ = timeout
	
	flag.Parse()
	
	if *url == "" {
		fmt.Println("SQL注入检测工具")
		fmt.Println("用法:")
		flag.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  go run . -url=\"http://example.com/page.php?id=1\" -method=union")
		fmt.Println("  go run . -url=\"http://example.com/page.php?id=1\" -method=error")
		fmt.Println("  go run . -url=\"http://example.com/page.php?id=1\" -method=boolean")
		os.Exit(1)
	}
	
	// 验证URL格式
	if !strings.HasPrefix(*url, "http://") && !strings.HasPrefix(*url, "https://") {
		fmt.Println("错误: URL必须以http://或https://开头")
		os.Exit(1)
	}
	
	fmt.Printf("目标URL: %s\n", *url)
	fmt.Printf("检测方法: %s\n", *method)
	fmt.Printf("超时时间: %d秒\n", *timeout)
	fmt.Println("开始检测...")
	fmt.Println(strings.Repeat("-", 50))
	
	var err error
	switch strings.ToLower(*method) {
	case "union":
		err = pkg.UnionInject(*url)
	case "error":
		err = pkg.ErrorInject(*url)
	case "boolean":
		err = pkg.BooleanInject(*url)
	case "time":
		err = pkg.TimeBasedInject(*url)
	case "auto":
		err = pkg.Request(*url)
	default:
		fmt.Printf("错误: 不支持的检测方法 '%s'\n", *method)
		fmt.Println("支持的方法: union, error, boolean, time, auto")
		os.Exit(1)
	}
	
	if err != nil {
		fmt.Printf("检测过程中发生错误: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("检测完成!")
}
