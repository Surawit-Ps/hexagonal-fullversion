package ports

import("hexagonal2/core/entity")

type UserServices interface{
	GetAllUser()([]entity.UserRes,error)
	GetUser(string)(*entity.UserRes,error)
	AddUser(entity.User)error
	Login(string,string)(string,error)
}

