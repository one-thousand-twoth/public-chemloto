package hub

import (
	"reflect"
	"testing"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
)

func Test_usersState_Add(t *testing.T) {
	type args struct {
		user *User
	}
	user := NewUser("Andrei", "test_apikey", "uuid_conn", common.Player_Role, []string{"default"})
	tests := []struct {
		name    string
		rs      *usersState
		args    args
		wantErr bool
	}{
		{
			name: "Add single user",
			rs:   &usersState{state: make(map[string]*User)},
			args: args{user: user},
		},
		{
			name:    "Add same user",
			rs:      &usersState{state: map[string]*User{user.Name: user}},
			args:    args{user: user},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rs.Add(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("usersState.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.rs.state, map[string]*User{"Andrei": user}) {
				t.Errorf("usersState.Add() have = %v, want %v", tt.rs.state, tt.wantErr)
			}
		})
	}
}
