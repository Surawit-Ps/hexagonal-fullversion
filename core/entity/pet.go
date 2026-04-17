package entity

import "time"

type Pet struct {
    ID        string
    PetID     string
    OwnerID   string
    Name      string
    Species   string // dog, cat
    Breed     string
    Age       int
    Weight    float64
    CreatedAt time.Time
}

type PetRes struct {
    PetID     string
	OwnerID   string
	Name      string
	Species   string // dog, cat
	Breed     string
	Age       int
	Weight    float64
}