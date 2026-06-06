import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Overview from './pages/Overview';
import HostDetail from './pages/HostDetail';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Overview />} />
          <Route path="host/:hostname" element={<HostDetail />} />
          <Route path="hosts" element={<div className="p-8 font-bold text-2xl">Hosts Management (Coming Soon)</div>} />
          <Route path="alerts" element={<div className="p-8">Alerts Page (Coming Soon)</div>} />
          <Route path="settings" element={<div className="p-8">Settings Page (Coming Soon)</div>} />
          <Route path="*" element={<div>Not Found</div>} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
