package ports

import "hgo/core/entity"

type HumanRepository interface{
	GetPeoples()([]entity.Humans,error)
	GetPerson(string)(*entity.Humans,error)
	AddPerson(entity.Humans)error
}