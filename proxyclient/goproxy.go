package proxyclient

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
	Module      string
	TimeoutInMs uint32
	ProxyUrl    string
	Version     *string
}

const proxybase = "https://proxy.golang.org"
const timeout = 200

var ErrFehlendeModulversion = errors.New("Die Version des Moduls muss angegeben werden")

//TODO: https://dev.azure.com/nortal-de-dev-training/_git/learning-go?path=/doc/1.3.md&version=GBmain&_a=preview

func DefaultParam(modulename string, version *string) *proxyparam {
	return &proxyparam{
		Module:      modulename,
		TimeoutInMs: timeout,
		ProxyUrl:    proxybase,
		Version:     version,
	}
}

func GetVersions(ctx context.Context, param *proxyparam) ([]Info, error) {
	var result = make([]Info, 0)
	client := http.Client{}
	url := fmt.Sprintf("%s/%s/@v/list", param.ProxyUrl, param.Module)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

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

func GetInfo(ctx context.Context, param *proxyparam) (*Info, error) {
	if param.Version == nil {
		return nil, ErrFehlendeModulversion
	}
	url := fmt.Sprintf("%s/%s/@v/%s.info", param.ProxyUrl, param.Module, *param.Version)
	client := http.Client{}
	myContext, cancel := context.WithTimeout(ctx, time.Duration(param.TimeoutInMs)*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(myContext, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occured, status was %q", res.Status)
	}

	return getInfoFromBody(res)
}

func GetLatest(ctx context.Context, param *proxyparam) (*Info, error) {

	url := fmt.Sprintf("%s/%s/@latest", param.ProxyUrl, param.Module)
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occured, status was %q", res.Status)
	}

	return getInfoFromBody(res)
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
