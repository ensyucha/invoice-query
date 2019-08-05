package controller

import (
	"github.com/kataras/iris"
	"invoice/auth"
	"invoice/dbop"
	"invoice/model"
)

func IndexAdd(ctx iris.Context) {

	auth.CheckToken(ctx)

	username := ctx.GetCookie("username")

	if username == "admin" {
		ctx.RemoveCookie("token")
		ctx.RemoveCookie("username")
		ctx.Redirect("/", 302)
	}

	user := &model.User{Username:username}

	usage, err := dbop.UCGetUserUsage(user)

	if err != nil {
		ctx.ViewData("Usage", err.Error())
	} else {
		ctx.ViewData("Usage", usage)
	}

	nickName := ctx.GetCookie("nickname")

	ctx.ViewData("NickName", nickName)


	if err := ctx.View("add.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}

func AddData(ctx iris.Context) {

	auth.CheckToken(ctx)

	username := ctx.GetCookie("username")

	if queryArray, ok := getQueryArrayJSON(ctx, "读取查询组信息"); ok {
		_, _ = ctx.JSON(dbop.AddDataToDB(queryArray, &model.User{Username:username}))
	}
}