package domain

type Strategy interface {
	Evaluate() Action
}

type Action string

const (
	Buy       Action = "BUY"
	Sell      Action = "SELL"
	DoNothing Action = "DONOTHING"
)

type doNothingStrategy struct{}

// NewDoNothingStrategy - creates a DoNothing Strategy
func NewDoNothingStategy() Strategy {
	return &doNothingStrategy{}
}

// Evaluate - unsurprisingly does nothing...
func (d *doNothingStrategy) Evaluate() Action {

	// You say it best... when you say nothing at all...
	return DoNothing
}
