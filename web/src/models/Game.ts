import { Role } from './User'

// Type definition for any object that might have a timer
type WithTimer = {
  StateStruct: {
    Timer: number
    TimerStatus: string
  }
}

// Type guard to check if state has timer structure
export const hasTimer = (state: States): state is States & WithTimer => {
  return typeof (state as WithTimer).StateStruct?.Timer === 'number'
}

// Simple timer utility
export const getStateTimer = (state: States): number | null => {
  return hasTimer(state) ? state.StateStruct.Timer : null
}

// Timer update utility with type narrowing
export const updateTimer = (
  state: States,
  newTimer: number
): state is States & WithTimer => {
  if (!hasTimer(state)) {
    return false
  }
  state.StateStruct.Timer = newTimer
  return true
}

export interface Bag {
  Elements: { [id: string]: number }
  LastElements: Array<string>
  RemainingElements: { [id: string]: number }
  DraftedElements: { [id: string]: number }
}
export interface Player {
  Name: string
  Role: Role
  Score: number
  RaisedHand: boolean
  Bag: { [id: string]: number }
  CompletedFields: Array<string>
}
export interface Hand {
  Player: Player
  Field: string
  Name: string
  Structure: { [id: string]: number }
  Checked: boolean
}
export type GameInfo = {
  Bag: Bag
  Players: Array<Player>
  RaisedHands: Array<Hand>
  Status: 'STATUS_WAITING' | 'STATUS_STARTED' | 'STATUS_COMPLETED'
  // State
  Fields: { [id: string]: { Score: number } }
} & States
export type States = StateOBTAIN | StateTRADE | StateCOMPLETED | StateHAND
export type StateNames = States['State']
export interface StateOBTAIN {
  State: 'OBTAIN'
  StateStruct: { Timer: number; TimerStatus: string }
}
export interface RequestEntity {
  ID: string
  Player: string
  Accept: boolean
}

export interface StockEntity {
  ID: string
  Owner: string
  Element: string
  ToElement: string
  Requests: { [id: string]: RequestEntity }
}

export interface StateTRADE {
  State: 'TRADE'
  StateStruct: {
    Timer: number
    TimerStatus: string
    StockExchange: {
      StockList: StockEntity[]
      TradeLog: TradeLog[]
    }
  }
  // Trade(): void,
}
export interface TradeLog {
  User: string
  GetElement: string
  GaveElement: string
}
export interface StateCOMPLETED {
  State: 'COMPLETED'
}
export interface StateHAND {
  State: 'HAND'
}
