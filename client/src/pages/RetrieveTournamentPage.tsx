import { useState, type FC } from 'react';
import { Input, Button, Snackbar } from '@telegram-apps/telegram-ui';
import { Page } from '@/components/Page.tsx';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';

export interface Person {
  id: string
}

export interface Team {
  person1: Person
  person2: Person
  gender: number
}

export interface Match {

  teamA: Team
  teamB: Team
  matchStatus: number
  courtId: number
}

export interface Matches {
  matches: Match[]
}

export interface TournamentData {
  name: string;
  date: string;
  teams: Team[];
  rounds: Matches[];
}

export interface Tournaments {
  date: string
  tournaments: TournamentData[]
}

interface FormData {
  date: string;
}

export const RetrieveTournamentPage: FC = () => {
  const { bearerToken } = useAuth()
  const [formData, setFormData] = useState<FormData>({ date: '' });
  const navigate = useNavigate();

  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {

      const response = await fetch(
        `/api/tournaments?date=${formData.date}`,
        {
          headers: { Authorization: `Bearer ${bearerToken}` }
        })

      const result = await response.json() as Tournaments
      if (response.ok) {
        navigate('/available-tournaments', { state: result });
      } else {
        return <Snackbar
        duration={3}
        onClose={() => {}}>Server not reachable</Snackbar>
      }
    } catch (error) {
      console.error("Submission error:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (

    <Page>
      <form onSubmit={handleSubmit}>
        <Input header="Insert date"
          type="date"
          required
          onChange={(e) => setFormData({ date: e.target.value })} />
        <Button
          type="submit"
          disabled={isLoading}>
          {isLoading ? 'Submitting...' : 'Submit'}
        </Button>
      </form>
    </Page>
  );

};
