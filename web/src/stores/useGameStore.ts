import { WEBSOCKET_EVENT, WebsocketConnector } from '@/api/websocket/websocket'
import { GameInfo, getStateTimer, hasTimer, Player, updateTimer } from '@/models/Game'
import { acceptHMRUpdate, defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'
import { useUserStore } from './useUserStore'


const ROOMNAME_LOCAL_STORAGE_KEY = "roomname"
const getName = () => {
    const value = localStorage.getItem(ROOMNAME_LOCAL_STORAGE_KEY)
    return value ?? "";
}

let lastTimerID: NodeJS.Timeout

export function StartTimer() {
    const gameStore = useGameStore()
    const state = gameStore.gameState
    console.log("Hei")
    clearInterval(lastTimerID)
    if (hasTimer(state) && state.StateStruct.Timer > 1) {
        console.log("Cool")
        lastTimerID = setInterval(() => {
            if (state.StateStruct.Timer == null || state.StateStruct.TimerStatus == "Stopped") {
                return
            }
            state.StateStruct.Timer--; // Decrement the timer count
        }, 1000);

    }
    console.log("Hi")
}


export const useGameStore = defineStore('game', () => {
    const userStore = useUserStore()
    const fetching = ref(true)
    const roomname = ref(getName())
    const name = computed({
        get: () => {
            return roomname.value
        },
        set: (v) => {
            roomname.value = v
            if (v === "") {
                roomname.value = ""
            }
            localStorage.setItem(ROOMNAME_LOCAL_STORAGE_KEY, roomname.value);
        }
    })
    // TODO: Сделать getter`om
    const timer = computed({
        get: () => getStateTimer(gameState.value),
        set: (v) => { updateTimer(gameState.value, v as number) }
    })
    const gameState = ref<GameInfo>({
        Bag: {
            Elements: {},
            LastElements: [],
            RemainingElements: {},
            DraftedElements: {}
        },
        Players: [],
        Status: 'STATUS_WAITING',
        State: 'OBTAIN',
        StateStruct:{
            Timer: 0,
            TimerStatus: ''
        },
        RaisedHands: [],
        Fields: {

        }
    })

    const LastElements = computed(() => {
        const elems: Array<string> = Object.assign([], gameState.value.Bag.LastElements);
        return elems.
            reverse().
            concat(Array.from({ length: 6 - elems.length }, () => "UNDEFINED"));
    })

    const currElement = computed(() => {
        const elems: Array<string> = gameState.value.Bag.LastElements;
        return elems[elems.length - 1] ? gameState.value.Bag.LastElements[elems.length - 1] : "UNDEFINED"
    })

    const SelfPlayer = computed(() => {
        return gameState.value.Players.find((player) => player.Name == userStore.UserCreds?.username)
    })

    function EngineInfo(e: WEBSOCKET_EVENT) {

        const b = e.Body["engine"] as GameInfo
        gameState.value = b
        fetching.value = false

        StartTimer()
    }


    return {
        fetching,
        name,
        timer,
        currElement,
        SelfPlayer,
        LastElements,
        gameState,
        EngineInfo,
        // Trade,
    }

})


if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useGameStore, import.meta.hot))
}
