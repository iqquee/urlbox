package urlbox

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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
	DefaultWidth = 1280
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
func (c *Client) newRequest(method, reqURL string, reqBody interface{}) ([]byte, error) {
	newURL := c.BaseUrl + reqURL
	var body io.Reader

	if reqBody != nil {
		bb, err := json.Marshal(reqBody)
		if err != nil {
			return nil, errors.Wrap(err, "http client ::: unable to marshal request struct")
		}
		body = bytes.NewReader(bb)
	}

	// bearer := fmt.Sprintf("Bearer %v", c.SecretKey)
	req, err := http.NewRequest(method, newURL, body)
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", bearer)

	if err != nil {
		return nil, errors.Wrap(err, "http client ::: unable to create request body")
	}

	res, err := c.Http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http client ::: client failed to execute request")
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "http client ::: client failed to read file")
	}

	return b, nil
}
