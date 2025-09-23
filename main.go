package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 定义命令处理函数类型
type CommandFunc func(args []string) error

// 命令结构体
type Command struct {
	Name        string
	Description string
	Func        CommandFunc
}

// 命令注册表
var commands = make(map[string]Command)

// 注册命令
func registerCommand(cmd Command) {
	commands[cmd.Name] = cmd
}

// 初始化命令
func initCommands() {
	registerCommand(Command{
		Name:        "help",
		Description: "显示帮助信息",
		Func:        helpCommand,
	})
	registerCommand(Command{
		Name:        "exit",
		Description: "退出控制台",
		Func:        exitCommand,
	})
	registerCommand(Command{
		Name:        "info",
		Description: "显示系统信息",
		Func:        infoCommand,
	})
	registerCommand(Command{
		Name:        "set",
		Description: "设置变量 (用法: set <变量名> <值>)",
		Func:        setCommand,
	})
}

// 全局变量存储
var variables = make(map[string]string)

// 帮助命令处理函数
func helpCommand(args []string) error {
	fmt.Println("可用命令:")
	fmt.Println("====================")
	for _, cmd := range commands {
		fmt.Printf("%-10s %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

// 退出命令处理函数
func exitCommand(args []string) error {
	fmt.Println("退出控制台...")
	os.Exit(0)
	return nil
}

// 信息命令处理函数
func infoCommand(args []string) error {
	fmt.Println("安全控制台 v1.0")
	fmt.Println("一个类似msfconsole的交互式命令行工具")
	return nil
}

// 设置变量命令处理函数
func setCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: set <变量名> <值>")
	}

	varName := args[0]
	varValue := strings.Join(args[1:], " ")
	variables[varName] = varValue
	fmt.Printf("已设置 %s = %s\n", varName, varValue)
	return nil
}

// 解析并执行命令
func executeCommand(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	// 分割命令和参数
	parts := strings.Fields(input)
	cmdName := parts[0]
	args := parts[1:]

	// 查找并执行命令
	cmd, exists := commands[cmdName]
	if !exists {
		return fmt.Errorf("未知命令: %s. 输入 'help' 查看可用命令", cmdName)
	}

	return cmd.Func(args)
}

// 显示欢迎信息
func printWelcome() {
	fmt.Println("======================================")
	fmt.Println("        安全控制台 v1.0")
	fmt.Println("  输入 'help' 查看可用命令")
	fmt.Println("  输入 'exit' 退出控制台")
	fmt.Println("======================================")
}

func main() {
	// 初始化命令
	initCommands()

	// 显示欢迎信息
	printWelcome()

	// 创建扫描器读取用户输入
	scanner := bufio.NewScanner(os.Stdin)

	// 主循环
	for {
		// 显示命令提示符
		fmt.Print("console > ")

		// 读取输入
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		// 执行命令
		if err := executeCommand(input); err != nil {
			fmt.Printf("错误: %v\n", err)
		}
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "读取输入时出错: %v\n", err)
		os.Exit(1)
	}
}
