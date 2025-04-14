import { WEBSOCKET_EVENT, WebsocketConnector } from '@/api/websocket/websocket'
import { GameInfo } from '@/models/Game'
import { Role } from '@/models/User'
import { acceptHMRUpdate, defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'
import { useUserStore } from './useUserStore'


const ROOMNAME_LOCAL_STORAGE_KEY = "roomname"
// const CONNECTED_LOCAL_STORAGE_KEY = "connected"
const getName = () => {
    const value = localStorage.getItem(ROOMNAME_LOCAL_STORAGE_KEY)
    return value ?? "";
}
// const getConnected = () => {
//     const value = localStorage.getItem(CONNECTED_LOCAL_STORAGE_KEY)
//     return value === "true";
// }
// Base interface for state handlers
interface GameStateHandler {
    getState(): string;
}

// Trade state handler
export class TradeStateHandler implements GameStateHandler {
    constructor(private ws: any) { }

    getState() {
        return 'TRADE';
    }

    trade(element: string, toElement: string) {
        this.ws.Send({
            Type: "ENGINE_ACTION",
            Action: "TradeOffer",
            Element: element,
            ToElement: toElement
        });
    }

    cancelTrade() {
        this.ws.Send({
            Type: "ENGINE_ACTION",
            Action: "RemoveTradeOffer",
            // StockId: stockId
        });
    }

    requestTrade(stockId: string, accept: boolean) {
        this.ws.Send({
            Type: "ENGINE_ACTION",
            Action: "TradeRequest",
            StockId: stockId,
            Accept: accept
        });
    }
    ackTrade(requestID: string) {
        this.ws.Send({
            Type: 'ENGINE_ACTION',
            Action: "TradeAck",
            TargetID: requestID,
        })
    }
    sendContinue() {
        this.ws.Send({
            Type: 'ENGINE_ACTION',
            Action: 'Continue'
        })
    }
}

// Obtain state handler
export class ObtainStateHandler implements GameStateHandler {
    constructor(private ws: WebsocketConnector) { }
    getState() {
        return 'OBTAIN';
    }
    getElement() {
        this.ws.Send({
            Type: 'ENGINE_ACTION',
            Action: 'GetElement'
        });
    }
    sendContinue() {
        this.ws.Send({
            Type: 'ENGINE_ACTION',
            Action: 'Continue'
        })
    }


}
// Factory to create state handlers
class GameStateFactory {
    constructor(private readonly ws: WebsocketConnector) { }
    public createHandler(state: string): GameStateHandler {
        switch (state) {
            case 'TRADE':
                return new TradeStateHandler(this.ws);
            case 'OBTAIN':
                return new ObtainStateHandler(this.ws);
            default:
                throw new Error(`Unknown state: ${state}`);
        }
    }
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
        State: undefined,
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
        return gameState.value.Players.find((player) => player.Name == userStore.UserCreds?.username) as Player
    })

    function EngineInfo(e: WEBSOCKET_EVENT) {
        const b = e.Body["engine"] as GameInfo
        gameState.value = b

        StartTimer()
    }

    const handlerFactory = new GameStateFactory(inject('connector')!);

    const currentStateHandler = computed(() => {
        if (gameState.value.State === undefined) {
            return null
        }
        return handlerFactory.createHandler(gameState.value.State)
    })

    return {
        fetching,
        name,
        timer,
        currElement,
        SelfPlayer,
        LastElements,
        gameState,
        EngineInfo,
        currentStateHandler,
        // Trade,
    }

})


if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useGameStore, import.meta.hot))
}
