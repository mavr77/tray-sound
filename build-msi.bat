@echo off
rem Builds tray-sound.exe and the per-machine MSI installer.
rem Usage: build-msi.bat [version]     (default 1.0.0)
rem Output: dist\tray-sound-<version>.msi, plus dist\tray-sound.msi (same file,
rem         versionless name for the releases/latest/download link) and
rem         dist\tray-sound.exe
rem
rem One-time setup (WiX v7 refuses to run until its EULA is accepted):
rem   wix eula accept wix7
rem   wix extension add -g WixToolset.Util.wixext/7.0.1

setlocal
cd /d "%~dp0"

set VERSION=%~1
if "%VERSION%"=="" set VERSION=1.0.0

if not exist dist mkdir dist

rem Build a separate exe in dist\ so a running dev copy in the project root
rem does not block the build.
go mod tidy || exit /b 1
go build -ldflags "-H=windowsgui -s -w" -o dist\tray-sound.exe . || exit /b 1

wix build installer\tray-sound.wxs ^
  -arch x64 ^
  -ext WixToolset.Util.wixext ^
  -d Version=%VERSION% ^
  -d ExePath="%~dp0dist\tray-sound.exe" ^
  -o dist\tray-sound-%VERSION%.msi || exit /b 1

copy /y dist\tray-sound-%VERSION%.msi dist\tray-sound.msi >nul || exit /b 1

echo Built dist\tray-sound-%VERSION%.msi (also dist\tray-sound.msi)
