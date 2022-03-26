package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type (
	YourLibrary struct {
		Tracks       []TrackItem        `json:"tracks"`
		BannedTracks []BannedTrackItem  `json:"bannedTracks"`
		Albums       []AlbumItem        `json:"albums"`
		Artists      []ArtistItem       `json:"artists"`
		BannedArtist []BannedArtistItem `json:"bannedArtists"`
		Shows        []ShowItem         `json:"shows"`
		Episodes     []EpisodeItem      `json:"episodes"`
		Other        []OtherItem        `json:"other"`
	}

	TrackItem struct {
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Track  string `json:"track"`
		Uri    string `json:"uri"`
	}

	BannedTrackItem struct {
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Track  string `json:"track"`
		Uri    string `json:"uri"`
	}

	AlbumItem struct {
		Artist string `json:"artist"`
		Album  string `json:"album"`
		Uri    string `json:"uri"`
	}

	ArtistItem struct {
		Name string `json:"name"`
		Uri  string `json:"uri"`
	}

	BannedArtistItem struct {
		Name string `json:"name"`
		Uri  string `json:"uri"`
	}

	ShowItem    struct{}
	EpisodeItem struct{}
	OtherItem   struct{}

	Exit struct {
		Code    int
		Message string
	}
)

var (
	fs   *flag.FlagSet
	file string
)

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func handleExit() {
	if r := recover(); r != nil {
		if exit, ok := r.(Exit); ok {
			fmt.Print(exit.Message)

			fs.Usage()
			os.Exit(exit.Code)
		}
		panic(r)
	}
}

func main() {
	defer handleExit()

	args := os.Args

	fs = flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)
	fs.StringVar(&file, "f", "YourLibrary.json", "spotify library file (required)")

	if len(args) < 2 {
		args = append(args, "-f=YourLibrary.json")
	}

	if args[1] == "-h" {
		panic(Exit{Code: 1})
	}

	err := fs.Parse(args[1:])
	if err != nil {
		panic(Exit{Code: 1})
	}

	if fs.Parsed() {
		if file, err = filepath.Abs(file); err != nil {
			panic(Exit{Code: 1})
		}

		if _, err := os.Stat(file); err != nil {
			panic(Exit{
				Code:    1,
				Message: "Spotify `YourLibrary.json` file required. Please supply the file\n\n",
			})
		}
	}

	r, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	var library YourLibrary
	if err := json.NewDecoder(r).Decode(&library); err != nil {
		log.Fatal(err)
	}

	w, err := os.Create(fileNameWithoutExt(file) + ".html")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = w.WriteString(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>YourLibrary to html converter</title>
</head>
<body>`)

	//TODO

	_, _ = w.WriteString(`
</body>
</html>
`)

	log.Println("Done!")
}
