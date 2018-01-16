nohup ./dbserver &
nohup ./rcenterserver&
nohup ./tcenterserver&

ping 127.0.0.1 -n 30 -w 1000 > nul

nohup ./gatewayserver&
nohup ./loginserver&
nohup ./roomserver&