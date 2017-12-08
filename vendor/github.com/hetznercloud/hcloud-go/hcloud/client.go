package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// Endpoint is the base URL of the API.
const Endpoint = "https://api.hetzner.cloud/v1"

// UserAgent is the value for the User-Agent header sent with each request.
const UserAgent = "hcloud-go/" + Version

// A BackoffFunc returns the duration to wait before performing the
// next retry. The retries argument specifies how many retries have
// already been performed. When called for the first time, retries is 0.
type BackoffFunc func(retries int) time.Duration

// ConstantBackoff returns a BackoffFunc which backs off for
// constant duration d.
func ConstantBackoff(d time.Duration) BackoffFunc {
	return func(_ int) time.Duration {
		return d
	}
}

// ExponentialBackoff returns a BackoffFunc which implements an exponential
// backoff using the formula: b^retries * d
func ExponentialBackoff(b float64, d time.Duration) BackoffFunc {
	return func(retries int) time.Duration {
		return time.Duration(math.Pow(b, float64(retries))) * d
	}
}

// Client is a client for the Hetzner Cloud API.
type Client struct {
	endpoint    string
	token       string
	backoffFunc BackoffFunc
	httpClient  *http.Client

	Action     ActionClient
	FloatingIP FloatingIPClient
	Server     ServerClient
	SSHKey     SSHKeyClient
}

// A ClientOption is used to configure a Client.
type ClientOption func(*Client)

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) {
		client.endpoint = strings.TrimRight(endpoint, "/")
	}
}

// WithToken configures a Client to use the specified token for authentication.
func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.token = token
	}
}

// WithBackoffFunc configures a Client to use the specified backoff function.
func WithBackoffFunc(f BackoffFunc) ClientOption {
	return func(client *Client) {
		client.backoffFunc = f
	}
}

// NewClient creates a new client.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		endpoint:    Endpoint,
		httpClient:  &http.Client{},
		backoffFunc: ExponentialBackoff(2, 500*time.Millisecond),
	}

	for _, option := range options {
		option(client)
	}

	client.Action = ActionClient{client: client}
	client.FloatingIP = FloatingIPClient{client: client}
	client.Server = ServerClient{client: client}
	client.SSHKey = SSHKeyClient{client: client}

	return client
}

// NewRequest creates an HTTP request against the API. The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.endpoint + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req = req.WithContext(ctx)
	return req, nil
}

// Do performs an HTTP request against the API.
func (c *Client) Do(r *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	response := &Response{Response: resp}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	if err := response.readMeta(body); err != nil {
		return response, fmt.Errorf("hcloud: error reading response meta data: %s", err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		err := errorFromResponse(resp, body)
		if err == nil {
			err = fmt.Errorf("hcloud: server responded with status code %d", resp.StatusCode)
		}
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bytes.NewReader(body))
		} else {
			err = json.Unmarshal(body, v)
		}
	}

	return response, err
}

func (c *Client) backoff(retries int) {
	time.Sleep(c.backoffFunc(retries))
}

func (c *Client) all(f func(int) (*Response, error)) (*Response, error) {
	var (
		retries = 0
		page    = 1
	)
	for {
		resp, err := f(page)
		if err != nil {
			if err, ok := err.(Error); ok {
				if err.Code == ErrorCodeLimitReached {
					c.backoff(retries)
					retries++
					continue
				}
			}
			return nil, err
		}
		retries = 0
		if resp.Meta.Pagination == nil || resp.Meta.Pagination.NextPage == 0 {
			return resp, nil
		}
		page = resp.Meta.Pagination.NextPage
	}
}

func errorFromResponse(resp *http.Response, body []byte) error {
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil
	}

	var apiError struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &apiError); err != nil {
		return nil
	}
	if apiError.Error.Code == "" && apiError.Error.Message == "" {
		return nil
	}
	return Error{
		Code:    ErrorCode(apiError.Error.Code),
		Message: apiError.Error.Message,
	}
}

// Response represents a response from the API. It embeds http.Response.
type Response struct {
	*http.Response
	Meta Meta
}

// ReadBody reads and returns the response's body. After reading the response's body
// is recovered so it can be read again.
//
// TODO(thcyron): Does this method really need to be exported?
func (r *Response) ReadBody() ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	return body, err
}

func (r *Response) readMeta(body []byte) error {
	if h := r.Header.Get("RateLimit-Limit"); h != "" {
		r.Meta.Ratelimit.Limit, _ = strconv.Atoi(h)
	}
	if h := r.Header.Get("RateLimit-Remaining"); h != "" {
		r.Meta.Ratelimit.Remaining, _ = strconv.Atoi(h)
	}
	if h := r.Header.Get("RateLimit-Reset"); h != "" {
		if ts, err := strconv.ParseInt(h, 10, 64); err == nil {
			r.Meta.Ratelimit.Reset = time.Unix(ts, 0)
		}
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var s schema.MetaResponse
		if err := json.Unmarshal(body, &s); err != nil {
			return err
		}
		if s.Meta.Pagination != nil {
			p := PaginationFromSchema(*s.Meta.Pagination)
			r.Meta.Pagination = &p
		}
	}

	return nil
}

// Meta represents meta information included in an API response.
type Meta struct {
	Pagination *Pagination
	Ratelimit  Ratelimit
}

// Pagination represents pagination meta information.
type Pagination struct {
	Page         int
	PerPage      int
	PreviousPage int
	NextPage     int
	LastPage     int
	TotalEntries int
}

// Ratelimit represents ratelimit information.
type Ratelimit struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// ListOpts specifies options for listing resources.
type ListOpts struct {
	Page    int // Page (starting at 1)
	PerPage int // Items per page (0 means default)
}

func valuesForListOpts(opts ListOpts) url.Values {
	vals := url.Values{}
	if opts.Page > 0 {
		vals.Add("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		vals.Add("per_page", strconv.Itoa(opts.PerPage))
	}
	return vals
}
