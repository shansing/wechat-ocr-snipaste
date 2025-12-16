
# download wxocr.dll first
# apt install gcc-mingw-w64-x86-64-win32
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags "-s -w" -o "wechat-ocr-snipaste.exe" wechat-ocr-snipaste
