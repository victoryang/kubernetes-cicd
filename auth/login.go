package auth

// LoginParam is used for login
type LoginParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login check if user has authority and generage token
func Login(param *LoginParam) (user *User, err error) {
	user = &User{}
	authErr := user.AuthByLdap(param.Username, param.Password)
	if authErr != nil {
		return user, authErr
	}
	user.Name = param.Username
	err = user.GetInfo()
	if err != nil {
		return
	}
	err = user.Store()
	return
}