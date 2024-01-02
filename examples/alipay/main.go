package main

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
	"log"
	"net/http"
)

const (
	kAppId      = "2021004128689728"
	kPrivateKey = "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCZEdwKkPYQOzfDCSMBksPfjGeOn0kWWwql0nGVa7QkGcCnzUfBD4Z4Zq3Krhu5w3TFSOHVyH/JHQ3R1KxZGgVWsfD3hCkQNyQCYgP1qTAHvrFsQwbIJyTwfJYgqFftR+USpbXuGy026yexX3EyvgJ+XBc/Th0r3K72fFImg4Hxz/VzoU7c2xb9MZm1d+4W1H6yn7vPjK+/Gsvjrvo6GsWGira9Thg+Y7ScbNY2wvJCi4hgSSAYsSUu5YdlhTSvenOp1U8KnmDxCunxV751pEfwidkH9HStUvuxClO/c+9J5xS0w+fbWfvux8SE3gZvObEdMc5NLLwKJQ2Tbo6AfW17AgMBAAECggEAGiSDuUZC0Ejc4DGaSfWAAJkhQqmPuQK5kdKcVZG8hYHkdoRH4gA9zihzPC96SsLIGb213GZO9NFCf/jbqqgYC1N+vTdUMBHK06Fb3cQUkO4PrVbRPLP6yhvtJAy8X6ksxX+Jz+3SThPhhpivY6QyFtSxn70+nDQnqa9X9H33Xo4LAE+EzivbOmZeRC1JXsFzZSbLIvMwibXg+dTf0bWl17rmePDJtxLW7c+20iVFMynJETI7uIBN2gHkp7vLAOFVfsMSK0rYNViNKA4hUhJw1SedbO3BwPdsxOZouYm0MY+w3Y5N1OBEZQ8jjMOVBlTfd0GCodPtYYchehbC0DwvmQKBgQDautxdM2QGPQiQ0kb1QMDbKABn+O86rm+TDKJh6x+dQps/hgLnAcvX57xNVmIQca+/WSAcXffnmV2czTugmRnfdy0Y89g9tUPurFyu09SfRV5Vrj0bQY6jdXFArZkkAsNec2bdWOso8Xnl7wtH00TD75x2ureNi714xz8J2qdA7wKBgQCzJtxLBxvGwXNIdSToAHwFR73yHVN62pdamtrhZGx+2RX/2oVuTn0l8rmff+G7W+rOeMdVIolW9trPaMXvvaDcpwOVkaOXHmSZmDQMT6O5mk2xT0KBoqGnnLS8EWkJpGlv64XnSIsEujXFnKLJugj3a+iYv3qTEmq26rDQCDrENQKBgCkazKbHLZjuh1mP6r3UOWn2Dn17jpmchmNAEJQON5a6GarKaGk5MTGV3xE5lpw4gSqYeSxbjGb9r1X0S6xWmUIhh1wVFyIhmm6T/abtMBvuUVgQsnMY0tFtFKdu+ESIMGbjkQUv3KGJH7tSPPB2h4m60dCOLkhvZl/4MaSMbroJAoGASPnjcoyKvAPBOhq91eOcoWn/7cgUYU75qGa8EmQd7e3wEDCreatvPy4IfvhQs0lV9JUuXXecCliz+RjsyCOuizNdOmgBA2XWBNsDGKC4SLqaO0fWB4h/4Q7scE+HQe4/JOADw5rBRkOz87NCfHnTfTXvoYkeHRq7bZdcPuGbTqECgYB/EvNv9/0piSv5VTkLZHUBsSSnaYzIbI9FjUSHSgMYdUB2Q83JKkvJh549p9p1OzSLxVcvu8G3BllWgvOhEzrl97Nhi9I/8tfDbD2oXlrGJyumcslMtsasSuWjiLQUoLOc8lkTEtIW4MNCCvDsUg7SshAZFq+RDlTgrIWiAZKcQQ=="
	kServerPort = "9989"
	// TODO 设置回调地址域名
	kServerDomain = ""
)

var client *alipay.Client

func main() {
	var err error
	if client, err = alipay.New(kAppId, kPrivateKey, true); err != nil {
		log.Println("初始化支付宝失败", err)
		return
	}

	http.HandleFunc("/alipay/pay", pay)
	http.HandleFunc("/alipay/callback", callback)
	http.HandleFunc("/alipay/notify", notify)

	http.ListenAndServe(":"+kServerPort, nil)
}

func pay(writer http.ResponseWriter, request *http.Request) {
	var tradeNo = fmt.Sprintf("%d", xid.Next())

	var p = alipay.TradePagePay{}
	p.NotifyURL = kServerDomain + "/alipay/notify"
	p.ReturnURL = kServerDomain + "/alipay/callback"
	p.Subject = "支付测试:" + tradeNo
	p.OutTradeNo = tradeNo
	p.TotalAmount = "0.01"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, _ := client.TradePagePay(p)
	http.Redirect(writer, request, url.String(), http.StatusTemporaryRedirect)
}

func callback(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	if err := client.VerifySign(request.Form); err != nil {
		log.Println("回调验证签名发生错误", err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("回调验证签名发生错误"))
		return
	}

	log.Println("回调验证签名通过")

	// 示例一：使用已有接口进行查询
	var outTradeNo = request.Form.Get("out_trade_no")
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo

	rsp, err := client.TradeQuery(p)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())))
		return
	}

	if rsp.IsFailure() {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Msg, rsp.SubMsg)))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("订单 %s 支付成功", outTradeNo)))
}

func notify(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	var notification, err = client.DecodeNotification(request.Form)
	if err != nil {
		log.Println("解析异步通知发生错误", err)
		return
	}

	log.Println("解析异步通知成功:", notification.NotifyId)

	// 示例一：使用自定义请求进行查询
	var p = alipay.NewPayload("alipay.trade.query")
	p.AddBizField("out_trade_no", notification.OutTradeNo)

	var rsp *alipay.TradeQueryRsp
	if err = client.Request(p, &rsp); err != nil {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s \n", notification.OutTradeNo, err.Error())
		return
	}
	if rsp.IsFailure() {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s-%s \n", notification.OutTradeNo, rsp.Msg, rsp.SubMsg)
		return
	}

	log.Printf("订单 %s 支付成功 \n", notification.OutTradeNo)

	client.ACKNotification(writer)
}
