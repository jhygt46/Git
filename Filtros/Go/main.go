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
	Filtros   map[uint32][]byte `json:"Filtros"`
	CountMem  uint32            `json:"CountMem"`
	CountDisk uint32            `json:"CountDisk"`
}

func main() {

	var port string
	if runtime.GOOS == "windows" {
		port = ":81"
	} else {
		port = ":8080"
	}

	pass := &MyHandler{
		DDoS:    utils.DDoS{Start: true, Ips: &utils.IPs{Ip: make(map[uint32]uint8, 0)}, BlackList: make([]uint32, 0)},
		Filtros: make(map[uint32][]byte, 0),
		Db:      LedisConfig(1),
	}

	pass.DDoS.BlackList = append(pass.DDoS.BlackList, 825307441)

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

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 10 * time.Second
	fmt.Println(utils.PrintMemUsage())
	fmt.Println(utils.GetMonitoringsCpu())
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

	uri := fmt.Sprintf("%v:%v", ops.Host, ops.Port)

	req := fasthttp.AcquireRequest()
	req.SetBody(data)
	req.Header.SetMethod(ops.Method)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(uri)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		panic("handle error")
	}
	fasthttp.ReleaseRequest(req)
	body := res.Body()
	fasthttp.ReleaseResponse(res)
	return body
}
func LedisConfig(path int) *ledis.DB {
	cfg := lediscfg.NewConfigDefault()
	cfg.DataDir = "C:/Go/LedisDB/Init"
	//cfg.DataDir = fmt.Sprintf("/var/Go/LedisDB/filtro/%v", path)
	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)
	return db
}
