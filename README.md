# urlbox
The Urlbox Go library provides easy access to the Urlbox website screenshot API from your Go application.

Now there's no need to muck around with http clients, etc...

Just initialise the UrlboxClient and make a screenshot of a URL in seconds.

# Documentation
See the [Urlbox API Docs](https://www.urlbox.io/docs/overview)

# Installation
This package can in installed using the go command below.
```sh
go get github.com/iqquee/urlbox
```
# Quick start
```sh
# assume the following codes in example.go file
$ touch example.go
# open the just created example.go file in the text editor of your choice
```

# Screenshot
The Screenshot method takes the screenshot of a website synchronously. 
Inother words, whenever you make a request using this method, you wait to get the screenshotted data([]byte) from the server.

This method takes in the Request{} struct as a parameter.
### Below is a list of all the fields available in the Request{} object.
```go
type Request struct {
		Url     string  `json:"url"`     // url of website to screenshot
		Format  string  `json:"fornmat"` // screenshot file format
		Options Options // optional params for the request
	}
	Options struct {
		FullPage        bool // for full page screenshot
		Width           int
		BlockingOptions Blocking // options for blocking or dismissing certain page elements, such as cookie banners
		SelectorOption  Selector // selector parameter
		ImageOption     Image    // options relating to the outputted PNG, WebP or JPEG file
		WaitOption      Wait
	}
	Blocking struct {
		BlockAds          bool `json:"block_ads"`           // remove ads from page
		HideCookieBanners bool `json:"hide_cookie_banners"` // remove cookie banners if any
		ClickAccept       bool `json:"click_accept"`        // click accept buttons to dismiss pop-upsSelector
	}
	Selector struct {
		Selector              string `json:"selector"`                 // for css selectors e.g #playground for id of playground
		FailIfSelectorMissing bool   `json:"fail_if_selector_missing"` // fail the request when the selector is not found
	}

	Image struct {
		Retina  bool `json:"retina"`  // take a 'retina' or high-definition screenshot, equivalent to setting a device pixel ratio of 2.0 or @2x. Please note that retina screenshots will be double the normal dimensions and will normally take slightly longer to process due to the much bigger image size.
		Quality int  `json:"quality"` // the image quality of the resulting screenshot (JPEG/WebP only)
	}

	Wait struct {
		Delay   int // the amount of time to wait before Urlbox takes the screenshot or PDF, in milliseconds.
		TimeOut int // the amount of time to wait for the requested URL to respond, in milliseconds.
	}
```

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/iqquee/urlbox"
)

func main() {
	apiKey := ""
	secreteKey := ""
	client := urlbox.New(*http.DefaultClient, apiKey, secreteKey)

	request := urlbox.Request{
		Url:    "https://urlbox.io",
		Format: urlbox.FileFormatPng,
		Options: urlbox.Options{
			ImageOption: urlbox.Image{
				Retina:  true,
				Quality: 100,
			},
		},
	}

	data, err := client.Screenshot(request)
	if err != nil {
		fmt.Println("an error occured: ", err)
		return
	}

	fmt.Println("This is the response byte: ", data)

	filename := fmt.Sprintf("%s-%d.%s", strings.Replace(request.Url, "/", "-", -1), time.Now().UTC().Unix(), request.Format)


	if err := os.WriteFile(filename, data, 0666); err != nil {
		fmt.Println("error writing to disk: ", err)
		return
	}
}
```

# ScreenshotAsync
ScreenshotAsync allow your application to receive information when a screenshot has been rendered.
This allows you to render screenshots asynchronously.

This method takes in the RequestAsync{} struct as a parameter.
### Use this object payload to implement the ScreenshotAsync() method
```go
type RequestAsync struct {
    Url        string `json:"url"`         // url of website to screenshot
    WebhookUrl string `json:"webhook_url"` // Pass a webhook URL in as the webhook_url option and Urlbox will send a POST request back to that URL with data about the screenshot in JSON format once it has completed rendering
}
```

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/iqquee/urlbox"
)

func main() {
	apiKey := "LFZu26GbSWgr0qgn"
	secreteKey := "cd3a89b805c740bdb62bef748bdcebfd"
	client := urlbox.New(*http.DefaultClient, apiKey, secreteKey)

	request := urlbox.RequestAsync{
		Url:        "https://urlbox.io",
		WebhookUrl: "https://example.com/webhooks/urlbox",
	}

	data, err := client.ScreenshotAsync(request)
	if err != nil {
		fmt.Println("an error occured: ", err)
		return
	}

	fmt.Println("This is the response: ", data)

}
```
