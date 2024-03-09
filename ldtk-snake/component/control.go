package component

type Control struct {
	MoveSpeed float64 // travel speed.
}

func NewControl(moveSpeed float64) Control {
	return Control{
		MoveSpeed: moveSpeed,
	}
}
