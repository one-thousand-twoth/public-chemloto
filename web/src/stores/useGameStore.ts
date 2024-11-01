import { WEBSOCKET_EVENT } from '@/api/websocket/websocket'
import { Role } from '@/models/User'
import { acceptHMRUpdate, defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useUserStore } from './useUserStore'


const ROOMNAME_LOCAL_STORAGE_KEY = "roomname"
const CONNECTED_LOCAL_STORAGE_KEY = "connected"
const getName = () => {
    const value = localStorage.getItem(ROOMNAME_LOCAL_STORAGE_KEY)
    return value ?? "";
}
const getConnected = () => {
    const value = localStorage.getItem(CONNECTED_LOCAL_STORAGE_KEY)
    return value === "true";
}

export const useGameStore = defineStore('game', () => {
    const userStore = useUserStore()
    const fetching = ref(false)
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
    const timer = ref(0)
    const gameState = ref<GameInfo>({
        Bag: {
            Elements: {},
            LastElements: []
        },
        Players: [],
        Started: false,
        State: undefined,
        RaisedHands: [],
        Fields: {

        }
    })

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
        const b = e.Body["engine"] as GameInfo
        gameState.value = b
        if (gameState.value.State === "OBTAIN")
            if (gameState.value.StateStruct) {
                timer.value = gameState.value.StateStruct.Timer
            }

    }
    return {
        fetching,
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
    Player: Player,
    Field: string,
    Name: string,
    Structure: { [id: string]: number; }
}
export type GameInfo = {
    Bag: Bag,
    Players: Array<Player>,
    RaisedHands: Array<Hand>
    Started: boolean,
    // State
    Fields: { [id: string]: { Score: number }; }
} & State
export type State = StateOBTAIN | StateTRADE | StateCOMPLETED | StateHAND | StateUndefined
export interface StateOBTAIN {
    State: "OBTAIN"
    StateStruct?: { Timer: number }
}
export interface StateTRADE {
    State: "TRADE",
    StateStruct?: {
        StockExchange: {
            StockList: {
                Owner: string, Element: string, ToElement: string,
                Request: {
                    ID: string
                    Player: string
                    Accept: boolean
                }[]
            }[]
        }
    }
}
export interface StateCOMPLETED {
    State: "COMPLETED"
}
export interface StateHAND {
    State: "HAND"
}
export interface StateUndefined {
    State: undefined
}

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useGameStore, import.meta.hot))
}
