package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cmdrkeene/booking"
)

var flagDataPath = flag.String("data", "/tmp/booking.db", "Path to data file")
var flagHttp = flag.String("http", ":3000", "HTTP address e.g. :3000")
var flagProcessorToken = flag.String("processor-token", "test", "Payment processor token")

func main() {
	flag.Parse()

	service := booking.NewService(
		*flagDataPath,
		*flagProcessorToken,
	)

	fmt.Println("listening on", *flagHttp)
	app := booking.NewWebApp(service)
	log.Fatal(http.ListenAndServe(*flagHttp, app))
}
