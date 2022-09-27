package proxyclient

import (
	"encoding/json"
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

const proxybase = "https://proxy.golang.org"

func GetVersions(module string) ([]Info, error) {
	var result = make([]Info, 0)

	url := fmt.Sprintf("%s/%s/@v/list", proxybase, module)

	res, err := http.Get(url)

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

func GetInfo(module string, version string) (*Info, error) {
	url := fmt.Sprintf("%s/%s/@v/%s.info", proxybase, module, version)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error occured, status was %q", res.Status)
	}

	return getInfoFromBody(res)
}

func GetLatest(module string) (*Info, error) {

	url := fmt.Sprintf("%s/%s/@latest", proxybase, module)

	res, err := http.Get(url)

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
