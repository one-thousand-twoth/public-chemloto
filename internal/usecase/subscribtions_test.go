package usecase

import (
	"testing"

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
		fn   func()
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
				fn: func() {
					panic("TODO")
				},
			},
			wantErr: false,
		},
		{
			name: "Create Channel duplicate should error",
			args: args{
				name: "test_channel",
				fn: func() {
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

	id, err := AddRegularChannel(repo, "test_channel", func() { results <- struct{}{} })
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
				id: id,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetRegularChannel(repo, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("GetRegularChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := results; ok != true {

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
