# osc-watch

> 转发 bluetooth hrp 数据到 VRChat OSC


## 已验证设备

- 小米手环 5 NFC

## 使用说明

克隆项目后使用 `go build .` 编译项目，然后执行 `osc-watch -mac <你的蓝牙 mac 地址>` , 启动后会在控制台打印心率，同时发送到 vrchat 的 osc 端口 ，osc 接口模仿了 [github.com/vard88508/vrc-osc-miband-hrm](https://github.com/vard88508/vrc-osc-miband-hrm?tab=readme-ov-file#what-is-this) ，如果你有适配了这个的模型则可以直接使用。

### 配置指南

- [小米手环](./docs/mi-band.md)

### 参数列表

```bash
~ $ osc-watch -h
Usage of osc-watch:
  -addr string
        VRChat OSC 地址 (default "localhost")
  -mac string
        过滤的MAC地址
  -port int
        VRChat OSC 端口 (default 9000)

```

## LICENSE

项目使用 [Apache-2.0](./LICENSE)