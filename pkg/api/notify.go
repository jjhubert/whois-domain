package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"whois-domain/pkg/entity"
)

type notification struct {
	Config       *entity.Config
	Url          string
	AlertTitle   string
	UnknownTitle string
}

const (
	METHODDINGDING = "dingding"
	METHODQIWEI    = "qiwei"
)

var objectNotification *notification

func GetNotificationObject() *notification {
	if objectNotification == nil {
		objectNotification = &notification{}
	}
	return objectNotification
}

func (n *notification) getWebhookInfo() {
	if n.Config == nil {
		return
	}
	switch n.Config.Global.NotifyMethod {
	case METHODDINGDING:
		n.Url = n.Config.DingDing.AlertNotifyUrl
		n.AlertTitle = n.Config.DingDing.AlertTitle
		n.UnknownTitle = n.Config.DingDing.UnknownTitle
	case METHODQIWEI:
		n.Url = n.Config.QiWei.AlertNotifyUrl
		n.AlertTitle = n.Config.QiWei.AlertTitle
		n.UnknownTitle = n.Config.QiWei.UnknownTitle
	}
}

func (n *notification) constructContent(i int, v entity.WhoisResult, m string) (msg string) {
	msg = m
	switch n.Config.Global.NotifyMethod {
	case METHODDINGDING:
		func() {
			if i == 0 && v.ExpireDate == "" {
				msg = n.constructMsg(n.Config.DingDing.UnknownTemplate, v.Domain)
				return
			}
			if i == 0 {
				msg = n.constructMsg(n.Config.DingDing.AlertTemplate, v.Domain, v.ExpireDate)
				return
			}
			msg += "\n\n"
			if v.ExpireDate == "" {
				msg += n.constructMsg(n.Config.DingDing.UnknownTemplate, v.Domain)
			}
			msg += n.constructMsg(n.Config.DingDing.AlertTemplate, v.Domain, v.ExpireDate)
		}()
	case METHODQIWEI:
		func() {
			if i == 0 && v.ExpireDate == "" {
				msg = n.constructMsg(n.Config.QiWei.UnknownTemplate, v.Domain)
				return
			}
			if i == 0 {
				msg = n.constructMsg(n.Config.QiWei.AlertTemplate, v.Domain, v.ExpireDate)
				return
			}
			if v.ExpireDate == "" {
				msg += n.constructMsg(n.Config.QiWei.UnknownTemplate, v.Domain)
			}
			msg += n.constructMsg(n.Config.QiWei.AlertTemplate, v.Domain, v.ExpireDate)
		}()
	}
	return msg
}

func (n *notification) SendExpireMsg(dl []entity.WhoisResult) {
	msg := ""
	n.getWebhookInfo()
	for i, v := range dl {
		msg = n.constructContent(i, v, msg)
	}
	jsonBody, err := n.constructRequestBody(n.AlertTitle, msg)
	if err != nil {
		log.Println(err)
		return
	}
	err = n.postRequest(jsonBody, n.Url)
	if err != nil {
		log.Println(err)
		return
	}
}

func (n *notification) SendUnknownMsg(dl []entity.WhoisResult) {
	msg := ""
	n.getWebhookInfo()
	for i, v := range dl {
		msg = n.constructContent(i, v, msg)
	}

	jsonBody, err := n.constructRequestBody(n.UnknownTitle, msg)
	if err != nil {
		log.Println(err)
		return
	}
	err = n.postRequest(jsonBody, n.Url)
	if err != nil {
		log.Println(err)
		return
	}
}

func (n *notification) postRequest(jsonStr []byte, url string) error {
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (n *notification) constructMsg(template string, args ...interface{}) string {
	tmpStr := fmt.Sprintf(template, args...)
	return tmpStr
}

func (n *notification) constructRequestBody(title string, msg string) ([]byte, error) {
	requestBody := make(map[string]interface{})
	mdContent := make(map[string]string)
	switch n.Config.Global.NotifyMethod {
	case METHODDINGDING:
		mdContent["title"] = title
		mdContent["text"] = msg
		requestBody["msgtype"] = n.Config.DingDing.MsgType
		requestBody["markdown"] = mdContent
	case METHODQIWEI:
		mdContent["content"] = fmt.Sprintf("%s\n%s", title, msg)
		requestBody["msgtype"] = n.Config.DingDing.MsgType
		requestBody["markdown"] = mdContent
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	return jsonBody, nil
}
