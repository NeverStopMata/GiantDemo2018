package conf

import (
	"base/env"
	"base/glog"
	"common"
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"usercmd"
)

const (
	ODTYPE_RANDONE = 0
	ODTYPE_RANDALL = 1
)

var (
	Config_Udp = false
)

// 队伍名
type XmlItem struct {
	Name uint32 `xml:"name,attr"`
}
type XmlTeamName struct {
	XMLName xml.Name  `xml:"config"`
	Items   []XmlItem `xml:"item"`
}

// 房间配置
type XmlRoomModel struct {
	PlayerNum uint32 `xml:"playernum,attr"`
	LastTm    uint32 `xml:"lasttime,attr"`
}

//刷新点
type XmlFoodPoint struct {
	ObjType int32                 `xml:"objtype,attr"`
	Foods   []XmlFoodPointRefresh `xml:"food"`
}

type XmlFoodPointRefresh struct {
	ID       uint16  `xml:"id,attr"`
	Num      uint32  `xml:"num,attr"`
	Interval float64 `xml:"interval,attr"`
}

//ai相关
type XmlAIInitDataLevelAni struct {
	Id uint32 `xml:"id,attr"`
}

type XmlAIInitDataLevel struct {
	Id      uint32                  `xml:"id,attr"`
	Animals []XmlAIInitDataLevelAni `xml:"animal"`
}

type XmlAIInitData struct {
	Levels []XmlAIInitDataLevel `xml:"level"`
}

//XmlAIBehave
type XmlAIBehaveData struct {
	Id     uint32 `xml:"id,attr"`
	Aifile string `xml:"aifile,attr"`

	AttackRange float64 `xml:"attackRange,attr"`
	FollowTime  int64   `xml:"followTime,attr"`
	FollowRange float64 `xml:"followRange,attr"`
	EvoTime     float64 `xml:"evoTime,attr"`
}

type XmlAIBehave struct {
	Datas []XmlAIBehaveData `xml:"bev"`
}

//XmlAISc2data
type XmlAISc2dataLevel struct {
	Id        uint32  `xml:"id,attr"`
	MinScore  uint32  `xml:"min,attr"`
	MaxScore  uint32  `xml:"max,attr"`
	Initdata  uint32  `xml:"initdata,attr"`
	Behave    string  `xml:"behave,attr"`
	SpeedRate float64 `xml:"speedrate,attr"`

	Power string `xml:"power,attr"`

	//计算数据
	Behaves  []uint32
	Powers   []int
	PowerSum int
}

type XmlAISc2data struct {
	Levels []*XmlAISc2dataLevel `xml:"data"`
}

type XmlAI struct {
	MaxExp   uint32        `xml:"maxExp,attr"`
	Initdata XmlAIInitData `xml:"initdata"`
	Behave   XmlAIBehave   `xml:"behave"`
	Sc2data  XmlAISc2data  `xml:"sc2data"`
}

//ai相关end

type XmlTreePosInfo struct {
	Id uint16  `xml:"id,attr"`
	X  float32 `xml:"x,attr"`
	Y  float32 `xml:"y,attr"`
}

type XmlTreePos struct {
	Pos []XmlTreePosInfo `xml:"pos"`
}

type XmlReLevel struct {
	Id   uint16 `xml:"id,attr"`
	Food uint16 `xml:"foodid,attr"`
}

type XmlTreeRe struct {
	BeTime uint16       `xml:"begintime,attr"`
	Inter  uint16       `xml:"refreshtime,attr"`
	Num    uint16       `xml:"refreshnum,attr"`
	ReMind uint16       `xml:"remindtime,attr"`
	Level  []XmlReLevel `xml:"level"`
}

type XmlTreeCfg struct {
	Level XmlTreeRe  `xml:"tree"`
	Pos   XmlTreePos `xml:"position"`
}

type XmlGlobal struct {
	XMLName  xml.Name     `xml:"global"`
	QRoom    XmlRoomModel `xml:"troom"`
	TeamRoom XmlRoomModel `xml:"teamroom"`
	AllMap   int          `xml:"allmap"`
	Pystress uint32       `xml:"pystress"`
	Norobot  uint32       `xml:"norobot"`
}

type XmlGameSet struct {
	WormholeWater int32 `xml:"wormholeWater,attr"` //虫洞耗水
}

type XmlAITable struct {
	XMLName xml.Name `xml:"config"`
	AIData  *XmlAI   `xml:"ai"` //无限
}

type XmlRoom struct {
	Id uint32 `xml:"id,attr"`
}

type XmlDeadFoods struct {
	DropPro   int32         `xml:"droppro,attr"`
	KillPro   int32         `xml:"killpro,attr"`
	DropNum   int32         `xml:"dropnum,attr"`
	DeadPro   int32         `xml:"deadpro,attr"`
	DeadFoods []XmlDeadFood `xml:"food"`
}

type XmlDeadFood struct {
	Index  uint32 `xml:"id,attr"`
	FoodId uint16 `xml:"foodid,attr"`
	Min    uint32 `xml:"min,attr"`
	Max    uint32 `xml:"max,attr"`
	Type   uint32 `xml:"modetype,attr"`
}

type XmlMap struct {
	XMLName xml.Name  `xml:"config"`
	Scenes  []XmlRoom `xml:"scene"`
}

// 机器人配置
type XmlRobotLevel struct {
	Val    uint32 `xml:"val,attr"`
	Robot  uint32 `xml:"robot,attr"`
	Player uint32 `xml:"player,attr"`
	Good   uint32 `xml:"good,attr"`
	Bad    uint32 `xml:"bad,attr"`
	Up     uint32 `xml:"up,attr"`
	Down   uint32 `xml:"down,attr"`
}

type XmlRobot struct {
	XMLName xml.Name                  `xml:"Robot"`
	Default uint32                    `xml:"default,attr"`
	Old     uint32                    `xml:"old,attr"`
	Max     uint32                    `xml:"max,attr"`
	Normal  uint32                    `xml:"normal,attr"`
	StartId uint64                    `xml:"startId,attr"`
	L       []XmlRobotLevel           `xml:"level"`
	Levels  map[uint32]*XmlRobotLevel `xml:"-"`
}

type XmlRobotNameItem struct {
	Names string `xml:"NameKey,attr"`
}

type XmlForeign struct {
	Countries    string   `xml:"countries,attr"`
	NameDefault  uint32   `xml:"namedefault,attr"`
	CountryCodes []uint32 `xml:"-"`
}

type XmlCountryName struct {
	Country uint32             `xml:"country,attr"`
	Items   []XmlRobotNameItem `xml:"item"`
}

type XmlRobotName struct {
	XMLName      xml.Name            `xml:"Config"`
	Foreign      XmlForeign          `xml:"foreign"`
	CountryNames []XmlCountryName    `xml:"rname"`
	Names        map[uint32][]string `xml:"-"`
}

type XmlRobotPlayer struct {
	Id       uint64 `xml:"id,attr"`
	Acc      string `xml:"acc,attr"`
	Sex      uint8  `xml:"sex,attr"`
	Icon     uint32 `xml:"icon,attr"`
	PassIcon string `xml:"passIcon,attr"`
	Location uint32 `xml:"location,attr"`
}

type XmlRobotData struct {
	XMLName xml.Name         `xml:"data"`
	Players []XmlRobotPlayer `xml:"player"`
}

type XmlLevelItem struct {
	Id  uint16 `xml:"id,attr"`
	Exp uint32 `xml:"nextexp,attr"`
}

type XmlLevelItems struct {
	MapId uint32         `xml:"mapid,attr"`
	Items []XmlLevelItem `xml:"item"`
}
type XmlLevelCfg struct {
	XMLName xml.Name        `xml:"config"`
	Levels  []XmlLevelItems `xml:"level"`
}

//动物配置表
type XmlAnimalFood struct {
	FoodId uint16 `xml:"foodid,attr"`
	Exp    uint32 `xml:"exp,attr"`
}

type XmlAnimalSpeed struct {
	LandId uint16  `xml:"landid,attr"`
	Speed  float64 `xml:"speed,attr"`
}

type XmlAnimalAngular struct {
	LandId uint16  `xml:"landid,attr"`
	Speed  float64 `xml:"speed,attr"`
}

type XmlAnimalFoods struct {
	Foods []XmlAnimalFood `xml:"item"`
}

type XmlAnimalWalkBlocks struct {
	WalkBlocks []XmlWalkBlock `xml:"item"`
}

type XmlWalkBlock struct {
	Type uint16 `xml:"type,attr"`
}

type XmlAnimalSpeeds struct {
	Speeds []XmlAnimalSpeed `xml:"item"`
}
type XmlAnimalAngulars struct {
	Angulars []XmlAnimalAngular `xml:"item"`
}
type XmlAnimaConsume struct {
	Landid    uint16 `xml:"landid,attr"`
	Persecond uint16 `xml:"persecond,attr"`
}

type XmlAnimaConsumeItems struct {
	Items []XmlAnimaConsume `xml:"item"`
}

type XmlAnimalRestore struct {
	Landid    uint16 `xml:"landid,attr"`
	Persecond uint16 `xml:"persecond,attr"`
}

type XmlAnimalRestoreItems struct {
	Items []XmlAnimalRestore `xml:"item"`
}

type XmlAnima struct {
	Id    uint16  `xml:"id,attr"`
	Scale float64 `xml:"scale,attr"`
}

type XmlAnimaCfg struct {
	XMLName xml.Name   `xml:"config"`
	Animals []XmlAnima `xml:"animal"`
}

type XmlWaterItem struct {
	Landid    uint16 `xml:"landid,attr"`
	Persecond uint16 `xml:"persecond,attr"`
}

type XmlWaterItems struct {
	Items []XmlWaterItem `xml:"item"`
}

type XmlWaterRestoreItems struct {
	Items []XmlWaterItem `xml:"item"`
}

type XmlFoodItem struct {
	FoodId    uint16                `xml:"id,attr"`
	FoodType  uint16                `xml:"type,attr"`
	Size      float32               `xml:"size,attr"`
	BirthTime float64               `xml:"birthTime,attr"`
	MapNum    uint16                `xml:"mapnum,attr"`
	HP        uint32                `xml:"hp,attr"`
	Buffstr   string                `xml:"buff,attr"`
	Area      int                   `xml:"area,attr"`
	Exp       uint32                `xml:"exp,attr"`
	LiveTime  int64                 `xml:"time,attr"`
	Children  []XmlFoodPointRefresh `xml:"child"`
	Buff      []uint32
	Rate      []uint16 // 概率
	Sum       uint16   // 基数
}

type XmlFoodItems struct {
	MapId  uint32 `xml:"mapid,attr"`
	ShotID uint16 `xml:"shotid,attr"`
	//ShotExp   uint32        `xml:"shotexp,attr"`
	//ShotDis   float64       `xml:"shotdis,attr"`
	//ShotSpeed float64       `xml:"shotspeed,attr"`
	//ShotFrame uint32        `xml:"shotframe,attr"`
	Items []XmlFoodItem `xml:"item"`
}

type XmlFoodCfg struct {
	XMLName xml.Name       `xml:"config"`
	Foods   []XmlFoodItems `xml:"food"`
	//QFood   XmlFoodItems `xml:"tfood"`
}

type XmlFRankItem struct {
	Id         uint16 `xml:"id,attr"`
	ScoreStart uint64 `xml:"scorestart,attr"`
	ScoreEnd   uint64 `xml:"scoreend,attr"`
	ItemId     uint32 `xml:"itemid,attr"`
	ItemNum    uint32 `xml:"num,attr"`
	AddHideExp uint32 `xml:"add_hide_exp,attr"`
	DropGroup  uint32 `xml:"dropgroup,attr"`
}

type XmlFRankItems struct {
	Items []XmlFRankItem `xml:"item"`
}

type XmlTRankItem struct {
	Id         uint16 `xml:"id,attr"`
	RankStart  uint32 `xml:"rankstart,attr"`
	RankeEnd   uint32 `xml:"rankend,attr"`
	ItemId     uint32 `xml:"itemid,attr"`
	ItemNum    uint32 `xml:"num,attr"`
	AddHideExp uint32 `xml:"add_hide_exp,attr"`
	DropGroup  uint32 `xml:"dropgroup,attr"`
}

type XmlTRankItems struct {
	Items []XmlTRankItem `xml:"item"`
}

type XmlRankCfg struct {
	XMLName   xml.Name      `xml:"config"`
	FRank     XmlFRankItems `xml:"frank"`
	TRank     XmlTRankItems `xml:"trank"`
	TeamRank  XmlTRankItems `xml:"teamrank"`
	WeekFRank XmlTRankItems `xml:"weekfrank"`
	WeekTRank XmlTRankItems `xml:"weekrtrank"`
}

//////////////////////////////////////////////////////////////////////////////////
//配置类
type ConfigMgr struct {
	Global      *XmlGlobal
	Map         *XmlMap
	RobotNams   *XmlRobotName
	RobotDatas  *XmlRobotData
	LevelDatas  *XmlLevelCfg
	AnimalDatas *XmlAnimaCfg
	FoodDatas   *XmlFoodCfg
	RankDatas   *XmlRankCfg
	AIDatas     *XmlAITable
	v_tnames    []uint32
}

var (
	configm      *ConfigMgr
	configmMutex sync.RWMutex
)

func NewConfigMgr() *ConfigMgr {
	c := &ConfigMgr{
		v_tnames: make([]uint32, 0),
	}
	return c
}

func ConfigMgr_GetMe() (c *ConfigMgr) {
	if configm == nil {
		configm = NewConfigMgr()
	}
	configmMutex.RLock()
	c = configm
	configmMutex.RUnlock()
	return
}

func ReloadConfig() bool {
	c := NewConfigMgr()
	if !c.Init() {
		return false
	}
	configmMutex.Lock()
	configm = c
	configmMutex.Unlock()
	return true
}

// 全局配置
func (this *ConfigMgr) LoadGlobal() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "global.xml")
	if err != nil {
		glog.Error("[配置] 打开配置失败 ", err)
		return false
	}
	xmlGlobal := &XmlGlobal{}
	err = xml.Unmarshal(content, xmlGlobal)
	if err != nil {
		glog.Error("[配置] 解析配置失败 ", err)
		return false
	}
	this.Global = xmlGlobal
	return true
}

// 全局ai
func (this *ConfigMgr) LoadAI() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "ai.xml")
	if err != nil {
		glog.Error("[配置] 打开ai配置失败 ", err)
		return false
	}
	xmlmap := &XmlAITable{}
	err = xml.Unmarshal(content, xmlmap)
	if err != nil {
		glog.Error("[配置] 解析ai配置失败 ", err)
		return false
	}
	this.AIDatas = xmlmap
	this.InitAI(this.AIDatas.AIData)
	return true
}

func (this *ConfigMgr) InitAI(aidata *XmlAI) {
	if aidata == nil {
		glog.Error("[配置] 解析ai配置失败  -  InitAI")
		return
	}
	for _, data := range aidata.Sc2data.Levels {
		//计算bev
		bevs := strings.Split(data.Behave, ",")
		for _, bs := range bevs {
			v, _ := strconv.Atoi(bs)
			data.Behaves = append(data.Behaves, uint32(v))
		}

		//计算power
		powers := strings.Split(data.Power, ",")
		var sum = 0
		for _, bs := range powers {
			v, _ := strconv.Atoi(bs)
			sum += v
			data.Powers = append(data.Powers, sum)
		}
		data.PowerSum = 100
		if sum > 0 {
			data.PowerSum = sum
		}

		glog.Info("InitAI:", data.Id, "  ", data.Behaves, "  ", data.Powers, "  sum:", data.PowerSum)
	}
}

// 全局map
func (this *ConfigMgr) LoadMap() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "map.xml")
	if err != nil {
		glog.Error("[配置] 打开map配置失败 ", err)
		return false
	}
	xmlmap := &XmlMap{}
	err = xml.Unmarshal(content, xmlmap)
	if err != nil {
		glog.Error("[配置] 解析map配置失败 ", err)
		return false
	}
	this.Map = xmlmap
	glog.Info("LoadMap:", len(this.Map.Scenes))
	return true
}

func (this *ConfigMgr) LoadRobotName() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "robotname.xml")
	if err != nil {
		glog.Error("[配置] 打开配置robot Name失败 ", err)
		return false
	}
	XmlRobotName := &XmlRobotName{}
	err = xml.Unmarshal(content, XmlRobotName)
	if err != nil {
		glog.Error("[配置] 解析配置robot Name失败 ", err)
		return false
	}

	countries := strings.Split(XmlRobotName.Foreign.Countries, "|")
	for _, v := range countries {
		codeid, _ := strconv.Atoi(v)
		XmlRobotName.Foreign.CountryCodes = append(XmlRobotName.Foreign.CountryCodes, uint32(codeid))
	}

	glog.Info("[配置] 解析机器人国家列表: ", XmlRobotName.Foreign.CountryCodes)

	XmlRobotName.Names = make(map[uint32][]string)
	for _, cns := range XmlRobotName.CountryNames {
		country := cns.Country
		for _, item := range cns.Items {
			names := strings.Split(item.Names, "|")
			XmlRobotName.Names[country] = append(XmlRobotName.Names[country], names...)
		}
	}

	glog.Info("[配置] 解析机器人名字列表: ", len(XmlRobotName.Names))

	this.RobotNams = XmlRobotName
	return true
}

func (this *ConfigMgr) LoadRobotData() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "robotdata.xml")
	if err != nil {
		glog.Error("[配置] 打开配置robot Data失败 ", err)
		return false
	}
	robotData := &XmlRobotData{}
	err = xml.Unmarshal(content, robotData)
	if err != nil {
		glog.Error("[配置] 解析配置robot Data失败 ", err)
		return false
	}

	this.RobotDatas = robotData
	return true
}

func (this *ConfigMgr) LoadLevelCfg() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "level.xml")
	if err != nil {
		glog.Error("[配置] 打开配置level.xml失败 ", err)
		return false
	}
	levelData := &XmlLevelCfg{}
	err = xml.Unmarshal(content, levelData)
	if err != nil {
		glog.Error("[配置] 解析配置 level.xml 失败", err)
		return false
	}
	this.LevelDatas = levelData
	return true
}

func (this *ConfigMgr) LoadAnimalCfg() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "animal.xml")
	if err != nil {
		glog.Error("[配置] 打开配置animal.xml失败 ", err)
		return false
	}
	animalData := &XmlAnimaCfg{}
	err = xml.Unmarshal(content, animalData)
	if err != nil {
		glog.Error("[配置] 解析配置 animal.xml 失败", err)
		return false
	}
	this.AnimalDatas = animalData
	return true
}

func (this *ConfigMgr) LoadFoodCfg() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "food.xml")
	if err != nil {
		glog.Error("[配置] 打开配置 food.xml失败 ", err)
		return false
	}
	foodData := &XmlFoodCfg{}
	err = xml.Unmarshal(content, foodData)
	if err != nil {
		glog.Error("[配置] 解析配置 food.xml 失败", err)
		return false
	}
	for index, food := range foodData.Foods {
		for innr, item := range food.Items {
			if len(item.Buffstr) == 0 {
				continue
			}
			strs := strings.Split(item.Buffstr, "|")
			for _, v := range strs {
				tmp := strings.Split(v, ":")
				if len(tmp) != 2 {
					glog.Error("food.xml bufferandrate error", v)
					return false
				}
				buff, ok := strconv.Atoi(tmp[0])
				if nil != ok {
					glog.Error("food.xml buff error: ", strs, ",", item.Buffstr)
					return false
				}
				rate, ok := strconv.Atoi(tmp[1])
				foodData.Foods[index].Items[innr].Buff = append(foodData.Foods[index].Items[innr].Buff, uint32(buff))
				foodData.Foods[index].Items[innr].Rate = append(foodData.Foods[index].Items[innr].Rate, uint16(rate))
				foodData.Foods[index].Items[innr].Sum += uint16(rate)
			}
			if 0 == foodData.Foods[index].Items[innr].Sum ||
				len(foodData.Foods[index].Items[innr].Buff) != len(foodData.Foods[index].Items[innr].Rate) {
				glog.Error("food.xml sum is zero or len error", strs)
				return false
			}
		}
	}
	this.FoodDatas = foodData
	//	for _, v := range this.FoodDatas.Foods {
	//		for _, v2 := range v.Items {
	//			glog.Info("LoadFoodCfg ", v2.Buff, ",", v2.Buffstr, ",", v2.FoodId, ",", v.MapId, ",", v2.Rate, ",", v2.Sum)
	//		}
	//	}
	return true
}

func (this *ConfigMgr) LoadRankCfg() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "rank.xml")
	if err != nil {
		glog.Error("[配置] 打开配置 rank.xml失败 ", err)
		return false
	}
	rankData := &XmlRankCfg{}
	err = xml.Unmarshal(content, rankData)
	if err != nil {
		glog.Error("[配置] 解析配置 rank.xml 失败", err)
		return false
	}
	this.RankDatas = rankData
	for _, data := range this.RankDatas.FRank.Items {
		glog.Info("Frank:", data.Id, ",", data.ScoreStart, ",", data.ScoreEnd, ",", data.ItemId, ",", data.ItemNum)
	}
	for _, data1 := range this.RankDatas.TRank.Items {
		glog.Info("Trank:", data1.Id, ",", data1.RankStart, ",", data1.RankeEnd, ",", data1.ItemId, ",", data1.ItemNum)
	}
	for _, data2 := range this.RankDatas.WeekFRank.Items {
		glog.Info("WFrank:", data2.Id, ",", data2.RankStart, ",", data2.RankeEnd, ",", data2.ItemId, ",", data2.ItemNum)
	}
	for _, data3 := range this.RankDatas.WeekTRank.Items {
		glog.Info("WTrank:", data3.Id, ",", data3.RankStart, ",", data3.RankeEnd, ",", data3.ItemId, ",", data3.ItemNum)
	}
	for _, data4 := range this.RankDatas.TeamRank.Items {
		glog.Info("TeamRank:", data4.Id, ",", data4.RankStart, ",", data4.RankeEnd, ",", data4.ItemId, ",", data4.ItemNum)
	}
	return true
}

func (this *ConfigMgr) loadTeamNameConfig() bool {
	content, err := ioutil.ReadFile(env.Get("global", "xmlcfg") + "teamname.xml")
	if err != nil {
		glog.Error("[配置] 打开配置失败 ", err)
		return false
	}
	teamname := &XmlTeamName{}
	err = xml.Unmarshal(content, teamname)
	if err != nil {
		glog.Error("[配置] 解析配置失败 ", err)
		return false
	}
	for _, v := range teamname.Items {
		this.v_tnames = append(this.v_tnames, v.Name)
	}
	glog.Info("[加载队伍名字配置]", this.v_tnames)
	return true
}

func (this *ConfigMgr) Init() bool {
	ok := this.LoadGlobal()
	if !ok {
		return false
	}

	if !this.LoadMap() {
		return false
	}

	if !this.LoadRobotName() {
		return false
	}

	if !this.LoadRobotData() {
		return false
	}
	if !this.LoadLevelCfg() {
		return false
	}
	if !this.LoadAnimalCfg() {
		return false
	}
	if !this.LoadFoodCfg() {
		return false
	}
	if !this.LoadRankCfg() {
		return false
	}

	if !this.LoadAI() {
		return false
	}
	if !this.loadTeamNameConfig() {
		return false
	}

	// DEL
	//	if !this.LoadSurviveConfig() {
	//		return false
	//	}
	glog.Info("[配置] 加载配置成功 ")
	return true
}

var defaultnames = []string{"西红柿蛋汤", "放开那个女孩", "殇不患", "新古龙群侠传", "0.0", "风吹乱了我的发型", "Jack", "Rose", "David",
	"黑猫警长", "炫斗三国志", "叫个鸭子", "仙侠世界", "大主宰", "武极天下", "新征途口袋", "奥斯卡", "骷髅", "专治各种不服", "小怪兽爱上奥特曼",
	"拿着试卷唱忐忑", "对着作业唱算你狠", "年少如歌", "手捧阳光", "哽住了喉", "明天你好", "爱妃接旨", "瘋言瘋语", "那時年少", "叽里呱啦", "迩芣慬", "莂說",
	"泡沫之夏", "醉相思", "凉夏", "浅黛梨妆", "木槿花開", "木槿暖夏", "乱了心", "萧萧暮雨", "︶浅笑ζ嫣然", "静若安然", "水墨青杉", "青纱挽妆", "恋雪", "空白",
	"烟花巷陌", "繁花、梦影", "心凉怎暖", "是梦终空", "幽竹烟雨", "若有来生", "想和你闹", "念念不忘", "浅暮流殇", "任心荒芜", "淺淺笑", "柠夏初开", "半盏流年", "訫洳訨氺",
	"樱花飞", "多啦ā梦", "浅夏〆淡殇", "ΩωΩ喵喵", "ロ觜角よ揚", "心悦君兮", "清风挽心", "陪你听风", "墨城烟柳", "且聼凨吟", "听弦断", "心之所向便是光", "在你的世界我称霸",
	"相見不如懷念", "骑着蜗牛找妞", "冷言冷语冷坚强", "轉角ジ撞到蘠", "那年╮沵笑靥如花", "长得丑活得久", "心已碎、情已断", "含个奶嘴闯天下", "城已空,人已散", "你这磨人的小妖精",
	"我对", "小賣部坑我們dě錢", "正宗纯天然学渣ァ", "时光催人老", "作业是个不可数名词", "蒙牛没我牛", "爱丶如履薄冰", "作业虐我千百遍T_T", "祝作业入土为安", "指着心脏说不痛",
	"承诺碎了一地", "心之所向便是光", "在你的世界我称霸", "相見不如懷念", "骑着蜗牛找妞", "冷言冷语冷坚强", "轉角ジ撞到蘠", "那年╮沵笑靥如花", "长得丑活得久", "心已碎、情已断",
	"含个奶嘴闯天下", "城已空, 人已散", "你这磨人的小妖精", "我对爱情过敏", "皇上，这是喜脉啊", "ㄡ冇誰會吢疼√", "獨愛伱一個ヤ", "不过一场少年梦", "别看了你帅不过我",
	"谁动朕江山", "回不去的年少时光", "奈何桥上唱小苹果", "人生苦短必须性感", "泡八喝九说十话", "路还长别猖狂", "没有钻石的王老五", "施主、你的贞节掉了", "星期⑧娶你",
	"你是我的生死劫", "别低头除非地上有钱", "╰朕赦迩无罪╯", "月亮是我踹弯的", "拿回忆下酒", "别闹快吃药", "虎背熊腰小蒸包", "露了馅的逗包", "拯救地球好累", "番茄你个西红柿",
	"喜欢天黑却怕鬼i", "国民男神经", "纯天然野生帅哥", "一岁时就很帅", "不帅你报警", "萌你一脸血", "闹钟你别闹", "超人不会飞", "掉毛的天使", "孟婆，来瓶加多宝",
	"闪耀的电灯泡", "剩下的盛夏", "哥不帅但很坏", "心碎谁买单", "数学这货太傲娇", "你与氧气平起平坐", "╭丅輩孒续約つ", "画个句号给昨天", "择一城终老", "咿呀咿呀哟~",
	"蛋蛋的忧伤", "奥特曼的蛋", "下个路口，放狗咬你", "1.2.3.木头人", "你在我心里迷了路", "花样作死冠军", "喂，放开那帅比", "吃素的蚊子", "开启作死模式", "奇葩哚哚″向阳开",
	"哎呀，你别闹", "一回到家就变乖", "总有佞臣想篡位", "没对象、省流量", "请叫我倍儿坚强", "疯人院vip用户", "趴下ゝ打劫棒棒糖", "若.只如初見", "时光偷走初心",
	"娇气的小奶包〆", "对着月亮说晚安", ">时间煮雨我煮饭ペ", "流年渲染了谁的容颜", "打听幸福的下落", "被风吹散の约定", "散场的拥抱", "丅①站-→幸福", "海水是鱼的眼泪√",
	"薄荷糖べ微微凉", "贪恋゛迩的温柔≈", "半城樱花﹌半城雨", "初夏ぃ蔷薇花开", "蒲公英dě约定ァ", "说好不沋伤〆", "微笑是莪的保护sè", "最美不过初遇见", "半夏琉璃ソ空人心",
	"１点点Dē小任性", "微微一笑醉倾城", "自找的痛ぃ何必喊疼", "偏偏喜欢迩、", "一笑倾城ぃ二笑倾心", "說好ㄋ吥見面", "花未落丶心已亡", "一个人的坚强", "半颗糖、也甜入心",
	"╭草编の戒指ㄨ", "戒不掉尓の味道ゝ", "谁能给我一世安稳", "尐哭尐闹尐任性", "嘴上逞强丶心里落泪", "吃醋是最诚实的告白", "爱很短ゝ回忆很长", "蝶舞﹌櫻婲落"}

func (this *ConfigMgr) GetRobotRandName(country uint32) string {
	countries, ok := this.RobotNams.Names[country]
	if !ok {
		countries, ok = this.RobotNams.Names[this.RobotNams.Foreign.NameDefault]
		if !ok {
			return defaultnames[rand.Intn(len(defaultnames))]
		}
	}

	return countries[rand.Intn(len(countries))]
}

func (this *ConfigMgr) RandForeignCountry() uint32 {
	return this.RobotNams.Foreign.CountryCodes[rand.Intn(len(this.RobotNams.Foreign.CountryCodes))]
}

/**
 * 获取随机队伍名
 * ids 已经被占用的队伍名字
 */
func (this *ConfigMgr) GetTeamName(ids map[uint32]bool) uint32 {
	nsize := len(this.v_tnames)
	if nsize == 0 {
		return 0
	}
	for i := 0; i < len(this.v_tnames); i++ {
		index := uint32(common.RandBetween(int64(i), int64(nsize-1)))
		if this.v_tnames[i] != this.v_tnames[index] {
			this.v_tnames[i], this.v_tnames[index] = this.v_tnames[index], this.v_tnames[i]
		}
	}
	var tname uint32 = 0
	for i := 0; i < len(this.v_tnames); i++ {
		tname = this.v_tnames[i]
		if _, ok := ids[tname]; !ok {
			break
		}
	}
	if tname == 0 {
		glog.Error("[ConfigMgr] GetTeamName tname=", tname, ",", ids)
	}
	return tname
}

func (this *ConfigMgr) GetRoomById(sceneId uint32) *XmlRoom {
	for _, m := range this.Map.Scenes {
		if m.Id == sceneId {
			return &m
		}
	}
	/*
		if roomId == common.RoomTypeFree {
			return &this.Global.Room
		} else if roomId == common.RoomTypeQuick {
			return &this.Global.QRoom
		}
	*/
	return nil
}

func (this *ConfigMgr) GetRoomType(rtype int) *XmlRoomModel {
	if rtype == common.RoomTypeQuick {
		return &this.Global.QRoom
	} else if rtype == common.RoomTypeTeam {
		return &this.Global.TeamRoom
	}

	return nil
}
func (this *ConfigMgr) GetAnimalSize(animalid uint16) (float64, bool) {
	data := this.GetAnimal(animalid)
	if data != nil {
		return data.Scale, true
	}
	return 0, false
}
func (this *ConfigMgr) GetAnimalExp(sceneId uint32, animalid uint16) uint32 {
	for _, val := range this.LevelDatas.Levels {
		if val.MapId == sceneId {
			for _, it := range val.Items {
				if it.Id == animalid {
					return it.Exp
				}
			}
		}
	}
	return 0
}

func (this *ConfigMgr) GetAnimal(animalid uint16) *XmlAnima {
	for index, _ := range this.AnimalDatas.Animals {
		if this.AnimalDatas.Animals[index].Id == animalid {
			return &this.AnimalDatas.Animals[index]
		}
	}
	return nil
}

func (this *ConfigMgr) GetLevelData(sceneId uint32) *XmlLevelItems {
	for _, v := range this.LevelDatas.Levels {
		if v.MapId == sceneId {
			return &v
		}
	}
	return nil
}

func (this *ConfigMgr) GetNextExp(sceneId uint32, id uint16) (uint32, bool) {
	v := this.GetLevelData(sceneId)

	for index, _ := range v.Items {
		if v.Items[index].Id == id {
			return v.Items[index].Exp, true
		}
	}

	return 0, false
}

func (this *ConfigMgr) GetMinAniLevel(sceneId uint32) uint16 {
	v := this.GetLevelData(sceneId)
	item := v.Items
	return item[0].Id
}

func (this *ConfigMgr) GetAnimalIdByExp(sceneId uint32, exp uint32) uint16 {
	v := this.GetLevelData(sceneId)
	var id uint16
	item := v.Items
	count := len(v.Items)

	if exp <= 0 {
		return item[0].Id
	}

	if exp >= item[count-1].Exp {
		return item[count-1].Id
	}

	for index, _ := range item {
		if exp < item[index].Exp {
			break
		}
		id = item[index].Id
	}

	if id == 0 {
		return item[0].Id
	}
	return id + 1
}

func (this *ConfigMgr) GetXmlFoodItems(sceneId uint32) *XmlFoodItems {
	for _, m := range this.FoodDatas.Foods {
		if m.MapId == sceneId {
			return &m
		}
	}
	return nil
}

//food 表
func (this *ConfigMgr) GetFood(t uint32, foodid uint16) *XmlFoodItem {
	item := this.GetXmlFoodItems(t)
	for _, val := range item.Items {
		if val.FoodId == foodid {
			return &val
		}
	}
	return nil
}

func (this *ConfigMgr) GetBuff(t uint32, foodid uint16) uint32 {
	val := this.GetFood(t, foodid)
	if nil == val || 0 == val.Sum || 0 == len(val.Buff) || 0 == len(val.Rate) {
		return 0
	}
	rate := common.RandBetweenInt32(0, int32(val.Sum))
	sum := uint16(0)
	for index := 0; index < len(val.Rate); index++ {
		sum += uint16(val.Rate[index])
		if sum >= uint16(rate) {
			return val.Buff[index]
		}
	}
	return 0
}

func (this *ConfigMgr) GetFoodSize(t uint32, typeId uint16) float32 {
	if food := this.GetFood(t, typeId); food != nil {
		return food.Size
	}
	return 0.5
}

func (this *ConfigMgr) GetFoodMapNum(t uint32, typeId uint16) uint16 {
	if food := this.GetFood(t, typeId); food != nil {
		return food.MapNum
	}
	return 0
}

func (this *ConfigMgr) GetFoodHP(t uint32, typeId uint16) uint32 {
	if food := this.GetFood(t, typeId); food != nil {
		return food.HP
	}
	return 0
}

func (this *ConfigMgr) GetFoodExp(t uint32, typeId uint16) uint32 {
	if food := this.GetFood(t, typeId); food != nil {
		return food.Exp
	}
	return 0
}

func (this *ConfigMgr) GetFoodTime(t uint32, typeId uint16) int64 {
	if food := this.GetFood(t, typeId); food != nil {
		return food.LiveTime
	}
	return 0
}

func (this *ConfigMgr) GetFoodBallType(t uint32, typeId uint16) usercmd.BallType {
	if food := this.GetFood(t, typeId); food != nil {
		return usercmd.BallType(food.FoodType)
	}
	return 0
}

func (this *ConfigMgr) GetRankByScore(t uint32, score uint64) (uint32, uint32, uint32) {
	return 0, 0, 0
}
func (this *ConfigMgr) GetDropGroupIdByScore(t uint32, score uint64) uint32 {
	return 0
}

func (this *ConfigMgr) GetDropGroupIdByRank(t, rank uint32) uint32 {
	if common.RoomTypeQuick == t {
		for index, _ := range this.RankDatas.TRank.Items {
			if rank >= this.RankDatas.TRank.Items[index].RankStart && rank <= this.RankDatas.TRank.Items[index].RankeEnd {
				return this.RankDatas.TRank.Items[index].DropGroup
			}
		}
	}
	if common.RoomTypeTeam == t {
		for index, _ := range this.RankDatas.TeamRank.Items {
			if rank >= this.RankDatas.TeamRank.Items[index].RankStart && rank <= this.RankDatas.TeamRank.Items[index].RankeEnd {
				return this.RankDatas.TeamRank.Items[index].DropGroup
			}
		}
	}
	return 0
}

func (this *ConfigMgr) GetRoomAI() *XmlAI {
	return this.AIDatas.AIData
}

func (this *ConfigMgr) GetAiInitDataLevel(id uint32, aidata *XmlAI) *XmlAIInitDataLevel {
	//glog.Info("GetAiInitDataLevel:", len(this.Global.AIData.Initdata.Levels))
	for _, data := range aidata.Initdata.Levels {
		if data.Id == id {
			return &data
		}
	}
	return nil
}

func (this *ConfigMgr) GetAiBevData(id uint32, aidata *XmlAI) *XmlAIBehaveData {
	//glog.Info("GetAiBevData:", len(this.Global.AIData.Behave.Datas))
	for _, data := range aidata.Behave.Datas {
		if data.Id == id {
			return &data
		}
	}
	return nil
}

func (this *ConfigMgr) GetAiBevDataByRoomType(bevid uint32) *XmlAIBehaveData {
	aidata := this.GetRoomAI()
	for _, data := range aidata.Behave.Datas {
		if data.Id == bevid {
			return &data
		}
	}
	return nil
}

/*
func (this *ConfigMgr) GetAIDataByScore(score uint32, t uint32) (ani *XmlAIInitDataLevel,
	bev *XmlAIBehaveData,
	speedRate float64) {
	//glog.Info("GetAIDataByScore:", len(this.Global.AIData.Sc2data.Levels))
	aidata := this.GetRoomAI(t)
	for _, data := range aidata.Sc2data.Levels {
		if score >= data.MinScore && score <= data.MaxScore {
			//glog.Info("GetAIDataByScore:", score, " id=", data)
			ani = this.GetAiInitDataLevel(data.Initdata, aidata)
			bev = this.GetAiBevData(data.Behave, aidata)
			speedRate = data.SpeedRate
			return
		}
	}
	glog.Error("GetAIDataByScore ", score, ", ", t)
	return
}
*/
func (this *ConfigMgr) GetAISrcDataByScore(score uint32) *XmlAISc2dataLevel {
	//glog.Info("GetAIDataByScore:", len(this.Global.AIData.Sc2data.Levels))
	aidata := this.GetRoomAI()
	for _, data := range aidata.Sc2data.Levels {
		if score >= data.MinScore && score <= data.MaxScore {
			//glog.Info("GetAISrcDataByScore:", score, " id=", data)
			return data
		}
	}
	glog.Error("GetAIDataByScore ", score)
	return nil
}
