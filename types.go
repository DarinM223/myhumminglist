package main

const (
	Hummingbird = iota
	MyAnimeList = iota
)

const (
	StatusWatching    = iota
	StatusCompleted   = iota
	StatusOnHold      = iota
	StatusDropped     = iota
	StatusPlanToWatch = iota
)

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
	Add(anime Anime)
	Edit(anime Anime)
	Get(id int) (Anime, error)
	Remove(anime Anime)
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
