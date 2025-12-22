package services

import (
	"os"
	"testing"
)

func TestCreatePDFTournament_VerifyCreation(t *testing.T) {

	data := TemplateData{
		Tournament: TournamentData{
			Data:   "",
			Rounds: []Round{},
		},
	}

	pdfPath, err := CreatePdfTournament(data, "../../template/template_schedule.html")
	if err != nil {
		t.Fatalf("CreatePdfTournament returned an error: %v", err)
	}
	t.Run("Assertion_1_PDFCreation", func(t *testing.T) {
		if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
			t.Errorf("Expected PDF file to be created at %s, but it does not exist", pdfPath)
		}
	})

	// Clean up the generated PDF file after test
	defer os.Remove(pdfPath)
}
