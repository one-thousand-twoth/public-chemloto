package usecase

import (
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/polymers"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/anrew1002/Tournament-ChemLoto/internal/hub/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoom(t *testing.T) {

	// var db *sql.DB = sqlite.MustInitDB()
	repo := repository.NewRoomRepo(db)
	t.Cleanup(cleanup)

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
			got, err := CreateRoom(repo, tt.args.req, MockLogger)
			if err != nil {
				if tt.wantErr != 0 {
					t.Logf("Error: %+v", err)
					assert.Equal(t, true, enerr.KindIs(tt.wantErr, err))
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
