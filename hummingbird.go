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

func (ha HummingbirdAnime) ID(listType int) (int, error) {
	switch listType {
	case Hummingbird:
		return ha.data.id, nil
	case MyAnimeList:
		return ha.data.malID, nil
	default:
		return -1, errors.New("Invalid anime list type")
	}
}

func (ha HummingbirdAnime) Title() string {
	return ha.data.title
}

func (ha HummingbirdAnime) Status() (int, error) {
	switch ha.status {
	case "currently-watching":
		return StatusWatching, nil
	case "plan-to-watch":
		return StatusPlanToWatch, nil
	case "completed":
		return StatusCompleted, nil
	case "on-hold":
		return StatusOnHold, nil
	case "dropped":
		return StatusDropped, nil
	default:
		return -1, errors.New("Invalid status")
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

func StatusToHummingbirdString(status int) (string, error) {
	switch status {
	case StatusWatching:
		return "currently-watching", nil
	case StatusPlanToWatch:
		return "plan-to-watch", nil
	case StatusCompleted:
		return "completed", nil
	case StatusOnHold:
		return "on-hold", nil
	case StatusDropped:
		return "dropped", nil
	default:
		return "", errors.New("Invalid status")
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
	id, err := anime.ID(hal.Type())
	if err != nil {
		return err
	}

	malID, err := anime.ID(MyAnimeList)
	if err != nil {
		return err
	}

	status, err := anime.Status()
	if err != nil {
		return err
	}

	statusStr, err := StatusToHummingbirdString(status)
	if err != nil {
		return err
	}

	hal.anime[id] = HummingbirdAnime{
		episodesWatched: anime.EpisodesWatched(),
		rewatchedTimes:  anime.RewatchedTimes(),
		status:          statusStr,
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
	id, err := anime.ID(hal.Type())
	if err != nil {
		return err
	}

	oldAnime := hal.anime[id]

	change := &Change{
		ListEdit,
		[]interface{}{},
	}
	change.Args = append(change.Args, id)
	change.PopulateDiff(oldAnime, anime)

	hal.changes = append(hal.changes, change)
	return nil
}

func (hal *HummingbirdAnimeList) Get(id int) Anime {
	return hal.anime[id]
}

func (hal *HummingbirdAnimeList) Remove(anime Anime) error {
	id, err := anime.ID(hal.Type())
	if err != nil {
		return err
	}
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
	changeRequests := make([]*http.Request, len(hal.changes))
	for _, change := range hal.changes {
		req, err := hal.GenerateChange(change)
		if err != nil {
			return err
		}

		changeRequests = append(changeRequests, req)
	}

	// TODO(DarinM223): send the change requests asynchronously and block

	hal.pastChanges = append(hal.pastChanges, hal.changes...)
	hal.changes = []*Change{}
	return nil
}

// GenerateChange returns a HTTP request that applies the change
func (hal *HummingbirdAnimeList) GenerateChange(change *Change) (*http.Request, error) {
	switch change.ChangeType {
	case ListAdd:
		anime := change.Args[0].(Anime)

		animeID, err := anime.ID(Hummingbird)
		if err != nil {
			return nil, err
		}

		status, err := anime.Status()
		if err != nil {
			return nil, err
		}

		statusStr, err := StatusToHummingbirdString(status)
		if err != nil {
			return nil, err
		}

		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())
		form.Add("status", statusStr)
		form.Add("rewatching", fmt.Sprintf("%t", anime.Rewatching()))
		form.Add("rewatched_times", fmt.Sprintf("%d", anime.RewatchedTimes()))
		form.Add("episodes_watched", fmt.Sprintf("%d", anime.EpisodesWatched()))

		req, err := http.NewRequest("POST", fmt.Sprintf(hummingbirdAddURL, animeID), strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}
		return req, nil
	case ListEdit:
		editID := change.Args[0].(int)
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
				return nil, errors.New("Change type not found")
			}
		}
		req, err := http.NewRequest("POST", fmt.Sprintf(hummingbirdEditURL, editID), strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}

		return req, nil
	case ListRemove:
		anime := change.Args[0].(Anime)

		animeID, err := anime.ID(Hummingbird)
		if err != nil {
			return nil, err
		}

		form := url.Values{}
		form.Add("auth_token", hal.AuthToken())

		req, err := http.NewRequest("POST", fmt.Sprintf(hummingbirdDeleteURL, animeID), strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}

		return req, nil
	default:
		return nil, errors.New("Invalid change type")
	}
}

// GenerateUndo returns a HTTP request that undos the change
func (hal *HummingbirdAnimeList) GenerateUndo(change *Change) (*http.Request, error) {
	// TODO(DarinM223): implement this
	return nil, nil
}
