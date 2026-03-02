import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { PropsWithChildren } from "react";

function Section({ children, title }: PropsWithChildren<{ title: string }>) {
  return (
    <Box
      component="section"
      sx={{
        width: "100%",
        py: 4,
        px: 2,
        boxSizing: "border-box",
      }}
    >
      <Stack spacing={3} alignItems="stretch" sx={{ width: "100%" }}>
        <Typography
          variant="h5"
          component="h1"
          align="center"
          sx={{
            fontWeight: 700,
            letterSpacing: "0.02em",
            textTransform: "capitalize",
            mb: 1,
          }}
        >
          {title}
        </Typography>

        <Box sx={{ width: "100%" }}>{children}</Box>
      </Stack>
    </Box>
  );
}

export default Section;
