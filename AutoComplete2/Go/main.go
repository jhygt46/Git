package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"unicode/utf8"
	"utils/utils"

	"github.com/valyala/fasthttp"

	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
)

type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Db        *ledis.DB      `json:"Db"`
	Conf      Config         `json:"Conf"`
	DDoS      utils.DDoS     `json:"DDoS"`
	Auto      map[uint8]Auto `json:"Auto"`
	CountMem  uint32         `json:"CountMem"`
	CountDisk uint32         `json:"CountDisk"`
}

type Auto struct {
	Auto map[string][]byte `json:"Auto"`
}

func main() {

	var port string
	var dbname string
	if runtime.GOOS == "windows" {
		port = ":83"
		dbname = "C:/Go/LedisDB/Auto"
	} else {
		port = ":8083"
		dbname = "/var/Go/LedisDB/Auto"
	}

	pass := &MyHandler{
		DDoS: utils.DDoS{Start: false, Ips: &utils.IPs{Ip: make(map[uint32]uint8, 0)}, BlackList: make([]uint32, 0)},
		Auto: make(map[uint8]Auto, 0),
		Db:   LedisConfig(dbname),
	}

	//pass.SaveMemoryDb()

	con := context.Background()
	con, cancel := context.WithCancel(con)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					pass.Conf.init()
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-con.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()
	go func() {
		fasthttp.ListenAndServe(port, pass.HandleFastHTTP)
	}()
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	//ctx.Response.Header.Set("Content-Type", "application/octet-stream")
	ctx.Response.Header.Set("Content-Type", "application/json")

	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/a":
			// SEARCH 1 CARACTER
			if !h.DDoS.Start || utils.VerificarIp(&h.DDoS, utils.Ip_str_u32(ctx.RemoteAddr().String())) {
				var pais uint8 = utils.ParamUint8(ctx.QueryArgs().Peek("p"))
				var search []byte = ctx.QueryArgs().Peek("s")
				if Auto, Found := h.Auto[pais]; Found {
					if Auto2, Found2 := Auto.Auto[string(search)]; Found2 {
						h.CountMem++
						ctx.SetBody(Auto2)
					} else {
						val, _ := h.Db.Get(utils.KeySearch(search))
						ctx.SetBody(val)
					}
				} else {
					fmt.Println("PAIS NO ENCONTRADO")
					ctx.SetBody([]byte{})
				}
			} else {
				Send(utils.SendParamPostJson(), []byte{})
				ctx.SetBody([]byte{})
			}
		case "/al":
			// SEARCH MAX 3 CARACTER
			if !h.DDoS.Start || utils.VerificarIp(&h.DDoS, utils.Ip_str_u32(ctx.RemoteAddr().String())) {
				var pais uint8 = utils.ParamUint8(ctx.QueryArgs().Peek("p"))
				var search []byte = ctx.QueryArgs().Peek("s")
				var leng int = utils.LengMax3(ctx.QueryArgs().Peek("l"))
				if leng > 0 && len(search) > leng {
					if Auto, Found := h.Auto[pais]; Found {
						for i := 0; i < leng; i++ {
							s := utils.Unicode(search, leng-i)
							if Auto2, Found2 := Auto.Auto[string(s)]; Found2 {
								ctx.Write(Auto2)
							} else {
								val, _ := h.Db.Get(utils.KeySearch(s))
								ctx.Write(val)
							}
						}
					} else {
						fmt.Println("PAIS NO ENCONTRADO")
						ctx.SetBody([]byte{})
					}
				} else {
					// WAF
				}
			} else {
				Send(utils.SendParamPostJson(), []byte{})
				ctx.SetBody([]byte{})
			}
		case "/ac":

		case "/bl":
			ctx.SetBody(h.BlackList())
		case "/unicode":

			n0 := 0
			n1 := 1114111

			var min1ba uint8 = 255
			var max1ba uint8 = 0

			var min2ba uint8 = 255
			var max2ba uint8 = 0
			var min2bb uint8 = 255
			var max2bb uint8 = 0

			var min3ba uint8 = 255
			var max3ba uint8 = 0
			var min3bb uint8 = 255
			var max3bb uint8 = 0
			var min3bc uint8 = 255
			var max3bc uint8 = 0

			var min4ba uint8 = 255
			var max4ba uint8 = 0
			var min4bb uint8 = 255
			var max4bb uint8 = 0
			var min4bc uint8 = 255
			var max4bc uint8 = 0
			var min4bd uint8 = 255
			var max4bd uint8 = 0

			for i := n0; i < n1; i++ {
				r := rune(int32(i))
				buf := make([]byte, 4)
				n := utf8.EncodeRune(buf, r)
				if n == 1 {
					if buf[0] < min1ba {
						min1ba = buf[0]
					}
					if buf[0] > max1ba {
						max1ba = buf[0]
					}
				}
				if n == 2 {
					if buf[0] < min2ba {
						min2ba = buf[0]
					}
					if buf[0] > max2ba {
						max2ba = buf[0]
					}
					if buf[1] < min2bb {
						min2bb = buf[1]
					}
					if buf[1] > max2bb {
						max2bb = buf[1]
					}
				}
				if n == 3 {
					if buf[0] < min3ba {
						min3ba = buf[0]
					}
					if buf[0] > max3ba {
						max3ba = buf[0]
					}
					if buf[1] < min3bb {
						min3bb = buf[1]
					}
					if buf[1] > max3bb {
						max3bb = buf[1]
					}
					if buf[2] < min3bc {
						min3bc = buf[2]
					}
					if buf[2] > max3bc {
						max3bc = buf[2]
					}
				}
				if n == 4 {
					if buf[0] < min4ba {
						min4ba = buf[0]
					}
					if buf[0] > max4ba {
						max4ba = buf[0]
					}
					if buf[1] < min4bb {
						min4bb = buf[1]
					}
					if buf[1] > max4bb {
						max4bb = buf[1]
					}
					if buf[2] < min4bc {
						min4bc = buf[2]
					}
					if buf[2] > max4bc {
						max4bc = buf[2]
					}
					if buf[3] < min4bd {
						min4bd = buf[3]
					}
					if buf[3] > max4bd {
						max4bd = buf[3]
					}
				}
			}
			fmt.Fprintf(ctx, "1 Byte MinA(%v) - MaxA(%v)\n", min1ba, max1ba)
			fmt.Fprintf(ctx, "2 Byte MinA(%v) - MaxA(%v) / MinB(%v) - MaxB(%v)\n", min2ba, max2ba, min2bb, max2bb)
			fmt.Fprintf(ctx, "3 Byte MinA(%v) - MaxA(%v) / MinB(%v) - MaxB(%v) / MinC(%v) - MaxC(%v)\n", min3ba, max3ba, min3bb, max3bb, min3bc, max3bc)
			fmt.Fprintf(ctx, "4 Byte MinA(%v) - MaxA(%v) / MinB(%v) - MaxB(%v) / MinC(%v) - MaxC(%v) / MinD(%v) - MaxD(%v)\n", min4ba, max4ba, min4bb, max4bb, min4bc, max4bc, min4bd, max4bd)
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}

func (h *MyHandler) BlackList() []byte {

	len := len(h.DDoS.BlackList)
	if len > 0 {
		var i int = 0
		resp := make([]byte, len*4)
		for _, ip := range h.DDoS.BlackList {
			i += copy(resp[i:], utils.Int32_by_Min4(ip))
		}
		h.DDoS.BlackList = make([]uint32, 0)
		return resp
	} else {
		return []byte{}
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 10 * time.Second
	fmt.Println("DAEMON")
}
func (c *Config) init() {
	var tick = flag.Duration("tick", 1*time.Second, "Ticking interval")
	c.Tiempo = *tick
}
func run(con context.Context, c *MyHandler, stdout io.Writer) error {
	c.Conf.init()
	log.SetOutput(os.Stdout)
	for {
		select {
		case <-con.Done():
			return nil
		case <-time.Tick(c.Conf.Tiempo):
			c.StartDaemon()
		}
	}
}

// DDoS //
func Send(ops utils.SendData, data []byte) []byte {

	uri := fmt.Sprintf("%v:%v%v", ops.Host, ops.Port, ops.Uri)

	req := fasthttp.AcquireRequest()
	req.SetBody(data)
	req.Header.SetMethod(ops.Method)
	req.Header.SetContentType(ops.ContentType)
	req.SetRequestURI(uri)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		fmt.Println(err)
		return []byte{}
	}
	fasthttp.ReleaseRequest(req)
	body := res.Body()
	fasthttp.ReleaseResponse(res)
	return body
}
func LedisConfig(path string) *ledis.DB {
	cfg := lediscfg.NewConfigDefault()
	cfg.DataDir = path
	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)
	return db
}

// SAVE MEMORY
func (h *MyHandler) SaveMemoryDb() {

	v1 := make([][]byte, 0)
	v1 = [][]byte{[]byte{97}, []byte{98}, []byte{99}, []byte{100}, []byte{101}, []byte{102}, []byte{103}, []byte{104}, []byte{105}, []byte{106}, []byte{107}, []byte{108}, []byte{109}, []byte{110}, []byte{195, 177}, []byte{111}, []byte{112}, []byte{113}, []byte{114}, []byte{115}, []byte{116}, []byte{117}, []byte{118}, []byte{119}, []byte{120}, []byte{121}, []byte{122}}
	fmt.Println(len(v1))

	for x := 0; x <= 60; x++ {
		h.Auto[uint8(x)] = Auto{Auto: make(map[string][]byte, 0)}
	}

	var z1 int = 0
	for x := 0; x <= 60; x++ {
		for i := 0; i < len(v1); i++ {
			for j := 0; j < len(v1); j++ {
				h.Auto[uint8(x)].Auto[b2(v1[i][0], v1[j][0])] = GetBytes(450)
				z1++
			}
		}
	}
	fmt.Println("Nivel1:", z1)

	var z2 int = 0
	for x := 0; x <= 60; x++ {
		for i := 0; i < len(v1); i++ {
			for j := 0; j < len(v1); j++ {
				for k := 0; k < len(v1); k++ {
					h.Auto[uint8(x)].Auto[b3(v1[i][0], v1[j][0], v1[k][0])] = GetBytes(450)
					z2++
				}
			}
		}
	}
	fmt.Println("Nivel2:", z2)
}

// TEST DELETE //
func GetBytes(n int) []byte {
	by := make([]byte, n)
	for i := 0; i < n; i++ {
		by[i] = 49
	}
	return by
}
func b2(b1 byte, b2 byte) string {
	b := make([]byte, 2)
	b[0] = b1
	b[1] = b2
	return string(b)
}
func b3(b1 byte, b2 byte, b3 byte) string {
	b := make([]byte, 3)
	b[0] = b1
	b[1] = b2
	b[2] = b3
	return string(b)
}
func b4(b1 byte, b2 byte, b3 byte, b4 byte) string {
	b := make([]byte, 4)
	b[0] = b1
	b[1] = b2
	b[2] = b3
	b[3] = b4
	return string(b)
}
func b5(b1 byte, b2 byte, b3 byte, b4 byte, b5 byte) string {
	b := make([]byte, 5)
	b[0] = b1
	b[1] = b2
	b[2] = b3
	b[3] = b4
	b[4] = b4
	return string(b)
}
