import { Placeholder, AppRoot } from '@telegram-apps/telegram-ui';
import { retrieveLaunchParams, isColorDark, isRGB } from '@tma.js/sdk-react';
import { FC, useMemo } from 'react';
import errorImage from '../../assets/error.svg'

export const GenericErrorPage: FC<{ error: unknown }> = ({ error }) => {

  const [platform, isDark] = useMemo(() => {
    try {
      const lp = retrieveLaunchParams();
      const { bg_color: bgColor } = lp.tgWebAppThemeParams;
      return [lp.tgWebAppPlatform, bgColor && isRGB(bgColor) ? isColorDark(bgColor) : false];
    } catch {
      return ['android', false];
    }
  }, []);

  return (
    <AppRoot
      appearance={isDark ? 'dark' : 'light'}
      platform={['macos', 'ios'].includes(platform) ? 'ios' : 'base'}
    >
      <Placeholder
        header="Oops"
        description={error instanceof Error
          ? error.message
          : typeof error === 'string'
            ? error
            : JSON.stringify(error)}
      >
        <img
          alt="Error image"
          src={errorImage}
          style={{ display: 'block', width: '144px', height: '144px' }}
        />
      </Placeholder>
    </AppRoot>
  );
}