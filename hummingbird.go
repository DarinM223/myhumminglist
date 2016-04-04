package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// HummingbirdAnime represents the JSON data of a Hummingbird library entry
type HummingbirdAnime struct {
	episodesWatched int                  `json:"episodes_watched"`
	status          string               `json:"status"`
	rewatchedTimes  int                  `json:"rewatched_times"`
	rewatching      bool                 `json:"rewatching"`
	data            HummingbirdAnimeData `json:"anime"`
}

// HummingbirdAnimeData represents the JSON data of a Hummingbird anime
type HummingbirdAnimeData struct {
	id            int    `json:"id"`
	malID         int    `json:"mal_id"`
	title         string `json:"title"`
	episodeCount  int    `json:"episode_count"`
	episodeLength int    `json:"episode_length"`
}

func (ha HummingbirdAnime) ID() AnimeID {
	return AnimeID{
		Hummingbird: ha.data.id,
		MyAnimeList: ha.data.malID,
	}
}

func (ha HummingbirdAnime) Title() string {
	return ha.data.title
}

func (ha HummingbirdAnime) Status() int {
	switch ha.status {
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
	return ha.episodesWatched
}

func (ha HummingbirdAnime) RewatchedTimes() int {
	return ha.rewatchedTimes
}

func (ha HummingbirdAnime) Rewatching() bool {
	return ha.rewatching
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

const (
	hummingbirdAddURL    = "https://hummingbird.me/api/v1/libraries/%d"
	hummingbirdEditURL   = "https://hummingbird.me/api/v1/libraries/%d"
	hummingbirdDeleteURL = "https://hummingbird.me/api/v1/libraries/%d/remove"
)

type HummingbirdAnimeList struct {
	anime       map[int]HummingbirdAnime
	changes     []*Change
	pastChanges []*Change
	authToken   string
}

func (hal HummingbirdAnimeList) Type() int {
	return Hummingbird
}

func (hal HummingbirdAnimeList) AuthToken() string {
	return hal.authToken
}

func (hal *HummingbirdAnimeList) Add(anime HummingbirdAnime) error {
	id := anime.ID().Get(Hummingbird)
	malID := anime.ID().Get(MyAnimeList)
	status := StatusToHummingbirdString(anime.Status())

	hal.anime[id] = HummingbirdAnime{
		episodesWatched: anime.EpisodesWatched(),
		rewatchedTimes:  anime.RewatchedTimes(),
		status:          status,
		data: HummingbirdAnimeData{
			id:    id,
			malID: malID,
			title: anime.Title(),
		},
	}
	change := &Change{
		ListAdd,
		[]interface{}{},
	}
	change.Args = append(change.Args, anime)
	hal.changes = append(hal.changes, change)
	return nil
}

func (hal *HummingbirdAnimeList) Edit(anime Anime) error {
	id := anime.ID().Get(hal.Type())
	oldAnime := hal.anime[id]

	change := &Change{
		ListEdit,
		[]interface{}{},
	}
	change.Args = append(change.Args, anime.ID())
	change.PopulateDiff(oldAnime, anime)

	hal.changes = append(hal.changes, change)
	return nil
}

func (hal *HummingbirdAnimeList) Get(id int) Anime {
	return hal.anime[id]
}

func (hal *HummingbirdAnimeList) Remove(anime Anime) error {
	id := anime.ID().Get(hal.Type())
	delete(hal.anime, id)

	change := &Change{
		ListRemove,
		[]interface{}{},
	}
	change.Args = append(change.Args, anime)
	hal.changes = append(hal.changes, change)
	return nil
}

func (hal *HummingbirdAnimeList) Push() error {
	mergedChanges := MergeChanges(hal.changes, Hummingbird)

	var changeRequests []*http.Request
	for _, change := range mergedChanges {
		request, err := hal.GenerateChange(change)
		if err != nil {
			return err
		}

		changeRequests = append(changeRequests, request)
	}

	// TODO(DarinM223): send asynchronously using goroutines
	client := &http.Client{}
	for _, request := range changeRequests {
		if _, err := client.Do(request); err != nil {
			return err
		}
	}

	hal.pastChanges = append(hal.pastChanges, mergedChanges...)
	hal.changes = []*Change{}
	return nil
}

func (hal *HummingbirdAnimeList) Undo() error {
	if len(hal.changes) <= 0 {
		return errors.New("Cannot undo from empty changelist")
	}

	change := hal.changes[len(hal.changes)-1]
	undoRequest, err := hal.GenerateUndo(change)
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
func (hal *HummingbirdAnimeList) GenerateChange(change *Change) (*http.Request, error) {
	switch change.ChangeType {
	case ListAdd:
		anime := change.Args[0].(Anime)
		animeID := anime.ID().Get(Hummingbird)

		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())
		form.Add("status", StatusToHummingbirdString(anime.Status()))
		form.Add("rewatching", fmt.Sprintf("%t", anime.Rewatching()))
		form.Add("rewatched_times", fmt.Sprintf("%d", anime.RewatchedTimes()))
		form.Add("episodes_watched", fmt.Sprintf("%d", anime.EpisodesWatched()))

		return http.NewRequest("POST", fmt.Sprintf(hummingbirdAddURL, animeID), strings.NewReader(form.Encode()))
	case ListEdit:
		editID := change.Args[0].(AnimeID)
		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())
		for i := 1; i < len(change.Args); i += 3 {
			title := change.Args[i].(string)
			new := change.Args[i+2]

			switch t := new.(type) {
			case int:
				form.Add(title, fmt.Sprintf("%d", t))
			case bool:
				form.Add(title, fmt.Sprintf("%t", t))
			case string:
				form.Add(title, t)
			default:
				panic("Change type not found")
			}
		}

		return http.NewRequest("POST", fmt.Sprintf(hummingbirdEditURL, editID.Get(Hummingbird)), strings.NewReader(form.Encode()))
	case ListRemove:
		anime := change.Args[0].(Anime)
		animeID := anime.ID().Get(Hummingbird)

		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())

		return http.NewRequest("POST", fmt.Sprintf(hummingbirdDeleteURL, animeID), strings.NewReader(form.Encode()))
	default:
		panic("Invalid change type")
	}
}

// GenerateUndo returns a HTTP request that undos the change
func (hal *HummingbirdAnimeList) GenerateUndo(change *Change) (*http.Request, error) {
	switch change.ChangeType {
	case ListAdd:
		anime := change.Args[0].(Anime)
		animeID := anime.ID().Get(Hummingbird)

		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())

		return http.NewRequest("POST", fmt.Sprintf(hummingbirdDeleteURL, animeID), strings.NewReader(form.Encode()))
	case ListEdit:
		editID := change.Args[0].(AnimeID)
		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())
		for i := 1; i < len(change.Args); i += 3 {
			title := change.Args[i].(string)
			old := change.Args[i+1]

			switch t := old.(type) {
			case int:
				form.Add(title, fmt.Sprintf("%d", t))
			case bool:
				form.Add(title, fmt.Sprintf("%t", t))
			case string:
				form.Add(title, t)
			default:
				panic("Change type not found")
			}
		}

		return http.NewRequest("POST", fmt.Sprintf(hummingbirdEditURL, editID.Get(Hummingbird)), strings.NewReader(form.Encode()))
	case ListRemove:
		anime := change.Args[0].(Anime)
		animeID := anime.ID().Get(Hummingbird)

		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())
		form.Add("status", StatusToHummingbirdString(anime.Status()))
		form.Add("rewatching", fmt.Sprintf("%t", anime.Rewatching()))
		form.Add("rewatched_times", fmt.Sprintf("%d", anime.RewatchedTimes()))
		form.Add("episodes_watched", fmt.Sprintf("%d", anime.EpisodesWatched()))

		return http.NewRequest("POST", fmt.Sprintf(hummingbirdAddURL, animeID), strings.NewReader(form.Encode()))
	default:
		panic("Invalid change type")
	}
}
