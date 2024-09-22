
export interface WEBSOCKET_EVENT {
    Type: string
    Ok: boolean
    Errors: Array<string>
    Body: { [id: string]: any; }
}
