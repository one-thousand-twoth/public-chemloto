package hub

import (
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
)

type channelsState struct {
	state map[string]mapset.Set[string]
	mutex sync.RWMutex
}

func (channels *channelsState) Add(channel string, connection string) {
	channels.mutex.Lock()
	defer channels.mutex.Unlock()

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

	channelEntry, channelExists := channels.state[channel]
	if channelExists {
		channelEntry.Remove(connection)
	}
}

func (channels *channelsState) Get(channel string) ([]string, bool) {
	channels.mutex.RLock()
	defer channels.mutex.RUnlock()

	channelEntry, channelExists := channels.state[channel]
	return channelEntry.ToSlice(), channelExists
}
