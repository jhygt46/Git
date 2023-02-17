package bytes

import (
	"fmt"
	"testing"
)

type MyHandler struct {
	Busqueda  map[uint8]BusquedaCat `json:"Busqueda"`
	Busqueda2 map[uint64][]byte     `json:"Busqueda2"`
}
type BusquedaCat struct {
	Cat map[uint32]BusquedaCuad `json:"Cat"`
}
type BusquedaCuad struct {
	Cuad       map[uint32][]byte `json:"Cuad"`
	CacheProds [][]byte          `json:"CacheProds"`
}

func Benchmark_Map1(b *testing.B) {

	h := &MyHandler{}
	h.Save1()
	for m := 0; m < b.N; m++ {
		if s1, Found := h.Busqueda[uint8(m%100)]; Found {
			if s2, Found := s1.Cat[uint32(m%100)]; Found {
				if _, Found := s2.Cuad[uint32(m%100)]; Found {

				}
			}
		}
	}
}
func Benchmark_Map2(b *testing.B) {

	h := &MyHandler{}
	h.Save2()
	for m := 0; m < b.N; m++ {
		if _, Found := h.Busqueda2[uint64(m%1000000)]; Found {

		}
	}
}

func (h *MyHandler) Save1() {
	h.Busqueda = make(map[uint8]BusquedaCat, 0)
	for i := 0; i < 100; i++ {
		h.Busqueda[uint8(i)] = BusquedaCat{Cat: make(map[uint32]BusquedaCuad, 0)}
		for j := 0; j < 100; j++ {
			h.Busqueda[uint8(i)].Cat[uint32(j)] = BusquedaCuad{Cuad: make(map[uint32][]byte, 0), CacheProds: make([][]byte, 0)}
			for z := 0; z < 100; z++ {
				h.Busqueda[uint8(i)].Cat[uint32(j)].Cuad[uint32(z)] = RandBytes(100)
			}
		}
	}
}
func (h *MyHandler) Save2() {
	h.Busqueda2 = make(map[uint64][]byte, 0)
	for i := 0; i < 1000000; i++ {
		h.Busqueda2[uint64(i)] = RandBytes(100)
	}
}
func aux() {
	fmt.Println("AUX")
}
func RandBytes(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[uint8(i%256)] = uint8(i % 256)
	}
	return bytes
}
