package usecase

import (
	"database/sql"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sqlite/sqlitetest"
)

var db, cleanup = sqlitetest.GetTestDatabase()

func TestAddRegularChannel(t *testing.T) {

	var repo *repository.ChannelsRepository = repository.NewChannelsRepo(db)

	t.Cleanup(cleanup)

	type args struct {
		name string
		fn   func(chan common.Message)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create Channel",
			args: args{
				name: "test_channel",
				fn: func(chan common.Message) {
					panic("TODO")
				},
			},
			wantErr: false,
		},
		{
			name: "Create Channel duplicate should error",
			args: args{
				name: "test_channel",
				fn: func(chan common.Message) {
					panic("TODO")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := AddRegularChannel(repo, tt.args.name, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("AddRegularChannel() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestGetRegularChannel(t *testing.T) {

	// var db *sql.DB = sqlite.MustInitDB()

	var repo *repository.ChannelsRepository = repository.NewChannelsRepo(db)

	t.Cleanup(cleanup)

	results := make(chan struct{})

	channel, err := AddRegularChannel(repo, "test_channel", func(chan common.Message) { results <- struct{}{} })
	if err != nil {
		t.Fatal("fail init add", err)
	}

	type args struct {
		id entities.ID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Get",
			args: args{
				id: channel.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetRegularChannel(repo, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("GetRegularChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscribeToChannel(t *testing.T) {

	// var db *sql.DB = sqlite.MustInitDB()

	var repo *repository.ChannelsRepository = repository.NewChannelsRepo(db)

	t.Cleanup(cleanup)

	// id, err := AddRegularChannel(repo, "test_channel", func() {})
	// if err != nil {
	// 	t.Fatal("fail init add", err)
	// }

	type args struct {
		id   entities.ID
		user entities.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Subscribe",
			args: args{
				id: 1,
				user: entities.User{
					ID: 10,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SubscribeToChannel(repo, tt.args.id, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeToChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscribeToRoom(t *testing.T) {
	t.Cleanup(cleanup)
	var channelsRepo *repository.ChannelsRepository = repository.NewChannelsRepo(db)
	var roomRepo *repository.RoomRepository = repository.NewRoomRepo(db)
	var userRepo *repository.UserRepository = repository.NewUserRepo(db)
	var user = &entities.User{ID: 1, Name: "test_user"}

	_, err := userRepo.CreateUser(database.InsertUserParams{
		Name:   user.Name,
		Apikey: "api",
		Room:   sql.NullString{},
		Role:   int64(common.Player_Role),
	})
	if err != nil {
		t.Fatal("Failed init")
	}
	_, err = CreateRoom(roomRepo, CreateRoomRequest{
		Name: "test_room",
		Type: "polymers",
		EngineConfig: map[string]any{
			"Elements":    map[string]int{},
			"Checks":      polymers.Checks{},
			"TimerInt":    0,
			"Unicast":     nil,
			"Broadcast":   nil,
			"MaxPlayers":  0,
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
			err := SubscribeToRoom(channelsRepo, roomRepo, tt.roomName, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeToRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
