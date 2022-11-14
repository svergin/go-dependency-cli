package goproxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Info struct {
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

type proxyparam struct {
	ProxyUrl string
}

type Client struct {
	http.Client
	params *proxyparam
}

const proxybase = "https://proxy.golang.org"

var ErrFehlendeModulversion = errors.New("Die Version des Moduls muss angegeben werden")

func (gpc *Client) WithParams(proxybase string) {
	gpc.params = &proxyparam{
		ProxyUrl: proxybase,
	}
}

func (gpc *Client) GetVersions(ctx context.Context, module string) ([]Info, error) {
	var result = make([]Info, 0)
	req, err := gpc.erstelleRequest(ctx, http.MethodGet, "%s/%s/@v/list", module, nil)
	if err != nil {
		return nil, err
	}

	res, err := gpc.Do(req)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occured, status was %q", res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	versions := strings.Split(string(body), "\n")
	for _, v := range versions {
		result = append(result, Info{Version: v})
	}

	return result, nil
}

func (gpc *Client) GetInfo(ctx context.Context, module string, version string) (*Info, error) {

	req, err := gpc.erstelleRequest(ctx, http.MethodGet, "%s/%s/@v/%s.info", module, &version)
	if err != nil {
		return nil, err
	}

	res, err := gpc.Do(req)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occured, status was %q", res.Status)
	}

	return getInfoFromBody(res)
}

func (gpc *Client) GetLatest(ctx context.Context, module string) (*Info, error) {

	req, err := gpc.erstelleRequest(ctx, http.MethodGet, "%s/%s/@latest", module, nil)
	if err != nil {
		return nil, err
	}

	res, err := gpc.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occured, status was %q", res.Status)
	}

	return getInfoFromBody(res)
}

func (gpc *Client) erstelleRequest(ctx context.Context, httpMethod string, resturl string, module string, version *string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, gpc.erzeugeUrl(resturl, module, version), nil)
}

func (gpc *Client) erzeugeUrl(resturl string, module string, version *string) string {
	if version != nil {
		return fmt.Sprintf(resturl, gpc.getProxyUrl(), module, *version)
	}
	return fmt.Sprintf(resturl, gpc.getProxyUrl(), module)

}

func (gpc *Client) getProxyUrl() string {
	if gpc.params == nil {
		return proxybase
	} else {
		return gpc.params.ProxyUrl
	}
}

func getInfoFromBody(res *http.Response) (*Info, error) {
	var result = Info{}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
