// service 发送消息
package service

import (
	"fmt"
	"strings"
	"time"

	cc "DingApiGo/api/message/corpconversation"
	"DingApiGo/config"
)

// prependTimestamp 为消息内容添加时间戳前缀
func prependTimestamp(text string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] %s", timestamp, text)
}

// SendTextNotification 发送文本消息给指定用户列表
//
// 参数：
//
// - content: 文本消息内容
//
// - userIDs: 接收者用户 ID 列表
func SendTextNotification(content string, userIDs []string) error {
	msg := cc.MessageBody{
		MsgType: "text",
		Content: cc.TextContent{Content: prependTimestamp(content)},
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendImageNotification 发送图片消息
//
// 参数：
//
// - mediaID: 钉钉媒体文件 ID（图片）
//
// - userIDs: 接收者用户 ID 列表
func SendImageNotification(mediaID string, userIDs []string) error {
	msg := cc.MessageBody{
		MsgType: "image",
		Content: cc.ImageContent{MediaID: mediaID},
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendVoiceNotification 发送语音消息
//
// 参数：
//
// - mediaID: 媒体文件 ID（语音）
//
// - duration: 语音时长（秒）
//
// - userIDs: 接收者用户 ID 列表
func SendVoiceNotification(mediaID string, duration int, userIDs []string) error {
	msg := cc.MessageBody{
		MsgType: "voice",
		Content: cc.VoiceContent{MediaID: mediaID, Duration: fmt.Sprintf("%d", duration)},
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendFileNotification 发送文件消息
//
// 参数：
//
// - mediaID: 钉钉文件媒体 ID
//
// - userIDs: 接收者用户 ID 列表
func SendFileNotification(mediaID string, userIDs []string) error {
	msg := cc.MessageBody{
		MsgType: "file",
		Content: cc.FileContent{MediaID: mediaID},
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendLinkNotification 发送链接消息
//
// 参数：
//
// - title: 消息标题
//
// - text: 消息正文
//
// - messageURL: 跳转链接
//
// - picURL: 链接附带的图片 URL
//
// - userIDs: 接收者用户 ID 列表
func SendLinkNotification(title, text, messageURL, picURL string, userIDs []string) error {
	msg := cc.MessageBody{
		MsgType: "link",
		Content: cc.LinkContent{
			Title:      prependTimestamp(title),
			Text:       prependTimestamp(text),
			MessageURL: messageURL,
			PicURL:     picURL,
		},
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendMarkdownNotification 发送 Markdown 格式消息
// 参数：
// - title: Markdown 标题
// - markdown: Markdown 正文
// - userIDs: 接收者用户 ID 列表
func SendMarkdownNotification(title, markdown string, userIDs []string) error {
	msg := cc.MessageBody{
		MsgType: "markdown",
		Content: cc.MarkdownContent{
			Title: prependTimestamp(title),
			Text:  prependTimestamp(markdown),
		},
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendActionCardNotification 发送单按钮卡片消息
//
// 参数：
//
// - title: 卡片标题
//
// - text: Markdown 格式正文
//
// - singleTitle: 按钮文字
//
// - singleURL: 按钮跳转链接
//
// - userIDs: 接收者用户 ID 列表
func SendActionCardNotification(title, text, singleTitle, singleURL string, userIDs []string) error {
	msg := cc.MessageBody{
		MsgType: "action_card",
		Content: cc.ActionCardContent{
			Title:       prependTimestamp(title),
			Markdown:    prependTimestamp(text),
			SingleTitle: singleTitle,
			SingleURL:   singleURL,
		},
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendOANotification 发送 OA 消息，需要传入完整的 OAContent
//
// 参数：
//
// - oa: 完整的 OAContent 结构体
//
// - userIDs: 接收者用户 ID 列表
func SendOANotification(oa cc.OAContent, userIDs []string) error {
	// OA 消息内容结构复杂，不统一加时间戳，调用方自行处理
	msg := cc.MessageBody{
		MsgType: "oa",
		Content: oa,
	}
	return SendMsg(msg, userIDs, nil, false)
}

// SendMsg 内部通用发送函数
//
// 参数：
//
// - msg: 消息结构体（已封装）
//
// - userIDs: 接收者用户 ID 列表
//
// - deptIDs: 接收者部门 ID 列表
//
// - toAll: 是否发送给所有人
func SendMsg(msg cc.MessageBody, userIDs, deptIDs []string, toAll bool) error {
	// 获取 accessToken
	accessToken := GetAccessToken()

	// 读取 AgentID
	cfg := config.GetConfig()

	req := cc.AsyncSendV2Request{
		AgentID:       cfg.DingDing.AgentID,
		ToAllUser:     toAll,
		UserIDList:    strings.Join(userIDs, ","),
		DeptIDList:    strings.Join(deptIDs, ","),
		Msg:           msg,
		EnableIDTrans: false,
	}

	resp, err := cc.SendAsyncMessageV2(accessToken, req)
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}
	fmt.Printf("发送成功，任务ID: %d\n", resp.TaskID)
	return nil
}
