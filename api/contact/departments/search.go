package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DeptSearchRequest 钉钉部门搜索请求结构体
type DeptSearchRequest struct {
	QueryWord string `json:"queryWord"` // 查询关键词
	Offset    int    `json:"offset"`    // 分页起始位置，1 表示第一页
	Size      int    `json:"size"`      // 每页大小
}

// DeptSearchResponse 钉钉部门搜索响应结构体
type DeptSearchResponse struct {
	HasMore    bool    `json:"hasMore"`    // 是否有更多数据
	TotalCount int     `json:"totalCount"` // 总记录数
	List       []int64 `json:"list"`       // 部门ID列表
}

// SearchDepartments 根据关键词查询部门ID列表
func SearchDepartments(accessToken, keyword string, offset, size int) (*DeptSearchResponse, error) {
	url := "https://api.dingtalk.com/v1.0/contact/departments/search"

	// 构造请求体
	reqBody := DeptSearchRequest{
		QueryWord: keyword,
		Offset:    offset,
		Size:      size,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP请求失败，状态码 %d，响应内容: %s", resp.StatusCode, string(data))
	}

	// 解析响应体
	var searchResp DeptSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &searchResp, nil
}
