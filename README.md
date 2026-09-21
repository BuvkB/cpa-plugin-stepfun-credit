# CPA StepFun Credit Tracker

CLIProxyAPI 插件：实时统计 **StepFun Step Plan** 的 Credit 消耗。

StepFun 官方没有开放订阅月池的查询接口，所以这个插件在 CPA 侧按每个请求的真实 token 用量、
结合官方单价折算 Credit，让你随时知道这个月花了多少、还剩多少。

## 功能

- 每次请求完成后实时记账（输入 / 缓存命中 / 输出 token、延迟、成败）
- 左侧菜单新增「StepFun Credit」页面，随 CPAMC 主题自动切换深浅色
- 时间范围（1 小时 ~ 全部）与粒度（分钟 / 小时 / 日 / 周 / 月）可切换，鼠标悬停看单点明细
- 支持登记多个订阅（各自档位或自定义额度），自动累加为总月池，显示已用 / 剩余 / 占比
- 内置 StepFun 官方单价；接入方式自动识别，无需手动配置 provider 名

## 安装

从 Releases 下载对应平台的 zip，解压后把动态库放进 CPA 的插件目录：

```text
plugins/linux/amd64/stepfun-credit-tracker.so
```

然后在 `config.yaml` 里启用：

```yaml
plugins:
  enabled: true
  configs:
    stepfun-credit-tracker:
      enabled: true
```

重启或热加载 CPA 后，左侧菜单会出现「StepFun Credit」。

## 前提

CPA 里已经通过 `openai-compatibility` 接入了 StepFun，例如：

```yaml
openai-compatibility:
  - name: stepfun
    base-url: https://api.stepfun.com/step_plan/v1
    api-key-entries:
      - api-key: <你的 StepFun API Key>
    models:
      - name: step-3.7-flash
      - name: step-5-preview
```

插件会从真实请求里自动识别 StepFun 接入（按 `base-url` 与模型名前缀），不需要额外配置。
若尚未检测到，插件页会直接给出上面这段可复制的配置。

## 关于额度

月池总额度需要你自己填（StepFun 不提供查询接口）：打开插件页右上角「设置额度」，
按实际订阅选档位即可。有多个订阅就逐个添加，会自动累加。此设置保存在浏览器本地。

## 说明

- Credit 折算口径：`1 元 = 1,000,000 Credit`
- 记账发生在 CPA 侧，是精确的 token 折算值；官方月池余额仍以 StepFun 控制台为准
- 插件不保存 prompt、请求正文或响应正文

## 许可

MIT
