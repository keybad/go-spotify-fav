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

const (
	pageBefore = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>YourLibrary to html converter</title>
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
<link href="https://cdn.datatables.net/1.11.5/css/dataTables.bootstrap5.min.css" rel="stylesheet"/>
<link href="https://cdn.datatables.net/select/1.3.4/css/select.dataTables.min.css" rel="stylesheet"/>
<link href="https://cdn.datatables.net/buttons/2.2.2/css/buttons.bootstrap5.min.css" rel="stylesheet"/>
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
<script src="https://cdn.datatables.net/1.11.5/js/jquery.dataTables.min.js"></script>
<script src="https://cdn.datatables.net/1.11.5/js/dataTables.bootstrap5.min.js"></script>
<script src="https://cdn.datatables.net/select/1.3.4/js/dataTables.select.min.js"></script>
<script src="https://cdn.datatables.net/buttons/2.2.2/js/dataTables.buttons.min.js"></script>
<script src="https://cdn.datatables.net/buttons/2.2.2/js/buttons.bootstrap5.min.js"></script>
<script>
$(document).ready(function(){ 
	$("table").DataTable({
		paging: false,
		select: true,
		dom: 'Bfrtip',
        buttons: [
            'copy', 'csv', 'excel', 'pdf', 'print'
        ],
	});
});
</script>
</head>
<body>
<div class="container">
`
	pageAfter = `</div>
</body>
</html>
`
	tableBefore = `<div class="table-responsive">
<table class="table" data-order='[[0,"asc"]]'>
<thead>
<tr>
<th scope="col">#</th>
<th scope="col">Artist</th>
<th scope="col">Track</th>
<th scope="col">Album</th>
</tr>
</thead>
<tbody>
`
	tableAfter = `</tbody>
</table>
</div>
`
	tableTracksRow = `<tr><th scope="col">%d</th><td>%s</td><td>%s</td><td>%s</td></tr>
`
	titleRow = `<h2>%s</h2>
`
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

	_, _ = w.WriteString(pageBefore)

	if len(library.Tracks) > 0 {
		_, _ = w.WriteString(fmt.Sprintf(titleRow, "Tracks"))

		_, _ = w.WriteString(tableBefore)

		for k, v := range library.Tracks {
			_, _ = w.WriteString(fmt.Sprintf(tableTracksRow, k+1, v.Artist, v.Track, v.Album))
		}

		_, _ = w.WriteString(tableAfter)
	}

	_, _ = w.WriteString(pageAfter)

	log.Println("Done!")
}
