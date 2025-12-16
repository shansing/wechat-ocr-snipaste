# WeChat OCR for Snipaste

Use wechat OCR in Snipaste.

## Usage

Currently only Snipaste PRO supports OCR.

This program is built for Windows x64 with [AVX2](https://github.com/swigger/wechat-ocr/issues/23#issuecomment-2566790847) CPU. Arm64 and x86 are not supported.

Go to [releases](https://github.com/shansing/wechat-ocr-snipaste/releases/) to download a pre-build `wechat-ocr-snipaste.exe`, or compile from source yourself.

1. Make sure WeChat PC is installed.
  - Prepare a param called `ocrBin`, which is full path (including file name) of `wxocr.dll` or `wechatocr.exe`. Examples:
    - `C:\Users\name\AppData\Roaming\Tencent\WeChat\XPlugin\Plugins\WeChatOCR\7061\extracted\WeChatOCR.exe`
    - `C:\Users\name\AppData\Roaming\Tencent\xwechat\XPlugin\Plugins\WeChatOcr\8075\extracted\wxocr.dll`
  - Prepare a param called `wechatDir`, which is WeChat path (usually with a version in path, typically where `mmmojo_64.dll` exists). Examples:
      - `C:\Program Files\Tencent\WeChat\[3.9.8.25]`
      - `C:\Program Files\Tencent\Weixin\4.1.5.30`
  - It is recommended to copy the parent directory of these two paths to another location.
2. Download `wcocr.dll`, and save it to the same folder as `wechat-ocr-snipaste.exe`.
  - You can get the dll from [swigger/wechat-ocr](https://github.com/swigger/wechat-ocr/releases/tag/demo-7) or [fanchenggang/wechat-ocr-go](https://github.com/fanchenggang/wechat-ocr-go/raw/refs/heads/main/wcocr.dll). 
3. Configure Snipaste:
  - Go to Snipaste -> Preferences... -> Output -> Text Recognition
  - Set OCR Engine to Tesseract
  - Set Executable to the `wechat-ocr-snipaste.exe`
  - Set Options to `-ocrBin "..." -wechatDir "..."` following step 1.
    - The default paths are `ocrBin\wxocr.dll` and `wechat\wechatDir`, where `wechat-ocr-snipaste.exe` is located.
    - Do not use a backslash `\` at the end when enclosing quotation marks `"`.

This program creates temporary files when running, and delete them when exited normally. That is due to limitations imposed by the upstream class library.

## Debug

If the configuration does not work, you can try to run `cmd` and execute the following command to see error logs:

```bat
chcp 65001
type "C:\your\image.png" | "C:\your\wechat-ocr-snipaste\wechat-ocr-snipaste.exe" stdin stdout -ocrBin "..." -wechatDir "..."
```

## Special Thanks

This project will not exist without them:

- [swigger/wechat-ocr](https://github.com/swigger/wechat-ocr/)
- [fanchenggang/wechat-ocr-go](https://github.com/fanchenggang/wechat-ocr-go/)
- [Coxxs/sp-oneocr](https://github.com/Coxxs/sp-oneocr)

## Disclaimer

The author has no affiliation with WeChat/WeiXin. This project is for personal study only, and the user is solely responsible for all consequences.
