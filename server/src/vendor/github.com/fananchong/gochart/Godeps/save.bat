rd /q /s D:\tmp
mkdir D:\tmp
move /y src D:\tmp\src
set CURDIR=%~dp0
set BASEDIR=%~dp0
set BASEDIR=%BASEDIR:\src\github.com\fananchong\gochart\Godeps\=\%
set GOPATH=%BASEDIR%;D:\tmp
cd %CURDIR%\..
godep.exe save -v ./...
rd /q /s vendor
cd %CURDIR%
gen.bat