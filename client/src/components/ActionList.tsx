import { List, ListItem, ListItemText, IconButton } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";

interface ActionListProps<T> {
  items: T[];
  renderKey: (item: T, index: number) => string | number;
  getPrimaryText: (item: T) => string;
  onDelete: (index: number) => void;
  emptyMessage?: string;
}

export const ActionList = <T,>({
  items,
  renderKey,
  getPrimaryText,
  onDelete,
  emptyMessage = "No items added yet",
}: ActionListProps<T>) => {
  return (
    <List
      sx={{
        bgcolor: "background.paper",
        borderRadius: 2,
        mb: 2,
        border: "1px solid",
        borderColor: "divider",
        overflow: "hidden",
      }}
    >
      {items.length === 0 ? (
        <ListItem>
          <ListItemText
            primary={emptyMessage}
            sx={{ color: "text.secondary", textAlign: "center" }}
          />
        </ListItem>
      ) : (
        items.map((item, i) => (
          <ListItem
            key={renderKey(item, i)}
            divider={i !== items.length - 1}
            sx={{ px: 2 }}
          >
            <ListItemText primary={getPrimaryText(item)} />
            <IconButton onClick={() => onDelete(i)} size="small">
              <CloseIcon fontSize="small" />
            </IconButton>
          </ListItem>
        ))
      )}
    </List>
  );
};
