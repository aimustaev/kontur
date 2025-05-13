import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { TicketList } from './components/TicketList/TicketList';
import { TicketPage } from './pages/TicketPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<TicketList />} />
        <Route path="/tickets/:ticketId" element={<TicketPage />} />
      </Routes>
    </Router>
  );
}

export default App; 