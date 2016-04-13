package main

import (
	"reflect"
	"testing"
)

var diffHummingbirdTests = []struct {
	oldList         *HummingbirdAnimeList
	newList         *HummingbirdAnimeList
	expectedChanges []Change
}{
	{
		oldList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				0: HummingbirdAnime{
					NumEpisodesWatched: 11,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  0,
					IsRewatching:       false,
					Data:               HummingbirdAnimeData{Id: 0},
				},
				2: HummingbirdAnime{
					NumEpisodesWatched: 5,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  1,
					IsRewatching:       true,
					Data:               HummingbirdAnimeData{Id: 2},
				},
			},
		},
		newList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				1: HummingbirdAnime{
					NumEpisodesWatched: 1,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  0,
					IsRewatching:       false,
					Data:               HummingbirdAnimeData{Id: 1},
				},
				2: HummingbirdAnime{
					NumEpisodesWatched: 6,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  1,
					IsRewatching:       true,
					Data:               HummingbirdAnimeData{Id: 2},
				},
			},
		},
		expectedChanges: []Change{
			DeleteChange{
				Anime: HummingbirdAnime{
					NumEpisodesWatched: 11,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  0,
					IsRewatching:       false,
					Data:               HummingbirdAnimeData{Id: 0},
				},
			},
			EditChange{
				OldAnime: HummingbirdAnime{
					NumEpisodesWatched: 5,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  1,
					IsRewatching:       true,
					Data:               HummingbirdAnimeData{Id: 2},
				},
				NewAnime: HummingbirdAnime{
					NumEpisodesWatched: 6,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  1,
					IsRewatching:       true,
					Data:               HummingbirdAnimeData{Id: 2},
				},
			},
			AddChange{
				Anime: HummingbirdAnime{
					NumEpisodesWatched: 1,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  0,
					IsRewatching:       false,
					Data:               HummingbirdAnimeData{Id: 1},
				},
			},
		},
	},
	{
		oldList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				0: HummingbirdAnime{
					NumEpisodesWatched: 11,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  0,
					IsRewatching:       true,
					Data:               HummingbirdAnimeData{Id: 0},
				},
				1: HummingbirdAnime{
					NumEpisodesWatched: 12,
					AnimeStatus:        "completed",
					NumRewatchedTimes:  1,
					IsRewatching:       false,
					Data:               HummingbirdAnimeData{Id: 1},
				},
			},
		},
		newList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				0: HummingbirdAnime{
					NumEpisodesWatched: 11,
					AnimeStatus:        "currently-watching",
					NumRewatchedTimes:  0,
					IsRewatching:       true,
					Data:               HummingbirdAnimeData{Id: 0},
				},
				1: HummingbirdAnime{
					NumEpisodesWatched: 12,
					AnimeStatus:        "completed",
					NumRewatchedTimes:  1,
					IsRewatching:       false,
					Data:               HummingbirdAnimeData{Id: 1},
				},
			},
		},
		expectedChanges: nil,
	},
}

func TestDiffHummingbirdLists(t *testing.T) {
	for _, test := range diffHummingbirdTests {
		changes := DiffHummingbirdLists(test.oldList, test.newList)
		for _, expectedChange := range test.expectedChanges {
			changeEqual := false
			for _, actualChange := range changes {
				if reflect.DeepEqual(actualChange, expectedChange) {
					changeEqual = true
					break
				}
			}

			if !changeEqual {
				t.Errorf("TestDiffHummingbirdLists failed: expected %+v got %+v", test.expectedChanges, changes)
			}
		}
	}
}

func TestHummingbirdAnimeList_FetchFromEmpty(t *testing.T) {
	list := NewHummingbirdAnimeList("darin_minamoto", "")
	if err := list.Fetch(); err != nil {
		t.Errorf("TestHummingbirdAnimeList_FetchFromEmpty failed: " + err.Error())
	}
	if len(list.changes) <= 0 {
		t.Errorf("TestHummingbirdAnimeList_FetchFromEmpty failed: list changes haven't been populated")
	}
	if len(list.anime) <= 0 {
		t.Errorf("TestHummingbirdAnimeList_FetchFromEmpty failed: anime changes haven't been populated")
	}

	for _, change := range list.changes {
		switch ch := change.(type) {
		case AddChange:
			if ch.Anime.ID().Get(Hummingbird) == 0 {
				t.Errorf("TestHummingbirdAnimeList_FetchFromEmpty failed: anime id should not be zero")
			}
		default:
			t.Errorf("TestHummingbirdAnimeList_FetchFromEmpty failed: there should only be Add changes")
		}
	}
}

func TestHummingbirdAnimeList_AddHummingbirdAnime(t *testing.T) {
	list := NewHummingbirdAnimeList("darin_minamoto", "")

	anime := HummingbirdAnime{
		NumEpisodesWatched: 12,
		NumRewatchedTimes:  2,
		AnimeStatus:        "currently-watching",
		Data: HummingbirdAnimeData{
			Id:    50,
			MalID: 20,
			Title: "Sample text",
		},
	}
	list.Add(anime)

	if !reflect.DeepEqual(list.anime[50], anime) {
		t.Errorf("TestHummingbirdAnimeList_AddHummingbirdAnime failed: expected: %+v, got %+v", anime, list.anime[50])
	}
	if !reflect.DeepEqual(list.changes[len(list.changes)-1], AddChange{anime}) {
		t.Errorf("TestHummingbirdAnimeList_AddHummingbirdAnime failed: expected change %+v", AddChange{anime})
	}
}

func TestHummingbirdAnimeList_EditHummingbirdAnime(t *testing.T) {
	list := NewHummingbirdAnimeList("darin_minamoto", "")
	oldAnime := HummingbirdAnime{
		NumEpisodesWatched: 11,
		NumRewatchedTimes:  2,
		AnimeStatus:        "currently-watching",
		Data: HummingbirdAnimeData{
			Id:    50,
			MalID: 20,
			Title: "Sample text",
		},
	}
	newAnime := HummingbirdAnime{
		NumEpisodesWatched: 12,
		NumRewatchedTimes:  3,
		AnimeStatus:        "completed",
		Data: HummingbirdAnimeData{
			Id:    50,
			MalID: 20,
			Title: "Sample text",
		},
	}
	list.Add(oldAnime)
	list.Edit(newAnime)

	if !reflect.DeepEqual(list.anime[50], newAnime) {
		t.Errorf("TestHummingbirdAnimeList_EditHummingbirdAnime failed: expected: %+v, got %+v",
			newAnime, list.anime[50])
	}
	if !reflect.DeepEqual(list.changes[len(list.changes)-1], EditChange{OldAnime: oldAnime, NewAnime: newAnime}) {
		t.Errorf("TestHummingbirdAnimeList_EditHummingbirdAnime failed: expected change %+v",
			EditChange{OldAnime: oldAnime, NewAnime: newAnime})
	}
}

func TestHummingbirdAnimeList_RemoveHummingbirdAnime(t *testing.T) {
	list := NewHummingbirdAnimeList("darin_minamoto", "")
	anime := HummingbirdAnime{
		NumEpisodesWatched: 12,
		NumRewatchedTimes:  2,
		AnimeStatus:        "currently-watching",
		Data: HummingbirdAnimeData{
			Id:    50,
			MalID: 20,
			Title: "Sample text",
		},
	}
	list.Add(anime)
	list.Remove(anime)

	if _, ok := list.anime[50]; ok {
		t.Errorf("TestHummingbirdAnimeList_RemoveHummingbirdAnime failed: expected anime to not exist after removed")
	}
	if !reflect.DeepEqual(list.changes[len(list.changes)-1], DeleteChange{Anime: anime}) {
		t.Errorf("TestHummingbirdAnimeList_RemoveHummingbirdAnime failed: expected change %+v",
			DeleteChange{Anime: anime})
	}
}
