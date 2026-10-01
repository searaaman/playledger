package services

import (
	"github.com/searaaman/playledger/internal/domain"
	"testing"
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
		CourtsBooked: 1,
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
		CourtsBooked: 1,
		Players: []domain.Player{
			rahul,
			anand,
		},
	}

	slot2 := domain.TimeSlot{
		CourtsBooked: 1,
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
		CourtsBooked: 2,
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
		CourtsBooked: 1,
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

// The session we verified by hand through the API.
func TestCalculateSessionBillsRealSession(t *testing.T) {
	rahul := domain.Player{ID: 1, Name: "Rahul"}
	akhil := domain.Player{ID: 2, Name: "Akhil"}
	anand := domain.Player{ID: 3, Name: "Anand"}

	session := domain.Session{
		CourtPrice: 300,
		TimeSlots: []domain.TimeSlot{
			{CourtsBooked: 2, Players: []domain.Player{rahul, akhil, anand}},
			{CourtsBooked: 1, Players: []domain.Player{rahul, anand}},
		},
	}

	bills := CalculateSessionBills(session)

	expected := []domain.PlayerBill{
		{PlayerID: 2, Name: "Akhil", Amount: 200},
		{PlayerID: 3, Name: "Anand", Amount: 350},
		{PlayerID: 1, Name: "Rahul", Amount: 350},
	}
	if len(bills) != len(expected) {
		t.Fatalf("expected %d bills, got %d", len(expected), len(bills))
	}
	// Bills come back sorted by name.
	for i := range expected {
		if bills[i] != expected[i] {
			t.Errorf("bill %d: expected %+v, got %+v", i, expected[i], bills[i])
		}
	}
}
