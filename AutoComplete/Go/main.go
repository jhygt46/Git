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
	"utils/utils"

	"github.com/valyala/fasthttp"

	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
)

type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Db        *ledis.DB         `json:"Db"`
	Conf      Config            `json:"Conf"`
	DDoS      utils.DDoS        `json:"DDoS"`
	Auto      map[string][]byte `json:"Auto"`
	CountMem  uint32            `json:"CountMem"`
	CountDisk uint32            `json:"CountDisk"`
}

func main() {

	var port string
	var dbname string
	if runtime.GOOS == "windows" {
		port = ":82"
		dbname = "C:/Go/LedisDB/Auto"
	} else {
		port = ":8082"
		dbname = "/var/Go/LedisDB/Auto"
	}

	pass := &MyHandler{
		DDoS: utils.DDoS{Start: true, Ips: &utils.IPs{Ip: make(map[uint32]uint8, 0)}, BlackList: make([]uint32, 0)},
		Auto: make(map[string][]byte, 0),
		Db:   LedisConfig(dbname),
	}

	v1 := make([][]byte, 0)
	v1 = [][]byte{[]byte{97}, []byte{98}, []byte{99}, []byte{100}, []byte{101}, []byte{102}, []byte{103}, []byte{104}, []byte{105}, []byte{106}, []byte{107}, []byte{108}, []byte{109}, []byte{110}, []byte{195, 177}, []byte{111}, []byte{112}, []byte{113}, []byte{114}, []byte{115}, []byte{116}, []byte{117}, []byte{118}, []byte{119}, []byte{120}, []byte{121}, []byte{122}}

	var z1, z2, z3 int = 0, 0, 0
	by := []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 195, 177, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 114, 115, 116, 117, 118, 119, 120, 121, 122, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 114, 115, 116, 117, 118, 119, 120, 121, 122, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 114, 115, 116, 117, 118, 119, 120, 121, 122, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57}

	for i := 0; i < len(v1); i++ {
		pass.Auto[string(v1[i])] = by
		z1++
	}

	for i := 0; i < len(v1); i++ {
		for j := 0; j < len(v1); j++ {
			key := append(v1[i], v1[j]...)
			pass.Auto[string(key)] = by
			z2++
		}
	}

	for i := 0; i < len(v1); i++ {
		for j := 0; j < len(v1); j++ {
			key := append(v1[i], v1[j]...)
			for k := 0; k < len(v1); k++ {
				key = append(key, v1[k]...)
				pass.Auto[string(key)] = by
				z3++
			}
		}
	}

	fmt.Println(z1, z2, z3)

	//pass.DDoS.BlackList = append(pass.DDoS.BlackList, 825307441)
	//pass.DDoS.BlackList = append(pass.DDoS.BlackList, 825307442)

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

	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/a":

			if !h.DDoS.Start || utils.VerificarIp(&h.DDoS, utils.Ip_str_u32(ctx.RemoteAddr().String())) {

				var search []byte = ctx.QueryArgs().Peek("s")
				var cuad []byte = ctx.QueryArgs().Peek("c")

				if Auto, Found := h.Auto[string(search)]; Found {
					h.CountMem++
					ctx.SetBody(Auto)
				} else {
					h.CountDisk++
					key1 := utils.KeySearch(search)
					val1, _ := h.Db.Get(key1)
					if len(val1) > 0 {
						ctx.SetBody(val1)
					} else {
						key2 := utils.KeySearchCuad(search, cuad)
						val2, _ := h.Db.Get(key2)
						if len(val2) > 0 {
							ctx.SetBody(val2)
						}
					}
				}

			} else {
				Send(utils.SendParamPostJson(), []byte{})
				ctx.SetBody([]byte{})
			}

		case "/al":

			ctx.SetBody([]byte{49, 50})
			ctx.SetBody([]byte{51, 52})

			/*
				if leng == 0 {

				} else if leng < lensearch {

					for i := 0; i < leng; i++ {

					}

				} else {
					// WAF
				}
			*/

			//fmt.Println("LENG:", leng)
			//fmt.Println("SEARCH:", lensearch)
			//fmt.Println("LEN BYTES:", ctx.QueryArgs().Peek("len"))
			/*
				id := utils.ParamUint32(ctx.QueryArgs().Peek("i"))
				if Filtro, Found := h.Filtros[id]; Found {
					h.CountMem++
					ctx.SetBody(Filtro)
				} else {
					h.CountDisk++
					val, _ := h.Db.Get(utils.Int32_by(id))
					if len(val) > 0 {
						ctx.SetBody(val)
					}
				}
			*/

		case "/bl":

			var i int = 0
			resp := make([]byte, len(h.DDoS.BlackList)*4)
			for _, ip := range h.DDoS.BlackList {
				i += copy(resp[i:], utils.Int32_by_Min4(ip))
			}
			h.DDoS.BlackList = make([]uint32, 0)
			ctx.SetBody(resp)

		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
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
