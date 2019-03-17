package main

import (
	"os"

	"github.com/grunions/steamcmd"
)

// AppHandler is the type of handler required to fetch the app, or gather information
// about it.
type AppHandler int

// Types of supported app handlers
const (
	Unknown AppHandler = iota
	Download
	Steam
	Minecraft
	Factorio
	Mumble
	Teamspeak
)

var games = []struct {
	Name     string
	Handler  AppHandler
	AppID    int
	Launcher string
}{
	{
		// https://community.bistudio.com/wiki/Arma_3_Dedicated_Server#Configuration
		Name:     "Arma 3",
		Handler:  Steam,
		AppID:    233780,
		Launcher: "./arma3server",
	},
	{
		// https://developer.valvesoftware.com/wiki/Rust_Dedicated_Server#Configuration_.26_running
		Name:     "Rust",
		Handler:  Steam,
		AppID:    258550,
		Launcher: "./RustDedicated",
	},
	{
		// https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Dedicated_Servers#Linux_Scripts
		Name:     "Counter-Strike: Global Offensive",
		Handler:  Steam,
		AppID:    740,
		Launcher: "./srcds_run -game csgo",
	},
	{
		// https://dontstarve.fandom.com/wiki/Guides/Don%E2%80%99t_Starve_Together_Dedicated_Servers
		// To use an alternative configuration directory,
		// `-conf_dir DoNotStarveServerDirectory`.
		// changes cfgdir to: "~/.klei/DoNotStarveServerDirectory"
		Name:     "Dont Starve Together",
		Handler:  Steam,
		AppID:    343050,
		Launcher: "./dontstarve_dedicated_server_nullrenderer",
	},
	{
		// https://wiki.teamfortress.com/wiki/Linux_dedicated_server
		Name:     "Team Fortress 2",
		Handler:  Steam,
		AppID:    232250,
		Launcher: "./srcds_run -game tf",
	},
	{
		// https://wiki.garrysmod.com/page/Hosting_A_Dedicated_Server
		Name:     "Garry's Mod",
		Handler:  Steam,
		AppID:    4020,
		Launcher: "./srcds_run -game garrysmod",
	},
	{
		Name:     "Counter-Strike: Source",
		Handler:  Steam,
		AppID:    232330,
		Launcher: "./srcds_run -game cstrike",
	},
	{
		// https://steamcommunity.com/sharedfiles/filedetails/?id=764005219
		Name:     "Left 4 Dead 2",
		Handler:  Steam,
		AppID:    222860,
		Launcher: "./srcds_run -game left4dead2",
	},
	{
		Name:     "Left 4 Dead",
		Handler:  Steam,
		AppID:    222840,
		Launcher: "./srcds_run -game left4dead",
	},
}

func main() {
	steam := steamcmd.New("", "", "/tmp/Steam")
	defer os.RemoveAll(steam.AppBasePath)

	steam.Debug = true

	if err := steam.EnsureInstalled(); err != nil {
		panic(err)
	}

	const dstServerAppID = 343050
	if err := steam.InstallUpdateApp(dstServerAppID); err != nil {
		panic(err)
	}

	// cleanup
}
