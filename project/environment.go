package project

import (
	"fmt"

	"github.com/victoryang/kubernetes-cicd/kubernetes"
)

// Environment is a project runtime info in a environment
type Environment struct {
	Name        string
	ProjName    string
	Region      string
	CodeBranch  string
	CodeVersion string
	NodeNum     int32
	CPU         int
	Memory      int
	Traffic     string
	UpTime      string
}

// Create a runtime environment of a project
func (e *Environment) Create() error {
	k8s := kubernetes.GetCli(e.Region)
	if k8s == nil {
		return fmt.Errorf("Get cluster failed: no region %s", e.Region)
	}
	return k8s.CreateDeployment(e.ProjName, e.Name, e.CodeBranch, e.NodeNum)
}

// UpdateCodeVersion by updating image
func (e *Environment) UpdateCodeVersion(image string) error {
	k8s := kubernetes.GetCli(e.Region)
	if k8s == nil {
		return fmt.Errorf("Get cluster failed: no region %s", e.Region)
	}
	return k8s.UpdateDeploymentImage(e.ProjName, e.Name, image)
}

// GetNodeNum get pod number of a deployment
func (e *Environment) GetNodeNum() (int32, error) {
	k8s := kubernetes.GetCli(e.Region)
	if k8s == nil {
		return 0, fmt.Errorf("Get cluster failed: no region %s", e.Region)
	}
	dm, err := k8s.GetDeployment(e.ProjName, e.Name)
	return *dm.Spec.Replicas, err
}

// SetNodeNum set pod number of a deployment
func (e *Environment) SetNodeNum() error {
	k8s := kubernetes.GetCli(e.Region)
	if k8s == nil {
		return fmt.Errorf("Get cluster failed: no region %s", e.Region)
	}
	return k8s.UpdateDeploymentPodNum(e.ProjName, e.Name, e.NodeNum)
}
