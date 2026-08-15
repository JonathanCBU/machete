package main

import (
	"io"
	"log"
	"os"

	"github.com/JonathanCBU/machete/internal/aws"
	"github.com/JonathanCBU/machete/internal/config"
	"github.com/JonathanCBU/machete/internal/tui"
	"github.com/rivo/tview"
)

// renderProfile is now a CardRenderer[aws.Profile] — T is inferred from
// the slice passed to NewCardScreen below, so nothing in tui needs to
// change to accept a different concrete type here.
func renderProfile(p aws.Profile) (title, summary, detail string) {
	title = p.Name
	summary = p.Attributes["region"]
	detail = p.ToString()
	return
}

func main() {
	logFile, _ := os.OpenFile("/tmp/tview-app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	defer logFile.Close()

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("failed to read config", err)
	}

	awsConf, err := aws.GetAwsConfig(conf.Aws)
	if err != nil {
		log.Fatal("failed to read aws config", err)
	}

	app := tui.NewApp()

	home := tview.NewTextView().
		SetText("Welcome home.\n\nPress a hotkey on the left to switch views.")
	home.SetBorder(true).SetTitle(" Home ")
	app.Register(tui.NewStaticScreen(app.Application(), "Home", 'h', home))

	profiles := tui.NewCardScreen(
		app.Application(), app.Pages(),
		"Profiles", 'p',
		awsConf.Profiles, renderProfile,
	)
	app.Register(profiles)

	logs := tview.NewTextView().SetDynamicColors(true)
	log.SetOutput(io.MultiWriter(logFile, tui.NewLogWriter(app.Application(), logs)))

	logs.SetBorder(true).SetTitle(" Logs ")
	app.Register(tui.NewStaticScreen(app.Application(), "Logs", 'l', logs))

	log.SetOutput(tui.NewLogWriter(app.Application(), logs))
	log.SetFlags(log.Ltime)

	log.Printf("loaded %d profiles", len(awsConf.Profiles))
	logs.SetBorder(true).SetTitle(" Logs ")
	app.Register(tui.NewStaticScreen(app.Application(), "Logs", 'l', logs))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
