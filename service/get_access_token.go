package service

import (
	"DingApiGo/api/oauth2"
	"DingApiGo/config"
	"log"
)

// GetAccessTokenFromService 获取 AccessToken（封装 api 层）
// 推荐在调用前已经调用 config.LoadConfig() 读取配置
func GetAccessToken() string {
	cfg := config.GetConfig()

	token, err := oauth2.GetAccessToken(cfg.DingDing.AppKey, cfg.DingDing.AppSecret)
	if err != nil {
		log.Fatalf("获取 AccessToken 失败: %v", err)
	}

	return token
}
