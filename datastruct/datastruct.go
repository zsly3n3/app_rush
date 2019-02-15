package datastruct

import "time"

const AdminLevelID = 1
const NormalLevelID = 2
const MaxLevel = 3
const MinInterval = 20 //获奖最小间隔20秒
const AgencyIdentifier = "common"
const DefaultMemberIdentifier = "普通会员"
const TmpDataTimes = time.Hour * 4

const MaxTimeOutForLotteryQueue = 8 //最大排队等待时间
const DefaultId = 1

const EXPIRETIME = 300    //Redis过期时间300秒
const WXOrderMaxSec = 600 //微信订单5分钟超时，这里设置600秒过期
const GameStateAdd = 10   //暂时没用到
const UserIdStart = 39834039

const ProductDesc = "闯关购-金币充值"
const WXPayCallRoute = "/wxpaycall"
const WXPayUrl = "https://api.mch.weixin.qq.com/pay/unifiedorder"
const WXPayeeUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"

const MoneyFactor = 100

// const AddGoldCountForDownLoadApp = 4 //首次app登录送6个金币
// const AddGoldCountForAppraised = 2   //评价送2个金币

type HttpStatusCode int //错误码
const (
	Maintenance CodeType = 900 //服务器维护中
)

type PurchaseType string //购买类型
const (
	GoodsType PurchaseType = "goods" //商品类型
	GoldType               = "gold"  //金币类型
	VipType                = "vip"   //会员类型
)

type CodeType int //错误码
const (
	NULLError                   CodeType = iota //无错误
	ParamError                                  //参数错误,数据为空或者类型不对等
	LoginFailed                                 //登录失败,如无此账号或者密码错误等
	JsonParseFailedFromPostBody                 //来自post请求中的Body解析json失败
	GetDataFailed                               //获取数据失败
	UpdateDataFailed                            //修改数据失败
	VersionError                                //客户端与服务器版本不一致
	TokenError                                  //没有Token或者值为空,或者不存在此Token
	JsonParseFailedFromPutBody                  //来自put请求中的Body解析json失败
	WXCodeInvalid                               //无效的微信code
	PlatformInvalid                             //无效的平台参数
	GoldLess                                    //金币不足
	FastPayModeSucceed                          //中奖过快,此次不记录
	DepositFailed                               //充值失败
	NothasgameRecord                            //没有此存档记录
	GetIpAdrrError                              //获取客户端ip失败
	NotEnoughBalance                            //没有足够的佣金提现
	NotEnoughMinCharge                          //没有达到最低提现额
	OverMaxDrawCount                            //超过今日提现次数
	PayeeOnlyInApp                              //只能在APP内提现
	WeChatPayeeError                            //微信提现出错
	PurchaseFailed                              //购买失败
	AppRedirect                                 //app重定向
	TodayCheckedIn                              //今天已签到
	HeaderParamError                            //header参数错误
	RefreshMemberList                           //刷新当前会员列表
	LotteryTimeOut                              //抽奖排队超时
	CheatUser                                   //作弊玩家
	NotExist                                    //数据不存在
	NotHasSendGoods                             //商品还未发货
	IsAppraised                                 //商品已评价
	NotHasSignedGoods                           //商品未签收不能评价
	NotExistUser                                //不存在此用户
	NotExistGoods                               //不存在此商品
	UpperLimit                                  //今日你的提交数已达到上限
	BlackList                                   //已成为黑名单用户
	PayeeReview                                 //金额数额过大,提现审核中
	DateTooLong                                 //查询日期过长
	WebPermissionDenied                         //web权限拒绝
)

type Platform int //平台
const (
	APP Platform = iota
	H5
)

type PayPlatform int //付费平台
const (
	WXPay PayPlatform = iota
)

type OrderType int //订单状态
const (
	NotApply OrderType = iota
	Apply
)

type SendGoodsType int8 //发货状态
const (
	NotSend SendGoodsType = iota
	Sended
)

type SignForType int8 //商品签收状态
const (
	NotSignGoods SignForType = iota
	GoodsSigned
)

type AgentLevelType int8 //代理级别
const (
	AgentLevel1 AgentLevelType = iota + 1
	AgentLevel2
	AgentLevel3
)

type GoldChangeType int //金币变化类型
const (
	DepositType     GoldChangeType = iota //充值金币类型
	RushConsumeType                       //闯关消耗的金币类型
	ProxyRewardType                       //获取佣金提成的金币
	GrantType                             //客户赠送的金币
	DeductType                            //客户收回的金币
	CheckInType                           //签到赠送的金币
	DownLoadAppType                       //下载app送的金币
	AppraisedType                         //评价商品送的金币
	RegisterType                          //新人福利
)

type UserAppraiseType int //用户评价类型
const (
	TAndP UserAppraiseType = iota //图文
	OnlyP                         //只有图
	OnlyT                         //只有文字
)

type GoodsDataType int //商品类型
const (
	RushGoods    GoodsDataType = iota //闯关商品
	LotteryGoods                      //抽奖商品
)

type DrawCashState int //商品类型
const (
	DrawCashReview  DrawCashState = iota //审核中
	DrawCashSucceed                      //提现成功
	DrawCashFailed                       //提现失败
)

type DrawCashArrivalType int //到账类型
const (
	DrawCashArrivalWX  DrawCashArrivalType = iota //微信钱包
	DrawCashArrivalZFB                            //支付宝
)

const UserIdField = "UserId"
const IsBlackListField = "IsBlackList"

type WXUserData struct {
	OpenId     string   `json:"openid"`
	NickName   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgUrl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionId    string   `json:"unionid"`
}

type OrderForm struct {
	Id           string       //充值订单号
	UserId       int          //用户id
	Money        int64        //充值额度
	CreatedAt    int64        //订单创建时间
	Desc         string       //商品描述
	OpenId       string       //用户openid
	PurchaseType PurchaseType //购买类型
	PurchaseId   int          //商品id或者充值id
	Platform     Platform     //充值平台
}

type WX_PayInfo struct {
	Appid      string
	Mch_id     string
	Trade_type string
	PaySecret  string
	OrderForm  *OrderForm
}

const GameStateKey = "GameState_"

type GameState struct {
	LevelId   string //关卡id
	UserId    int    //用户id
	GameTime  int    //游戏时间
	StartTime int64  //游戏开始时间
}

type TmpUser struct {
	NickName string //昵称
	Avatar   string //头像
	Desc     string //描述
}

//body
type H5LoginBody struct {
	Code     string `json:"code"`     //身份标识
	Referrer int    `json:"referrer"` //推荐人id,为UseId
}

type AppLoginBody struct {
	Code string `json:"code"` //身份标识
}

type LevelPassBody struct {
	Id       int      `json:"id"`       //关卡标识符
	Platform Platform `json:"platform"` //平台
}

type LevelPassFailedBody struct {
	Id     int `json:"id"`     //关卡标识符
	Number int `json:"number"` //第几只口红失败的
}

type ApplySendBody struct {
	OrderNumber string `json:"number"` //订单号
	LinkMan     string `json:"linkman"`
	PhoneNumber string `json:"phone"`
	Address     string `json:"addr"`
	Remark      string `json:"remark"`
}

type AgentLevelNBody struct {
	Ids []int `json:"ids"`
}

type PayeeBody struct {
	Amount float64 `json:"amount"`
}

type AppPurchaseBody struct {
	Id    string `json:"id"`
	Class string `json:"class"`
}

type WXPayResultNoticeBody struct {
	AppId          string `xml:"appid"`
	Bank_type      string `xml:"bank_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Fee_type       string `xml:"fee_type"`
	Is_subscribe   string `xml:"is_subscribe"`
	Mch_id         string `xml:"mch_id"`
	Nonce_str      string `xml:"nonce_str"`
	Openid         string `xml:"openid"`
	Out_trade_no   string `xml:"out_trade_no"`
	Result_code    string `xml:"result_code"`
	Return_code    string `xml:"return_code"`
	Sign           string `xml:"sign"`
	Time_end       string `xml:"time_end"`
	Total_fee      int    `xml:"total_fee"`
	Trade_type     string `xml:"trade_type"`
	Transaction_id string `xml:"transaction_id"`
}

type UserAppraiseBody struct {
	Number    string        `json:"number"` //订单编号
	ImgNames  []string      `json:"imgnames"`
	Desc      string        `json:"desc"`
	GoodsType GoodsDataType `json:"goodstype"`
}

type SuggestBody struct {
	Desc string `json:"desc"`
}
type ComplaintBody struct {
	ComplaintType string `json:"type"`
	Desc          string `json:"desc"`
}

//Response
type ResponseGoodsData struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	ImgUrl    string   `json:"imgurl"`     //图片地址
	ReImgUrl  string   `json:"reimgurl"`   //推荐类型图片地址
	RushPrice int64    `json:"rushprice"`  //闯关价
	Avatar    []string `json:"avatarlist"` //用户头像列表
	SendedOut int64    `json:"sendedout"`  //已发货数量
}

type ResponseAdData struct {
	ImgUrl string `json:"imgurl"` //图片地址
	IsJump int    `json:"isjump"` //是否跳转
	JumpTo string `json:"jumpto"` //跳转标识
}

type ResponseSaveGame struct {
	GoodsId int `json:"goodsid"` //商品id
}

type ResponseHomeData struct {
	Ad []*ResponseAdData `json:"ad"`
}

type ResponseFreeRougeGame struct {
	Level      int `json:"level"`      //关卡编号
	RougeCount int `json:"rougecount"` //口红数
	Difficulty int `json:"difficulty"` //难度系数
	GameTime   int `json:"gametime"`   //游戏时间,单位为秒
}

type ResponseOrderInfo struct {
	Id        int    `json:"id"`        //商品id
	Name      string `json:"name"`      //商品名
	ImgUrl    string `json:"imgurl"`    //图片地址
	Price     int64  `json:"price"`     //价格
	PriceDesc string `json:"pricedesc"` //价格描述
	Count     int    `json:"count"`     //数量
	Number    string `json:"number"`    //订单号
	Remark    string `json:"remark"`    //备注
}

type ResponseNotSendGoods struct {
	ResponseOrderInfo
	IsRemind int `json:"isremind"` //提醒发货
}

type ResponseHasSendedGoods struct {
	ResponseOrderInfo
	ExpressNumber string `json:"expressnumber"` //快递单号
	ExpressAgency string `json:"expressagency"`
}

type ResponseAppraiseOrder struct {
	ResponseOrderInfo
	IsAppraised int `json:"isappraised"` //是否已评价
}

type ResponseUserCommission struct {
	Avatar   string  `json:"avatar"`
	NickName string  `json:"nickname"`
	Total    float64 `json:"total"` //个人赚取的佣金数
}
type ResponseSelfCommission struct {
	Base ResponseUserCommission `json:"base"`
	Rank int                    `json:"rank"`
}
type ResponseCommissionRank struct {
	List []*ResponseUserCommission `json:"list"`
	Self *ResponseSelfCommission   `json:"self"`
}

type ResponseCommissionInfo struct {
	AgencyLevel int8    `json:"agencylevel"`
	CreatedAt   int64   `json:"time"`
	EarnBalance float64 `json:"balance"`
}

type ResponseAgent struct {
	Avatar      string  `json:"avatar"`
	NickName    string  `json:"nickname"`
	CreatedAt   int64   `json:"time"`
	EarnBalance float64 `json:"balance"`
	UserId      int     `json:"id"`
}

type ResponseAgentInfo struct {
	Total   float64          `json:"total"`
	Balance float64          `json:"balance"`
	Agent   []*ResponseAgent `json:"agent"`
}

type ResponseDrawCashInfo struct {
	Charge      float64       `json:"money"`
	Poundage    float64       `json:"poundage"`
	State       DrawCashState `json:"state"`
	PaymentTime string        `json:"time"`
	ArrivalType string        `json:"type"` //到账类型
}

type ResponseDrawCashRule struct {
	Balance       float64 `json:"balance"`       //当前佣金
	MinCharge     float64 `json:"mincharge"`     //最低提现额度
	MinPoundage   float64 `json:"minpoundage"`   //最低提现手续费
	MaxDrawCount  int     `json:"times"`         //每日最大提现次数
	PoundagePer   int     `json:"poundageper"`   //提现手续费百分比 0~100值
	RequireVerify float64 `json:"requireverify"` //超过多少钱需要审核
}

type ResponseDepositParams struct {
	GoldCount int64            `json:"goldcount"` //当前佣金
	List      []*DepositParams `json:"list"`
}

type ResponsSaveGame struct {
	Goods     ResponseGoodsData `json:"goods"`
	CurrentId int               `json:"currentid"` //当前关卡标识符id
}

type ResponseGameContinue struct {
	GameRecord *ResponsSaveGame    `json:"gamerecord"`
	List       []*PayModeRougeGame `json:"list"`
	Label      int                 `json:"label"` //干扰项,1为干扰,0为不干扰
}

type DownLoadInfo struct {
	Link                 string `json:"link"`
	ImgUrl               string `json:"imgurl"`
	IsGotDownLoadAppGift int    `json:"isgotdownloadappgift"`
	GoldCount            int64  `json:"gold"`
}
type RegisterGift struct {
	IsGotRegisterGift int   `json:"isgotregistergift"`
	GoldCount         int64 `json:"gold"`
}
type AppraiseAdInfo struct {
	AppraiseAd string `json:"appraisead"`
	GoldCount  int64  `json:"gold"`
}

type ResponseUserInfo struct {
	GoldCount    int64           `json:"goldcount"`
	Avatar       string          `json:"avatar"`
	NickName     string          `json:"nickname"`
	UserId       int             `json:"userid"`
	DLQCount     int64           `json:"dlqcount"`
	DFHCount     int64           `json:"dfhcount"`
	YFHCount     int64           `json:"yfhcount"`
	DPJCount     int64           `json:"dpjCount"`
	DownLoad     *DownLoadInfo   `json:"download"`
	RegisterGift *RegisterGift   `json:"registergift"`
	AppraiseAd   *AppraiseAdInfo `json:"appraisead"`
}

type ResponsGoldData struct {
	GoldCount int64                `json:"goldcount"`
	List      []*ResponsGoldChange `json:"list"`
}

type ResponsGoldChange struct {
	ChangeType GoldChangeType `json:"changetype"`
	GoldChange int64          `json:"goldchange"`
	CreatedAt  int64          `json:"created_at"`
}

type ResponsH5Pay struct {
	AppId     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	TimeStamp string `json:"timeStamp"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type ResponsAppPay struct {
	PartnerId string `json:"partnerid"`
	NonceStr  string `json:"noncestr"`
	Package   string `json:"package"`
	TimeStamp string `json:"timestamp"`
	PrepayId  string `json:"prepayid"`
	Sign      string `json:"sign"`
}

type ResponseShareData struct {
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Link   string `json:"link"`
	ImgUrl string `json:"imgurl"`
}

type ResponseCheckInData struct {
	IsCheckedIn int   `json:"ischeckedin"`
	Gold        int64 `json:"gold"`
}

type ResponseMemberData struct {
	VipId   int    `json:"vipid"`
	Name    string `json:"name"`
	A1Money int    `json:"a1money"`
	A2Money int    `json:"a2money"`
	A3Money int    `json:"a3money"`
	A1Gold  int    `json:"a1gold"`
	A2Gold  int    `json:"a2gold"`
	A3Gold  int    `json:"a3gold"`
	Price   int64  `json:"price"`
	Enable  int    `json:"enable"` //是否可用,0为不可用,1为可用
	Owned   int    `json:"owned"`  //是否拥有,0为没拥有,1为拥有
}

type ResponseLotteryGoods struct {
	ImgUrl string `json:"imgurl"`
	Price  int64  `json:"price"`
}

type ResponseLotteryGoodsInfo struct {
	IsGotReward int                     `json:"isgotreward"` //0为没中,1为中奖
	List        []*ResponseLotteryGoods `json:"list"`
}

type ResponseLotteryGoodsSucceedHistory struct {
	NickName string `json:"nickname"` //昵称
	Avatar   string `json:"avatar"`   //头像
	DescInfo string `json:"desc"`     //用户中奖描述
}

type LotteryGoodsOrderInfo struct {
	Name      string `json:"name"`      //商品名称
	ImgUrl    string `json:"imgurl"`    //商品图片地址
	PriceDesc string `json:"pricedesc"` //价格
	Count     int    `json:"count"`     //商品标签
	Number    string `json:"number"`    //订单号
	Remark    string `json:"remark"`    //备注
	IsSended  int    `json:"issended"`  //是否已发货
}

type ResponseLotteryOrderNotSend struct {
	LotteryGoodsOrderInfo
	PostAge string `json:"postage"` //邮费
}
type ResponseLotteryOrderHasSended struct {
	LotteryGoodsOrderInfo
	ExpressNumber string `json:"expressnumber"`
	ExpressAgency string `json:"expressagency"`
	IsAppraised   int    `json:"isappraised"`
}

type ResponsePosters struct {
	ImgUrl   string `json:"imgurl"`
	Icon     string `json:"icon"`
	Location int    `json:"location"`
}

type ResponseSharePosters struct {
	Posters []*ResponsePosters `json:"posters"`
	QRcode  string             `json:"qrcode"`
}

type ResponseRewardHistory struct {
	Avatar  string `json:"avatar"`
	Desc    string `json:"desc"`
	GoodsId int    `json:"goodsid"`
}

type ResponseGoodsClass struct {
	Name   string `json:"name"`
	ImgUrl string `json:"imgurl"`
	Id     int    `json:"id"`
}

type ResponseUserAppraise struct {
	ImgUrls   []string `json:"imgurls"` //图片url
	Desc      string   `json:"desc"`
	NickName  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	GoodsName string   `json:"goodsname"`
	CreatedAt int64    `json:"time"`
}

type GoodsDetailAppraise struct {
	Total        int64                 `json:"total"`
	UserAppraise *ResponseUserAppraise `json:"userappraise"`
}

type ResponseGoodsDetail struct {
	GoodsCover []string             `json:"cover"`     //封面图片
	Price      int64                `json:"price"`     //购买价
	RushPrice  int64                `json:"rushprice"` //闯关价
	PriceDesc  string               `json:"pricedesc"` //价格秒速
	SendedOut  int64                `json:"sendedout"` //发货量
	ReIcon     string               `json:"icon"`      //活动图标
	GoodsName  string               `json:"name"`      //商品名称
	Words16    string               `json:"words16"`   //16字
	Percent    int                  `json:"percent"`   //百分比
	TmpUsers   []*TmpUser           `json:"users"`     //获奖名单
	Appraise   *GoodsDetailAppraise `json:"appraise"`
	DetailImgs []string             `json:"detail"`
}

type MemberInfo struct {
	Name    string `json:"name"`
	A1Money int    `json:"a1money"`
	A2Money int    `json:"a2money"`
	A3Money int    `json:"a3money"`
}

type ResponseAgencyPage struct {
	Member    *MemberInfo `json:"member"`
	Avatar    []string    `json:"avatar"`
	InviteUrl string      `json:"inviteurl"`
}

func ResponseLoginData(u_data *UserInfo) map[string]interface{} {
	if u_data == nil {
		return nil
	}
	mp := make(map[string]interface{})
	mp["token"] = u_data.Token
	return mp
}
