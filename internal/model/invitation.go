package model

type Invitation struct {
	UserId      string `json:"user_id"`
	ContainerId string `json:"container_id"`
	Recipient   string `json:"recipient"`
}
