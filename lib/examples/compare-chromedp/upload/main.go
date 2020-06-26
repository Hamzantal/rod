package main

import (
	"flag"
	"fmt"
	"github.com/ysmood/rod"
	"github.com/ysmood/rod/lib/launcher"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var flagPort = flag.Int("port", 8544, "port")

// This example demonstrates how to upload a file on a form.
func main() {
	flag.Parse()

	// get wd
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filepath := wd + "/main.go"

	// get some info about the file
	fi, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// start upload server
	result := make(chan int, 1)
	go uploadServer(fmt.Sprintf(":%d", *flagPort), result)

	url := launcher.New().Headless(false).Launch()
	browser := rod.New().ControlURL(url).Connect()

	page := browser.Page(fmt.Sprintf("http://localhost:%d", *flagPort))

	page.Element(`input[name="upload"]`).Input(filepath)
	page.Element(`input[name="submit"]`).Click()

	page.Element("#result").Text()

	log.Printf("original size: %d, upload size: %d", fi.Size(), <-result)
}

func uploadServer(addr string, result chan int) error {
	// create http server and result channel
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, uploadHTML)
	})
	mux.HandleFunc("/upload", func(res http.ResponseWriter, req *http.Request) {
		f, _, err := req.FormFile("upload")
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(res, resultHTML, len(buf))

		result <- len(buf)
	})
	return http.ListenAndServe(addr, mux)
}

const (
	uploadHTML = `<!doctype html>
<html>
<body>
  <form method="POST" action="/upload" enctype="multipart/form-data">
    <input name="upload" type="file"/>
    <input name="submit" type="submit"/>
  </form>
</body>
</html>`

	resultHTML = `<!doctype html>
<html>
<body>
  <div id="result">%d</div>
</body>
</html>`
)
