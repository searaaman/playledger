import { Link } from "react-router-dom";
import { money, timeRange } from "../lib/format";
import { sessionCost, sessionPlayers, slotsOf } from "../lib/session";
import type { Session } from "../types";

function SessionList({ sessions }: { sessions: Session[] }) {
  if (sessions.length === 0) {
    return <p className="empty">No sessions yet.</p>;
  }
  return (
    <div className="session-list">
      {sessions.map((session) => {
        const start = new Date(session.start_time);
        const slots = slotsOf(session).length;
        const players = sessionPlayers(session).length;
        return (
          <Link key={session.id} to={`/sessions/${session.id}`} className="session-item">
            <div className="date-badge">
              <div className="day">{start.getDate()}</div>
              <div className="mon">
                {start.toLocaleDateString("en-IN", { month: "short", weekday: "short" })}
              </div>
            </div>
            <div>
              <div style={{ fontWeight: 600 }}>{timeRange(session.start_time, session.end_time)}</div>
              <div className="muted small">
                {slots} slot{slots === 1 ? "" : "s"} · {players} player{players === 1 ? "" : "s"} ·{" "}
                {money(session.court_price)}/court/hr
              </div>
            </div>
            <div className="num" style={{ fontWeight: 600 }}>
              {money(sessionCost(session))}
            </div>
          </Link>
        );
      })}
    </div>
  );
}

export default SessionList;
