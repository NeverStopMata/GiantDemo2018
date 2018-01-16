# 目录说明

* `base`		base 库
* `common`		common 库, 依赖 base
* `usercmd`		proto 生成消息代码
* `vender`		第3方库
* `...server`	各个服务代码，如 roomserver, 依赖所有其他库

Todo: base 库中有些第3方库，应该移到 vender 中去。
有些是第3方库的改编版本，这些可以从原库fork出来一个版本并修改，
然后也放到 vender 中。

mgo：https://github.com/go-mgo/mgo/archive/r2016.08.01.tar.gz

(*TcpTask)CheckAndSend() 没在 RoomServer 中用到，其他服有用吗，是否可删？
(*TcpTask) AsyncSendWithHead() 也没用？

