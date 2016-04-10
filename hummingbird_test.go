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
