@echo off 

set srcPath=%cd%\

set distCsPath=%srcPath%..\..\..\mope_client\unity\Assets\Script\Network\ProtoMsg

set binPath=%srcPath%\bin

::for /r "%srcPath%" %%i in (*.proto) do ( %binPath%\protoGen -i:%%i -o:%distCsPath%\%%~ni.cs  )

%binPath%\protoGen -i:%srcPath%\player.proto -o:%distCsPath%\player.cs
%binPath%\protoGen -i:%srcPath%\team.proto -o:%distCsPath%\team.cs
%binPath%\protoGen -i:%srcPath%\wilds.proto -o:%distCsPath%\wilds.cs

rem for /r "%srcPath%" %%i in (*.proto) do ( %binPath%\protoc --gogofaster_out=%distGoPath% %%i --proto_path=%srcPath% )
 
echo "ok"

pause