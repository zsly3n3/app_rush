package tools

import (
	"app/conf"
	"app/datastruct"
	"app/log"
	"app/thirdParty"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/holdno/snowFlakeByGo"
)

var IdWorker *snowFlakeByGo.Worker

func CreateIdWorker() {
	IdWorker, _ = snowFlakeByGo.NewWorker(0)
}
func UniqueId() string {
	return Int64ToString(IdWorker.GetId())
}

func Int64ToString(tmp int64) string {
	return strconv.FormatInt(tmp, 10)
}

func StringToInt64(tmp string) int64 {
	rs, _ := strconv.ParseInt(tmp, 10, 64)
	return rs
}

func StringToFloat64(tmp string) float64 {
	rs, _ := strconv.ParseFloat(tmp, 64)
	return rs
}

func BoolToInt(tf bool) int {
	rs := 0
	if tf {
		rs = 1
	}
	return rs
}

func IntToBool(tf int) bool {
	rs := false
	if tf == 1 {
		rs = true
	}
	return rs
}

func IntToString(tmp int) string {
	return strconv.Itoa(tmp)
}

func StringToInt(tmp string) int {
	rs, _ := strconv.Atoi(tmp)
	return rs
}

func BoolToString(tmp bool) string {
	if tmp == false {
		return "0"
	} else {
		return "1"
	}
}

func StringToBool(tmp string) bool {
	if tmp == "0" {
		return false
	} else {
		return true
	}
}

func InterfaceToString(tmp interface{}) (string, bool) {
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
		log.Debug("PlayerSoilToString error:%s", err.Error())
		return "", true
	}
	return string(jsons), false
}

func BytesToOrderForm(bytes []byte) (*datastruct.OrderForm, bool) {
	tmp := new(datastruct.OrderForm)
	err := json.Unmarshal(bytes, tmp)
	if err != nil {
		log.Debug("BytesToOrderForm error:%s", err.Error())
		return nil, true
	}
	return tmp, false
}

func BytesToGameState(bytes []byte) (*datastruct.GameState, bool) {
	tmp := new(datastruct.GameState)
	err := json.Unmarshal(bytes, tmp)
	if err != nil {
		log.Debug("BytesToGameState error:%s", err.Error())
		return nil, true
	}
	return tmp, false
}

func GetWXUserData(openid string, access_token string, platform datastruct.Platform) (*datastruct.WXUserData, bool) {
	wx_user := thirdParty.GetWXUserData(openid, access_token, platform)
	if wx_user == nil || wx_user.OpenId == "" || wx_user.UnionId == "" {
		return nil, true
	}
	return wx_user, false
}

func GetOpenidAndAccessToken(code, appid, secret string) (string, string, bool) {
	wx_data := thirdParty.GetWXData(code, appid, secret)
	if wx_data == nil || wx_data.AccessToken == "" || wx_data.OpenId == "" {
		return "", "", true
	}
	return wx_data.OpenId, wx_data.AccessToken, false
}

func CreateAuthLink(token string, appid string, addr string) string {
	var buf bytes.Buffer
	buf.WriteString("https://open.weixin.qq.com/connect/oauth2/authorize?")
	buf.WriteString("appid=" + appid)
	buf.WriteString("&redirect_uri=" + url.QueryEscape(addr))
	buf.WriteString("&response_type=code")
	buf.WriteString("&scope=snsapi_userinfo")
	var value string
	if conf.Common.Mode == conf.Debug {
		value = token + ",app_dev"
	} else {
		value = token
	}
	parms := "&state=" + value + "#wechat_redirect"
	buf.WriteString(parms)
	return buf.String()
}

func CreateInviteLink(token string, url string) string {
	var buf bytes.Buffer
	buf.WriteString(url)
	buf.WriteString("?referrer=" + token)
	if conf.Common.Mode == conf.Debug {
		buf.WriteString("&app=app_dev")
	}
	return buf.String()
}

func CreateDownLoadLink(token string, url string) string {
	var buf bytes.Buffer
	buf.WriteString(url)
	buf.WriteString("?uid=" + token)
	return buf.String()
}

func GetShareImgUrl(imgName string) string {
	var path string
	if conf.Common.Mode == conf.Debug {
		path = "/devshare/"
	} else {
		path = "/share/"
	}
	url := conf.Server.Domain + path + imgName + ".png"
	return url
}

func GetKfQRcode(imgName string) string {
	var path string
	if conf.Common.Mode == conf.Debug {
		path = "/devkfqrcode/"
	} else {
		path = "/kfqrcode/"
	}
	url := conf.Server.Domain + path + imgName + ".png"
	return url
}

func GetShareText() (string, string) {
	title := "999口红机"
	desc := "999口红机分享描述"
	return title, desc
}

func CreateGoodsImgUrl(imgName string) string {
	url := conf.Server.Domain + getGoodsImgPath() + imgName
	return url
}

func getGoodsImgPath() string {
	var path string
	if conf.Common.Mode == conf.Debug {
		path = "/devgoodsimg/"
	} else {
		path = "/goodsimg/"
	}
	return path
}

func GetImgPath() string {
	url := "assets" + getGoodsImgPath()
	return url
}

//店长推荐
func CreateDZTJImgUrl(imgName string) string {
	url := conf.Server.Domain + "/dztj/" + imgName + ".png"
	return url
}

func CreateUserDescInfo(nickname string, goodsName string) string {
	str := nickname + "闯关成功,已获得" + goodsName
	return str
}

//保留两位小数
func Decimal2(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func CreateImgFromBase64(base64_str *string, path string) bool {
	isError := false
	dist, err := base64.StdEncoding.DecodeString(*base64_str)

	if err != nil {
		log.Debug("CreateImgFromBase64  base64.StdEncoding.DecodeString err:%s", err)
		return true
	}
	var f *os.File
	f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	if err != nil {
		log.Debug("CreateImgFromBase64  os.OpenFile err:%s", err)
		return true
	}
	_, err = f.Write(dist)
	if err != nil {
		log.Debug("CreateImgFromBase64  f.Write err:%s", err)
		return true
	}
	return isError
}

func GetTodayTomorrowTime() (int64, int64) {
	now_Time := time.Now()
	tomorrow := now_Time.Add(24 * time.Hour)
	year, month, day := tomorrow.Date()
	tomorrow_Time := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	today_Time := time.Date(now_Time.Year(), now_Time.Month(), now_Time.Day(), 0, 0, 0, 0, time.Local)

	tomorrow_unix := tomorrow_Time.Unix()
	today_unix := today_Time.Unix()
	return today_unix, tomorrow_unix
}

func GetYesterdayTodayTime() (int64, int64) {
	now_Time := time.Now()
	yesterday := now_Time.AddDate(0, 0, -1)
	year, month, day := yesterday.Date()
	yesterday_Time := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	today_Time := time.Date(now_Time.Year(), now_Time.Month(), now_Time.Day(), 0, 0, 0, 0, time.Local)

	yesterday_unix := yesterday_Time.Unix()
	today_unix := today_Time.Unix()
	return yesterday_unix, today_unix
}

func DeleteFile(filePath string) bool {
	tf := true
	err := os.Remove(filePath)
	if err != nil {
		tf = false
	}
	return tf
}

// func SaveActivityData(activity []*datastruct.TodayUserActivityInfo, oldusers int, newusers int, fileName string) datastruct.CodeType {
// 	xlsx := excelize.NewFile()
// 	count := len(activity)

// 	SheetName := "Sheet1"
// 	index := xlsx.NewSheet(SheetName)
// 	xlsx.SetCellValue(SheetName, "A1", "UserId")              //用户id
// 	xlsx.SetCellValue(SheetName, "B1", "StartTime")           //今天开始访问时间
// 	xlsx.SetCellValue(SheetName, "C1", "EndTime")             //今天结束访问时间
// 	xlsx.SetCellValue(SheetName, "D1", IntToString(oldusers)) //老用户数
// 	xlsx.SetCellValue(SheetName, "E1", IntToString(newusers)) //新用户数
// 	start_index := 2
// 	for i := 0; i < count; i++ {
// 		row_name_A := fmt.Sprintf("A%d", start_index)
// 		row_name_B := fmt.Sprintf("B%d", start_index)
// 		row_name_C := fmt.Sprintf("C%d", start_index)
// 		value_A := IntToString(activity[i].UserId)
// 		value_B := Int64ToString(activity[i].StartTime)
// 		value_C := Int64ToString(activity[i].EndTime)
// 		xlsx.SetCellValue(SheetName, row_name_A, value_A)
// 		xlsx.SetCellValue(SheetName, row_name_B, value_B)
// 		xlsx.SetCellValue(SheetName, row_name_C, value_C)
// 		start_index++
// 	}

// 	// Set active sheet of the workbook.
// 	xlsx.SetActiveSheet(index)
// 	// Save xlsx file by the given path.
// 	var dir string
// 	if conf.Common.Mode == conf.Debug {
// 		dir = "/devdata/"
// 	} else {
// 		dir = "/data/"
// 	}
// 	path := fmt.Sprintf("assets%sactivitydata/%s.xlsx", dir, fileName)
// 	err := xlsx.SaveAs(path)
// 	if err != nil {
// 		log.Debug("SaveActivityData err:%s", err.Error())
// 		return datastruct.UpdateDataFailed
// 	}
// 	return datastruct.NULLError
// }

// func GetActivityData(fileName string) (int, int) {
// 	var dir string
// 	if conf.Common.Mode == conf.Debug {
// 		dir = "/devdata/"
// 	} else {
// 		dir = "/data/"
// 	}
// 	path := fmt.Sprintf("assets%sactivitydata/%s.xlsx", dir, fileName)
// 	xlsx, err := excelize.OpenFile(path)
// 	if err != nil {
// 		return 0, 0
// 	}
// 	SheetName := "Sheet1"
// 	cell_old := xlsx.GetCellValue(SheetName, "D1")
// 	cell_new := xlsx.GetCellValue(SheetName, "E1")
// 	return StringToInt(cell_old), StringToInt(cell_new)
// }

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func Remove(slice []int, i int) []int {
	return append(slice[:i], slice[i+1:]...)
}

func UnixToString(unix int64, timeLayout string) string {
	dataTimeStr := time.Unix(unix, 0).Format(timeLayout)
	return dataTimeStr
}
