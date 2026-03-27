package ports

import("hgo/core/entity")

type DogServices interface{
	GetAllDogs()([]entity.DogRes,error)
	GetDog(string)(*entity.DogRes,error)
	AddDog(entity.Dogs,string)error
}