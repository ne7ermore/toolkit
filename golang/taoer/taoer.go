package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/ne7ermore/taoer/t"
)

var (
	httpU          = flag.String("url", "", "url")
	method         = flag.String("method", "get", "method")
	form           = flag.String("form", "", "http form data")
	duration       = flag.Int("duration", 60, "times of req")
	qps            = flag.Int("qps", 100, "Qps")
	disKeepA       = flag.Bool("disKA", true, "DisableKeepAlives, if true, prevents re-use of TCP connections between different HTTP requests")
	disCompression = flag.Bool("disComp", true, "DisableCompression, if true, prevents the Transport from requesting compression with an Accept-Encoding: gzip")
	timeout        = flag.Int("timeout", 100, "time for handshake")
	cpus           = flag.Int("cpus", runtime.GOMAXPROCS(-1), "number of CPUs")
	bodyShow       = flag.Bool("bodyShow", true, "show response body")
)

func getHttpsClient() *http.Client {
	var httpClient http.Client
	return &httpClient
}

func main() {
	flag.Parse()

	if *httpU == "" || *qps <= 0 || *cpus <= 0 || *duration <= 0 || *qps%t.InnerQps != 0 {
		t.PrintHelp()
		os.Exit(1)
	}

	httpMethod := strings.ToUpper(*method)
	switch httpMethod {
	case "GET", "POST", "PUT":
	default:
		fmt.Printf("Url method error\n")
		t.PrintHelp()
		os.Exit(1)
	}

	if httpMethod != "GET" && *form == "" {
		t.PrintHelp()
		os.Exit(1)
	}

	runtime.GOMAXPROCS(*cpus)

	t.SetupTao(t.GetClient(*disKeepA,
		*disCompression,
		*timeout),
		*duration,
		*qps,
		&t.TaoR{
			Method: httpMethod,
			Url:    *httpU,
			Form:   *form,
		},
		*bodyShow).
		Run()
}
