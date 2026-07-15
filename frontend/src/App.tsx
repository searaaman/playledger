import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import CreateSession from "./pages/CreateSession";
import Session from "./pages/Session";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/create-session" element={<CreateSession />} />
        <Route path="/sessions/:id" element={<Session />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;