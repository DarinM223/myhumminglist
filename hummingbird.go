package main

import "errors"

// HummingbirdAnime represents the JSON data of a Hummingbird library entry
type HummingbirdAnime struct {
	episodesWatched int                  `json:"episodes_watched"`
	status          string               `json:"status"`
	rewatchedTimes  int                  `json:"rewatched_times"`
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
	case HUMMINGBIRD:
		return ha.data.id, nil
	case MYANIMELIST:
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
		return STATUS_WATCHING, nil
	case "plan-to-watch":
		return STATUS_PLANTOWATCH, nil
	case "completed":
		return STATUS_COMPLETED, nil
	case "on-hold":
		return STATUS_ONHOLD, nil
	case "dropped":
		return STATUS_DROPPED, nil
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

type HummingbirdAnimeList struct {
	anime   map[int]HummingbirdAnime
	changes []Change
}

func (hal HummingbirdAnimeList) Type() int {
	return HUMMINGBIRD
}

func StatusToHummingbirdString(status int) (string, error) {
	switch status {
	case STATUS_WATCHING:
		return "currently-watching", nil
	case STATUS_PLANTOWATCH:
		return "plan-to-watch", nil
	case STATUS_COMPLETED:
		return "completed", nil
	case STATUS_ONHOLD:
		return "on-hold", nil
	case STATUS_DROPPED:
		return "dropped", nil
	default:
		return "", errors.New("Invalid status")
	}
}

func (hal *HummingbirdAnimeList) Add(anime Anime) error {
	id, err := anime.ID(hal.Type())
	if err != nil {
		return err
	}

	malID, err := anime.ID(MYANIMELIST)
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
	change := Change{
		LIST_ADD,
		[]interface{}{},
	}
	change.Args = append(change.Args, id)
	hal.changes = append(hal.changes, change)
	return nil
}

func (hal *HummingbirdAnimeList) Edit(anime Anime) error {
	id, err := anime.ID(hal.Type())
	if err != nil {
		return err
	}

	oldAnime := hal.anime[id]

	change := Change{
		LIST_EDIT,
		[]interface{}{},
	}
	// TODO(DarinM223) calculate diff and add to change list
	if oldAnime.Title() != anime.Title() {
		change.Args = append(change.Args, "title")
		change.Args = append(change.Args, oldAnime.Title())
		change.Args = append(change.Args, anime.Title())
	}
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

	change := Change{
		LIST_REMOVE,
		[]interface{}{},
	}
	change.Args = append(change.Args, anime)
	return nil
}

func (hal HummingbirdAnimeList) Push() {
	// TODO(DarinM223): implement this
}
