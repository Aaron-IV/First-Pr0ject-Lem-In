package helpers

type Ant struct {
	Name    string
	Current int
	Path    []string
	X       int
	PrevX   int
	Y       int
	PrevY   int
}

var (
	N         int
	Err       error
	RoomNames []string
	Matrix    [][]string
	BestGroup [][]string
)

type VizAnt struct {
	So    []Ant
	Path  [][]int
	Index int
}
