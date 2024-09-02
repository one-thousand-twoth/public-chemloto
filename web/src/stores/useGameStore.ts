import { WEBSOCKET_EVENT } from '@/api/websocket/websocket'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { Role, useUserStore } from './useUserStore'

const ROOMNAME_LOCAL_STORAGE_KEY = "roomname"
const getName = () => {
    const value = localStorage.getItem(ROOMNAME_LOCAL_STORAGE_KEY)
    return value ? value : "";
}

export const useGameStore = defineStore('game', () => {

    const fetching = ref(false)
    const connected = ref(false)
    const roomname = ref(getName())
    const name  = computed({
        get: () => {
           return roomname.value
        },
        set: (v) => {
           roomname.value = v
           if (v === ""){
            roomname.value = ""
           }
           localStorage.setItem(ROOMNAME_LOCAL_STORAGE_KEY, roomname.value);
        }
      })
    const timer = ref(0)
    const gameState = ref<GameInfo>({
        Bag: {
            Elements: {},
            LastElements: []
        },
        Players: [],
        Started: false,
        State: "none",
        RaisedHands: [],
        Fields: {

        }
    })
    const userStore = useUserStore()
    // const LastElements = ref<Array<string>>(Array.from({ length: 5 }, () => "UNDEFINED"))
    // const currElement = ref('UNDEFINED')
    const LastElements = computed(() => {
        const elems: Array<string> = Object.assign([], gameState.value.Bag.LastElements);
        return elems.
            reverse().
            concat(Array.from({ length: 5 - elems.length }, () => "UNDEFINED"));
    })

    const currElement = computed(() => {
        const elems: Array<string> = gameState.value.Bag.LastElements;
        return elems[elems.length - 1] ? gameState.value.Bag.LastElements[elems.length - 1] : "UNDEFINED"
    })

    const SelfPlayer = computed(() => {

        return gameState.value.Players.find((player) => player.Name == userStore.UserCreds?.username) as Player
    })

    function EngineInfo(e: WEBSOCKET_EVENT) {
        // console.log("changing to room")
        const b = e.Body["engine"] as GameInfo
        gameState.value = b
        if (gameState.value.StateStruct) {
            timer.value = gameState.value.StateStruct.Timer
        }

    }
    // function Timer(e: WEBSOCKET_EVENT) {
    //     const b = e.Body["Value"] as number
    //     timer.value = b

    // }
    return {
        fetching,
        connected,
        name,
        timer,
        currElement,
        SelfPlayer,
        LastElements,
        gameState,
        EngineInfo
    }

})

export interface Bag {
    Elements: { [id: string]: number; },
    LastElements: Array<string>
}
export interface Player {
    Name: string,
    Role: Role,
    Score: number,
    RaisedHand: boolean,
    Bag: { [id: string]: number; }
}
export interface Hand {
    Player: Player
    Field: string,
    Name: string,
    Structure: { [id: string]: number; }
}
export interface GameInfo {
    Bag: Bag,
    Players: Array<Player>,
    RaisedHands: Array<Hand>
    Started: boolean,
    State: String,
    StateStruct?: { Timer: number }
    Fields: { [id: string]: { Score: number }; }
} 