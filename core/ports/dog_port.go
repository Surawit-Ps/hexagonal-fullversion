package ports
import "hexagonal2/core/entity"

type DogsRepository interface{
	GetDogs()([]entity.Dogs,error)
	GetADogs(string)(*entity.Dogs,error)
	AddDog(entity.Dogs,string)error
}

type DogServices interface{
	GetAllDogs()([]entity.DogRes,error)
	GetDog(string)(*entity.DogRes,error)
	AddDog(entity.Dogs,string)error
}