# Tray Sound

A small Windows tray app that checks the Kyiv air-alert status every minute
and plays a sound when it changes.

- **Alert active:** the bell icon takes the color of the danger level
  (yellow, orange, red, purple), plays the alert sound and shows a notification.
- **No alert:** green icon, plays the all-clear sound.
- **Server unreachable:** gray icon, no sound.

The tray menu lets you check now, test both sounds, or exit.
Only one copy runs per user session.

## Download

Latest version (these links always point to the newest release):

- **[tray-sound.msi](https://github.com/mavr77/tray-sound/releases/latest/download/tray-sound.msi)**:
  installs to `C:\Program Files\tray-sound\` and starts automatically at logon
  for every user. Needs administrator rights.
- **[tray-sound.exe](https://github.com/mavr77/tray-sound/releases/latest/download/tray-sound.exe)**:
  the program alone, no installation. Run it directly.

Release notes and older versions: [latest release](https://github.com/mavr77/tray-sound/releases/latest) ·
[all releases](https://github.com/mavr77/tray-sound/releases).

The files are not code-signed, so Windows SmartScreen may warn on first run
(choose **More info → Run anyway**).

Logs are written to `%LOCALAPPDATA%\tray-sound\tray-sound.log`.

## Build from source

Requirements: Windows, [Go](https://go.dev/dl/) 1.25+.

```bat
build.bat
```

This produces `tray-sound.exe` (no console window).

Run the tests:

```bat
go test ./...
```

### Building the MSI

Requires the .NET 8 SDK and [WiX Toolset](https://wixtoolset.org/) v7.
One-time setup:

```bat
dotnet tool install --global wix
wix eula accept wix7
wix extension add -g WixToolset.Util.wixext/7.0.1
```

Then:

```bat
build-msi.bat 1.0.0
```

Output: `dist\tray-sound-1.0.0.msi`, the same file as `dist\tray-sound.msi`,
and `dist\tray-sound.exe`.
Keep the `UpgradeCode` in `installer\tray-sound.wxs` unchanged so new versions
upgrade old ones.

### Publishing a release

Upload the versionless names so the "latest" links in this README keep working:

```bat
build-msi.bat 1.0.1
gh release create v1.0.1 dist\tray-sound.msi dist\tray-sound.exe --title "Tray Sound 1.0.1" --generate-notes
```

## Project layout

| Path | What it is |
| --- | --- |
| `main.go` | Tray setup, polling loop, logging |
| `alerts.go` | Alert API URL, JSON parsing, icon colors |
| `sound_windows.go` | Plays the embedded WAV files |
| `instance_windows.go` | Single-instance check (named mutex) |
| `sounds/` | `sound1.wav` (alert) and `sound2.wav` (all clear) |
| `installer/` | WiX installer definition and icon generator |
| `systray/` | Patched copy of [gogpu/systray](https://github.com/gogpu/systray) (MIT) with silent balloon notifications |
