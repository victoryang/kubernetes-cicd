// Package kubernetes is a wapper of kubernetes client
package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/victoryang/kubernetes-cicd/config"
)

var clients map[string]*KubeCli

func init() {
	clients = map[string]*KubeCli{}
	for _, region := range config.Regions {
		c, _ := newCli(region.Addr)
		clients[region.Name] = c
	}
}

// GetCli return a client by cluste name
func GetCli(clusterName string) *KubeCli {
	if _, ok := clients[clusterName]; !ok {
		return nil
	}
	return clients[clusterName]
}

// GetAll return all clients
func GetAll() map[string]*KubeCli {
	return clients
}

// KubeCli is client of a k8s cluster
type KubeCli struct {
	cliSet *kubernetes.Clientset
}

// newCli return a k8s cli by master url
func newCli(masterURL string) (*KubeCli, error) {
	config, err := clientcmd.BuildConfigFromFlags(masterURL, "")
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return &KubeCli{cliSet: clientset}, nil
}

// GetDeployments show deployments info in k8s cluster of a project
func (c *KubeCli) GetDeployments(projName string) (*appsv1.DeploymentList, error) {
	return c.cliSet.AppsV1().Deployments(projName).List(context.TODO(), metav1.ListOptions{})
}

// CreateDeployment create a deployment of a project
func (c *KubeCli) CreateDeployment(projName, envName, codeBranch string, podNum int32) error {
	//envName = projName + "-" + envName
	env := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: envName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(podNum),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":    envName,
						"branch": codeBranch,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  envName,
							Image: "nginx", // for default image
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{},
							},
						},
					},
				},
			},
		},
	}
	_, err := c.cliSet.AppsV1().Deployments(projName).Create(context.TODO(), env, metav1.CreateOptions{})
	return err
}

// GetDeployment by namespace and deployment name
func (c *KubeCli) GetDeployment(projName, envName string) (*appsv1.Deployment, error) {
	return c.cliSet.AppsV1().Deployments(projName).Get(context.TODO(), envName, metav1.GetOptions{})
}

// UpdateDeploymentImage update image of a deployment
func (c *KubeCli) UpdateDeploymentImage(projName, envName, image string) error {
	dm, err := c.cliSet.AppsV1().Deployments(projName).Get(context.TODO(), envName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for i := range dm.Spec.Template.Spec.Containers {
		dm.Spec.Template.Spec.Containers[i].Image = image
	}
	_, err = c.cliSet.AppsV1().Deployments(projName).Update(context.TODO(), dm, metav1.UpdateOptions{})
	return err
}

// UpdateDeploymentPodNum pod number of a deployment
func (c *KubeCli) UpdateDeploymentPodNum(projName, envName string, podNum int32) error {
	dm, err := c.cliSet.AppsV1().Deployments(projName).Get(context.TODO(), envName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	dm.Spec.Replicas = int32Ptr(podNum)
	_, err = c.cliSet.AppsV1().Deployments(projName).Update(context.TODO(), dm, metav1.UpdateOptions{})
	return err
}

// CreateNamespace create a project in k8s cluster
// A project is a namespace.
func (c *KubeCli) CreateNamespace(projName string) error {
	proj := v1.Namespace{}
	proj.Name = projName
	_, err := c.cliSet.CoreV1().Namespaces().Create(context.TODO(), &proj, metav1.CreateOptions{})
	return err
}

func int32Ptr(i int32) *int32 { return &i }