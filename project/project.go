package project

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/victoryang/kubernetes-cicd/kubernetes"
	"github.com/victoryang/kubernetes-cicd/models"
	"github.com/victoryang/kubernetes-cicd/orm"
)

func InitProjectModule() {
	orm.MySQL.AutoMigrate(&Project{})
}

// ConfigProjectParam is used for create project
type ConfigProjectParam struct {
	ProjectName     string `json:"projectName" binding:"required"`
	GitURL          string `json:"gitURL" binding:"required"`
	HTTPPort        string `json:"httpPort"`
	ProjectDesc     string `json:"projectDesc"`
	BuildCmd        string `json:"buildCmd" binding:"required"`
	TargetZip       string `json:"targetZip" binding:"required"`
	UnzipDir        string `json:"unzipDir" binding:"required"`
	ProgramLanguage string `json:"programLanguage" binding:"required"`
	BuildDependency string `json:"buildDependency"`
	StartCmd        string `json:"startCmd" binding:"required"`
	StopCmd         string `json:"stopCmd" binding:"required"`
	PreCmd          string `json:"preCmd"`
}

// Project is a rolling project info
type Project struct {
	gorm.Model      `json:"-"`
	Name            string `json:"projectName" binding:"required" gorm:"size:64;index;not null;unique"`
	GitURL          string `json:"gitURL" binding:"required" gorm:"size:128;not null"`
	HTTPPort        int    `json:"httpPort"`
	ProjectDesc     string `json:"projectDesc" gorm:"type:longtext"`
	BuildCmd        string `json:"buildCmd" binding:"required" gorm:"type:longtext;not null"`
	TargetZip       string `json:"targetZip" binding:"required" gorm:"size:128;not null"`
	UnzipDir        string `json:"unzipDir" binding:"required" gorm:"size:128;not null"`
	ProgramLanguage string `json:"programLanguage" binding:"required" gorm:"size:128;not null"`
	BuildDependency string `json:"buildDependency" gorm:"type:longtext"`
	StartCmd        string `json:"startCmd" binding:"required" gorm:"size:128;not null"`
	StopCmd         string `json:"stopCmd" binding:"required" gorm:"size:128;not null"`
	PreCmd          string `json:"preCmd" gorm:"type:longtext"`
}

// RuntimeInfo is info of a project runtime
type RuntimeInfo struct {
	Regions []region
}

// region is a project runtime info in a region
type region struct {
	Name string
	Envs []Environment
}

// CreateProject create a project
func CreateProject(param *ConfigProjectParam) error {
	proj := &Project{}
	// 1. 判断输入参数格式是否正确
	if !strings.HasPrefix(param.GitURL, "git@git.snowballfinance.com:") &&
		!strings.HasSuffix(param.GitURL, ".git") {
		return errors.New("Git 地址应该以 git@git.snowballfinance.com 开始，以 .git 结尾")
	}
	if len(param.HTTPPort) > 0 {
		port, err := strconv.Atoi(strings.TrimSpace(param.HTTPPort))
		if err != nil {
			return err
		}
		proj.HTTPPort = port
	}
	// 2. 将信息存入数据库，存入之前 TrimSpace，如果项目已经存在，则报错
	proj.transParamToProj(param)
	if err := orm.MySQL.Create(proj).Error; err != nil {
		return err
	}

	models.Logger.Info("insert mysql succed")
	// 3. 在 kubernetes 中创建 namespace #应该在项目创建后生成实例之前做
	for _, k8s := range kubernetes.GetAll() {
		err := k8s.CreateNamespace(proj.Name)
		if err != nil {
			models.Logger.Error("k8s create namespace failed", err)
			return err
		}
	}

	models.Logger.Info("kbs create namespace succed")

	// TODO 4. 在 gitlab 中自动添加 hook 及触发镜像构建
	return nil
}

// UpdateProject update a project
func UpdateProject(param *ConfigProjectParam) error {
	proj := &Project{}
	// 1. 判断输入参数格式是否正确
	if !strings.HasPrefix(param.GitURL, "git@git.snowballfinance.com:") ||
		!strings.HasSuffix(param.GitURL, ".git") {
		return errors.New("Git 地址应该以 git@git.snowballfinance.com 开始，以 .git 结尾")
	}
	if len(param.HTTPPort) > 0 {
		port, err := strconv.Atoi(strings.TrimSpace(param.HTTPPort))
		if err != nil {
			return err
		}
		proj.HTTPPort = port
	}
	// 2. 将信息存入数据库，存入之前 TrimSpace
	proj.transParamToProj(param)
	proj.updateConfig()
	// 在 kubernetes 无需变更 namespace 因为项目名没变
	// TODO 3. 在 gitlab 中自动添加 hook 及触发镜像构建
	return nil
}

func (proj *Project) transParamToProj(param *ConfigProjectParam) {
	proj.Name = strings.TrimSpace(param.ProjectName)
	proj.GitURL = strings.TrimSpace(param.GitURL)
	proj.ProjectDesc = strings.TrimSpace(param.ProjectDesc)
	proj.BuildCmd = strings.TrimSpace(param.BuildCmd)
	proj.TargetZip = strings.TrimSpace(param.TargetZip)
	proj.UnzipDir = strings.TrimSpace(param.UnzipDir)
	proj.ProgramLanguage = strings.TrimSpace(param.ProgramLanguage)
	proj.BuildDependency = strings.TrimSpace(param.BuildDependency)
	proj.StartCmd = strings.TrimSpace(param.StartCmd)
	proj.StopCmd = strings.TrimSpace(param.StopCmd)
	proj.PreCmd = strings.TrimSpace(param.PreCmd)

	if len(param.HTTPPort) > 0 {
		port, err := strconv.Atoi(strings.TrimSpace(param.HTTPPort))
		if err == nil {
			proj.HTTPPort = port
		}
	}
}

func (proj *Project) updateConfig() error {
	oldProj, err := GetProject(proj.Name)
	if err != nil {
		return err
	}
	proj.Model = oldProj.Model
	return orm.MySQL.Save(proj).Error
}

// GetAllProjects return all projects
func GetAllProjects() ([]Project, error) {
	projects := []Project{}
	err := orm.MySQL.Find(&projects).Error

	return projects, err
}

// GetProject by project name
func GetProject(projName string) (*Project, error) {
	proj := &Project{}
	err := orm.MySQL.First(proj, "name = ?", projName).Error
	return proj, err
}

// GetRuntimeInfo by project name
func GetRuntimeInfo(projName string) (*RuntimeInfo, error) {
	runtimeInfo := &RuntimeInfo{}
	for regionName, k8s := range kubernetes.GetAll() {
		// 1. 获取环境信息：namespace 下 deployment 信息
		envs, err := k8s.GetDeployments(projName)
		if err != nil {
			return runtimeInfo, err
		}

		// 2. 获取每个环境下的运行信息：deployment 下的 pod 信息、资源使用
		rgn := region{}
		rgn.Name = regionName
		for _, e := range envs.Items {
			env := Environment{}
			env.Name = e.Name
			env.ProjName = projName
			branch, ok := e.Spec.Template.ObjectMeta.Labels["branch"]
			env.CodeBranch = branch
			if !ok {
				env.CodeBranch = "未知"
			}
			env.CodeVersion = "无"
			env.Region = rgn.Name
			env.CPU = -1
			env.Memory = -1
			env.NodeNum = *e.Spec.Replicas
			env.Traffic = "待开发"
			env.UpTime = "待开发"
			rgn.Envs = append(rgn.Envs, env)
		}
		runtimeInfo.Regions = append(runtimeInfo.Regions, rgn)
		sort.Slice(runtimeInfo.Regions, func(i, j int) bool {
			return len(runtimeInfo.Regions[i].Envs) > len(runtimeInfo.Regions[j].Envs)
		})
		// TODO 3. 获取 nginx 中流量调度信息
	}
	return runtimeInfo, nil
}