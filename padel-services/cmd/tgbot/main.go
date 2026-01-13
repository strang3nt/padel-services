package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"padelservices/pkg/services"
	"padelservices/pkg/tournament"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"padelservices/pkg/database"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"
)

var defaultMessage string

var ctx = context.Background()

func main() {
	conn, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = database.CreateDatabaseTables(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while initializing database: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	telegramBotToken := os.Getenv("BOT_TOKEN")

	state := make(map[int64]StateMachine, 0)

	opts := []bot.Option{
		bot.WithDebug(),
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler("/crea_torneo", bot.MatchTypeExact, createTournamentHandler(&state)),
		bot.WithMessageTextHandler("/cerca_torneo", bot.MatchTypeExact, searchTournament(conn)),
	}

	b, err := bot.New(telegramBotToken, opts...)

	// Register this in your main function after creating the bot instance
	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		var chatId int64
		if update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil {
			chatId = update.CallbackQuery.Message.Message.Chat.ID
		} else if update.Message != nil {
			chatId = update.Message.Chat.ID
		} else {
			return false
		}
		if s, ok := state[chatId]; ok && s == TournamentCreated {
			return true
		}
		return false
	}, printToPdf(&state, conn))

	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func searchTournament(conn *pgxpool.Pool) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		kb := datepicker.New(b, selectTournament(conn))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Inserisci la data del torneo da cercare",
			ReplyMarkup: kb,
		})

		if err != nil {
			log.Printf("error while sending message: %v", err)
			return
		}
	}
}

func selectTournament(conn *pgxpool.Pool) datepicker.OnSelectHandler {

	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, date time.Time) {

		fmt.Printf("selectTournament %v\n", date)

		tournaments, err := database.GetTournamentsByDate(ctx, conn, date)

		if err != nil {
			log.Printf("error while trying to retrieve tournaments: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   "Non è possibile trovare un torneo per via di un errore interno :(",
			})
			if err != nil {
				log.Printf("error while sending message: %v", err)
			}
			return
		}

		kb := inline.New(b)

		for i, t := range tournaments {
			kb = kb.Button(fmt.Sprintf("%v", i), []byte(""), selectedTournament(t))
		}

		var msg string

		if len(tournaments) == 0 {
			msg = fmt.Sprintf(`In data %v non ho trovato tornei`, date.Format("2006-01-02"))
		} else {
			msg = fmt.Sprintf(`In data %v ho trovato i seguenti tornei:\n`, date.Format("2006-01-02"))
		}
		for i, t := range tournaments {
			txt := fmt.Sprintf("\nTorneo %v, con le seguenti squadre", i)
			for _, tms := range t.Teams {
				txt = txt + fmt.Sprintf("%v, %v\n", tms.Person_1.Id, tms.Person_2.Id)
			}
			msg += txt
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      mes.Message.Chat.ID,
			Text:        msg,
			ReplyMarkup: kb})

		if err != nil {
			log.Printf("error while sending message: %v", err)
			return
		}
	}
}

func selectedTournament(t tournament.TournamentData) inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		printPdf(ctx, services.FromTournamentDataToTemplateData(t), b, mes.Message.Chat.ID)
	}
}

func createTournamentHandler(state *map[int64]StateMachine) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		(*state)[update.Message.Chat.ID] = TournamentCreated
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: `Ho bisogno dei dettagli del torneo, che mi devi inviare in un messaggio,
rispettando il seguente formato:

<campi disponibili per il torneo>
<round del torneo>
<team 1 partecipante 1>, <team 1 partecipante 2>
<team 2 partecipante 1>, <team 2 partecipante 2>,
...

Ad esempio:

6
5
Marco Rossi, Luigi Blu,
Tizio Caio, Sempro Nio
`,
		})

		if err != nil {
			log.Printf("error while sending message: %v", err)
			return
		}
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      defaultMessage,
		ParseMode: models.ParseModeMarkdown,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: bot.True(),
		},
	})

	if err != nil {
		log.Printf("error while sending message: %v", err)
		return
	}
}

func printToPdf(state *map[int64]StateMachine, conn *pgxpool.Pool) bot.HandlerFunc {

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		tournamentData := update.Message.Text

		msgScanner := bufio.NewScanner(strings.NewReader(tournamentData))

		msgScanner.Scan()
		availableCourts, err := strconv.Atoi(msgScanner.Text())
		if err != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Errore: la prima riga deve contenere il numero di campi disponibili per il torneo.",
			})

			if err != nil {
				log.Printf("error while sending message: %v", err)
				return
			}

			return
		}

		msgScanner.Scan()
		roundsNumber, err := strconv.Atoi(msgScanner.Text())
		if err != nil {
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Errore: la seconda riga deve contenere il numero di round da disputare durante il torneo.",
			})
			if err != nil {
				log.Printf("error while sending message: %v", err)
				return
			}
			return
		}

		teams, err := services.MakeTeamsFromMessage(msgScanner)
		if err != nil {
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Ho incontrato un errore leggendo le squadre: %v.", err),
			})
			if err != nil {
				log.Printf("error while sending message: %v", err)

			}
			return
		}
		log.Print("creating tournament...")
		tournament := services.CreateTournament("Rodeo", time.Now(), teams, roundsNumber, availableCourts)
		log.Print("tournament created")
		template_data := services.FromTournamentToTemplateData(tournament)
		err = database.CreateTournament(ctx, conn, &tournament)
		if err != nil {
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Errore: non è stato possibile salvare il torneo nel database.",
			})
			if err != nil {
				log.Printf("error while sending message: %v", err)

			}
			return
		}
		printPdf(ctx, template_data, b, update.Message.Chat.ID)
		(*state)[update.Message.Chat.ID] = Started
	}
}

func printPdf(ctx context.Context, template_data services.TemplateData, b *bot.Bot, chatId int64) {
	pdfPath, err := services.CreatePdfTournament(template_data, "template/template_schedule.html")
	if err != nil {
		log.Printf("Error creating PDF: %v", err)
		return
	}
	pdfFile, err := os.Open(pdfPath)
	if err != nil {
		log.Printf("Error opening generated PDF: %v", err)
		return
	}

	defer func() {
		err := os.Remove(pdfFile.Name())
		if err != nil {
			log.Printf("error while removing file: %v", err)
		}
	}()

	_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   chatId,
		Document: &models.InputFileUpload{Filename: "tournament_schedule.pdf", Data: bufio.NewReader(pdfFile)},
		Caption:  fmt.Sprintf("Document '%s'.", "tournament_schedule.pdf"),
	})

	if err != nil {
		log.Printf("error while sending document: %v", err)
		return
	}

}
