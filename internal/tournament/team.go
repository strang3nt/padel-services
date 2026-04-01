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

func (p Person) IsNil() bool {
	return p.Id == ""
}

type Team struct {
	Person1    Person `json:"person1"`
	Person2    Person `json:"person2"`
	TeamGender Gender `json:"gender"`
}

func NewTeam(person1, person2 Person, teamGender Gender) *Team {
	return &Team{
		Person1:    person1,
		Person2:    person2,
		TeamGender: teamGender,
	}
}

func MakeTeam(person1, person2 Person, teamGender Gender) Team {
	return Team{
		Person1:    person1,
		Person2:    person2,
		TeamGender: teamGender,
	}
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

func GetPeople(teams []Team) []Person {
	res := make([]Person, 0)

	for _, t := range teams {
		ps := []Person{t.Person1, t.Person2}

		for _, p := range ps {
			if !p.IsNil() {
				res = append(res, p)
			}
		}
	}
	return res
}
