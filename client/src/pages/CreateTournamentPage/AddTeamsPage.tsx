import { useState, type FC } from "react";
import { useLoaderData, useLocation, useNavigate } from "react-router-dom";
import { genderToString } from "./utils";
import { Team } from "@/pages/RetrieveTournamentPage";
import { useAuth } from "@/components/AuthProvider";
import { TournamentSetupData } from "./CreateTournamentPage";

import Divider from "@mui/material/Divider";
import Box from "@mui/material/Box";
import InputLabel from "@mui/material/InputLabel";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import List from "@mui/material/List";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import ListItemText from "@mui/material/ListItemText";
import ListItem from "@mui/material/ListItem";
import CloseIcon from "@mui/icons-material/Close";
import FormControl from "@mui/material/FormControl";
import IconButton from "@mui/material/IconButton";
import Snackbar from "@mui/material/Snackbar";
import { Link } from "@/components/Link/Link.tsx";

interface NotificationContent {
  title: string;
  description: string;
  onClose?: () => void;
}

export const tournamentStore = {
  teams: [] as Team[],
  addTeam: (team: Team) => {
    tournamentStore.teams.push(team);
  },
  removeTeam: (index: number) => {
    tournamentStore.teams.splice(index, 1);
  },
  getTeams: () => tournamentStore.teams,
};

export function teamsLoader() {
  const teams = tournamentStore.getTeams();
  return { teams };
}

export const AddPlayersPage: FC = () => {
  const loaderData = useLoaderData() as { teams: Team[] } | null;
  const teams = loaderData?.teams || [];
  const [_, setRefresh] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { bearerToken } = useAuth();
  const [open, setOpen] = useState<null | NotificationContent>(null);

  const config = location.state as TournamentSetupData | null;

  const handleDeleteTeam = (index: number) => {
    tournamentStore.removeTeam(index);
    setRefresh((prev) => !prev);
  };

  if (!config) {
    return (
      <Box sx={{ p: 4, textAlign: "center" }}>
        <ListItemText
          primary="Missing tournament configuration."
          secondary="Please go back and fill out the setup page first."
        />
        <Button
          variant="contained"
          onClick={() => navigate("/create-tournament")}
        >
          Go to Setup
        </Button>
      </Box>
    );
  }

  const isRosterFull = teams.length >= config.numberOfTeams;

  const handleSendTournament = () => {
    const dateStart = new Date(config.tournamentDate);

    fetch(
      `/api/create-tournament?tournamentType=${config.selectedTournament}&dateStart=${dateStart.toISOString()}&totalRounds=${config.roundsNumber}&availableCourts=${config.availableCourts}`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${bearerToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(teams),
      },
    )
      .then((response) => {
        if (response.ok) {
          setOpen({
            title: "Success",
            description: "Tournament created successfully",
            onClose: () => navigate("/"),
          });
        } else {
          setOpen({ title: "Failed", description: "Try again later" });
        }
      })
      .catch((error) => {
        console.error("Network error:", error);
        setOpen({ title: "Error", description: "Could not reach the server." });
      });
  };

  return (
    <>
      <Snackbar
        open={open != null}
        autoHideDuration={3000}
        onClose={() => {
          setOpen(null);
          open?.onClose?.();
        }}
        message={`${open?.description}`}
      />

      {!isRosterFull ? (
        <Link to="/create-tournament/add-team" state={config}>
          <Button type="submit" size="large" fullWidth variant="outlined">
            Add Team
          </Button>
        </Link>
      ) : (
        <Box
          sx={{
            p: 2,
            textAlign: "center",
            bgcolor: "#f5f5f5",
            borderRadius: 1,
          }}
        >
          <ListItemText
            primary="Roster is full!"
            secondary="Proceed to finalize the tournament."
          />
        </Box>
      )}
      <Divider textAlign="center">
        Teams Added ({teams.length} / {config.numberOfTeams})
      </Divider>

      <List>
        {teams.length === 0 ? (
          <ListItem>
            <ListItemText primary="No teams added yet" />
          </ListItem>
        ) : (
          teams.map(
            (
              {
                person1: { id: teammate1 },
                person2: { id: teammate2 },
                gender,
              },
              i,
            ) => (
              <ListItem key={i}>
                <ListItemText
                  primary={`${teammate1}, ${teammate2}`}
                  secondary={`${genderToString(gender)} team`}
                />
                <IconButton onClick={() => handleDeleteTeam(i)}>
                  <CloseIcon />
                </IconButton>
              </ListItem>
            ),
          )
        )}
      </List>

      <Button
        variant="contained"
        fullWidth
        disabled={teams.length < config.numberOfTeams}
        onClick={handleSendTournament}
      >
        Save Tournament
      </Button>
    </>
  );
};

export const AddTeamPage: FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const data = location.state as TournamentSetupData | null;

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const formData = new FormData(event.currentTarget);

    const teammate1 = formData.get("teammate1") as string;
    const teammate2 = formData.get("teammate2") as string;
    const gender = formData.get("gender") as string;

    if (!teammate1 || !teammate2 || gender === "") return;

    const newTeam: Team = {
      person1: { id: teammate1 },
      person2: { id: teammate2 },
      gender: parseInt(gender, 10),
    };

    tournamentStore.addTeam(newTeam);

    navigate("/create-tournament/add-players", {
      state: data,
      replace: true,
    });
  };

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{ width: "100%", display: "flex", flexDirection: "column", gap: 2 }}
    >
      <TextField
        label="First teammate"
        name="teammate1"
        variant="outlined"
        required
      />
      <TextField
        label="Second teammate"
        name="teammate2"
        variant="outlined"
        required
      />
      <FormControl required>
        <InputLabel id="gender-label">Gender</InputLabel>
        <Select
          labelId="gender-label"
          label="Gender"
          name="gender"
          defaultValue=""
        >
          <MenuItem value={0}>Male</MenuItem>
          <MenuItem value={1}>Female</MenuItem>
          <MenuItem value={2}>Mixed</MenuItem>
        </Select>
      </FormControl>

      <Button type="submit" size="large" fullWidth variant="outlined">
        Add Team
      </Button>
    </Box>
  );
};
