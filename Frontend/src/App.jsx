import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";
import RatingPage from "./pages/RatingPage";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import AdminPanel from "./pages/AdminPanel";
import CleanerPanel from "./pages/CleanerPanel";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<HomePage />} />
        <Route path="/admin" element={<AdminPanel />} />
        <Route path="/cleaner" element={<CleanerPanel />} />
        <Route path="/rating" element={<RatingPage />} />
        <Route path="/rating/:toiletId" element={<RatingPage />} />
      </Routes>
    </Router>
  );
}

export default App;
