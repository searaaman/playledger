import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../api/axios";

function Session() {
  const { id } = useParams();

  const [session, setSession] = useState<any>(null);

  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [courtsBooked, setCourtsBooked] = useState("");

  const fetchSession = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await api.get(`/sessions/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      setSession(response.data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchSession();
  }, []);

  const handleCreateTimeSlot = async () => {
    try {
      const token = localStorage.getItem("token");

      await api.post(
        `/sessions/${id}/timeslots`,
        {
          StartTime: startTime,
          EndTime: endTime,
          CourtsBooked: Number(courtsBooked),
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      
      await fetchSession();

      
      setStartTime("");
      setEndTime("");
      setCourtsBooked("");
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div>
      <h1>Session {id}</h1>

      <h2>Create TimeSlot</h2>

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
        placeholder="Courts Booked"
        value={courtsBooked}
        onChange={(e) => setCourtsBooked(e.target.value)}
      />

      <br />
      <br />

      <button onClick={handleCreateTimeSlot}>
        Create TimeSlot
      </button>

      <hr />

      <h2>Time Slots</h2>

      {session?.TimeSlots?.length === 0 && (
        <p>No TimeSlots created yet.</p>
      )}

      {session?.TimeSlots?.map((slot: any) => (
        <div
          key={slot.ID}
          style={{
            border: "1px solid #ccc",
            padding: "10px",
            marginBottom: "10px",
          }}
        >
          <p>
            <strong>Start:</strong> {slot.StartTime}
          </p>

          <p>
            <strong>End:</strong> {slot.EndTime}
          </p>

          <p>
            <strong>Courts Booked:</strong> {slot.CourtsBooked}
          </p>
        </div>
      ))}
    </div>
  );
}

export default Session;