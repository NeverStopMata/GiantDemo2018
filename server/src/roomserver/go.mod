module roomserver

go 1.15

replace (
	base => ../base
	common => ../common
	usercmd => ../usercmd
)

require (
	base v0.0.0-00010101000000-000000000000
	common v0.0.0-00010101000000-000000000000
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/fananchong/gochart v0.0.0-20180117141114-a0d1b57622da
	github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/consul v1.9.3 // indirect
	github.com/shirou/gopsutil v3.21.2+incompatible
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	usercmd v0.0.0-00010101000000-000000000000
)
