set curdir=%cd%
git pull
call stop.bat
cd /D %curdir%
call build.bat
cd /D %curdir%
call run.bat
cd /D %curdir%