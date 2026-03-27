package ports

import("hgo/core/entity")

type HumanServices interface{
	GetAllUser()([]entity.HumanRes,error)
	GetUser(string)(*entity.HumanRes,error)
	AddUser(entity.Humans)error
}

