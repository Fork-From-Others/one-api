package controller

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/pkg/xpem"
	"github.com/go-pay/gopay/pkg/xrsa"
	"github.com/skip2/go-qrcode"
	"net/http"
	"one-api/model"
	"strconv"
	"time"
)

const (
	AppId                = "2021004128689728"
	alipayPublicKey_RSA2 = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyDUDy+1yUkqZ8ukNRxO/pkxiCdMhftDYi6Ek92iVq6ronKnQAmR7HDmMIw8oOpDdWxa6T05uHcDpH7ZTbFrPGBq4w2ivFlbr63jxjFJ4Xdz2LfShQP0hfYHNQrpiV+z2613JLCjRp3FA3IzAut6n6eirWLBh9xYX8EXELB6V8jqtlxco4v5uR4F3go8rWoE+Dmf2J0CGJEKN/5jP49usuJqYw2Zvz6cE/la5JCocnx/yn/pJUSGe0a8ktvYmO/UNBioXBpcPB7GMsky9/H1JeZRWrp4PNF+a3z+WYzuRu4+tlxmf5+Nx3Hh45a1SSkrZwG/ow+ymPaw7GQ23Y+WPyQIDAQAB"
	PrivateKey           = "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCZEdwKkPYQOzfDCSMBksPfjGeOn0kWWwql0nGVa7QkGcCnzUfBD4Z4Zq3Krhu5w3TFSOHVyH/JHQ3R1KxZGgVWsfD3hCkQNyQCYgP1qTAHvrFsQwbIJyTwfJYgqFftR+USpbXuGy026yexX3EyvgJ+XBc/Th0r3K72fFImg4Hxz/VzoU7c2xb9MZm1d+4W1H6yn7vPjK+/Gsvjrvo6GsWGira9Thg+Y7ScbNY2wvJCi4hgSSAYsSUu5YdlhTSvenOp1U8KnmDxCunxV751pEfwidkH9HStUvuxClO/c+9J5xS0w+fbWfvux8SE3gZvObEdMc5NLLwKJQ2Tbo6AfW17AgMBAAECggEAGiSDuUZC0Ejc4DGaSfWAAJkhQqmPuQK5kdKcVZG8hYHkdoRH4gA9zihzPC96SsLIGb213GZO9NFCf/jbqqgYC1N+vTdUMBHK06Fb3cQUkO4PrVbRPLP6yhvtJAy8X6ksxX+Jz+3SThPhhpivY6QyFtSxn70+nDQnqa9X9H33Xo4LAE+EzivbOmZeRC1JXsFzZSbLIvMwibXg+dTf0bWl17rmePDJtxLW7c+20iVFMynJETI7uIBN2gHkp7vLAOFVfsMSK0rYNViNKA4hUhJw1SedbO3BwPdsxOZouYm0MY+w3Y5N1OBEZQ8jjMOVBlTfd0GCodPtYYchehbC0DwvmQKBgQDautxdM2QGPQiQ0kb1QMDbKABn+O86rm+TDKJh6x+dQps/hgLnAcvX57xNVmIQca+/WSAcXffnmV2czTugmRnfdy0Y89g9tUPurFyu09SfRV5Vrj0bQY6jdXFArZkkAsNec2bdWOso8Xnl7wtH00TD75x2ureNi714xz8J2qdA7wKBgQCzJtxLBxvGwXNIdSToAHwFR73yHVN62pdamtrhZGx+2RX/2oVuTn0l8rmff+G7W+rOeMdVIolW9trPaMXvvaDcpwOVkaOXHmSZmDQMT6O5mk2xT0KBoqGnnLS8EWkJpGlv64XnSIsEujXFnKLJugj3a+iYv3qTEmq26rDQCDrENQKBgCkazKbHLZjuh1mP6r3UOWn2Dn17jpmchmNAEJQON5a6GarKaGk5MTGV3xE5lpw4gSqYeSxbjGb9r1X0S6xWmUIhh1wVFyIhmm6T/abtMBvuUVgQsnMY0tFtFKdu+ESIMGbjkQUv3KGJH7tSPPB2h4m60dCOLkhvZl/4MaSMbroJAoGASPnjcoyKvAPBOhq91eOcoWn/7cgUYU75qGa8EmQd7e3wEDCreatvPy4IfvhQs0lV9JUuXXecCliz+RjsyCOuizNdOmgBA2XWBNsDGKC4SLqaO0fWB4h/4Q7scE+HQe4/JOADw5rBRkOz87NCfHnTfTXvoYkeHRq7bZdcPuGbTqECgYB/EvNv9/0piSv5VTkLZHUBsSSnaYzIbI9FjUSHSgMYdUB2Q83JKkvJh549p9p1OzSLxVcvu8G3BllWgvOhEzrl97Nhi9I/8tfDbD2oXlrGJyumcslMtsasSuWjiLQUoLOc8lkTEtIW4MNCCvDsUg7SshAZFq+RDlTgrIWiAZKcQQ=="
	isPro                = true
)

var ctx = context.Background()

func GetPaymentQrcode(c *gin.Context) {
	alipayPreTradeParams := model.AlipayGetPaymentQrcodeReq{}
	err := c.ShouldBindJSON(&alipayPreTradeParams)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	client, err := alipay.NewClient(AppId, PrivateKey, isPro)
	if err != nil {
		xlog.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	//配置公共参数
	client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetNotifyUrl("/alipay/notify")

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("app_id", AppId)
	bm.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	bm.Set("charset", "utf-8")
	bm.Set("version", "1.0")
	bm.Set("notify_url", "")
	bm.Set("method", "alipay.trade.precreate")
	bm.Set("subject", alipayPreTradeParams.Subject)
	bm.Set("out_trade_no", time.Now().Format("20060102150405"))
	bm.Set("total_amount", alipayPreTradeParams.Money)
	bm.Set("biz_content", bm.JsonBody())

	key := xrsa.FormatAlipayPrivateKey(PrivateKey)
	priKey, err := xpem.DecodePrivateKey([]byte(key))
	if err != nil {
		xlog.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	sign, err := alipay.GetRsaSign(bm, alipay.RSA2, priKey)
	if err != nil {
		xlog.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	bm.Set("sign", sign)
	bm.Set("sign_type", alipay.RSA2)

	//创建订单 https://opendocs.alipay.com/open/02ekfg
	aliRsp, err := client.TradePrecreate(ctx, bm)
	if err != nil {
		xlog.Error("err:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	var png []byte
	png, err = qrcode.Encode(aliRsp.Response.QrCode, qrcode.Medium, 300)
	if err != nil {
		xlog.Error("err:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	rsp := &model.AlipayGetPaymentQrcodeRsp{
		ProductId:    alipayPreTradeParams.ProductId,
		OutTradeNo:   aliRsp.Response.OutTradeNo,
		QrCodeUrl:    aliRsp.Response.QrCode,
		QrCodeBase64: "data:image/png;base64," + base64.StdEncoding.EncodeToString(png),
	}
	xlog.Debug("aliRsp:", *aliRsp)
	xlog.Debug("aliRsp.QrCode:", aliRsp.Response.QrCode)
	xlog.Debug("aliRsp.OutTradeNo:", aliRsp.Response.OutTradeNo)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    rsp,
	})
	return
}

func GetPaymentStatus(c *gin.Context) {
	productId, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	redemption, err := model.GetRedemptionById(productId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	outTradeNo := c.Query("out_trade_no")

	client, err := alipay.NewClient(AppId, PrivateKey, true)
	if err != nil {
		xlog.Error(err)
		return
	}
	//配置公共参数
	client.SetCharset("utf-8").
		SetSignType(alipay.RSA2)
	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("app_id", AppId)
	bm.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	bm.Set("charset", "utf-8")
	bm.Set("version", "1.0")
	bm.Set("notify_url", "")
	bm.Set("method", "alipay.trade.query")
	bm.Set("subject", redemption.Name)
	bm.Set("out_trade_no", outTradeNo)
	bm.Set("total_amount", redemption.Price)
	bm.Set("biz_content", bm.JsonBody())
	//查询订单
	aliRsp, err := client.TradeQuery(ctx, bm)
	if err != nil {
		xlog.Error("err:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    aliRsp,
	})
	return
}
