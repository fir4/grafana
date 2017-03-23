package imguploader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/grafana/grafana/pkg/util"
)

type WebdavUploader struct {
	url        string
	username   string
	password   string
	public_url string
}

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 60 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

var netClient = &http.Client{
	Timeout:   time.Second * 60,
	Transport: netTransport,
}

func (u *WebdavUploader) Upload(pa string) (string, error) {
	url, _ := url.Parse(u.url)
	filename := util.GetRandomString(20) + ".png"
	url.Path = path.Join(url.Path, filename)

	imgData, err := ioutil.ReadFile(pa)
	req, err := http.NewRequest("PUT", url.String(), bytes.NewReader(imgData))

	if u.username != "" {
		req.SetBasicAuth(u.username, u.password)
	}

	res, err := netClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("Failed to upload image. Returned statuscode %v body %s", res.StatusCode, body)
	}

	if u.public_url != "" {
		return (u.public_url + filename), nil
	} else {
		return url.String(), nil
	}
}

func NewWebdavImageUploader(url, username, password, public_url string) (*WebdavUploader, error) {
	return &WebdavUploader{
		url:        url,
		username:   username,
		password:   password,
		public_url: public_url,
	}, nil
}
