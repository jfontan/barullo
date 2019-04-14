package barullo

type Node interface {
	Get(offset int, buf []float64)
}
