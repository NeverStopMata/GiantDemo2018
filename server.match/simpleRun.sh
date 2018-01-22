#!/bin/sh

CURRDIR=$PWD
cd $PWD/bin

startwork()
{
	#rm -rf /log/*
	nohup $CURRDIR/bin/dbserver -config=./config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $CURRDIR/bin/tcenterserver -config=./config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $CURRDIR/bin/rcenterserver -config=./config.json > /dev/null 2>&1 &
	sleep 20
	nohup $CURRDIR/bin/gatewayserver -config=./config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $CURRDIR/bin/loginserver -config=./config.json > /dev/null 2>&1 &
	sleep 0.5
	nohup $CURRDIR/bin/roomserver -config=./config.json > /dev/null 2>&1 &
	sleep 2

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

