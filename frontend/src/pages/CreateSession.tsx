import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api/axios";

function CreateSession() {
  const navigate = useNavigate();

  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [courtPrice, setCourtPrice] = useState("");

  const handleCreate = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await api.post(
        "/sessions",
        {
          StartTime: startTime,
          EndTime: endTime,
          court_price: Number(courtPrice),
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      console.log("Session Created:", response.data);

      console.log(response.data);
      console.log(response.data.ID);
      console.log(response.data.id);
      navigate(`/sessions/${response.data.id}`);

      
      setStartTime("");
      setEndTime("");
      setCourtPrice("");
    } catch (err) {
      console.error("Failed to create session:", err);
    }
  };

  return (
    <div>
      <h1>Create Session</h1>

      <input
        type="datetime-local"
        value={startTime}
        onChange={(e) => setStartTime(e.target.value)}
      />

      <br />
      <br />

      <input
        type="datetime-local"
        value={endTime}
        onChange={(e) => setEndTime(e.target.value)}
      />

      <br />
      <br />

      <input
        type="number"
        placeholder="Court Price"
        value={courtPrice}
        onChange={(e) => setCourtPrice(e.target.value)}
      />

      <br />
      <br />

      <button onClick={handleCreate}>Create Session</button>
    </div>
  );
}

export default CreateSession;