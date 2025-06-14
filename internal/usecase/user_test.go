package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/stretchr/testify/assert"
)

var buff = &bytes.Buffer{}
var MockLogger = slog.New(slog.NewTextHandler(buff, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
	return a
}}))

func TestLogin(t *testing.T) {
	t.Cleanup(cleanup)

	uc := NewUsecase(db)

	type args struct {
		req LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *LoginResponse
		wantErr *enerr.ApplicationError
	}{
		{
			name: "Create Player",
			args: args{
				req: LoginRequest{"Player", ""},
			},
			want: &LoginResponse{
				Token: "",
				Role:  common.Player_Role,
				Error: []string{},
			},
			wantErr: nil,
		},
		{
			name: "Create Duplicate Player",
			args: args{
				req: LoginRequest{
					Name: "Player",
					Code: "",
				},
			},
			want:    nil,
			wantErr: enerr.E(enerr.Exist),
		},
		{
			name: "Create Admin",
			args: args{
				req: LoginRequest{
					Name: "Admin",
					Code: "valid",
				},
			},
			want: &LoginResponse{
				Token: "",
				Role:  common.Admin_Role,
				Error: []string{},
			},
			wantErr: nil,
		},
		{
			name: "Create Admin - invalid code",
			args: args{
				req: LoginRequest{
					Name: "Admin",
					Code: "invalid",
				},
			},
			want: &LoginResponse{
				Token: "",
				Role:  common.Admin_Role,
				Error: []string{},
			},
			wantErr: enerr.E(enerr.Validation),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := uc.Login(MockLogger, tt.args.req, "valid")
			if err != nil {
				if tt.wantErr != nil {
					fmt.Printf("Error: %+v", err)
					assert.Equal(t, true, enerr.KindIs(tt.wantErr.Kind, err))
					return
				}

				t.Fatalf("Login() unexpected error = %v", err)
				return
			}
			// do not compare random token
			// got.Token = ""
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Login() = %v, want %v", got, tt.want)
			// }
			assert.Equal(t, got.Role, tt.want.Role)
			assert.Equal(t, got.Error, tt.want.Error)

			savedUser, err := uc.UserRepo.GetUserByApikey(got.Token)
			assert.NoError(t, err)
			assert.NotNil(t, savedUser.MessageChan)
		})
	}
}

func TestPatchUser(t *testing.T) {

	t.Cleanup(cleanup)

	// db := sqlite.MustInitDB()
	params := database.InsertUserParams{
		Name:   "TestUser",
		Apikey: "test_api",
		Room:   sql.NullString{},
		Role:   int64(common.Player_Role),
	}
	_, err := database.New(db).InsertUser(context.TODO(), params)
	if err != nil {
		panic(fmt.Sprintf("init failed: %s", err))
	}

	uc := NewUsecase(db)

	type args struct {
		req PatchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *entities.User
		wantErr bool
	}{
		{
			name: "Alter Role to Judge",
			args: args{
				req: PatchRequest{
					Name: params.Name,
					Role: common.Judge_Role,
				},
			},
			want: &entities.User{
				Name:   params.Name,
				Apikey: params.Apikey,
				Room:   "",
				Role:   common.Judge_Role,
			},
			wantErr: false,
		},
		{
			name: "Alter Role to Player",
			args: args{
				req: PatchRequest{
					Name: params.Name,
					Role: common.Player_Role,
				},
			},
			want: &entities.User{
				Name:   params.Name,
				Apikey: params.Apikey,
				Room:   "",
				Role:   common.Player_Role,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.PatchUserRole(context.TODO(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, err = uc.UserRepo.GetUserByApikey("test_api")
			if !tt.wantErr && err != nil {
				t.Error("shouldnt be error")
			}
		})
	}
}
