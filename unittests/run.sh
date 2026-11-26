#!/usr/bin/env bash
# 运行全部单元测试。
# zgg 进程即使断言失败也返回 0，因此这里根据输出中的标记判断结果。
set -u

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -n "${ZGG_BIN:-}" ]; then
    ZGG="$ZGG_BIN"
elif command -v zgg >/dev/null 2>&1; then
    ZGG="$(command -v zgg)"
elif [ -x "$DIR/../bin/zgg" ]; then
    ZGG="$DIR/../bin/zgg"
else
    echo "找不到 zgg：请设置 ZGG_BIN，或把 zgg 加入 PATH" >&2
    exit 2
fi

export CAN_EVAL=1
out="$("$ZGG" "$DIR/run_all.zgg" 2>&1)"
code=$?
printf '%s\n' "$out"

if printf '%s\n' "$out" | grep -q '^FAILED'; then
    exit 1
fi
if printf '%s\n' "$out" | grep -q '^ALL PASSED'; then
    exit 0
fi
echo "测试未正常结束（zgg 退出码 $code）" >&2
exit 2
