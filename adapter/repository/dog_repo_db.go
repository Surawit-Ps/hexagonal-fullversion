package repository

import ("hgo/core/entity"
"gorm.io/gorm"
"github.com/google/uuid")

type dogsRepositoryDB struct{
	db *gorm.DB
}

type DogsModel struct{
	Id string `gorm:"primaryKey"`
	Name string
	Age uint
	Colour string
	HumanID string
}

func NewDogsRepositoryDB(db *gorm.DB) *dogsRepositoryDB{
	return &dogsRepositoryDB{db:db}
}

func EnToGorm(d entity.Dogs)DogsModel{
	return DogsModel{
		Id: d.Id,
		Name: d.Name,
		Age: d.Age,
		Colour: d.Colour,
		HumanID: d.HumanID,
	}
} 

func GormToEn(d DogsModel)entity.Dogs{
	return entity.Dogs{
		Id: d.Id,
		Name: d.Name,
		Age: d.Age,
		Colour: d.Colour,
		HumanID: d.HumanID,
	}
}

func (r dogsRepositoryDB)GetDogs()([]entity.Dogs,error){
	var dogs []DogsModel
	result := r.db.Find(&dogs)
	if result.Error != nil {
		return nil, result.Error
	}
	var dogEntities []entity.Dogs
	for _, d := range dogs {
		dogEntities = append(dogEntities, GormToEn(d))
	}

	return dogEntities,nil

}

func (r dogsRepositoryDB)GetADogs(id string)(*entity.Dogs,error){
	var dog DogsModel
	result := r.db.Find(&dog,"id = ?",id)
	if result.Error != nil{
		return nil,result.Error
	}
	edog := GormToEn(dog)
	return &edog,nil
}

func(r dogsRepositoryDB)AddDog(d entity.Dogs,humanID string)error{
	d.Id = uuid.New().String()
	d.HumanID = humanID
	dog := EnToGorm(d)
	result := r.db.Create(&dog)
	if result.Error != nil{
		return result.Error
	}
	return nil
}