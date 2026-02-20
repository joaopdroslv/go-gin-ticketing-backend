package enums

type UserStatusID int64

const (
	Active            UserStatusID = 1
	Inactive          UserStatusID = 2
	PasswordCreation  UserStatusID = 3
	EmailConfirmation UserStatusID = 4
	DeletedAccount    UserStatusID = 5
)
