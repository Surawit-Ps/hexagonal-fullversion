package sql

import ("hexagonal2/core/entity"
"gorm.io/gorm"
"github.com/google/uuid"
"hexagonal2/core/middleware"
e "hexagonal2/pkg/errors"
"strconv"
"fmt")

type userRepositoryDB struct{
	db *gorm.DB
}

type User struct{
	Id string `gorm:"primaryKey"`
	UserID string
	Name string 
	LastName string 
	Age int 
	Email string 
	Tel string 
	Password string
	Role string // Admin, User
}

func NewUserRepositoryDB(db *gorm.DB) userRepositoryDB {
	return userRepositoryDB{db: db}
}

func userEnToGorm(u entity.User)User{
	return User{
		Id: u.Id,
		UserID: u.UserID,
		Name: u.Name,
		LastName: u.LastName,
		Age: u.Age,
		Email: u.Email,
		Tel: u.Tel,
		Password: u.Password,
		Role: u.Role,
	}
}

func userGormToEn(u User)entity.User{
	return entity.User{
		Id: u.Id,
		UserID: u.UserID,
		Name: u.Name,
		LastName: u.LastName,
		Age: u.Age,
		Email: u.Email,
		Tel: u.Tel,
		Password: u.Password,
		Role: u.Role,
	}
}
    // GetPeoples()([]entity.Humans,error)
	// GetPerson(id string)(*entity.Humans,error)
	// AddPerson(p entity.Humans)error
func (r userRepositoryDB) GetUsers()([]entity.User,error){
	var pe []User
	result := r.db.Find(&pe)
	if  result.Error != nil{
		return nil,e.ErrUserNotFound
	}
	var peo []entity.User
	for _, u := range pe{
		peo = append(peo, userGormToEn(u))
	}
	return  peo,nil
}

func(r userRepositoryDB)GetUser(id string)(*entity.User,error){
	var pe User
	result := r.db.Find(&pe,"id = ?",id)
	if result.Error != nil{
		return nil,e.ErrUserNotFound
	}
	peo := userGormToEn(pe)
	return &peo,nil
}

func (r userRepositoryDB) AddUser(p entity.User)error{
	p.Id = uuid.New().String()
	p.UserID = r.GenerateAccountID()
	Password, err := middleware.HashPassword(p.Password)
	if err != nil {
		return e.ErrInternalServer
	}
	p.Password = Password
	hu := userEnToGorm(p)

	result := r.db.Create(&hu)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r userRepositoryDB) GetUserByEmail(email string)(*entity.User,error){
	var pe User
	result := r.db.Find(&pe,"email = ?",email)
	if result.Error != nil{
		return nil,e.ErrUserNotFound
	}
	m := userGormToEn(pe)
	return &m,nil
}

func (r userRepositoryDB) GenerateAccountID() string {
	var user entity.User
	seq := 477631 // default start

	err := r.db.Last(&user).Error
	if err == nil && user.UserID != "" {
		if n, parseErr := strconv.Atoi(user.UserID); parseErr == nil {
			seq = n + 1
		}
	}

	return fmt.Sprintf("%06d", seq)
}