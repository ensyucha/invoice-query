package main

import (
	"github.com/kataras/iris"
	"invoice/dbop"
	"invoice/server"
	"log"
)

func main() {

	// 新建服务器
	app := server.NewApp()

	log.Println("监听地址   : http://" + dbop.GetIP() + dbop.GetPort() + "\n")

	// 启动服务器
	log.Fatal(app.Run(iris.Addr(dbop.GetPort()), iris.WithoutStartupLog))
}
