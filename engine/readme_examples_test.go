package engine

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/DTrader-store/formula-go/types"
)

type readmeFormulaExample struct {
	name    string
	formula string
}

func TestReadmeFormulaExamplesExecute(t *testing.T) {
	content, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}

	examples := readmeFormulaExamples(string(content))
	if len(examples) < 8 {
		t.Fatalf("expected README formula examples, got %d", len(examples))
	}

	engine := NewFormulaEngine()
	marketData := createReadmeExampleData(40)
	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			result, err := engine.Run(example.formula, marketData)
			if err != nil {
				t.Fatalf("README formula example failed: %v\nformula:\n%s", err, example.formula)
			}
			if len(result.Outputs) == 0 && len(result.Variables) == 0 && len(result.Drawings) == 0 {
				t.Fatalf("README formula example produced an empty result:\n%s", example.formula)
			}
			if strings.Contains(example.formula, "DRAWTEXT(") && len(result.Drawings) == 0 {
				t.Fatalf("README drawing example produced no drawing events:\n%s", example.formula)
			}
			if strings.Contains(example.formula, "COLOR") || strings.Contains(example.formula, "LINETHICK") {
				assertReadmeExampleHasStyle(t, result)
			}
		})
	}
}

func readmeFormulaExamples(content string) []readmeFormulaExample {
	blocks := readmeGoCodeBlocks(content)
	examples := make([]readmeFormulaExample, 0)
	seen := make(map[string]bool)

	for blockIndex, block := range blocks {
		for _, formula := range extractFormulaAssignments(block) {
			if seen[formula] {
				continue
			}
			seen[formula] = true
			examples = append(examples, readmeFormulaExample{
				name:    fmt.Sprintf("block_%02d_%s", blockIndex+1, readmeExampleName(formula)),
				formula: formula,
			})
		}
		for _, formula := range extractInlineEngineRunFormulas(block) {
			if seen[formula] {
				continue
			}
			seen[formula] = true
			examples = append(examples, readmeFormulaExample{
				name:    fmt.Sprintf("block_%02d_%s", blockIndex+1, readmeExampleName(formula)),
				formula: formula,
			})
		}
	}
	return examples
}

func readmeGoCodeBlocks(content string) []string {
	const fence = "```"
	blocks := make([]string, 0)
	offset := 0
	for {
		start := strings.Index(content[offset:], "```go")
		if start < 0 {
			break
		}
		start += offset + len("```go")
		if strings.HasPrefix(content[start:], "\r\n") {
			start += len("\r\n")
		} else if strings.HasPrefix(content[start:], "\n") {
			start++
		}

		end := strings.Index(content[start:], fence)
		if end < 0 {
			break
		}
		blocks = append(blocks, content[start:start+end])
		offset = start + end + len(fence)
	}
	return blocks
}

func extractFormulaAssignments(block string) []string {
	return extractStringLiteralsAfter(block, "formula :=")
}

func extractInlineEngineRunFormulas(block string) []string {
	return extractStringLiteralsAfter(block, "engine.Run(")
}

func extractStringLiteralsAfter(source, marker string) []string {
	values := make([]string, 0)
	offset := 0
	for {
		index := strings.Index(source[offset:], marker)
		if index < 0 {
			break
		}
		start := offset + index + len(marker)
		value, end, ok := readGoStringLiteral(source, start)
		if ok {
			values = append(values, value)
			offset = end
			continue
		}
		offset = start
	}
	return values
}

func readGoStringLiteral(source string, start int) (string, int, bool) {
	for start < len(source) && (source[start] == ' ' || source[start] == '\t' || source[start] == '\n' || source[start] == '\r') {
		start++
	}
	if start >= len(source) {
		return "", start, false
	}

	quote := source[start]
	if quote != '"' && quote != '`' {
		return "", start, false
	}

	escaped := false
	for i := start + 1; i < len(source); i++ {
		if quote == '"' {
			if escaped {
				escaped = false
				continue
			}
			if source[i] == '\\' {
				escaped = true
				continue
			}
		}
		if source[i] == quote {
			literal := source[start : i+1]
			value, err := strconv.Unquote(literal)
			if err != nil {
				return "", i + 1, false
			}
			return value, i + 1, true
		}
	}
	return "", len(source), false
}

func readmeExampleName(formula string) string {
	for _, line := range strings.Split(formula, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.NewReplacer(" ", "_", "\t", "_", ":", "", "=", "", ",", "", "(", "_", ")", "").Replace(line)
		if len(line) > 48 {
			line = line[:48]
		}
		return line
	}
	return "empty"
}

func createReadmeExampleData(count int) []*types.MarketData {
	data := make([]*types.MarketData, count)
	for i := range data {
		open := 100.0 + float64(i)
		close := open + 1
		if i%2 == 1 {
			close = open - 1
		}
		high := open + 3
		if close+2 > high {
			high = close + 2
		}
		low := open - 3
		if close-2 < low {
			low = close - 2
		}
		volume := 1000.0 + float64(i)*100
		amount := close * volume
		data[i] = types.NewMarketData(open, close, high, low, volume, amount)
	}
	return data
}

func assertReadmeExampleHasStyle(t *testing.T, result *types.FormulaResult) {
	t.Helper()
	for _, output := range result.Outputs {
		if output.Style != nil {
			return
		}
	}
	t.Fatal("expected README style example to produce at least one styled output")
}
