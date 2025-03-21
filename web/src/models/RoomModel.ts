export type Room = {
  name: string,
  type: string,
  engine: {Status: EngineStatus},
}

export type EngineStatus = "STATUS_WAITING" | "STATUS_STARTED" | "STATUS_COMPLETED"

