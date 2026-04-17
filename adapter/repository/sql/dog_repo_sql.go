package sql

import ("hexagonal2/core/entity"
"gorm.io/gorm"
"github.com/google/uuid"
e "hexagonal2/pkg/errors")

type dogsRepositoryDB struct{
	db *gorm.DB
}

type dogs struct{
	Id string `gorm:"primaryKey"`
	Name string
	Age uint
	Colour string
	UserID string
}

func NewDogsRepositoryDB(db *gorm.DB) *dogsRepositoryDB{
	return &dogsRepositoryDB{db:db}
}

func EnToGorm(d entity.Dogs)dogs{
	return dogs{
		Id: d.Id,
		Name: d.Name,
		Age: d.Age,
		Colour: d.Colour,
		UserID: d.UserID,
	}
} 

func GormToEn(d dogs)entity.Dogs{
	return entity.Dogs{
		Id: d.Id,
		Name: d.Name,
		Age: d.Age,
		Colour: d.Colour,
		UserID: d.UserID,
	}
}

func (r dogsRepositoryDB)GetDogs()([]entity.Dogs,error){
	var dogs []dogs
	result := r.db.Find(&dogs)
	if result.Error != nil {
		return nil, e.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return nil, e.ErrInternalServer
	}
	var dogEntities []entity.Dogs
	for _, d := range dogs {
		dogEntities = append(dogEntities, GormToEn(d))
	}

	return dogEntities,nil

}

func (r dogsRepositoryDB)GetADogs(id string)(*entity.Dogs,error){
	var dog dogs
	result := r.db.Find(&dog,"id = ? OR human_id = ?",id,id)
	if result.Error != nil{
		return nil,e.ErrDogNotFound
	}
	edog := GormToEn(dog)
	return &edog,nil
}

func(r dogsRepositoryDB)AddDog(d entity.Dogs,userID string)error{
	d.Id = uuid.New().String()
	d.UserID = userID
	dog := EnToGorm(d)
	result := r.db.Create(&dog)
	if result.Error != nil{
		return e.ErrInternalServer
	}
	return nil
}