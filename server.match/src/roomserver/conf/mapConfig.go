package conf

import (
	"base/env"
	"base/glog"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
)

type MapNodeConfig struct {
	Id     uint32  `json:"id"`
	Type   int     `json:"type"`
	Px     float64 `json:"px"`
	Py     float64 `json:"py"`
	Radius float64 `json:"radius"`
}

//AreaNodeConfig 区域节点
type AreaNodeConfig struct {
	Id     uint32  `json:"id"`
	Type   int     `json:"type"`
	Layer  uint8   `json:"layer"`
	Px     float64 `json:"px"`
	Py     float64 `json:"py"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

//CheckIn 确认是否在范围内
func (t *AreaNodeConfig) CheckIn(x, y float64) bool {
	return (math.Abs(x-t.Px) < t.Width/2) && (math.Abs(y-t.Py) < t.Height/2)
}

//RandPos 随机一个范围内的位置
func (t *AreaNodeConfig) RandPos() (x, y float64) {
	x = t.Px + t.Width*rand.Float64() - t.Width/2
	y = t.Py + t.Height*rand.Float64() - t.Height/2
	return
}

type MapConfig struct {
	Title string           `json:"title"`
	Size  float64          `json:"size"`
	Nodes []*MapNodeConfig `json:"nodes"`
}

//var _mapConfig *MapConfig
//var _qmapConfig *MapConfig

var _mapConfigDic map[uint32]*MapConfig

func LoadMapConfig(path string) (*MapConfig, bool) {
	config := MapConfig{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Info("LoadMapConfig fail:", err.Error())
		return nil, false
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		glog.Info("LoadMapConfig ummarshal fail:", err.Error())
		return nil, false
	}

	glog.Info("LoadMapConfig:", config.Size, " nodes:", len(config.Nodes))
	return &config, true
}

func InitMapConfig() bool {
	_mapConfigDic = make(map[uint32]*MapConfig)
	for _, m := range ConfigMgr_GetMe().Map.Scenes {
		path := fmt.Sprintf("%s%d.json", env.Get("global", "terraincfg"), m.Id)
		glog.Info("LoadMapConfig:" + path)
		if config, ok := LoadMapConfig(path); ok {
			_mapConfigDic[m.Id] = config
		} else {
			return false
		}

	}

	return true
}

/*
func InitQMapConfig() bool {
	if config, ok := LoadMapConfig(env.Get("global", "terraincfg") + "1002.json"); ok {
		_qmapConfig = config
		return true
	}
	return false
}
*/
func GetMapConfigById(t uint32) *MapConfig {
	val := _mapConfigDic[t]
	return val
}
