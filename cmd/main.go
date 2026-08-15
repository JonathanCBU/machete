package main

import (
	"fmt"
	"log"

	"github.com/JonathanCBU/machete/internal/tui"
	"github.com/rivo/tview"
)

// Profile is domain data — tui knows nothing about this type. It only
// needs to satisfy whatever CardRenderer[Profile] below extracts from it.
type Profile struct {
	Name     string
	Role     string
	Region   string
	Status   string
	LastSeen string
}

// fetchProfiles stands in for wherever your real data comes from — an API
// call, a DB query, a file watch. Because CardScreen.SetItems can be
// called any time, this could just as easily run on a timer or in
// response to a keypress to refresh the screen live.
func fetchProfiles() []Profile {
	return []Profile{
		{"Alice Chen", "Backend", "us-east-1", "Active", "2m ago"},
		{"Bob Martins", "Frontend", "eu-west-1", "Active", "14m ago"},
		{"Carla Diaz", "SRE", "ap-south-1", "On call", "just now"},
	}
}

func renderProfile(p Profile) (title, summary, detail string) {
	title = p.Name
	summary = p.Role + " • " + p.Region
	detail = fmt.Sprintf(
		"\nRole: %s\nRegion: %s\nStatus: %s\nLast seen: %s",
		p.Role, p.Region, p.Status, p.LastSeen,
	)
	return
}

func main() {
	app := tui.NewApp()

	home := tview.NewTextView().
		SetText("Welcome home.\n\nPress a hotkey on the left to switch views.")
	home.SetBorder(true).SetTitle(" Home ")
	app.Register(tui.NewStaticScreen(app.Application(), "Home", 'h', home))

	profiles := tui.NewCardScreen(
		app.Application(), app.Pages(),
		"Profiles", 'p',
		fetchProfiles(), renderProfile,
	)
	app.Register(profiles)

	logs := tview.NewTextView().
		SetText("Logs view.\n\nTail output would go here.").
		SetDynamicColors(true)
	logs.SetBorder(true).SetTitle(" Logs ")
	app.Register(tui.NewStaticScreen(app.Application(), "Logs", 'l', logs))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
