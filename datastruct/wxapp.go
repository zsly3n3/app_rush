package datastruct

const WX_GZH_AppID = "wxe8d884b6dda2b66c"
const WX_GZH_AppSecret = "467ef3ffa7cdef9421f5d7d496638815"
const WX_GZH_Mch_Id = "1519982251" //公众号中微信商户号
const WX_GZH_PaySecret = "ASDIOVBLJMKJ213dNBAS23sas123saqw"

const WX_KFPT_AppID = "wx95b14104c5ae55f4"
const WX_KFPT_AppSecret = "59f84c2c8b1dfd31e8412613abff2d0a"
const WX_KFPT_Mch_Id = "1520240821" //开放平台中微信商户号
const WX_KFPT_PaySecret = "LJMKJ213dNBAS23sas123saqwASDIOVB"

const WX_Payee_Mch_Id = "1520240821" //开放平台中微信商户号

const WX_KFPT_CertPath = "wechatpayserver/cert/apiclient_cert.pem"
const WX_KFPT_KeyPath = "wechatpayserver/cert/apiclient_key.pem"

const WX_GZH_CertPath = "wechatpayserver/cert_gzh/apiclient_cert.pem"
const WX_GZH_KeyPath = "wechatpayserver/cert_gzh/apiclient_key.pem"

//微信支付要传入的参数
type WXOrderReq struct {
	Appid            string `xml:"appid"`            //公众账号ID或应用ID
	Body             string `xml:"body"`             //商品描述
	Mch_id           string `xml:"mch_id"`           //商户号
	Nonce_str        string `xml:"nonce_str"`        //随机字符串
	Notify_url       string `xml:"notify_url"`       //通知地址
	Trade_type       string `xml:"trade_type"`       //交易类型
	Spbill_create_ip string `xml:"spbill_create_ip"` //终端IP
	Total_fee        int    `xml:"total_fee"`        //总金额
	Out_trade_no     string `xml:"out_trade_no"`     //商户订单号
	OpenId           string `xml:"openid"`           //用户openid
	Sign             string `xml:"sign"`             //签名
}

//微信支付结果通知传入的参数
type WXOrderResult struct {
	Appid          string `xml:"appid"`          //公众账号ID或应用ID
	Bank_type      string `xml:"bank_type"`      //付款银行
	Cash_fee       int    `xml:"cash_fee"`       //现金支付金额
	Is_subscribe   string `xml:"is_subscribe"`   //是否关注公众账号
	Mch_id         string `xml:"mch_id"`         //商户号
	Nonce_str      string `xml:"nonce_str"`      //随机字符串
	Trade_type     string `xml:"trade_type"`     //交易类型
	Total_fee      int    `xml:"total_fee"`      //订单金额
	Transaction_id string `xml:"transaction_id"` //微信支付订单号
	Time_end       string `xml:"time_end"`       //支付完成时间
	Out_trade_no   string `xml:"out_trade_no"`   //商户订单号
	Openid         string `xml:"openid"`         //用户标识
	Result_code    string `xml:"result_code"`    //业务结果 SUCCESS/FAIL
	Sign           string `xml:"sign"`           //签名
}

type WXPayeeReq struct {
	Amount           int    `xml:"amount"`           //金额
	Check_name       string `xml:"check_name"`       //校验用户姓名选项
	Desc             string `xml:"desc"`             //企业付款备注
	Mch_appid        string `xml:"mch_appid"`        //公众账号ID或应用ID
	Mch_id           string `xml:"mchid"`            //商户号
	Nonce_str        string `xml:"nonce_str"`        //随机字符串
	OpenId           string `xml:"openid"`           //用户openid
	Partner_trade_no string `xml:"partner_trade_no"` //商户订单号
	Spbill_create_ip string `xml:"spbill_create_ip"` //终端IP
	Sign             string `xml:"sign"`             //签名
}
