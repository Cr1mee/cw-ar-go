package controllers

import (
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type RoomController struct {
	roomService app.RoomService
}

func NewRoomController(rs app.RoomService) RoomController {
	return RoomController{
		roomService: rs,
	}
}

func (c RoomController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		room, err := requests.Bind(r, requests.RoomRequest{}, domain.Room{})
		if err != nil {
			log.Printf("OrganizationController: %s", err)
			BadRequest(w, err)
			return
		}

		room, err = c.roomService.Save(room, user.Id)
		if err != nil {
			log.Printf("OrganizationController: %s", err)
			if err.Error() == "access denied" {
				Forbidden(w, err)
			} else {
				InternalServerError(w, err)
			}
			return
		}

		var roomDto resources.RoomDto
		Created(w, roomDto.DomainToDto(room))
	}
}
