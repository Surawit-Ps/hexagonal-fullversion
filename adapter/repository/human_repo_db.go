package repository

import ("hgo/core/entity"
"gorm.io/gorm"
"github.com/google/uuid")

type humanRepositoryDB struct{
	db *gorm.DB
}

type HumanDB struct{
	Id string `gorm:"primaryKey"`
	Name string 
	LastName string 
	Age int 
	Email string 
	Tel string 
}

func NewHumanReposityDB(db *gorm.DB) humanRepositoryDB {
	return humanRepositoryDB{db: db}
}

func humanEnToGorm(h entity.Humans)HumanDB{
	return HumanDB{
		Id: h.Id,
		Name: h.Name,
		LastName: h.LastName,
		Age: h.Age,
		Email: h.Email,
		Tel: h.Tel,
	}
}

func humanGormToEn(h HumanDB)entity.Humans{
	return entity.Humans{
		Id: h.Id,
		Name: h.Name,
		LastName: h.LastName,
		Age: h.Age,
		Email: h.Email,
		Tel: h.Tel,
	}
}
    // GetPeoples()([]entity.Humans,error)
	// GetPerson(id string)(*entity.Humans,error)
	// AddPerson(p entity.Humans)error
func (r humanRepositoryDB) GetPeoples()([]entity.Humans,error){
	var pe []HumanDB
	result := r.db.Find(&pe)
	if  result.Error != nil{
		return nil,result.Error
	}
	var peo []entity.Humans
	for _, hu := range pe{
		peo = append(peo, humanGormToEn(hu))
	}
	return  peo,nil
}

func(r humanRepositoryDB)GetPerson(id string)(*entity.Humans,error){
	var pe HumanDB
	result := r.db.Find(&pe,"id = ?",id)
	if result.Error != nil{
		return nil,result.Error
	}
	peo := humanGormToEn(pe)
	return &peo,nil
}

func (r humanRepositoryDB) AddPerson(p entity.Humans)error{
	p.Id = uuid.New().String()
	hu := humanEnToGorm(p)

	result := r.db.Create(&hu)
	if result.Error != nil {
		return result.Error
	}
	return nil
}