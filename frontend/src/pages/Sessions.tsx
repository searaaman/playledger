import { Link } from "react-router-dom";
import { listSessions } from "../api/endpoints";
import SessionList from "../components/SessionList";
import { ErrorBox, Loading } from "../components/Status";
import { useLoad } from "../lib/useLoad";

function Sessions() {
  const { data: sessions, error } = useLoad(listSessions);

  return (
    <div>
      <div className="page-head">
        <div>
          <h1>Sessions</h1>
          <p className="muted">Every booking, newest first.</p>
        </div>
        <Link to="/sessions/new" className="btn">
          + New session
        </Link>
      </div>
      <ErrorBox message={error} />
      {sessions ? (
        <section className="card card-flush">
          <SessionList sessions={sessions} />
        </section>
      ) : (
        !error && <Loading />
      )}
    </div>
  );
}

export default Sessions;
