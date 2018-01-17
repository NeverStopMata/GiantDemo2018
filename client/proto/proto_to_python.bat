@echo off 

set srcPath=%cd%\
 
set distGoPath=%srcPath%..\tools\py_guiclient
 
set binPath=%srcPath%\bin
 
%binPath%\protoc_python.exe --python_out=%distGoPath%\proto\ player.proto  -I ../src/vendor -I . -I ../src/vendor/github.com/gogo/protobuf/protobuf
%binPath%\protoc_python.exe --python_out=%distGoPath%\proto\ wilds.proto
%binPath%\protoc_python.exe --python_out=%distGoPath%\proto\ team.proto  -I ../src/vendor -I . -I ../src/vendor/github.com/gogo/protobuf/protobuf
%binPath%\protoc_python.exe --python_out=%distGoPath%\proto\ chat.proto
%binPath%\sed-win\sed.exe -i "s/import github.com.gogo.protobuf.gogoproto.gogo_pb2//g" %distGoPath%\proto\player_pb2.py
%binPath%\sed-win\sed.exe -i "s/import github.com.gogo.protobuf.gogoproto.gogo_pb2//g" %distGoPath%\proto\wilds_pb2.py
%binPath%\sed-win\sed.exe -i "s/import github.com.gogo.protobuf.gogoproto.gogo_pb2//g" %distGoPath%\proto\team_pb2.py

del /Q sed*

echo "ok"
pause