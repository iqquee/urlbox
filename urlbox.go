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
	FileFormatPng = "png"
	//FileFormatJpeg for jpeg file format
	FileFormatJpeg = "jpeg"
	// FileFormatAvif for avif file format
	FileFormatAvif = "avif"
	//FileFormatWebp for webp file format
	FileFormatWebp = "webp"
	//FileFormatWebm for webm file format
	FileFormatWebm = "webm"
	// FileFormatPdf for pdf file format
	FileFormatPdf = "pdf"
	//FileFormatSvg for svg file format
	FileFormatSvg = "svg"
	// FileFormatHtml for html file format
	FileFormatHtml = "html"
	// FileFormatMd for md file format
	FileFormatMd = "md"
	// FileFormatMp4 for mp4 file format
	FileFormatMp4 = "mp4"
	// DefaultWidth default width of the screen shots to be taken
	DefaultWidth = "1280"
)

// client config
type Client struct {
	Http        http.Client
	BaseUrl     string
	ApiKey      string
	BearerToken string
}

// NewClient
func NewClient(h http.Client, apiKey, bearerToken string) *Client {
	return &Client{
		BaseUrl:     "https://api.urlbox.io/v1/",
		Http:        h,
		ApiKey:      apiKey,
		BearerToken: bearerToken,
	}
}

/*
newRequest makes a http request to the urlbox server and decodes the server response into the reqBody parameter passed into the newRequest method
*/
func (c *Client) newRequest(method, reqURL string, reqBody interface{}, resp interface{}) error {
	newURL := c.BaseUrl + reqURL
	var body io.Reader

	if reqBody != nil {
		bb, err := json.Marshal(reqBody)
		if err != nil {
			return errors.Wrap(err, "http client ::: unable to marshal request struct")
		}
		body = bytes.NewReader(bb)
	}

	bearer := fmt.Sprintf("Bearer %v", c.BearerToken)
	req, err := http.NewRequest(method, newURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)

	if err != nil {
		return errors.Wrap(err, "http client ::: unable to create request body")
	}

	res, err := c.Http.Do(req)
	if err != nil {
		return errors.Wrap(err, "http client ::: client failed to execute request")
	}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return errors.Wrap(err, "http client ::: unable to unmarshal response body")
	}

	return nil
}
