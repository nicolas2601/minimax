package recurring

// RuleDTO is the wire shape. Mirrors Rule but with safe nullability.
type RuleDTO struct {
	*Rule
}

// RunDTO mirrors Run with executed_at serialized cleanly.
type RunDTO struct {
	*Run
}

func (r *Rule) ToDTO() *RuleDTO     { return &RuleDTO{Rule: r} }
func (r *Run) ToDTO() *RunDTO       { return &RunDTO{Run: r} }
