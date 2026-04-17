package ports

import "hexagonal2/core/entity"

type UserRepository interface{
	GetUsers()([]entity.User,error)
	GetUser(string)(*entity.User,error)
	AddUser(entity.User)error
	GetUserByEmail(string)(*entity.User,error)
}

type UserServices interface{
	GetAllUser()([]entity.UserRes,error)
	GetUser(string)(*entity.UserRes,error)
	AddUser(entity.User)error
	Login(string,string)(string,error)
}