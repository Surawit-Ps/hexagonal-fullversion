package ports

import "hexagonal2/core/entity"

type RoomRepository interface {
	GetAllRooms() ([]entity.RoomType, error)
	GetARoom(id string) (entity.RoomType, error)
	AddRoom(room entity.RoomType) (string, error)
	UpdateRoom(id string, room entity.RoomType) error
	DeleteRoom(id string) error
}

type RoomService interface {
	GetAllRooms() ([]entity.RoomType, error)
	GetARoom(id string) (entity.RoomType, error)
	AddRoom(room entity.RoomType) (string, error)
	UpdateRoom(id string, room entity.RoomType) error
	DeleteRoom(id string) error
}