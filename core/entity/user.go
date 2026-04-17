package entity

type User struct{
	Id string 
	UserID string
	Name string 
	LastName string 
	Age int 
	Email string 
	Tel string 
	Password string
	Role string // Admin, User
}

type UserRes struct{
	ID string `json:"id"`
	Name string `json:"name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Tel string `json:"tel"`
}