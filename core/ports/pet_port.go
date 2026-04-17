package ports

import("hexagonal2/core/entity")

type PetRepository interface{
	GetPets()([]entity.Pet,error)
	GetAPet(string)(*entity.Pet,error)
	AddPet(entity.Pet,string)error
	UpdatePet(string,entity.Pet)error
	DeletePet(string)error
}

type PetServices interface{
	GetAllPets()([]entity.PetRes,error)
	GetPet(string)(*entity.PetRes,error)
	AddPet(entity.Pet,string)error
	UpdatePet(string,entity.Pet)error
	DeletePet(string)error
}

