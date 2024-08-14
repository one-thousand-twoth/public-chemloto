import { EngineAction, Subscribe } from "./handlers";

export interface WEBSOCKET_EVENT {
    Type:   string
	Ok :    boolean
	Errors: Array<string>
	Body  : { [id: string]: any; }
}

export class WebsocketConnector {
    baseUrl: string;
    token: string;
    active: boolean;
    connection!: WebSocket | null;
    // gameStore: Store;
    // gameStore: 

    constructor(baseUrl: string, token: string = "") {
        this.token = token;
        this.baseUrl = baseUrl;
        this.active = false;
    }
    Run() {
        // const gameStore = useGameStore()
        this.connection = new WebSocket(`ws://${this.baseUrl}/api/v1/ws?token=${this.token}`)
        this.connection.onmessage = function (event) {
            const data = JSON.parse( event.data ) as WEBSOCKET_EVENT
            switch (data.Type){
                case "HUB_SUBSCRIBE": 
                    Subscribe(data)
                    break;
                
                case "ENGINE_ACTION":
                    EngineAction(data)
                    break;
                
            }
            console.log(data)
        }

        this.connection.onopen = function (event) {
            console.log(event)
            console.log("Successfully connected to the echo websocket server...")
        }
    }
    Send(msg: Object){
        if (!this.connection){
            return
        }
        this.connection.send(JSON.stringify(msg))
    }
}