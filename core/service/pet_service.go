package service

import (
"hexagonal2/core/ports"
"hexagonal2/core/entity")

type petService struct{
	repo ports.PetRepository
}

func NewPetService(repo ports.PetRepository) petService{
	return petService{repo: repo}
}

func(r petService)GetAllPets()([]entity.PetRes,error){
	pet,err := r.repo.GetPets()
	if err != nil{
		return nil,err
	}
	var petRes []entity.PetRes
	for _,p := range pet{
		pt := entity.PetRes{
			PetID: p.PetID,
			OwnerID: p.OwnerID,
			Name: p.Name,
			Species: p.Species,
			Breed: p.Breed,
			Age: p.Age,
			Weight: p.Weight,
		}
		petRes = append(petRes, pt)
	}
	return petRes,nil
}

func(r petService)GetPet(id string)(*entity.PetRes,error){
	pet,err := r.repo.GetAPet(id)			
	if err != nil{
		return nil,err
	}
	pt := entity.PetRes{
			PetID: pet.PetID,
			OwnerID: pet.OwnerID,
			Name: pet.Name,
			Species: pet.Species,
			Breed: pet.Breed,
			Age: pet.Age,
			Weight: pet.Weight,
		}
	return &pt,nil
}

func (r petService)AddPet(p entity.Pet,h string)error{       
	err := r.repo.AddPet(p,h)
	if err != nil{
		return err	
	}
	return nil
}	

func (r petService)UpdatePet(id string,p entity.Pet)error{
	err := r.repo.UpdatePet(id,p)
	if err != nil{
		return err
	}
	return nil
}

func (r petService)DeletePet(id string)error{
	err := r.repo.DeletePet(id)
	if err != nil{
		return err
	}
	return nil
}

	