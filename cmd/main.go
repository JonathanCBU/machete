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
func renderProfile(p aws.Profile) (title, detail string) {
	title = p.Name
	detail = p.ToString()
	return
}

func main() {
	// File logging catches anything that happens before app.Run() starts
	// (config/AWS load failures) or if the app exits early — the in-app
	// LogWriter can't show those since it needs the event loop running.
	logFile, err := os.OpenFile("/tmp/tview-app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ltime)

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("failed to read config", err)
	}
	log.Printf("aws config dir: %q", conf.Aws)

	awsConf, err := aws.GetAwsConfig(conf.Aws)
	if err != nil {
		log.Fatal("failed to read aws config", err)
	}
	log.Printf("loaded %d profiles", len(awsConf.Profiles))

	app := tui.NewApp()

	home := tview.NewTextView().
		SetText("Welcome home.\n\nPress a hotkey on the left to switch views.")
	home.SetBorder(true).SetTitle(" Home ")
	app.Register(tui.NewStaticScreen(app, "Home", 'H', home))

	profiles := tui.NewCardScreen(
		app, "Profiles", 'P',
		awsConf.Profiles, aws.RenderProfile,
	)
	app.Register(profiles)

	logs := tview.NewTextView().SetDynamicColors(true)
	logs.SetBorder(true).SetTitle(" Logs ")
	app.Register(tui.NewStaticScreen(app, "Logs", 'L', logs))

	// Once the Logs screen exists, also stream log output into it live,
	// on top of (not instead of) the file — see log_writer.go.
	log.SetOutput(io.MultiWriter(logFile, tui.NewLogWriter(app.Application(), logs)))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
