package hub

import (
	"encoding/json"
	"sync"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	mapset "github.com/deckarep/golang-set/v2"
)

// channelsState содержит карту со списками ID соединений подписанных на канал
type channelsState struct {
	// map key is a channel name
	state map[string]mapset.Set[string]
	// Функции каналов срабатывающие при первом подключении.
	// key: channel name
	initFunctions map[string]func(chan common.Message)
	mutex         sync.RWMutex
}

func (rs *channelsState) MarshalJSON() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.state)
}

func (channels *channelsState) GetChannelFunc(channel string) (func(chan common.Message), bool) {
	channels.mutex.Lock()
	defer channels.mutex.Unlock()
	fun, ok := channels.initFunctions[channel]
	return fun, ok
}
func (channels *channelsState) SetChannelFunc(channel string, fun func(chan common.Message)) {
	channels.mutex.Lock()
	defer channels.mutex.Unlock()
	channels.initFunctions[channel] = fun
}

// Add добаляет ID соединения в подписчики канала
func (channels *channelsState) Add(channel string, connection string) {
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

func (channels *channelsState) Remove(channel string, connection string) {
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

func (channels *channelsState) Get(channel string) ([]string, bool) {
	channels.mutex.RLock()
	defer channels.mutex.RUnlock()
	channelEntry, channelExists := channels.state[channel]
	if !channelExists {
		return nil, channelExists
	}
	return channelEntry.ToSlice(), channelExists
}
