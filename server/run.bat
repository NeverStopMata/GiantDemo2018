cd bin

start dbserver.exe
start rcenterserver.exe
start tcenterserver.exe

ping 127.0.0.1 -n 30 -w 1000 > nul

REM start gatewayserver.exe
start loginserver.exe
start roomserver.exe

