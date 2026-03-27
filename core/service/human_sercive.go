package service

import ("fmt"
"hgo/core/ports"
"hgo/core/entity")

type humanService struct{
	repo ports.HumanRepository
	
}

func NewHumanServive (repo ports.HumanRepository) *humanService {
	return &humanService{repo: repo}
}

func (r humanService) GetAllUser()([]entity.HumanRes,error){
	human,err := r.repo.GetPeoples()
	if err != nil{
		fmt.Print(err)
		return nil,err
	}
	var huRes []entity.HumanRes
	for _,h := range human{
		huResp := entity.HumanRes{
			Name: h.Name,
			LastName: h.LastName,
			Email: h.Email,
			Tel: h.Tel,
		}
		huRes = append(huRes,huResp)
	}
	return huRes,nil
}

func (r humanService) GetUser(id string)(*entity.HumanRes,error){
	human,err := r.repo.GetPerson(id)
	if err != nil{
		fmt.Print(err)
		return nil,err
	}
	huResp := entity.HumanRes{
			Name: human.Name,
			LastName: human.LastName,
			Email: human.Email,
			Tel: human.Tel,
		}
	return &huResp,nil
	
}

func (r humanService) AddUser(p entity.Humans)error{
	err := r.repo.AddPerson(p)
	if err != nil{
		fmt.Print(err)
		return err
	}
	return nil
}