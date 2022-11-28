package cmd

type LoginCmd struct {
	Name string `json:"name"` // 用户名
	Pwd  string `json:"pwd"`  // 密码
}

type RegisterCmd struct {
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Pwd      string `json:"pwd"`
}
