import { useState, type FormEvent } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { errorMessage } from "../api/client";
import {
  assignPlayer,
  createPayment,
  createPlayer,
  createTimeSlot,
  deletePayment,
  deleteSession,
  deleteTimeSlot,
  getLedger,
  getSession,
  getSessionBilling,
  listPayments,
  listPlayers,
  removePlayer,
  updateTimeSlotCourts,
} from "../api/endpoints";
import Balance from "../components/Balance";
import { ErrorBox, Loading } from "../components/Status";
import { date, dateTime, isZero, money, time, timeRange } from "../lib/format";
import { billedCost, perPlayerCost, playersOf, sessionPlayers, slotsOf } from "../lib/session";
import { useLoad } from "../lib/useLoad";
import type { Player, Session, TimeSlot } from "../types";

function SessionDetail() {
  const { id } = useParams();
  const sessionId = Number(id);
  const navigate = useNavigate();

  const { data, error: loadError, reload } = useLoad(
    () =>
      Promise.all([
        getSession(sessionId),
        getSessionBilling(sessionId),
        listPayments({ session_id: sessionId }),
        listPlayers(),
        getLedger(),
      ]),
    sessionId,
  );
  const [error, setError] = useState<string | null>(null);
  // Cells, slots or rows with a request in flight, so we can disable them.
  const [busy, setBusy] = useState<Set<string>>(new Set());
  // Players added to the grid who haven't been ticked into any slot yet.
  const [extraRows, setExtraRows] = useState<Player[]>([]);

  if (loadError) return <ErrorBox message={loadError} />;
  if (!data) return <Loading />;

  const [session, bills, payments, allPlayers, ledger] = data;
  const slots = slotsOf(session);

  const run = async (key: string, action: () => Promise<unknown>) => {
    setBusy((b) => new Set(b).add(key));
    setError(null);
    try {
      await action();
      reload();
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy((b) => {
        const next = new Set(b);
        next.delete(key);
        return next;
      });
    }
  };

  const inSlot = (slot: TimeSlot, playerId: number) =>
    playersOf(slot).some((p) => p.id === playerId);

  const toggle = (slot: TimeSlot, player: Player) =>
    run(`cell-${slot.id}-${player.id}`, () =>
      inSlot(slot, player.id) ? removePlayer(slot.id, player.id) : assignPlayer(slot.id, player.id),
    );

  const setRow = (player: Player, on: boolean) =>
    run(`row-${player.id}`, () =>
      Promise.all(
        slots
          .filter((slot) => inSlot(slot, player.id) !== on)
          .map((slot) => (on ? assignPlayer(slot.id, player.id) : removePlayer(slot.id, player.id))),
      ),
    );

  const gridPlayers = [...sessionPlayers(session)];
  for (const player of extraRows) {
    if (!gridPlayers.some((p) => p.id === player.id)) gridPlayers.push(player);
  }
  const available = allPlayers.filter((p) => !gridPlayers.some((g) => g.id === p.id));

  const paidInSession = payments.reduce((total, p) => total + p.amount, 0);
  const billedTotal = billedCost(session);
  const balanceOf = (playerId: number) =>
    ledger.find((entry) => entry.player_id === playerId)?.balance ?? 0;
  const paidBy = (playerId: number) =>
    payments.filter((p) => p.player_id === playerId).reduce((total, p) => total + p.amount, 0);

  const handleDeleteSession = async () => {
    if (!confirm("Delete this session, its slots and attendance?")) return;
    try {
      await deleteSession(session.id);
      navigate("/sessions");
    } catch (err) {
      setError(errorMessage(err));
    }
  };

  return (
    <div className="stack">
      <div className="page-head">
        <div>
          <p className="small">
            <Link to="/sessions">← Sessions</Link>
          </p>
          <h1>{date(session.start_time)}</h1>
          <p className="muted">
            {timeRange(session.start_time, session.end_time)} · {money(session.court_price)} per court per hour
          </p>
        </div>
        <button className="btn btn-danger btn-sm" onClick={handleDeleteSession}>
          Delete session
        </button>
      </div>

      <ErrorBox message={error} />

      <div className="stats">
        <div className="stat">
          <div className="stat-label">Players</div>
          <div className="stat-value">{sessionPlayers(session).length}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Billed</div>
          <div className="stat-value">{money(billedTotal)}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Collected</div>
          <div className="stat-value">{money(paidInSession)}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Still to collect</div>
          <div className={`stat-value ${billedTotal - paidInSession > 0.005 ? "owes" : ""}`}>
            {money(Math.max(billedTotal - paidInSession, 0))}
          </div>
        </div>
      </div>

      <section className="card card-flush">
        <div className="card-head">
          <div>
            <h2>Attendance</h2>
            <p className="muted small" style={{ margin: "2px 0 0" }}>
              Tap a cell to mark who played in each slot. Each slot's court cost is split between its players.
            </p>
          </div>
        </div>

        {slots.length === 0 ? (
          <p className="empty">No time slots yet. Add one below.</p>
        ) : (
          <div className="attendance" style={{ marginTop: 12 }}>
            <table>
              <thead>
                <tr>
                  <th className="player-col">Player</th>
                  {slots.map((slot) => (
                    <th key={slot.id}>
                      <SlotHeader
                        slot={slot}
                        disabled={busy.has(`slot-${slot.id}`)}
                        onCourts={(courts) =>
                          run(`slot-${slot.id}`, () => updateTimeSlotCourts(slot.id, courts))
                        }
                        onDelete={() => {
                          if (confirm(`Delete the ${timeRange(slot.start_time, slot.end_time)} slot?`)) {
                            run(`slot-${slot.id}`, () => deleteTimeSlot(slot.id));
                          }
                        }}
                      />
                    </th>
                  ))}
                  <th />
                </tr>
              </thead>
              <tbody>
                {gridPlayers.length === 0 && (
                  <tr>
                    <td className="empty" colSpan={slots.length + 2}>
                      Add the players who turned up.
                    </td>
                  </tr>
                )}
                {gridPlayers.map((player) => {
                  const played = slots.filter((slot) => inSlot(slot, player.id)).length;
                  const rowBusy = busy.has(`row-${player.id}`);
                  return (
                    <tr key={player.id}>
                      <td className="player-col">
                        <div className="player-name">
                          <Link to={`/players/${player.id}`}>{player.name}</Link>
                        </div>
                      </td>
                      {slots.map((slot) => {
                        const on = inSlot(slot, player.id);
                        return (
                          <td key={slot.id}>
                            <button
                              className={`cell ${on ? "on" : ""}`}
                              aria-pressed={on}
                              aria-label={`${player.name} ${timeRange(slot.start_time, slot.end_time)}`}
                              disabled={rowBusy || busy.has(`cell-${slot.id}-${player.id}`)}
                              onClick={() => toggle(slot, player)}
                            >
                              ●
                            </button>
                          </td>
                        );
                      })}
                      <td>
                        <button
                          className="btn btn-ghost btn-sm"
                          disabled={rowBusy}
                          onClick={() => setRow(player, played < slots.length)}
                        >
                          {played < slots.length ? "All" : "None"}
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
              <tfoot>
                <tr>
                  <td className="player-col">Each pays</td>
                  {slots.map((slot) => {
                    const cost = perPlayerCost(session, slot);
                    return (
                      <td key={slot.id} className="num" style={{ textAlign: "center" }}>
                        {cost === null ? "—" : money(cost)}
                        <div>{playersOf(slot).length} playing</div>
                      </td>
                    );
                  })}
                  <td />
                </tr>
              </tfoot>
            </table>
          </div>
        )}

        <AddPlayerRow
          available={available}
          onAdd={(player) => setExtraRows((rows) => [...rows, player])}
          onError={setError}
        />
      </section>

      <AddSlotForm session={session} onCreated={reload} onError={setError} />

      <section className="card card-flush">
        <div className="card-head">
          <div>
            <h2>Billing & payments</h2>
            <p className="muted small" style={{ margin: "2px 0 0" }}>
              Overall balance includes anything carried over from earlier sessions.
            </p>
          </div>
        </div>
        {bills.length === 0 ? (
          <p className="empty">Nobody has been marked as playing yet.</p>
        ) : (
          <div className="table-wrap" style={{ marginTop: 12 }}>
            <table>
              <thead>
                <tr>
                  <th>Player</th>
                  <th className="num">This session</th>
                  <th className="num">Paid here</th>
                  <th className="num">Overall</th>
                  <th className="num">Record payment</th>
                </tr>
              </thead>
              <tbody>
                {bills.map((bill) => {
                  const balance = balanceOf(bill.player_id);
                  return (
                    <tr key={bill.player_id}>
                      <td>
                        <Link to={`/players/${bill.player_id}`}>{bill.name}</Link>
                      </td>
                      <td className="num">{money(bill.amount)}</td>
                      <td className="num">{money(paidBy(bill.player_id))}</td>
                      <td className="num">
                        <Balance amount={balance} />
                      </td>
                      <td>
                        <PayForm
                          suggested={balance < 0 && !isZero(balance) ? -balance : 0}
                          disabled={busy.has(`pay-${bill.player_id}`)}
                          onPay={(amount) =>
                            run(`pay-${bill.player_id}`, () =>
                              createPayment(bill.player_id, session.id, amount),
                            )
                          }
                        />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}

        {payments.length > 0 && (
          <div style={{ borderTop: "1px solid var(--border)" }}>
            <div className="card-head">
              <h3>Payments recorded for this session</h3>
            </div>
            <div className="table-wrap">
              <table>
                <tbody>
                  {payments.map((payment) => (
                    <tr key={payment.id}>
                      <td>{payment.player?.name ?? `Player #${payment.player_id}`}</td>
                      <td className="muted small">{dateTime(payment.created_at)}</td>
                      <td className="num">{money(payment.amount)}</td>
                      <td className="num">
                        <button
                          className="btn btn-ghost btn-sm"
                          disabled={busy.has(`delpay-${payment.id}`)}
                          onClick={() => {
                            if (confirm(`Delete ${money(payment.amount)} payment?`)) {
                              run(`delpay-${payment.id}`, () => deletePayment(payment.id));
                            }
                          }}
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </section>
    </div>
  );
}

function SlotHeader({
  slot,
  disabled,
  onCourts,
  onDelete,
}: {
  slot: TimeSlot;
  disabled: boolean;
  onCourts: (courts: number) => void;
  onDelete: () => void;
}) {
  return (
    <div className="slot-head">
      <span className="slot-time">
        {time(slot.start_time)}–{time(slot.end_time)}
      </span>
      <span className="courts">
        <button
          className="icon-btn"
          aria-label="One court fewer"
          disabled={disabled || slot.courts_booked === 0}
          onClick={() => onCourts(slot.courts_booked - 1)}
        >
          −
        </button>
        {slot.courts_booked} court{slot.courts_booked === 1 ? "" : "s"}
        <button
          className="icon-btn"
          aria-label="One court more"
          disabled={disabled}
          onClick={() => onCourts(slot.courts_booked + 1)}
        >
          +
        </button>
      </span>
      <button className="btn btn-ghost btn-sm" disabled={disabled} onClick={onDelete}>
        Remove
      </button>
    </div>
  );
}

function AddPlayerRow({
  available,
  onAdd,
  onError,
}: {
  available: Player[];
  onAdd: (player: Player) => void;
  onError: (message: string) => void;
}) {
  const [selected, setSelected] = useState("");
  const [newName, setNewName] = useState("");
  const [busy, setBusy] = useState(false);

  const addExisting = () => {
    const player = available.find((p) => p.id === Number(selected));
    if (player) onAdd(player);
    setSelected("");
  };

  const addNew = async (e: FormEvent) => {
    e.preventDefault();
    if (!newName.trim()) return;
    setBusy(true);
    try {
      onAdd(await createPlayer(newName.trim(), ""));
      setNewName("");
    } catch (err) {
      onError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="add-player">
      <select value={selected} onChange={(e) => setSelected(e.target.value)} aria-label="Add a player">
        <option value="">Add a player…</option>
        {available.map((player) => (
          <option key={player.id} value={player.id}>
            {player.name}
          </option>
        ))}
      </select>
      <button className="btn btn-secondary" disabled={!selected} onClick={addExisting}>
        Add
      </button>
      <form className="row" style={{ flex: 1, flexWrap: "nowrap" }} onSubmit={addNew}>
        <input
          placeholder="or a new player's name"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
        />
        <button className="btn btn-secondary" disabled={busy || !newName.trim()}>
          Create
        </button>
      </form>
    </div>
  );
}

function AddSlotForm({
  session,
  onCreated,
  onError,
}: {
  session: Session;
  onCreated: () => void;
  onError: (message: string) => void;
}) {
  const [open, setOpen] = useState(false);
  const [start, setStart] = useState("");
  const [end, setEnd] = useState("");
  const [courts, setCourts] = useState("1");
  const [busy, setBusy] = useState(false);

  // Slot times are entered as times of day and placed on the session's date,
  // rolling past midnight if the session does.
  const onSessionDay = (value: string) => {
    const base = new Date(session.start_time);
    const [h, m] = value.split(":").map(Number);
    const d = new Date(base);
    d.setHours(h, m, 0, 0);
    if (d < base) d.setDate(d.getDate() + 1);
    return d;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setBusy(true);
    try {
      const startDate = onSessionDay(start);
      let endDate = onSessionDay(end);
      if (endDate <= startDate) endDate = new Date(endDate.getTime() + 86_400_000);
      await createTimeSlot(session.id, {
        start_time: startDate.toISOString(),
        end_time: endDate.toISOString(),
        courts_booked: Number(courts),
      });
      setStart("");
      setEnd("");
      setOpen(false);
      onCreated();
    } catch (err) {
      onError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  if (!open) {
    return (
      <div>
        <button className="btn btn-secondary btn-sm" onClick={() => setOpen(true)}>
          + Add time slot
        </button>
      </div>
    );
  }

  return (
    <form className="card" onSubmit={handleSubmit}>
      <div className="card-head">
        <h2>Add time slot</h2>
      </div>
      <div className="form-grid">
        <label>
          Starts
          <input type="time" required value={start} onChange={(e) => setStart(e.target.value)} />
        </label>
        <label>
          Ends
          <input type="time" required value={end} onChange={(e) => setEnd(e.target.value)} />
        </label>
        <label>
          Courts booked
          <input type="number" min="0" required value={courts} onChange={(e) => setCourts(e.target.value)} />
        </label>
      </div>
      <div className="form-actions">
        <button className="btn" disabled={busy}>
          Add slot
        </button>
        <button type="button" className="btn btn-secondary" onClick={() => setOpen(false)}>
          Cancel
        </button>
      </div>
    </form>
  );
}

function PayForm({
  suggested,
  disabled,
  onPay,
}: {
  suggested: number;
  disabled: boolean;
  onPay: (amount: number) => void;
}) {
  const [amount, setAmount] = useState("");

  const submit = (e: FormEvent) => {
    e.preventDefault();
    const value = Number(amount || suggested);
    if (value > 0) {
      onPay(value);
      setAmount("");
    }
  };

  return (
    <form className="pay-form" onSubmit={submit}>
      <input
        type="number"
        min="0"
        step="any"
        inputMode="decimal"
        placeholder={suggested > 0 ? String(Math.round(suggested * 100) / 100) : "₹"}
        value={amount}
        onChange={(e) => setAmount(e.target.value)}
        aria-label="Payment amount"
      />
      <button className="btn btn-sm" disabled={disabled || !(Number(amount) > 0 || suggested > 0)}>
        Paid
      </button>
    </form>
  );
}

export default SessionDetail;
