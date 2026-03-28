package tournament

type Gender int

const (
	Male Gender = iota
	Female
	Else
)

func GetAllGenders() []Gender {
	return []Gender{Male, Female, Else}
}

func GenderFromString(g string) Gender {
	switch g {
	case "Male", "M":
		return Male
	case "Female", "F":
		return Female
	default:
		return Else
	}
}

type Person struct {
	Id string `json:"id"`
}

type Team struct {
	Person_1   Person `json:"person1"`
	Person_2   Person `json:"person2"`
	TeamGender Gender `json:"gender"`
}

func NewTeam(person1, person2 Person, teamGender Gender) *Team {
	return &Team{
		Person_1:   person1,
		Person_2:   person2,
		TeamGender: teamGender,
	}
}

func MakeTeam(person1, person2 Person, teamGender Gender) Team {
	return Team{
		Person_1:   person1,
		Person_2:   person2,
		TeamGender: teamGender,
	}
}

func GenderCount(teams []Team, gender Gender) int {
	count := 0
	for _, t := range teams {
		if t.TeamGender == gender {
			count++
		}
	}
	return count
}

func GetTeamsByGender(teams []Team, gender Gender) []Team {
	var filteredTeams []Team
	for _, t := range teams {
		if t.TeamGender == gender {
			filteredTeams = append(filteredTeams, t)
		}
	}
	return filteredTeams
}
