package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database/stores"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoom(t *testing.T) {

	t.Cleanup(cleanup)
	uc := NewUsecase(db)

	type args struct {
		req CreateRoomRequest
	}
	configPolymers := map[string]any{
		"Elements":    map[string]int{},
		"Checks":      polymers.Checks{},
		"TimerInt":    0,
		"Unicast":     nil,
		"Broadcast":   nil,
		"MaxPlayers":  0,
		"IsAutoCheck": false,
	}

	tests := []struct {
		name    string
		args    args
		want    *entities.Room
		wantErr enerr.Kind
	}{
		{
			name: "Create Room",
			args: args{
				req: CreateRoomRequest{
					Name:         "test_room",
					Type:         "polymers",
					EngineConfig: configPolymers,
				},
			},
			want:    &entities.Room{Name: "test_room"},
			wantErr: 0,
		},
		{
			name: "Create Room duplicate should error",
			args: args{
				req: CreateRoomRequest{
					Name:         "test_room",
					Type:         "polymers",
					EngineConfig: configPolymers,
				},
			},
			want:    &entities.Room{Name: "test_room"},
			wantErr: enerr.Exist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.CreateRoom(tt.args.req, MockLogger)
			if err != nil {
				if tt.wantErr != 0 {
					t.Logf("Error: %+v", err)

					assert.Equal(t, true, enerr.KindIs(tt.wantErr, err), "Kind should be equal")
					return
				}
				t.Log(err)
				t.Fatalf("Login() unexpected error = %v", err)
				return
			}
			// do not assert engine
			got.Engine = nil
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSubscribeToRoom(t *testing.T) {
	t.Cleanup(cleanup)

	memStreams := stores.NewStreamStore()

	var userRepo *repository.UserRepository = repository.NewUserRepo(db, memStreams)
	var user = &entities.User{ID: 1, Apikey: "api", Name: "test_user"}

	_, err := userRepo.CreateUser(database.InsertUserParams{
		Name:   user.Name,
		Apikey: user.Apikey,
		Room:   sql.NullString{},
		Role:   int64(common.Player_Role),
	})

	uc := NewUsecase(db)

	if err != nil {
		t.Fatal("Failed init")
	}
	_, err = uc.CreateRoom(CreateRoomRequest{
		Name: "test_room",
		Type: "polymers",
		EngineConfig: map[string]any{
			"Elements":    map[string]int{},
			"Checks":      polymers.Checks{},
			"TimerInt":    0,
			"Unicast":     nil,
			"Broadcast":   nil,
			"MaxPlayers":  2,
			"IsAutoCheck": false,
		},
	}, MockLogger)
	if err != nil {
		t.Fatal("Failed init")
	}

	t.Cleanup(cleanup)

	tests := []struct {
		name     string
		roomName string
		user     *entities.User
		wantErr  bool
	}{
		{
			name:     "Subscribe to new room",
			roomName: "test_room",
			user:     user,
			wantErr:  false,
		},
		{
			name:     "Subscribe to the same room again",
			roomName: "test_room",
			user:     user,
			wantErr:  true, // предполагаем, что повторная подписка вызывает ошибку
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.SubscribeToRoom(context.TODO(), tt.roomName, tt.user.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeToRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
			user, err := userRepo.GetUserByApikey(user.Apikey)
			if err != nil {
				t.Fatalf("failed getting user")
			}
			assert.Equal(t, tt.roomName, user.Room)
		})
	}
}

func TestGetSubscribersFromRoom(t *testing.T) {
	t.Cleanup(cleanup)
	memStreams := stores.NewStreamStore()
	var userRepo *repository.UserRepository = repository.NewUserRepo(db, memStreams)

	var user1 = &entities.User{ID: 1, Apikey: "api_11", Name: "test_user1"}

	us1, err := userRepo.CreateUser(database.InsertUserParams{
		Name:   user1.Name,
		Apikey: user1.Apikey,
		Room:   sql.NullString{},
		Role:   int64(common.Player_Role),
	})
	if err != nil {
		t.Fatalf("Failed init: %s", err.Error())
	}
	var user2 = &entities.User{ID: 2, Apikey: "api_12", Name: "test_user2"}

	us2, err := userRepo.CreateUser(database.InsertUserParams{
		Name:   user2.Name,
		Apikey: user2.Apikey,
		Room:   sql.NullString{},
		Role:   int64(common.Player_Role),
	})
	if err != nil {
		t.Fatalf("Failed init: %s", err.Error())
	}

	uc := NewUsecase(db)

	_, err = uc.CreateRoom(CreateRoomRequest{
		Name: "test_room",
		Type: "polymers",
		EngineConfig: map[string]any{
			"Elements":    map[string]int{},
			"Checks":      polymers.Checks{},
			"TimerInt":    0,
			"Unicast":     nil,
			"Broadcast":   nil,
			"MaxPlayers":  2,
			"IsAutoCheck": false,
		},
	}, MockLogger)
	if err != nil {
		t.Fatalf("Failed init: %s", err.Error())
	}
	err = uc.SubscribeToRoom(context.TODO(), "test_room", us1.ID)
	assert.NoError(t, err)

	err = uc.SubscribeToRoom(context.TODO(), "test_room", us2.ID)
	assert.NoError(t, err)

	roomSubs, err := uc.UserRepo.GetRoomSubscribers("test_room")
	if err != nil {
		t.Error("failed getting roomSubs")
	}

	assert.NotEmpty(t, roomSubs)
	roomSubsIds := make([]entities.ID, 0)
	for _, v := range roomSubs {
		roomSubsIds = append(roomSubsIds, v.ID)
	}

	assert.Contains(t, roomSubsIds, us1.ID)
	assert.Contains(t, roomSubsIds, us2.ID)
}
