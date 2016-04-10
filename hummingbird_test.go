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
					JEpisodesWatched: 11,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 0},
				},
				2: HummingbirdAnime{
					JEpisodesWatched: 5,
					JStatus:          "currently-watching",
					JRewatchedTimes:  1,
					JRewatching:      true,
					JData:            HummingbirdAnimeData{Id: 2},
				},
			},
		},
		newList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				1: HummingbirdAnime{
					JEpisodesWatched: 1,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 1},
				},
				2: HummingbirdAnime{
					JEpisodesWatched: 6,
					JStatus:          "currently-watching",
					JRewatchedTimes:  1,
					JRewatching:      true,
					JData:            HummingbirdAnimeData{Id: 2},
				},
			},
		},
		expectedChanges: []Change{
			DeleteChange{
				Anime: HummingbirdAnime{
					JEpisodesWatched: 11,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 0},
				},
			},
			EditChange{
				OldAnime: HummingbirdAnime{
					JEpisodesWatched: 5,
					JStatus:          "currently-watching",
					JRewatchedTimes:  1,
					JRewatching:      true,
					JData:            HummingbirdAnimeData{Id: 2},
				},
				NewAnime: HummingbirdAnime{
					JEpisodesWatched: 6,
					JStatus:          "currently-watching",
					JRewatchedTimes:  1,
					JRewatching:      true,
					JData:            HummingbirdAnimeData{Id: 2},
				},
			},
			AddChange{
				Anime: HummingbirdAnime{
					JEpisodesWatched: 1,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 1},
				},
			},
		},
	},
	{
		oldList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				0: HummingbirdAnime{
					JEpisodesWatched: 11,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      true,
					JData:            HummingbirdAnimeData{Id: 0},
				},
				1: HummingbirdAnime{
					JEpisodesWatched: 12,
					JStatus:          "completed",
					JRewatchedTimes:  1,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 1},
				},
			},
		},
		newList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				0: HummingbirdAnime{
					JEpisodesWatched: 11,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      true,
					JData:            HummingbirdAnimeData{Id: 0},
				},
				1: HummingbirdAnime{
					JEpisodesWatched: 12,
					JStatus:          "completed",
					JRewatchedTimes:  1,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 1},
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
	err := list.Fetch()
	if err != nil {
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
