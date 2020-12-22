package gitlab

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

var git *gitlab.Client

var token string
var gitURL string
var userAgent string

func init() {
	token = "xxxxx"
	gitURL = "http://xxxxxxx.com/api/v4"
	userAgent = "xxxx"
	git,_ = gitlab.NewClient(token, gitlab.WithBaseURL(gitURL))
}

func getUserIDByName(userName string) (int, error) {
	name := &userName
	users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{Username: name})
	if len(users) != 1 || err != nil {
		return 0, fmt.Errorf("Get user from gitlab error: username wrong or %v", err)
	}
	return users[0].ID, nil
}

// GetProjectsByUser return all projects this user has privilages in GitLab
func GetProjectsByUser(userName string) ([]string, error) {
	opt := gitlab.ListProjectsOptions{}
	opt.Simple = func(b bool) *bool { return &b }(true)
	projs, _, listPorjsErr := git.Projects.ListUserProjects(userName, &opt)

	if listPorjsErr != nil {
		return nil, fmt.Errorf("Get user's project from gitlab error:  %v", listPorjsErr)
	}

	projects := []string{}
	for _, proj := range projs {
		projects = append(projects, proj.Name)
	}
	return projects, nil
}

// GetGroupsByUser get a user's groups
func GetGroupsByUser(userName string) (res []string, err error) {
	id, err := getUserIDByName(userName)
	if err != nil {
		return nil, err
	}

	groups, _, listGroupsErr := git.Groups.ListGroups(&gitlab.ListGroupsOptions{})

	if listGroupsErr != nil {
		return nil, fmt.Errorf("Get groups info from gitlab error: %v", listGroupsErr)
	}

	for _, group := range groups {
		member, _, getMemberErr := git.GroupMembers.GetGroupMember(group.ID, id)
		if getMemberErr != nil {
			continue
		}
		if userName == member.Username {
			res = append(res, group.Path)
		}
	}

	return
}