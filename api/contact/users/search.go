// package users 搜索用户 userId
package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserSearchRequest struct {
	FullMatchField int    `json:"fullMatchField,omitempty"` // 1 表示精确匹配，可省略表示模糊
	QueryWord      string `json:"queryWord"`                // 关键字
	Offset         int    `json:"offset"`                   // 分页偏移量
	Size           int    `json:"size"`                     // 返回条数
}

type UserInfo struct {
	UserID string `json:"userId"`
}

type UserSearchResponse struct {
	HasMore    bool     `json:"hasMore"`
	TotalCount int      `json:"totalCount"`
	List       []string `json:"list"`
}

// SearchUsersID 搜索用户 ID
//
// 参数：
//
// - accessToken 是钉钉的访问令牌。
//
// - queryWord 是查询的关键词，例如用户名或手机号。
//
// - fullMatch 设置为 true 表示精确匹配，false 表示模糊匹配（可选）。
//
// - offset 是分页起始位置, 表示偏移量，最好为0。
//
// - size 是返回的用户数量。
//
// 返回用户搜索结果 UserSearchResponse。
func SearchUsers(accessToken, queryWord string, fullMatch bool, offset, size int) (*UserSearchResponse, error) {
	url := "https://api.dingtalk.com/v1.0/contact/users/search"

	// 构造请求体
	reqBody := UserSearchRequest{
		QueryWord: queryWord,
		Offset:    offset,
		Size:      size,
	}
	if fullMatch {
		reqBody.FullMatchField = 1
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("请求体序列化失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP 状态码 %d，响应内容: %s", resp.StatusCode, string(data))
	}

	var result UserSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}
