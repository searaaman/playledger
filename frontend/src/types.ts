// Shapes returned by the Go API. Keep in sync with internal/domain.

export interface Player {
  id: number;
  name: string;
  phone: string;
}

export interface TimeSlot {
  id: number;
  session_id: number;
  start_time: string;
  end_time: string;
  courts_booked: number;
  players: Player[] | null;
}

export interface Session {
  id: number;
  start_time: string;
  end_time: string;
  court_price: number;
  time_slots: TimeSlot[] | null;
}

export interface PlayerBill {
  player_id: number;
  name: string;
  amount: number;
}

export interface LedgerEntry {
  player_id: number;
  name: string;
  phone: string;
  total_bill: number;
  total_paid: number;
  // total_paid - total_bill: negative means the player owes money.
  balance: number;
}

export interface SessionCharge {
  session_id: number;
  start_time: string;
  end_time: string;
  slot_count: number;
  bill: number;
  paid: number;
}

export interface PlayerLedger extends LedgerEntry {
  sessions: SessionCharge[];
  payments: Payment[];
}

export interface Payment {
  id: number;
  player_id: number;
  session_id: number;
  amount: number;
  created_at: string;
  player?: Player;
}

export interface User {
  id: number;
  name: string;
  email: string;
}
