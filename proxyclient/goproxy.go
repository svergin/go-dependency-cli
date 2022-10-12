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

	"github.com/rogpeppe/go-internal/modfile"
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

type GoProxyClient struct {
	http.Client
	params *proxyparam
}

type DependencyReport struct {
}

const proxybase = "https://proxy.golang.org"
const timeout = 200

var ErrFehlendeModulversion = errors.New("Die Version des Moduls muss angegeben werden")

//TODO: https://dev.azure.com/nortal-de-dev-training/_git/learning-go?path=/doc/1.3.md&version=GBmain&_a=preview

func (gpc *GoProxyClient) ErstelleReport() (DependencyReport, error) {
	dr := DependencyReport{}

	data, err := gpc.GetGoMod(context.Background())
	if err != nil {
		return dr, err
	}
	f, err := modfile.Parse("myfilename", data, nil)
	if err != nil {
		return dr, err
	}
	fmt.Println(f.Require)

	return dr, nil
}

func (gpc *GoProxyClient) WithParams(modulename string, version *string) {
	gpc.params = &proxyparam{
		Module:      modulename,
		TimeoutInMs: timeout,
		ProxyUrl:    proxybase,
		Version:     version,
	}
}

func (gpc *GoProxyClient) GetVersions(ctx context.Context) ([]Info, error) {
	var result = make([]Info, 0)
	myContext, cancel := context.WithTimeout(ctx, time.Duration(gpc.params.TimeoutInMs)*time.Millisecond)
	req, err := erstelleRequest(myContext, http.MethodGet, "%s/%s/@v/list", gpc.params)
	defer cancel()
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

func (gpc *GoProxyClient) GetInfo(ctx context.Context) (*Info, error) {
	if gpc.params.Version == nil {
		return nil, ErrFehlendeModulversion
	}
	myContext, cancel := context.WithTimeout(ctx, time.Duration(gpc.params.TimeoutInMs)*time.Millisecond)
	req, err := erstelleRequest(myContext, http.MethodGet, "%s/%s/@v/%s.info", gpc.params)
	defer cancel()
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

func (gpc *GoProxyClient) GetGoMod(ctx context.Context) ([]byte, error) {
	if gpc.params.Version == nil {
		return nil, ErrFehlendeModulversion
	}
	myContext, cancel := context.WithTimeout(ctx, time.Duration(gpc.params.TimeoutInMs)*time.Millisecond)
	req, err := erstelleRequest(myContext, http.MethodGet, "%s/%s/@v/%s.mod", gpc.params)
	defer cancel()
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

	return io.ReadAll(res.Body)
}

func (gpc *GoProxyClient) GetLatest(ctx context.Context) (*Info, error) {

	myContext, cancel := context.WithTimeout(ctx, time.Duration(gpc.params.TimeoutInMs)*time.Millisecond)
	req, err := erstelleRequest(myContext, http.MethodGet, "%s/%s/@latest", gpc.params)
	defer cancel()
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

func erstelleRequest(ctx context.Context, httpMethod string, resturl string, param *proxyparam) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, erzeugeUrl(resturl, param), nil)
}

func erzeugeUrl(resturl string, param *proxyparam) string {
	if param.Version != nil {
		return fmt.Sprintf(resturl, param.ProxyUrl, param.Module, *param.Version)
	}
	return fmt.Sprintf(resturl, param.ProxyUrl, param.Module)

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

func (dr *DependencyReport) String() string {
	return ""
}
