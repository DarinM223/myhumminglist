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
					episodesWatched: 11,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 0},
				},
				2: HummingbirdAnime{
					episodesWatched: 5,
					status:          "currently-watching",
					rewatchedTimes:  1,
					rewatching:      true,
					data:            HummingbirdAnimeData{id: 2},
				},
			},
		},
		newList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				1: HummingbirdAnime{
					episodesWatched: 1,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 1},
				},
				2: HummingbirdAnime{
					episodesWatched: 6,
					status:          "currently-watching",
					rewatchedTimes:  1,
					rewatching:      true,
					data:            HummingbirdAnimeData{id: 2},
				},
			},
		},
		expectedChanges: []Change{
			DeleteChange{
				Anime: HummingbirdAnime{
					episodesWatched: 11,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 0},
				},
			},
			EditChange{
				OldAnime: HummingbirdAnime{
					episodesWatched: 5,
					status:          "currently-watching",
					rewatchedTimes:  1,
					rewatching:      true,
					data:            HummingbirdAnimeData{id: 2},
				},
				NewAnime: HummingbirdAnime{
					episodesWatched: 6,
					status:          "currently-watching",
					rewatchedTimes:  1,
					rewatching:      true,
					data:            HummingbirdAnimeData{id: 2},
				},
			},
			AddChange{
				Anime: HummingbirdAnime{
					episodesWatched: 1,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 1},
				},
			},
		},
	},
	{
		oldList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				0: HummingbirdAnime{
					episodesWatched: 11,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      true,
					data:            HummingbirdAnimeData{id: 0},
				},
				1: HummingbirdAnime{
					episodesWatched: 12,
					status:          "completed",
					rewatchedTimes:  1,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 1},
				},
			},
		},
		newList: &HummingbirdAnimeList{
			anime: map[int]HummingbirdAnime{
				0: HummingbirdAnime{
					episodesWatched: 11,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      true,
					data:            HummingbirdAnimeData{id: 0},
				},
				1: HummingbirdAnime{
					episodesWatched: 12,
					status:          "completed",
					rewatchedTimes:  1,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 1},
				},
			},
		},
		expectedChanges: nil,
	},
}

func TestDiffHummingbirdLists(t *testing.T) {
	for _, test := range diffHummingbirdTests {
		changes := DiffHummingbirdLists(test.oldList, test.newList)
		if !reflect.DeepEqual(test.expectedChanges, changes) {
			t.Errorf("TestDiffHummingbirdLists failed: expected %+v got %+v", test.expectedChanges, changes)
		}
	}
}
