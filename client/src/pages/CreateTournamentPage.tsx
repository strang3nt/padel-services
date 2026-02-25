import { useState, type FC } from 'react';
import { Form, Outlet, useLoaderData, LoaderFunctionArgs, replace, useNavigate } from 'react-router-dom';
import { Page } from '@/components/Page';
import { Team } from '@/pages/RetrieveTournamentPage';
import { useAuth } from '@/components/AuthProvider';
import { Link } from '@/components/Link/Link';

import Divider from '@mui/material/Divider';
import Box from '@mui/material/Box';
import InputLabel from '@mui/material/InputLabel';
import Snackbar from '@mui/material/Snackbar';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import List from '@mui/material/List';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Section from '@/components/Section';
import AddIcon from '@mui/icons-material/Add';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import ListItem from '@mui/material/ListItem';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';
import CloseIcon from '@mui/icons-material/Close';
import FormControl from '@mui/material/FormControl';
import IconButton from '@mui/material/IconButton';

export const CreateTournamentPage: FC = () => {
  return <Section
    title="Tournament creation"
  >
    <Outlet />
  </Section>

};

function genderToString(g: number): string {
  switch (g) {
    case 0:
      return "Male"
    case 1:
      return "Female"
    case 2:
      return "Mixed"
    default:
      throw Error("Team gender not recognized");
  }
}

export const AddTeamsPage: FC = () => {
  const loaderData = useLoaderData() as { teams: Team[] } | null;
  const teams = loaderData?.teams || [];
  const [_, setRefresh] = useState(false);

  const handleDeleteTeam = (index: number) => {
    tournamentStore.removeTeam(index);
    setRefresh((prev) => !prev);
  };

  return (
    <Page>
      <List>
        <Link to="/create-tournament/add-team">
          <ListItem disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <AddIcon />
              </ListItemIcon>
              <ListItemText primary="Add Team" />
            </ListItemButton>
          </ListItem>
        </Link>
        <Link to="/create-tournament/tournament-type" state={teams}>
          <ListItem disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <NavigateNextIcon />
              </ListItemIcon>
              <ListItemText primary="Next" />
            </ListItemButton>
          </ListItem>
        </Link>
      </List>
      <Divider textAlign="center">Teams Added</Divider>
      <List>
        {teams.length === 0 ? (
          <ListItem>
            <ListItemText primary="No teams added yet" />
          </ListItem>
        ) : (
          teams.map(({ person1: { id: teammate1 }, person2: { id: teammate2 }, gender }, i) => (
            <ListItem key={i}>
              <ListItemText
                primary={`${teammate1}, ${teammate2}`}
                secondary={`${genderToString(gender)} team`}
              />
              <IconButton onClick={() => handleDeleteTeam(i)}>
                <CloseIcon />
              </IconButton>
            </ListItem>
          ))
        )}
      </List>
    </Page >
  );
};

export const tournamentStore = {
  teams: [] as Team[],
  addTeam: (team: Team) => {
    tournamentStore.teams.push(team);
  },
  removeTeam: (index: number) => {
    tournamentStore.teams.splice(index, 1)
  },
  getTeams: () => tournamentStore.teams,
};

export async function teamsLoader() {
  const teams = tournamentStore.getTeams();
  return { teams };
}

export async function addTeamAction({ request }: LoaderFunctionArgs) {
  const formData = await request.formData();

  const newTeam: Team = {
    person1: { id: formData.get("teammate1") as string },
    person2: { id: formData.get("teammate2") as string },
    gender: parseInt(formData.get("gender") as string, 10),
  };

  tournamentStore.addTeam(newTeam);

  return replace('/create-tournament');
}

export const AddTeamPage: FC = () => {

  return (
    <Page>
      <Form method="post">
        <FormControl fullWidth variant="outlined">
          <Box sx={{ width: '100%', display: 'flex', flexDirection: 'column', gap: 2 }}>
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
            <FormControl fullWidth variant="outlined">
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
            <Button type="submit" variant="contained"
              color="primary"
              size="large"
              fullWidth>Add Team</Button>
          </Box>
        </FormControl>
      </Form>
    </Page >
  );
}

interface FormData {
  availableCourts: number;
  roundsNumber: number;
  tournamentDate: string;
}

interface NotificationContent {
  title: string;
  description: string;
  onClose?: () => void;
}

export const ChooseTournamentType: FC = () => {

  const getMatchesPerTeam = (
    teamsNumber: number,
    totalRounds: number,
    availableCourts: number): [number, number, number] => {

    let matchesPerTeam = totalRounds

    while (matchesPerTeam > 0) {

      const totalParticipations = teamsNumber * matchesPerTeam
      const totalMatchesFloat = totalParticipations / 2.0

      if (Number.isInteger(totalMatchesFloat)) {

        const matchesPerTurn = totalMatchesFloat / totalRounds

        if (matchesPerTurn <= availableCourts && teamsNumber > matchesPerTeam) {
          return [totalMatchesFloat, matchesPerTurn, matchesPerTeam]
        }
      }

      matchesPerTeam -= 1
    }

    return [0, 0.0, 0]
  }
  const [formData, setFormData] = useState<FormData>({
    availableCourts: 0,
    roundsNumber: 0,
    tournamentDate: "",
  })
  const { bearerToken } = useAuth()
  const [selectedTournament, setTournament] = useState("Rodeo");
  const navigate = useNavigate()
  const [open, setOpen] = useState<null | NotificationContent>(null);

  const renderSwitch = (): React.ReactNode => {
    switch (selectedTournament) {

      case "Rodeo":
        const sendRodeoTournament = (
          tournamentType: string,
          dateStart: Date,
          teams: Team[],
          totalRounds: number,
          availableCourts: number
        ) => {
          fetch(`/api/create-tournament?tournamentType=${tournamentType}&dateStart=${dateStart.toISOString()}&totalRounds=${totalRounds}&availableCourts=${availableCourts}`, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${bearerToken}`,
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(teams)
          }).then(response => {
            if (response.ok) {
              setOpen(
                {
                  title: "Tournament creation success",
                  description: "Tournament created and saved successfully",
                  onClose: () => { navigate("/") }
                }
              )
            } else {
              setOpen(
                {
                  title: "Tournament creation failed",
                  description: "Try again later",
                  onClose: () => { navigate("/") }
                }
              )
            }
          })
        }
        return <>
          <ListItem>
            <TextField
              label="Courts available"
              name="courtsAvailable"
              type="number"
              onChange={(e) => setFormData({
                ...formData,
                availableCourts: parseInt(e.target.value, 10)
              })}
              required
            />
          </ListItem>
          <ListItem>
            <TextField
              label="Number of rounds"
              name="roundsNumber"
              type="number"
              helperText=
              {
                (() => {
                  let totalMatches, matchesPerTurn, matchesPerTeam;
                  [totalMatches, matchesPerTurn, matchesPerTeam] = getMatchesPerTeam(tournamentStore.teams.length, formData.roundsNumber, formData.availableCourts)
                  if (totalMatches == 0) {
                    return "Configuration is not valid"
                  } else {
                    return `Current configuration translates to ${matchesPerTeam} matches per team and at most ${Math.ceil(matchesPerTurn)} matches per turn.`
                  }

                })()
              }
              onChange={(e) => setFormData({
                ...formData,
                roundsNumber: parseInt(e.target.value, 10)
              })}
              required
            />
          </ListItem>
          <ListItem>
            <Button fullWidth onClick={() => sendRodeoTournament(
              selectedTournament,
              new Date(formData.tournamentDate),
              tournamentStore.teams,
              formData.roundsNumber,
              formData.availableCourts

            )}>Send</Button>
          </ListItem>
        </>
      default:
        return "Tournament type not supported";
    }
  }

  return <Page>
    <Snackbar
      open={(() => open != null)()}
      autoHideDuration={3000}
      onClose={() => {
        setOpen(null);
        open?.onClose?.();
      }}
      message={`${open?.description}`}
    />
    <Form>
      <FormControl fullWidth variant="outlined">
        <List>
          <ListItem>
            <Select
              labelId="gender-label"
              label="Tournament Type"
              name="tournamentType"
              defaultValue="Rodeo"
              value={selectedTournament}
              onChange={e => setTournament(e.target.value)}
            >
              <MenuItem value={"Rodeo"}>Rodeo</MenuItem>
            </Select>
          </ListItem>
          <ListItem>
            <TextField
              label="Tournament date"
              name="tournamentDate"
              type="date"
              slotProps={{ inputLabel: { shrink: true } }}
              onChange={(e) => setFormData({
                ...formData,
                tournamentDate: e.target.value
              })
              }
              required
            />
          </ListItem>
          {
            renderSwitch()
          }
        </List>
      </FormControl>
    </Form>
  </Page>
}
