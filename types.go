package main

const (
	Hummingbird = iota
	MyAnimeList = iota
)

const (
	ListAdd    = iota
	ListEdit   = iota
	ListRemove = iota
)

const (
	StatusWatching    = iota
	StatusCompleted   = iota
	StatusOnHold      = iota
	StatusDropped     = iota
	StatusPlanToWatch = iota
)

type Change struct {
	ChangeType int
	Args       []interface{}
}

func (change *Change) PopulateDiff(oldAnime Anime, anime Anime) {
	if oldAnime.Title() != anime.Title() {
		change.Args = append(change.Args, "title")
		change.Args = append(change.Args, oldAnime.Title())
		change.Args = append(change.Args, anime.Title())
	}
	if oldAnime.Status() != anime.Status() {
		change.Args = append(change.Args, "status")
		change.Args = append(change.Args, oldAnime.Status())
		change.Args = append(change.Args, anime.Status())
	}
	if oldAnime.EpisodesWatched() != anime.EpisodesWatched() {
		change.Args = append(change.Args, "episodes_watched")
		change.Args = append(change.Args, oldAnime.EpisodesWatched())
		change.Args = append(change.Args, anime.EpisodesWatched())
	}
	if oldAnime.RewatchedTimes() != anime.RewatchedTimes() {
		change.Args = append(change.Args, "rewatched_times")
		change.Args = append(change.Args, oldAnime.RewatchedTimes())
		change.Args = append(change.Args, anime.RewatchedTimes())
	}
}

// MergeChanges takes a list of changes and returns a
// smaller list with similar changes merged
func MergeChanges(changes []*Change, listType int) []*Change {
	addMap := make(map[int]Anime)
	editMap := make(map[int]*Change)

	for _, change := range changes {
		switch change.ChangeType {
		case ListAdd:
			anime := change.Args[0].(Anime)
			animeID := anime.ID().Get(listType)

			addMap[animeID] = anime
		case ListEdit:
			editID := change.Args[0].(AnimeID).Get(listType)

			if _, ok := addMap[editID]; ok {
				// TODO(DarinM223): apply edit changes to added anime
			} else {
				editMap[editID] = change
			}
		case ListRemove:
			animeID := change.Args[0].(Anime).ID().Get(listType)
			delete(addMap, animeID)
			delete(editMap, animeID)
		}
	}

	var newChanges []*Change
	for _, anime := range addMap {
		change := &Change{
			ListAdd,
			[]interface{}{},
		}
		change.Args = append(change.Args, anime)
		newChanges = append(newChanges, change)
	}
	for _, change := range editMap {
		newChanges = append(newChanges, change)
	}

	return newChanges
}

type AnimeID struct {
	Hummingbird int
	MyAnimeList int
}

func (id AnimeID) Get(listType int) int {
	switch listType {
	case Hummingbird:
		return id.Hummingbird
	case MyAnimeList:
		return id.MyAnimeList
	default:
		panic("Invalid anime list")
	}
}

type Anime interface {
	ID() AnimeID
	Title() string
	Status() int
	EpisodesWatched() int
	RewatchedTimes() int
	Rewatching() bool
}

type Animelist interface {
	Type() int
	Add(anime Anime) error
	Edit(anime Anime) error
	Get(id int) Anime
	Remove(anime Anime) error
	Push()
	Undo()
	Anime() []Anime
	Contains(id int) bool
	AuthToken() string
}

// Manages multiple anime lists by syncing changes to the others
type AnimelistManager struct {
	primary  Animelist
	replicas []Animelist
}

// Add adds an anime to all of the lists
func (m *AnimelistManager) Add(anime Anime) {
	m.primary.Add(anime)
	for _, replica := range m.replicas {
		replica.Add(anime)
	}
}

// Edit changes an anime to all of the lists
func (m *AnimelistManager) Edit(anime Anime) {
	m.primary.Edit(anime)
	for _, replica := range m.replicas {
		replica.Edit(anime)
	}
}

// Remove an anime from all of the lists
func (m *AnimelistManager) Remove(anime Anime) {
	m.primary.Remove(anime)
	for _, replica := range m.replicas {
		replica.Remove(anime)
	}
}

// Sync syncs the replica lists to the primary list
func (m *AnimelistManager) Sync() error {
	for id, anime := range m.primary.Anime() {
		for _, replica := range m.replicas {
			if replica.Contains(id) {
				replica.Edit(anime)
			} else {
				replica.Add(anime)
			}
		}
	}

	return nil
}
