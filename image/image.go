package image

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/victoryang/kubernetes-cicd/orm"
	"github.com/victoryang/kubernetes-cicd/project"
)

func InitImageModule() {
	orm.MySQL.AutoMigrate(&Image{})
}

// Image store image info in db
type Image struct {
	gorm.Model
	FullName     string `gorm:"size:200;not null"`
	Project      string `gorm:"size:50;index;not null"`
	Env          string `gorm:"size:42;index;not null"`
	Tag          string `gorm:"size:100"`
	CodeVersion  string `gorm:"size:42"`
	CommitTime   time.Time
	DeployTime   time.Time
	DeployStatus string `gorm:"size:15"`
}

// Create insert the value into database
func (i *Image) Create() error {
	tagInfo := strings.Split(i.Tag, "_")
	// tag formate is 1515568519_65e49180_prod_base-v21
	if len(tagInfo) != 4 {
		return fmt.Errorf("wrong tag format %v", i.Tag)
	}
	timestamp, timeErr := strconv.ParseInt(tagInfo[0], 10, 64)
	if timeErr != nil {
		return fmt.Errorf("wrong tag timestamp %v", tagInfo[0])
	}
	i.CommitTime = time.Unix(timestamp, 0)
	i.DeployTime = i.CommitTime
	i.CodeVersion = tagInfo[1]
	if len(i.DeployStatus) == 0 {
		i.DeployStatus = "Building"
	}
	return orm.MySQL.Create(i).Error
}

// GetImageList return image list of current project and environment
func GetImageList(projName, env string) ([]Image, error) {
	images := []Image{}
	if err := orm.MySQL.Where(
		&Image{Project: projName, Env: env}).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

// Get a image info from db
func (i *Image) Get() (*Image, error) {
	err := orm.MySQL.Where(i).First(i).Error
	return i, err
}

// Create image by proj, git address, branch, commit id
// Write image info into DB
// Send task to etcd
func Create(projName, gitAddr, branch, commitID, timestamp string) error {
	proj, err := project.GetProject(projName)
	if err != nil {
		return fmt.Errorf("Get Porject %v failed:%v", projName, err)
	}

	img := &Image{}
	img.Project = projName
	img.CodeVersion = commitID
	img.Tag = timestamp + "_" + commitID + "_" + branch + "_" + proj.ProgramLanguage

	// create image in DB
	if err := img.Create(); err != nil {
		return err
	}

	// send task to etcd
	return nil
}