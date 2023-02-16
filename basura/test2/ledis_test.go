package bytes

import (
	"fmt"
	"testing"
)

type MyHandler struct {
	Busqueda map[uint8]BusquedaCat `json:"Busqueda"`
}
type BusquedaCat struct {
	Cat map[uint32]BusquedaCuad `json:"Cat"`
}
type BusquedaCuad struct {
	Cuad       map[uint32][]byte `json:"Cuad"`
	CacheProds [][]byte          `json:"CacheProds"`
}

func Benchmark_GETIP(b *testing.B) {

	h := &MyHandler{}
	h.Save1()
	for m := 0; m < b.N; m++ {
		pais := uint8(m % 256)
		if _, Found := h.Busqueda[pais]; Found {

		}
	}
}

func (h *MyHandler) Save1() {
	h.Busqueda = make(map[uint8]BusquedaCat, 0)
	for i := 0; i < 256; i++ {
		h.Busqueda[uint8(i)] = BusquedaCat{Cat: make(map[uint32]BusquedaCuad, 0)}
		for j := 0; j < 200; j++ {
			h.Busqueda[uint8(i)].Cat[uint32(j)] = BusquedaCuad{Cuad: make(map[uint32][]byte, 0), CacheProds: make([][]byte, 0)}
			for z := 0; z < 200; z++ {
				h.Busqueda[uint8(i)].Cat[uint32(j)].Cuad[uint32(z)] = []byte{}
			}
		}
	}
}
func aux() {
	fmt.Println("AUX")
}
