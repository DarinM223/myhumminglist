package main

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

// TODO(DarinM223): add MyAnimeList tests
var changeURLTests = []struct {
	listType    int
	change      Change
	undo        bool
	expectedURL string
}{
	{
		Hummingbird,
		AddChange{HummingbirdAnime{data: HummingbirdAnimeData{id: 69}}},
		false,
		fmt.Sprintf(hummingbirdAddURL, 69),
	},
	{
		Hummingbird,
		AddChange{HummingbirdAnime{data: HummingbirdAnimeData{id: 69}}},
		true,
		fmt.Sprintf(hummingbirdDeleteURL, 69),
	},
	{
		Hummingbird,
		EditChange{
			HummingbirdAnime{data: HummingbirdAnimeData{id: 69}},
			HummingbirdAnime{data: HummingbirdAnimeData{id: 69}},
		},
		false,
		fmt.Sprintf(hummingbirdEditURL, 69),
	},
	{
		Hummingbird,
		EditChange{
			HummingbirdAnime{data: HummingbirdAnimeData{id: 69}},
			HummingbirdAnime{data: HummingbirdAnimeData{id: 69}},
		},
		true,
		fmt.Sprintf(hummingbirdEditURL, 69),
	},
	{
		Hummingbird,
		DeleteChange{HummingbirdAnime{data: HummingbirdAnimeData{id: 69}}},
		false,
		fmt.Sprintf(hummingbirdDeleteURL, 69),
	},
	{
		Hummingbird,
		DeleteChange{HummingbirdAnime{data: HummingbirdAnimeData{id: 69}}},
		true,
		fmt.Sprintf(hummingbirdAddURL, 69),
	},
}

func TestChangeURL(t *testing.T) {
	for _, test := range changeURLTests {
		if test.change.URL(test.listType, test.undo) != test.expectedURL {
			t.Errorf("TestChangeURL failed: want %s got %s", test.expectedURL, test.change.URL(test.listType, test.undo))
		}
	}
}

var defaultHummingbirdAnime = HummingbirdAnime{
	episodesWatched: 11,
	status:          "currently-watching",
	rewatchedTimes:  2,
	rewatching:      false,
	data:            HummingbirdAnimeData{id: 69},
}

// TODO(DarinM223): add MyAnimeList tests and edit tests for Hummingbird
var changeFillFormTests = []struct {
	listType     int
	change       Change
	undo         bool
	expectedForm url.Values
}{
	{
		listType: Hummingbird,
		change:   AddChange{defaultHummingbirdAnime},
		undo:     false,
		expectedForm: map[string][]string{
			"status":           []string{"currently-watching"},
			"rewatching":       []string{"false"},
			"rewatched_times":  []string{"2"},
			"episodes_watched": []string{"11"},
		},
	},
	{
		listType:     Hummingbird,
		change:       AddChange{defaultHummingbirdAnime},
		undo:         true,
		expectedForm: map[string][]string{},
	},
	{
		listType:     Hummingbird,
		change:       DeleteChange{defaultHummingbirdAnime},
		undo:         false,
		expectedForm: map[string][]string{},
	},
	{
		listType: Hummingbird,
		change:   DeleteChange{defaultHummingbirdAnime},
		undo:     true,
		expectedForm: map[string][]string{
			"status":           []string{"currently-watching"},
			"rewatching":       []string{"false"},
			"rewatched_times":  []string{"2"},
			"episodes_watched": []string{"11"},
		},
	},
}

func TestChangeFillForm(t *testing.T) {
	for _, test := range changeFillFormTests {
		form := url.Values{}
		test.change.FillForm(test.listType, &form, test.undo)
		if !reflect.DeepEqual(form, test.expectedForm) {
			t.Errorf("TestChangeFillForm failed: want %v got %v", test.expectedForm, form)
		}
	}
}

var mergeChangesTests = []struct {
	listType        int
	changes         []Change
	expectedChanges []Change
}{
	{
		listType: Hummingbird,
		changes: []Change{
			AddChange{defaultHummingbirdAnime},
			AddChange{
				HummingbirdAnime{
					episodesWatched: 11,
					status:          "currently-watching",
					rewatchedTimes:  2,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 420},
				},
			},
			EditChange{
				HummingbirdAnime{data: HummingbirdAnimeData{id: 420}},
				HummingbirdAnime{
					episodesWatched: 12,
					status:          "completed",
					rewatchedTimes:  3,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 420},
				},
			},
			EditChange{
				HummingbirdAnime{data: HummingbirdAnimeData{id: 71}},
				HummingbirdAnime{
					episodesWatched: 1,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 71},
				},
			},
			DeleteChange{HummingbirdAnime{data: HummingbirdAnimeData{id: 72}}},
			DeleteChange{defaultHummingbirdAnime},
		},
		expectedChanges: []Change{
			AddChange{
				HummingbirdAnime{
					episodesWatched: 12,
					status:          "completed",
					rewatchedTimes:  3,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 420},
				},
			},
			EditChange{
				HummingbirdAnime{data: HummingbirdAnimeData{id: 71}},
				HummingbirdAnime{
					episodesWatched: 1,
					status:          "currently-watching",
					rewatchedTimes:  0,
					rewatching:      false,
					data:            HummingbirdAnimeData{id: 71},
				},
			},
			DeleteChange{HummingbirdAnime{data: HummingbirdAnimeData{id: 72}}},
		},
	},
}

func TestMergeChanges(t *testing.T) {
	for _, test := range mergeChangesTests {
		result := MergeChanges(test.changes, test.listType)
		if !reflect.DeepEqual(result, test.expectedChanges) {
			t.Errorf("TestMergeChanges failed: want %+v got %+v", test.expectedChanges, result)
		}
	}
}
