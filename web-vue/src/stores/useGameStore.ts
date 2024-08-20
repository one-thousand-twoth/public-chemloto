import { WEBSOCKET_EVENT } from '@/api/websocket/websocket'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { Role, useUserStore } from './useUserStore'




export const useGameStore = defineStore('game', () => {

    const fetching = ref(false)
    const connected = ref(false)
    const name = ref("")
    const gameState = ref<GameInfo>({
        Bag: {
            Elements: {},
            LastElements: []
        },
        Players: [],
        Started: false,
        State: "none",
        RaisedHands: []
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

        return  gameState.value.Players.find((player) =>  player.Name == userStore.UserCreds?.username) as Player
    })

    function EngineInfo(e: WEBSOCKET_EVENT) {
        // console.log("changing to room")
        const b = e.Body["engine"] as GameInfo
        gameState.value = b

    }
    return {
        fetching,
        connected,
        name,
        currElement,
        SelfPlayer,
        LastElements,
        gameState,
        EngineInfo
    }

})

export interface Bag {
    Elements:  {[id: string] : number; },
    LastElements: Array<string>
}
export interface Player{
    Name: string,
    Role: Role,
    Score: number,
    RaisedHand: boolean,
    Bag:  {[id: string] : number; }
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
    State: String
} 