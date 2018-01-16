# 目录说明

本目录包含协议定义和用于协议代码生成的脚本。

* bin/ 协议生成代码的工具执行文件，如 protoc.exe
* github.com/ google/ 协议依赖文件
* *.proto 协议定义
	+ 与RoomServer相关的协议都定义在wilds.proto中
* *.bat 协议生成代码的脚本
	+ gogoproto.bat 生成服务器端 go 代码
		- protoc根据.proto文件，生成go代码到目录：mope\server\src\usercmd
	+ proto_to_python.bat 生成 py_guiclient 测试客户端所用的 python 代码
	+ protoClient.bat 生成客户端的 cs 文件
		- 使用时须先更改 distCsPath

