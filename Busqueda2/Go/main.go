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
	Db        *ledis.DB          `json:"Db"`
	Conf      Config             `json:"Conf"`
	DDoS      utils.DDoS         `json:"DDoS"`
	Busqueda  map[uint8]Busqueda `json:"Busqueda"`
	CountMem  uint32             `json:"CountMem"`
	CountDisk uint32             `json:"CountDisk"`
}
type Busqueda struct {
	Res map[uint64][]byte `json:"Res"`
}
type Params struct {
	Lat   float64 `json:"Lat"`
	Lon   float64 `json:"Lng"`
	Desde float64 `json:"Desde"`
}

func main() {

	var port string
	var dbname string
	if runtime.GOOS == "windows" {
		port = ":82"
		dbname = "C:/Go/LedisDB/Base"
	} else {
		port = ":8082"
		dbname = "/var/Go/LedisDB/Base"
	}

	pass := &MyHandler{
		DDoS:     utils.DDoS{Start: false, Ips: &utils.IPs{Ip: make(map[uint32]uint8, 0)}, BlackList: make([]uint32, 0)},
		Busqueda: make(map[uint8]Busqueda, 0),
		Db:       LedisConfig(dbname),
	}

	pass.SaveMemoryDb()

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

	ctx.Response.Header.Set("Content-Type", "application/octet-stream")

	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/":

			if !h.DDoS.Start || utils.VerificarIp(&h.DDoS, utils.Ip_str_u32(ctx.RemoteAddr().String())) {

				var pais uint8 = utils.ParamUint8(ctx.QueryArgs().Peek("p"))
				if Busqueda, Found := h.Busqueda[pais]; Found {

					var p Params
					var err bool

					p.Lat, err = utils.ParamAlphaLat(ctx.QueryArgs().Peek("lat"))
					p.Lon, err = utils.ParamAlphaLon(ctx.QueryArgs().Peek("lng"))
					p.Desde, err = utils.ParamDesde(ctx.QueryArgs().Peek("desde"))

					CkeckErr(err)

					if Res, Found := Busqueda.Res[5687]; Found {
						ctx.SetBody(Res)
					} else {

					}

					/*
						var key1 uint64 = 0
						var key2 []byte = []byte{0}
						if Res, Found2 := Busqueda.Res[key1]; Found2 {
							h.CountMem++
							ctx.SetBody(Res)
						} else {
							val, _ := h.Db.Get(key2)
							ctx.SetBody(val)
						}
					*/

				} else {

				}

			} else {
				Send(utils.SendParamPostJson(), []byte{49})
				ctx.SetBody([]byte{})
			}

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

/*
func (h *MyHandler) DecodeBytes(Res *Respuesta, bytes []byte, P NewParams) {

}
*/

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
func CkeckErr(b bool) {

}
func CkeckRes(b []byte) {

}

// SAVE MEMORY
func (h *MyHandler) SaveMemoryDb() {

	for x := 0; x <= 60; x++ {
		h.Busqueda[uint8(x)] = Busqueda{Res: make(map[uint64][]byte, 0)}
	}

	var z1 int = 0
	for x := 0; x <= 60; x++ {
		for i := 0; i < 10000; i++ {
			h.Busqueda[uint8(x)].Res[uint64(i)] = GetBytes(450)
			z1++
		}
	}
	fmt.Println("Nivel1:", z1)

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
