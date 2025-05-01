import { WebsocketConnector } from '@/api/websocket/websocket'
import { BaseStateHandler as BaseStateController } from './base'

// Trade state handler
export default class TradeStateController extends BaseStateController<'TRADE'> {
  constructor (protected ws: WebsocketConnector) {
    super(ws, 'TRADE')
  }

  getState () {
    return 'TRADE'
  }

  trade (element: string, toElement: string) {
    this.sendCommand({
      Type: 'ENGINE_ACTION',
      Action: 'TradeOffer',
      Element: element,
      ToElement: toElement
    })
  }

  cancelTrade () {
    this.sendCommand({
      Type: 'ENGINE_ACTION',
      Action: 'RemoveTradeOffer'
      // StockId: stockId
    })
  }

  requestTrade (stockId: string, accept: boolean) {
    this.sendCommand({
      Type: 'ENGINE_ACTION',
      Action: 'TradeRequest',
      StockId: stockId,
      Accept: accept
    })
  }
  ackTrade (requestID: string) {
    this.sendCommand({
      Type: 'ENGINE_ACTION',
      Action: 'TradeAck',
      TargetID: requestID
    })
  }
  sendContinue () {
    this.sendCommand({
      Type: 'ENGINE_ACTION',
      Action: 'Continue'
    })
  }
}
