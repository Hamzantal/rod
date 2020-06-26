# Rod comparison with chromedp

chromedp is one of the most popular libraries available for Go which controls the Chrome Devtools Protocol. 

We have emulated the examples they provide using our own library so you can compare which one you prefer.

Occasionally, some of these examples may break if the Rod API gets updated or if the specific websites these examples use get updated.
We suggest you create an [issue](https://github.com/ysmood/rod/issues/new/choose).

You can build and run these examples in the usual Go way:

```sh
# retrieve examples
$ go get -u -d github.com/ysmood/rod/

# run example <prog>
$ go run $GOPATH/src/github.com/ysmood/rod/lib/examples/compare-chromedp/<prog>/main.go

# build example <prog>
$ go build -o <prog> github.com/ysmood/rod/lib/examples/compare-chromedp/<prog>/ && ./<prog>
```

| Example                   | chromedp Example                                              | Description                                                                  |
|---------------------------|---------------------------------------------------------------|------------------------------------------------------------------------------|
| [click](/lib/examples/compare-chromedp/click)           | [click](https://github.com/chromedp/examples/click)           | use a selector to click on an element                                        |
| [cookie](/lib/examples/compare-chromedp/cookie)         | [cookie](https://github.com/chromedp/examples/cookie)         | set a HTTP cookie on requests                                                |
| [emulate](/lib/examples/compare-chromedp/emulate)       | [emulate](https://github.com/chromedp/examples/emulate)       | emulate a specific device such as an iPhone                                  |
| [eval](/lib/examples/compare-chromedp/eval)             | [eval](https://github.com/chromedp/examples/eval)             | evaluate javascript and retrieve the result                                  |
| [headers](/lib/examples/compare-chromedp/headers)       | [headers](https://github.com/chromedp/examples/headers)       | set a HTTP header on requests                                                |
| [keys](/lib/examples/compare-chromedp/keys)             | [keys](https://github.com/chromedp/examples/keys)             | send key events to an element                                                |
| [logic](/lib/examples/compare-chromedp/logic)           | [logic](https://github.com/chromedp/examples/logic)           | more complex logic beyond simple actions                                     |
| [remote](/lib/examples/compare-chromedp/remote)         | [remote](https://github.com/chromedp/examples/remote)         | connect to an existing Chrome DevTools instance using a remote WebSocket URL |
| [screenshot](/lib/examples/compare-chromedp/screenshot) | [screenshot](https://github.com/chromedp/examples/screenshot) | take a screenshot of a specific element and of the entire browser viewport   |
| [submit](/lib/examples/compare-chromedp/submit)         | [submit](https://github.com/chromedp/examples/submit)         | fill out and submit a form                                                   |
| [text](/lib/examples/compare-chromedp/text)             | [text](https://github.com/chromedp/examples/text)             | extract text from a specific element                                         |
| [upload](/lib/examples/compare-chromedp/upload)         | [upload](https://github.com/chromedp/examples/upload)         | upload a file on a form                                                      |
| [visible](/lib/examples/compare-chromedp/visible)       | [visible](https://github.com/chromedp/examples/visible)       | wait until an element is visible                                             |
