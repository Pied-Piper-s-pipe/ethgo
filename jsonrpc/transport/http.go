package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/umbracle/ethgo/jsonrpc/codec"
	// "github.com/valyala/fasthttp"
)

// HTTP is an http transport
type HTTP struct {
	seq  uint64
	addr string
	// client  *fasthttp.Client
	client  *http.Client
	headers map[string]string
}

func newHTTP(addr string, headers map[string]string) *HTTP {
	return &HTTP{
		addr: addr,
		client: &http.Client{
			Transport: &http.Transport{},
		},
		// client: &fasthttp.Client{
		// 	DialDualStack: true,
		// },
		headers: headers,
	}
}

func (h *HTTP) incSeq() uint64 {
	return atomic.AddUint64(&h.seq, 1)
}

// Close implements the transport interface
func (h *HTTP) Close() error {
	return nil
}

// Call implements the transport interface
func (h *HTTP) Call(method string, out interface{}, params ...interface{}) error {
	// Encode json-rpc request
	seq := h.incSeq()
	request := codec.Request{
		JsonRPC: "2.0",
		ID:      seq,
		Method:  method,
	}
	if len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return err
		}
		request.Params = data
	}
	raw, err := json.Marshal(request)
	// fmt.Printf("%s ----------> %s\n", h.addr, string(raw))
	if err != nil {
		return err
	}

	// req := fasthttp.AcquireRequest()
	// res := fasthttp.AcquireResponse()
	//
	// defer fasthttp.ReleaseRequest(req)
	// defer fasthttp.ReleaseResponse(res)
	//
	// req.SetRequestURI(h.addr)
	// req.Header.SetMethod("POST")
	// req.Header.SetContentType("application/json")
	// for k, v := range h.headers {
	// 	req.Header.Add(k, v)
	// }
	// req.SetBody(raw)
	//
	// if err := h.client.Do(req, res); err != nil {
	// 	return err
	// }
	// if sc := res.StatusCode(); sc != fasthttp.StatusOK {
	// 	return fmt.Errorf("status code is %d. response = %s", sc, string(res.Body()))
	// }

	req, err := http.NewRequest("POST", h.addr, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}

	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if sc := res.StatusCode; sc != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("status code is %d. response = %s", sc, string(body))
	}

	var response codec.Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		// if err := json.Unmarshal(res.Body(), &response); err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error
	}

	if err := json.Unmarshal(response.Result, out); err != nil {
		return err
	}
	return nil
}

// SetMaxConnsPerHost sets the maximum number of connections that can be established with a host
func (h *HTTP) SetMaxConnsPerHost(count int) {
	// h.client.MaxConnsPerHost = count
	h.client.Transport.(*http.Transport).MaxConnsPerHost = count
}
