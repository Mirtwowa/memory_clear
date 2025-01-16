## Readme.md文档

```markdown
# Memory Monitor

Memory Monitor 是一个使用 Go 语言编写的操作系统内存监控软件，旨在帮助用户实时监控内存使用情况，检测和清理异常进程，以保障系统应用的流畅运行。

 功能特点

- 实时内存监控：监控系统内存使用情况，当内存占用超过阈值时提供提示。
- 异常进程检测：检测占用内存过高或无响应的进程，并提供自动清理选项。
- 优先保障前台应用：尽量保证前台应用的运行，适当关闭不影响系统运行的进程。
- 通知功能：当发现异常进程时，通过邮件或桌面弹窗等方式进行通知。
- 配置文件支持：通过 `config.yaml` 文件配置监控的参数和规则。

-- 安装

1. 克隆项目：

   ```bash
   git clone git@github.com:Mirtwowa/memory_clear.git
   cd memory-monitor
   ```

2. 安装依赖：

   ```bash
   go mod tidy
   ```

3. 编译程序：

   ```bash
   go build -o memory-monitor ./cmd
   ```

4. 运行程序：

   ```bash
   ./memory-monitor
   ```

## 配置

配置文件位于 `config/config.yaml`。你可以通过该文件设置内存使用阈值、忽略的进程列表、通知方式等。

### 示例配置 (`config/config.yaml`)

```yaml
memory_threshold: 80 # 内存占用阈值，单位：百分比
process_ignore_list:
  - "systemd"
  - "bash"
notification:
  enabled: true
  method: "email"  # 支持 email 或 desktop
  email:
    recipient: "example@example.com"
    smtp_server: "smtp.example.com"
  desktop:
    enabled: true
    timeout: 5
```

## 开发

1. 克隆项目后，可以直接开始修改或扩展代码。
2. 项目使用 Go 1.18 及以上版本开发，确保你的开发环境支持 Go 语言。
3. 如果你想贡献代码，首先请确保编写相应的单元测试并通过测试。

## 测试

你可以使用以下命令运行测试：

```bash
go test ./tests
```

## 贡献

欢迎贡献代码，任何功能增强或修复都很受欢迎！如果你想贡献代码，请确保：
- 提交详细的 Pull Request 描述。
- 编写单元测试并通过所有测试。
## 修正点
```
1.Edge进程Microsoft Edge 被某些系统进程或者策略依赖，导致在终止其他进程时系统强制将它保活
2.当前窗口进程识别有待进一步优化
3.双屏或多屏情况下当前活跃窗口识别有待进一步优化
4.可通过可视化图形界面将kill进程的选择权交由用户选择
```

## License

该项目使用 [MIT License](LICENSE) 开源。
```

### 说明：
- 你可以修改 `README.md` 文件中的内容，以便与实际的功能和配置文件保持一致。
- `config.yaml` 配置文件示例也可以根据实际需求进行调整。
- 如果需要任何进一步的修改或添加其他功能文档，告诉我即可！