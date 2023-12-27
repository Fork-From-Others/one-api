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
	kAppId      = "2021004128689728"
	kPrivateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCF4G2a6hPp86CfvpptYDpOihroJ3SXnhn1yEJXPAScsfN8zYFbOnfNjyoX6cqNnCJiihmYY2Y3+W/w+IRBNiaBDrCunvyHj9glXZK9R00SFsFuYYGoGHpOVYpc0HoKiehX7wZRXN7K+1GVv/7XK8NvY2U6pfWgoGdcGm3gDbIPF7La2xGxocbqb8mtUKtRbPwGZwHp+xJi9EJK0ZTUJE9Z+8KiUhudiiAnwzNdKLsa3pcvxW7sJeYDryC6A8w2vX/8Jh3yeBsw3RkstWxUJsZ7HOtYfDe7IfrS7LnrJ4gBMD3dmVkxSP32LPwiKPrQ+JaJCHnAUKF8UjsPH3bAOjtNAgMBAAECggEAYUjdXR2Mqw2nQ242uhSbSkeBlgJV73esVbbYvpuWnmeSELclsS2jsXS/mfECiDVVp1XDk8FnnnVcqzdspBa9lDsgmURfLgORhlWhNHqDvwlaNuQUXBqthg8TJK86gD4G4R+I78cU/1vxxWhnv+TFeEQ6Y4wGGlt1wLBT9+T754q7EzdCNlbTDHLtnLnI4APYnuZFZ3VHK4a9Ya/KxeG4Nzui0u6TUBgzMtnkr1o1C7YohcJo7SWX1zW1pqz6AeJQgaKt3QavOPAp+9ig5+Fox79ubxN5OrpcV0tJbTXAxlLBTIvx4p3WiHGfb/RBW+14xIzYB4xoGbCPLQeT2rKJ4QKBgQDJ6GN2GsVt0c9FygAOKSlap3Jw0WtjLoXIH2x45r97Ybh7PjCQ7s499+OsYUXuKaY/fawMW0aTPagnzBcs4nrGf+kiTIUd5o0DZl/hoVoqSQGQJiEe8gkmTTkA0A58u9ygCoJL+371Ld6utL2YPnPG5t6bvcJIqCc3QFFL05/EBwKBgQCpvjSUP6Tds5RginkM0+isqoS8BpJFoHe31ZM6V+Oy2hi8wWXtQPhULEbH/6IZf+57lSRbpQyJRrT7RWEdgfM9vjBBs9udnAXYYxxoVfw2JGM5k3KTOetAXcD22zLxZ13PdFM40svrl6ACgb2DxSPmHQmXB9JWY0e8kF70rz75CwKBgDcspwleNXdWL1L96Vf1TZS1T2OfBr56txpB5A2B2O+pe5VKIFFqwLdUZ0Xy6v4zKXgOKpbR1o5j+fCuJ+MfHLfjulyiFnpeR0iXkDoDMrV9b8zpboGbWH67+YjMUjkpyRd+565F0qLDXyfUwj65SkAfVNOwXkgmk1jY3Z5f/Te/AoGBAJu/F3VzVC3MSXevdtSLV7GeoD359ZqHW3HFONrOq/F+ZjZDaeegtnpdSfDWoQCuvr6MIRkpvu/yfbsUdMBjbTGY4aRXiEn8T+y1O+qMAugWySiaHwXxneaoX5bhl8OiqZPhUz8PQ+Z+cGX8b0yJxZ+tww/eMGPtonPlMAtpBOnzAoGBAJUNrlEiXycYhs91GdQEClwAWlMTseOX0Zc/ArgVJGrShAQ4ZJ7oA1RzQnB3LhsFL5pAYD1tz+SSa3rELAdJeMqR0VQVGW5gSqwI/yvHYg5SjuoVMrJBQXk2aAlTIWGn6eeT282Uz7flXzryDINwO0VOonMqr8HXAOksNuFKoLT/"
	kServerPort = "9989"
	// TODO 设置回调地址域名
	kServerDomain = ""
)

func main() {
	var err error

	if client, err = alipay.New(kAppId, kPrivateKey, true); err != nil {
		log.Println("初始化支付宝失败", err)
		return
	}

	// 加载证书
	if err = client.LoadAppCertPublicKeyFromFile("examples/alipay/appCertPublicKey.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}
	if err = client.LoadAliPayRootCertFromFile("examples/alipay/alipayRootCert.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}
	if err = client.LoadAlipayCertPublicKeyFromFile("examples/alipay/alipayCertPublicKey_RSA2.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}

	if err = client.SetEncryptKey("9F7r76gjRfF7uW4+Jwxt5A=="); err != nil {
		log.Println("加载内容加密密钥发生错误", err)
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
	p.TotalAmount = "10.00"
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
