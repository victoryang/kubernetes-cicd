package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/victoryang/kubernetes-cicd/image"
	"github.com/victoryang/kubernetes-cicd/project"
)

func createProject(c *gin.Context) {
	// Get param
	param := &project.ConfigProjectParam{}
	err := c.BindJSON(param)
	if err != nil {
		err = fmt.Errorf("Wrong Create Project Param: %v", err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	c.Set("info", param.ProjectName)

	// Logic
	err = project.CreateProject(param)
	if err != nil {
		err = fmt.Errorf("Create Project %v Failed: %v", param.ProjectName, err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, nil)
}

func updateProject(c *gin.Context) {
	// Get param
	param := &project.ConfigProjectParam{}
	err := c.BindJSON(param)
	if err != nil {
		err = fmt.Errorf("Wrong Update Project Param: %v", err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	c.Set("info", param.ProjectName)

	// Logic
	err = project.UpdateProject(param)
	if err != nil {
		err = fmt.Errorf("Update Project %v Failed: %v", param.ProjectName, err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, nil)
}

func getProjectConfig(c *gin.Context) {
	// Get param
	projName := c.Query("proj")
	c.Set("info", projName)

	// Logic
	proj, err := project.GetProject(projName)
	if err != nil {
		err = fmt.Errorf("Get Project %v Failed: %v", projName, err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	old := project.ConfigProjectParam{}
	old.BuildCmd = proj.BuildCmd
	old.BuildDependency = proj.BuildDependency
	old.GitURL = proj.GitURL
	old.HTTPPort = strconv.Itoa(proj.HTTPPort)
	old.PreCmd = proj.PreCmd
	old.ProgramLanguage = proj.ProgramLanguage
	old.ProjectDesc = proj.ProjectDesc
	old.ProjectName = proj.Name
	old.StartCmd = proj.StartCmd
	old.StopCmd = proj.StopCmd
	old.TargetZip = proj.TargetZip
	old.UnzipDir = proj.UnzipDir
	HttpResponseWithSuccess(c, old)
}

func getProjectRuntime(c *gin.Context) {
	// Get param
	projName := c.Query("proj")
	c.Set("info", projName)

	// Logic
	// 从 project 库获取该项目运行时数据
	info, err := project.GetRuntimeInfo(projName)
	if err != nil {
		err = fmt.Errorf("Get Project %v runtime info Failed: %v", projName, err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, info)
}

func getEnvNodeNum(c *gin.Context) {
	// Get param
	env := project.Environment{}
	env.ProjName = c.Query("proj")
	c.Set("info", env.ProjName)
	env.Name = c.Query("env")
	env.Region = c.Query("region")

	num, err := env.GetNodeNum()
	if err != nil {
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, num)
}

func setEnvNodeNum(c *gin.Context) {
	// Get param
	type param struct {
		Proj    string `json:"proj"`
		Env     string `json:"env"`
		Region  string `json:"region"`
		NodeNum int32  `json:"node_num"`
	}
	p := param{}

	err := c.BindJSON(&p)
	if err != nil {
		HttpResponseWithBadRequest(c, err)
		return
	}
	env := project.Environment{}
	env.ProjName = p.Proj
	c.Set("info", env.ProjName)
	env.Name = p.Env
	env.Region = p.Region
	env.NodeNum = p.NodeNum

	err = env.SetNodeNum()
	if err != nil {
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, nil)
}

func createEnv(c *gin.Context) {
	// Get param
	env := &project.Environment{}
	err := c.BindJSON(env)
	if err != nil {
		err = fmt.Errorf("Wrong Create Environment Param: %v", err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	c.Set("info", env.ProjName)

	// Logic
	err = env.Create()
	if err != nil {
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, nil)
}

func updateEnvCodeVersion(c *gin.Context) {
	// Get param
	type param struct {
		Proj    string `json:"proj"`
		Env     string `json:"env"`
		Region  string `json:"region"`
		ImageID uint   `json:"image_id"`
	}
	p := param{}

	err := c.BindJSON(&p)
	if err != nil {
		HttpResponseWithBadRequest(c, err)
		return
	}
	env := project.Environment{}
	env.ProjName = p.Proj
	c.Set("info", env.ProjName)
	env.Name = p.Env
	env.Region = p.Region

	image := image.Image{}
	image.ID = p.ImageID
	_, err = image.Get()
	if err != nil {
		HttpResponseWithBadRequest(c, err)
		return
	}
	err = env.UpdateCodeVersion(image.FullName)
	if err != nil {
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, nil)
}

func receiveHook(c *gin.Context) {
	// Get param
	var projNames []string
	var gitAddr string
	var branch string
	var version string
	var timestamp string

	hookJSON := map[string]interface{}{}
	err := c.BindJSON(&hookJSON)
	if err != nil {
		err = fmt.Errorf("Wron Hook JSON: %v", err)
		HttpResponseWithBadRequest(c, err)
		return
	}

	projNamesStr := c.Query("projectname")
	branch = c.Query("branch")
	if len(hookJSON) > 0 {
		// 如果没有带projectname参数，默认使用git项目名
		if len(projNamesStr) < 1 {
			projNamesStr = hookJSON["repository"].(map[string]interface{})["name"].(string)
		}
		//get git address
		gitAddr = hookJSON["repository"].(map[string]interface{})["git_ssh_url"].(string)
		//get branch
		branch = strings.TrimPrefix(hookJSON["ref"].(string), "refs/heads/")
		//get version
		version = hookJSON["after"].(string)
		//get timestamp
		for _, commit := range hookJSON["commits"].([]interface{}) {
			if version == commit.(map[string]interface{})["id"].(string) {
				timestamp = commit.(map[string]interface{})["timestamp"].(string)
				t, timeErr := time.Parse("2006-01-02T15:04:05-07:00", timestamp)
				//如果序列化时间戳不合法，返回400
				if timeErr != nil {
					err = timeErr
					HttpResponseWithBadRequest(c, err)
					return
				}
				timestamp = strconv.FormatInt(t.Unix(), 10)
				break
			}
		}
		// 只取前8位
		version = version[:8]
	}
	projNames = strings.Split(projNamesStr, ",")

	// Logic
	for _, projName := range projNames {
		err = image.Create(projName, gitAddr, branch, version, timestamp)
		if err != nil {
			c.Error(fmt.Errorf("Create Image %v:%v:%v failed:%v", projName, branch, version, err))
		}
	}

	HttpResponseWithSuccess(c, nil)
}

func getImageList(c *gin.Context) {
	// Get param
	projName := c.Query("proj")
	c.Set("info", projName)
	env := c.Query("env")

	// Logic
	// 从 project 库获取该项目运行时数据
	images, err := image.GetImageList(projName, env)
	if err != nil {
		err = fmt.Errorf("Get %v %v  Image List Failed: %v", projName, env, err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, images)
}