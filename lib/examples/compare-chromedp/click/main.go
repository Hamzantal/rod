package main

import (
	"github.com/ysmood/rod"
	"log"
	"time"
)

// This example demonstrates how to use a selector to click on an element.
func main() {
	browser := rod.New().Timeout(15 * time.Second).Trace(true).Connect()
	defer browser.Close()

	page := browser.Page("https://golang.org/pkg/time/")
	// Element will wait till an element with the selector is found.
	page.Element(`body > footer`)
	// Click will expand the dropdown menu for the example.
	page.Element(`#pkg-examples > div`).Click()
	// Text will extract the example's content.
	example := page.Element(`#example_After .play .input textarea`).Text()

	log.Printf("Go's time.After example:\n%s", example)
}
