package service

type moveCommand struct {
	PlayerID     string `json:"player_id"`
	TargetTileID string `json:"target_tile_id"`
}

type playerJoinCommand struct {
	PlayerID string `json:"player_id"`
	Sprite   string `json:"sprite"`
	Name     string `json:"name"`
}

//QueueDispatcher publishes messages to AMQP
type QueueDispatcher interface {
	DispatchMessage(message interface{}) (err error)
}

type dispatcherMap map[string]QueueDispatcher
