package godingtalk

import (
	"net/url"
	"strconv"
)

type MsgType struct {
	Type string      `json:"msgtype"`
	Text interface{} `json:"text"`
	Link interface{} `json:"link"`
}

type Text struct {
	Content string `json:"content"`
}

type Link struct {
	MessageUrl string `json:"messageUrl"`
	PicUrl     string `json:"picUrl"`
	Title      string `json:"title"`
	Text       string `json:"text"`
}

/***********************************************工作通知************************************************************/

type MsgWork struct {
	AgentId    string      `json:"agent_id"`
	UserIdList string      `json:"userid_list"`
	Msg        interface{} `json:"msg"`
}

//SendAppMessage is 发送企业会话消息
func (c *DingTalkClient) SendWorkTextMessage(agentID string, touser string, msg string) (data OAPIResponse, err error) {
	request := MsgWork{
		AgentId:    agentID,
		UserIdList: touser,
		Msg: MsgType{
			Type: "text",
			Text: Text{
				Content: msg,
			},
		},
	}
	err = c.httpRPC("topapi/message/corpconversation/asyncsend_v2", nil, request, &data)
	return
}

//SendAppLinkMessage is 发送企业会话链接消息
func (c *DingTalkClient) SendWorkLinkMessage(agentID, touser string, title, text string, picUrl, url string) (data OAPIResponse, err error) {
	request := MsgWork{
		AgentId:    agentID,
		UserIdList: touser,
		Msg: MsgType{
			Type: "link",
			Link: Link{
				MessageUrl: url,
				PicUrl:     picUrl,
				Title:      title,
				Text:       text,
			},
		},
	}
	err = c.httpRPC("topapi/message/corpconversation/asyncsend_v2", nil, request, &data)
	return
}

/***************************************************对私信息**************************************************************/
type MsgUser struct {
	Sender string      `json:"sender"`
	Cid    string      `json:"cid"`
	Msg    interface{} `json:"msg"`
}

//SendTextMessage is 发送普通文本消息
func (c *DingTalkClient) SendUserTextMessage(sender string, cid string, msg string) (data OAPIResponse, err error) {
	request := MsgUser{
		Sender: sender,
		Cid:    cid,
		Msg: MsgType{
			Type: "text",
			Text: Text{
				Content: msg,
			},
		},
	}
	err = c.httpRPC("message/send_to_conversation", nil, request, &data)
	return
}

/*************************************************OA消息*******************************************************/

//OAMessage is the Message for OA
type OAMessage struct {
	URL   string `json:"message_url"`
	PcURL string `json:"pc_message_url"`
	Head  struct {
		BgColor string `json:"bgcolor,omitempty"`
		Text    string `json:"text,omitempty"`
	} `json:"head,omitempty"`
	Body struct {
		Title     string          `json:"title,omitempty"`
		Form      []OAMessageForm `json:"form,omitempty"`
		Rich      OAMessageRich   `json:"rich,omitempty"`
		Content   string          `json:"content,omitempty"`
		Image     string          `json:"image,omitempty"`
		FileCount int             `json:"file_count,omitempty"`
		Author    string          `json:"author,omitempty"`
	} `json:"body,omitempty"`
}

type OAMessageForm struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type OAMessageRich struct {
	Num  string `json:"num,omitempty"`
	Unit string `json:"body,omitempty"`
}

func (m *OAMessage) AppendFormItem(key string, value string) {
	f := OAMessageForm{Key: key, Value: value}

	if m.Body.Form == nil {
		m.Body.Form = []OAMessageForm{}
	}

	m.Body.Form = append(m.Body.Form, f)
}

//SendOAMessage is 发送OA消息
func (c *DingTalkClient) SendOAMessage(sender string, cid string, msg OAMessage) (data MessageResponse, err error) {
	request := map[string]interface{}{
		"chatid":  cid,
		"sender":  sender,
		"msgtype": "oa",
		"oa":      msg,
	}
	err = c.httpRPC("chat/send", nil, request, &data)
	return data, err
}

//GetMessageReadList is 获取已读列表
func (c *DingTalkClient) GetMessageReadList(messageID string, cursor int, size int) (data MessageReadListResponse, err error) {
	params := url.Values{}
	params.Add("messageId", messageID)
	params.Add("cursor", strconv.Itoa(cursor))
	params.Add("size", strconv.Itoa(size))
	err = c.httpRPC("chat/getReadList", params, nil, &data)
	return data, err
}