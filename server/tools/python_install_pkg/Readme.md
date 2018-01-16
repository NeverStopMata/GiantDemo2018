# Python环境安装介绍

1. 安装python-3.6.2-amd64.exe
2. 设置环境变量（第一步安装过程中也可以勾选设置）
3. 解压Lib.rar，覆盖Python安装目录中的Lib目录


## Lib.rar

内网不方便，pip install安装python第3方库，因此直接拷贝了Lib目录。

内包括了 wxPython、protobuf 的python模块、KCP模块

## KCP
https://github.com/sunzhaoping/python-ikcp

## 验证是否安装成功

打开CMD控制台，键python 回车，查看python版本号 为 python3.6

## 验证wxPython、protobuf、KCP是否正常

python命令行下，依次执行以下命令：
```
import wx
import google.protobuf
import ikcp
```
没有出错，则表示以上模块均安装正常



