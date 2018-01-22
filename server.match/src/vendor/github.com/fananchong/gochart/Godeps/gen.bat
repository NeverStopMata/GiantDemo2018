rd /q /s src
set GOPATH=%~dp0
cd %GOPATH%\..
godep.exe restore
pause