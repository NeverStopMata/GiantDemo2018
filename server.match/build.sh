export GOPATH=/home/mjc/go/mope:/home/mjc/go/mope/base:/home/mjc/go/mope/server

cd ./src/chatdbserver
echo 'build chatdbserver'
go build -o ../../bin/chatdbserver

echo 'build chatserver'
cd ../chatserver
go build -o ../../bin/chatserver

echo 'build dbserver'
cd ../dbserver
go build -o ../../bin/dbserver

echo 'build gatewayserver'
cd ../gatewayserver
go build -o ../../bin/gatewayserver

echo 'build gmserver'
cd ../gmserver
go build -o ../../bin/gmserver

echo 'build loginserver'
cd ../loginserver
go build -o ../../bin/loginserver

echo 'build mgrserver'
cd ../mgrserver
go build -o ../../bin/mgrserver

echo 'build qiniuuploadserver'
cd ../qiniuuploadserver
go build -o ../../bin/qiniuuploadserver

echo 'build rcenterserver'
cd ../rcenterserver
go build -o ../../bin/rcenterserver

echo 'build roomserver'
cd ../roomserver
go build -o ../../bin/roomserver

echo 'build tcenterserver'
cd ../tcenterserver
go build -o ../../bin/tcenterserver

echo 'build teamserver'
cd ../teamserver
go build -o ../../bin/teamserver

echo 'build voiceserver'
cd ../voiceserver
go build -o ../../bin/voiceserver

