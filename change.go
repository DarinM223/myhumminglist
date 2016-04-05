package main

import (
	"fmt"
	"net/url"
)

const (
	hummingbirdAddURL    = "https://hummingbird.me/api/v1/libraries/%d"
	hummingbirdEditURL   = "https://hummingbird.me/api/v1/libraries/%d"
	hummingbirdDeleteURL = "https://hummingbird.me/api/v1/libraries/%d/remove"
)

type Change interface {
	FillForm(listType int, form *url.Values, undo ...bool)
	URL(listType int, undo ...bool) string
}

type AddChange struct {
	Anime Anime
}

func (change AddChange) URL(listType int, undo ...bool) string {
	undoURL := false
	if len(undo) > 0 {
		if undo[0] == true {
			undoURL = true
		}
	}

	switch listType {
	case Hummingbird:
		if undoURL {
			return fmt.Sprintf(hummingbirdAddURL, change.Anime.ID().Get(Hummingbird))
		} else {
			return fmt.Sprintf(hummingbirdDeleteURL, change.Anime.ID().Get(Hummingbird))
		}
	case MyAnimeList:
		// TODO(DarinM223): set URL for MyAnimeList
		return ""
	default:
		panic("Invalid list type")
	}
}

func (change AddChange) FillForm(listType int, form *url.Values, undo ...bool) {
	undoForm := false
	if len(undo) > 0 {
		if undo[0] == true {
			undoForm = true
		}
	}

	if !undoForm {
		switch listType {
		case Hummingbird:
			// TODO(DarinM223): set form for Hummingbird request
			form.Add("status", StatusToHummingbirdString(change.Anime.Status()))
			form.Add("rewatching", fmt.Sprintf("%t", change.Anime.Rewatching()))
			form.Add("rewatched_times", fmt.Sprintf("%d", change.Anime.RewatchedTimes()))
			form.Add("episodes_watched", fmt.Sprintf("%d", change.Anime.EpisodesWatched()))
		case MyAnimeList:
			// TODO(DarinM223): set form for MyAnimeList request
		}
	}
}

type EditChange struct {
	OldAnime Anime
	NewAnime Anime
}

func (change EditChange) URL(listType int, undo ...bool) string {
	switch listType {
	case Hummingbird:
		return fmt.Sprintf(hummingbirdEditURL, change.NewAnime.ID().Get(Hummingbird))
	case MyAnimeList:
		// TODO(DarinM223): set URL for MyAnimeList
		return ""
	default:
		panic("Invalid list type")
	}
}

func (change EditChange) FillForm(listType int, form *url.Values, undo ...bool) {
	undoForm := false
	if len(undo) > 0 {
		if undo[0] == true {
			undoForm = true
		}
	}

	_ = undoForm

	switch listType {
	case Hummingbird:
		// TODO(DarinM223): set form for Hummingbird request
	case MyAnimeList:
		// TODO(DarinM223): set form for MyAnimeList request
	}
}

type DeleteChange struct {
	Anime Anime
}

func (change DeleteChange) URL(listType int, undo ...bool) string {
	// Reverse undo parameter
	undoURL := true
	if len(undo) > 0 {
		undoURL = !undo[0]
	}

	c := AddChange{Anime: change.Anime}
	return c.URL(listType, undoURL)
}

func (change DeleteChange) FillForm(listType int, form *url.Values, undo ...bool) {
	undoForm := false
	if len(undo) > 0 {
		if undo[0] == true {
			undoForm = true
		}
	}

	if undoForm {
		addChange := AddChange{Anime: change.Anime}
		addChange.FillForm(listType, form)
	}
}

// MergeChanges takes a list of changes and returns a
// smaller list with similar changes merged
func MergeChanges(changes []Change, listType int) []Change {
	addMap := make(map[int]Anime)
	editMap := make(map[int]Change)

	for _, change := range changes {
		switch c := change.(type) {
		case AddChange:
			addMap[c.Anime.ID().Get(listType)] = c.Anime
		case EditChange:
			animeID := c.NewAnime.ID().Get(listType)
			if _, ok := addMap[animeID]; ok {
				addMap[animeID] = c.NewAnime
			} else {
				editMap[animeID] = change
			}
		case DeleteChange:
			animeID := c.Anime.ID().Get(listType)
			delete(addMap, animeID)
			delete(editMap, animeID)
		}
	}

	var newChanges []Change
	for _, anime := range addMap {
		change := AddChange{
			Anime: anime,
		}
		newChanges = append(newChanges, change)
	}
	for _, change := range editMap {
		newChanges = append(newChanges, change)
	}

	return newChanges
}
