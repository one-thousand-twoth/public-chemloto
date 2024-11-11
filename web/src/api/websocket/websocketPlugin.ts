// websocketPlugin.ts
import { WebsocketConnector } from '@/api/websocket/websocket'
import { PiniaPluginContext } from 'pinia'
import { Plugin } from 'vue'

export interface StoreWithWS {
    $ws: WebsocketConnector
}

export const websocketPlugin: Plugin = {
    install: (app, wsConnection: WebsocketConnector) => {
        app.provide('connector', wsConnection)
    }
}

export const piniaWebsocketPlugin = (wsConnection: WebsocketConnector) => {
    return ({ store }: PiniaPluginContext) => {
        store.$ws = wsConnection
    }
}
