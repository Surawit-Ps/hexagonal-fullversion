package entity

type Dogs struct{
	Id string 
	Name string
	Age uint
	Colour string
	HumanID string
}

type DogRes struct{
	Name string `json:"name"`
	Age uint `json:"age"`
	HumanID string `json:"human_id"`
}