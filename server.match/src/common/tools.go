package common

import (
	"base/env"
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	monthDay = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

// 返回当天天数 (从时间点January 1,1970 UTC 起)
func NowDays() uint32 {
	return uint32((time.Now().Unix() + 28800) / (3600 * 24))
}

// 返回当前小时数 (从时间点January 1,1970 UTC 起)
func NowHours() uint32 {
	return uint32((time.Now().Unix() + 28800) / 3600)
}

// 返回月数(从时间点January 1,1970 UTC 起)
func Months(t int64) int64 {
	y, m, _ := time.Unix(t, 0).Date()
	return int64(y)*12 + int64(m)
}

// 返回天数(从时间点January 1,1970 UTC 起)
func Days(t int64) int64 {
	return (t + 28800) / 86400
}

func Weeks(t int64) int64 {
	return (t - 345600) / 806400
}

// 检查是否同一天
func InSameDay(t1, t2 int64) bool {
	y1, m1, d1 := time.Unix(t1, 0).Date()
	y2, m2, d2 := time.Unix(t2, 0).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// 检查是否同一月
func InSameMonth(t1, t2 int64) bool {
	y1, m1, _ := time.Unix(t1, 0).Date()
	y2, m2, _ := time.Unix(t2, 0).Date()
	return y1 == y2 && m1 == m2
}

// 是否闰年
func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// 判断两个时间点是不是经过了星期天
func InSameWeekSunday(t1, t2 int64) bool {
	if t1 == t2 {
		return true
	}
	if t1 > t2 {
		return false
	}
	tm1 := time.Unix(t1, 0)
	tm2 := time.Unix(t2, 0)
	y1, wd1, yd1 := tm1.Year(), tm1.Weekday(), tm1.YearDay()
	y2, wd2, yd2 := tm2.Year(), tm2.Weekday(), tm2.YearDay()
	if wd1 == 0 {
		wd1 = 7
	}
	if wd2 == 0 {
		wd2 = 7
	}

	//	fmt.Println("tm1= ", tm1, " tm2= ", tm2)
	if y2-y1 > 1 {
		return false
	}
	if y2-y1 == 1 {
		if isLeap(y1) {
			yd2 = yd2 + 366 - yd1 + 1
		} else {
			yd2 = yd2 + 365 - yd1 + 1
		}
	}
	if yd2-yd1 >= 7 {
		return false
	}
	if wd1 == 7 {
		return true
	}
	if wd2 == 7 {
		return false
	}
	if wd1 <= wd2 {
		return true
	}
	return false
}

// 检查是否同一周 (周一开始)
func InSameWeek(t1, t2 int64) bool {
	if t1 == t2 {
		return true
	}
	if t1 > t2 {
		t1, t2 = t2, t1
	}
	tm1 := time.Unix(t1, 0)
	tm2 := time.Unix(t2, 0)
	y1, wd1, yd1 := tm1.Year(), tm1.Weekday(), tm1.YearDay()
	y2, wd2, yd2 := tm2.Year(), tm2.Weekday(), tm2.YearDay()
	if wd1 == 0 {
		wd1 = 7
	}
	if wd2 == 0 {
		wd2 = 7
	}
	if y1 == y2 {
		if (yd2-yd1) < 7 && wd2 >= wd1 {
			return true
		}
	} else if (y2 - y1) == 1 {
		if isLeap(y1) {
			if (366+yd2-yd1) < 7 && wd2 >= wd1 {
				return true
			}
		} else {
			if (365+yd2-yd1) < 7 && wd2 >= wd1 {
				return true
			}
		}
	}
	return false
}

// 获取当天0点
func TodayZero() int64 {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
}

// 获取明天0点
func TomorrowZero() int64 {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix() + 86400
}

// 获取下周一0点
func NextWeekZero() int64 {
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	weekday := timeNow.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix() + (8-int64(weekday))*86400
}

// 获取下下周天0点
func NextWeekTZero() int64 {
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	weekday := timeNow.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix() + (14-int64(weekday))*86400
}

//下月1号0点0分钟
func NextMonthZero() int64 {
	timeNow := time.Now()
	year, month, _ := timeNow.AddDate(0, 1, 0).Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
}

// 获取下一分钟
func NextMinute() int64 {
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	return time.Date(year, month, day, timeNow.Hour(), timeNow.Minute(), 0, 0, time.Local).Unix() + 60
}

// 获取本小时0分0秒
func HourZero() int64 {
	timeNow := time.Now()
	year, month, day := timeNow.Date()
	return time.Date(year, month, day, timeNow.Hour(), 0, 0, 0, time.Local).Unix()
}

//计算每个月的天数
func GetDaysOfMonth(year, month int) int {
	if month < 1 || month > 12 {
		return 0
	}
	if month == 2 && ((year%4 == 0 && year%100 != 0) || year%400 == 0) {
		return 29
	}
	return monthDay[month-1]
}

// 产生[min,max]之间的数
func RandBetween(min, max int64) int64 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + rand.Int63n(max-min+1)
}

// 获取几分之几的几率
func SelectByOdds(upNum, downNum int) bool {
	if downNum < 1 || upNum < 1 {
		return false
	}
	if upNum > downNum-1 {
		return true
	}
	return RandBetween(1, int64(downNum)) <= int64(upNum)
}

// 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 用gob进行数据解码
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// md5编码
func EncMd5(params string) string {
	h := md5.New()
	h.Write([]byte(params))
	return hex.EncodeToString(h.Sum(nil))
}

// md5编码 二进制
func EncMd5Bytes(params string) []byte {
	h := md5.New()
	h.Write([]byte(params))
	return h.Sum(nil)
}

// 产生唯一key
func GenerateKey(uid uint64) string {
	return EncMd5(time.Now().String() + strconv.FormatInt(int64(uid), 10) + env.Get("global", "key"))
}

func IsMoneyItem(itemid uint32) bool {
	return false
}

// 获取货币名
func GetMoneyName(mtype uint32) string {
	switch mtype {
	case MONEY_ID_1:
		return "彩豆"
	case MONEY_ID_2:
		return "棒棒糖"
	case MONEY_ID_TICKET:
		return "入场券"
	case MONEY_ID_SPEAKER:
		return "喇叭"
	}
	return "未知类型"
}

func GetWeekCnt() int64 {
	return (int64(NowDays()) + 3) / 7
}

// 随机数产生
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// 产生[min,max]之间的数
func RandBetweenInt32(min, max int32) int32 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + random.Int31n(max-min+1)
}

func StringToInt64(s string) int64 {
	if s == "" {
		return 0
	}
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i64
}

func StringToUint64(s string) uint64 {
	if s == "" {
		return 0
	}
	ui64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return ui64
}

func StringToInt32(s string) int32 {
	if s == "" {
		return 0
	}
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int32(i64)
}

func StringToInt(s string) int {
	if s == "" {
		return 0
	}
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int(i64)
}

func StringToUint16(s string) uint16 {
	if s == "" {
		return 0
	}
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint16(i64)
}

func StringToUint32(s string) uint32 {
	if s == "" {
		return 0
	}
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint32(i64)
}

func StringToFloat32(s string) float32 {
	if s == "" {
		return 0
	}
	f64, err := strconv.ParseFloat(s, 10)
	if err != nil {
		return 0
	}
	return float32(f64)
}

//列表乱序
func GetChaosListUint64(list []uint64) []uint64 {
	maxCount := len(list)
	var chaosList []uint64
	chaosList = append(chaosList, list[:]...)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < maxCount; i++ {
		curI := r.Intn(maxCount)
		chaosList[i], chaosList[curI] = chaosList[curI], chaosList[i]
	}
	return chaosList
}

//判断值是否在数组中存在
func IsInArrayUInt64(val uint64, arr []uint64) bool {
	ishave := false
	for _, v := range arr {
		if val == v {
			ishave = true
			break
		}
	}
	return ishave
}

//判断值是否在数组中存在
func IsInArrayUInt16(val uint16, arr []uint16) bool {
	ishave := false
	for _, v := range arr {
		if val == v {
			ishave = true
			break
		}
	}
	return ishave
}

//判断值是否在数组中存在
func IsInArrayUInt32(val uint32, arr []uint32) bool {
	ishave := false
	for _, v := range arr {
		if val == v {
			ishave = true
			break
		}
	}
	return ishave
}

//判断值是否在数组中存在 不区分大小写
func IsInArrayString(val string, arr []string) bool {
	ishave := false
	for _, v := range arr {
		if strings.ToLower(val) == strings.ToLower(v) {
			ishave = true
			break
		}
	}
	return ishave
}

//func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
//	math.Pi
//	math.Sqrt()
//	lat1 = lat1 * 0.1745
//	lng1 = lng1 * 0.1745
//	lat2 = lat2 * 0.1745
//	lng2 = lng2 * 0.1745

//	theta := lng2 - lng1
//	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))

//	return dist * float64(6378137)
//}

func Rad(num float64) float64 {
	return num * math.Pi / 180
}

func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	var radLat1 float64 = Rad(lat1)
	var radLat2 float64 = Rad(lat2)
	var a float64 = radLat1 - radLat2
	var b float64 = Rad(lng1) - Rad(lng2)

	var s float64 = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+
		math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	return s * 6378137
}

// 设置状态
func SetState(datas, state uint32) uint32 {
	return datas | (uint32(1) << state)
}

// 清除状态
func ClrState(datas, state uint32) uint32 {
	return datas & ^(uint32(1) << state)
}

// 检查状态
func IsState(datas, state uint32) bool {
	return datas&(uint32(1)<<state) != 0
}

func GetNetIpSlice(addr string) (bgpip, otherip string, cnet uint8) {
	if idx := strings.Index(addr, "/"); idx > 0 && idx < len(addr)-1 {
		if netidx := strings.LastIndex(addr, "_"); netidx > 0 {
			typ := addr[netidx+1:]
			addr = addr[:netidx]

			typd, _ := strconv.Atoi(typ)
			return addr[:idx], addr[idx+1:], byte(typd)
		}
		return addr[:idx], addr[idx+1:], 4
	}
	return addr, addr, 0
}

// 双线区分
func GetBgpTelAddr(addr string, cnet uint8) string {
	if idx := strings.Index(addr, "/"); idx > 0 && idx < len(addr)-1 {
		if netidx := strings.LastIndex(addr, "_"); netidx > 0 {
			typ := addr[netidx+1:]
			addr = addr[:netidx]
			if typ == strconv.Itoa(int(cnet)) {
				return addr[idx+1:]
			} else {
				return addr[:idx]
			}
		}
		if cnet == 4 {
			return addr[idx+1:]
		} else {
			return addr[:idx]
		}
	}
	return addr
}

// 获取ip地址
func GetIP(addr string) string {
	return strings.Split(addr, ":")[0]
}

func ParseVec(str string, delimiter string) []uint32 {
	ret := make([]uint32, 0)
	if str != "" {
		item := strings.Split(str, delimiter)
		for _, items := range item {
			param, _ := strconv.Atoi(items)
			ret = append(ret, uint32(param))
		}
	}
	return ret
}

func MinUint32(a, b uint32) uint32 {
	if a > b {
		return b
	}
	return a
}

func MaxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func MaxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func MinUint64(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

func CanInviteTeamUserState(userstate uint32) bool {
	if userstate == UserStateOnline || userstate == UserStateFPlaying || userstate == UserStateQPlaying {
		return true
	}
	if userstate == UserStateTeam || userstate == UserStateTeamInv {
		return false
	}
	return false
}
