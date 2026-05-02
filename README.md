# Formula-Go

一个用 Go 语言实现的通达信公式解析器和执行引擎，为开发者和量化交易者提供解析和执行通达信技术指标公式的能力。

## 项目状态

🚧 **持续完善中** - 核心解析、执行、85 个内置函数、基础绘图事件输出和常用通达信兼容语法已实现并测试通过

本项目参考 [formula-ts](https://github.com/DTrader-store/formula-ts) TypeScript 实现，使用 Go 语言重新实现。

## 特性

- ✅ **类型安全**: 使用 Go 的强类型系统，确保代码安全性
- ✅ **完整实现**: 词法分析、语法分析、解释执行全流程
- ✅ **丰富的内置函数**: 85 个内置函数，覆盖常用技术指标、数学统计、引用、逻辑和基础绘图事件输出
- ✅ **通达信兼容语法**: 支持 `:` 输出声明、常用样式后缀、行情字段别名和字符串字面量
- ✅ **易于集成**: 简洁的 API 设计，易于集成到现有项目
- ✅ **测试覆盖**: 注册表、README 函数清单和内置函数覆盖用例已建立一致性校验

## 安装

```bash
go get github.com/DTrader-store/formula-go
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/DTrader-store/formula-go"
)

func main() {
    // 创建市场数据
    data := []*formula.MarketData{
        formula.NewMarketData(100, 105, 107, 99, 1000, 100000),
        formula.NewMarketData(105, 103, 108, 102, 1100, 110000),
        formula.NewMarketData(103, 107, 109, 101, 1200, 120000),
        formula.NewMarketData(107, 110, 112, 106, 1300, 130000),
        formula.NewMarketData(110, 108, 113, 107, 1400, 140000),
    }

    // 创建公式引擎
    engine := formula.NewFormulaEngine()

    // 执行公式
    result, err := engine.Run("MA5 := MA(CLOSE, 5)", data)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 输出结果
    for _, output := range result.Outputs {
        fmt.Printf("%s: %v\n", output.Name, output.Data)
    }
}
```

## 支持的功能

### 1. 语法特性

- ✅ 变量声明: `MA5 := MA(CLOSE, 5)`
- ✅ 算术运算: `+`, `-`, `*`, `/`
- ✅ 比较运算: `>`, `<`, `>=`, `<=`, `=`, `<>`
- ✅ 逻辑运算: `AND`, `OR`
- ✅ 函数调用: `MA(CLOSE, 5)`
- ✅ 括号表达式: `(a + b) * c`
- ✅ 一元运算: `-x`
- ✅ 通达信输出声明: `DIF: EMA(CLOSE, 12), COLORWHITE, LINETHICK2`
- ✅ 输出样式后缀: `COLOR*`, `LINETHICK*`, `DOTLINE`, `STICK`, `COLORSTICK`, `VOLSTICK`, `NODRAW`
- ✅ 字符串字面量: `'UP'`, `"UP"`，可用于 `DRAWTEXT`
- ✅ 外部指标引用字面量: `"MACD.DIF#WEEK"`，当前仅在执行环境已存在同名变量时可解析

### 2. 内置函数

现已支持 **85 个内置函数**！

**数学统计函数**
- `MA(data, period)` - 简单移动平均
- `SMA(data, period)` - 简单移动平均（MA 的别名）
- `SMA(data, period, weight)` - 通达信递推平滑移动平均
- `EMA(data, period)` - 指数移动平均
- `WMA(data, period)` - 加权移动平均
- `DMA(data, alpha)` - 动态移动平均
- `SUM(data, period)` - 求和
- `STD(data, period)` - 标准差
- `VAR(data, period)` - 方差
- `AVEDEV(data, period)` - 平均绝对偏差
- `MAX(a, b)` - 最大值
- `MIN(a, b)` - 最小值
- `ABS(value)` - 绝对值
- `SQRT(value)` - 平方根
- `POW(base, exponent)` - 幂运算
- `EXP(value)` - 自然指数
- `LN(value)` - 自然对数
- `LOG(value)` - 常用对数
- `MOD(a, b)` - 取模
- `CEILING(value)` - 向上取整
- `FLOOR(value)` - 向下取整
- `INTPART(value)` - 整数部分
- `FRACPART(value)` - 小数部分
- `ROUND(value)` - 四舍五入
- `ROUND2(value, digits)` - 指定位数四舍五入
- `SIGN(value)` - 符号函数
- `SIN(value)` - 正弦
- `COS(value)` - 余弦
- `TAN(value)` - 正切
- `ASIN(value)` - 反正弦
- `ACOS(value)` - 反余弦
- `ATAN(value)` - 反正切
- `STDP(data, period)` - 总体标准差
- `STDDEV(data, period)` - 样本标准差
- `VARP(data, period)` - 总体方差
- `DEVSQ(data, period)` - 离差平方和
- `FORCAST(data, period)` - 线性回归预测值
- `SLOPE(data, period)` - 线性回归斜率
- `COVAR(a, b, period)` - 协方差
- `RELATE(a, b, period)` - 相关系数
- `BETA(a, b, period)` - Beta 系数

**引用函数**
- `REF(data, n)` - 引用 n 期前的数据
- `REFV(data, n)` - `REF` 的兼容变体
- `REFX(data, n)` - 引用 n 期后的数据
- `REFXV(data, n)` - `REFX` 的兼容变体
- `HHV(data, period)` - 周期内最高值
- `LLV(data, period)` - 周期内最低值
- `HHVBARS(data, period)` - 周期内最高值距离当前的周期数
- `LLVBARS(data, period)` - 周期内最低值距离当前的周期数
- `CURRBARSCOUNT()` - 到最后交易日的周期数
- `TOTALBARSCOUNT()` - 总周期数
- `ISLASTBAR()` - 是否最后一个周期
- `BARSTATUS()` - 当前周期状态
- `SUMBARS(data, target)` - 累计值达到目标所需周期数

**条件和逻辑函数**
- `IF(condition, trueValue, falseValue)` - 条件判断
- `IFF(condition, trueValue, falseValue)` - `IF` 的兼容别名
- `IFN(condition, falseValue, trueValue)` - 反向条件判断
- `NOT(value)` - 逻辑取反
- `COUNT(condition, period)` - 统计满足条件的周期数
- `EVERY(condition, period)` - 检查是否所有周期都满足条件
- `EXIST(condition, period)` - 检查是否存在满足条件的周期
- `EXISTR(condition, from, to)` - 检查指定历史区间是否存在满足条件
- `BETWEEN(value, lower, upper)` - 检查值是否在范围内
- `RANGE(value, lower, upper)` - `BETWEEN` 的兼容别名
- `CONST(value)` - 使用最后一个值填充全序列
- `VALUEWHEN(condition, value)` - 条件满足时记录并保持对应值
- `DRAWNULL()` - 返回空值用于断线/隐藏

**技术分析函数**
- `CROSS(a, b)` - 交叉检测（a 上穿 b）
- `LONGCROSS(a, b, period)` - 持续低于后上穿检测
- `BARSLAST(condition)` - 距离最后一次满足条件的周期数
- `BARSCOUNT(data)` - 有效数据周期计数
- `BARSSINCE(condition)` - 距离第一次满足条件的周期数
- `BARSLASTCOUNT(condition)` - 当前连续满足条件的周期数
- `UPNDAY(data, period)` - 连续上涨检测
- `DOWNNDAY(data, period)` - 连续下跌检测
- `NDAY(a, b, period)` - 连续大于检测
- `LAST(condition, from, to)` - 指定历史区间持续满足检测
- `FILTER(condition, period)` - 过滤信号，防止频繁触发

**绘图事件函数**
- `DRAWTEXT(condition, price, text)` - 条件满足时输出文字标注事件
- `DRAWICON(condition, price, iconType)` - 条件满足时输出图标标注事件
- `DRAWNUMBER(condition, price, number)` - 条件满足时输出数值标注事件
- `STICKLINE(condition, price1, price2, width, empty)` - 条件满足时输出柱线事件
- `DRAWLINE(cond1, price1, cond2, price2, expand)` - 输出由起点条件和终点条件连接的线段事件
- `POLYLINE(condition, price)` - 条件满足时输出折线点事件
- `DRAWBAND(upper, upperColor, lower, lowerColor)` - 输出上下轨之间的带状区域事件
- `DRAWKLINE(high, open, low, close)` - 输出自定义 K 线事件

### 3. 内置变量

- `OPEN` - 开盘价
- `CLOSE` - 收盘价
- `HIGH` - 最高价
- `LOW` - 最低价
- `VOLUME` - 成交量
- `AMOUNT` - 成交额
- `O`, `C`, `H`, `L`, `V`, `VOL`, `AMO` - 通达信常用行情字段别名

### 4. 当前边界

- `MarketData` 当前只包含 `Open`, `Close`, `High`, `Low`, `Volume`, `Amount` 六个数值字段。
- 暂不支持需要日期或时间索引的函数，例如 `REFDATE`。如需实现，需要先扩展行情数据模型。
- 暂不支持需要财务、盘口、成本分布或逐笔数据的函数，例如 `FINANCE`, `DYNAINFO`, `COST`, `WINNER`。
- 绘图事件函数只返回结构化 `DrawingEvent`，不直接负责渲染图表；调用方可按自身图表库适配 `Function`, `BarIndex`, `Values`, `Text`, `Meta`。

## 使用示例

### 简单移动平均

```go
formula := "MA5 := MA(CLOSE, 5)"
result, _ := engine.Run(formula, marketData)
```

### MACD 指标

```go
formula := `
    EMA12 := EMA(CLOSE, 12)
    EMA26 := EMA(CLOSE, 26)
    DIF := EMA12 - EMA26
    DEA := EMA(DIF, 9)
    MACD := (DIF - DEA) * 2
`
result, _ := engine.Run(formula, marketData)
```

### 金叉检测

```go
formula := `
    MA5 := MA(CLOSE, 5)
    MA10 := MA(CLOSE, 10)
    SIGNAL := CROSS(MA5, MA10)
`
result, _ := engine.Run(formula, marketData)
```

### 条件选股

```go
formula := `
    MA5 := MA(CLOSE, 5)
    MA10 := MA(CLOSE, 10)
    GOLDEN := CROSS(MA5, MA10)
    STRONG := CLOSE > MA5 AND EVERY(CLOSE > OPEN, 3)
    SELECT := GOLDEN AND STRONG
`
result, _ := engine.Run(formula, marketData)
// SELECT 中为 1 的位置表示满足选股条件
```

### 布林带指标

```go
formula := `
    MID := MA(CLOSE, 20)
    STDEV := STD(CLOSE, 20)
    UPPER := MID + 2 * STDEV
    LOWER := MID - 2 * STDEV
    BREAK_UP := CROSS(CLOSE, UPPER)
    BREAK_DOWN := CROSS(LOWER, CLOSE)
`
result, _ := engine.Run(formula, marketData)
```

### KDJ 指标

```go
formula := `
    LOW9 := LLV(LOW, 9)
    HIGH9 := HHV(HIGH, 9)
    RSV := (CLOSE - LOW9) / (HIGH9 - LOW9) * 100
    K := SMA(RSV, 3, 1)
    D := SMA(K, 3, 1)
    J := 3 * K - 2 * D
`
result, _ := engine.Run(formula, marketData)
```

### 信号过滤

```go
formula := `
    MA5 := MA(CLOSE, 5)
    MA10 := MA(CLOSE, 10)
    GOLDEN := CROSS(MA5, MA10)
    FILTERED := FILTER(GOLDEN, 10)
`
result, _ := engine.Run(formula, marketData)
// FILTERED 会过滤掉 10 个周期内的重复信号
```

### 通达信输出样式

```go
formula := `
    DIF: EMA(CLOSE, 12) - EMA(CLOSE, 26), COLORWHITE, LINETHICK2
    DEA: EMA(DIF, 9), COLORYELLOW
    MACD: (DIF - DEA) * 2, COLORSTICK
`
result, _ := engine.Run(formula, marketData)
// OutputLine.Style 中会保留颜色、线宽、绘制方式和隐藏标记
```

### 绘图事件

```go
formula := `
    UP := CLOSE > OPEN
    MARK := DRAWTEXT(UP, LOW, 'UP')
`
result, _ := engine.Run(formula, marketData)
// result.Drawings 中包含 DRAWTEXT 事件，调用方可自行适配图表渲染
```

```go
formula := `
    START := BARSTATUS() = 1
    END := ISLASTBAR()
    TREND := DRAWLINE(START, LOW, END, HIGH, 0)
    BAND := DRAWBAND(HIGH, 1, LOW, 2)
    KLINE := DRAWKLINE(HIGH, OPEN, LOW, CLOSE)
`
result, _ := engine.Run(formula, marketData)
// result.Drawings 中包含 DRAWLINE、DRAWBAND、DRAWKLINE 等结构化事件
```

## 项目结构

```
formula-go/
├── engine/              # 公式引擎
│   ├── engine.go       # FormulaEngine 主类
│   └── engine_test.go  # 集成测试
├── errors/              # 错误类型定义
│   ├── errors.go       # 各类错误
│   └── errors_test.go
├── interpreter/         # 解释器
│   ├── interpreter.go  # 解释执行
│   ├── functions.go    # 内置函数
│   └── registry.go     # 函数注册
├── lexer/              # 词法分析器
│   ├── lexer.go        # 词法分析主逻辑
│   ├── token.go        # Token 定义
│   ├── token_type.go   # Token 类型
│   └── *_test.go
├── parser/             # 语法分析器
│   ├── parser.go       # 语法分析主逻辑
│   ├── parser_test.go
│   └── ast/            # 抽象语法树
│       └── nodes.go
├── types/              # 类型定义
│   ├── market_data.go  # 市场数据
│   ├── formula_result.go # 公式结果
│   └── *_test.go
├── formula.go          # 主入口，导出 API
├── go.mod
└── README.md
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test ./... -cover

# 运行详细测试
go test ./... -v

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

测试约束：
- `engine/builtin_registry_test.go` 会校验注册表中的每个内置函数都有覆盖用例。
- README 的内置函数清单会和 `interpreter.FunctionRegistry.Names()` 自动比对。
- 新增函数如果只注册、不补测试或不更新 README，`go test ./...` 会失败。

说明：`interpreter` 包的函数主要通过 `engine` 集成测试覆盖；`go test ./... -cover` 会按包独立统计，因此 `interpreter` 会显示为无直接测试文件。

## API 文档

### FormulaEngine

```go
type FormulaEngine struct {}

// 创建新引擎
func NewFormulaEngine() *FormulaEngine

// 编译公式为 AST
func (e *FormulaEngine) Compile(formula string) (*Program, error)

// 执行已编译的程序
func (e *FormulaEngine) Execute(program *Program, marketData []*MarketData) (*FormulaResult, error)

// 一步编译并执行
func (e *FormulaEngine) Run(formula string, marketData []*MarketData) (*FormulaResult, error)
```

### MarketData

```go
type MarketData struct {
    Open   float64
    Close  float64
    High   float64
    Low    float64
    Volume float64
    Amount float64
}

func NewMarketData(open, close, high, low, volume, amount float64) *MarketData
func (m *MarketData) Validate() error
```

### FormulaResult

```go
type FormulaResult struct {
    Outputs   []*OutputLine
    Variables map[string]float64
    Drawings  []*DrawingEvent
}

type OutputLine struct {
    Name  string
    Data  []float64
    Style *LineStyle
}

type DrawingEvent struct {
    Function string
    BarIndex int
    Values   map[string]float64
    Text     string
    Meta     map[string]string
}
```

## 开发路线图

### ✅ Phase 1: 基础类型系统

- [x] 错误处理系统
- [x] Token 系统
- [x] AST 节点定义
- [x] 市场数据类型
- [x] 公式结果类型

### ✅ Phase 2: 词法分析器和语法分析器

- [x] 实现 Lexer（词法分析器）
- [x] 实现 Parser（语法分析器）
- [x] 支持基础语法规则
- [x] 完整的错误报告

### ✅ Phase 3: 解释器和内置函数

- [x] 实现 Interpreter（解释器）
- [x] 实现函数注册机制和首批核心内置函数
- [x] 变量管理和求值
- [x] 数组和标量运算

### ✅ Phase 4: 通达信兼容增强

- [x] 通达信输出声明、样式后缀和行情别名
- [x] 常用数学、统计、引用、逻辑、技术分析和基础绘图事件输出函数
- [x] 内置函数注册表、README 清单和覆盖用例一致性校验

### 🚧 后续完善方向

- [ ] `REFDATE`：需要日期字段或时间索引支持
- [ ] `FINANCE` / `DYNAINFO` / `COST` / `WINNER`：需要扩展数据模型
- [ ] 函数参数边界、错误语义和 NaN 传播规则专项测试
- [ ] 更多绘图事件函数和下游图表库适配示例
- [ ] 增量计算与性能优化
- [ ] 格式化器和更完整的示例文档

## 性能

项目已提供可重复运行的 benchmark 基线，用于后续优化前后对比：

```bash
go test ./engine -bench=. -benchmem
```

当前覆盖场景：
- `BenchmarkCompileMACD`：只测试公式编译。
- `BenchmarkRunMACD`：测试 MACD 公式的完整编译 + 执行。
- `BenchmarkExecuteCompiledMACD`：测试已编译 AST 的重复执行。
- `BenchmarkRunRollingFunctions`：测试 `MA/SUM/STD` 等滚动函数。
- `BenchmarkRunDrawingEvents`：测试基础绘图事件生成。

优化性能时应先记录修改前后的 `ns/op`, `B/op`, `allocs/op`，并保留 `go test ./...` 通过结果。

## 参考项目

- [formula-ts](https://github.com/DTrader-store/formula-ts) - TypeScript 实现版本

## 技术要求

- Go 版本以 `go.mod` 为准
- 公式解析和执行核心无外部运行时服务依赖

## 开发

```bash
# 克隆仓库
git clone https://github.com/DTrader-store/formula-go.git
cd formula-go

# 运行测试
go test ./...

# 构建
go build

# 格式化代码
go fmt ./...

# 静态检查
go vet ./...
```

### 新增内置函数流程

1. 在 `interpreter/functions.go` 或 `interpreter/functions_ext.go` 实现函数。
2. 在 `interpreter/registry.go` 注册函数名。
3. 在 `engine/builtin_registry_test.go` 的 `builtinCoverageCases()` 中新增精确覆盖 case。
4. 在 README 的“内置函数”小节加入函数签名和说明。
5. 运行 `go test ./...`，确认注册表、测试覆盖和 README 清单一致。

## 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 许可证

ISC License

## 联系方式

- GitHub: https://github.com/DTrader-store/formula-go
- Issues: https://github.com/DTrader-store/formula-go/issues

---

**最后更新**: 2026-05-02
**状态**: 核心解析执行、85 个内置函数、基础绘图事件输出和常用通达信兼容语法已实现；数据模型扩展和更多函数仍在推进
