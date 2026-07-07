package domain

type PlayerBill struct{
	PlayerID uint `json:"player_id"`
	Name string `json:"Name"`
	Amount float64 `json:"amount"`

}