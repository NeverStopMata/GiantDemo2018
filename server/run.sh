#!/bin/sh

startwork()
{
	#rm -rf /log/*
	nohup $PWD/bin/dbserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/voiceserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/tcenterserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/gatewayserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/rcenterserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/loginserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/teamserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/roomserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $PWD/bin/gmserver -config=./bin/config.json > /dev/null 2>&1 &
	sleep 1
	nohup $PWD/bin/mgrserver -config=./bin/config.json &
	sleep 1
	nohup $PWD/bin/chatdbserver -config=./bin/config.json &
	sleep 1
	nohup $PWD/bin/chatserver -config=./bin/config.json &

	echo "ps x | grep \"server\""
	ps x|grep "server"   
}

stopwork()
{
	SERVERLIST='roomserver teamserver loginserver rcenterserver gatewayserver tcenterserver dbserver voiceserver gmserver mgrserver chatserver chatdbserver' 

	for serv in $SERVERLIST
	do
		echo -n "stop $serv "
		ps aux|grep "$PWD/.*server" | sed -e "/grep/d" | grep "$serv"|awk '{print $2}'|xargs kill 2&>/dev/null
        while test -f run.sh
        do
			echo -n "."
			#count=`ps x|grep -w $serv|sed -e '/grep/d'|wc -l`
			count=`ps x |grep "$PWD/.*server" | sed -e '/grep/d' |grep -c "$serv"`
            if [ $count -eq 0 ]; then
                break
            fi
            sleep 0.5
        done
#echo "ok"
		echo -e "\033[1;40;32mok\033[0m"
	done
    echo "running server:"`ps x|grep "server -c"|sed -e '/grep/d'|wc -l`
}

echo "-------------------start-----------------------"

case $1 in
stop)
    stopwork
;;
start)
    startwork
;;
*)
    stopwork
    sleep 0.5
    startwork
;;
esac

echo "-------------------end-----------------------"

