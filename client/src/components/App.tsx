import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { useLaunchParams, useSignal, miniApp } from '@tma.js/sdk-react';
import { AppRoot } from '@telegram-apps/telegram-ui';

import { routes } from '@/navigation/routes.tsx';
import { AuthProvider } from './AuthProvider';
import { ErrorBoundary } from './ErrorBoundary';
import { ErrorBoundaryError } from './Root';

export function App() {
  const lp = useLaunchParams();
  const isDark = useSignal(miniApp.isDark);

  return (
    <AppRoot
      appearance={isDark ? 'dark' : 'light'}
      platform={['macos', 'ios'].includes(lp.tgWebAppPlatform) ? 'ios' : 'base'}
    >
      <ErrorBoundary fallback={ErrorBoundaryError}>
        <AuthProvider>
          <RouterProvider router={createBrowserRouter(routes)} />
        </AuthProvider>
      </ErrorBoundary>
    </AppRoot>
  );
}
