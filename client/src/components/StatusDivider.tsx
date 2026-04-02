import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";

import { type FC } from "react";

interface StatusDividerProps {
  label: string;
  current: number;
  total: number;
  isFull: boolean;
}

export const StatusDivider: FC<StatusDividerProps> = ({
  label,
  current,
  total,
  isFull,
}) => {
  return (
    <Divider
      textAlign="center"
      sx={{
        my: 4,
        "&::before, &::after": { borderColor: "divider" },
      }}
    >
      <Box
        sx={{
          px: 2,
          py: 0.5,
          borderRadius: "20px",
          bgcolor: isFull ? "success.light" : "action.selected",
          color: isFull ? "success.contrastText" : "text.secondary",
          fontSize: "0.75rem",
          fontWeight: "bold",
          textTransform: "uppercase",
          letterSpacing: 1,
          display: "flex",
          alignItems: "center",
          gap: 1,
          border: "1px solid",
          borderColor: isFull ? "success.main" : "divider",
        }}
      >
        {label}
        <Box component="span" sx={{ opacity: 0.8, ml: 0.5 }}>
          ({current} / {total})
        </Box>
      </Box>
    </Divider>
  );
};
