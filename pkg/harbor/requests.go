package harbor

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

func request(method, url, username, password string, rbody []byte) (int, *[]byte, error) {

	transport := http.DefaultTransport
	transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: !*hbTLS}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(*hbTimeout) * time.Second,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(rbody))
	if err != nil {
		return 0, nil, err
	}
	/*
		for _, cookie := range rcookies {
			req.AddCookie(cookie)
		}

		{'Content-type': 'application/json', 'Accept': 'application/json'}
	*/

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, &body, nil
}
