package services

import(
	"testing"
	"github.com/searaaman/playledger/internal/domain"
)

func TestCalculateSessionBillsTwoPlayers(t *testing.T) {
	rahul := domain.Player{
		ID:   1,
		Name: "Rahul",
	}

	anand := domain.Player{
		ID:   2,
		Name: "Anand",
	}

	slot := domain.TimeSlot{
		Courtsbooked: 1,
		Players: []domain.Player{
			rahul,
			anand,
		},
	}

	session := domain.Session{
		CourtPrice: 300,
		TimeSlots: []domain.TimeSlot{
			slot,
		},
	}

	bills := CalculateSessionBills(session)

	
	if len(bills) != 2 {
		t.Fatalf("expected 2 bills, got %d", len(bills))
	}

	billMap := make(map[uint]float64)

	for _, bill := range bills {
		billMap[bill.PlayerID] = bill.Amount
	}

	if billMap[rahul.ID] != 150 {
		t.Errorf("expected Rahul to owe 150, got %.2f", billMap[rahul.ID])
	}

	if billMap[anand.ID] != 150 {
		t.Errorf("expected Anand to owe 150, got %.2f", billMap[anand.ID])
	}
}

func TestCalculateSessionBillsMultipleTimeSlots(t *testing.T) {
	rahul := domain.Player{
		ID:   1,
		Name: "Rahul",
	}

	anand := domain.Player{
		ID:   2,
		Name: "Anand",
	}

	slot1 := domain.TimeSlot{
		Courtsbooked: 1,
		Players: []domain.Player{
			rahul,
			anand,
		},
	}

	slot2 := domain.TimeSlot{
		Courtsbooked: 1,
		Players: []domain.Player{
			rahul,
		},
	}

	session := domain.Session{
		CourtPrice: 300,
		TimeSlots: []domain.TimeSlot{
			slot1,
			slot2,
		},
	}

	bills := CalculateSessionBills(session)

	if len(bills) != 2 {
		t.Fatalf("expected 2 bills, got %d", len(bills))
	}

	billMap := make(map[uint]float64)

	for _, bill := range bills {
		billMap[bill.PlayerID] = bill.Amount
	}

	if billMap[rahul.ID] != 450 {
		t.Errorf("expected Rahul to owe 450, got %.2f", billMap[rahul.ID])
	}

	if billMap[anand.ID] != 150 {
		t.Errorf("expected Anand to owe 150, got %.2f", billMap[anand.ID])
	}
}

func TestCalculateSessionBillsMultipleCourts(t *testing.T) {
	
	rahul := domain.Player{
		ID:   1,
		Name: "Rahul",
	}

	anand := domain.Player{
		ID:   2,
		Name: "Anand",
	}

	slot := domain.TimeSlot{
		Courtsbooked: 2,
		Players: []domain.Player{
			rahul,
			anand,
		},
	}

	session := domain.Session{
		CourtPrice: 300,
		TimeSlots: []domain.TimeSlot{
			slot,
		},
	}

	
	bills := CalculateSessionBills(session)

	
	if len(bills) != 2 {
		t.Fatalf("expected 2 bills, got %d", len(bills))
	}

	billMap := make(map[uint]float64)

	for _, bill := range bills {
		billMap[bill.PlayerID] = bill.Amount
	}

	if billMap[rahul.ID] != 300 {
		t.Errorf("expected Rahul to owe 300, got %.2f", billMap[rahul.ID])
	}

	if billMap[anand.ID] != 300 {
		t.Errorf("expected Anand to owe 300, got %.2f", billMap[anand.ID])
	}
}

func TestCalculateSessionBillsEmptyTimeSlot(t *testing.T) {
	
	slot := domain.TimeSlot{
		Courtsbooked: 1,
		Players:      []domain.Player{},
	}

	session := domain.Session{
		CourtPrice: 300,
		TimeSlots: []domain.TimeSlot{
			slot,
		},
	}

	
	bills := CalculateSessionBills(session)

	
	if len(bills) != 0 {
		t.Fatalf("expected 0 bills, got %d", len(bills))
	}
}

func TestCalculateSessionBillsEmptySession(t *testing.T) {
	
	session := domain.Session{
		CourtPrice: 300,
		TimeSlots:  []domain.TimeSlot{},
	}

	
	bills := CalculateSessionBills(session)

	
	if len(bills) != 0 {
		t.Fatalf("expected 0 bills, got %d", len(bills))
	}
}