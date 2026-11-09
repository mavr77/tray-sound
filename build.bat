@echo off
go mod tidy || exit /b 1
go build -ldflags "-H=windowsgui -s -w" -o tray-sound.exe . || exit /b 1
echo Built tray-sound.exe
