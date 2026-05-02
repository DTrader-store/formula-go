package types

// LineStyle represents line style configuration for output visualization
type LineStyle struct {
	Color      string // Color of the line (e.g., '#FF0000', 'red')
	LineWidth  int    // Width of the line in pixels
	LineStyle  string // Style of the line ('solid', 'dashed', 'dotted', etc.)
	DrawMethod string // TDX draw method such as stick/colorstick/volstick
	Hidden     bool   // Whether output is hidden from chart drawing
}

// OutputLine represents a single output line representing calculated data
type OutputLine struct {
	Name  string     // Name/identifier of the output line
	Data  []float64  // Data points for the output line
	Style *LineStyle // Optional style configuration for visualization
}

// DrawingEvent represents a rendering-agnostic drawing event emitted by formulas.
type DrawingEvent struct {
	Function string             // Builtin function that emitted the event, such as DRAWTEXT or DRAWLINE
	BarIndex int                // Primary bar index for the event
	Values   map[string]float64 // Numeric payload used by downstream chart adapters
	Text     string             // Optional text payload
	Meta     map[string]string  // Optional non-numeric payload for future adapters
}

// FormulaResult represents the result of formula calculation containing outputs and variables
type FormulaResult struct {
	Outputs   []*OutputLine      // Array of output lines from the formula calculation
	Variables map[string]float64 // Calculated variables and their values
	Drawings  []*DrawingEvent    // Rendering-agnostic drawing events emitted by formulas
}

// NewFormulaResult creates a new FormulaResult instance
func NewFormulaResult() *FormulaResult {
	return &FormulaResult{
		Outputs:   make([]*OutputLine, 0),
		Variables: make(map[string]float64),
		Drawings:  make([]*DrawingEvent, 0),
	}
}

// AddOutput adds an output line to the result
func (f *FormulaResult) AddOutput(name string, data []float64, style *LineStyle) {
	f.Outputs = append(f.Outputs, &OutputLine{
		Name:  name,
		Data:  data,
		Style: style,
	})
}

// AddDrawing adds a rendering-agnostic drawing event to the result.
func (f *FormulaResult) AddDrawing(event *DrawingEvent) {
	f.Drawings = append(f.Drawings, event)
}

// SetVariable sets a variable value in the result
func (f *FormulaResult) SetVariable(name string, value float64) {
	f.Variables[name] = value
}

// GetVariable gets a variable value from the result
func (f *FormulaResult) GetVariable(name string) (float64, bool) {
	value, exists := f.Variables[name]
	return value, exists
}
