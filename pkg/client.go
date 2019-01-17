package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/h2non/filetype"
	"io/ioutil"
	"net/http"
)

func upload(url string, path string, contentType string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(file))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", string(len(file)))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return errors.New("Bad response from S3")
	}
	return nil
}

// Client is a progimage client
type Client struct {
	BaseURL string
}

// NewClient creates a new progimage client
func NewClient() *Client {
	return &Client{
		BaseURL: "https://if5mras1pf.execute-api.us-east-1.amazonaws.com/dev",
	}
}

// UploadParams is the the parameters for uploading an image
type UploadParams struct {
	ContentType string `json:"contentType"`
}

// UploadResponse comes
type signResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// ErrUnknownContentType means the content type is unsupported or could not be detected
var ErrUnknownContentType = errors.New("Unknown content type")

// Upload will upload an image to progimage
func (c *Client) Upload(path string, params *UploadParams) (string, error) {
	if params.ContentType == "" {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		kind, unknown := filetype.Match(buf)
		if unknown != nil {
			return "", ErrUnknownContentType
		}
		params.ContentType = kind.MIME.Value
	}
	var signResp *signResponse
	err := c.postJSON("/", params, &signResp)
	if err != nil {
		return "", err
	}
	err = upload(signResp.URL, path, params.ContentType)
	return signResp.ID, err
}

func (c *Client) postJSON(path string, body interface{}, re interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.BaseURL, path), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("Server is unresponsive")
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&re); err != nil {
		return err
	}
	return nil
}
