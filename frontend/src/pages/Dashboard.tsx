import { useNavigate } from "react-router-dom";

function Dashboard() {
  const navigate = useNavigate();

  return (
    <div>
      <h1>PlayLedger Dashboard</h1>

      <button onClick={() => navigate("/create-session")}>
        Create Session
      </button>
    </div>
  );
}

export default Dashboard;