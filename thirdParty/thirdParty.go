package thirdParty

import (
	"app/conf"
	"app/datastruct"
	"app/log"
	"bytes"
	"crypto/md5"
	crypto_rand "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json" //json封装解析
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type WX_Data struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

func GetWXData(code, appid, secret string) *WX_Data {
	var buf bytes.Buffer
	buf.WriteString("https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appid)
	buf.WriteString("&secret=" + secret)
	buf.WriteString("&code=" + code)
	buf.WriteString("&grant_type=authorization_code")
	url := buf.String()
	p_body := httpGet(url)
	wx_data := new(WX_Data)
	if json_err := json.Unmarshal(*p_body, wx_data); json_err == nil {
		return wx_data
	}
	return nil
}
func GetWXUserData(openId string, accessToken string, platform datastruct.Platform) *datastruct.WXUserData {
	var buf bytes.Buffer
	buf.WriteString("https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken)
	buf.WriteString("&openid=" + openId)
	if platform == datastruct.H5 {
		buf.WriteString("&lang=zh_CN")
	}
	url := buf.String()
	p_body := httpGet(url)
	wx_user := new(datastruct.WXUserData)
	if json_err := json.Unmarshal(*p_body, wx_user); json_err == nil {
		wx_user.HeadImgUrl = wx_user.HeadImgUrl + datastruct.TmpAvatarPostfix
		return wx_user
	}
	return nil
}

func httpGet(url string) *[]byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return &body
}

func WXpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	//log.Debug("sorted_keys------:%s", sorted_keys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

func WXunifyChargeReq(pay *datastruct.WX_PayInfo, ip_addr string) (*ResponsWXPayXML, error) {
	req := new(datastruct.WXOrderReq)
	req.Appid = pay.Appid
	req.Body = pay.OrderForm.Desc
	req.Mch_id = pay.Mch_id
	req.Nonce_str = GetRandomString()
	req.Notify_url = fmt.Sprintf("%s%s", conf.Server.OutHttpServer, datastruct.WXPayCallRoute)
	req.Trade_type = pay.Trade_type
	req.Spbill_create_ip = ip_addr
	req.Total_fee = int(pay.OrderForm.Money * datastruct.MoneyFactor)
	req.Out_trade_no = pay.OrderForm.Id
	req.OpenId = pay.OrderForm.OpenId
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = req.Appid
	m["body"] = req.Body
	m["mch_id"] = req.Mch_id
	m["notify_url"] = req.Notify_url
	m["nonce_str"] = req.Nonce_str
	m["trade_type"] = req.Trade_type
	m["spbill_create_ip"] = req.Spbill_create_ip
	m["total_fee"] = req.Total_fee
	m["out_trade_no"] = req.Out_trade_no
	m["openid"] = req.OpenId

	req.Sign = WXpayCalcSign(m, pay.PaySecret) // 这个是计算wxpay签名的函数上面已贴出
	bytes_req, err := xml.Marshal(req)
	if err != nil {
		log.Error("WXunifyChargeReq(): xml.Marshal error:%s", err)
		return nil, err
	}
	str_req := string(bytes_req)

	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "WXOrderReq", "xml", -1)

	//log.Debug("old-----------:%s", str_req)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	request, err := http.NewRequest("POST", datastruct.WXPayUrl, bytes.NewReader(bytes_req))
	if err != nil {
		log.Error("wxUnifyChargeReq(): http.NewRequest error:%s", err)
		return nil, err
	}
	request.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	request.Header.Set("Content-Type", "application/xml;charset=utf-8")
	c := http.Client{}
	resp, _err := c.Do(request)
	if _err != nil {
		log.Error("wxUnifyChargeReq(): http.Do error:%s", _err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	retXml := string(body)
	//log.Debug("WXunifyChargeReq retXml-----------:%s", retXml)
	if err != nil {
		return nil, err
	}
	resp_wx_pay := new(ResponsWXPayXML)
	err = xml.Unmarshal([]byte(retXml), resp_wx_pay)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_wx_pay, nil
}

type ResponsWXPayXML struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	AppId       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Result_code string `xml:"result_code"`
	Trade_type  string `xml:"trade_type"`
	NonceStr    string `xml:"nonce_str"`
	PaySign     string `xml:"sign"`
	Prepay_id   string `xml:"prepay_id"`
}

type ResponsWXPayeeXML struct {
	Return_code      string `xml:"return_code"`
	Result_code      string `xml:"result_code"`
	Partner_trade_no string `xml:"partner_trade_no"`
	Payment_no       string `xml:"payment_no"`
	Payment_time     string `xml:"payment_time"`
	Err_code         string `xml:"err_code"`
	Err_code_des     string `xml:"err_code_des"`
}

func WXunifyPayeeReq(amount float64, openid string, ip_addr string, trade_no string, isOnlyApp bool) (*ResponsWXPayeeXML, error) {
	req := new(datastruct.WXPayeeReq)
	req.Amount = int(amount * float64(datastruct.MoneyFactor))
	req.OpenId = openid
	req.Check_name = "NO_CHECK"
	req.Desc = "提现"
	Mch_appid := ""
	Mch_id := ""
	PaySecret := ""
	if isOnlyApp {
		Mch_appid = datastruct.WX_KFPT_AppID
		Mch_id = datastruct.WX_KFPT_Mch_Id
		PaySecret = datastruct.WX_KFPT_PaySecret
	} else {
		Mch_appid = datastruct.WX_GZH_AppID
		Mch_id = datastruct.WX_GZH_Mch_Id
		PaySecret = datastruct.WX_GZH_PaySecret
	}
	req.Mch_appid = Mch_appid
	req.Mch_id = Mch_id
	req.Nonce_str = GetRandomString()
	req.Partner_trade_no = trade_no
	req.Spbill_create_ip = ip_addr

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["amount"] = req.Amount
	m["check_name"] = req.Check_name
	m["desc"] = req.Desc
	m["mch_appid"] = req.Mch_appid
	m["mchid"] = req.Mch_id
	m["nonce_str"] = req.Nonce_str
	m["openid"] = req.OpenId
	m["partner_trade_no"] = req.Partner_trade_no
	m["spbill_create_ip"] = req.Spbill_create_ip

	req.Sign = WXpayCalcSign(m, PaySecret) // 这个是计算wxpay签名的函数上面已贴出
	bytes_req, err := xml.Marshal(req)
	if err != nil {
		log.Error("WXunifyChargeReq(): xml.Marshal error:%s", err)
		return nil, err
	}
	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "WXOrderReq", "xml", -1)

	bytes_req = []byte(str_req)

	resp, _err := securePost(datastruct.WXPayeeUrl, bytes_req, isOnlyApp)
	if _err != nil {
		log.Error("wxUnifyChargeReq(): http.Do error:%s", _err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	retXml := string(body)
	if err != nil {
		return nil, err
	}
	//log.Debug("retXml-------------------%s", retXml)
	resp_wx_payee := new(ResponsWXPayeeXML)
	err = xml.Unmarshal([]byte(retXml), resp_wx_payee)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_wx_payee, nil
}

//生成32位md5字串
func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成32随机字符串
func GetRandomString() string {
	// 生成节点实例
	b := make([]byte, 48)
	if _, err := io.ReadFull(crypto_rand.Reader, b); err != nil {
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}

//ca证书的位置，需要绝对路径
var _tlsConfig *tls.Config

//采用单例模式初始化ca
func getTLSConfig(isOnlyApp bool) (*tls.Config, error) {
	if _tlsConfig != nil {
		return _tlsConfig, nil
	}
	var certFile, keyFile string
	if isOnlyApp {
		certFile = datastruct.WX_KFPT_CertPath
		keyFile = datastruct.WX_KFPT_KeyPath
	} else {
		certFile = datastruct.WX_GZH_CertPath
		keyFile = datastruct.WX_GZH_KeyPath
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	_tlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            pool,
	}
	return _tlsConfig, nil
}

//携带ca证书的安全请求
func securePost(url string, xmlContent []byte, isOnlyApp bool) (*http.Response, error) {
	tlsConfig, err := getTLSConfig(isOnlyApp)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}
	return client.Post(
		url,
		"application/xml",
		bytes.NewBuffer(xmlContent))
}
