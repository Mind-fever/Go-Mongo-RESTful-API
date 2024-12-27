package responses

type UserInfo struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Codigo   string `json:"codigo"`
	Rol      string `json:"rol"`
}
