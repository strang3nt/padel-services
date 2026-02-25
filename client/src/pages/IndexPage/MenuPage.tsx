import { type FC } from 'react';

import { Link } from '@/components/Link/Link.tsx';
import { Page } from '@/components/Page.tsx';

import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import Section from '@/components/Section';
export const MenuPage: FC = () => (
  <Page back={false}>
    <Section title="Menu">
      <List>
        <ListItemButton>
          <Link to="/create-tournament">
            <ListItemText
              primary="Create tournament"
              secondary="Input data to create tournament" />
          </Link>
        </ListItemButton>

        <ListItemButton>
          <Link to="/retrieve-tournament">
            <ListItemText
              primary="Retrieve tournament"
              secondary="Retrieve tournament pairings and, if available, results" />
          </Link>
        </ListItemButton>
      </List>
    </Section>
  </Page >
);
