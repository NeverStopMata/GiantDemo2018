@echo off 

set srcPath=%cd%\

set distCsPath=%srcPath%..\..\..\mope_client\unity\Assets\Script\Network\ProtoMsg

svn ci %distCsPath% -m protol