package tournament

type Gender int

const (
	Male Gender = iota
	Female
	Else
)

func GenderFromString(g string) Gender {
	switch g {
	case "Male":
		return Male
	case "Female":
		return Female
	default:
		return Else
	}
}

type Person struct {
	Id string
}

type Team struct {
	Person_1   Person
	Person_2   Person
	TeamGender Gender
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
