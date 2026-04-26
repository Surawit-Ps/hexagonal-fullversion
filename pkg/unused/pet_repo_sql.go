package unused

// import ("hexagonal2/core/entity"
// "gorm.io/gorm"
// "github.com/google/uuid"
// e "hexagonal2/pkg/errors"
// "time"
// "strings"
// "strconv"
// "fmt")

// type petRepositoryDB struct{
// 	db *gorm.DB
// }

// func NewPetRepositoryDB(db *gorm.DB) petRepositoryDB{
// 	return petRepositoryDB{db:db}
// }

// type Pet struct{
// 	ID string `gorm:"primaryKey"`
// 	OwnerID string
// 	Name string
// 	PetID string
// 	Species string
// 	Breed string	
// 	Age int
// 	Weight float64
// 	CreatedAt time.Time
// }

// func EnToGormPet(p entity.Pet)Pet{
// 	return Pet{
// 		ID: p.ID,
// 		PetID: p.PetID,
// 		OwnerID: p.OwnerID,
// 		Name: p.Name,
// 		Species: p.Species,	
// 		Breed: p.Breed,
// 		Age: p.Age,
// 		Weight: p.Weight,
// 		CreatedAt: p.CreatedAt,
// 	}
// }

// func GormToEnPet(p Pet)entity.Pet{
// 	return entity.Pet{
// 		ID: p.ID,
// 		OwnerID: p.OwnerID,
// 		Name: p.Name,
// 		PetID: p.PetID,
// 		Species: p.Species,
// 		Breed: p.Breed,
// 		Age: p.Age,
// 		Weight: p.Weight,
// 		CreatedAt: p.CreatedAt,
// 	}
// }

// func (r petRepositoryDB)GetPets()([]entity.Pet,error){
// 	var pets []Pet
// 	result := r.db.Find(&pets)
// 	if result.Error != nil {
// 		return nil, e.ErrInternalServer
// 	}
// 	if result.RowsAffected == 0 {
// 		return nil, e.ErrInternalServer
// 	}
// 	var petEntities []entity.Pet
// 	for _, p := range pets {
// 		petEntities = append(petEntities, GormToEnPet(p))
// 	}
// 	return petEntities, nil
// }

// func (r petRepositoryDB)GetAPet(id string)(*entity.Pet,error){
// 	var pet Pet
// 	result := r.db.Find(&pet,"id = ? OR owner_id = ? OR pet_id = ?",id,id,id)
// 	if result.Error != nil {
// 		return nil, e.ErrInternalServer
// 	}
// 	if result.RowsAffected == 0 {
// 		return nil, e.ErrInternalServer
// 	}
// 	petEn := GormToEnPet(pet)
// 	return &petEn, nil
// }



// func (r petRepositoryDB)AddPet(p entity.Pet,h string)error{
// 	p.ID = uuid.New().String()
// 	p.PetID = r.GenerateID()
// 	p.CreatedAt = time.Now()
	
// 	pet := EnToGormPet(p)
// 	result := r.db.Create(&pet)
// 	if result.Error != nil {
// 		return e.ErrInternalServer
// 	}
// 	return nil
// }

// func (r petRepositoryDB)UpdatePet(id string,p entity.Pet)error{
// 	var pet Pet
// 	result := r.db.Find(&pet,"id = ?",id)

// 	if result.Error != nil {
// 		return e.ErrInternalServer
// 	}
// 	if result.RowsAffected == 0 {
// 		return e.ErrInternalServer
// 	}

// 	pet.Name = p.Name
// 	pet.Species = p.Species
// 	pet.Breed = p.Breed
// 	pet.Age = p.Age
// 	pet.Weight = p.Weight
// 	result = r.db.Save(&pet)
// 	if result.Error != nil {
// 		return e.ErrInternalServer
// 	}
// 	return nil
// }

// func (r petRepositoryDB)DeletePet(id string)error{
// 	result := r.db.Delete(&Pet{},"id = ?",id)
// 	if result.Error != nil {
// 		return e.ErrInternalServer
// 	}
// 	if result.RowsAffected == 0 {
// 		return e.ErrInternalServer
// 	}
// 	return nil
// }

// func (r petRepositoryDB) GenerateID() string {
// 	var petID entity.Pet
// 	today := time.Now().Format("20060102")
// 	seq := 1
// 	err := r.db.Last(&petID)
// 	if err == nil {
// 		num := strings.Split(petID.PetID, "-")
// 		if len(num) == 2 {
// 			if n, parseErr := strconv.Atoi(num[1]); parseErr == nil {
// 				seq = n + 1
// 			}
// 		}
// 	}
// 	return fmt.Sprintf("PT%s-%03d", today, seq)
// }