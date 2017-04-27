package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	ServiceURL *url.URL
	httpClient *http.Client
}

func LocalClient() *Client {
	c, _ := New("http://localhost:8082")
	return c
}

func New(u string) (*Client, error) {
	addr, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &Client{
		ServiceURL: addr,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}, nil
}

type Task struct {
	Content  string
	RunAt    time.Time
	RevertAt time.Time
	Region   string
}

type Form struct {
	Region, RunIn, RevertIn string
	Template                string
}

func (c *Client) List() ([]*Task, error) {
	var tasks []*Task

	addr, _ := url.Parse(c.ServiceURL.String())
	addr.Path = "tasks"

	resp, err := c.httpClient.Get(addr.String())
	if err != nil {
		return tasks, err
	}
	defer resp.Body.Close()

	if err := notOKStatus(addr.String(), resp); err != nil {
		return tasks, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (c *Client) Post(f Form) error {
	addr, _ := url.Parse(c.ServiceURL.String())
	addr.Path = "tasks"
	query := addr.Query()
	query.Add("region", f.Region)
	query.Add("run", f.RunIn)
	query.Add("revert", f.RevertIn)
	addr.RawQuery = query.Encode()

	resp, err := c.httpClient.Post(
		addr.String(),
		"application/text",
		strings.NewReader(f.Template),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := notOKStatus(addr.String(), resp); err != nil {
		return err
	}

	return nil
}

func notOKStatus(addr string, resp *http.Response) error {
	if code := resp.StatusCode; code != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Got %d status instead of 200 from '%s': %q", code, addr, body)
	}

	return nil
}
