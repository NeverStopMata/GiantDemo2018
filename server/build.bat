set GOPATH=E:\\GitLab\\mope;E:\\GitLab\\mope\\base;E:\\GitLab\\mope\\server;D:\\gopath;D:\project\Giant\_DEMO_PROJECT\GiantDemo2018\server

cd src\chatdbserver
go build -o ..\..\bin\chatdbserver.exe

cd ..\chatserver
go build -o ..\..\bin\chatserver.exe

cd ..\dbserver
go build -o ..\..\bin\dbserver.exe

cd ..\gatewayserver
go build -o ..\..\bin\gatewayserver.exe

cd ..\gmserver
go build -o ..\..\bin\gmserver.exe

cd ..\loginserver
go build -o ..\..\bin\loginserver.exe

cd ..\mgrserver
go build -o ..\..\bin\mgrserver.exe

cd ..\qiniuuploadserver
go build -o ..\..\bin\qiniuuploadserver.exe

cd ..\rcenterserver
go build -o ..\..\bin\rcenterserver.exe

cd ..\roomserver
go build -o ..\..\bin\roomserver.exe

cd ..\tcenterserver
go build -o ..\..\bin\tcenterserver.exe

cd ..\teamserver
go build -o ..\..\bin\teamserver.exe

cd ..\voiceserver
go build -o ..\..\bin\voiceserver.exe

pause
