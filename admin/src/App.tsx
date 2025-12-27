import {  Routes, Route, Navigate} from 'react-router-dom';
import LoginPage from './pages/Login';
import DashboardPage from './pages/Dashboard';
import UsersPage from './pages/Users/UsersPage'; 
import RolesPage from './pages/Roles/RolesPage'; // Placeholder
import { useAuthStore } from './hooks/useAuth';
import ProtectedRoute from './guards/ProtectedRoute';
import ProfilePage from './pages/Profile';

const App = () => {
const { isAuthenticated } = useAuthStore();
console.log('AUTH STATE:', isAuthenticated);
console.log('CURRENT PATH:', window.location.pathname);
console.log('AUTH:', isAuthenticated);
  return (
      <Routes>
        <Route path="/login" element={!isAuthenticated ? <LoginPage /> : <Navigate to="/dashboard" />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/roles" element={<RolesPage />} />
          <Route path="/profile" element={<ProfilePage />} />
        </Route>
        <Route path="*" element={<Navigate to={isAuthenticated ? "/dashboard" : "/login"} />} />
      </Routes>
   
  );
};

export default App;