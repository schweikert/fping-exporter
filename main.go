package main

import (
	"log"
	"net/http"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var opts struct {
	Listen string `short:"l" long:"listen" description:"Listen address" value-name:"[HOST]:PORT" default:":8080"`
	Period uint   `short:"p" long:"period" description:"Period in seconds, should match Prometheus scrape interval" value-name:"SECS" default:"60"`
	Fping  string `short:"f" long:"fping"  description:"Fping binary path" value-name:"PATH" default:"/usr/bin/fping"`
	Count  uint   `short:"c" long:"count"  description:"Number of pings to send at each period" value-name:"N" default:"20"`
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	targetParam := r.URL.Query().Get("target")
	if targetParam == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html>
		    <head><title>Fping Exporter</title></head>
			<body>
			<b>ERROR: missing target parameter</b>
			</body>`))
		return
	}

	target := GetTarget(
		WorkerSpec{
			period: time.Second * 60,
		},
		TargetSpec{
			host: targetParam,
		},
	)

	h := promhttp.HandlerFor(target.registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(0)
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", probeHandler)
	log.Fatal(http.ListenAndServe(opts.Listen, nil))
}
