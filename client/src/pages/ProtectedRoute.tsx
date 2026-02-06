import { Outlet } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';
import { Placeholder, Spinner } from '@telegram-apps/telegram-ui';

export const ProtectedRoute = () => {
  const { loading, error } = useAuth();

  if (error) throw new Error(`${error}`)
  if (loading) return <Placeholder
    description="Description"
    header="Validating session"
  >
    <Spinner size="l" />
  </Placeholder>

  return <Outlet />;
};