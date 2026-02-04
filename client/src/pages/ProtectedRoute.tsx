import { Outlet } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';

export const ProtectedRoute = () => {
  const { loading, error } = useAuth();

  if (error) throw new Error(`${error}`)
  if (loading) return <div>Validating session...</div>;
  
  return <Outlet />; // Success case: Render the child routes
};