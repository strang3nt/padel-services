import { type FC } from 'react';
import { List } from '@telegram-apps/telegram-ui';
import { TournamentData, Tournaments } from './RetrieveTournamentPage';
import { Page } from '@/components/Page.tsx';
import { useLocation } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';
import { isTMA } from '@tma.js/sdk-react';
import { DisplayData } from '@/components/DisplayData/DisplayData';

export const AvailableTournamentsPage: FC = () => {
  const { bearerToken } = useAuth()
  const location = useLocation();

  const onClick = (tournamentData: TournamentData) => {

    // This triggers the browser's native download manager
    const apiRequest = fetch(
      `/api/tournament/generate-link`,
      {
        method: "POST",
        headers: {
          'Authorization': `Bearer ${bearerToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(tournamentData),
      }
    )
    Promise.all([isTMA('complete'), apiRequest]).then(([_, response]) => {
      response.json().then(({ token }) => {
        const url = `${window.location.origin}/api/tournament/download?token=${token}`
        window.location.href = url;

      }
      )

    });


  }

  // Cast the state to our interface. 
  // We use optional chaining because state could be 'null' if someone types the URL directly.
  const { date, tournaments } = location.state as Tournaments;
  return (
    <Page>
      <List>
        <DisplayData
          header={`Tournaments at ${date}`}
          rows={
            tournaments
              .map((tournamentData) => (
                {
                  title: tournamentData.name,
                  value: <span onClick={() => onClick(tournamentData)}> 
                    Participants: {tournamentData.teams.length}
                  </span>
                }
              ))
          }
        />
      </List>
    </Page>
  );
};
