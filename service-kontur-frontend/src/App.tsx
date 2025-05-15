import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { TicketList } from './components/TicketList/TicketList';
import { TicketPage } from './pages/TicketPage';
import WorkflowBuilder from './pages/workflow';
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<TicketList />} />
        <Route path="/tickets/:ticketId" element={<TicketPage />} />
        <Route path="/workflow" element={<WorkflowBuilder />} />
      </Routes>
    </Router>
  );
}

export default App; 