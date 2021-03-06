package pipeline

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	yamlv2 "gopkg.in/yaml.v2"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-yaml/yaml"
)

const (
	PipelineKind = "pipeline"
	PipelineRunnerExec = "exec"
	PackagingWorkspace string = "/data/kubernetes-build/projects/"
)

type DroneBuildInfo struct {
	Project 	string
	Env 		string
	Timestamp 	string
	Version 	string
	Tag 		string
}

func ProcessRepoAndEventInfo(repoInfo *drone.Repo, buildInfo *drone.Build) *DroneBuildInfo {
	name := repoInfo.Name
	branch := strings.TrimPrefix(buildInfo.Ref, "refs/heads/")

	var env string
	switch branch {
	case "staging":
		env = "staging"
	case "release":
		env = "release"
	case "prod":
		env = "rc"
	case "sep":
		env = "sep"
	default:
		fmt.Println("env not supported: ", env)
		return nil
	}

	timestamp := strconv.FormatInt(buildInfo.Created, 10)
	version := buildInfo.After[:8]

	tag := timestamp + "_" + version + "_" + env + "_" + "base-go"

	return &DroneBuildInfo {
		Project: name,
		Env: env,
		Timestamp: timestamp,
		Version: version,
		Tag: tag,
	}
}

type BuildPipeline struct {
	DroneBuildInfo
	CIBuildInfo
	ImageName 	string
}

func NewBuildPipeline(repoInfo drone.Repo, buildInfo drone.Build) (*BuildPipeline,error) {

	droneInfo := ProcessRepoAndEventInfo(&repoInfo, &buildInfo)
	if droneInfo == nil {
		return nil, errors.New("Language not supported now")
	}

	ciBuildInfo := CDServer.GetBuildInfo(droneInfo.Project)
	if ciBuildInfo == nil {
		return nil, errors.New("CI build info not found")
	}

	switch ciBuildInfo.Lang {
	case "Java", "Go", "Node":
	default:
		return nil, errors.New("Language not supported now")
	}

	imageName := BuildImageName(droneInfo.Project, droneInfo.Tag)

	return &BuildPipeline {
		DroneBuildInfo: *droneInfo,
		CIBuildInfo: *ciBuildInfo,
		ImageName: imageName,
	}, nil
}

func (p *BuildPipeline) Compile() (string, error) {

	steps := p.CreateSteps()
	if steps==nil {
		return "", errors.New("create step fail")
	}

	pipeline := &yaml.Pipeline {
		Kind: PipelineKind,
		Type: PipelineRunnerExec,
		Name: p.Project,
		Steps: steps,
	}

	d, err := yamlv2.Marshal(pipeline)
	if err!=nil {
		fmt.Println("marshall error:", err)
		return "", err
	}

	return string(d), nil
}

func (p *BuildPipeline) CreateSteps() []*yaml.Container {
	steps := make([]*yaml.Container, 0)

	// Build step
	buildStep := p.CreateBuildStep()
	steps = append(steps, buildStep)

	// Packaging step
	packagingStep := p.CreatePackagingStep()
	steps = append(steps, packagingStep)

	// Publish step
	publishStep := p.CreatePublishStep()
	steps = append(steps, publishStep)

	// Clean up
	cleanUpStep := p.CreateCleanupStep()
	steps = append(steps, cleanUpStep)

	return steps
}

func (p *BuildPipeline) CreateEnvCommands() []string {

	return []string {
		"export",
	}
}

func (p *BuildPipeline) CreateBuildStep() *yaml.Container {

	buildCommands := []string {
		p.BuildCmd,
	}
	postBuildCommands := p.CreatePostBuildCommands()
	buildCommands = append(buildCommands, postBuildCommands...)

	return &yaml.Container {
		Name: "build",
		Commands: buildCommands,
	}
}

func (p *BuildPipeline) CreatePostBuildCommands() []string {
	var from string
	switch p.Lang {
	case "Java":
		from = path.Join(p.Project, p.Target)
	case "Go":
		from = p.Target
	}

	to := path.Join(PackagingWorkspace, p.Project, "release-"+p.Env)

	return []string {
		"cp -f " + from + " " + to,
	}
}

func (p *BuildPipeline) CreatePackagingStep() *yaml.Container {

	packagingCommand := []string {
		"cd " + path.Join(PackagingWorkspace, p.Project, "release-"+p.Env),
		"docker build -t " + p.ImageName + " .",
	}

	return &yaml.Container {
		Name: "packaging",
		Commands: packagingCommand,
	}
}

func (p *BuildPipeline) CreatePublishStep() *yaml.Container {

	publishCommand := []string {
		"echo $CI_JOB_TOKEN | docker login --username $CI_USER --password-stdin $CI_REGISTRY",
		"docker push " + p.ImageName,
	}

	return &yaml.Container {
		Name: "publish",
		Commands: publishCommand,
	}
}

func (p *BuildPipeline) CreateCleanupStep() *yaml.Container {

	cleanUpCommand := []string {
		"docker rmi " + p.ImageName,
	}

	return &yaml.Container {
		Name: "cleanup",
		Commands: cleanUpCommand,
	}
}