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

func BenchmarkMultipleInt(b *testing.B) {
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
}
func BenchmarkDecode1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < LOOKUPS; i++ {
			Decode1(454345)
		}
	}
}
func BenchmarkDecode2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < LOOKUPS; i++ {
			Decode2(454345)
		}
	}
}

func Decode1(num int) []int {

	v1 := num / 8 // CANT ARR CUADS
	num = num % 8
	v2 := num / 4 // CANT BYTE CATEGORIA 2-3-4
	num = num % 4
	return []int{v1, v2, num / 2, num % 2}
}
func Decode2(num int) [4]int {

	v1 := num / 8 // CANT ARR CUADS
	num = num % 8
	v2 := num / 4 // CANT BYTE CATEGORIA 2-3-4
	num = num % 4
	return [4]int{v1, v2, num / 2, num % 2}
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
