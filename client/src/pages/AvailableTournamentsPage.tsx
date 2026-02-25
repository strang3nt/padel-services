import { type FC } from 'react';
import { TournamentData, Tournaments } from './RetrieveTournamentPage';
import { Page } from '@/components/Page.tsx';
import { useLocation } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';
import { postEvent } from '@tma.js/sdk-react';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import Section from '@/components/Section';
import ListItemText from '@mui/material/ListItemText';
import ListItem from '@mui/material/ListItem';
export const AvailableTournamentsPage: FC = () => {
  const { bearerToken } = useAuth()
  const location = useLocation();

  const downloadTournament = (tournamentData: TournamentData) => {

    fetch(
      `/api/tournament/generate-link`,
      {
        method: "POST",
        headers: {
          'Authorization': `Bearer ${bearerToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(tournamentData),
      }
    ).then((response) => {
      response.json().then(({ user, token }) => {
        const url = `${window.location.origin}/api/tournament/download?user=${user}&token=${token}`
        if (import.meta.env.DEV) {
          window.location.href = url;
        } else {
          postEvent('web_app_request_file_download', {
            url: url,
            file_name: `${tournamentData.date}_${tournamentData.name}.pdf`
          })
        }
      }
      )
    });
  }

  const { date, tournaments } = location.state as Tournaments;
  return (
    <Page>
      <Section
        title={`Tournaments at date ${date}`}
      >

        <List>
          {
            tournaments == null || tournaments.length == 0 ?
              <ListItem>
                <ListItemText primary={"No tournaments found at selected date"} />
              </ListItem> :
              (
                tournaments.map(
                  (tournamentData) =>

                    <ListItemButton
                      onClick={() => downloadTournament(tournamentData)}
                    >
                      <ListItemText primary={tournamentData.name} secondary={`Participants: ${tournamentData.teams.length}`} />
                    </ListItemButton>
                )
              )}
        </List>

      </Section>
    </Page >
  );
};
