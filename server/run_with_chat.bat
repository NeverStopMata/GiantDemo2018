cd bin

start mgrserver.exe

ping 127.0.0.1 -n 10 -w 1000 > nul

start dbserver.exe
start rcenterserver.exe
start tcenterserver.exe

start chatdbserver.exe

ping 127.0.0.1 -n 30 -w 1000 > nul

start gatewayserver.exe
start loginserver.exe
start roomserver.exe


ping 127.0.0.1 -n 10 -w 1000 > nul

start chatserver.exe
