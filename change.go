package main

import (
	"fmt"
	"net/url"
)

type Change interface {
	FillForm(listType int, form *url.Values, undo ...bool)
	URL(listType int, undo ...bool) string
}

type AddChange struct {
	Anime Anime
}

func (change AddChange) URL(listType int, undo ...bool) string {
	undoURL := len(undo) > 0 && undo[0]

	switch listType {
	case Hummingbird:
		if undoURL {
			return fmt.Sprintf(HummingbirdDeleteURL, change.Anime.ID().Get(Hummingbird))
		} else {
			return fmt.Sprintf(HummingbirdAddURL, change.Anime.ID().Get(Hummingbird))
		}
	case MyAnimeList:
		// TODO(DarinM223): set URL for MyAnimeList
		return ""
	default:
		panic("Invalid list type")
	}
}

func (change AddChange) FillForm(listType int, form *url.Values, undo ...bool) {
	if len(undo) <= 0 || !undo[0] {
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
		return fmt.Sprintf(HummingbirdEditURL, change.NewAnime.ID().Get(Hummingbird))
	case MyAnimeList:
		// TODO(DarinM223): set URL for MyAnimeList
		return ""
	default:
		panic("Invalid list type")
	}
}

func (change EditChange) FillForm(listType int, form *url.Values, undo ...bool) {
	undoForm := len(undo) > 0 && undo[0]

	switch listType {
	case Hummingbird:
		if change.NewAnime.Status() != change.OldAnime.Status() {
			var status string
			if undoForm {
				status = StatusToHummingbirdString(change.OldAnime.Status())
			} else {
				status = StatusToHummingbirdString(change.NewAnime.Status())
			}
			form.Add("status", status)
		}
		if change.NewAnime.EpisodesWatched() != change.OldAnime.EpisodesWatched() {
			if undoForm {
				form.Add("episodes_watched", fmt.Sprintf("%d", change.OldAnime.EpisodesWatched()))
			} else {
				form.Add("episodes_watched", fmt.Sprintf("%d", change.NewAnime.EpisodesWatched()))
			}
		}
		if change.NewAnime.RewatchedTimes() != change.OldAnime.RewatchedTimes() {
			if undoForm {
				form.Add("rewatched_times", fmt.Sprintf("%d", change.OldAnime.RewatchedTimes()))
			} else {
				form.Add("rewatched_times", fmt.Sprintf("%d", change.NewAnime.RewatchedTimes()))
			}
		}
		if change.NewAnime.Rewatching() != change.OldAnime.Rewatching() {
			if undoForm {
				form.Add("rewatching", fmt.Sprintf("%t", change.OldAnime.Rewatching()))
			} else {
				form.Add("rewatching", fmt.Sprintf("%t", change.NewAnime.Rewatching()))
			}
		}
	case MyAnimeList:
		// TODO(DarinM223): set form for MyAnimeList request
	}
}

type DeleteChange struct {
	Anime Anime
}

func (change DeleteChange) URL(listType int, undo ...bool) string {
	undoURL := true
	if len(undo) > 0 {
		undoURL = !undo[0]
	}

	c := AddChange{Anime: change.Anime}
	return c.URL(listType, undoURL)
}

func (change DeleteChange) FillForm(listType int, form *url.Values, undo ...bool) {
	if len(undo) > 0 && undo[0] {
		addChange := AddChange{Anime: change.Anime}
		addChange.FillForm(listType, form)
	}
}

// MergeChanges takes a list of changes and returns a
// smaller list with similar changes merged
func MergeChanges(changes []Change, listType int) []Change {
	addMap := NewLRUMap()
	editMap := NewLRUMap()
	deleteMap := NewLRUMap()

	for _, change := range changes {
		switch c := change.(type) {
		case AddChange:
			animeID := c.Anime.ID().Get(listType)
			addMap.Add(animeID, c.Anime)
			if deleteMap.Contains(animeID) {
				deleteMap.Remove(animeID)
			}
		case EditChange:
			animeID := c.NewAnime.ID().Get(listType)
			if addMap.Contains(animeID) {
				addMap.Add(animeID, c.NewAnime)
			} else {
				editMap.Add(animeID, change)
			}
			if deleteMap.Contains(animeID) {
				deleteMap.Remove(animeID)
			}
		case DeleteChange:
			animeID := c.Anime.ID().Get(listType)
			if !addMap.Contains(animeID) && !editMap.Contains(animeID) {
				deleteMap.Add(animeID, change)
			} else {
				addMap.Remove(animeID)
				editMap.Remove(animeID)
			}
		}
	}

	var newChanges []Change
	for _, key := range addMap.Keys() {
		anime, _ := addMap.Get(key)
		change := AddChange{Anime: anime.(Anime)}
		newChanges = append(newChanges, change)
	}
	for _, key := range editMap.Keys() {
		change, _ := editMap.Get(key)
		newChanges = append(newChanges, change.(Change))
	}
	for _, key := range deleteMap.Keys() {
		change, _ := deleteMap.Get(key)
		newChanges = append(newChanges, change.(Change))
	}

	return newChanges
}
