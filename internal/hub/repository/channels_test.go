package repository

import (
	"reflect"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_channelsState_Add(t *testing.T) {
	type args struct {
		channel    string
		connection string
	}
	state := &ChannelsState{state: make(map[string]mapset.Set[string], 0)}
	wantState := &ChannelsState{state: make(map[string]mapset.Set[string], 0)}
	wantState.state["default"] = mapset.NewThreadUnsafeSet("1033")
	tests := []struct {
		name     string
		channels *ChannelsState
		args     args
		want     map[string]mapset.Set[string]
	}{
		{
			"Add new should add",
			state,
			args{channel: "default", connection: "1033"},
			wantState.state,
		},
		{
			"Add existed should not add",
			state,
			args{channel: "default", connection: "1033"},
			wantState.state,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.channels.Add(tt.args.channel, tt.args.connection)
			if !reflect.DeepEqual(state.state, tt.want) {
				t.Fatalf("\n\nHave %+v,\n want %+v\n\n", state.state, wantState.state)
			}
		})
	}
}
