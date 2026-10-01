package services

import (
	"sort"

	"github.com/searaaman/playledger/internal/domain"
)

// BuildLedger totals every player's bills across all sessions and their payments.
// Because it works over all history, any unpaid amount naturally carries forward.
func BuildLedger(players []domain.Player, sessions []domain.Session, payments []domain.Payment) []domain.LedgerEntry {
	entries := make(map[uint]*domain.LedgerEntry, len(players))
	for _, player := range players {
		entries[player.ID] = &domain.LedgerEntry{
			PlayerID: player.ID,
			Name:     player.Name,
			Phone:    player.Phone,
		}
	}

	for _, session := range sessions {
		for _, bill := range CalculateSessionBills(session) {
			if entry, ok := entries[bill.PlayerID]; ok {
				entry.TotalBill += bill.Amount
			}
		}
	}
	for _, payment := range payments {
		if entry, ok := entries[payment.PlayerID]; ok {
			entry.TotalPaid += payment.Amount
		}
	}

	ledger := make([]domain.LedgerEntry, 0, len(entries))
	for _, entry := range entries {
		entry.Balance = entry.TotalPaid - entry.TotalBill
		ledger = append(ledger, *entry)
	}
	// Biggest debts first, then alphabetical.
	sort.Slice(ledger, func(i, j int) bool {
		if ledger[i].Balance != ledger[j].Balance {
			return ledger[i].Balance < ledger[j].Balance
		}
		return ledger[i].Name < ledger[j].Name
	})
	return ledger
}

// BuildPlayerLedger is the ledger for a single player, broken down by session.
// Sessions appear if the player played in them or paid towards them.
func BuildPlayerLedger(player domain.Player, sessions []domain.Session, payments []domain.Payment) domain.PlayerLedger {
	paidBySession := make(map[uint]float64)
	var totalPaid float64
	for _, payment := range payments {
		if payment.PlayerID != player.ID {
			continue
		}
		paidBySession[payment.SessionID] += payment.Amount
		totalPaid += payment.Amount
	}

	charges := []domain.SessionCharge{}
	var totalBill float64
	for _, session := range sessions {
		charge := domain.SessionCharge{
			SessionID: session.ID,
			StartTime: session.StartTime,
			EndTime:   session.EndTime,
			Paid:      paidBySession[session.ID],
		}
		for _, bill := range CalculateSessionBills(session) {
			if bill.PlayerID == player.ID {
				charge.Bill = bill.Amount
			}
		}
		for _, slot := range session.TimeSlots {
			for _, p := range slot.Players {
				if p.ID == player.ID {
					charge.SlotCount++
				}
			}
		}
		if charge.SlotCount == 0 && charge.Paid == 0 {
			continue
		}
		totalBill += charge.Bill
		charges = append(charges, charge)
	}
	// Most recent session first.
	sort.Slice(charges, func(i, j int) bool {
		return charges[i].StartTime.After(charges[j].StartTime)
	})

	playerPayments := []domain.Payment{}
	for _, payment := range payments {
		if payment.PlayerID == player.ID {
			playerPayments = append(playerPayments, payment)
		}
	}

	return domain.PlayerLedger{
		LedgerEntry: domain.LedgerEntry{
			PlayerID:  player.ID,
			Name:      player.Name,
			Phone:     player.Phone,
			TotalBill: totalBill,
			TotalPaid: totalPaid,
			Balance:   totalPaid - totalBill,
		},
		Sessions: charges,
		Payments: playerPayments,
	}
}
