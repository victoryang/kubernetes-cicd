package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	CDServer *CDServerCli
)

func init() {
	CDServer = NewCDServerClient("8080")
}

type CDServerCli struct {
	Addr 	string
}

func NewCDServerClient(addr string) *CDServerCli {
	return &CDServerCli{Addr: addr}
}

type CIBuildInfo struct {
	BuildCmd 	string 		`json:"buildcmd,omitempty"`
	Target		string 		`json:"from,omitempty"`
	Lang		string 		`json:"lang,omitempty"`
}

func (this *CDServerCli) GetBuildInfo(project string) *CIBuildInfo {
	respBody, err := this.Do("GET", this.Addr + "/projects/" + project + "/build_info", nil)
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

func (this *CDServerCli) CreateImage(project string, tag string) error {
	data := []byte(`{"Project":"` + project + `","Tag":"` + tag + `"}`)

	_, err := this.Do("POST", this.Addr + "/image/create_image", data)
	if err!=nil {
		fmt.Println("Create Image error:", err)
	}

	return err
}

func (this *CDServerCli) UpdateImage(project string, tag string, env string, deployStatus string, failLog string) error {
	data := []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `"}`)
	if len(failLog) > 0 {
		failLog = url.QueryEscape(failLog)
		data = []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `","MaintainPlan":"` + failLog + `"}`)
	}

	_,err := this.Do("POST", this.Addr + "/image/update_image", data)
	if err!=nil {
		fmt.Println("Update Image error:", err)
	}

	return err
}

func (this *CDServerCli) Do(method string, url string, data []byte) ([]byte,error) {
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