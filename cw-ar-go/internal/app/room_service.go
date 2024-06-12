package app

import (
	"errors"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type RoomService interface {
	Save(r domain.Room, uId uint64) (domain.Room, error)
	FindByOrgId(oId uint64) ([]domain.Room, error)
	Find(id uint64) (interface{}, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomService struct {
	roomRepo database.RoomRepository
	orgRepo  database.OrganizationRepository
}

func NewRoomService(rr database.RoomRepository, or database.OrganizationRepository) RoomService {
	return roomService{
		roomRepo: rr,
		orgRepo:  or,
	}
}

func (s roomService) Save(r domain.Room, uId uint64) (domain.Room, error) {
	org, err := s.orgRepo.FindById(r.OrganizationId)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	if org.UserId != uId {
		err = errors.New("access denied")
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	r, err = s.roomRepo.Save(r)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	return r, nil
}

func (s roomService) FindByOrgId(oId uint64) ([]domain.Room, error) {
	rooms, err := s.roomRepo.FindByOrgId(oId)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return nil, err
	}

	return rooms, nil
}

func (s roomService) Find(id uint64) (interface{}, error) {
	room, err := s.roomRepo.FindById(id)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return nil, err
	}

	return room, nil
}

func (s roomService) Update(r domain.Room) (domain.Room, error) {
	room, err := s.roomRepo.Update(r)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	return room, nil
}

func (s roomService) Delete(id uint64) error {
	err := s.roomRepo.Delete(id)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return err
	}

	return nil
}
