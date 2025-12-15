import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '../hooks/useAuth';
import MainLayout from '../components/common/MainLayout';

const ProtectedRoute = () => {
  const { isAuthenticated, isInitialized } = useAuthStore();

  if (!isInitialized) {
    return (
      <div className="h-screen flex items-center justify-center">
        <span>Loading authentication...</span>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return (
    <MainLayout>
      <Outlet />
    </MainLayout>
  );
};

export default ProtectedRoute;
