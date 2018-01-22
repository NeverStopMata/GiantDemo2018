package common

// 检查Session是否存在
type CheckSession struct {
	Uid     uint64
	Session string
}

type RetCheckSession struct {
	Ret bool
}
