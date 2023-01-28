package bytes

import (
	"testing"
)

func Benchmark_GETIP(b *testing.B) {
	var by string = "127.0.0.1:7653"
	for m := 0; m < b.N; m++ {
		GetIp(by)
	}
}
func Benchmark_GETIP2(b *testing.B) {
	var by string = "127.0.0.1:7653"
	for m := 0; m < b.N; m++ {
		GetIp2(by)
	}
}
func Benchmark_Unicode(b *testing.B) {
	var by []byte = []byte{240, 159, 144, 167, 240, 159, 144, 167, 240, 159, 144, 167}
	for m := 0; m < b.N; m++ {
		Unicode(by, 1)
	}
}
func Benchmark_Unicode2(b *testing.B) {
	var by []byte = []byte{240, 159, 144, 167, 240, 159, 144, 167, 240, 159, 144, 167}
	for m := 0; m < b.N; m++ {
		Unicode2(by, 1)
	}
}
func Benchmark_ParamAlphaLat(b *testing.B) {
	var by []byte = []byte{51, 83, 86, 111, 69, 115}
	for m := 0; m < b.N; m++ {
		ParamAlphaLat(by)
	}
}
func Benchmark_ParamAlphaLat2(b *testing.B) {
	var by []byte = []byte{51, 54, 48, 48, 48, 48, 48, 48, 48, 48}
	for m := 0; m < b.N; m++ {
		ParamAlphaLat2(by)
	}
}

func GetIp(ip string) uint32 {
	var b uint32
	var x uint32 = 16777216
	var n uint32 = 0
	for _, i := range ip {
		if i == 46 || i == 58 {
			b = b + x*n
			n = 0
			x = x / 256
			if i == 58 {
				return b
			}
		} else {
			n = n*10 + uint32(i-'0')
		}
	}
	return b
}
func GetIp2(ip string) uint32 {
	var b uint32
	var x uint32 = 16777216
	var n uint32 = 0
	for _, i := range ip {
		if i > 47 && i < 58 {
			n = n*10 + uint32(i-'0')
		} else {
			b = b + x*n
			n = 0
			x = x / 256
		}
	}
	return b
}

func LengMax31(b []byte) int {
	if len(b) > 0 {
		if b[0] > 48 && b[0] < 52 {
			return int(b[0] - '0')
		} else {
			return 0
		}
	} else {
		return 0
	}
}
func LengMax32(b []byte) int {
	if len(b) > 0 {
		if b[0] == 49 {
			return 1
		} else if b[0] == 50 {
			return 2
		} else if b[0] == 51 {
			return 3
		} else {
			return 0
		}
	} else {
		return 0
	}
}

func Unicode(str []byte, n int) []byte {

	var l int = len(str)
	var count int = 0
	var x int = 0

	if n > 0 {
		for i := 1; i <= l; i++ {
			if str[l-i] < 128 {
				x++
			} else if str[l-i] > 127 && str[l-i] < 192 {

			} else if str[l-i] > 193 && str[l-i] < 224 {
				x++
			} else if str[l-i] > 223 && str[l-i] < 240 {
				x++
			} else {
				x++
			}
			count++
			if x == n {
				return str[:l-count]
			}
		}
	}
	return str
}
func Unicode2(str []byte, n int) []byte {

	var l int = len(str)
	var count int = 0
	var x int = 0
	if n > 0 {
		for i := 1; i <= l; i++ {
			if str[l-i] < 128 || str[l-i] > 191 {
				x++
			}
			count++
			if x == n {
				return str[:l-count]
			}
		}
	}
	return str
}

func ParamAlphaLat2(b []byte) (float64, bool) {

	if len(b) > 10 {
		return 0, false
	}

	var x float64
	for _, c := range b {
		if c > 47 && c < 58 {
			x = x*10 + float64(c-'0')
		} else {
			return 0, false
		}
	}
	return (x - 1800000000) / 10000000, true
}
func ParamAlphaLat(b []byte) (float64, bool) {

	var leng int = len(b)

	if leng > 6 {
		return 0, false
	}

	var res float64 = 0
	var pon float64 = 1
	var w uint8

	for i := 0; i < leng; i++ {
		w = b[leng-i-1]
		if w > 47 && w < 58 {
			res = res + float64(w-48)*pon
		} else if w > 96 && w < 123 {
			res = res + float64(w-87)*pon
		} else if w > 64 && w < 91 {
			res = res + float64(w-29)*pon
		} else {
			return 0, false
		}
		pon = pon * 62
	}
	return (res - 1800000000) / 10000000, true
}
