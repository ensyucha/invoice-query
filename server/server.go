package server

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"invoice/controller"
)

func NewApp() *iris.Application {

	app := iris.New()

	registerHTMLTemplate(app)
	registerMiddleware(app)
	registerRouter(app)

	return app
}

func registerHTMLTemplate(app *iris.Application) {
	app.RegisterView(iris.HTML("./assets", ".html"))
}

func registerMiddleware(app *iris.Application) {
	app.Use(recover.New()) // 使用恢复中间件
	app.Use(logger.New())  // 使用日志中间件
	app.Logger().SetLevel("warn")
}

func registerRouter(app *iris.Application) {

	// 根路径重定向到登录页面
	app.Get("/", controller.IndexLogin)

	// 登录与登出
	app.Post("/login", controller.Login)
	app.Get("/logout", controller.Logout)

	// 新增路由
	app.Get("/add", controller.IndexAdd)
	app.Post("/add", controller.AddData)

	// 查询路由
	app.Post("/query", controller.AcceptQuery)

	// 归档封存路由
	app.Get("/result", controller.IndexResult)
	app.Post("/result/data", controller.ProcessResultData)
	app.Get("/result/sealed", controller.IndexSealed)

	// 系统管理路由
	app.Get("/system", controller.IndexSystem)
	app.Get("/system/user", controller.ListUser)
	app.Get("/system/outputalldata", controller.OutputAllData)
	app.Post("/system/user", controller.AddUser)
	app.Put("/system/user", controller.UpdateUser)
	app.Delete("/system/user", controller.RemoveUser)
	app.Post("/system/password", controller.UpdateSystemPassword)
	app.Post("/system/unusedusage", controller.UpdateSystemUnusedUsage)
	app.Post("/system/apicode", controller.UpdateSystemApiCode)

	// 静态文件服务
	app.StaticWeb("/assets", "./assets")
}
