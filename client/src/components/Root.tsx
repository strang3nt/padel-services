import { App } from '@/components/App.tsx';
import { ErrorBoundary } from '@/components/ErrorBoundary.tsx';
import { GenericErrorPage } from '@/pages/GenericErrorPage';
import { IconContext } from 'react-icons/lib';

export function Root() {
  return (
    <IconContext.Provider value={{ size: "20px", className: "global-class-name" }}>
      <ErrorBoundary fallback={GenericErrorPage}>
        <App />
      </ErrorBoundary>
    </IconContext.Provider>
  );
}
