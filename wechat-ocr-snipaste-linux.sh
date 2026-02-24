#!/bin/bash

# AI assistanted

# 参数默认值
STDIN_FLAG=false
STDOUT_FLAG=false
OCR_BIN=""
WECHAT_DIR=""
TEST_CLI=""

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -stdin|stdin)
            STDIN_FLAG=true
            shift
            ;;
        -stdout|stdout)
            STDOUT_FLAG=true
            shift
            ;;
        -ocrBin)
            OCR_BIN="$2"
            shift 2
            ;;
        -wechatDir)
            WECHAT_DIR="$2"
            shift 2
            ;;
        -testCli)
            TEST_CLI="$2"
            shift 2
            ;;
        *)
            echo "[Error] Unknown parameter: $1" >&2
            exit 1
            ;;
    esac
done

# 检查必需参数
if [ "$STDIN_FLAG" != true ] || [ "$STDOUT_FLAG" != true ]; then
    echo "[Error] \"stdin stdout\" must be set. This program is for Snipaste." >&2
    exit 1
fi

if [ -z "$OCR_BIN" ]; then
    echo "[Error] ocrBin parameter is required" >&2
    exit 1
fi

if [ -z "$WECHAT_DIR" ]; then
    echo "[Error] wechatDir parameter is required" >&2
    exit 1
fi

# 如果未指定 test_cli，尝试在脚本同目录或 PATH 中查找
if [ -z "$TEST_CLI" ]; then
    if [ -x "$SCRIPT_DIR/test_cli" ]; then
        TEST_CLI="$SCRIPT_DIR/test_cli"
    elif command -v test_cli &> /dev/null; then
        TEST_CLI="test_cli"
    else
        echo "[Error] test_cli not found. Please specify with -testCli parameter" >&2
        exit 1
    fi
fi

# 检查路径是否存在
if [ ! -e "$OCR_BIN" ]; then
    echo "[Error] ocrBin seems not to exist: $OCR_BIN" >&2
    exit 1
fi

if [ ! -d "$WECHAT_DIR" ]; then
    echo "[Error] wechatDir seems not to exist: $WECHAT_DIR" >&2
    exit 1
fi

if ! command -v "$TEST_CLI" &> /dev/null && [ ! -x "$TEST_CLI" ]; then
    echo "[Error] test_cli not executable: $TEST_CLI" >&2
    exit 1
fi

echo "[Info] test_cli: $TEST_CLI" >&2
echo "[Info] ocrBin: $OCR_BIN" >&2
echo "[Info] wechatDir: $WECHAT_DIR" >&2

# 生成随机文件名
UUID=$(cat /proc/sys/kernel/random/uuid)
TEMP_FILE="${TMPDIR:-/tmp}/${UUID}.png"

echo "[Info] tempFile: $TEMP_FILE" >&2

# 设置清理陷阱，确保临时文件被删除
cleanup() {
    if [ -f "$TEMP_FILE" ]; then
        rm -f "$TEMP_FILE"
        echo "[Info] Cleaned up temp file" >&2
    fi
}
trap cleanup EXIT INT TERM

# 从标准输入读取图片并保存到临时文件
cat > "$TEMP_FILE"
if [ $? -ne 0 ] || [ ! -f "$TEMP_FILE" ]; then
    echo "[Error] Failed to write stdin to temp file" >&2
    exit 1
fi

# 调用 test_cli 并捕获输出
# 格式：./test_cli <wechatocr_exe> <wechat_dir> <test_png>
OCR_OUTPUT=$("$TEST_CLI" "$OCR_BIN" "$WECHAT_DIR" "$TEMP_FILE" 2>&1)
OCR_EXIT_CODE=$?
echo "[OCR_EXIT_CODE]$OCR_EXIT_CODE" >&2
echo "[---OCR_OUTPUT---]$OCR_OUTPUT[/--OCR_OUTPUT---]" >&2

# 检查 test_cli 是否成功执行
if [ $OCR_EXIT_CODE -ne 0 ]; then
    echo "[Error] test_cli failed with exit code $OCR_EXIT_CODE" >&2
    echo "$OCR_OUTPUT" >&2
    exit 1
fi

# 检查 errcode
if ! echo "$OCR_OUTPUT" | grep -q "OCR errcode=0"; then
    echo "[Error] OCR errcode is not 0" >&2
    echo "$OCR_OUTPUT" >&2
    exit 1
fi

# 解析输出，只保留 OCR 识别的文本
# 匹配格式：[...] r=数字 文本
echo "$OCR_OUTPUT" | grep -E '^\[.*\] r=[0-9.]+' | sed -E 's/^\[.*\] r=[0-9.]+ //'

exit 0
