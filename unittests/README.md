# ZGG 单元测试

用 ZGG 脚本编写的语言测试集，覆盖**基本语法**、**内置函数与类型方法**、**内置标准库**。

## 运行

```bash
# 仓库根目录
./unittests/run.sh

# 或指定解释器
ZGG_BIN=/path/to/zgg ./unittests/run.sh

# 直接运行（此时退出码不可靠，见下）
zgg unittests/run_all.zgg
```

`run.sh` 会设置 `CAN_EVAL=1`（`eval` 测试需要），并根据输出里的 `ALL PASSED` / `FAILED` 返回 0/1。

> **为什么需要 run.sh**：当前 zgg 解释器在断言失败、未捕获异常、语法错误时仍返回退出码 0，所以不能用进程退出码判断测试是否通过。

## 目录结构

```text
unittests/
├── framework.zgg              # 极简测试框架（注册、断言、汇总）
├── run_all.zgg                # 汇总运行器
├── run.sh                     # 包装脚本（负责退出码）
├── test_syntax_basics.zgg     # 变量、字面量、作用域、类型、空值
├── test_operators.zgg         # 运算符、优先级、??、?.、!、%
├── test_strings.zgg           # 字符串字面量与方法
├── test_collections.zgg       # 数组、对象、推导式、容器方法
├── test_control_flow.zgg      # if/when/switch/循环/defer/use
├── test_functions.zgg         # 函数、闭包、参数、隐式返回
├── test_classes.zgg           # 类、继承、魔术方法、extend
├── test_errors.zgg            # try/catch/finally、assert、!、??
├── test_modules.zgg           # import/export、导入类型、热重载
├── test_concurrency.zgg       # future、all/map、Mutex、Chan、Limiter
├── test_builtins.zgg          # 内置函数
├── test_type_methods.zgg      # Str/Array/Object/Bytes 方法
├── test_stdlib_json_file.zgg  # json / yaml / file
├── test_stdlib_time_regex.zgg # time / regex / regex2
├── test_stdlib_http.zgg       # http 客户端、服务端、WebSocket
├── test_stdlib_db.zgg         # SQLite、QueryResult、ActiveRecord、事务
├── test_stdlib_redis.zgg      # Redis（自动探测 127.0.0.1:6379，不可用则整组跳过）
├── test_stdlib_data.zgg       # ptable / xlsx / template / url / base64 / hash / compress / binary / msgpack
├── test_stdlib_cli.zgg        # sys / sh / flag / path / math / random / psutil
└── test_stdlib_network.zgg    # tcp / udp / kv（nsq/ip2region 需外部依赖，跳过）
```

## 编写新测试

```zgg
t := import('./framework')

t.test('模块: 用例名', () => {
    t.assertEq(1 + 1, 2)
    t.assertTrue(true)
    t.assertContains('abc', 'b')
    t.assertThrows(() => {
        x := 1 / 0
    }, 'division by zero')
})

if isMain {
    t.run()
}
```

把新文件加入 `run_all.zgg` 的 `import` 列表即可。

### 断言 API

| 方法 | 说明 |
|---|---|
| `assertTrue(v[, msg])` / `assertFalse(v[, msg])` | 真值断言 |
| `assertEq(actual, expected[, msg])` | 相等（结构相等） |
| `assertNe(actual, expected[, msg])` | 不等 |
| `assertNil(v[, msg])` / `assertNotNil(v[, msg])` | nil 判断 |
| `assertUndefined(v[, msg])` | undefined 判断 |
| `assertClose(actual, expected, eps[, msg])` | 浮点近似 |
| `assertContains(haystack, needle[, msg])` | 字符串/数组包含 |
| `assertLen(v, n[, msg])` | 长度 |
| `assertThrows(fn[, substr[, msg]])` | 期望抛出异常（可校验消息） |
| `test(name, fn)` / `skip(name, reason)` | 注册用例 / 跳过用例 |

## 注意

- 测试共享一个解释器进程，模块级状态（如 `extend Str`）只在其模块作用域内生效。
- HTTP/TCP/UDP 用例使用固定端口（19180、19200、19201、19202），如被占用请修改。
- `test_stdlib_redis.zgg` 会先探测 `127.0.0.1:6379`：能 `PING` 通且无需密码时执行 11 个用例，否则整组 `skip` 并打印原因（连接失败 / NOAUTH 等）。
- Redis 用例只读写 `zgg:unittest:` 前缀的键，并在用例内清理，不会清空你的库。
- `test_stdlib_network.zgg` 中 nsq / ip2region 依赖外部服务或数据文件，默认 `skip`。
- 依赖外部服务/数据文件（NSQ、ip2region）的用例标记为 `skip`；实现缺陷修复后，请同步取消对应 `skip`。
- `db` 用例会在 `/tmp` 下创建 SQLite 文件；如需干净环境可先删除 `/tmp/zgg_unittest.db`。
