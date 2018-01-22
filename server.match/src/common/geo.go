package common

type Coordinate struct {
	UID       uint64
	Latitude  float64
	Longitude float64
	Desc      string
}

//附近玩家
/*type GeoUser struct {
	Id           uint64 // 玩家id
	Account      string // 帐号
	Icon         uint32 // 图标
	Sex          uint8  // 性别
	State        uint32 // 状态
	PassIcon     string // 头像
	City         string
	Desc         string
	Longitude    float64
	Latitude     float64
	Level        uint32 // 当前段位
	Scores       uint32 // 当前星数
	Age          uint32
	Sign         string
	AudienceUrl  string
	AudienceTime uint32
	RelType      uint32
	YaoPlay      uint32
	Distance     float64
}*/

type GeoResult struct {
	Address string `json:"formatted_address"`
	//AddressComponent AddressComponent `json:"addressComponent"`
}

type AddressComponent struct {
	City          string `json:"city"`
	Country       string `json:"country"`
	Direction     string `json:"direction"`
	Distance      string `json:"distance"`
	District      string `json:"district"`
	Province      string `json:"province"`
	Street        string `json:"street"`
	Street_number string `json:"street_number"`
	Country_code  int    `json:"country_code"`
}

type AddressResult struct {
	Status int       `json:"status"`
	Result GeoResult `json:"result"`
}

type UserLocaltion struct {
	User       GeoUser
	Coordinate Coordinate
}

type CityGeo struct {
	City       string
	Coordinate Coordinate
	Pois       []GaodePoi
	Roads      []GaodeRoad
}

//高德接口相关的结构体定义
type GaodePoi struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type GaodeRoad struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type GaodeReGeoCode struct {
	FormattedAddress string      `json:"formatted_address"`
	Roads            []GaodeRoad `json:"roads"`
	Pois             []GaodePoi  `json:"pois"`
}

type GaodeGeo struct {
	Status     string           `json:"status"`
	Info       string           `json:"info"`
	InfoCode   string           `json:"infocode"`
	ReGeoCodes []GaodeReGeoCode `json:"regeocodes"`
}

//附近的自建房间
type NearRoom struct {
	UserId  uint64
	Account string // 账号
	Name    string // 房间名
	Sex     uint8  // 性别
	RType   uint32 // 房间类型 0自由模式 1组队模式 2闪电战模式
	Priv    uint32 // 房间权限 0只有邀请可以进入 1所有人可以进入 2输入密码可以进入
	MemNum  uint32 // 成员数量
	NowNum  uint32 // 当前成员数量
	City    string
	Coordinate Coordinate
	Time    uint32
}


/////////////////////////
