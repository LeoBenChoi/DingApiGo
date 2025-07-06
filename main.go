package main

import (
	_ "DingApiGo/api/contact/departments"
	_ "DingApiGo/api/oauth2"
	"DingApiGo/config"
	_ "DingApiGo/config"
	"DingApiGo/service"
	_ "DingApiGo/service"
	"fmt"
)

func main() {
	// 加载 config 配置
	config.LoadConfig()

	// 获取 token
	accessToken := service.GetAccessToken()
	fmt.Println(accessToken)

	// 搜索 userid
	userid := service.GetUserID(accessToken, "c", false)
	fmt.Println(userid)

	service.GetAccessToken()

	contact := "测试"
	service.SendTextNotification(contact, userid)
}
