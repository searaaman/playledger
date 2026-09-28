import type { Player, Session, TimeSlot } from "../types";

export function slotsOf(session: Session): TimeSlot[] {
  return session.time_slots ?? [];
}

export function playersOf(slot: TimeSlot): Player[] {
  return slot.players ?? [];
}

// Everyone who played at least one slot, sorted by name.
export function sessionPlayers(session: Session): Player[] {
  const byId = new Map<number, Player>();
  for (const slot of slotsOf(session)) {
    for (const player of playersOf(slot)) byId.set(player.id, player);
  }
  return [...byId.values()].sort((a, b) => a.name.localeCompare(b.name));
}

export function sessionCost(session: Session): number {
  return slotsOf(session).reduce(
    (total, slot) => total + slot.courts_booked * session.court_price,
    0,
  );
}

// Cost of slots that someone actually played in; empty slots aren't billed.
export function billedCost(session: Session): number {
  return slotsOf(session)
    .filter((slot) => playersOf(slot).length > 0)
    .reduce((total, slot) => total + slot.courts_booked * session.court_price, 0);
}

export function perPlayerCost(session: Session, slot: TimeSlot): number | null {
  const count = playersOf(slot).length;
  if (count === 0) return null;
  return (slot.courts_booked * session.court_price) / count;
}
