package main

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
	"log"
	"net/http"
)

var client *alipay.Client

const (
	kAppId       = "2021004130630630"
	kPrivateKey  = "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCIFnLZE4qgFeqkiIl/XxlhWwuVY69iNpEeyRpSPylJE90FU8QJmD0dMcHOBe0btF2WM1LNDm0xfiNYo7puGZKdnZH3eg+BobrRF84Yd7EV2KOgBnJyFdG6Z+WCtsCctynrxHP35rvsH0GnWMB50frGaca79rgXBiaSQDXUc/lr9BK7wHq9AzNsM44GfG8ZV7E4vJuT8Q/UALjWRJjPwMUyrWqkkZ/sXVcBHjQRYBMsxs4QGS5f4x8+zH3660j6xQu0bkHXuF6MFBYI8D9V6RImFb+t1TNVTRbsroRDuuuKAYG7Msz5rZMk3G8+WvMBNLD6gE0J3DHVWsRWDybZmNcrAgMBAAECggEATxTcgJD+ibCyvhIp9L6KvSxvOszG6kfTZRRwG0Ng9np4gdP/o6O4P/LeMJ74/pR9nf8NKuQrSAuI9sWOXNS3gqhCXpGr2STmrwMqXMWRPqS50RBuCSXszmn50U/p9ifqUHvE+FY49injuR+2GhSPqiVlu2zP5XFJkMuHfII4eyZmtpeqesNZ49XIIfnHPpB3B+aQeJlsU36sFlP2TxbZlvvf5WJwy4AHWJwee3osWwIdPyb18rV1mh702i3Z3Sh5Ob7gvwnHnoS7slqaiLhWvo9YWAOBVLBSsLxqW4tVtWYzz8/DW9pBmK1FP1fNtDOJeamhCjZt4B2Fg5HGwhZbYQKBgQDwDlfUjkWzs7dh671eDOuYaxvku4ULY+Sbl8bR3NifoSI7KcyxkgqDWMYUvRLAOaPVJXN1XVj+fXv3qA5WZvloBhLqx9mEIT73a2pdl+VH4xY2QyMA5u7AJ27hNK892dDeIXWWpr/n1/TEsIInWdRJ+gYvuC29o0VoMndukAGI1QKBgQCRIFa1K3KVsYzXz/wstsJemhKprQn6S4Wm9RZBG2tNMGNDsNohiM6nSchwetP622DCU1j6Je2VlgVvARQl8J+wP+1nzOAuAogmxLysZrC8ow4Chjr3hLOGKrPoDbCZwmTLI5+PdcOpE20ZHPbZZL1sXxMkeOoeHDpEHus+umDf/wKBgBwfsg6G9IePIIbqVW81WEytD8GGbpndBCVubK6djwt0l0wTI5YSJAUrW1mGpTG8DwOjtZkkbI60KNfk6nkY61NSktjKvSMLuhLGlNmOOCBp7GpDB1DNvV7pv8XGpFk8sYm0pdAWjRkeZeC9RSJTFdns3tisXT+AZ1tDvlZHrMZZAoGAGvsAI75oFxxjKtwn7cgsapoKTjE1YasYtelqsb//OuJ8EeGXLBTbFo7JDOBI+KJAYuBL8nWKrfyuFe0FaehKR+IaqOmV4/fkiBCbYxHUWb2WpTF/VPT+yzq1J7cj1fIl+v4sc+dY8N4Dsl+IJPJtpPAoBufT3rUwv+lfotHToNECgYBIr/y4BcJMhWjhTChw+Yd2zTTyFfF3qZC2if+A6v6cw1N6aS8WpguMo1N7Cev4nP6CwnoXoM4/OihXaE/Ry1wd6y+M1wdy6txoZbpnm6XHQNK/k83GfOZcfVybLj19K/JfdVp5UI3EIekD1z3y+W4/9E7l5JXl8S/L/m80zUEW6g=="
	aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyDUDy+1yUkqZ8ukNRxO/pkxiCdMhftDYi6Ek92iVq6ronKnQAmR7HDmMIw8oOpDdWxa6T05uHcDpH7ZTbFrPGBq4w2ivFlbr63jxjFJ4Xdz2LfShQP0hfYHNQrpiV+z2613JLCjRp3FA3IzAut6n6eirWLBh9xYX8EXELB6V8jqtlxco4v5uR4F3go8rWoE+Dmf2J0CGJEKN/5jP49usuJqYw2Zvz6cE/la5JCocnx/yn/pJUSGe0a8ktvYmO/UNBioXBpcPB7GMsky9/H1JeZRWrp4PNF+a3z+WYzuRu4+tlxmf5+Nx3Hh45a1SSkrZwG/ow+ymPaw7GQ23Y+WPyQIDAQAB"
	kServerPort  = "9989"
	// TODO 设置回调地址域名
	kServerDomain = ""
)

func main() {
	var err error
	print(len("你好Go"))
	if client, err = alipay.New(kAppId, kPrivateKey, false); err != nil {
		log.Println("初始化支付宝失败", err)
		return
	}
	if err = client.LoadAliPayPublicKey(aliPublicKey); err != nil {
		log.Println("加载支付宝公钥失败", err)
		return
	}
	http.HandleFunc("/alipay/pay", pay)
	http.HandleFunc("/alipay/callback", callback)
	http.HandleFunc("/alipay/notify", notify)

	http.ListenAndServe(":"+kServerPort, nil)
}

func pay(writer http.ResponseWriter, request *http.Request) {
	var tradeNo = fmt.Sprintf("%d", xid.Next())

	var p = alipay.TradePreCreate{}
	p.OutTradeNo = tradeNo
	p.Subject = "测试订单"
	p.TotalAmount = "0.01"

	rsp, err := client.TradePreCreate(p)
	if err != nil {
		log.Println(err)
	}

	if rsp.IsFailure() {
		log.Println(rsp.Msg, rsp.SubMsg)
	}
	log.Println(rsp.QRCode)
	writer.Write([]byte(rsp.QRCode))
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
