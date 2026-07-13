package services

import(
	"github.com/searaaman/playledger/internal/domain"
)




func CalculateSessionBills(session domain.Session) []domain.PlayerBill{

	
	playerBills := make(map[uint]*domain.PlayerBill)
	for _,slot:=range session.TimeSlots{
		if len(slot.Players)==0{
		continue
	}
		slotcost:=float64(slot.Courtsbooked)*session.CourtPrice
		costPerPlayer:=slotcost/float64(len(slot.Players))
		for _,player:=range slot.Players{
			bill,exists:=playerBills[player.ID]

			if exists{
				bill.Amount+=costPerPlayer
			}else{
				playerBills[player.ID]=&domain.PlayerBill{
					PlayerID:player.ID,
					Name:player.Name,
					Amount:costPerPlayer,
				}
			}
		}
	}
	var bills []domain.PlayerBill
	for _,bill:=range playerBills{
		bills=append(bills,*bill)
	}
	return bills
}
