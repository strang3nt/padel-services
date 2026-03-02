import { Outlet } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';
import LoadingPage from '@/components/LoadingPage';

export const ProtectedRoute = () => {
  const { loading, error } = useAuth();

  if (error) throw new Error(`${error}`)
  if (loading) return <LoadingPage
    title=""
    message="Validating session"
  />

  return <Outlet />;
};

