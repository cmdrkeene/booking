package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cmdrkeene/booking"
)

var flagDb = flag.String("db", "", "Path to data file")
var flagHttp = flag.String("http", ":3000", "HTTP address e.g. :3000")

func main() {
	flag.Parse()

	if *flagDb == "" {
		file, err := ioutil.TempFile(os.TempDir(), "booking")
		if err != nil {
			panic(err)
		}
		defer os.Remove(file.Name())
		*flagDb = file.Name()
	}

	app := booking.NewApplication(*flagDb)
	defer app.Close()

	if *flagHttp != "" {
		server := app.NewServer(*flagHttp)
		fmt.Println("listening on", *flagHttp)
		log.Fatal(server.ListenAndServe())
	}
}
