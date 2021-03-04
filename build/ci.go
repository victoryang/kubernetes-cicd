package build

import (
	"fmt"
	"bytes"
	"net/http"
	"net/url"
)

type RollingCli struct {
	Addr 	string
}

func NewRollingClient(addr string) *RollingCli {
	return &RollingCli{Addr: addr}
}

type BuildInfo struct {
	Commands 	string 		`json:"command"`
	Target		string 		`json:"from"`
	Lang		string 		`json:"lang"`
}

func (rc *RollingCli) GetBuildInfo(project string) *BuildInfo {
	fmt.Println("get build info from rolling")

	url := rc.Addr + "/projects/" + project + "/build_info"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	} else {
		resp.Body.Close()
	}
	return nil
}

func (rc *RollingCli) GetRuntimeInfo(project string, env string) error {
	fmt.Println("get runtime info from rolling")

	url := rc.Addr + "/projects/" + project + "/runtime_info"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		resp.Body.Close()
	}
	return nil
}

func (rc *RollingCli) CreateImage(project string, tag string) error {
	jsonStr := []byte(`{"Project":"` + project + `","Tag":"` + tag + `"}`)
	req, err := http.NewRequest("POST", rc.addr+"/image/create_image", bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		resp.Body.Close()
	}
	return nil
}

func (this *RollingCli) UpdateImage(project string, tag string, env string, deployStatus string, failLog string) error {
	jsonStr := []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `"}`)
	if len(failLog) > 0 {
		failLog = url.QueryEscape(failLog)
		jsonStr = []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `","MaintainPlan":"` + failLog + `"}`)
	}
	req, err := http.NewRequest("POST", "http://rolling.snowballfinance.com/image/update_image", bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		resp.Body.Close()
	}
	return nil
}