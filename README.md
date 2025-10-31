# InSQL - SQL注入检测工具

一个用Go语言编写的SQL注入检测工具，支持多种注入方式的检测和数据提取。

## 功能特性

- **联合注入 (Union-based)**: 通过UNION查询获取数据库信息
- **错误注入 (Error-based)**: 利用数据库错误信息提取数据
- **布尔盲注 (Boolean-based)**: 通过响应差异判断注入点
- **时间盲注 (Time-based)**: 通过响应时间判断注入点
- **自动检测**: 自动尝试多种注入方式

## 安装

```bash
git clone https://github.com/jingyi/insql.git
cd insql
go build -o insql.exe
```

## 使用方法

### 基本用法

```bash
# 使用联合注入检测
./insql.exe -url="http://example.com/page.php?id=1" -method=union

# 使用错误注入检测
./insql.exe -url="http://example.com/page.php?id=1" -method=error

# 使用布尔盲注检测
./insql.exe -url="http://example.com/page.php?id=1" -method=boolean

# 使用时间盲注检测
./insql.exe -url="http://example.com/page.php?id=1" -method=time

# 自动检测（推荐）
./insql.exe -url="http://example.com/page.php?id=1" -method=auto
```

### 命令行参数

- `-url`: 目标URL（必需）
- `-method`: 注入方法，可选值：union, error, boolean, time, auto（默认：union）
- `-timeout`: 请求超时时间，单位秒（默认：10）
- `-verbose`: 详细输出模式（默认：false）

### 示例

```bash
# 检测目标网站
./insql.exe -url="http://testphp.vulnweb.com/artists.php?artist=1" -method=auto

# 使用详细模式
./insql.exe -url="http://example.com/page.php?id=1" -method=union -verbose=true
```

## 支持的注入类型

### 1. 联合注入 (Union-based)
- 自动检测列数
- 提取数据库名、表名、列名
- 适用于有回显的注入点

### 2. 错误注入 (Error-based)
- 使用updatexml函数触发错误
- 从错误信息中提取数据
- 适用于显示错误信息的应用

### 3. 布尔盲注 (Boolean-based)
- 通过真/假条件响应差异判断
- 逐字符提取数据
- 适用于无回显但有响应差异的注入点

### 4. 时间盲注 (Time-based)
- 通过响应时间判断注入点
- 使用sleep函数控制响应时间
- 适用于完全盲注的情况

## 输出示例

```
目标URL: http://example.com/page.php?id=1
检测方法: union
超时时间: 10秒
开始检测...
--------------------------------------------------
Starting Union-based SQL injection...
Column count: 3
Database name: testdb
Table names: [users, products, orders]
Column names for table users: [id, username, password, email]
--------------------------------------------------
检测完成!
```

## 注意事项

1. **仅用于授权测试**: 请确保您有权限测试目标系统
2. **遵守法律法规**: 不要用于非法用途
3. **测试环境**: 建议在测试环境中使用
4. **网络延迟**: 时间盲注可能受网络延迟影响


