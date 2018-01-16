@echo off 

set srcPath=%cd%\
 
set distGoPath=%srcPath%..\src\usercmd
 
set binPath=%srcPath%\bin
 
%binPath%\protoc --gogofaster_out=%distGoPath% wilds.proto
%binPath%\protoc --gogofaster_out=%distGoPath% gateway.proto

%binPath%\protoc --gogofaster_out=%distGoPath% player.proto
%binPath%\protoc --gogofaster_out=%distGoPath% team.proto
%binPath%\protoc --gogofaster_out=%distGoPath% server.proto
%binPath%\protoc --gogofaster_out=%distGoPath% sns.proto
%binPath%\protoc --gogofaster_out=%distGoPath% chat.proto
 
echo "ok"
pause