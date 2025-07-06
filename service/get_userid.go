package service

import (
	"DingApiGo/api/contact/users"
	"log"
)

// GetUserID 获取，从用户名关键字获取id
// - accessToken
func GetUserID(accessToken, keyword string, fullMatch bool) []string {

	// 查询用户id
	resp, err := users.SearchUsers(accessToken, keyword, fullMatch, 0, 10)
	if err != nil {
		log.Fatalf("搜索用户失败: %v", err)
	}

	if resp.TotalCount == 0 || len(resp.List) == 0 {
		log.Fatalf("未找到匹配的用户: %s", keyword)
	}

	// 返回第一个匹配用户的 userId
	return resp.List
}
