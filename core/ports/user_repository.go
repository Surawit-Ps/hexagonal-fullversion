package ports

import "hexagonal2/core/entity"

type UserRepository interface{
	GetUsers()([]entity.User,error)
	GetUser(string)(*entity.User,error)
	AddUser(entity.User)error
	GetUserByEmail(string)(*entity.User,error)
}