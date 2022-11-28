package shuntCore

type UserStatus uint8 // 用户状态

const (
	UserNormal UserStatus = 1 << iota
	UserFrozen
)

var UserStatusGenesis = map[UserStatus]string{
	UserNormal: "正常",
	UserFrozen: "冻结",
}
