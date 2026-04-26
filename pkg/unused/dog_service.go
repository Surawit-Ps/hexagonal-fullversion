package unused

// import ("fmt"
// "hexagonal2/core/ports"
// "hexagonal2/core/entity")

// type dogService struct{
// 	repo ports.DogsRepository
// }

// func NewDogService(repo ports.DogsRepository) dogService{
// 	return dogService{repo: repo}
// }

// func(r dogService)GetAllDogs()([]entity.DogRes,error){
// 	dog,err := r.repo.GetDogs()
// 	if err != nil{
// 		fmt.Print(err)
// 		return nil,err
// 	}
// 	var dosRes []entity.DogRes
// 	for _,d := range dog{
// 		ds := entity.DogRes{
// 			ID: d.Id,
// 			Name: d.Name,
// 			Age: d.Age,
// 			UserID: d.UserID,
// 		}
// 		dosRes = append(dosRes, ds)
// 	}
// 	return dosRes,nil
// }

// func(r dogService)GetDog(id string)(*entity.DogRes,error){
// 	dog,err := r.repo.GetADogs(id)
// 	if err != nil{
// 		fmt.Print(err)
// 		return nil,err
// 	}
// 	ds := entity.DogRes{
// 			ID: dog.Id,
// 			Name: dog.Name,
// 			Age: dog.Age,
// 			UserID: dog.UserID,
// 		}
// 	return &ds,nil
// }

// func (r dogService)AddDog(d entity.Dogs,h string)error{
// 	err := r.repo.AddDog(d,h)
// 	if err != nil{
// 		fmt.Print(err)
// 		return err
// 	}
// 	return nil
// }