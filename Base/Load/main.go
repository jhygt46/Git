package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {
	// Define los servidores de destino
	destinations := []*url.URL{
		{
			Scheme: "http",
			Host:   "localhost:8081",
		},
		{
			Scheme: "http",
			Host:   "localhost:8082",
		},
		{
			Scheme: "http",
			Host:   "localhost:8083",
		},
	}

	// Crea el balanceador de carga
	balancer := &roundRobinBalancer{
		destinations: destinations,
	}

	// Crea un servidor Fasthttp que utilizará el balanceador de carga
	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			balancer.proxyRequest(ctx)
		},
	}

	// Inicia el servidor en el puerto 9090
	if err := server.ListenAndServe(":8081"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %s", err)
	}
}

// roundRobinBalancer es un balanceador de carga que utiliza el algoritmo Round Robin
type roundRobinBalancer struct {
	destinations []*url.URL
	currentIndex uint32
}

// proxyRequest redirige una solicitud HTTP a uno de los servidores de destino utilizando el algoritmo Round Robin
func (b *roundRobinBalancer) proxyRequest(ctx *fasthttp.RequestCtx) {
	// Obtiene el índice del servidor de destino actual y lo incrementa para la próxima solicitud
	index := int(atomic.AddUint32(&b.currentIndex, 1) % uint32(len(b.destinations)))

	// Obtiene el servidor de destino actual
	dest := b.destinations[index]

	// Crea una nueva solicitud HTTP utilizando la URL del servidor de destino actual
	req := &fasthttp.Request{}
	req.SetRequestURI(dest.String())
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.Header.SetContentLength(ctx.Request.Header.ContentLength())
	req.SetBody(ctx.Request.Body())

	// Envía la solicitud HTTP al servidor de destino actual y devuelve la respuesta
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return net.DialTimeout("tcp", addr, 5*time.Second)
		},
	}

	if err := client.Do(req, resp); err != nil {
		ctx.SetStatusCode(http.StatusBadGateway)
		fmt.Fprintf(ctx, "Error al procesar la solicitud: %s", err)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetContentLength(resp.Header.ContentLength())
	ctx.SetBody(resp.Body())
}
