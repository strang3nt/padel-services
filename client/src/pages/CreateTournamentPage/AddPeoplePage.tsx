import { useState, type FC } from "react";
import { useLoaderData, useLocation, useNavigate } from "react-router-dom";
import { Person, Team } from "@/pages/tournament";
import { useAuth } from "@/components/AuthProvider";
import { TournamentSetupData } from "./CreateTournamentPage";

import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import ListItemText from "@mui/material/ListItemText";
import Snackbar from "@mui/material/Snackbar";
import { Link } from "@/components/Link/Link.tsx";
import { StatusDivider } from "@/components/StatusDivider";
import { ActionList } from "@/components/ActionList";

interface NotificationContent {
  title: string;
  description: string;
  onClose?: () => void;
}

export const peopleStore = {
  people: [] as Person[],
  addPerson: (person: Person) => {
    peopleStore.people.push(person);
  },
  removePerson: (index: number) => {
    peopleStore.people.splice(index, 1);
  },
  getPeople: () => peopleStore.people,
};

export function peopleLoader() {
  const people = peopleStore.getPeople();
  return { people };
}

function peopleToTeams(people: Person[]): Team[] {
  const emptyPerson: Person = { id: "" };
  const teams = [] as Team[];

  for (let i = 0; i < people.length; i += 2) {
    teams.push({
      person1: people[i],
      person2: people[i + 1],
      gender: 1,
    });
  }

  if (people.length % 2 == 1) {
    teams.push({
      person1: people[people.length - 1],
      person2: emptyPerson,
      gender: 1,
    });
  }
  return teams;
}

export const AddPeoplePage: FC = () => {
  const loaderData = useLoaderData() as { people: Person[] } | null;
  const people = loaderData?.people || [];
  const [_, setRefresh] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { bearerToken } = useAuth();
  const [open, setOpen] = useState<null | NotificationContent>(null);

  const config = location.state as TournamentSetupData | null;

  const handleDeletePerson = (index: number) => {
    peopleStore.removePerson(index);
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

  const isRosterFull = people.length >= config.numberOfTeams;

  const handleSendTournament = () => {
    const dateStart = new Date(config.tournamentDate);

    fetch(
      `/api/create-tournament?tournamentType=SinglePlayerRodeo&dateStart=${dateStart.toISOString()}&totalRounds=${config.roundsNumber}&availableCourts=${config.availableCourts}`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${bearerToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(peopleToTeams(people)),
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
        <Link to="/create-tournament/add-person" state={config}>
          <Button type="submit" size="large" fullWidth variant="outlined">
            Add Person
          </Button>
        </Link>
      ) : (
        <Box
          sx={{
            p: 2,
            textAlign: "center",
            borderRadius: 1,
          }}
        >
          <ListItemText
            primary="Roster is full!"
            secondary="Proceed to finalize the tournament."
          />
        </Box>
      )}
      <StatusDivider
        label="People"
        current={people.length}
        total={config.numberOfTeams}
        isFull={isRosterFull}
      />
      <ActionList
        items={people}
        renderKey={(person) => person.id}
        getPrimaryText={(person) => person.id}
        onDelete={(index) => handleDeletePerson(index)}
        emptyMessage="No players added."
      />

      <Button
        variant="contained"
        fullWidth
        disabled={people.length < config.numberOfTeams}
        onClick={handleSendTournament}
      >
        Save Tournament
      </Button>
    </>
  );
};

export const AddPersonPage: FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const data = location.state as TournamentSetupData | null;

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const formData = new FormData(event.currentTarget);

    const person = formData.get("person") as string;

    if (!person) return;

    const newPerson: Person = { id: person };
    peopleStore.addPerson(newPerson);

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
      <TextField label="Player" name="person" variant="outlined" required />
      <Button type="submit" size="large" fullWidth variant="outlined">
        Add Player
      </Button>
    </Box>
  );
};
