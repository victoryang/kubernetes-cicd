package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"

	"gopkg.in/ldap.v2"

	"github.com/victoryang/kubernetes-cicd/gitlab"
	"github.com/victoryang/kubernetes-cicd/orm"
	"github.com/victoryang/kubernetes-cicd/project"
)

var (
	LDAP_ADDR string
	LDAP_PWD string
)

func InitAuthModule(address, password string) {
	LDAP_ADDR = address
	LDAP_PWD = password

	orm.MySQL.AutoMigrate(&User{})
}

// User to store user's infomation
type User struct {
	gorm.Model `json:"-"`
	Name       string   `json:"name" gorm:"size:64;index;not null;unique"`
	Token      string   `json:"token" gorm:"size:64;index;not null;unique"`
	Groups     []string `json:"groups" gorm:"-"`
	Projects   []string `json:"projects" gorm:"-"`
}

// AuthByLdap check if user has authority in ldap
func (u *User) AuthByLdap(username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	//l, err := ldap.Dial("tcp", "10.10.130.3:389")
	l, err := ldap.Dial("tcp", LDAP_ADDR)
	if err != nil {
		return err
	}
	defer l.Close()

	//bindusername := "cn=readonly,dc=snowballfinance,dc=com"
	//bindpassword := "a76JzLCUKLYGcDlU"
	bindusername := "CN=bot,CN=users,DC=snowball,DC=com"
	bindpassword := LDAP_PWD
	err = l.Bind(bindusername, bindpassword)

	if err != nil {
		fmt.Println(err)
		return err
	}
	searchRequest := ldap.NewSearchRequest(
		//"ou=users,dc=snowballfinance,dc=com",
		"DC=snowball,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(cn=%s)(project~=rolling))", username),
		[]string{"rolling"},
		nil,
	)

	sr, searchErr := l.Search(searchRequest)
	if searchErr != nil {
		return searchErr
	}
	if len(sr.Entries) != 1 {
		return errors.New("无此用户")
	}

	userdn := sr.Entries[0].DN
	err = l.Bind(userdn, password)
	if err != nil {
		return err
	}

	return nil
}

// generateToken generate a token for a user
func (u *User) generateToken() {
	b := make([]byte, 18)
	rand.Read(b)
	u.Token = hex.EncodeToString(b)
}

// GetInfo get groups and projects info with this user
func (u *User) GetInfo() error {
	// get projects info
	projects, err := project.GetAllProjects()
	if err != nil {
		return err
	}
	// TODO merge projects from GitLab
	// get user's projects from GitLab
	userGitlabProjs, getProjsErr := gitlab.GetProjectsByUser(u.Name)
	if getProjsErr != nil {
		return getProjsErr
	}
	userGitlabGroups, getGroupsErr := gitlab.GetGroupsByUser(u.Name)
	if getGroupsErr != nil {
		return getGroupsErr
	}
	u.Groups = userGitlabGroups
	u.Projects = make([]string, 0)
	for _, proj := range projects {
		shouldAdd := false
		for _, group := range userGitlabGroups {
			if strings.Contains(proj.GitURL, ":"+group+"/") ||
				group == "ops" || group == "qa" {
				shouldAdd = true
				break
			}
		}
		for _, projName := range userGitlabProjs {
			// Sometimes projName different from proj.Name, but
			// they are the same project. So we check GitURL.
			if strings.Contains(strings.ToLower(proj.GitURL),
				strings.ToLower(projName)) {
				shouldAdd = true
				break
			}
		}
		if shouldAdd {
			u.Projects = append(u.Projects, proj.Name)
		}
	}
	sort.Strings(u.Projects)
	return nil
}

// Store user info and token into mysql
func (u *User) Store() error {
	u.generateToken()
	q := orm.MySQL.Where(User{Name: u.Name}).FirstOrCreate(u)
	// FirstOrUpdate
	//orm.MySQL.Where(User{Name: u.Name}).Assign(User{Projects:xx,Groups:xx}) .FirstOrCreate(u)
	return q.Error
}

// GetUserByToken check token and return user info
func GetUserByToken(token string) (u *User, err error) {
	u = &User{Token: token}
	return u, orm.MySQL.First(&u, "token = ?", token).Error
}