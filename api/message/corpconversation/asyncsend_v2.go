// package corpconversation 用于发送通知消息
package corpconversation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ==== 请求体结构 ====

// MessageBody 通用消息体封装，子字段根据 msgtype 填充
type MessageBody struct {
	MsgType string      `json:"msgtype"` // 消息类型：text、image、voice、file、link、markdown、action_card、oa
	Content interface{} `json:"-"`
}

// MarshalJSON 自定义序列化，将 content 嵌入到 msg 下
func (m MessageBody) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		MsgType string      `json:"msgtype"`
		Msg     interface{} `json:"-"`
	}
	// 用 map 构造动态字段
	base := map[string]interface{}{
		"msgtype": m.MsgType,
		m.MsgType: m.Content, // e.g. "text": TextContent
	}
	return json.Marshal(base)
}

// TextContent 文本消息内容
type TextContent struct {
	Content string `json:"content"` // 文本内容
}

type ImageContent struct {
	MediaID string `json:"media_id"`
}

type VoiceContent struct {
	MediaID  string `json:"media_id"`
	Duration string `json:"duration"` // 注意类型为 string
}

type FileContent struct {
	MediaID string `json:"media_id"`
}

type LinkContent struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageURL string `json:"messageUrl"`
	PicURL     string `json:"picUrl"`
}

type MarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ActionCardContent struct {
	Title       string `json:"title"`
	Markdown    string `json:"markdown"`
	SingleTitle string `json:"single_title"`
	SingleURL   string `json:"single_url"`
}

type OAContent struct {
	MessageURL   string `json:"message_url"`
	PCMessageURL string `json:"pc_message_url"`
	Head         OAHead `json:"head"`
	Body         OABody `json:"body"`
}

type OAHead struct {
	BgColor string `json:"bgcolor"`
	Text    string `json:"text"`
}

type OABody struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Author    string   `json:"author"`
	Image     string   `json:"image"`
	FileCount string   `json:"file_count"`
	Form      []OAForm `json:"form"`
	Rich      OARich   `json:"rich"`
}

type OAForm struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OARich struct {
	Num  string `json:"num"`
	Unit string `json:"unit"`
}

// AsyncSendV2Request 异步发送工作通知请求体
type AsyncSendV2Request struct {
	AgentID       int64       `json:"agent_id"`                  // 必填：应用的 AgentID
	ToAllUser     bool        `json:"to_all_user"`               // 是否全部用户
	UserIDList    string      `json:"userid_list,omitempty"`     // 用户 ID 列表，逗号分隔
	DeptIDList    string      `json:"dept_id_list,omitempty"`    // 部门 ID 列表，逗号分隔
	Msg           MessageBody `json:"msg"`                       // 消息体
	EnableIDTrans bool        `json:"enable_id_trans,omitempty"` // 是否开启 ID 转换
}

// ==== 响应体结构 ====

// AsyncSendV2Response 异步发送工作通知响应
type AsyncSendV2Response struct {
	ErrCode   int    `json:"errcode"`    // 错误码，0 表示成功
	ErrMsg    string `json:"errmsg"`     // 错误信息
	TaskID    int64  `json:"task_id"`    // 任务 ID，可用于后续查询
	RequestID string `json:"request_id"` // 请求唯一标识
}

// ==== 发送函数 ====

// SendAsyncMessageV2 调用钉钉异步发送工作通知消息接口（v2 版本）
//
// accessToken: 钉钉接口 access_token
// reqBody:     构造的请求体
//
// 返回 AsyncSendV2Response 及 error
func SendAsyncMessageV2(accessToken string, reqBody AsyncSendV2Request) (*AsyncSendV2Response, error) {
	// 构造 URL，access_token 作为 query 参数
	url := fmt.Sprintf(
		"https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=%s",
		accessToken,
	)

	// 序列化请求体
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("请求体序列化失败: %v", err)
	}

	// 构造 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP 状态码 %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应体
	var result AsyncSendV2Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 业务错误
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("钉钉返回错误: %s", result.ErrMsg)
	}

	return &result, nil
}
