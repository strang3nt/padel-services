import { Section, Cell } from '@telegram-apps/telegram-ui';
import { type FC } from 'react';

import { Link } from '@/components/Link/Link.tsx';
import { Page } from '@/components/Page.tsx';

export const MenuPage: FC = () => {

  return (
    <Page back={false}>
        <Section
          header="Menu"
          footer="Select item to access its functionality"
        >

          <Link to="/create-tournament">
            <Cell 
            subtitle="Input data to create tournament"
            >Create tournament
            </Cell>
          </Link>

          <Link to="/retrieve-tournament">
            <Cell 
            subtitle="Retrieve tournament pairings and, if available, results"
            >Retrieve tournament
            </Cell>
          </Link>
        </Section>
    </Page>
  );
};
