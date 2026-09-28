package services

import (
	"sort"

	"github.com/searaaman/playledger/internal/domain"
)

// CalculateSessionBills splits each slot's court cost evenly between the players
// in that slot and totals the result per player. Empty slots are not billed.
func CalculateSessionBills(session domain.Session) []domain.PlayerBill {
	playerBills := make(map[uint]*domain.PlayerBill)
	for _, slot := range session.TimeSlots {
		if len(slot.Players) == 0 {
			continue
		}
		slotCost := float64(slot.CourtsBooked) * session.CourtPrice
		costPerPlayer := slotCost / float64(len(slot.Players))
		for _, player := range slot.Players {
			bill, exists := playerBills[player.ID]
			if exists {
				bill.Amount += costPerPlayer
			} else {
				playerBills[player.ID] = &domain.PlayerBill{
					PlayerID: player.ID,
					Name:     player.Name,
					Amount:   costPerPlayer,
				}
			}
		}
	}

	bills := make([]domain.PlayerBill, 0, len(playerBills))
	for _, bill := range playerBills {
		bills = append(bills, *bill)
	}
	// Map iteration order is random; sort so API responses are stable.
	sort.Slice(bills, func(i, j int) bool {
		if bills[i].Name != bills[j].Name {
			return bills[i].Name < bills[j].Name
		}
		return bills[i].PlayerID < bills[j].PlayerID
	})
	return bills
}
