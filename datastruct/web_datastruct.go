package datastruct

type UpdateSendInfoBody struct {
	OrderNumber   string `json:"number"`        //订单号
	ExpressNumber string `json:"expressnumber"` //快递单号
	ExpressAgency string `json:"agency"`        //快递单号
}

type UpdateSignForBody struct {
	OrderNumber string `json:"number"` //订单号
}

//body
type WebLoginBody struct {
	Account string `json:"account"`
	Pwd     string `json:"pwd"`
}

type EditGoodsBody struct {
	Goodsid       int              `json:"goodsid"`
	Name          string           `json:"name"`
	Classid       int              `json:"classid"`
	Base64str     []string         `json:"base64str"`
	Sortid        int              `json:"sortid"`
	Price         int64            `json:"price"`
	Pricedesc     string           `json:"pricedesc"`
	Rushprice     int64            `json:"rushprice"`
	Rushpricedesc string           `json:"rushpricedesc"`
	Goodsdesc     string           `json:"goodsdesc"`
	Brand         string           `json:"brand"`
	LevelData     []GoodsLevelInfo `json:"level"`
	IsHidden      int              `json:"ishidden"`
	Count         int              `json:"count"`
	Type          int              `json:"type"`
	OriginalPrice int64            `json:"originalprice"`
	PostAge       int              `json:"postage"`
	LimitAmount   int64            `json:"limit"`
	SendedOut     int64            `json:"sendedout"` //已发数量
	ReClassid     int              `json:"reclassid"` //推荐类型id
	Words16       string           `json:"words16"`   //16字
	Percent       int              `json:"percent"`   //多少百分比闯关成功
}

type EditRandomLotteryGoodsBody struct {
	Goodsid     int    `json:"goodsid"`
	Name        string `json:"name"`
	Classid     int    `json:"classid"`
	Base64str   string `json:"base64str"`
	Price       int64  `json:"price"`
	Goodsdesc   string `json:"goodsdesc"`
	Brand       string `json:"brand"`
	Probability int    `json:"probability"`
	IsHidden    int    `json:"ishidden"`
}

type GoodsLevelInfo struct {
	Time       int `json:"time"`
	Count      int `json:"count"`
	Difficulty int `json:"difficulty"`
}

type WebGetGoodsBody struct {
	Name      string `json:"name"`
	Classid   int    `json:"classid"`
	IsHidden  int    `json:"ishidden"`
	PageIndex int    `json:"pageindex"`
	PageSize  int    `json:"pagesize"`
}

type WebEditDomainBody struct {
	EntryDomain       string `json:"entrydomain"`
	EntryPageUrl      string `json:"entrypageurl"`
	AuthDomain        string `json:"authdomain"`
	AppDomain         string `json:"appdomain"`
	DownLoadUrl       string `json:"download"`
	DirectDownLoadUrl string `json:"directdownload"`
	IOSApp            string `json:"iosapp"`
	AndroidApp        string `json:"androidapp"`
}

type BlackListJumpBody struct {
	BLJumpTo string `json:"blacklist"`
	PCJumpTo string `json:"pcjumpto"`
}

type RushOrderState int8 //闯关订单状态
const (
	RushPayed RushOrderState = iota
	RushFinishedApply
	RushFinishedNotApply
	RushFinishedFailed
)

type GetRushOrderBody struct {
	UserName  string         `json:"username"`
	GoodsName string         `json:"goodsname"`
	State     RushOrderState `json:"state"`
	StartTime int64          `json:"starttime"`
	EndTime   int64          `json:"endtime"`
	PageIndex int            `json:"pageindex"`
	PageSize  int            `json:"pagesize"`
}

type GetPurchaseBody struct {
	UserName  string    `json:"username"`
	GoodsName string    `json:"goodsname"`
	State     OrderType `json:"state"`
	StartTime int64     `json:"starttime"`
	EndTime   int64     `json:"endtime"`
	PageIndex int       `json:"pageindex"`
	PageSize  int       `json:"pagesize"`
}

type GetSendGoodsBody struct {
	OrderId   string        `json:"orderid"`
	UserName  string        `json:"username"`
	GoodsName string        `json:"goodsname"`
	State     SendGoodsType `json:"state"`
	StartTime int64         `json:"starttime"`
	EndTime   int64         `json:"endtime"`
	PageIndex int           `json:"pageindex"`
	PageSize  int           `json:"pagesize"`
}

type DefaultAgencyBody struct {
	Agent1Gold  int `json:"agent1Gold"`
	Agent1Money int `json:"agent1Money"`
	Agent2Gold  int `json:"agent2Gold"`
	Agent2Money int `json:"agent2Money"`
	Agent3Gold  int `json:"agent3Gold"`
	Agent3Money int `json:"agent3Money"`
}

type EditLevelDataBody struct {
	Id       int    `json:"id"`       //id
	Name     string `json:"name"`     //等级名称
	Level    int    `json:"level"`    //级别
	Price    int64  `json:"price"`    //购买价格
	IsHidden int    `json:"ishidden"` //0为显示,1为隐藏
	DefaultAgencyBody
}

type GetMemberLevelDataBody struct {
	Name string `json:"name"` //等级名称
}

type WebGetMembersBody struct {
	NickName    string `json:"name"`        //昵称
	IsBlacklist int    `json:"isblacklist"` //是否查询黑名单
	LevelId     int    `json:"levelid"`     //等级id
	PageIndex   int    `json:"pageindex"`
	PageSize    int    `json:"pagesize"`
}

type WebUpdateUserBlBody struct {
	UserId      int `json:"userid"`
	IsBlacklist int `json:"isblacklist"` //是否查询黑名单
}

type WebUpdateUserLevelBody struct {
	UserId  int `json:"userid"`
	LevelId int `json:"levelid"`
}

type WebAddGoldBody struct {
	UserId int   `json:"userid"`
	Gold   int64 `json:"gold"`
}

type WebGetAgencyInfoBody struct {
	UserId    int    `json:"userid"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	StartTime int64  `json:"starttime"`
	EndTime   int64  `json:"endtime"`
	PageIndex int    `json:"pageindex"`
	PageSize  int    `json:"pagesize"`
}

type WebServerInfoBody struct {
	Version    string `json:"version"`
	IsMaintain int    `json:"ismaintain"`
}

type WebEditGoodsClassBody struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsHidden int    `json:"ishidden"`
	SortId   int    `json:"sortid"`
	Icon     string `json:"icon"`
}

type WebQueryGoodsClassBody struct {
	Name     string `json:"name"`
	IsHidden int    `json:"ishidden"`
}

type WebQueryDepositInfoBody struct {
	Name      string `json:"nickname"`
	Platform  int    `json:"platform"`
	StartTime int64  `json:"starttime"`
	EndTime   int64  `json:"endtime"`
	PageIndex int    `json:"pageindex"`
	PageSize  int    `json:"pagesize"`
}

type WebQueryDrawInfoBody struct {
	Name      string        `json:"nickname"`
	State     DrawCashState `json:"state"`
	StartTime int64         `json:"starttime"`
	EndTime   int64         `json:"endtime"`
	PageIndex int           `json:"pageindex"`
	PageSize  int           `json:"pagesize"`
}

type WebUpdateMemberLevelBody struct {
	Id       int `json:"id"`
	IsHidden int `json:"ishidden"`
}

type WebNewsUserBody struct {
	StartTime   int64    `json:"starttime"`
	EndTime     int64    `json:"endtime"`
	RPlatform   Platform `json:"registerplatform"`
	PayPlatform Platform `json:"payplatform"`
}

type WebActiveUserBody struct {
	StartTime int64    `json:"starttime"`
	EndTime   int64    `json:"endtime"`
	RPlatform Platform `json:"registerplatform"`
}

type WebCommissionStatisticsBody struct {
	StartTime int64 `json:"starttime"`
	EndTime   int64 `json:"endtime"`
}

type WebMemberOrderBody struct {
	Name      string `json:"nickname"`
	LevelName string `json:"levelname"`
	StartTime int64  `json:"starttime"`
	EndTime   int64  `json:"endtime"`
	PageIndex int    `json:"pageindex"`
	PageSize  int    `json:"pagesize"`
}

type WebDeleteMemberBody struct {
	Id int `json:"id"`
}

type WebResponseLotteryGoodsSendStateBody struct {
	OrderNumber   string `json:"number"`        //订单号
	ExpressNumber string `json:"expressnumber"` //快递单号
	ExpressAgency string `json:"agency"`        //快递机构
	LinkMan       string `json:"linkman"`       //联系人
	PhoneNumber   string `json:"phone"`         //联系号码
	Address       string `json:"addr"`          //详细地址
}

type WebRushLimitSettingBody struct {
	LotteryCount int    `json:"count"`      //用户达到中奖次数后
	Diff2        int    `json:"diff2"`      //改变第二关难度系数
	Diff2r       int    `json:"diff2r"`     //改变第二关口红数
	Diff2t       int    `json:"diff2t"`     //改变第二关游戏时间
	Diff3        int    `json:"diff3"`      //改变第三关难度系数
	Diff3r       int    `json:"diff3r"`     //改变第三关口红数
	Diff3t       int    `json:"diff3t"`     //改变第三关游戏时间
	CheatCount   int    `json:"cheatcount"` //达到作弊条件
	CheatTips    string `json:"tips"`       //提示语
}

type WebEditReClassBody struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
type WebQueryReClassBody struct {
	Name string `json:"name"`
}
type WebQuerySharePostersBody struct {
	IsHidden int `json:"ishidden"`
}

type WebEditSharePosterBody struct {
	Id       int    `json:"id"`
	ImgUrl   string `json:"imgurl"`
	Icon     string `json:"icon"`
	IsHidden int    `json:"ishidden"` //0为显示,1为隐藏
	SortId   int    `json:"sortid"`   //排序编号,根据该字段降序查询
	Location int    `json:"location"` //位置
}

type WebHiddenPostersBody struct {
	Id       int `json:"id"`
	IsHidden int `json:"ishidden"`
}

type WebEditUserAppraiseBody struct {
	Id        int           `json:"id"`
	ImgNames  []string      `json:"imgnames"`
	UserId    int           `json:"userid"`
	Desc      string        `json:"desc"`
	GoodsId   int           `json:"goodsid"`
	IsPassed  int           `json:"ispassed"`
	GoodsType GoodsDataType `json:"goodstype"`
	TimeStamp int64         `json:"timestamp"`
}

type WebDeleteUserAppraiseBody struct {
	Id int `json:"id"`
}

type WebQueryUserAppraiseBody struct {
	UserName  string        `json:"nickname"`
	GoodsType GoodsDataType `json:"goodstype"`
	GoodsName string        `json:"goodsname"`
	Desc      string        `json:"desc"`
	IsPassed  int           `json:"ispassed"`
	PageIndex int           `json:"pageindex"`
	PageSize  int           `json:"pagesize"`
}

type WebEditGoodsDetailBody struct {
	ImgNames []string `json:"imgnames"`
	GoodsId  int      `json:"goodsid"`
}

type WebQuerySuggestBody struct {
	PageIndex int    `json:"pageindex"`
	PageSize  int    `json:"pagesize"`
	UserName  string `json:"nickname"`
	Desc      string `json:"desc"`
}

type WebQueryComplaintBody struct {
	WebQuerySuggestBody
	ComplaintType string `json:"type"`
}

type WebQueryAdBody struct {
	Location AdLocation `json:"location"`
	Platform Platform   `json:"platform"`
	IsHidden int        `json:"ishidden"`
}

type WebEditPermissionUserBody struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Pwd           string `json:"pwd"`
	LoginName     string `json:"loginname"`
	PermissionIds []int  `json:"permission"`
}

type WebUserPwdBody struct {
	OldPwd string `json:"oldpwd"`
	NewPwd string `json:"newpwd"`
}

type WebCreateGoldPosterBody struct {
	StartTime int64 `json:"starttime"`
	EndTime   int64 `json:"endtime"`
	GoldCount int   `json:"gold"`
}

//Web Response
type WebResponseGoods struct {
	Total int                     `json:"total"`
	List  []*WebResponseGoodsData `json:"list"`
}

type WebResponseGoodsData struct {
	Id            int               `json:"id"`            //商品id
	Name          string            `json:"name"`          //商品名称
	ImgUrls       []string          `json:"imgurls"`       //商品图片地址
	Price         int64             `json:"price"`         //直接购买价格
	PriceDesc     string            `json:"pricedesc"`     //价格描述
	RushPrice     int64             `json:"rushprice"`     //闯关价
	RushPriceDesc string            `json:"rushpricedesc"` //闯关价描述
	GoodsDesc     string            `json:"goodsdesc"`     //商品别名
	Brand         string            `json:"brand"`         //商品标签
	SortId        int               `json:"sortid"`        //排序编号,根据该字段降序查询
	GoodsClass    string            `json:"goodsclass"`    //商品类型
	IsHidden      int               `json:"ishidden"`      //是否隐藏商品
	Count         int               `json:"count"`         //库存
	Type          int               `json:"type"`          //1为实体商品,0为虚拟商品
	OriginalPrice int64             `json:"originalprice"` //商品原价
	PostAge       int               `json:"postage"`       //邮费
	LimitAmount   int64             `json:"limitamount"`   //限制金额
	FailedTotal   int64             `json:"failedtotal"`   //失败金额统计
	Level         []*GoodsLevelInfo `json:"level"`         //关卡数据
	SendedOut     int64             `json:"sendedout"`     //已发数量
	ReClassIcon   string            `json:"reclassicon"`   //推荐类型icon
	Words16       string            `json:"words16"`       //16字
	Percent       int               `json:"percent"`       //多少百分比闯关成功
}

type RushOrderInfo struct {
	UserName    string         `json:"username"`
	Avatar      string         `json:"avatar"`
	GoodsImgUrl string         `json:"goodsimgurl"`
	GoodsName   string         `json:"goodsname"`
	RushPrice   int64          `json:"price"`
	State       RushOrderState `json:"state"`
	ReMark      string         `json:"remark"`
	Time        int64          `json:"time"`
}

type WebResponseRushOrderInfo struct {
	Total              int              `json:"total"`              //总记录数
	TotalAmount        int64            `json:"totalamount"`        //总金额
	TodayCount         int              `json:"todaycount"`         //今日记录数
	TodayAmount        int64            `json:"todayamount"`        //今日金额
	CurrentTotal       int              `json:"currenttotal"`       //当前状态总记录数
	CurrentTotalAmount int64            `json:"currenttotalamount"` //当前状态总金额
	List               []*RushOrderInfo `json:"list"`
}

type PurchaseOrderInfo struct {
	UserName    string    `json:"username"`
	Avatar      string    `json:"avatar"`
	GoodsImgUrl string    `json:"goodsimgurl"`
	GoodsName   string    `json:"goodsname"`
	Price       int64     `json:"price"`
	State       OrderType `json:"state"`
	ReMark      string    `json:"remark"`
	Time        int64     `json:"time"`
}

type WebResponsePurchaseOrderInfo struct {
	Total              int                  `json:"total"`              //总记录数
	TotalAmount        int64                `json:"totalamount"`        //总金额
	TodayCount         int                  `json:"todaycount"`         //今日记录数
	TodayAmount        int64                `json:"todayamount"`        //今日金额
	CurrentTotal       int                  `json:"currenttotal"`       //当前状态总记录数
	CurrentTotalAmount int64                `json:"currenttotalamount"` //当前状态总金额
	List               []*PurchaseOrderInfo `json:"list"`
}

type SendGoodsInfo struct {
	OrderId      string                `json:"orderid"`
	UserName     string                `json:"username"`
	Avatar       string                `json:"avatar"`
	GoodsImgUrl  string                `json:"goodsimgurl"`
	GoodsName    string                `json:"goodsname"`
	Count        int                   `json:"count"`
	Price        int64                 `json:"price"`
	State        SendGoodsType         `json:"state"`
	ReMark       string                `json:"remark"`
	Receiver     *ReceiverForSendGoods `json:"receiver"`
	Sender       *SenderForSendGoods   `json:"sender"`
	Time         int64                 `json:"time"`
	SignForState SignForType           `json:"signstate"`
}

type SenderForSendGoods struct {
	ExpressNumber string `json:"number"` //快递单号
	ExpressAgency string `json:"name"`   //快递机构
}

type ReceiverForSendGoods struct {
	LinkMan     string `json:"linkman"` //联系人
	PhoneNumber string `json:"phone"`   //联系号码
	Address     string `json:"addr"`    //详细地址
	Remark      string `json:"remark"`  //备注
}

type WebResponseSendGoodsInfo struct {
	Total              int              `json:"total"`              //总记录数
	TotalAmount        int64            `json:"totalamount"`        //总金额
	TodayCount         int              `json:"todaycount"`         //今日记录数
	TodayAmount        int64            `json:"todayamount"`        //今日金额
	CurrentTotal       int              `json:"currenttotal"`       //当前状态总记录数
	CurrentTotalAmount int64            `json:"currenttotalamount"` //当前状态总金额
	List               []*SendGoodsInfo `json:"list"`
}

type WebResponseMembersInfo struct {
	Members      []*WebResponseMember `json:"members"`
	TodayCreated int                  `json:"todaycreated"`
	Total        int                  `json:"total"`
	CurrentTotal int                  `json:"currenttotal"`
}

type WebResponseMember struct {
	Id           int     `json:"id"`
	NickName     string  `json:"name"`
	Avatar       string  `json:"avatar"`
	GoldCount    int64   `json:"goldcount"`
	Balance      float64 `json:"balance"`      //可提现佣金
	BalanceTotal float64 `json:"balancetotal"` //佣金总额
	LevelName    string  `json:"levelname"`    //等级名称
	LevelId      int     `json:"levelid"`      //等级id
	IsBlacklist  int     `json:"isblacklist"`  //是否为黑名单
	CreateTime   int64   `json:"createtime"`   //创建时间
}

type WebAgencyStatistics struct {
	Count         int     `json:"count"`
	DepositTotal  int64   `json:"deposittotal"`
	PayRushTotal  int64   `json:"payrushtotal"`
	PurchaseTotal float64 `json:"purchasetotal"`
}

type WebAgencyUser struct {
	NickName      string  `json:"name"`
	Avatar        string  `json:"avatar"`
	DepositTotal  int64   `json:"deposittotal"`
	PayRushTotal  int64   `json:"payrushtotal"`
	PurchaseTotal float64 `json:"purchasetotal"`
	GoldCount     int64   `json:"gold_count"` //金币余额
	Balance       float64 `json:"balance"`    //佣金余额
	CreatedAt     int64   `json:"time"`       //时间
}

type WebResponseAgencyData struct {
	Statistics   []*WebAgencyStatistics `json:"statistics"`
	Users        []*WebAgencyUser       `json:"list"`
	CurrentTotal int                    `json:"currenttotal"`
}

type WebResponseEditDomain struct {
	EntryDomain       string `json:"entrydomain"`
	EntryPageUrl      string `json:"entrypageurl"`
	AuthDomain        string `json:"authdomain"`
	AppDomain         string `json:"appdomain"`
	DownLoadUrl       string `json:"download"`
	DirectDownLoadUrl string `json:"directdownload"`
	IOSApp            string `json:"iosapp"`
	AndroidApp        string `json:"androidapp"`
	Version           string `json:"version"`
}

type WebDepositAgency struct {
	Gold  int64   `json:"gold"`
	Money float64 `json:"money"`
}

type WebDepositUser struct {
	NickName   string              `json:"name"`
	Avatar     string              `json:"avatar"`
	Pay        int64               `json:"pay"`
	CreateTime int64               `json:"createtime"`
	Agency     []*WebDepositAgency `json:"agency"`
	Platform   int                 `json:"platform"`
}

type WebResponseDepositInfo struct {
	Count        int64             `json:"count"`        //总记录数
	Amount       float64           `json:"amount"`       //总充值金额
	TodayCount   int64             `json:"todaycount"`   //今日记录数
	TodayAmount  float64           `json:"todayamount"`  //今日金额
	EarnGold     float64           `json:"earngold"`     //总提成的金币
	EarnMoney    float64           `json:"earnmoney"`    //总提成的佣金
	CurrentTotal int64             `json:"currenttotal"` //当前状态总记录数
	Users        []*WebDepositUser `json:"list"`
}

type WebDrawUser struct {
	Id          int           `json:"id"` //提现id
	NickName    string        `json:"name"`
	Avatar      string        `json:"avatar"`
	PaymentNo   string        `json:"paymentno"` //提现成功单号
	CreateTime  string        `json:"time"`      //提现时间
	Charge      float64       `json:"charge"`
	Poundage    float64       `json:"poundage"`
	State       DrawCashState `json:"state"`
	ArrivalType string        `json:"type"` //到账类型
}

type WebResponseDrawInfo struct {
	Amount       float64        `json:"amount"`       //已提现金额
	TodayAmount  float64        `json:"todayamount"`  //今日已提现金额
	Poundage     float64        `json:"poundage"`     //总手续费
	CurrentTotal int64          `json:"currenttotal"` //当前状态总记录数
	Users        []*WebDrawUser `json:"list"`
}

type WebResponseNotHiddenMember struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// type WebResponseOldAndNewUsers struct {
// 	OldUsers int `json:"oldusers"`
// 	NewUsers int `json:"newusers"`
// }

type WebMemberOrder struct {
	Id        int    `json:"id"`
	NickName  string `json:"name"`
	Avatar    string `json:"avatar"`
	LevelName string `json:"levelname"`
	Price     int64  `json:"price"`
	CreatedAt int64  `json:"time"`
}

type WebResponseMemberOrder struct {
	Count        int64             `json:"count"`        //总记录数
	Amount       float64           `json:"amount"`       //总金额
	TodayCount   int64             `json:"todaycount"`   //今日记录数
	TodayAmount  float64           `json:"todayamount"`  //今日金额
	CurrentTotal int64             `json:"currenttotal"` //当前状态总记录数
	List         []*WebMemberOrder `json:"list"`
}

type WebResponseRandomLotteryGoodsData struct {
	Id          int    `json:"id"`          //商品id
	Name        string `json:"name"`        //商品名称
	ImgUrl      string `json:"imgurl"`      //商品图片地址
	Price       int64  `json:"price"`       //直接购买价格
	GoodsDesc   string `json:"goodsdesc"`   //商品别名
	Brand       string `json:"brand"`       //商品标签
	IsHidden    int    `json:"ishidden"`    //是否隐藏商品
	Probability int    `json:"probability"` //概率
	GoodsClass  string `json:"goodsclass"`  //商品类型
}

type WebResponseRandomLotteryGoods struct {
	Total int                                  `json:"total"`
	List  []*WebResponseRandomLotteryGoodsData `json:"list"`
}

type WebResponseRandomLotteryPool struct {
	Current            float64 `json:"current"`      //该商品在当前池中的金额
	RandomLotteryCount int     `json:"lotterycount"` //用户达到中奖次数后
	Probability        int     `json:"probability"`  //概率
}

type WebResponseLotteryOrderInfo struct {
	Total        int              `json:"total"`        //总记录数
	TodayCount   int              `json:"todaycount"`   //今日记录数
	CurrentTotal int              `json:"currenttotal"` //当前状态总记录数
	List         []*SendGoodsInfo `json:"list"`
}

type WebResponseStatistics struct {
	GrowthCount      int     `json:"growthcount"`      //中奖用户的下级用户增长数
	AgentPayRate     float64 `json:"agentpayrate"`     //中奖用户的下级用户付费率
	UsersForPurchase int     `json:"usersforpurchase"` //直接购买用户数
	PurchaseAmount   int64   `json:"purchaseamount"`   //直接购买总金额
	UsersForDeposit  int     `json:"usersfordeposit"`  //充值用户数
	DepositAmount    float64 `json:"depositamount"`    //充值总金额
	UsersForPay      int     `json:"usersforpay"`      //付费总用户数
	UserPayAmount    float64 `json:"payamount"`        //付费总金额
	Date             string  `json:"date"`             //日期
}

type WebResponseActiveUsers struct {
	NewUsers    int64  `json:"newusers"`    //新用户数
	ActiveUsers int64  `json:"activeusers"` //活跃用户数
	Date        string `json:"date"`        //日期
}

type WebResponseCommissionStatistics struct {
	DepositTotal       float64 `json:"deposittotal"`    //充值总额
	BalanceTotal       float64 `json:"commissiontotal"` //总提成佣金
	DrawTotal          float64 `json:"drawtotal"`       //已提现总额
	RemainingDrawTotal float64 `json:"rdt"`             //待提现总额
	Profit             float64 `json:"profit"`          //营收利润
	Date               string  `json:"date"`            //日期
}

type WebResponseGoodsClass struct {
	Id       int    `json:"classid"`   //自增id
	Name     string `json:"classname"` //商品类型名称
	IsHidden int    `json:"ishidden"`  //0为显示,1为隐藏
	SortId   int    `json:"sortid"`    //排序编号,根据该字段降序查询
	Icon     string `json:"icon"`      //类型图标
}

type WebResponseReClass struct {
	Id   int    `json:"classid"`   //自增id
	Name string `json:"classname"` //商品类型名称
	Icon string `json:"icon"`      //类型图标
}

type WebResponseAvailableGoodsClass struct {
	Id   int    `json:"classid"`
	Name string `json:"name"`
}

type WebUserAppraise struct {
	Id        int           `json:"id"`      //自增id
	ImgUrls   []string      `json:"imgurls"` //图片url
	Desc      string        `json:"desc"`
	UserId    int           `json:"userid"`
	NickName  string        `json:"nickname"`
	Avatar    string        `json:"avatar"`
	GoodsName string        `json:"goodsname"`
	IsPassed  int           `json:"ispassed"`
	CreatedAt int64         `json:"time"`
	GoodsType GoodsDataType `json:"goodstype"`
	GoodsId   int           `json:"goodsid"`
}

type WebResponseUserAppraise struct {
	CurrentTotal int64              `json:"currenttotal"` //用于分页
	List         []*WebUserAppraise `json:"list"`
}

type WebResponseSCP struct {
	SCFD  int `json:"scfd"`  //每天建议次数
	CCFD  int `json:"ccfd"`  //每天投诉次数
	CCFBL int `json:"ccfbl"` //达到此投诉次数,自动进入黑名单
}

type WebResponseSuggest struct {
	Avatar    string `json:"avatar"`
	NickName  string `json:"nickname"`
	Desc      string `json:"desc"`
	CreatedAt int64  `json:"time"`
	Id        int    `json:"id"`
}

type WebResponseComplaint struct {
	WebResponseSuggest
	ComplaintType string `json:"type"`
}

type WebResponseSuggestInfo struct {
	CurrentTotal int64                 `json:"currenttotal"` //用于分页
	List         []*WebResponseSuggest `json:"list"`
}

type WebResponseComplaintInfo struct {
	CurrentTotal int64                   `json:"currenttotal"` //用于分页
	List         []*WebResponseComplaint `json:"list"`
}

type WebResponseDrawCashParams struct {
	MinCharge     float64 `json:"mincharge"`     //最低提现额度
	MinPoundage   float64 `json:"minpoundage"`   //最低提现手续费
	MaxDrawCount  int     `json:"times"`         //每日最大提现次数
	PoundagePer   int     `json:"poundageper"`   //提现手续费百分比 0~100值
	RequireVerify float64 `json:"requireverify"` //超过多少钱需要审核
}

type WebResponseAdInfo struct {
	Id       int        `json:"id"`       //自增id
	ImgUrl   string     `json:"imgurl"`   //图片地址
	IsJump   int        `json:"isjump"`   //是否跳转
	JumpTo   string     `json:"jumpto"`   //跳转标识
	SortId   int        `json:"sortid"`   //排序id
	Location AdLocation `json:"location"` //放置位置
	Platform Platform   `json:"platform"` //平台
	IsHidden int        `json:"ishidden"` //是否隐藏,1为隐藏,0为显示
}

type WebResponseGoldCoinGift struct {
	DownLoadAppGoldGift  int64 `json:"download"`          //下载app赠送金币数
	AppraisedGoldGift    int64 `json:"appraised"`         //评价后赠送金币数
	RegisterGoldGift     int64 `json:"register"`          //新人用户赠送金币数
	IsEnableRegisterGift int   `json:"isenableregister"`  //0关闭新人福利,1开启新人福利
	IsDrawCashOnlyApp    int   `json:"isdrawcashonlyapp"` //0关闭 ,1开启。是否只在app内提现
}

type MasterInfo struct {
	MasterId  int              `json:"masterid"`
	Name      string           `json:"name"`
	Secondary []*SecondaryInfo `json:"Secondary"`
}

type SecondaryInfo struct {
	SecondaryId int    `json:"secondaryid"`
	Name        string `json:"name"`
}

type WebResponsePermissionUser struct {
	Name       string        `json:"name"` //名称
	Token      string        `json:"token"`
	Permission []*MasterInfo `json:"permission"`
}

type WebResponseAllWebUser struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	LoginName     string `json:"loginname"`
	PermissionIds []int  `json:"permission"`
}
