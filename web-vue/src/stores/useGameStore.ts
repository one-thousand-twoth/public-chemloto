import { WEBSOCKET_EVENT } from '@/api/websocket/websocket'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'




export const useGameStore = defineStore('game', () => {

    const fetching = ref(false)
    const connected = ref(false)
    const name = ref("")
    const gameState = ref<GameInfo>({
        Bag: {
            Elements: undefined,
            LastElements: []
        },
        Players: [],
        Started: false,
        State: "none"
    })
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
        LastElements,
        gameState,
        EngineInfo
    }

})

export interface Bag {
    Elements: any,
    LastElements: Array<string>
}
export interface Player{
    Name: string,
    Role: number,
    RaisedHand: boolean,

}
export interface GameInfo {
    Bag: Bag,
    Players: Array<Player>,
    Started: boolean,
    State: String
} 