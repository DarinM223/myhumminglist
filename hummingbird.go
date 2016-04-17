package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	HummingbirdAddURL     = "https://hummingbird.me/api/v1/libraries/%d"
	HummingbirdEditURL    = "https://hummingbird.me/api/v1/libraries/%d"
	HummingbirdDeleteURL  = "https://hummingbird.me/api/v1/libraries/%d/remove"
	HummingbirdLibraryURL = "https://hummingbird.me/api/v1/users/%s/library?include_mal_id=true"
)

// HummingbirdAnime represents the JSON data of a Hummingbird library entry
type HummingbirdAnime struct {
	NumEpisodesWatched int                  `json:"episodes_watched"`
	AnimeStatus        string               `json:"status"`
	NumRewatchedTimes  int                  `json:"rewatched_times"`
	IsRewatching       bool                 `json:"rewatching"`
	Data               HummingbirdAnimeData `json:"anime"`
}

// HummingbirdAnimeData represents the JSON data of a Hummingbird anime
type HummingbirdAnimeData struct {
	Id            int    `json:"id"`
	MalID         int    `json:"mal_id"`
	Title         string `json:"title"`
	EpisodeCount  int    `json:"episode_count"`
	EpisodeLength int    `json:"episode_length"`
}

func (ha HummingbirdAnime) ID() AnimeID {
	return AnimeID{
		Hummingbird: ha.Data.Id,
		MyAnimeList: ha.Data.MalID,
	}
}

func (ha HummingbirdAnime) Title() string {
	return ha.Data.Title
}

func (ha HummingbirdAnime) Status() int {
	switch ha.AnimeStatus {
	case "currently-watching":
		return StatusWatching
	case "plan-to-watch":
		return StatusPlanToWatch
	case "completed":
		return StatusCompleted
	case "on-hold":
		return StatusOnHold
	case "dropped":
		return StatusDropped
	default:
		panic("Invalid status")
	}
}

func (ha HummingbirdAnime) EpisodesWatched() int {
	return ha.NumEpisodesWatched
}

func (ha HummingbirdAnime) RewatchedTimes() int {
	return ha.NumRewatchedTimes
}

func (ha HummingbirdAnime) Rewatching() bool {
	return ha.IsRewatching
}

func StatusToHummingbirdString(status int) string {
	switch status {
	case StatusWatching:
		return "currently-watching"
	case StatusPlanToWatch:
		return "plan-to-watch"
	case StatusCompleted:
		return "completed"
	case StatusOnHold:
		return "on-hold"
	case StatusDropped:
		return "dropped"
	default:
		panic("Invalid status")
	}
}

func AnimeToHummingbird(anime Anime) HummingbirdAnime {
	return HummingbirdAnime{
		NumEpisodesWatched: anime.EpisodesWatched(),
		NumRewatchedTimes:  anime.RewatchedTimes(),
		AnimeStatus:        StatusToHummingbirdString(anime.Status()),
		Data: HummingbirdAnimeData{
			Id:    anime.ID().Get(Hummingbird),
			MalID: anime.ID().Get(MyAnimeList),
			Title: anime.Title(),
		},
	}
}

type HummingbirdAnimeList struct {
	username    string
	anime       map[int]HummingbirdAnime
	changes     []Change
	pastChanges []Change
	authToken   string
}

func NewHummingbirdAnimeList(username string, authToken string) *HummingbirdAnimeList {
	return &HummingbirdAnimeList{
		username:    username,
		authToken:   authToken,
		anime:       make(map[int]HummingbirdAnime),
		changes:     []Change{},
		pastChanges: []Change{},
	}
}

func (hal HummingbirdAnimeList) Type() int {
	return Hummingbird
}

func (hal HummingbirdAnimeList) AuthToken() string {
	return hal.authToken
}

// Fetch fetches the animelist from the api and adds the changes to the change lists
func (hal *HummingbirdAnimeList) Fetch() error {
	// TODO(DarinM223): http get request is extremely slow, maybe put inside goroutine?
	resp, err := http.Get(fmt.Sprintf(HummingbirdLibraryURL, hal.username))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Status code for response is not 200")
	}

	decoder := json.NewDecoder(resp.Body)
	animeMap := make(map[int]HummingbirdAnime)

	// read the first token (the bracket token)
	if _, err := decoder.Token(); err != nil {
		return err
	}

	for decoder.More() {
		var anime HummingbirdAnime
		if err := decoder.Decode(&anime); err != nil {
			return err
		}

		animeMap[anime.ID().Get(Hummingbird)] = anime
	}

	// read the last token (the bracket token)
	if _, err = decoder.Token(); err != nil {
		return err
	}

	newHummingbirdList := &HummingbirdAnimeList{
		username:    hal.username,
		authToken:   hal.authToken,
		anime:       animeMap,
		changes:     []Change{},
		pastChanges: []Change{},
	}

	changes := DiffHummingbirdLists(hal, newHummingbirdList)
	hal.anime = animeMap
	hal.changes = append(hal.changes, changes...)
	return nil
}

func (hal *HummingbirdAnimeList) Add(anime Anime) {
	id := anime.ID().Get(Hummingbird)
	hal.anime[id] = AnimeToHummingbird(anime)

	change := AddChange{Anime: anime}
	hal.changes = append(hal.changes, change)
}

func (hal *HummingbirdAnimeList) Edit(anime Anime) {
	animeID := anime.ID().Get(Hummingbird)
	oldAnime := hal.anime[animeID]
	hal.anime[animeID] = AnimeToHummingbird(anime)

	change := EditChange{
		OldAnime: oldAnime,
		NewAnime: anime,
	}

	hal.changes = append(hal.changes, change)
}

func (hal *HummingbirdAnimeList) Get(id int) (Anime, error) {
	if anime, ok := hal.anime[id]; ok {
		return anime, nil
	}
	return nil, errors.New(fmt.Sprintf("Anime with ID %d is not in the anime list", id))
}

func (hal *HummingbirdAnimeList) Remove(anime Anime) {
	delete(hal.anime, anime.ID().Get(Hummingbird))
	change := DeleteChange{Anime: anime}
	hal.changes = append(hal.changes, change)
}

func (hal *HummingbirdAnimeList) Push() error {
	mergedChanges := MergeChanges(hal.changes, Hummingbird)

	changeRequests := make([]*http.Request, len(mergedChanges))
	for i, change := range mergedChanges {
		request, err := hal.GenerateChange(change)
		if err != nil {
			return err
		}

		changeRequests[i] = request
	}

	// sends changes asynchronously
	resultCh := make(chan error)
	go SendRequests(changeRequests, resultCh, 3, func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return errors.New("Status code is not 200")
		}
		return nil
	})

	// return error if sending changes failed
	if err := <-resultCh; err != nil {
		return err
	}

	hal.pastChanges = append(hal.pastChanges, mergedChanges...)
	hal.changes = []Change{}
	return nil
}

func (hal *HummingbirdAnimeList) Undo() error {
	if len(hal.changes) <= 0 {
		return errors.New("Cannot undo from empty changelist")
	}

	change := hal.changes[len(hal.changes)-1]
	undoRequest, err := hal.GenerateChange(change, true)
	if err != nil {
		return err
	}

	client := &http.Client{}
	if _, err := client.Do(undoRequest); err != nil {
		return err
	}

	hal.changes = hal.changes[:len(hal.changes)-1]
	return nil
}

// GenerateChange returns a HTTP request that applies the change
func (hal *HummingbirdAnimeList) GenerateChange(change Change, undo ...bool) (*http.Request, error) {
	undoForm := len(undo) > 0 && undo[0]

	form := url.Values{}
	form.Add("auth_token", hal.AuthToken())

	change.FillForm(Hummingbird, &form, undoForm)
	return http.NewRequest("POST", change.URL(Hummingbird, undoForm), strings.NewReader(form.Encode()))
}

// DiffHummingbirdLists creates a list of changes from diffing two Hummingbird anime lists
func DiffHummingbirdLists(oldList *HummingbirdAnimeList, newList *HummingbirdAnimeList) []Change {
	var changes []Change
	for id, oldAnime := range oldList.anime {
		if newAnime, ok := newList.anime[id]; ok {
			if !reflect.DeepEqual(newAnime, oldAnime) {
				changes = append(changes, EditChange{OldAnime: oldAnime, NewAnime: newAnime})
			}
		} else {
			changes = append(changes, DeleteChange{Anime: oldAnime})
		}
	}
	for id, anime := range newList.anime {
		if _, ok := oldList.anime[id]; !ok {
			changes = append(changes, AddChange{Anime: anime})
		}
	}
	return changes
}
