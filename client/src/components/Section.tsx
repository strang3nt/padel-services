import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { PropsWithChildren } from 'react';

function Section({ children, title }: PropsWithChildren<{ title: String }>) {
  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      sx={{ backgroundColor: 'background.default' }}
    >
      <Stack spacing={2} alignItems="center">
        <Box textAlign="center">
          <Typography variant="h6" component="h1" color="text.primary">
            {title}
          </Typography>
          {children}
        </Box>
      </Stack>
    </Box>
  );
};

export default Section
