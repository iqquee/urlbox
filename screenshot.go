package urlbox

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrUrlRequired          = errors.New("a url must be passed in")
	ErrImageQualityExceeded = errors.New("image quality cannot be greater than 100")
)

type (
	Request struct {
		Url     string  // url of website to screenshot
		Format  *string // screenshot file format
		Options Options // optional params for the request
	}
	Options struct {
		FullPage        *bool // for full page screenshot
		Width           *string
		BlockingOptions Blocking // options for blocking or dismissing certain page elements, such as cookie banners
		SelectorOption  Selector // selector parameter
		ImageOption     Image    // options relating to the outputted PNG, WebP or JPEG file
		Download        Download // pass in a filename which sets the content-disposition header on the response. E.g. download=myfilename.png This will make the Urlbox link downloadable, and will prompt the user to save the file as myfilename.png
		WaitOption      Wait
	}
	Blocking struct {
		BlockAds          *bool // remove ads from page
		HideCookieBanners *bool // remove cookie banners if any
		ClickAccept       *bool // click accept buttons to dismiss pop-upsSelector
	}
	Selector struct {
		Selector              *string // for css selectors e.g #playground for id of playground
		FailIfSelectorMissing *bool   // fail the request when the selector is not found
	}

	Image struct {
		Retina  *bool // take a 'retina' or high-definition screenshot, equivalent to setting a device pixel ratio of 2.0 or @2x. Please note that retina screenshots will be double the normal dimensions and will normally take slightly longer to process due to the much bigger image size.
		Quality int   // the image quality of the resulting screenshot (JPEG/WebP only)
	}

	Wait struct {
		Delay   *int // the amount of time to wait before Urlbox takes the screenshot or PDF, in milliseconds.
		TimeOut *int // the amount of time to wait for the requested URL to respond, in milliseconds.
	}

	Download struct {
		DownloadFile bool
		FileName     *string
	}

	Response struct {
		File io.ReadCloser `json:"file"`
	}
)

// parse() sets up default values if the user doesnt pass any params in
func (r Request) parse() Request {
	// if a file format is not provided, set file format as png
	if r.Format == nil {
		png := FileFormatPng
		r.Format = &png
	}
	// if FullPage Options field is not passed set to false
	if r.Options.FullPage == nil {
		fullPage := false
		r.Options.FullPage = &fullPage
	}
	// if Width Options field is not passed set to DefaultWidth
	if r.Options.Width == nil {
		width := DefaultWidth
		r.Options.Width = &width
	}
	// if BlockAds Options field is not passed set to true
	if r.Options.BlockingOptions.BlockAds == nil {
		blockAds := true
		r.Options.BlockingOptions.BlockAds = &blockAds
	}
	// if HideCookieBanners Options field is not passed set to true
	if r.Options.BlockingOptions.HideCookieBanners == nil {
		cookie := true
		r.Options.BlockingOptions.HideCookieBanners = &cookie
	}
	// if ClickAccept Options field is not passed set to true
	if r.Options.BlockingOptions.ClickAccept == nil {
		accept := true
		r.Options.BlockingOptions.ClickAccept = &accept
	}
	// by default FailIfSelectorMissing should be false. Even if the selector is not found, it should not return any error
	if r.Options.SelectorOption.FailIfSelectorMissing == nil {
		failSelector := false
		r.Options.SelectorOption.FailIfSelectorMissing = &failSelector
	}
	// by default the Retina is set to false
	if r.Options.ImageOption.Retina == nil {
		retina := false
		r.Options.ImageOption.Retina = &retina
	}
	// by default the Quality is set to 80
	if r.Options.ImageOption.Quality == 0 {
		quality := 80
		r.Options.ImageOption.Quality = quality
	}
	// by default the Delay is set to 0
	if r.Options.WaitOption.Delay == nil {
		delay := 0
		r.Options.WaitOption.Delay = &delay
	}
	// by default the TimeOut is set to 30000 in milliseconds(3 seconds)
	if r.Options.WaitOption.TimeOut == nil {
		timeOut := 30000
		r.Options.WaitOption.TimeOut = &timeOut
	}

	return r
}

func (c *Client) SynchronousScreenshot(r Request) (*Response, error) {
	// the function shouldnt run if there was no url provided
	if r.Url == "" {
		return nil, ErrUrlRequired
	}
	// check if the Image quality is not above 100
	if r.Options.ImageOption.Quality > 100 {
		return nil, ErrImageQualityExceeded
	}

	r.parse()

	var downloadFileName string
	if r.Options.Download.DownloadFile {
		fileName := fmt.Sprintf("%v.%v", r.Options.Download.FileName, r.Format)
		downloadFileName = fileName
	}
	// setup the url
	url := fmt.Sprintf("%s/%v?url=%s&width=%v&full_page=%v&block_ads=%v&hide_cookie_banners=%v&click_accept=%v&selector=%v&retina=%v&quality=%v&download=%v&delay=%v&timeout=%v	",
		c.ApiKey, &r.Format, r.Url, &r.Options.Width, &r.Options.FullPage, &r.Options.BlockingOptions.BlockAds,
		&r.Options.BlockingOptions.HideCookieBanners, &r.Options.BlockingOptions.ClickAccept, &r.Options.SelectorOption.Selector,
		&r.Options.ImageOption.Retina, r.Options.ImageOption.Quality, downloadFileName, r.Options.WaitOption.Delay, r.Options.WaitOption.TimeOut,
	)

	var response Response

	if err := c.newRequest(http.MethodGet, url, nil, response); err != nil {
		return nil, err
	}

	return &response, nil
}
