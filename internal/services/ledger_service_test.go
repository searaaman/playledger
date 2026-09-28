package services

import (
	"testing"
	"time"

	"github.com/searaaman/playledger/internal/domain"
)

var (
	rahul = domain.Player{ID: 1, Name: "Rahul"}
	akhil = domain.Player{ID: 2, Name: "Akhil"}
	anand = domain.Player{ID: 3, Name: "Anand"}
	kiran = domain.Player{ID: 4, Name: "Kiran"}
)

func ledgerSessions() []domain.Session {
	day1 := time.Date(2026, 9, 1, 19, 0, 0, 0, time.UTC)
	day2 := day1.AddDate(0, 0, 7)
	return []domain.Session{
		{
			ID: 1, StartTime: day1, EndTime: day1.Add(2 * time.Hour), CourtPrice: 300,
			TimeSlots: []domain.TimeSlot{
				{CourtsBooked: 2, Players: []domain.Player{rahul, akhil, anand}},
				{CourtsBooked: 1, Players: []domain.Player{rahul, anand}},
			},
		},
		{
			ID: 2, StartTime: day2, EndTime: day2.Add(time.Hour), CourtPrice: 250,
			TimeSlots: []domain.TimeSlot{
				{CourtsBooked: 1, Players: []domain.Player{rahul}},
			},
		},
	}
}

func TestBuildLedgerCarriesForwardUnpaidBalance(t *testing.T) {
	payments := []domain.Payment{
		{PlayerID: 1, SessionID: 1, Amount: 300},
		{PlayerID: 2, SessionID: 1, Amount: 200},
	}

	ledger := BuildLedger([]domain.Player{rahul, akhil, anand, kiran}, ledgerSessions(), payments)

	byID := make(map[uint]domain.LedgerEntry)
	for _, entry := range ledger {
		byID[entry.PlayerID] = entry
	}

	// Rahul: 350 + 250 billed, 300 paid, so the 50 from session 1 carries into session 2.
	if got := byID[1]; got.TotalBill != 600 || got.TotalPaid != 300 || got.Balance != -300 {
		t.Errorf("Rahul: unexpected ledger %+v", got)
	}
	if got := byID[2]; got.Balance != 0 {
		t.Errorf("Akhil: expected settled, got %+v", got)
	}
	if got := byID[3]; got.Balance != -350 {
		t.Errorf("Anand: expected -350, got %+v", got)
	}
	// Players who never played still appear, with a zero balance.
	if got, ok := byID[4]; !ok || got.Balance != 0 {
		t.Errorf("Kiran: expected zero entry, got %+v (present=%v)", got, ok)
	}

	// Biggest debt first.
	if ledger[0].PlayerID != 3 {
		t.Errorf("expected Anand (owes 350) first, got %s", ledger[0].Name)
	}
}

func TestBuildPlayerLedgerBreaksDownBySession(t *testing.T) {
	payments := []domain.Payment{
		{PlayerID: 1, SessionID: 1, Amount: 300},
		{PlayerID: 2, SessionID: 1, Amount: 200},
	}

	ledger := BuildPlayerLedger(rahul, ledgerSessions(), payments)

	if ledger.TotalBill != 600 || ledger.TotalPaid != 300 || ledger.Balance != -300 {
		t.Errorf("unexpected totals %+v", ledger.LedgerEntry)
	}
	if len(ledger.Sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(ledger.Sessions))
	}
	// Newest session first.
	latest, earliest := ledger.Sessions[0], ledger.Sessions[1]
	if latest.SessionID != 2 || latest.Bill != 250 || latest.Paid != 0 || latest.SlotCount != 1 {
		t.Errorf("unexpected session 2 charge %+v", latest)
	}
	if earliest.SessionID != 1 || earliest.Bill != 350 || earliest.Paid != 300 || earliest.SlotCount != 2 {
		t.Errorf("unexpected session 1 charge %+v", earliest)
	}
	if len(ledger.Payments) != 1 {
		t.Errorf("expected only Rahul's payment, got %d", len(ledger.Payments))
	}
}

func TestBuildPlayerLedgerSkipsSessionsPlayerWasNotIn(t *testing.T) {
	ledger := BuildPlayerLedger(akhil, ledgerSessions(), nil)

	if len(ledger.Sessions) != 1 || ledger.Sessions[0].SessionID != 1 {
		t.Errorf("expected only session 1, got %+v", ledger.Sessions)
	}
	if ledger.Balance != -200 {
		t.Errorf("expected -200, got %.2f", ledger.Balance)
	}
}
