package service

import (
	"fmt"
	"hexagonal2/core/entity"
	"hexagonal2/core/middleware"
	"hexagonal2/core/ports"
	e "hexagonal2/pkg/errors"
)

type userService struct{
	repo ports.UserRepository
}

func NewUserService (repo ports.UserRepository) *userService {
	return &userService{repo: repo}
}

func (r userService) GetAllUser()([]entity.UserRes,error){
	user,err := r.repo.GetUsers()
	if err != nil{
		fmt.Print(err)
		return nil,err
	}
	var usRes []entity.UserRes
	for _,u := range user{
		usResp := entity.UserRes{
			Name: u.Name,
			LastName: u.LastName,
			Email: u.Email,
			Tel: u.Tel,
		}
		usRes = append(usRes,usResp)
	}
	return usRes,nil
}

func (r userService) GetUser(id string)(*entity.UserRes,error){
	user,err := r.repo.GetUser(id)
	if err != nil{
		fmt.Print(err)
		return nil,err
	}
	usResp := entity.UserRes{
			Name: user.Name,
			LastName: user.LastName,
			Email: user.Email,
			Tel: user.Tel,
		}
	
	return &usResp,nil
	
}

func (r userService) AddUser(p entity.User)error{
	err := r.repo.AddUser(p)
	if err != nil{
		fmt.Print(err)
		return err
	}
	return nil
}

func (r userService) Login(email string,password string)(string,error){
	user,err := r.repo.GetUserByEmail(email)
	if err != nil{
		return "",err
	}
	ok := middleware.CheckPasswordHash([]byte(password), []byte(user.Password))
	if !ok {
		return "", e.ErrInvalidCredentials
	}

	jwtWrapper := middleware.JwtWrapper{
		SecretKey:       "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
		Issuer:          "authService",
		ExpirationHours: 24,
	}

	token, err := jwtWrapper.GenerateToken(user.Id, user.Role)
	if err != nil {
		return "", e.ErrInternalServer
	}

	return token,nil
}