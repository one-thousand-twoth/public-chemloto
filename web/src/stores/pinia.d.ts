import 'pinia';
import { WebsocketConnector } from '@/api/websocket/websocket';

declare module 'pinia' {
    export interface PiniaCustomProperties {
        $ws: WebsocketConnector
    }
}