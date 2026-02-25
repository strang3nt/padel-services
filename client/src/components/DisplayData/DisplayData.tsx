import Checkbox from '@mui/material/Checkbox';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import ListItem from '@mui/material/ListItem';
import type { FC, ReactNode } from 'react';
import { Link } from '@/components/Link/Link.tsx';
import { bem } from '@/css/bem.ts';
import './DisplayData.css';

const [, e] = bem('display-data');

export type DisplayDataRow =
  & { title: string }
  & (
    | { type: 'link'; value?: string }
    | { value: ReactNode }
  )

export interface DisplayDataProps {
  header?: ReactNode;
  footer?: ReactNode;
  rows: DisplayDataRow[];
}

export const DisplayData: FC<DisplayDataProps> = ({ header, rows }) => (
  <List
    component="nav"
    aria-labelledby="nested-list-subheader"
    subheader={header}
  >
    {rows.map((item, _) => {
      let valueNode: ReactNode;

      if (item.value === undefined) {
        valueNode = <i>empty</i>;
      } else {
        if ('type' in item) {
          valueNode = <Link to={item.value}>Open</Link>;
        } else if (typeof item.value === 'string') {
          valueNode = item.value;
        } else if (typeof item.value === 'boolean') {
          valueNode = <Checkbox checked={item.value} disabled />;
        } else {
          valueNode = item.value;
        }
      }

      return (
        <ListItemButton
          classes={e('line')}
        >
          <ListItemText>{item.title}</ListItemText>
          <ListItem classes={e('line-value')}>
            {valueNode}
          </ListItem>
        </ListItemButton>
      );
    })}
  </List>
);
