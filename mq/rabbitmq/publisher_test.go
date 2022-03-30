/*
 * @Author: ZhaoMingJun
 * @Date: 2021/7/22 11:50 上午
 */
package rabbitmq

import (
	"encoding/json"
	"testing"
)

type KeyWordInfo struct {
	CategoryId int `json:"category_id"`
	//CategoryType int `json:"category_type"`
	KeyWord  string  `json:"key_word"`
	//KeyWordType int `json:"key_word_type"`
}
type MQMessage struct {
	SourceId   string `json:"source_id"`
	GroupId string `json:"source_group_id"`
	KeyWordList  []KeyWordInfo `json:"key_word_list"`
	Attribute int `json:"attribute"`
	GroupName string `json:"group_name"`
	GroupUserCount string `json:"group_user_count"`
	MerchtypeId string `json:"merchtype_id"`
	MerchtypeName string `json:"merchtype_name"`
	OrderId string `json:"order_id"`
	SuborderId string `json:"suborder_id"`
	MainSiteId string `json:"main_site_id"`
	MainSiteName string `json:"main_site_name"`
	MerchandiseId string `json:"merchandise_id"`
	MerchandiseName string `json:"merchandise_name"`
	ResponsibilityId string `json:"responsibility_id"`
	HighRisk  int `json:"high_risk"`
	Quantity  int `json:"quantity"`  // 数量
	Amount   int  `json:"amount"`  // 金额
	WorkflowSn string `json:"workflow_sn"`
	MsgContent string `json:"msg_content"`


}
func TestSend(t *testing.T) {
	RabbitMqUserName := "admin"
	RabbitMqPassWord := "admin"
	RabbitMqHost := "127.0.0.1:5672"
	RabbitMqExchange := "sht.sco.alarmMessage"
	RabbitMqExchangeType := "topic"
	RabbitMqRoutingKey := "scoAlarmMessage"
	// 构造链接地址
	addr := "amqp://" + RabbitMqUserName + ":" + RabbitMqPassWord + "@" + RabbitMqHost
	publisher := NewPublisher(addr, RabbitMqExchange, RabbitMqExchangeType, RabbitMqRoutingKey)

	keywordList := []KeyWordInfo{
		{
			CategoryId: 1072,
			KeyWord: "退货",
		},
		{
			CategoryId: 1097,
			KeyWord: "退货",
		},
	}

	msg := MQMessage{
		SourceId: "4",
		GroupId: "",
		KeyWordList: keywordList,
		Attribute: 1,
		GroupName: "",
		GroupUserCount: "",
		MerchtypeId: "2409246",
		MerchtypeName: " 500g±50g",
		OrderId:"524679151274688463",
		SuborderId:"524679151278882768",
		MainSiteId:"240",
		MainSiteName:"郑州市",
		MerchandiseId:"2033338",
		MerchandiseName:"乐观鲜 云南鹰嘴芒",
		ResponsibilityId:"0",
		HighRisk:1,
		Quantity:1,
		Amount:199,
		WorkflowSn:"",
		MsgContent:"",
	}
	jsonmsg, _ := json.Marshal(msg)
	publisher.Send(string(jsonmsg))
	// 在项目中使用,Globalpublisher 作为全局变量
	//Globalpublisher := GetConfig()
	//jsonmsg, _ := json.Marshal(msg)
	//Globalpublisher.Send(string(jsonmsg))
}


