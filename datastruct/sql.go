package datastruct

import "time"

const NULLSTRING = ""
const NULLID = -1

//用户信息表
type UserInfo struct {
	Id                   int      `xorm:"not null pk autoincr INT(11)"`     //自增id
	NickName             string   `xorm:"VARCHAR(255) not null"`            //昵称
	Avatar               string   `xorm:"VARCHAR(255) not null"`            //头像
	Sex                  int      `xorm:"TINYINT(1) not null"`              //1为男,2为女
	IsCheat              int      `xorm:"not null INT(11)"`                 //是否作弊,默认为0
	LotterySucceed       int      `xorm:"not null INT(11)"`                 //中奖次数
	Token                string   `xorm:"VARCHAR(255) not null"`            //标识符
	DepositTotal         int64    `xorm:"bigint not null"`                  //充值总额
	GoldCount            int64    `xorm:"bigint not null"`                  //当前金币数,默认为0
	BalanceTotal         float64  `xorm:"decimal(16,2) not null"`           //佣金总额
	Balance              float64  `xorm:"decimal(16,2) not null"`           //可提现佣金
	CreatedAt            int64    `xorm:"bigint not null"`                  //创建用户的时间
	LoginTime            int64    `xorm:"bigint not null"`                  //最近一次离开或者登陆的时间
	MemberIdentifier     string   `xorm:"VARCHAR(255) not null"`            //会员标识符id用于分销比例
	IsBlackList          int      `xorm:"TINYINT(1) not null default 0"`    //1为黑名单,0为正常
	PayRushTotal         int64    `xorm:"bigint not null default 0 "`       //闯关总消费
	PurchaseTotal        float64  `xorm:"decimal(16,2) not null default 0"` //直接购买总消费
	RandomLotterySucceed int      `xorm:"not null INT(11) default 0 "`      //随机抽奖中奖次数
	Platform             Platform `xorm:"TINYINT(1) not null default 1"`    //注册平台
	IsGotRegisterGift    int      `xorm:"TINYINT(1) not null default 1"`    //1为已领取新用户金币奖励,0为未领取
	IsGotDownLoadAppGift int      `xorm:"TINYINT(1) not null default 1"`    //1为已领取下载app金币奖励,0为未领取
}

//微信平台
type WXPlatform struct {
	UserId           int    `xorm:"not null pk INT(11)"`   //关联用户id
	WXUUID           string `xorm:"VARCHAR(255) not null"` //微信用户uuid
	PayOpenidForGZH  string `xorm:"VARCHAR(255) null"`     //用于微信公众号充值
	PayOpenidForKFPT string `xorm:"VARCHAR(255) null"`     //用于微信开放平台充值
	PayeeOpenid      string `xorm:"VARCHAR(255) null"`     //用于用户提现
}

//商品推荐类型
type RecommendedClass struct {
	Id   int    `xorm:"not null pk autoincr INT(11)" json:"id"` //自增id
	Name string `xorm:"VARCHAR(255) not null" json:"name"`      //名称
	Icon string `xorm:"VARCHAR(255) not null" json:"icon"`      //图片名称
}

//获取分享海报图片
type SharePosters struct {
	Id       int    `xorm:"not null pk autoincr INT(11)" json:"id"`        //自增id
	ImgName  string `xorm:"VARCHAR(255) not null" json:"img"`              //大图名称
	IconName string `xorm:"VARCHAR(255) not null" json:"icon"`             //小图名称
	SortId   int    `xorm:"not null INT(11) default 0" json:"sortid"`      //排序编号,根据该字段降序查询
	IsHidden int    `xorm:"TINYINT(1) not null default 0" json:"ishidden"` //0为显示,1为隐藏
	Location int    `xorm:"not null INT(11)" json:"location"`              //二维码位置
}

//商品类型
type GoodsClass struct {
	Id       int    `xorm:"not null pk autoincr INT(11)" json:"classid"`   //自增id
	Name     string `xorm:"VARCHAR(255) not null" json:"classname"`        //商品类型名称
	IsHidden int    `xorm:"TINYINT(1) not null default 0" json:"ishidden"` //0为显示,1为隐藏
	SortId   int    `xorm:"not null INT(11) default 0" json:"sortid"`      //排序编号,根据该字段降序查询
	Icon     string `xorm:"VARCHAR(255) not null" json:"icon"`             //类型图标
}

type AdLocation int

const (
	HomeAd     AdLocation = iota //首页
	DownLoadAd                   //下载
	AppraiseAd                   //评价
)

//广告栏
type AdInfo struct {
	Id       int        `xorm:"not null pk autoincr INT(11)"`  //自增id
	ImgName  string     `xorm:"VARCHAR(255) not null"`         //图片地址
	IsJump   int        `xorm:"not null INT(11)"`              //是否跳转
	JumpTo   string     `xorm:"VARCHAR(255) null"`             //跳转标识
	SortId   int        `xorm:"not null INT(11) default 0"`    //排序id
	Location AdLocation `xorm:"not null INT(11) default 0"`    //放置位置
	Platform Platform   `xorm:"not null INT(11) default 0"`    //平台
	IsHidden int        `xorm:"TINYINT(1) not null default 0"` //是否隐藏,1为隐藏,0为显示
}

//商品信息
type Goods struct {
	Id                  int    `xorm:"not null pk autoincr INT(11)"`  //自增id
	Name                string `xorm:"VARCHAR(255) not null"`         //商品名称
	ImgName             string `xorm:"VARCHAR(50) not null"`          //图片名称
	Price               int64  `xorm:"bigint not null"`               //直接购买价格
	PriceDesc           string `xorm:"VARCHAR(255) not null"`         //价格描述
	RushPrice           int64  `xorm:"bigint not null"`               //闯关价
	RushPriceDesc       string `xorm:"VARCHAR(255) not null"`         //闯关价描述
	GoodsDesc           string `xorm:"VARCHAR(255) not null"`         //商品描述
	Brand               string `xorm:"VARCHAR(255) not null"`         //所属品牌
	SortId              int    `xorm:"not null INT(11)"`              //排序编号,根据该字段降序查询
	SellerRecommendedId int    `xorm:"null INT(11)"`                  //店长推荐id
	GoodsClassId        int    `xorm:"null INT(11)"`                  //商品类型id
	IsHidden            int    `xorm:"TINYINT(1) not null default 0"` //是否隐藏商品
	Count               int    `xorm:"not null INT(11) default 0"`    //库存
	Type                int    `xorm:"TINYINT(1) not null default 1"` //1为实体商品,0为虚拟商品
	OriginalPrice       int64  `xorm:"bigint not null default 0"`     //商品原价
	PostAge             int    `xorm:"not null INT(11) default 0"`    //邮费
	SendedOut           int64  `xorm:"bigint not null default 0"`     //已发数量
	ReClassid           int    `xorm:"null INT(11) default 0"`        //推荐类型id
	Words16             string `xorm:"VARCHAR(255) not null"`         //16字
	Percent             int    `xorm:"not null INT(11) default 0"`    //多少百分比闯关成功
}

//商品池水
type GoodsRewardPool struct {
	GoodsId     int   `xorm:"not null pk INT(11)"`       //商品id
	Current     int64 `xorm:"bigint not null default 0"` //该商品在当前池中的金额
	LimitAmount int64 `xorm:"bigint not null default 0"` //限制金额
}

//口红体验模式
type FreeModeRougeGame struct {
	Id         int `xorm:"not null pk autoincr INT(11)"` //自增id
	Level      int `xorm:"INT(11) not null"`             //关卡编号
	RougeCount int `xorm:"INT(11) not null"`             //口红数
	Difficulty int `xorm:"INT(11) not null"`             //难度系数
	GameTime   int `xorm:"INT(11) not null"`             //游戏时间,单位为秒
}

//口红付费模式
type PayModeRougeGame struct {
	Id         int `xorm:"not null pk autoincr INT(11)" json:"id"` //自增id
	Level      int `xorm:"INT(11) not null" json:"level"`          //关卡编号
	RougeCount int `xorm:"INT(11) not null" json:"rougecount"`     //口红数
	Difficulty int `xorm:"INT(11) not null" json:"difficulty"`     //难度系数
	GameTime   int `xorm:"INT(11) not null" json:"gametime"`       //游戏时间,单位为秒
	GoodsId    int `xorm:"INT(11) not null" json:"goodsid"`        //商品id
}

//用户关卡记录,用户通关后或者失败后删除
type SaveGameInfo struct {
	UserId    int   `xorm:"not null pk INT(11)"` //关联用户id
	LevelId   int   `xorm:"INT(11) not null"`    //关卡自增id
	GoodsId   int   `xorm:"INT(11) not null"`    //商品id
	CreatedAt int64 `xorm:"bigint not null"`     //创建时间
}

//邀请表
type InviteInfo struct {
	Id        int   `xorm:"not null pk autoincr INT(11)"` //自增id
	Sender    int   `xorm:"not null INT(11)"`             //发送邀请者
	Receiver  int   `xorm:"not null INT(11)"`             //接受邀请者
	CreatedAt int64 `xorm:"bigint not null"`
}

//代理参数信息,各级代理用户充值时检查此表
type AgencyParams struct {
	Identifier          string `xorm:"VARCHAR(255) pk not null"`   //可为用户id,当为用户id时,即表示当前用户的代理提成;不为用户id,则是全局的提成
	Agency1MoneyPercent int    `xorm:"not null INT(11) default 0"` //1级代理提成 取值0～100
	Agency2MoneyPercent int    `xorm:"not null INT(11) default 0"` //2级代理提成 取值0～100
	Agency3MoneyPercent int    `xorm:"not null INT(11) default 0"` //3级代理提成 取值0～100
	Agency1GoldPercent  int    `xorm:"not null INT(11) default 0"` //1级代理提成 取值0～100
	Agency2GoldPercent  int    `xorm:"not null INT(11) default 0"` //2级代理提成 取值0～100
	Agency3GoldPercent  int    `xorm:"not null INT(11) default 0"` //3级代理提成 取值0～100
}

//用户付费口红游戏关卡失败记录
type PayModeRougeGameFailed struct {
	UserId             int   `xorm:"not null INT(11)"` //关联用户id
	CreatedAt          int64 `xorm:"bigint not null"`  //数据创建时间
	PayModeRougeGameId int   `xorm:"INT(11) not null"` //关联关卡数据
	RougeNumber        int   `xorm:"INT(11) not null"` //失败在第几只口红
}

//用户付费口红游戏关卡成功记录
type PayModeRougeGameSucceed struct {
	UserId             int   `xorm:"not null INT(11)"` //关联用户id
	CreatedAt          int64 `xorm:"bigint not null"`  //数据创建时间
	PayModeRougeGameId int   `xorm:"INT(11) not null"` //关卡标识符id
}

//用户付费口红游戏关卡成功历史
type PayModeRougeGameSucceedHistory struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"` //自增id
	NickName string `xorm:"VARCHAR(255) not null"`        //昵称
	Avatar   string `xorm:"VARCHAR(255) not null"`        //头像
	GoodsId  int    `xorm:"not null INT(11)"`
}

//用户订单表
type OrderInfo struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`  //自增id
	UserId     int       `xorm:"INT(11) not null"`              //用户id
	Number     string    `xorm:"VARCHAR(50) not null"`          //订单号
	GoodsId    int       `xorm:"INT(11) not null"`              //商品id
	OrderState OrderType `xorm:"TINYINT(1) not null"`           //订单状态 0为未申请,1为申请
	IsPurchase int       `xorm:"TINYINT(1) not null default 0"` //是否直接购买
	CreatedAt  int64     `xorm:"bigint not null"`               //创建时间
	Platform   Platform  `xorm:"TINYINT(1) not null default 1"` //平台
	IsRemind   int       `xorm:"TINYINT(1) not null"`           //0为未提醒发货,1为已提醒
}

//用户发货信息表
type SendGoods struct {
	OrderId        int64         `xorm:"bigint not pk null"`    //订单唯一标识符
	SendGoodsState SendGoodsType `xorm:"TINYINT(1) not null"`   //发货状态 0为未发货,1为已发货
	LinkMan        string        `xorm:"VARCHAR(255) not null"` //联系人
	PhoneNumber    string        `xorm:"VARCHAR(40) not null"`  //联系号码
	// Province       string        `xorm:"VARCHAR(255)  null"`            //省
	// City           string        `xorm:"VARCHAR(255)  null"`            //市
	// District       string        `xorm:"VARCHAR(255)  null"`            //区
	Address        string      `xorm:"VARCHAR(255)  not null"`        //详细地址
	ExpressNumber  string      `xorm:"VARCHAR(100) null"`             //快递单号
	ExpressAgency  string      `xorm:"VARCHAR(100) null"`             //快递机构
	CreatedAt      int64       `xorm:"bigint not null"`               //创建时间
	IsLotteryGoods int         `xorm:"TINYINT(1) not null default 0"` //是否抽奖商品,0不是,1是
	SignForState   SignForType `xorm:"TINYINT(1) not null default 0"` //签收状态 0为未签收,1为已签收
	IsAppraised    int         `xorm:"TINYINT(1) not null default 0"` //0为未评价,1为已评价
	Remark         string      `xorm:"VARCHAR(255) null"`
}

//用户赚取佣金记录表
type BalanceInfo struct {
	Id          int     `xorm:"not null pk autoincr INT(11)"` //自增id
	DepositId   int     `xorm:"not null INT(11)"`             //充值id
	AgencyLevel int8    `xorm:"TINYINT(1) not null"`          //代理级别取值为1,2,3
	FromUserId  int     `xorm:"INT(11)  not null"`            //充值用户id
	ToUserId    int     `xorm:"INT(11)  not null"`            //赚取佣金用户id
	EarnBalance float64 `xorm:"decimal(16,2) not null"`       //返现多少佣金
	EarnGold    int64   `xorm:"bigint not null default 0"`    //返多少金币
	CreatedAt   int64   `xorm:"INT(11)  not null"`
}

//用户提现参数
type DrawCashParams struct {
	Id            int     `xorm:"not null autoincr INT(11)"` //自增id
	MinCharge     float64 `xorm:"decimal(16,2) not null"`    //最低提现额度
	MinPoundage   float64 `xorm:"decimal(16,2) not null"`    //最低提现手续费
	MaxDrawCount  int     `xorm:"INT(11)  not null"`         //每日最大提现次数
	PoundagePer   int     `xorm:"INT(11)  not null"`         //提现手续费百分比 0~100值
	RequireVerify float64 `xorm:"decimal(16,2) not null"`    //超过多少钱需要审核
}

//用户提现记录
type DrawCashInfo struct {
	Id          int                 `xorm:"not null pk autoincr INT(11)"` //自增id
	UserId      int                 `xorm:"INT(11)  not null"`
	Charge      float64             `xorm:"decimal(16,2) not null"`
	Poundage    float64             `xorm:"decimal(16,2) not null"`
	TradeNo     string              `xorm:"VARCHAR(50) not null"` //自定义交易号
	CreatedAt   int64               `xorm:"bigint not null"`
	State       DrawCashState       `xorm:"TINYINT(1) not null"` //0为提现中,1为提现成功,2为提现失败
	PaymentNo   string              `xorm:"VARCHAR(100) null"`
	PaymentTime string              `xorm:"VARCHAR(100) null"`
	ArrivalType DrawCashArrivalType `xorm:"TINYINT(1) not null default 0"` //到账类型
	IpAddr      string              `xorm:"VARCHAR(100) null"`
	Origin      float64             `xorm:"decimal(16,2) not null"` //用户提款数目
}

//充值列表参数
type DepositParams struct {
	Id    int     `xorm:"not null pk autoincr INT(11)" json:"id"` //自增id
	Money float64 `xorm:"decimal(16,2) not null" json:"amount"`
}

//用户操作金币记录
type GoldChangeInfo struct {
	Id         int            `xorm:"not null pk autoincr INT(11)"` //自增id
	ChangeType GoldChangeType `xorm:"not null INT(11)"`             //金币变化类型
	UserId     int            `xorm:"not null INT(11)"`             //关联用户id
	VarGold    int64          `xorm:"bigint not null"`              //当前操作的金币数
	CreatedAt  int64          `xorm:"bigint not null"`              //用户充值时间
}

type EntryAddr struct {
	Url       string `xorm:"VARCHAR(255) pk not null"` //中转地址(非炮灰)
	CreatedAt int64  `xorm:"bigint not null"`
	PageUrl   string `xorm:"VARCHAR(255) not null"` //二维码落地页地址(炮灰)
}

type AuthAddr struct {
	Url       string `xorm:"VARCHAR(255) pk not null"` //授权中转地址(非炮灰)
	CreatedAt int64  `xorm:"bigint not null"`
}

type AppAddr struct {
	Url       string `xorm:"VARCHAR(255) pk not null"` //app落地页地址(炮灰)
	CreatedAt int64  `xorm:"bigint not null"`
}

type AppDownloadAddr struct {
	IosUrl            string `xorm:"VARCHAR(255) not null"`    //ios下载地址
	AndroidUrl        string `xorm:"VARCHAR(255) not null"`    //android下载地址
	DownLoadUrl       string `xorm:"VARCHAR(255) pk not null"` //下载中转地址(非炮灰)
	DirectDownLoadUrl string `xorm:"VARCHAR(255) not null"`    //直接下载落地页地址(炮灰)
	IosUrlSchemes     string `xorm:"VARCHAR(255) null"`
	AndroidUrlSchemes string `xorm:"VARCHAR(255) null"`
	CreatedAt         int64  `xorm:"bigint not null"`
}

type KfInfo struct {
	WX     string `xorm:"VARCHAR(50) pk not null"` //微信号
	QQ     string `xorm:"VARCHAR(50) not null"`    //qq号
	QRcode string `xorm:"VARCHAR(100) not null"`   //二维码图片url
}

type CheckInInfo struct {
	UserId          int   `xorm:"not null pk INT(11)"` //用户id
	LastCheckIn     int64 `xorm:"bigint not null"`     //上一次签到时间
	ContinuousCount int   `xorm:"not null INT(11)"`    //最后签到日期
}

type CheckInReward struct {
	DayIndex   int   `xorm:"not null pk INT(11)"` //从1开始
	RewardGold int64 `xorm:"bigint not null"`     //奖励的金币数
}

//------------------web
type ServerVersion struct {
	Version    string `xorm:"VARCHAR(50) pk not null"` //版本号
	IsMaintain int    `xorm:"TINYINT(1) not null"`     //是否维护,1为维护中,0为正常
	CreatedAt  int64  `xorm:"bigint not null"`
}

type WebUser struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"` //自增id
	LoginName string    `xorm:"VARCHAR(50) not null"`         //登录名
	Pwd       string    `xorm:"VARCHAR(50) not null"`         //密码
	RoleId    int       `xorm:"not null INT(11)"`             //权限id
	CreatedAt time.Time `xorm:"created"`
}

type Role struct {
	Id   int    `xorm:"not null pk INT(11)"`
	Desc string `xorm:"VARCHAR(32) not null"` //权限名称
}

type BlackListJump struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	BLJumpTo  string `xorm:"VARCHAR(255) not null"` //黑名单跳转
	PCJumpTo  string `xorm:"VARCHAR(255) not null"` //pc版跳转
	CreatedAt int64  `xorm:"bigint not null"`
}

type MemberLevelData struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"`
	Name     string `xorm:"VARCHAR(255) not null"`         //等级名称
	Level    int    `xorm:"not null INT(11)"`              //级别
	Price    int64  `xorm:"bigint not null"`               //购买价格
	IsHidden int    `xorm:"TINYINT(1) not null default 0"` //0为显示,1为隐藏
}

type AppCustomShare struct {
	Title     string `xorm:"VARCHAR(255) not null"` //标题
	Desc      string `xorm:"VARCHAR(255) not null"` //描述
	ImgName   string `xorm:"VARCHAR(255) not null"` //图片名称
	CreatedAt int64  `xorm:"bigint not null"`
}

type GZHCustomShare struct {
	Title     string `xorm:"VARCHAR(255) not null"` //标题
	Desc      string `xorm:"VARCHAR(255) not null"` //描述
	ImgName   string `xorm:"VARCHAR(255) not null"` //图片名称
	CreatedAt int64  `xorm:"bigint not null"`
}

type MemberLevelOrder struct {
	Id            int   `xorm:"not null pk autoincr INT(11)"`
	UserId        int   `xorm:"not null INT(11)"` //用户id
	MemberLevelId int   `xorm:"not null INT(11)"`
	CreatedAt     int64 `xorm:"bigint not null"`
}

type UserDepositInfo struct {
	Id          int         `xorm:"not null pk autoincr INT(11)"`
	UserId      int         `xorm:"not null INT(11)"` //用户id
	Money       float64     `xorm:"decimal(16,2) not null "`
	Platform    Platform    `xorm:"not null INT(11) "` //0为app,1为h5
	PayPlatform PayPlatform `xorm:"not null INT(11) "` //用户使用的付费平台,比如0(微信)
	CreatedAt   int64       `xorm:"bigint not null"`
}

// //今日用户活跃信息
// type TodayUserActivityInfo struct {
// 	Id        int   `xorm:"not null pk autoincr INT(11)"`
// 	UserId    int   `xorm:"not null INT(11)"` //用户id
// 	StartTime int64 `xorm:"bigint not null"`  //今天开始访问时间
// 	EndTime   int64 `xorm:"bigint not null"`  //今天结束访问时间
// }

type RushLimitSetting struct {
	Id           int    `xorm:"not null pk INT(11)"`
	LotteryCount int    `xorm:"not null INT(11)"`      //用户达到中奖次数后
	Diff2        int    `xorm:"not null INT(11)"`      //改变第二关难度系数
	Diff2r       int    `xorm:"not null INT(11)"`      //改变第二关口红数
	Diff2t       int    `xorm:"not null INT(11)"`      //改变第二关游戏时间
	Diff3        int    `xorm:"not null INT(11)"`      //改变第三关难度系数
	Diff3r       int    `xorm:"not null INT(11)"`      //改变第三关口红数
	Diff3t       int    `xorm:"not null INT(11)"`      //改变第三关游戏时间
	CheatCount   int    `xorm:"not null INT(11)"`      //达到作弊条件
	CheatTips    string `xorm:"VARCHAR(255) not null"` //提示语
}

type RandomLotteryGoods struct {
	Id           int    `xorm:"not null pk autoincr INT(11)"`                  //自增id
	Name         string `xorm:"VARCHAR(255) not null"`                         //商品名称
	ImgName      string `xorm:"VARCHAR(50) not null"`                          //图片名称
	Price        int64  `xorm:"bigint not null"`                               //价格
	GoodsDesc    string `xorm:"VARCHAR(255) not null"`                         //商品描述
	Brand        string `xorm:"VARCHAR(255) not null"`                         //所属品牌
	GoodsClassId int    `xorm:"null INT(11)"`                                  //商品类型id
	Probability  int    `xorm:"not null INT(11)"`                              //中奖概率
	IsHidden     int    `xorm:"TINYINT(1) not null default 0" json:"ishidden"` //0为显示,1为隐藏
	CreatedAt    int64  `xorm:"bigint not null"`                               //创建时间
}

//抽奖商品池水
type RandomLotteryGoodsPool struct {
	Id                 int     `xorm:"not null pk INT(11)"`
	Current            float64 `xorm:"decimal(16,2) not null "` //该商品在当前池中的金额
	RandomLotteryCount int     `xorm:"not null INT(11)"`        //用户达到中奖次数后
	Probability        int     `xorm:"not null INT(11)"`        //概率
}

//抽奖成功记录表
type RandomLotteryGoodsSucceed struct {
	OrderId        int64 `xorm:"bigint not pk null"` //商品订单号
	LotteryGoodsId int   `xorm:"not null INT(11)"`   //奖品id
	UserId         int   `xorm:"not null INT(11)"`
	CreatedAt      int64 `xorm:"bigint not null"`
}

//抽奖成功历史表
type RandomLotteryGoodsSucceedHistory struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"` //自增id
	NickName string `xorm:"VARCHAR(255) not null"`        //昵称
	Avatar   string `xorm:"VARCHAR(255) not null"`        //头像
	DescInfo string `xorm:"VARCHAR(255) not null"`        //用户中奖描述
}

//商品图片(如封面)
type GoodsImgs struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"` //自增id
	GoodsId  int    `xorm:"not null INT(11)"`             //对应商品id
	ImgName  string `xorm:"VARCHAR(255) not null"`        //名称
	ImgIndex int    `xorm:"not null INT(11)"`             //图片顺序,从0开始
}

//商品详情
type GoodsDetail struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"` //自增id
	GoodsId  int    `xorm:"not null INT(11)"`             //对应商品id
	ImgName  string `xorm:"VARCHAR(255) not null"`        //名称
	ImgIndex int    `xorm:"not null INT(11)"`             //图片顺序,从0开始
}

type UserAppraise struct {
	Id        int              `xorm:"not null pk autoincr INT(11)"` //自增id
	UserId    int              `xorm:"not null INT(11)"`
	Desc      string           `xorm:"VARCHAR(255) not null"` //文本描述 最大为100字
	ShowType  UserAppraiseType `xorm:"TINYINT(1) not null"`   //展示类型
	GoodsId   int              `xorm:"not null INT(11)"`
	CreatedAt int64            `xorm:"bigint not null"`
	IsPassed  int              `xorm:"TINYINT(1) not null"` //0未审核,1审核通过
	GoodsType GoodsDataType    `xorm:"TINYINT(1) not null"` //商品类型
}

type UserAppraisePic struct {
	Id             int    `xorm:"not null pk autoincr INT(11)"` //自增id
	UserAppraiseId int    `xorm:"not null INT(11)"`
	ImgName        string `xorm:"VARCHAR(255) not null"`
	ImgIndex       int    `xorm:"not null INT(11)"`
}

type UserShippingAddress struct {
	UserId  int    `xorm:"not null pk INT(11)"` //用户id
	Linkman string `xorm:"VARCHAR(255) not null"`
	Phone   string `xorm:"VARCHAR(25) not null"`
	Addr    string `xorm:"VARCHAR(255) not null"`
	Remark  string `xorm:"VARCHAR(255) null"`
}

type SuggestionComplaintParams struct {
	Id                   int `xorm:"not null pk INT(11)"` //自增id
	SuggestCountForDay   int `xorm:"not null INT(11)"`    //每天建议次数
	ComplaintCountForDay int `xorm:"not null INT(11)"`    //每天投诉次数
	ComplaintCountForBL  int `xorm:"not null INT(11)"`    //达到此投诉次数,自动进入黑名单
}

type Suggestion struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"` //自增id
	UserId    int    `xorm:"not null INT(11)"`             //用户id
	Desc      string `xorm:"VARCHAR(255) not null"`
	CreatedAt int64  `xorm:"bigint not null"`
}

type Complaint struct {
	Id            int    `xorm:"not null pk autoincr INT(11)"` //自增id
	UserId        int    `xorm:"not null INT(11)"`             //用户id
	ComplaintType string `xorm:"VARCHAR(100) not null"`
	Desc          string `xorm:"VARCHAR(255) not null"`
	CreatedAt     int64  `xorm:"bigint not null"`
}

type GoldCoinGift struct {
	Id                   int   `xorm:"not null pk INT(11)"`
	DownLoadAppGoldGift  int64 `xorm:"not null INT(11)"`    //下载app赠送金币数
	AppraisedGoldGift    int64 `xorm:"not null INT(11)"`    //评价后赠送金币数
	RegisterGoldGift     int64 `xorm:"not null INT(11)"`    //新人用户赠送金币数
	IsEnableRegisterGift int   `xorm:"TINYINT(1) not null"` //0关闭新人福利,1开启新人福利
	IsDrawCashOnlyApp    int   `xorm:"TINYINT(1) not null"` //0关闭 ,1开启。是否只在app内提现
}

type TmpData struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"` //自增id
	NickName string `xorm:"VARCHAR(100) not null"`
	Avatar   string `xorm:"VARCHAR(500) not null"`
}

type TmpDataForGoods struct {
	Id        int   `xorm:"not null pk autoincr INT(11)"` //自增id
	TmpUserId int   `xorm:"not null INT(11)"`
	GoodsId   int   `xorm:"not null INT(11)"`
	UpdateAt  int64 `xorm:"bigint not null"`
}
