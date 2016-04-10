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
		AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}}},
		false,
		fmt.Sprintf(HummingbirdAddURL, 69),
	},
	{
		Hummingbird,
		AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}}},
		true,
		fmt.Sprintf(HummingbirdDeleteURL, 69),
	},
	{
		Hummingbird,
		EditChange{
			HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}},
			HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}},
		},
		false,
		fmt.Sprintf(HummingbirdEditURL, 69),
	},
	{
		Hummingbird,
		EditChange{
			HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}},
			HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}},
		},
		true,
		fmt.Sprintf(HummingbirdEditURL, 69),
	},
	{
		Hummingbird,
		DeleteChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}}},
		false,
		fmt.Sprintf(HummingbirdDeleteURL, 69),
	},
	{
		Hummingbird,
		DeleteChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}}},
		true,
		fmt.Sprintf(HummingbirdAddURL, 69),
	},
}

var defaultHummingbirdAnime = HummingbirdAnime{
	JEpisodesWatched: 11,
	JStatus:          "currently-watching",
	JRewatchedTimes:  2,
	JRewatching:      false,
	JData:            HummingbirdAnimeData{Id: 69},
}

// TODO(DarinM223): add MyAnimeList tests
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
	{
		listType: Hummingbird,
		change: EditChange{
			HummingbirdAnime{
				JEpisodesWatched: 2,
				JStatus:          "currently-watching",
				JRewatchedTimes:  0,
				JRewatching:      false,
			},
			HummingbirdAnime{
				JEpisodesWatched: 3,
				JStatus:          "completed",
				JRewatchedTimes:  1,
				JRewatching:      true,
			},
		},
		undo: false,
		expectedForm: map[string][]string{
			"episodes_watched": []string{"3"},
			"status":           []string{"completed"},
			"rewatched_times":  []string{"1"},
			"rewatching":       []string{"true"},
		},
	},
	{
		listType: Hummingbird,
		change: EditChange{
			defaultHummingbirdAnime,
			defaultHummingbirdAnime,
		},
		undo:         false,
		expectedForm: map[string][]string{},
	},
	{
		listType: Hummingbird,
		change: EditChange{
			HummingbirdAnime{
				JEpisodesWatched: 2,
				JStatus:          "currently-watching",
				JRewatchedTimes:  0,
				JRewatching:      false,
			},
			HummingbirdAnime{
				JEpisodesWatched: 3,
				JStatus:          "completed",
				JRewatchedTimes:  1,
				JRewatching:      true,
			},
		},
		undo: true,
		expectedForm: map[string][]string{
			"episodes_watched": []string{"2"},
			"status":           []string{"currently-watching"},
			"rewatched_times":  []string{"0"},
			"rewatching":       []string{"false"},
		},
	},
	{
		listType: Hummingbird,
		change: EditChange{
			defaultHummingbirdAnime,
			defaultHummingbirdAnime,
		},
		undo:         true,
		expectedForm: map[string][]string{},
	},
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
					JEpisodesWatched: 11,
					JStatus:          "currently-watching",
					JRewatchedTimes:  2,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 420},
				},
			},
			EditChange{
				HummingbirdAnime{JData: HummingbirdAnimeData{Id: 420}},
				HummingbirdAnime{
					JEpisodesWatched: 12,
					JStatus:          "completed",
					JRewatchedTimes:  3,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 420},
				},
			},
			EditChange{
				HummingbirdAnime{JData: HummingbirdAnimeData{Id: 71}},
				HummingbirdAnime{
					JEpisodesWatched: 1,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 71},
				},
			},
			DeleteChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 72}}},
			DeleteChange{defaultHummingbirdAnime},
		},
		expectedChanges: []Change{
			AddChange{
				HummingbirdAnime{
					JEpisodesWatched: 12,
					JStatus:          "completed",
					JRewatchedTimes:  3,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 420},
				},
			},
			EditChange{
				HummingbirdAnime{JData: HummingbirdAnimeData{Id: 71}},
				HummingbirdAnime{
					JEpisodesWatched: 1,
					JStatus:          "currently-watching",
					JRewatchedTimes:  0,
					JRewatching:      false,
					JData:            HummingbirdAnimeData{Id: 71},
				},
			},
			DeleteChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 72}}},
		},
	},
	{
		listType: Hummingbird,
		changes: []Change{
			AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 420}}},
			AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}}},
			AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 620}}},
		},
		expectedChanges: []Change{
			AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 420}}},
			AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 69}}},
			AddChange{HummingbirdAnime{JData: HummingbirdAnimeData{Id: 620}}},
		},
	},
}

func TestChangeURL(t *testing.T) {
	for _, test := range changeURLTests {
		if test.change.URL(test.listType, test.undo) != test.expectedURL {
			t.Errorf("TestChangeURL failed: want %s got %s", test.expectedURL, test.change.URL(test.listType, test.undo))
		}
	}
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

func TestMergeChanges(t *testing.T) {
	for _, test := range mergeChangesTests {
		result := MergeChanges(test.changes, test.listType)
		if !reflect.DeepEqual(result, test.expectedChanges) {
			t.Errorf("TestMergeChanges failed: want %+v got %+v", test.expectedChanges, result)
		}
	}
}
