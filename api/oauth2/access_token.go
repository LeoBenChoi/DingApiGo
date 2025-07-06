// package oauth2 用于获取 Access Token
package oauth2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 钉钉 API 响应结构体
type AccessTokenResponse struct {
	AccessToken string `json:"accessToken"` // accessToken
	ExpireIn    int    `json:"expireIn"`    // 过期时间（秒）
}

// 全局变量存储 accessToken 和过期时间
var (
	token      string
	expireAt   time.Time
	tokenMutex sync.Mutex
)

// GetAccessToken 获取有效的 accessToken，自动刷新过期的 token
//
// 参数：
//
// - appKey:     钉钉开放平台应用的 AppKey。
//
// - appSecret:  钉钉开放平台应用的 AppSecret，用于鉴权。
//
// 返回：
//
// - string:     获取到的 accessToken 字符串。
//
// - error:      如果请求或解析失败，返回错误信息。
func GetAccessToken(appKey, appSecret string) (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// 如果 token 还没过期，直接返回
	if time.Now().Before(expireAt) && token != "" {
		return token, nil
	}

	// 钉钉 accessToken 接口 URL
	url := "https://api.dingtalk.com/v1.0/oauth2/accessToken"

	// 构建请求体
	requestBody := map[string]string{
		"appKey":    appKey,
		"appSecret": appSecret,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("构建请求体失败: %v", err)
	}

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("请求access token失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 如果钉钉返回错误
	if result.AccessToken == "" {
		return "", fmt.Errorf("钉钉返回错误: 没有获取到 accessToken")
	}

	// 缓存 token 和过期时间
	token = result.AccessToken
	expireAt = time.Now().Add(time.Duration(result.ExpireIn-30) * time.Second) // 提前 30 秒刷新

	return token, nil
}
