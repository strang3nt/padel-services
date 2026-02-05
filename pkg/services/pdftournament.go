package services

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/strang3nt/padel-services/pkg/tournament"
)

type Match struct {
	Court       string
	TeamA       string
	ScoreA      string
	TeamB       string
	ScoreB      string
	RoundNumber int
}

type Round struct {
	RoundNumber int
	Matches     []Match
}

type TournamentData struct {
	Name      string
	StartDate string
	Rounds    []Round
}

type TemplateData struct {
	Tournament TournamentData
}

type TournamentType int

const (
	Rodeo TournamentType = iota
)

var tournamentType = map[TournamentType]string{
	Rodeo: "Rodeo",
}

func (tt TournamentType) String() string {
	return tournamentType[tt]
}

type TournamentPdfGenerator struct {
	chromeExecutable string
	templatesDirs    map[TournamentType]*template.Template
}

//go:embed templates/*
var templates embed.FS
var style, _ = templates.ReadFile("templates/style.css")
var templateRodeoSchedule, _ = template.ParseFS(templates, "templates/template_rodeo_schedule.html")

func MakeTournamentPdfGenerator() TournamentPdfGenerator {
	var chromeExecutable string = "google-chrome"
	if value, ok := os.LookupEnv("CHROME_EXECUTABLE"); ok {
		chromeExecutable = value
	}

	return TournamentPdfGenerator{
		chromeExecutable: chromeExecutable,
		templatesDirs: map[TournamentType]*template.Template{
			Rodeo: templateRodeoSchedule,
		},
	}

}

func (t TournamentPdfGenerator) CreatePdfTournament(
	data TemplateData,
	tt TournamentType,
	outputFileName string,
) (string, error) {

	tempHTMLFile, err := t.runTemplate(data, tt)
	if err != nil {
		return "", fmt.Errorf("error executing and saving template: %v", err)
	}
	defer func() {
		err := os.Remove(tempHTMLFile)
		if err != nil {
			log.Printf("error while removing temp file: %v", err)
		}
	}()

	outputFile := fmt.Sprint(outputFileName, ".pdf")

	if err := t.generatePDFWithHeadlessChrome(tempHTMLFile, outputFile); err != nil {
		return "", fmt.Errorf("error generating PDF: %v", err)
	}

	return outputFile, nil
}

func (tt TournamentPdfGenerator) runTemplate(
	data TemplateData, tournamentType TournamentType) (string, error) {

	tournamentTemplate, ok := tt.templatesDirs[tournamentType]
	if !ok {
		return "", fmt.Errorf("unexpected tournament type: %s", tournamentType)
	}

	templateData := map[string]any{
		"Tournament": data.Tournament,
		"Style":      string(style),
	}

	var buf bytes.Buffer
	if err := tournamentTemplate.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	tempFile, _ := os.CreateTemp("./", "schedule-*.html")

	defer func() {
		err := tempFile.Close()
		if err != nil {
			log.Printf("error while closing temp file: %v", err)
		}
	}()

	if _, err := tempFile.Write(buf.Bytes()); err != nil {
		return "", fmt.Errorf("writing to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

func (tt TournamentPdfGenerator) generatePDFWithHeadlessChrome(inputHTMLPath, outputPath string) error {

	absPath, _ := filepath.Abs(inputHTMLPath)
	inputURL := "file://" + absPath

	args := []string{
		"--headless=new",
		fmt.Sprintf("--print-to-pdf=%s", outputPath),
		"--print-to-pdf-no-header",
		"--landscape",
		"--no-margins",
		"--no-sandbox",
		"--disable-gpu",
		"--font-render-hinting=none",
		"--virtual-time-budget=1000",
		"--default-page-settings={\"pageSize\":\"A4\"}",
		inputURL,
	}

	cmd := exec.Command(tt.chromeExecutable, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("chrome executable failed: %w. Stderr: %s", err, stderr.String())
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF output file was not created. Chrome executable error log: %s", stderr.String())
	}

	return nil
}

func FromTournamentToTemplateData(tournament tournament.Tournament) TemplateData {
	return TemplateData{
		Tournament: TournamentData{
			Name:      tournament.GetName(),
			StartDate: tournament.GetDateStart().Format("2006-01-02"),
			Rounds: func() []Round {
				var rounds []Round
				for roundIndex, round := range tournament.GetRounds() {
					var matches []Match
					for _, match := range round.Matches {

						surnamePerson1TeamA := strings.Split(match.TeamA.Person_1.Id, " ")
						surnamePerson2TeamA := strings.Split(match.TeamA.Person_2.Id, " ")
						surnamePerson1TeamB := strings.Split(match.TeamB.Person_1.Id, " ")
						surnamePerson2TeamB := strings.Split(match.TeamB.Person_2.Id, " ")

						matches = append(matches, Match{
							Court:       strconv.Itoa(match.CourtId),
							TeamA:       surnamePerson1TeamA[len(surnamePerson1TeamA)-1] + ", " + surnamePerson2TeamA[len(surnamePerson2TeamA)-1],
							ScoreA:      "",
							TeamB:       surnamePerson1TeamB[len(surnamePerson1TeamB)-1] + ", " + surnamePerson2TeamB[len(surnamePerson2TeamB)-1],
							ScoreB:      "",
							RoundNumber: roundIndex + 1,
						})
					}
					rounds = append(rounds, Round{
						RoundNumber: roundIndex + 1,
						Matches:     matches,
					})
				}
				return rounds
			}(),
		},
	}
}

func FromTournamentDataToTemplateData(tournament tournament.TournamentData) TemplateData {
	return TemplateData{
		Tournament: TournamentData{
			Name:      tournament.Name,
			StartDate: tournament.Date.Format("2006-01-02"),
			Rounds: func() []Round {
				var rounds []Round
				for roundIndex, round := range tournament.Rounds {
					var matches []Match
					for _, match := range round.Matches {

						surnamePerson1TeamA := strings.Split(match.TeamA.Person_1.Id, " ")
						surnamePerson2TeamA := strings.Split(match.TeamA.Person_2.Id, " ")
						surnamePerson1TeamB := strings.Split(match.TeamB.Person_1.Id, " ")
						surnamePerson2TeamB := strings.Split(match.TeamB.Person_2.Id, " ")

						matches = append(matches, Match{
							Court:       strconv.Itoa(match.CourtId),
							TeamA:       surnamePerson1TeamA[len(surnamePerson1TeamA)-1] + ", " + surnamePerson2TeamA[len(surnamePerson2TeamA)-1],
							ScoreA:      "",
							TeamB:       surnamePerson1TeamB[len(surnamePerson1TeamB)-1] + ", " + surnamePerson2TeamB[len(surnamePerson2TeamB)-1],
							ScoreB:      "",
							RoundNumber: roundIndex + 1,
						})
					}
					rounds = append(rounds, Round{
						RoundNumber: roundIndex + 1,
						Matches:     matches,
					})
				}
				return rounds
			}(),
		},
	}
}
