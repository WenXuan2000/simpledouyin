package main

import (
	"github.com/gin-gonic/gin"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

func main() {
	err := dao.InitMySql()
	if err != nil {
		panic(err)
	}
	//程序退出关闭数据库连接
	defer dao.Close()
	//绑定模型
	//dao.SqlSession.AutoMigrate(&entity.User{})
	dao.SqlSession.AutoMigrate(&entity.Follow{})
	//注册路由

	r := gin.Default()
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	//10.0.2.2是一个特殊IP，在模拟器上用10.0.2.2就是访问你的电脑本机
}
