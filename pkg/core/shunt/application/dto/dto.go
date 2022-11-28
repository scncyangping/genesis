package dto

type UserDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Type     uint8  `json:"type"`
	Status   uint8  `json:"status"`
	Token    string `json:"token"`
}
