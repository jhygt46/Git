package main

// run with: $ go test --bench=. -test.benchmem .
// @see https://twitter.com/karlseguin/status/524452778093977600
import (
	"fmt"
	"testing"
)

type MyHandler struct {
	Busqueda   map[uint64][]byte `json:"Busqueda2"`
	CacheProds [][]byte          `json:"CacheProds"`
}
type Params struct {
	Lat          float64      `json:"Lat"`   //LISTO
	Lon          float64      `json:"Lng"`   // LISTO
	Desde        float64      `json:"Desde"` // LISTO
	Largo        int          `json:"Largo"` // LISTO
	DistanciaMax uint64       `json:"DistanciaMax"`
	Opciones     [4]uint64    `json:"Opciones"` // LISTO
	Filtros      [64][]uint64 `json:"Filtros"`  // LISTO
	Evals        []uint64     `json:"Evals"`    // LISTO
	Cuads        [][3]byte    `json:"Cuads"`    // LISTO
	PaisCat      [7]byte      `json:"PaisCat"`  //LISTO
}

func Benchmark_ParamBusqueda4(b *testing.B) {
	//by := []byte("9f7mY4Ou-pCjkgf-PoK27Bb4JZ904iWmGKO8SYp11kx5i9H8JidrFXJTFY98vc9G4ecGMEk2bJuRK6gB3yt95v4jjaP0eVM7765hUbj94CK9f_1CrB58ViazIIt1NcbCX1a99eI4Q5cgubQrfli1DNa-IZbeKW8Lqdwb3GR")
	for n := 0; n < b.N; n++ {
		Uint64_To_7Bytes(uint64(n))
	}
}
func auxs() {
	fmt.Println("BUENA")
}
func RandBytes(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[uint8(i%256)] = uint8(i % 256)
	}
	return bytes
}
func ConvertByteInt2(b []byte) ([]byte, bool) {

	for i, x := range b {
		if x > 47 && x < 58 {
			b[i] = x - 48
		} else if x > 96 && x < 123 {
			b[i] = x - 87
		} else if x > 64 && x < 91 {
			b[i] = x - 29
		} else if x == 95 {
			b[i] = 62
		} else if x == 45 {
			b[i] = 63
		} else {
			return b, true
		}
	}
	return b, false
}
func ParamBusqueda(b []byte) (Params, bool) {

	var leng int = len(b)
	var p Params

	if leng < 3 {
		return p, true
	}

	var count int = 2
	var j int = 0
	aux := DecodeParams(AlphaValueInt(b[j]))
	aux2 := DecodeParams2(AlphaValueInt(b[j+1]))
	j += 2

	count += aux2[4]*4 + aux2[0]*2 + aux2[3] + aux[1] + aux[2] + aux[3] + aux[0] + 18

	if leng < count {
		return p, true
	}

	var CantEvals int
	if aux2[2] == 1 {
		CantEvals = AlphaValueInt(b[j])
		count += CantEvals + 1
		j++
	}

	CantCuads := AlphaValueInt(b[j])
	p.Cuads = make([][3]byte, CantCuads)
	j += 1

	// EVALS
	if aux2[2] == 1 {
		p.Evals = make([]uint64, CantEvals)
		for i := 0; i < CantEvals; i++ {
			p.Evals[i] = AlphaValue(b[j])
			j++
		}
	}
	// DESDE
	if aux2[4] == 1 {
		p.Desde = AlphaValueFloat(b[j])*262144 + AlphaValueFloat(b[j+1])*4096 + AlphaValueFloat(b[j+2])*64 + AlphaValueFloat(b[j+3])
		j += 4
	}

	// LARGO
	if aux2[3] == 1 {
		p.Largo = AlphaValueInt(b[j])
		if p.Largo > 30 {
			p.Largo = 20
		}
		j += 1
	} else {
		p.Largo = 20
	}

	// DISTANCIA MAX
	if aux2[0] == 1 {
		p.DistanciaMax = AlphaValue(b[j])*64 + AlphaValue(b[j+1])
		j += 2
	}

	// OPCIONES
	v1 := AlphaValue(b[j])
	v2 := AlphaValue(b[j+1])
	v3 := AlphaValue(b[j+2])

	p.Opciones = [4]uint64{v1, v2, v3, v1 + v2 + v3}
	j += 3

	// PAIS CATEGORIA
	var cat int
	p.PaisCat[0] = AlphaValueUint8(b[j])
	j += 1

	if aux[1] == 0 {
		cat = AlphaValueInt(b[j])*64 + AlphaValueInt(b[j+1])
		j += 2
	} else if aux[1] == 1 {
		cat = AlphaValueInt(b[j])*4096 + AlphaValueInt(b[j+1])*64 + AlphaValueInt(b[j+2])
		j += 3
	} else {
		cat = AlphaValueInt(b[j])*262144 + AlphaValueInt(b[j+1])*4096 + AlphaValueInt(b[j+2])*64 + AlphaValueInt(b[j+3])
		j += 4
	}

	copy(p.PaisCat[1:], Int_by_Min3(cat))

	// LAT Y LON
	if aux[3] == 0 {
		p.Lat = (AlphaValueFloat(b[j])*16777216 + AlphaValueFloat(b[j+1])*262144 + AlphaValueFloat(b[j+2])*4096 + AlphaValueFloat(b[j+3])*64 + AlphaValueFloat(b[j+4]) - 1800000000) / 10000000
		j += 5
	} else {
		p.Lat = (AlphaValueFloat(b[j])*1073741824 + AlphaValueFloat(b[j+1])*16777216 + AlphaValueFloat(b[j+2])*262144 + AlphaValueFloat(b[j+3])*4096 + AlphaValueFloat(b[j+4])*64 + AlphaValueFloat(b[j+5]) - 1800000000) / 10000000
		j += 6
	}
	if aux[2] == 0 {
		p.Lon = (AlphaValueFloat(b[j])*16777216 + AlphaValueFloat(b[j+1])*262144 + AlphaValueFloat(b[j+2])*4096 + AlphaValueFloat(b[j+3])*64 + AlphaValueFloat(b[j+4]) - 900000000) / 10000000
		j += 5
	} else {
		p.Lon = (AlphaValueFloat(b[j])*1073741824 + AlphaValueFloat(b[j+1])*16777216 + AlphaValueFloat(b[j+2])*262144 + AlphaValueFloat(b[j+3])*4096 + AlphaValueFloat(b[j+4])*64 + AlphaValueFloat(b[j+5]) - 900000000) / 10000000
		j += 6
	}

	// CUADS
	var cuad int
	var xc int
	for i := 0; i < aux[0]; i++ {

		aux3 := DecodeParamsCuads(AlphaValueInt(b[j]))
		count += (aux3[0] + 1) * (aux3[1] + 1)
		j += 1

		if count < leng {
			for k := 0; k < aux3[0]+1; k++ {

				if aux3[1] == 0 {
					cuad = AlphaValueInt(b[j])
					j += 1
				} else if aux3[1] == 1 {
					cuad = AlphaValueInt(b[j])*64 + AlphaValueInt(b[j+1])
					j += 2
				} else if aux3[1] == 2 {
					cuad = AlphaValueInt(b[j])*4096 + AlphaValueInt(b[j+1])*64 + AlphaValueInt(b[j+2])
					j += 3
				} else {
					cuad = AlphaValueInt(b[j])*262144 + AlphaValueInt(b[j+1])*4096 + AlphaValueInt(b[j+2])*64 + AlphaValueInt(b[j+3])
					j += 4
				}
				if xc < CantCuads {
					p.Cuads[xc] = ParamCuad(cuad)
					xc++
				} else {
					return p, true
				}
			}
		} else {
			return p, true
		}
	}

	if aux2[1] == 1 {

		CantFiltros := AlphaValueInt(b[j]) + 1
		count += CantFiltros * 2
		j += 1

		for i := 0; i < CantFiltros; i++ {

			num := AlphaValueInt(b[j])
			decfiltro := DecodeParamsFiltro(AlphaValueInt(b[j+1]))
			count += (decfiltro[0] + 1) * (decfiltro[1] + 1)
			j += 2

			if count <= leng {
				p.Filtros[num] = make([]uint64, decfiltro[0]+1)
				for m := 0; m < decfiltro[0]+1; m++ {
					if decfiltro[1] == 0 {
						p.Filtros[num][m] = AlphaValue(b[j])
						j += 1
					} else if decfiltro[1] == 1 {
						p.Filtros[num][m] = AlphaValue(b[j])*64 + AlphaValue(b[j+1])
						j += 2
					} else {
						p.Filtros[num][m] = AlphaValue(b[j])*4096 + AlphaValue(b[j+1])*64 + AlphaValue(b[j+2])
						j += 3
					}
				}
			} else {
				return p, true
			}
		}
	}

	if leng != count {
		return p, true
	}
	return p, false
}

func AlphaValue(b byte) uint64 {
	if b > 47 && b < 58 {
		return uint64(b - 48)
	} else if b > 96 && b < 123 {
		return uint64(b - 87)
	} else if b > 64 && b < 91 {
		return uint64(b - 29)
	} else if b == 95 {
		return 62
	} else if b == 45 {
		return 63
	}
	return 0
}
func AlphaValueFloat(b byte) float64 {
	if b > 47 && b < 58 {
		return float64(b - 48)
	} else if b > 96 && b < 123 {
		return float64(b - 87)
	} else if b > 64 && b < 91 {
		return float64(b - 29)
	} else if b == 95 {
		return 62
	} else if b == 45 {
		return 63
	}
	return 0
}
func AlphaValueInt(b byte) int {
	if b > 47 && b < 58 {
		return int(b - 48)
	} else if b > 96 && b < 123 {
		return int(b - 87)
	} else if b > 64 && b < 91 {
		return int(b - 29)
	} else if b == 95 {
		return 62
	} else if b == 45 {
		return 63
	}
	return 0
}
func DecodeParams(num int) []int {

	var res []int = make([]int, 4)
	res[0] = num / 12 // CANT ARR CUADS
	num = num % 12
	res[1] = num / 4 // CANT BYTE CATEGORIA 2-3-4
	num = num % 4
	res[2] = num / 2 // CANT BYTE LAT 5-6
	res[3] = num % 2 // CANT BYTE LON 5-6
	return res
}
func DecodeParams2(num int) []int {

	var res []int = make([]int, 5)
	res[0] = num / 16 // DISTANCIA MAX
	num = num % 16
	res[1] = num / 8 // SI TIENE FILTROS
	num = num % 8
	res[2] = num / 4 // SI TIENE EVAL
	num = num % 4
	res[3] = num / 2 // SI TIENE LARGO
	res[4] = num % 2 // SI TIENE DESDE
	return res
}
func DecodeParamsCuads(num int) [2]int {

	v1 := num / 4 // CANT CUADS
	v2 := num % 4 // CANT BYTE CUADS
	return [2]int{v1, v2}
}
func DecodeParamsFiltro(num int) [2]int {

	v1 := num / 3 // CANT
	v2 := num % 3 // CANT BYTE
	return [2]int{v1, v2}
}
func ParamCuadrantesBytes(b []byte) [][]byte {

	var leng int = len(b)
	res := make([][]byte, leng/4)
	if leng%4 == 0 {
		for i := 0; i < leng/4; i++ {
			res[i] = ParamCategoriaBytes(b[i*4 : (i+1)*4])
		}
	}
	return res
}
func AlphaValueUint8(b byte) uint8 {
	if b > 47 && b < 58 {
		return uint8(b - 48)
	} else if b > 96 && b < 123 {
		return uint8(b - 87)
	} else if b > 64 && b < 91 {
		return uint8(b - 29)
	} else if b == 95 {
		return 62
	} else if b == 45 {
		return 63
	}
	return 0
}
func Int_by_Min3(num int) []byte {

	b := make([]byte, 3)
	var r int = num % 16777216
	b[0] = uint8(r / 65536)
	r = r % 65536
	b[1], b[2] = uint8(r/256), uint8(r%256)
	return b
}
func ParamCategoriaBytes(b []byte) []byte {

	res := make([]byte, 3)
	var aux uint64
	if len(b) == 4 {
		aux = AlphaValue(b[0])*262144 + AlphaValue(b[1])*4096 + AlphaValue(b[2])*64 + AlphaValue(b[3])
		res[0] = uint8(aux / 65536)
		aux = aux % 65536
		res[1], res[2] = uint8(aux/256), uint8(aux%256)
	}
	return res
}
func ParamCuad(n int) [3]byte {

	v1 := uint8(n / 65536)
	n = n % 65536
	v2, v3 := uint8(n/256), uint8(n%256)
	return [3]byte{v1, v2, v3}
}
func ParamBusqueda3(b []byte) (Params, bool) {

	var leng int = len(b)
	var p Params

	if leng < 2 {
		return p, true
	}

	var j int = 0
	aux := DecodeParams(AlphaValueInt(b[j]))
	aux2 := DecodeParams2(AlphaValueInt(b[j+1]))
	j += 2

	var CantEvals int
	if leng > j {
		if aux2[2] == 1 {
			CantEvals = AlphaValueInt(b[j]) + 1
			j++
		}
	} else {
		return p, true
	}

	var CantCuads int
	if leng > j {
		CantCuads = AlphaValueInt(b[j])
		if CantCuads > 14 {
			return p, true
		}
		p.Cuads = make([][3]byte, CantCuads)
		j += 1
	} else {
		return p, true
	}

	// EVALS
	if aux2[2] == 1 {
		p.Evals = make([]uint64, CantEvals)
		for i := 0; i < CantEvals; i++ {
			if leng > j {
				p.Evals[i] = AlphaValue(b[j])
				j++
			} else {
				return p, true
			}
		}
	}

	// DESDE
	if aux2[4] == 1 {
		if leng > j+3 {
			p.Desde = AlphaValueFloat(b[j])*262144 + AlphaValueFloat(b[j+1])*4096 + AlphaValueFloat(b[j+2])*64 + AlphaValueFloat(b[j+3])
			j += 4
		} else {
			return p, true
		}
	}

	// LARGO
	if aux2[3] == 1 {
		if leng > j {
			p.Largo = AlphaValueInt(b[j])
			if p.Largo > 30 {
				p.Largo = 20
			}
			j += 1
		} else {
			return p, true
		}
	} else {
		p.Largo = 20
	}

	// DISTANCIA MAX
	if aux2[0] == 1 {
		if leng > j+1 {
			p.DistanciaMax = AlphaValue(b[j])*64 + AlphaValue(b[j+1])
			j += 2
		} else {
			return p, true
		}
	}

	// OPCIONES
	if leng > j+2 {

		v1 := AlphaValue(b[j])
		v2 := AlphaValue(b[j+1])
		v3 := AlphaValue(b[j+2])
		p.Opciones = [4]uint64{v1, v2, v3, v1 + v2 + v3}
		j += 3

	} else {
		return p, true
	}

	// PAIS CATEGORIA
	var cat int
	if leng > j {
		p.PaisCat[0] = AlphaValueUint8(b[j])
		j += 1
	} else {
		return p, true
	}

	if aux[1] == 0 {
		if leng > j+1 {
			cat = AlphaValueInt(b[j])*64 + AlphaValueInt(b[j+1])
			j += 2
		} else {
			return p, true
		}
	} else if aux[1] == 1 {
		if leng > j+2 {
			cat = AlphaValueInt(b[j])*4096 + AlphaValueInt(b[j+1])*64 + AlphaValueInt(b[j+2])
			j += 3
		} else {
			return p, true
		}
	} else {
		if leng > j+3 {
			cat = AlphaValueInt(b[j])*262144 + AlphaValueInt(b[j+1])*4096 + AlphaValueInt(b[j+2])*64 + AlphaValueInt(b[j+3])
			j += 4
		} else {
			return p, true
		}
	}

	copy(p.PaisCat[1:], Int_by_Min3(cat))

	// LAT Y LON
	if aux[3] == 0 {
		if leng > j+4 {
			p.Lat = (AlphaValueFloat(b[j])*16777216 + AlphaValueFloat(b[j+1])*262144 + AlphaValueFloat(b[j+2])*4096 + AlphaValueFloat(b[j+3])*64 + AlphaValueFloat(b[j+4]) - 1800000000) / 10000000
			j += 5
		} else {
			return p, true
		}
	} else {
		if leng > j+5 {
			p.Lat = (AlphaValueFloat(b[j])*1073741824 + AlphaValueFloat(b[j+1])*16777216 + AlphaValueFloat(b[j+2])*262144 + AlphaValueFloat(b[j+3])*4096 + AlphaValueFloat(b[j+4])*64 + AlphaValueFloat(b[j+5]) - 1800000000) / 10000000
			j += 6
		} else {
			return p, true
		}
	}

	if aux[2] == 0 {
		if leng > j+4 {
			p.Lon = (AlphaValueFloat(b[j])*16777216 + AlphaValueFloat(b[j+1])*262144 + AlphaValueFloat(b[j+2])*4096 + AlphaValueFloat(b[j+3])*64 + AlphaValueFloat(b[j+4]) - 900000000) / 10000000
			j += 5
		} else {
			return p, true
		}
	} else {
		if leng > j+5 {
			p.Lon = (AlphaValueFloat(b[j])*1073741824 + AlphaValueFloat(b[j+1])*16777216 + AlphaValueFloat(b[j+2])*262144 + AlphaValueFloat(b[j+3])*4096 + AlphaValueFloat(b[j+4])*64 + AlphaValueFloat(b[j+5]) - 900000000) / 10000000
			j += 6
		} else {
			return p, true
		}
	}

	// CUADS
	var cuad int
	var xc int
	for i := 0; i < aux[0]+1; i++ {

		if leng > j {

			aux3 := DecodeParamsCuads(AlphaValueInt(b[j]))
			j += 1

			for k := 0; k < aux3[0]+1; k++ {
				if aux3[1] == 0 {
					if leng > j {
						cuad = AlphaValueInt(b[j])
						j += 1
					} else {
						return p, true
					}
				} else if aux3[1] == 1 {
					if leng > j+1 {
						cuad = AlphaValueInt(b[j])*64 + AlphaValueInt(b[j+1])
						j += 2
					} else {
						return p, true
					}
				} else if aux3[1] == 2 {
					if leng > j+2 {
						cuad = AlphaValueInt(b[j])*4096 + AlphaValueInt(b[j+1])*64 + AlphaValueInt(b[j+2])
						j += 3
					} else {
						return p, true
					}
				} else {
					if leng > j+3 {
						cuad = AlphaValueInt(b[j])*262144 + AlphaValueInt(b[j+1])*4096 + AlphaValueInt(b[j+2])*64 + AlphaValueInt(b[j+3])
						j += 4
					} else {
						return p, true
					}
				}
				if xc < CantCuads {
					p.Cuads[xc] = ParamCuad(cuad)
					xc++
				} else {
					return p, true
				}
			}

		} else {
			return p, true
		}
	}

	if aux2[1] == 1 {

		if leng > j {

			CantFiltros := AlphaValueInt(b[j]) + 1
			j += 1

			for i := 0; i < CantFiltros; i++ {

				if leng > j+1 {

					num := AlphaValueInt(b[j])
					decfiltro := DecodeParamsFiltro(AlphaValueInt(b[j+1]))
					j += 2

					p.Filtros[num] = make([]uint64, decfiltro[0]+1)
					for m := 0; m < decfiltro[0]+1; m++ {
						if decfiltro[1] == 0 {
							if leng > j {
								p.Filtros[num][m] = AlphaValue(b[j])
								j += 1
							} else {
								return p, true
							}
						} else if decfiltro[1] == 1 {
							if leng > j+1 {
								p.Filtros[num][m] = AlphaValue(b[j])*64 + AlphaValue(b[j+1])
								j += 2
							} else {
								return p, true
							}
						} else {
							if leng > j+2 {
								p.Filtros[num][m] = AlphaValue(b[j])*4096 + AlphaValue(b[j+1])*64 + AlphaValue(b[j+2])
								j += 3
							} else {
								return p, true
							}
						}
					}

				} else {
					return p, true
				}
			}

		} else {
			return p, true
		}

	}

	if leng == j-1 {
		return p, true
	}
	return p, false
}
func ParamBusqueda4(b []byte) (Params, bool) {

	c, err := ConvertByteInt2(b)
	var p Params

	if !err {

		var leng int = len(b)

		if leng < 2 {
			return p, true
		}

		var j int = 0
		aux := DecodeParams(int(c[j]))
		aux2 := DecodeParams2(int(c[j+1]))
		j += 2

		var CantCuads int
		if leng > j {
			CantCuads = int(c[j])
			if CantCuads > 14 {
				return p, true
			}
			p.Cuads = make([][3]byte, CantCuads)
			j += 1
		} else {
			return p, true
		}

		// EVALS
		if aux2[2] == 1 {
			p.Evals = make([]uint64, int(c[j])+1)
			j++
			for i := 0; i < len(p.Evals); i++ {
				if leng > j {
					p.Evals[i] = uint64(c[j])
					j++
				} else {
					return p, true
				}
			}
		}

		// DESDE
		if aux2[4] == 1 {
			if leng > j+3 {
				p.Desde = float64(c[j])*262144 + float64(c[j+1])*4096 + float64(c[j+2])*64 + float64(c[j+3])
				j += 4
			} else {
				return p, true
			}
		}

		// LARGO
		p.Largo = 20
		if aux2[3] == 1 {
			if leng > j {
				p.Largo = int(c[j])
				if p.Largo > 30 && p.Largo < 20 {
					p.Largo = 20
				}
				j += 1
			} else {
				return p, true
			}
		}

		// DISTANCIA MAX
		if aux2[0] == 1 {
			if leng > j+1 {
				p.DistanciaMax = uint64(c[j])*64 + uint64(c[j+1])
				j += 2
			} else {
				return p, true
			}
		}

		// OPCIONES
		if leng > j+2 {
			p.Opciones = [4]uint64{uint64(c[j]), uint64(c[j+1]), uint64(c[j+2]), uint64(c[j+0]) + uint64(c[j+1]) + uint64(c[j+2])}
			j += 3
		} else {
			return p, true
		}

		// PAIS CATEGORIA
		if leng > j+2+aux[1] {
			p.PaisCat[0] = uint8(c[j])
			if aux[1] == 0 {
				copy(p.PaisCat[1:], Int_by_Min3(int(c[j+1])*64+int(c[j+2])))
			} else if aux[1] == 1 {
				copy(p.PaisCat[1:], Int_by_Min3(int(c[j+1])*4096+int(c[j+2])*64+int(c[j+3])))
			} else {
				copy(p.PaisCat[1:], Int_by_Min3(int(c[j+1])*262144+int(c[j+2])*4096+int(c[j+3])*64+int(c[j+4])))
			}
			j += 3 + aux[1]
		} else {
			return p, true
		}

		if leng > j+4+aux[3]+4+aux[2] {
			if aux[3] == 0 {
				p.Lat = (float64(c[j])*16777216 + float64(c[j+1])*262144 + float64(c[j+2])*4096 + float64(c[j+3])*64 + float64(c[j+4]) - 1800000000) / 10000000
				j += 5
			} else {
				p.Lat = (float64(c[j])*1073741824 + float64(c[j+1])*16777216 + float64(c[j+2])*262144 + float64(c[j+3])*4096 + float64(c[j+4])*64 + float64(c[j+5]) - 1800000000) / 10000000
				j += 6
			}
			if aux[2] == 0 {
				p.Lon = (float64(c[j])*16777216 + float64(c[j+1])*262144 + float64(c[j+2])*4096 + float64(c[j+3])*64 + float64(c[j+4]) - 900000000) / 10000000
				j += 5
			} else {
				p.Lon = (float64(c[j])*1073741824 + float64(c[j+1])*16777216 + float64(c[j+2])*262144 + float64(c[j+3])*4096 + float64(c[j+4])*64 + float64(c[j+5]) - 900000000) / 10000000
				j += 6
			}
		} else {
			return p, true
		}

		var xc int
		for i := 0; i < aux[0]+1; i++ {
			if leng > j {
				aux3 := DecodeParamsCuads(int(c[j]))
				j += 1
				if leng >= j+(aux3[0]+1)*(aux3[1]+1) {
					for k := 0; k < aux3[0]+1; k++ {
						if CantCuads > xc {
							if aux3[1] == 0 {
								p.Cuads[xc] = ParamCuad(int(c[j]))
							} else if aux3[1] == 1 {
								p.Cuads[xc] = ParamCuad(int(c[j])*64 + int(c[j+1]))
							} else if aux3[1] == 2 {
								p.Cuads[xc] = ParamCuad(int(c[j])*4096 + int(c[j+1])*64 + int(c[j+2]))
							} else {
								p.Cuads[xc] = ParamCuad(int(c[j])*262144 + int(c[j+1])*4096 + int(c[j+2])*64 + int(c[j+3]))
							}
							xc++
							j += aux3[1] + 1
						}
					}
				}
			} else {
				return p, true
			}
		}

		if aux2[1] == 1 {
			if leng > j {
				var num int
				var decfiltro [2]int
				CantFiltros := int(c[j]) + 1
				j += 1
				for i := 0; i < CantFiltros; i++ {
					if leng > j+1 {
						num = int(c[j])
						decfiltro = DecodeParamsFiltro(int(c[j+1]))
						j += 2
					}
					if leng >= j+(decfiltro[0]+1)*(decfiltro[1]+1) {
						p.Filtros[num] = make([]uint64, decfiltro[0]+1)
						for m := 0; m < decfiltro[0]+1; m++ {
							if decfiltro[1] == 0 {
								p.Filtros[num][m] = uint64(c[j])
							} else if decfiltro[1] == 1 {
								p.Filtros[num][m] = uint64(c[j])*64 + uint64(c[j+1])
							} else {
								p.Filtros[num][m] = uint64(c[j])*4096 + uint64(c[j+1])*64 + uint64(c[j+2])
							}
							j += decfiltro[1] + 1
						}
					}
				}
			} else {
				return p, true
			}
		}
		if leng == j {
			return p, false
		}
	} else {
		return p, true
	}

	return p, false
}

/*
func ParamBusqueda(b []byte) (utils.Params, bool) {

	//fmt.Println("---------- DECODE PARAMS ----------")

	var debug_min int = -1
	var debug_max int = -1
	var debug bool = false

	var leng int = len(b)
	var p utils.Params

	if leng < 2 {
		fmt.Println("Err1")
		return p, true
	}

	var count int = 0
	var j int = 0
	aux := utils.DecodeParams(utils.AlphaValueInt(b[j]))
	aux2 := utils.DecodeParams2(utils.AlphaValueInt(b[j+1]))
	//fmt.Printf("[0:1] %v%v Aux 1-2\n", string(b[j]), string(b[j+1]))
	j += 2

	if debug_max >= 0 && debug_min <= 0 {
		fmt.Printf("01- Byte(%v) %v %v\n", j-2, b[j-2], string(b[j-2]))
		if debug {
			fmt.Printf("Lat %v - Lon %v - Cat %v - ListCuad %v\n", aux[2], aux[3], aux[1], aux[0])
		}
	}
	if debug_max >= 1 && debug_min <= 1 {
		fmt.Printf("02- Byte(%v) %v %v\n", j-1, b[j-1], string(b[j-1]))
		if debug {
			fmt.Printf("Desde %v - Largo %v - Evals %v - DistMAx %v - Filtro %v\n", aux2[4], aux2[3], aux2[2], aux2[0], aux2[1])
		}
	}

	count += 2
	count += aux[0] + 1
	count += aux2[2] // CANTEVALS
	//fmt.Printf("CantEva %v %v\n", count, aux2[2])
	count += 1 // CANTCUADS
	//fmt.Printf("CantCua %v 1\n", count)
	count += aux2[4] * 4 // DESDE
	//fmt.Printf("Desde   %v %v\n", count, aux2[4]*4)
	count += aux2[3] // LARGO
	//fmt.Printf("Largo   %v %v\n", count, aux2[3])
	count += aux2[0] * 2 // DISTANCIA MAX
	//fmt.Printf("DistMax %v %v\n", count, aux2[0]*2)
	count += 3 // OPCIONES
	//fmt.Printf("Opcione %v 3\n", count)
	count += aux[1] + 3 // PAIS CATEGORIA
	//fmt.Printf("PaisCat %v %v\n", count, aux[1]+3)
	count += aux[2] + 5 // LAITUTD
	//fmt.Printf("Latitud %v %v\n", count, aux[2]+5)
	count += aux[3] + 5 // LONGITUD
	//fmt.Printf("Longitu %v %v\n", count, aux[3]+5)
	count += aux2[1] // CantFiltros
	//fmt.Printf("CantFil %v %v\n", count, aux2[1])

	//
	if leng < count {
		fmt.Println("Err2")
		return p, true
	}

	var CantEvals int
	if aux2[2] == 1 {
		//fmt.Printf("[%v] %v CantEvals\n", j, string(b[j]))
		CantEvals = utils.AlphaValueInt(b[j]) + 1
		if debug_max >= 2 && debug_min <= 2 {
			fmt.Printf("03- Byte(%v) %v %v\n", j, b[j], string(b[j]))
			if debug {
				fmt.Printf("CantEvals %v\n", CantEvals)
			}
		}
		count += CantEvals
		j++
	}

	CantCuads := utils.AlphaValueInt(b[j])
	//fmt.Printf("[%v] %v CantCuads\n", j, string(b[j]))
	if CantCuads > 14 {
		fmt.Println("Err Max CantCuads 14")
		return p, true
	}
	if debug_max >= 3 && debug_min <= 3 {
		fmt.Printf("04- Byte(%v) %v %v\n", j, b[j], string(b[j]))
		if debug {
			fmt.Printf("CantCuads %v\n", CantCuads)
		}
	}
	p.Cuads = make([][3]byte, CantCuads)
	j += 1

	// EVALS
	if aux2[2] == 1 {
		//fmt.Printf("CantEvals (%v): ", CantEvals)
		p.Evals = make([]uint64, CantEvals)
		if debug_max >= 4 && debug_min <= 4 {
			fmt.Printf("05- Byte(%v) - ", j)
		}
		for i := 0; i < CantEvals; i++ {
			if debug_max >= 4 && debug_min <= 4 {
				fmt.Printf("%v", string(b[j]))
				//fmt.Printf("05- Byte(%v): %v %v ", j, b[j], string(b[j]))
				if debug {
					fmt.Printf("Evals: %v %v\n", i, utils.AlphaValue(b[j]))
				}
			}
			p.Evals[i] = utils.AlphaValue(b[j])
			j++
		}
		if debug_max >= 4 && debug_min <= 4 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
	// DESDE
	if aux2[4] == 1 {
		p.Desde = utils.AlphaValueFloat(b[j])*262144 + utils.AlphaValueFloat(b[j+1])*4096 + utils.AlphaValueFloat(b[j+2])*64 + utils.AlphaValueFloat(b[j+3])
		if debug_max >= 5 && debug_min <= 5 {
			fmt.Printf("06- Byte(%v) %v %v\n", j, b[j:j+4], string(b[j:j+4]))
			if debug {
				fmt.Printf("DESDE %v %v\n", b[j:j+4], p.Desde)
			}
		}
		//fmt.Printf("%v DESDE\n", string(b[j:j+4]))
		j += 4
	}

	// LARGO
	if aux2[3] == 1 {
		p.Largo = utils.AlphaValueInt(b[j])
		if p.Largo > 30 {
			p.Largo = 20
		}
		if debug_max >= 6 && debug_min <= 6 {
			fmt.Printf("07- Byte(%v) %v %v\n", j, b[j], string(b[j]))
			if debug {
				fmt.Printf("LARGO %v %v\n", b[j], p.Largo)
			}
		}
		//fmt.Printf("%v Largo\n", string(b[j]))
		j += 1
	} else {
		p.Largo = 20
	}

	// DISTANCIA MAX
	if aux2[0] == 1 {
		p.DistanciaMax = utils.AlphaValue(b[j])*64 + utils.AlphaValue(b[j+1])
		if debug_max >= 7 && debug_min <= 7 {
			fmt.Printf("08- Byte(%v) %v %v\n", j, b[j:j+2], string(b[j:j+2]))
			if debug {
				fmt.Printf("DistanciaMax %v %v\n", b[j:j+2], p.DistanciaMax)
			}
		}
		//fmt.Printf("%v DistMax\n", string(b[j:j+2]))
		j += 2
	}

	// OPCIONES
	v1 := utils.AlphaValue(b[j])
	v2 := utils.AlphaValue(b[j+1])
	v3 := utils.AlphaValue(b[j+2])

	if debug_max >= 8 && debug_min <= 8 {
		fmt.Printf("09- Byte(%v) %v %v\n", j, b[j:j+3], string(b[j:j+3]))
		if debug {
			fmt.Printf("Opciones [%v-%v-%v]\n", b[j], b[j+1], b[j+2])
		}
	}
	p.Opciones = [4]uint64{v1, v2, v3, v1 + v2 + v3}
	j += 3

	// PAIS CATEGORIA
	var cat int
	p.PaisCat[0] = utils.AlphaValueUint8(b[j])
	if debug_max >= 9 && debug_min <= 9 {
		fmt.Printf("10- Byte(%v) %v %v\n", j, b[j], string(b[j]))
		if debug {
			fmt.Printf("Pais %v\n", b[j])
		}
	}
	//fmt.Printf("%v Pais\n", string(b[j]))
	j += 1

	if aux[1] == 0 {
		cat = utils.AlphaValueInt(b[j])*64 + utils.AlphaValueInt(b[j+1])
		if debug_max >= 10 && debug_min <= 10 {
			fmt.Printf("11A- Byte(%v) %v %v\n", j, b[j:j+2], string(b[j:j+2]))
			if debug {
				fmt.Printf("Categoria 2b %v\n", cat)
			}
		}
		//fmt.Printf("%v Categoria2\n", string(b[j:j+2]))
		j += 2
	} else if aux[1] == 1 {
		cat = utils.AlphaValueInt(b[j])*4096 + utils.AlphaValueInt(b[j+1])*64 + utils.AlphaValueInt(b[j+2])
		if debug_max >= 10 && debug_min <= 10 {
			fmt.Printf("11B- Byte(%v) %v %v\n", j, b[j:j+3], string(b[j:j+3]))
			if debug {
				fmt.Printf("Categoria 2b %v\n", cat)
			}
		}
		//fmt.Printf("%v Categoria3\n", string(b[j:j+3]))
		j += 3
	} else {
		cat = utils.AlphaValueInt(b[j])*262144 + utils.AlphaValueInt(b[j+1])*4096 + utils.AlphaValueInt(b[j+2])*64 + utils.AlphaValueInt(b[j+3])
		if debug_max >= 10 && debug_min <= 10 {
			fmt.Printf("11C- Byte(%v) %v %v\n", j, b[j:j+4], string(b[j:j+4]))
			if debug {
				fmt.Printf("Categoria 2b %v\n", cat)
			}
		}
		//fmt.Printf("%v Categoria4\n", string(b[j:j+4]))
		j += 4
	}

	copy(p.PaisCat[1:], utils.Int_by_Min3(cat))

	// LAT Y LON
	if aux[3] == 0 {
		p.Lat = (utils.AlphaValueFloat(b[j])*16777216 + utils.AlphaValueFloat(b[j+1])*262144 + utils.AlphaValueFloat(b[j+2])*4096 + utils.AlphaValueFloat(b[j+3])*64 + utils.AlphaValueFloat(b[j+4]) - 1800000000) / 10000000
		if debug_max >= 11 && debug_min <= 11 {
			fmt.Printf("12A- Byte(%v) %v %v\n", j, b[j:j+5], string(b[j:j+5]))
			if debug {
				fmt.Printf("Latitud %v\n", p.Lat)
			}
		}
		//fmt.Printf("%v Lat5\n", string(b[j:j+5]))
		j += 5
	} else {
		p.Lat = (utils.AlphaValueFloat(b[j])*1073741824 + utils.AlphaValueFloat(b[j+1])*16777216 + utils.AlphaValueFloat(b[j+2])*262144 + utils.AlphaValueFloat(b[j+3])*4096 + utils.AlphaValueFloat(b[j+4])*64 + utils.AlphaValueFloat(b[j+5]) - 1800000000) / 10000000
		if debug_max >= 11 && debug_min <= 11 {
			fmt.Printf("12B- Byte(%v) %v %v\n", j, b[j:j+6], string(b[j:j+6]))
			if debug {
				fmt.Printf("Latitud %v\n", p.Lat)
			}
		}
		//fmt.Printf("%v Lat6\n", string(b[j:j+6]))
		j += 6
	}

	if aux[2] == 0 {
		p.Lon = (utils.AlphaValueFloat(b[j])*16777216 + utils.AlphaValueFloat(b[j+1])*262144 + utils.AlphaValueFloat(b[j+2])*4096 + utils.AlphaValueFloat(b[j+3])*64 + utils.AlphaValueFloat(b[j+4]) - 900000000) / 10000000
		if debug_max >= 12 && debug_min <= 12 {
			fmt.Printf("13A- Byte(%v) %v %v\n", j, b[j:j+5], string(b[j:j+5]))
			if debug {
				fmt.Printf("Latitud %v\n", p.Lon)
			}
		}
		//fmt.Printf("%v Lon5\n", string(b[j:j+5]))
		j += 5
	} else {
		p.Lon = (utils.AlphaValueFloat(b[j])*1073741824 + utils.AlphaValueFloat(b[j+1])*16777216 + utils.AlphaValueFloat(b[j+2])*262144 + utils.AlphaValueFloat(b[j+3])*4096 + utils.AlphaValueFloat(b[j+4])*64 + utils.AlphaValueFloat(b[j+5]) - 900000000) / 10000000
		if debug_max >= 12 && debug_min <= 12 {
			fmt.Printf("13B- Byte(%v) %v %v\n", j, b[j:j+6], string(b[j:j+6]))
			if debug {
				fmt.Printf("Latitud %v\n", p.Lon)
			}
		}
		//fmt.Printf("%v Lon6\n", string(b[j:j+6]))
		j += 6
	}

	// CUADS
	var cuad int
	var xc int

	for i := 0; i < aux[0]+1; i++ {

		aux3 := utils.DecodeParamsCuads(utils.AlphaValueInt(b[j]))
		if debug_max >= 13 && debug_min <= 13 {
			fmt.Printf("14- Byte(%v) %v %v\n", j, b[j], string(b[j]))
			if debug {
				fmt.Printf("DecodeParamsCuads CantCuad %v Cantbyte %v\n", aux3[0], aux3[1])
			}
		}
		j += 1

		count += (aux3[0] + 1) * (aux3[1] + 1)
		//fmt.Printf("Cuads %v %v %v\n", i, count, (aux3[0]+1)*(aux3[1]+1))

		if debug_max >= 14 && debug_min <= 14 {
			fmt.Printf("15A- Byte(%v) ", j)
		}
		if count <= leng {
			for k := 0; k < aux3[0]+1; k++ {

				if aux3[1] == 0 {
					cuad = utils.AlphaValueInt(b[j])
					if debug_max >= 14 && debug_min <= 14 {
						//fmt.Printf("15A- Byte(%v) %v %v\n", j, b[j], string(b[j]))
						fmt.Printf("%v ", string(b[j]))
						if debug {

						}
					}
					j += 1
				} else if aux3[1] == 1 {
					cuad = utils.AlphaValueInt(b[j])*64 + utils.AlphaValueInt(b[j+1])
					if debug_max >= 14 && debug_min <= 14 {
						//fmt.Printf("15B- Byte(%v) %v %v\n", j, b[j:j+2], string(b[j:j+2]))
						fmt.Printf("%v ", string(b[j:j+2]))
						if debug {

						}
					}
					j += 2
				} else if aux3[1] == 2 {
					cuad = utils.AlphaValueInt(b[j])*4096 + utils.AlphaValueInt(b[j+1])*64 + utils.AlphaValueInt(b[j+2])
					if debug_max >= 14 && debug_min <= 14 {
						//fmt.Printf("15C- Byte(%v) %v %v\n", j, b[j:j+3], string(b[j:j+3]))
						fmt.Printf("%v ", string(b[j:j+3]))
						if debug {

						}
					}
					j += 3
				} else {
					cuad = utils.AlphaValueInt(b[j])*262144 + utils.AlphaValueInt(b[j+1])*4096 + utils.AlphaValueInt(b[j+2])*64 + utils.AlphaValueInt(b[j+3])
					if debug_max >= 14 && debug_min <= 14 {
						//fmt.Printf("15D- Byte(%v) %v %v\n", j, b[j:j+4], string(b[j:j+4]))
						fmt.Printf("%v ", string(b[j:j+4]))
						if debug {

						}
					}
					j += 4
				}

				if xc < CantCuads {
					p.Cuads[xc] = utils.ParamCuad(cuad)
					xc++
				} else {
					fmt.Println("Err3")
					return p, true
				}
			}
		} else {
			fmt.Println("Count-Length-J:", count, leng, j)
			fmt.Println("Err4")
			return p, true
		}
		if debug_max >= 14 && debug_min <= 14 {
			fmt.Printf("\n")
		}
	}

	if aux2[1] == 1 {

		CantFiltros := utils.AlphaValueInt(b[j]) + 1
		if debug_max >= 15 && debug_min <= 15 {
			fmt.Printf("16- Byte(%v) %v %v\n", j, b[j], string(b[j]))
			if debug {

			}
		}
		j += 1
		count += (CantFiltros) * 2

		if leng < count {
			fmt.Println("Err Max Filtros")
			return p, true
		}

		for i := 0; i < CantFiltros; i++ {
			num := utils.AlphaValueInt(b[j])
			decfiltro := utils.DecodeParamsFiltro(utils.AlphaValueInt(b[j+1]))

			if debug_max >= 16 && debug_min <= 16 {
				fmt.Printf("17- Byte(%v) %v %v\n", j, b[j], string(b[j]))
				if debug {
					fmt.Printf("CantFiltros Num %v\n", num)
				}
			}
			j += 1

			if debug_max >= 17 && debug_min <= 17 {
				fmt.Printf("18- Byte(%v) %v %v\n", j, b[j+1], string(b[j+1]))
				if debug {
					fmt.Printf("DecodeParamsFiltro Cant %v Bytes %v\n", decfiltro[0], decfiltro[1])
				}
			}
			j += 1
			count += (decfiltro[0] + 1) * (decfiltro[1] + 1)

			if count <= leng {

				p.Filtros[num] = make([]uint64, decfiltro[0]+1)
				for m := 0; m < decfiltro[0]+1; m++ {
					if decfiltro[1] == 0 {
						p.Filtros[num][m] = utils.AlphaValue(b[j])
						if debug_max >= 18 && debug_min <= 18 {
							fmt.Printf("19A- Byte(%v) %v %v\n", j, b[j], string(b[j]))
							if debug {

							}
						}
						j += 1
					} else if decfiltro[1] == 1 {
						p.Filtros[num][m] = utils.AlphaValue(b[j])*64 + utils.AlphaValue(b[j+1])
						if debug_max >= 18 && debug_min <= 18 {
							fmt.Printf("19B- Byte(%v) %v %v\n", j, b[j:j+2], string(b[j:j+2]))
							if debug {

							}
						}
						j += 2
					} else {
						p.Filtros[num][m] = utils.AlphaValue(b[j])*4096 + utils.AlphaValue(b[j+1])*64 + utils.AlphaValue(b[j+2])
						if debug_max >= 19 && debug_min <= 19 {
							fmt.Printf("19B- Byte(%v) %v %v\n", j, b[j:j+3], string(b[j:j+3]))
							if debug {

							}
						}
						j += 3
					}
				}
			} else {
				fmt.Println("Error 8")
				fmt.Println("Count-Length-J:", count, leng, j)
				return p, true
			}
		}
	}

	if leng != count {
		return p, true
	}
	return p, false
}
func ParamBusqueda2(b []byte) (utils.Params, bool) {

	var leng int = len(b)
	var p utils.Params

	if leng < 2 {
		return p, true
	}

	var count int = 0
	var j int = 0
	aux := utils.DecodeParams(utils.AlphaValueInt(b[j]))
	aux2 := utils.DecodeParams2(utils.AlphaValueInt(b[j+1]))
	j += 2

	count += (aux2[0] * 2) + (aux2[4] * 4) + aux2[1] + aux[1] + aux[2] + aux[3] + aux2[2] + aux[0] + aux2[3] + 20

	if leng < count {
		return p, true
	}

	var CantEvals int
	if aux2[2] == 1 {
		CantEvals = utils.AlphaValueInt(b[j]) + 1
		count += CantEvals
		j++
	}

	CantCuads := utils.AlphaValueInt(b[j])
	if CantCuads > 14 {
		return p, true
	}
	p.Cuads = make([][3]byte, CantCuads)
	j += 1

	// EVALS
	if aux2[2] == 1 {
		p.Evals = make([]uint64, CantEvals)
		for i := 0; i < CantEvals; i++ {
			p.Evals[i] = utils.AlphaValue(b[j])
			j++
		}
	}

	// DESDE
	if aux2[4] == 1 {
		p.Desde = utils.AlphaValueFloat(b[j])*262144 + utils.AlphaValueFloat(b[j+1])*4096 + utils.AlphaValueFloat(b[j+2])*64 + utils.AlphaValueFloat(b[j+3])
		j += 4
	}

	// LARGO
	if aux2[3] == 1 {
		p.Largo = utils.AlphaValueInt(b[j])
		if p.Largo > 30 {
			p.Largo = 20
		}
		j += 1
	} else {
		p.Largo = 20
	}

	// DISTANCIA MAX
	if aux2[0] == 1 {
		p.DistanciaMax = utils.AlphaValue(b[j])*64 + utils.AlphaValue(b[j+1])
		j += 2
	}

	// OPCIONES
	v1 := utils.AlphaValue(b[j])
	v2 := utils.AlphaValue(b[j+1])
	v3 := utils.AlphaValue(b[j+2])

	p.Opciones = [4]uint64{v1, v2, v3, v1 + v2 + v3}
	j += 3

	// PAIS CATEGORIA
	var cat int
	p.PaisCat[0] = utils.AlphaValueUint8(b[j])
	j += 1

	if aux[1] == 0 {
		cat = utils.AlphaValueInt(b[j])*64 + utils.AlphaValueInt(b[j+1])
		j += 2
	} else if aux[1] == 1 {
		cat = utils.AlphaValueInt(b[j])*4096 + utils.AlphaValueInt(b[j+1])*64 + utils.AlphaValueInt(b[j+2])
		j += 3
	} else {
		cat = utils.AlphaValueInt(b[j])*262144 + utils.AlphaValueInt(b[j+1])*4096 + utils.AlphaValueInt(b[j+2])*64 + utils.AlphaValueInt(b[j+3])
		j += 4
	}

	copy(p.PaisCat[1:], utils.Int_by_Min3(cat))

	// LAT Y LON
	if aux[3] == 0 {
		p.Lat = (utils.AlphaValueFloat(b[j])*16777216 + utils.AlphaValueFloat(b[j+1])*262144 + utils.AlphaValueFloat(b[j+2])*4096 + utils.AlphaValueFloat(b[j+3])*64 + utils.AlphaValueFloat(b[j+4]) - 1800000000) / 10000000
		j += 5
	} else {
		p.Lat = (utils.AlphaValueFloat(b[j])*1073741824 + utils.AlphaValueFloat(b[j+1])*16777216 + utils.AlphaValueFloat(b[j+2])*262144 + utils.AlphaValueFloat(b[j+3])*4096 + utils.AlphaValueFloat(b[j+4])*64 + utils.AlphaValueFloat(b[j+5]) - 1800000000) / 10000000
		j += 6
	}

	if aux[2] == 0 {
		p.Lon = (utils.AlphaValueFloat(b[j])*16777216 + utils.AlphaValueFloat(b[j+1])*262144 + utils.AlphaValueFloat(b[j+2])*4096 + utils.AlphaValueFloat(b[j+3])*64 + utils.AlphaValueFloat(b[j+4]) - 900000000) / 10000000
		j += 5
	} else {
		p.Lon = (utils.AlphaValueFloat(b[j])*1073741824 + utils.AlphaValueFloat(b[j+1])*16777216 + utils.AlphaValueFloat(b[j+2])*262144 + utils.AlphaValueFloat(b[j+3])*4096 + utils.AlphaValueFloat(b[j+4])*64 + utils.AlphaValueFloat(b[j+5]) - 900000000) / 10000000
		j += 6
	}

	// CUADS
	var cuad int
	var xc int
	for i := 0; i < aux[0]+1; i++ {

		aux3 := utils.DecodeParamsCuads(utils.AlphaValueInt(b[j]))
		count += (aux3[0] + 1) * (aux3[1] + 1)
		j += 1

		if count <= leng {
			for k := 0; k < aux3[0]+1; k++ {
				if aux3[1] == 0 {
					cuad = utils.AlphaValueInt(b[j])
					j += 1
				} else if aux3[1] == 1 {
					cuad = utils.AlphaValueInt(b[j])*64 + utils.AlphaValueInt(b[j+1])
					j += 2
				} else if aux3[1] == 2 {
					cuad = utils.AlphaValueInt(b[j])*4096 + utils.AlphaValueInt(b[j+1])*64 + utils.AlphaValueInt(b[j+2])
					j += 3
				} else {
					cuad = utils.AlphaValueInt(b[j])*262144 + utils.AlphaValueInt(b[j+1])*4096 + utils.AlphaValueInt(b[j+2])*64 + utils.AlphaValueInt(b[j+3])
					j += 4
				}
				if xc < CantCuads {
					p.Cuads[xc] = utils.ParamCuad(cuad)
					xc++
				} else {
					return p, true
				}
			}
		} else {
			return p, true
		}
	}

	if aux2[1] == 1 {

		CantFiltros := utils.AlphaValueInt(b[j]) + 1
		count += (CantFiltros) * 2
		j += 1

		if leng < count {
			return p, true
		}

		for i := 0; i < CantFiltros; i++ {

			num := utils.AlphaValueInt(b[j])
			decfiltro := utils.DecodeParamsFiltro(utils.AlphaValueInt(b[j+1]))
			count += (decfiltro[0] + 1) * (decfiltro[1] + 1)
			j += 2

			if count <= leng {
				p.Filtros[num] = make([]uint64, decfiltro[0]+1)
				for m := 0; m < decfiltro[0]+1; m++ {
					if decfiltro[1] == 0 {
						p.Filtros[num][m] = utils.AlphaValue(b[j])
						j += 1
					} else if decfiltro[1] == 1 {
						p.Filtros[num][m] = utils.AlphaValue(b[j])*64 + utils.AlphaValue(b[j+1])
						j += 2
					} else {
						p.Filtros[num][m] = utils.AlphaValue(b[j])*4096 + utils.AlphaValue(b[j+1])*64 + utils.AlphaValue(b[j+2])
						j += 3
					}
				}
			} else {
				return p, true
			}
		}
	}

	if leng != count {
		return p, true
	}
	return p, false
}
func ParamBusqueda5(b []byte) (utils.Params, bool) {

	var leng int = len(b)
	var p utils.Params

	if leng < 2 {
		return p, true
	}

	var count int = 0
	var j int = 0
	aux := utils.DecodeParams(utils.AlphaValueInt(b[j]))
	aux2 := utils.DecodeParams2(utils.AlphaValueInt(b[j+1]))
	j += 2

	count += (aux2[0] * 2) + (aux2[4] * 4) + aux2[1] + aux[1] + aux[2] + aux[3] + aux2[2] + aux[0] + aux2[3] + 20

	if leng < count {
		return p, true
	}

	var CantCuads int
	if leng > j {
		CantCuads = utils.AlphaValueInt(b[j])
		if CantCuads > 14 {
			return p, true
		}
		p.Cuads = make([][3]byte, CantCuads)
		j += 1
	} else {
		return p, true
	}

	// EVALS
	if aux2[2] == 1 {
		p.Evals = make([]uint64, utils.AlphaValueInt(b[j])+1)
		j++
		for i := 0; i < len(p.Evals); i++ {
			if leng > j {
				p.Evals[i] = utils.AlphaValue(b[j])
				j++
			} else {
				return p, true
			}
		}
	}

	// DESDE
	if aux2[4] == 1 {
		if leng > j+3 {
			p.Desde = utils.AlphaValueFloat(b[j])*262144 + utils.AlphaValueFloat(b[j+1])*4096 + utils.AlphaValueFloat(b[j+2])*64 + utils.AlphaValueFloat(b[j+3])
			j += 4
		} else {
			return p, true
		}
	}

	// LARGO
	if aux2[3] == 1 {
		if leng > j {
			p.Largo = utils.AlphaValueInt(b[j])
			if p.Largo > 30 {
				p.Largo = 20
			}
			j += 1
		} else {
			return p, true
		}
	} else {
		p.Largo = 20
	}

	// DISTANCIA MAX
	if aux2[0] == 1 {
		if leng > j+1 {
			p.DistanciaMax = utils.AlphaValue(b[j])*64 + utils.AlphaValue(b[j+1])
			j += 2
		} else {
			return p, true
		}
	}

	// OPCIONES
	if leng > j+2 {

		v1 := utils.AlphaValue(b[j])
		v2 := utils.AlphaValue(b[j+1])
		v3 := utils.AlphaValue(b[j+2])
		p.Opciones = [4]uint64{v1, v2, v3, v1 + v2 + v3}
		j += 3

	} else {
		return p, true
	}

	// PAIS CATEGORIA
	var cat int
	if leng > j {
		p.PaisCat[0] = utils.AlphaValueUint8(b[j])
		j += 1
	} else {
		return p, true
	}

	if aux[1] == 0 {
		if leng > j+1 {
			cat = utils.AlphaValueInt(b[j])*64 + utils.AlphaValueInt(b[j+1])
			j += 2
		} else {
			return p, true
		}
	} else if aux[1] == 1 {
		if leng > j+2 {
			cat = utils.AlphaValueInt(b[j])*4096 + utils.AlphaValueInt(b[j+1])*64 + utils.AlphaValueInt(b[j+2])
			j += 3
		} else {
			return p, true
		}
	} else {
		if leng > j+3 {
			cat = utils.AlphaValueInt(b[j])*262144 + utils.AlphaValueInt(b[j+1])*4096 + utils.AlphaValueInt(b[j+2])*64 + utils.AlphaValueInt(b[j+3])
			j += 4
		} else {
			return p, true
		}
	}

	copy(p.PaisCat[1:], utils.Int_by_Min3(cat))

	// LAT Y LON
	if aux[3] == 0 {
		if leng > j+4 {
			p.Lat = (utils.AlphaValueFloat(b[j])*16777216 + utils.AlphaValueFloat(b[j+1])*262144 + utils.AlphaValueFloat(b[j+2])*4096 + utils.AlphaValueFloat(b[j+3])*64 + utils.AlphaValueFloat(b[j+4]) - 1800000000) / 10000000
			j += 5
		} else {
			return p, true
		}
	} else {
		if leng > j+5 {
			p.Lat = (utils.AlphaValueFloat(b[j])*1073741824 + utils.AlphaValueFloat(b[j+1])*16777216 + utils.AlphaValueFloat(b[j+2])*262144 + utils.AlphaValueFloat(b[j+3])*4096 + utils.AlphaValueFloat(b[j+4])*64 + utils.AlphaValueFloat(b[j+5]) - 1800000000) / 10000000
			j += 6
		} else {
			return p, true
		}
	}

	if aux[2] == 0 {
		if leng > j+4 {
			p.Lon = (utils.AlphaValueFloat(b[j])*16777216 + utils.AlphaValueFloat(b[j+1])*262144 + utils.AlphaValueFloat(b[j+2])*4096 + utils.AlphaValueFloat(b[j+3])*64 + utils.AlphaValueFloat(b[j+4]) - 900000000) / 10000000
			j += 5
		} else {
			return p, true
		}
	} else {
		if leng > j+5 {
			p.Lon = (utils.AlphaValueFloat(b[j])*1073741824 + utils.AlphaValueFloat(b[j+1])*16777216 + utils.AlphaValueFloat(b[j+2])*262144 + utils.AlphaValueFloat(b[j+3])*4096 + utils.AlphaValueFloat(b[j+4])*64 + utils.AlphaValueFloat(b[j+5]) - 900000000) / 10000000
			j += 6
		} else {
			return p, true
		}
	}

	// CUADS
	var cuad int
	var xc int
	for i := 0; i < aux[0]+1; i++ {

		if leng > j {

			aux3 := utils.DecodeParamsCuads(utils.AlphaValueInt(b[j]))
			count += (aux3[0] + 1) * (aux3[1] + 1)
			j += 1

			for k := 0; k < aux3[0]+1; k++ {
				if aux3[1] == 0 {
					if leng > j {
						cuad = utils.AlphaValueInt(b[j])
						j += 1
					} else {
						return p, true
					}
				} else if aux3[1] == 1 {
					if leng > j+1 {
						cuad = utils.AlphaValueInt(b[j])*64 + utils.AlphaValueInt(b[j+1])
						j += 2
					} else {
						return p, true
					}
				} else if aux3[1] == 2 {
					if leng > j+2 {
						cuad = utils.AlphaValueInt(b[j])*4096 + utils.AlphaValueInt(b[j+1])*64 + utils.AlphaValueInt(b[j+2])
						j += 3
					} else {
						return p, true
					}
				} else {
					if leng > j+3 {
						cuad = utils.AlphaValueInt(b[j])*262144 + utils.AlphaValueInt(b[j+1])*4096 + utils.AlphaValueInt(b[j+2])*64 + utils.AlphaValueInt(b[j+3])
						j += 4
					} else {
						return p, true
					}
				}
				if xc < CantCuads {
					p.Cuads[xc] = utils.ParamCuad(cuad)
					xc++
				} else {
					return p, true
				}
			}

		} else {
			return p, true
		}
	}

	if aux2[1] == 1 {

		if leng > j {

			CantFiltros := utils.AlphaValueInt(b[j]) + 1
			count += (CantFiltros) * 2
			j += 1

			for i := 0; i < CantFiltros; i++ {

				if leng > j+1 {

					num := utils.AlphaValueInt(b[j])
					decfiltro := utils.DecodeParamsFiltro(utils.AlphaValueInt(b[j+1]))
					count += (decfiltro[0] + 1) * (decfiltro[1] + 1)
					j += 2

					p.Filtros[num] = make([]uint64, decfiltro[0]+1)
					for m := 0; m < decfiltro[0]+1; m++ {
						if decfiltro[1] == 0 {
							if leng > j {
								p.Filtros[num][m] = utils.AlphaValue(b[j])
								j += 1
							} else {
								return p, true
							}
						} else if decfiltro[1] == 1 {
							if leng > j+1 {
								p.Filtros[num][m] = utils.AlphaValue(b[j])*64 + utils.AlphaValue(b[j+1])
								j += 2
							} else {
								return p, true
							}
						} else {
							if leng > j+2 {
								p.Filtros[num][m] = utils.AlphaValue(b[j])*4096 + utils.AlphaValue(b[j+1])*64 + utils.AlphaValue(b[j+2])
								j += 3
							} else {
								return p, true
							}
						}
					}

				} else {
					return p, true
				}
			}

		} else {
			return p, true
		}

	}

	if leng == j-1 {
		return p, true
	}
	return p, false
}
*/
func Uint64_To_7Bytes(num uint64) []byte {
	var b []byte = make([]byte, 7)
	copy(b[0:], Uint32_To_3Bytes(num/4294967296))
	copy(b[3:], Uint32_To_4Bytes(num%4294967296))
	return b
}
func Uint32_To_3Bytes(num uint64) []byte {
	b := make([]byte, 3)
	var r uint64 = num % 16777216
	b[0] = uint8(r / 65536)
	r = r % 65536
	b[1], b[2] = uint8(r/256), uint8(r%256)
	return b
}
func Uint32_To_4Bytes(num uint64) []byte {
	b := make([]byte, 4)
	var r uint64 = num % 16777216
	b[0] = uint8(num / 16777216)
	b[1] = uint8(r / 65536)
	r = r % 65536
	b[2], b[3] = uint8(r/256), uint8(r%256)
	return b
}
