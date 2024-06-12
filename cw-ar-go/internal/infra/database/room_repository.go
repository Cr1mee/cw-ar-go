package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const RoomsTableName = "rooms"

type room struct {
	Id             uint64     `db:"id,omitempty"`
	OrganizationId uint64     `db:"organizationId"`
	Name           string     `db:"	name"`
	Description    string     `db:"	description"`
	CreatedDate    time.Time  `db:"created_date"`
	UpdatedDate    time.Time  `db:"updated_date"`
	DeletedDate    *time.Time `db:"deleted_date"`
}

type RoomRepository interface {
	Save(dr domain.Room) (domain.Room, error)
	FindByOrgId(oId uint64) ([]domain.Room, error)
	FindById(id uint64) (domain.Room, error)
	Update(o domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomRepository struct {
	coll db.Collection
	sess db.Session
}

func NewRoomRepository(dbSession db.Session) RoomRepository {
	return roomRepository{
		coll: dbSession.Collection(RoomsTableName),
		sess: dbSession,
	}
}

func (r roomRepository) Save(dr domain.Room) (domain.Room, error) {
	rm := r.mapDomainToModel(dr)
	rm.CreatedDate, rm.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&rm)
	if err != nil {
		return domain.Room{}, err
	}
	dr = r.mapModelToDomain(rm)
	return dr, nil
}

func (r roomRepository) FindByOrgId(oId uint64) ([]domain.Room, error) {
	var rooms []room
	err := r.coll.Find(db.Cond{"organiztion_id": oId, "deleted_date": nil}).All(&rooms)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(rooms)
	return res, nil
}

func (r roomRepository) FindById(id uint64) (domain.Room, error) {
	var rm room
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&rm)
	if err != nil {
		return domain.Room{}, err
	}
	o := r.mapModelToDomain(rm)
	return o, nil
}

func (r roomRepository) Update(o domain.Room) (domain.Room, error) {
	rm := r.mapDomainToModel(o)
	rm.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": rm.Id, "deleted_date": nil}).Update(&rm)
	if err != nil {
		return domain.Room{}, err
	}
	o = r.mapModelToDomain(rm)
	return o, nil
}

func (r roomRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r roomRepository) mapDomainToModel(d domain.Room) room {
	return room{
		Id:             d.Id,
		OrganizationId: d.OrganizationId,
		Name:           d.Name,
		Description:    d.Description,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomain(d room) domain.Room {
	return domain.Room{
		Id:             d.Id,
		OrganizationId: d.OrganizationId,
		Name:           d.Name,
		Description:    d.Description,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomainCollection(rs []room) []domain.Room {
	var rooms []domain.Room
	for _, rm := range rs {
		dr := r.mapModelToDomain(rm)
		rooms = append(rooms, dr)
	}
	return rooms
}
