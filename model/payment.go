package model

type AlipayGetPaymentQrcodeReq struct {
	ProductId int     `json:"product_id"`
	Subject   string  `json:"subject"`
	Money     float32 `json:"money"`
}

type AlipayGetPaymentQrcodeRsp struct {
	ProductId    int    `json:"product_id"`
	OutTradeNo   string `json:"out_trade_no"`
	QrCodeUrl    string `json:"qr_code_url"`
	QrCodeBase64 string `json:"qr_code_base64"`
}
