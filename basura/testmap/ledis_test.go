package main

// run with: $ go test --bench=. -test.benchmem .
// @see https://twitter.com/karlseguin/status/524452778093977600
import (
	"fmt"
	"testing"
)

const (
	SIZE    = 1000000
	LOOKUPS = 300
)

type MyHandler struct {
	Busqueda  map[uint8]BusquedaCat2 `json:"Busqueda"`
	Busqueda2 map[uint64][]byte      `json:"Busqueda2"`
}
type BusquedaCat struct {
	Cat map[uint32]BusquedaCuad `json:"Cat"`
}
type BusquedaCat2 map[uint32]BusquedaCuad

type BusquedaCuad struct {
	Cuad       map[uint32][]byte `json:"Cuad"`
	CacheProds [][]byte          `json:"CacheProds"`
}

func BenchmarkMapString(b *testing.B) {
	lookup := make(map[string][]byte, SIZE)
	for i := 0; i < SIZE; i += 1 {
		lookup[string(ParamIntto7Bytes2(i))] = RandBytes(500)
	}
	var x, y int
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LOOKUPS; i++ {
			if _, ok := lookup[string(ParamIntto7Bytes2(n+i))]; ok {
				x++
			} else {
				y++
			}
		}
	}
	//fmt.Println("RES", x, y)
}
func BenchmarkMapByte(b *testing.B) {
	lookup := make(map[[7]byte][]byte, SIZE)
	for i := 0; i < SIZE; i += 1 {
		lookup[ParamIntto7Bytes(i)] = RandBytes(500)
	}
	var x, y int
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LOOKUPS; i++ {
			if _, ok := lookup[ParamIntto7Bytes(n+i)]; ok {
				x++
			} else {
				y++
			}
		}
	}
	//fmt.Println("RES", x, y)
}
func BenchmarkMapUint64(b *testing.B) {
	lookup := make(map[uint64][]byte, SIZE)
	for i := 0; i < SIZE; i += 1 {
		lookup[uint64(i)] = RandBytes(500)
	}
	var x, y int
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LOOKUPS; i++ {
			if _, ok := lookup[uint64(n+i)]; ok {
				x++
			} else {
				y++
			}
		}
	}
	//fmt.Println("RES", x, y)
}

func Aux() {
	fmt.Println("AUX")
}
func ParamIntto3Bytes(n int) []byte {
	res := make([]byte, 3)
	res[0] = uint8(n / 65536)
	n = n % 65536
	res[1], res[2] = uint8(n/256), uint8(n%256)
	return res
}
func ParamIntto7Bytes(num int) [7]byte {
	b := [7]byte{}
	var r int = num % 16777216
	b[3] = uint8(num / 16777216)
	b[4] = uint8(r / 65536)
	r = r % 65536
	b[5], b[6] = uint8(r/256), uint8(r%256)
	return b
}
func ParamIntto7Bytes2(num int) []byte {
	b := make([]byte, 7)
	var r int = num % 16777216
	b[3] = uint8(num / 16777216)
	b[4] = uint8(r / 65536)
	r = r % 65536
	b[5], b[6] = uint8(r/256), uint8(r%256)
	return b
}
func RandBytes(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[uint8(i%256)] = uint8(i % 256)
	}
	return bytes
}
