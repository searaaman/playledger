import api from "./client";
import type {
  LedgerEntry,
  Payment,
  Player,
  PlayerBill,
  PlayerLedger,
  Session,
  TimeSlot,
  User,
} from "../types";

// Auth

export async function login(email: string, password: string) {
  const { data } = await api.post<{ token: string; user: User }>("/login", { email, password });
  return data;
}

export async function register(name: string, email: string, password: string) {
  const { data } = await api.post<{ user: User }>("/register", { name, email, password });
  return data;
}

// Sessions

export interface NewSession {
  start_time: string;
  end_time: string;
  court_price: number;
  hourly_slots: boolean;
  default_courts: number;
}

export async function listSessions() {
  return (await api.get<Session[]>("/sessions")).data;
}

export async function getSession(id: number) {
  return (await api.get<Session>(`/sessions/${id}`)).data;
}

export async function createSession(body: NewSession) {
  return (await api.post<Session>("/sessions", body)).data;
}

export async function deleteSession(id: number) {
  await api.delete(`/sessions/${id}`);
}

export async function getSessionBilling(id: number) {
  return (await api.get<PlayerBill[]>(`/sessions/${id}/billing`)).data;
}

// Time slots

export async function createTimeSlot(
  sessionId: number,
  body: { start_time: string; end_time: string; courts_booked: number },
) {
  return (await api.post<TimeSlot>(`/sessions/${sessionId}/timeslots`, body)).data;
}

export async function updateTimeSlotCourts(slotId: number, courts: number) {
  return (await api.put<TimeSlot>(`/timeslots/${slotId}`, { courts_booked: courts })).data;
}

export async function deleteTimeSlot(slotId: number) {
  await api.delete(`/timeslots/${slotId}`);
}

export async function assignPlayer(slotId: number, playerId: number) {
  await api.post(`/timeslots/${slotId}/players`, { player_id: playerId });
}

export async function removePlayer(slotId: number, playerId: number) {
  await api.delete(`/timeslots/${slotId}/players/${playerId}`);
}

// Players and ledger

export async function listPlayers() {
  return (await api.get<Player[]>("/players")).data;
}

export async function createPlayer(name: string, phone: string) {
  return (await api.post<Player>("/players", { name, phone })).data;
}

export async function getLedger() {
  return (await api.get<LedgerEntry[]>("/ledger")).data;
}

export async function getPlayerLedger(id: number) {
  return (await api.get<PlayerLedger>(`/players/${id}/ledger`)).data;
}

// Payments

export async function listPayments(filter: { player_id?: number; session_id?: number } = {}) {
  return (await api.get<Payment[]>("/payments", { params: filter })).data;
}

export async function createPayment(playerId: number, sessionId: number, amount: number) {
  return (
    await api.post<Payment>("/payments", {
      player_id: playerId,
      session_id: sessionId,
      amount,
    })
  ).data;
}

export async function deletePayment(id: number) {
  await api.delete(`/payments/${id}`);
}
