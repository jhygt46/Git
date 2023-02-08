package bytes

import (
	"testing"
)

type ResEmp struct {
	Id         int     `json:"Id"`
	Nombre     []byte  `json:"Nombre"` // Nombre
	Lat        float64 `json:"Lat"`    // Lat
	Lng        float64 `json:"Lng"`    // Lng
	TotalBytes int     `json:"TotalBytes"`
	Unique     uint8   `json:"Unique"`
}

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
func Benchmark_Repeat(b *testing.B) {
	var Emps []ResEmp = make([]ResEmp, 20)
	Emps[0] = ResEmp{Id: 70001, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[1] = ResEmp{Id: 70002, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[2] = ResEmp{Id: 70003, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[3] = ResEmp{Id: 70004, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[4] = ResEmp{Id: 70005, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[5] = ResEmp{Id: 70006, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[6] = ResEmp{Id: 70007, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[7] = ResEmp{Id: 70008, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[8] = ResEmp{Id: 70009, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[9] = ResEmp{Id: 70010, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[10] = ResEmp{Id: 70011, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[11] = ResEmp{Id: 70012, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[12] = ResEmp{Id: 70013, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[13] = ResEmp{Id: 70014, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[14] = ResEmp{Id: 70015, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[15] = ResEmp{Id: 70016, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[16] = ResEmp{Id: 70017, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[17] = ResEmp{Id: 70018, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[18] = ResEmp{Id: 70019, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[19] = ResEmp{Id: 70020, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	for m := 0; m < b.N; m++ {
		Repeat(Emps)
	}
}
func Benchmark_Repeat2(b *testing.B) {
	var Emps []ResEmp = make([]ResEmp, 20)
	Emps[0] = ResEmp{Id: 70001, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[1] = ResEmp{Id: 70002, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[2] = ResEmp{Id: 70003, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[3] = ResEmp{Id: 70004, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[4] = ResEmp{Id: 70005, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[5] = ResEmp{Id: 70006, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[6] = ResEmp{Id: 70007, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[7] = ResEmp{Id: 70008, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[8] = ResEmp{Id: 70009, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[9] = ResEmp{Id: 70010, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[10] = ResEmp{Id: 70011, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[11] = ResEmp{Id: 70012, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[12] = ResEmp{Id: 70013, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[13] = ResEmp{Id: 70014, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[14] = ResEmp{Id: 70015, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[15] = ResEmp{Id: 70016, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[16] = ResEmp{Id: 70017, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[17] = ResEmp{Id: 70018, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[18] = ResEmp{Id: 70019, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[19] = ResEmp{Id: 70020, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	for m := 0; m < b.N; m++ {
		Repeat2(Emps)
	}
}
func Benchmark_Repeat3(b *testing.B) {
	var Emps []ResEmp = make([]ResEmp, 20)
	Emps[0] = ResEmp{Id: 70001, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[1] = ResEmp{Id: 70002, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[2] = ResEmp{Id: 70003, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[3] = ResEmp{Id: 70004, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[4] = ResEmp{Id: 70005, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[5] = ResEmp{Id: 70006, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[6] = ResEmp{Id: 70007, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[7] = ResEmp{Id: 70008, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[8] = ResEmp{Id: 70009, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[9] = ResEmp{Id: 70010, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[10] = ResEmp{Id: 70011, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[11] = ResEmp{Id: 70012, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[12] = ResEmp{Id: 70013, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[13] = ResEmp{Id: 70014, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[14] = ResEmp{Id: 70015, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[15] = ResEmp{Id: 70016, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[16] = ResEmp{Id: 70017, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[17] = ResEmp{Id: 70018, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[18] = ResEmp{Id: 70019, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	Emps[19] = ResEmp{Id: 70020, Nombre: []byte{65, 65, 65, 66, 67, 68}, Lat: 231.3444, Lng: 150.3324, TotalBytes: 345}
	for m := 0; m < b.N; m++ {
		Repeat3(Emps)
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
func Repeat(Emps []ResEmp) {
	res := make(map[int]ResEmp)
	for _, emp := range Emps {
		if _, Found := res[emp.Id]; !Found {
			res[emp.Id] = emp
		}
	}
}
func Repeat2(Emps []ResEmp) {

	var j, i int = 0, 1
	var leng = len(Emps)
	for {
		if Emps[j].Id == Emps[i].Id {
			Emps[i].Unique = 1
		}
		i++
		if leng <= i {
			j++
			i = j + 1
		}
		if leng == j || leng == i {
			break
		}
	}
}
func Repeat3(Emps []ResEmp) {

	var arrInt []int = make([]int, len(Emps))
	var arrInt2 []int = make([]int, len(Emps))
	var insert bool
	var count int = 0
	var TotalBytes = 0
	for x, emp := range Emps {
		insert = true
		for _, v := range arrInt {
			if emp.Id == v {
				insert = false
			}
			if emp.Id == 0 {
				break
			}
		}
		if insert {
			arrInt[count] = emp.Id
			arrInt2[count] = x
			count++
			TotalBytes += emp.TotalBytes
		}
	}
}
