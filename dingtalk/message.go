package dingtalk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Message https://open.dingtalk.com/document/orgapp/custom-robot-access#title-72m-8ag-pqw
type Message struct {
	MsgType    string     `json:"msgtype"`              // 消息类型
	Text       Text       `json:"text,omitempty"`       // text类型
	Link       Link       `json:"link,omitempty"`       // link类型
	Markdown   Markdown   `json:"markdown,omitempty"`   // markdown类型
	ActionCard ActionCard `json:"actionCard,omitempty"` // actionCard类型
	FeedCard   FeedCard   `json:"feedCard,omitempty"`   // feedCard类型
	At         At         `json:"at,omitempty"`
	WebhookURL string     `json:"-"` // 钉钉群的webhook地址
}

// Text text类型
type Text struct {
	Content string `json:"content"` // 消息内容
}

// Link link类型
type Link struct {
	Title      string `json:"title"`      // [必填] 消息标题
	Text       string `json:"text"`       // [必填] 消息内容。如果太长只会部分展示
	MessageURL string `json:"messageUrl"` // [必填] 点击消息跳转的URL
	PicURL     string `json:"picUrl"`     // 图片URL
}

// Markdown markdown类型
type Markdown struct {
	Title string `json:"title"` // [必填] 首屏会话透出的展示内容
	Text  string `json:"text"`  // [必填] markdown格式的消息
}

// ActionCard actionCard类型
type ActionCard struct {
	Title          string `json:"title"`                 // [必填] 首屏会话透出的展示内容
	Text           string `json:"text"`                  // [必填] markdown格式的消息
	SingleTitle    string `json:"singleTitle,omitempty"` // [必填] 单个按钮的标题
	SingleURL      string `json:"singleURL,omitempty"`   // [必填] 点击消息跳转的URL
	BtnOrientation string `json:"btnOrientation"`        // 0：按钮竖直排列 1：按钮横向排列
	Btns           []Btn  `json:"btns,omitempty"`        // 独立跳转ActionCard类型使用
}

// FeedCard feedCard类型
type FeedCard struct {
	Links []FeedCardLink `json:"links"`
}
type FeedCardLink struct {
	Title      string `json:"title"`      // [必填] 单条信息文本
	MessageURL string `json:"messageURL"` // [必填] 点击单条信息到跳转链接。
	PicURL     string `json:"picURL"`     // [必填] 单条信息后面图片的URL
}

// Btn 独立跳转actionCard类型中的按钮
type Btn struct {
	Title     string `json:"title"`     // 按钮标题
	ActionURL string `json:"actionURL"` // 点击按钮触发的URL
}

// At 被@的人
type At struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

// SendDingTalk 发送消息到钉钉群
func (msg Message) SendDingTalk() error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	client := http.Client{}
	req, err := http.NewRequest("POST", msg.WebhookURL, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("http request status code %v", res.StatusCode)
	}
	if err := parseError(res.Body); err != nil {
		return err
	}
	return nil
}

type response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func parseError(body io.Reader) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	var resp response
	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}
	switch resp.Errcode {
	case 0:
		return nil
	case 400013:
		return fmt.Errorf("群已被解散")
	case 400101:
		return fmt.Errorf("access_token不存在")
	case 400102:
		return fmt.Errorf("机器人已停用")
	case 400105:
		return fmt.Errorf("不支持的消息类型")
	case 400106:
		return fmt.Errorf("机器人不存在")
	case 410100:
		return fmt.Errorf("发送速度太快而限流")
	case 430101:
		return fmt.Errorf("含有不安全的外链")
	case 430102:
		return fmt.Errorf("含有不合适的文本")
	case 430103:
		return fmt.Errorf("含有不合适的图片")
	case 430104:
		return fmt.Errorf("含有不合适的内容")
	case 310000:
		return fmt.Errorf(resp.Errmsg)
	default:
		return fmt.Errorf("未知错误")
	}
}
