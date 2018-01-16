package consts

// 这里罗列一些系统的缺省值
// 根据玩法不同，可能会把这些值写到配置中，不同情况会有不同值。

const (
	DefaultBallSize      = 0.5          // 缺省球大小（在配置中找不到玩家球大小时，使用该值。）
	DefaultBallSpeed     = 1.0          // 缺省球速度
	DefaultMaxHP         = 100          // 缺省球HP最大值
	DefaultMaxMP         = 100          // 缺省球MP最大值
	DefaultHpRecover     = 2            // 缺省每秒恢复HP
	DefaultMpRecover     = 2            // 缺省每秒恢复MP
	DefaultBallFoodExp   = 2            // 缺省食物球经验值
	DefaultBallPlayerExp = 20           // 缺省玩家球经验值
	DefaultRunRatio      = 2.0          // 缺省奔跑时的速度
	DefaultRunCostMP     = 3.0 / 1000.0 // 缺省奔跑时每毫秒消耗MP
	DefaultEatFoodRange  = 0.4          // 缺省吃食物距离
	DefaultAttack        = 10           // 缺省攻击伤害
)
