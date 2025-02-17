package repository

import (
	"encoding/json"
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	mapset "github.com/deckarep/golang-set/v2"
)

func NewChannelState() *ChannelsState {
	return &ChannelsState{
		state: make(map[string]mapset.Set[string]),
		initFunctions: map[string]func(chan common.Message){
			"default": func(c chan common.Message) {},
		},
	}
}

// ChannelsState содержит карту со списками ID соединений подписанных на канал
type ChannelsState struct {
	// map key is a channel name
	state map[string]mapset.Set[string]
	// Функции каналов срабатывающие при первом подключении.
	// key: channel name
	initFunctions map[string]func(chan common.Message)
	mutex         sync.RWMutex
}

func (rs *ChannelsState) MarshalJSON() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.state)
}

func (channels *ChannelsState) GetChannelFunc(channel string) (func(chan common.Message), bool) {
	channels.mutex.Lock()
	defer channels.mutex.Unlock()
	fun, ok := channels.initFunctions[channel]
	return fun, ok
}
func (channels *ChannelsState) SetChannelFunc(channel string, fun func(chan common.Message)) {
	channels.mutex.Lock()
	defer channels.mutex.Unlock()
	channels.initFunctions[channel] = fun
}

// Add добаляет ID соединения в подписчики канала
func (channels *ChannelsState) Add(channel string, connection string) {
	channels.mutex.Lock()
	defer channels.mutex.Unlock()
	if channel == "" || connection == "" {
		return
	}

	channelEntry, ok := channels.state[channel]
	if ok {
		channelEntry.Add(connection)
	} else {
		channels.state[channel] = mapset.NewThreadUnsafeSet[string](connection)
	}
}

func (channels *ChannelsState) Remove(channel string, connection string) {
	channels.mutex.Lock()
	defer channels.mutex.Unlock()
	if channel == "" || connection == "" {
		return
	}
	channelEntry, channelExists := channels.state[channel]
	if channelExists {
		channelEntry.Remove(connection)
	}
}

func (channels *ChannelsState) Get(channel string) ([]string, bool) {
	channels.mutex.RLock()
	defer channels.mutex.RUnlock()
	channelEntry, channelExists := channels.state[channel]
	if !channelExists {
		return nil, channelExists
	}
	return channelEntry.ToSlice(), channelExists
}
