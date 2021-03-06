package build

import (
	"fmt"
	"bytes"
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

func NewCDServerClient(addr string) *CDServer {
	return &CDServer{Addr: addr}
}

type CIBuildInfo struct {
	BuildCmd 	string 		`json:"buildcmd"`
	Target		string 		`json:"from"`
	Lang		string 		`json:"lang"`
}

func (this *CDServerCli) GetBuildInfo(project string) *CIBuildInfo {
	fmt.Println("get build info from cd server")

	url := this.Addr + "/projects/" + project + "/build_info"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Kubernetes Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	} else {
		resp.Body.Close()
	}
	return nil
}

func (this *CDServerCli) CreateImage(project string, tag string) error {
	jsonStr := []byte(`{"Project":"` + project + `","Tag":"` + tag + `"}`)
	req, err := http.NewRequest("POST", this.Addr+"/image/create_image", bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Kubernetes Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		resp.Body.Close()
	}
	return nil
}

func (this *CDServerCli) UpdateImage(project string, tag string, env string, deployStatus string, failLog string) error {
	jsonStr := []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `"}`)
	if len(failLog) > 0 {
		failLog = url.QueryEscape(failLog)
		jsonStr = []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `","MaintainPlan":"` + failLog + `"}`)
	}
	req, err := http.NewRequest("POST", this.Addr + "/image/update_image", bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Kubernetes Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		resp.Body.Close()
	}
	return nil
}