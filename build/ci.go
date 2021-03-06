package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	ApiVersion1 = "api/v1"
)

var (
	CDServer *CDServerCli
)

type CDServerCli interface {
	GetBuildInfo(string) *CIBuildInfo
	CreateImage(string,string) error
	UpdateImage(string, string, string, string, string) error
}

// CI Build Info
type CIBuildInfo struct {
	BuildCmd 	string 		`json:"buildcmd,omitempty"`
	Target		string 		`json:"from,omitempty"`
	Lang		string 		`json:"lang,omitempty"`
}

// call CD server locally
func InitCDServerClientWithLocaldMode() {
	CDServer = nil
}


// Call CD server remotely
func InitCDServerClientWithRemoteMode() {
	CDServer = NewRemoteCDServerClient("http://127.0.0.1:8080")
}

type RemoteCDServerCli struct {
	BaseUrl 	string
}

func NewRemoteCDServerClient(addr string) *RemoteCDServerCli {
	return &RemoteCDServerCli{BaseUrl: addr}
}

func (this *RemoteCDServerCli) GetBuildInfo(project string) *CIBuildInfo {
	url := path.Join(this.BaseUrl, ApiVersion1, "projects", project, "/build/config")
	respBody, err := this.Do("GET", url, nil)
	if err!=nil {
		fmt.Println("Get Build Info error:", err)
		return nil
	}

	var info CIBuildInfo
	err = json.Unmarshal(respBody, &info)
	if err!=nil {
		fmt.Println("UnMarshal Build Info error:", err)
		return nil
	}

	return &info
}

func (this *RemoteCDServerCli) CreateImage(project string, tag string) error {
	data := []byte(`{"Project":"` + project + `","Tag":"` + tag + `"}`)

	url := path.Join(this.BaseUrl, ApiVersion1, "/images/")
	_, err := this.Do("POST", url, data)
	if err!=nil {
		fmt.Println("Create Image error:", err)
	}

	return err
}

func (this *RemoteCDServerCli) UpdateImage(project string, tag string, env string, deployStatus string, failLog string) error {
	data := []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `"}`)
	if len(failLog) > 0 {
		failLog = url.QueryEscape(failLog)
		data = []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `","MaintainPlan":"` + failLog + `"}`)
	}

	url := path.Join(this.BaseUrl, ApiVersion1, "/images/")
	_,err := this.Do("POST", url, data)
	if err!=nil {
		fmt.Println("Update Image error:", err)
	}

	return err
}

func (this *RemoteCDServerCli) Do(method string, url string, data []byte) ([]byte,error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err!=nil {
		return nil, errors.New(fmt.Sprintf("Create http request error:", err))
	}

	req.Header.Set("User-Agent", "Kubernetes Build")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err!=nil {
		return nil, errors.New(fmt.Sprintf("Exec http request error:", err))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Response code not 200")
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}