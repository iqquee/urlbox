package urlbox

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUrlRequired          = errors.New("a url must be passed in")
	ErrImageQualityExceeded = errors.New("image quality cannot be greater than 100")
)

type (
	Request struct {
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
		DownloadOption  Download // pass in a filename which sets the content-disposition header on the response. E.g. download=myfilename.png This will make the Urlbox link downloadable, and will prompt the user to save the file as myfilename.png
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

	Download struct {
		DownloadFile bool
		FileName     string `json:"download"`
	}

	Wait struct {
		Delay   int // the amount of time to wait before Urlbox takes the screenshot or PDF, in milliseconds.
		TimeOut int // the amount of time to wait for the requested URL to respond, in milliseconds.
	}
)

// parse() sets up default values if the user doesn't pass any params in
func (r Request) parse() Request {
	// if a file format is not provided, set file format as png
	if r.Format == "" {
		png := FileFormatPng
		r.Format = png
	}
	// if FullPage Options field is not passed set to false
	if !r.Options.FullPage {
		fullPage := false
		r.Options.FullPage = fullPage
	}
	// if Width Options field is not passed set to DefaultWidth
	if r.Options.Width == 0 {
		width := DefaultWidth
		r.Options.Width = width
	}
	// if BlockAds Options field is not passed set to true
	if !r.Options.BlockingOptions.BlockAds {
		blockAds := true
		r.Options.BlockingOptions.BlockAds = blockAds
	}
	// if HideCookieBanners Options field is not passed set to true
	if !r.Options.BlockingOptions.HideCookieBanners {
		cookie := true
		r.Options.BlockingOptions.HideCookieBanners = cookie
	}
	// if ClickAccept Options field is not passed set to true
	if !r.Options.BlockingOptions.ClickAccept {
		accept := true
		r.Options.BlockingOptions.ClickAccept = accept
	}
	// by default FailIfSelectorMissing should be false. Even if the selector is not found, it should not return any error
	if !r.Options.SelectorOption.FailIfSelectorMissing {
		failSelector := false
		r.Options.SelectorOption.FailIfSelectorMissing = failSelector
	}
	// by default the Retina is set to false
	if !r.Options.ImageOption.Retina {
		retina := false
		r.Options.ImageOption.Retina = retina
	}
	// by default the Quality is set to 80
	if r.Options.ImageOption.Quality == 0 {
		quality := 80
		r.Options.ImageOption.Quality = quality
	}
	// by default the Delay is set to 0
	if r.Options.WaitOption.Delay == 0 {
		delay := 0
		r.Options.WaitOption.Delay = delay
	}
	// by default the TimeOut is set to 30000 in milliseconds(3 seconds)
	if r.Options.WaitOption.TimeOut == 0 {
		timeOut := 30000
		r.Options.WaitOption.TimeOut = timeOut
	}

	return r
}

func (c *Client) Screenshot(rq Request) ([]byte, error) {
	// the function shouldnt run if there was no url provided
	if rq.Url == "" {
		return nil, ErrUrlRequired
	}
	// check if the Image quality is not above 100
	if rq.Options.ImageOption.Quality > 100 {
		return nil, ErrImageQualityExceeded
	}

	r := rq.parse()

	// setup the url
	url := fmt.Sprintf("%s/%v?url=%s&width=%v&full_page=%v&block_ads=%v&hide_cookie_banners=%v&click_accept=%v&retina=%v&quality=%v&delay=%v&timeout=%v",
		c.ApiKey, r.Format, r.Url, r.Options.Width, r.Options.FullPage, r.Options.BlockingOptions.BlockAds,
		r.Options.BlockingOptions.HideCookieBanners, r.Options.BlockingOptions.ClickAccept,
		r.Options.ImageOption.Retina, r.Options.ImageOption.Quality, r.Options.WaitOption.Delay, r.Options.WaitOption.TimeOut,
	)

	res, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) ScreenshotAsync(rq Request) ([]byte, error) {
	// the function shouldnt run if there was no url provided
	if rq.Url == "" {
		return nil, ErrUrlRequired
	}
	// check if the Image quality is not above 100
	if rq.Options.ImageOption.Quality > 100 {
		return nil, ErrImageQualityExceeded
	}

	r := rq.parse()

	var downloadFileName string
	if r.Options.DownloadOption.DownloadFile {
		fileName := fmt.Sprintf("%v.%v", r.Options.DownloadOption.FileName, r.Format)
		downloadFileName = fileName
	}
	r.Options.DownloadOption.FileName = downloadFileName

	url := "render/sync"

	res, err := c.newRequest(http.MethodPost, url, r)
	if err != nil {
		return nil, err
	}

	return res, nil
}
