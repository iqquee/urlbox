package urlbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	//FileFormatPng for png file format
	FileFormatPng string = "png"
	//FileFormatJpeg for jpeg file format
	FileFormatJpeg string = "jpeg"
	// FileFormatAvif for avif file format
	FileFormatAvif string = "avif"
	//FileFormatWebp for webp file format
	FileFormatWebp string = "webp"
	//FileFormatWebm for webm file format
	FileFormatWebm string = "webm"
	// FileFormatPdf for pdf file format
	FileFormatPdf string = "pdf"
	//FileFormatSvg for svg file format
	FileFormatSvg string = "svg"
	// FileFormatHtml for html file format
	FileFormatHtml string = "html"
	// FileFormatMd for md file format
	FileFormatMd string = "md"
	// FileFormatMp4 for mp4 file format
	FileFormatMp4 string = "mp4"
	// DefaultWidth default width of the screen shot to be taken
	DefaultWidth int = 1280
	// DefaultWidth default height of the screen shot to be taken
	DefaultHeight   int = 1024
	DefaultSelector     = "random_default_selector_will_not_be_found" // DefaultSelector is set so that you can have the flexibility of adding a selector params if there is a need for it. This selector does not affect the response from urlbox as it is not valid. if the selector is not found, Urlbox will take a normal viewport screenshot as in this case.
)

// client config
type Client struct {
	Http      http.Client
	BaseUrl   string
	ApiKey    string
	SecretKey string
}

// New is the urlbox config initializer
func New(h http.Client, apiKey, secretKey string) *Client {
	return &Client{
		BaseUrl:   "https://api.urlbox.io/v1/",
		Http:      h,
		ApiKey:    apiKey,
		SecretKey: secretKey,
	}
}

/*
newRequest makes a http request to the urlbox server and decodes the server response into the reqBody parameter passed into the newRequest method
*/
func (c *Client) newRequest(method, reqURL string, reqBody interface{}) ([]byte, int, error) {
	newURL := c.BaseUrl + reqURL
	var body io.Reader

	if reqBody != nil {
		bb, err := json.Marshal(reqBody)
		if err != nil {
			return nil, 0, errors.Wrap(err, "http client ::: unable to marshal request struct")
		}
		body = bytes.NewReader(bb)
	}

	req, err := http.NewRequest(method, newURL, body)
	if reqBody != nil {
		bearer := fmt.Sprintf("Bearer %v", c.SecretKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", bearer)
	}

	if err != nil {
		return nil, 0, errors.Wrap(err, "http client ::: unable to create request body")
	}

	res, err := c.Http.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "http client ::: client failed to execute request")
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, errors.Wrap(err, "http client ::: client failed to read file")
	}

	return b, res.StatusCode, nil
}
