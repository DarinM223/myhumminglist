package main

const (
	HUMMINGBIRD = iota
	MYANIMELIST = iota
)

const (
	LIST_ADD    = iota
	LIST_EDIT   = iota
	LIST_REMOVE = iota
)

const (
	STATUS_WATCHING    = iota
	STATUS_COMPLETED   = iota
	STATUS_ONHOLD      = iota
	STATUS_DROPPED     = iota
	STATUS_PLANTOWATCH = iota
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
	oldStatus, err := oldAnime.Status()
	newStatus, err2 := anime.Status()
	if err == nil && err2 == nil && oldStatus != newStatus {
		change.Args = append(change.Args, "status")
		change.Args = append(change.Args, oldStatus)
		change.Args = append(change.Args, newStatus)
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

type Anime interface {
	ID(listType int) (int, error)
	Title() string
	Status() (int, error)
	EpisodesWatched() int
	RewatchedTimes() int
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
	for _, anime := range m.primary.Anime() {
		for _, replica := range m.replicas {
			id, err := anime.ID(replica.Type())
			if err != nil {
				return err
			}

			if replica.Contains(id) {
				replica.Edit(anime)
			} else {
				replica.Add(anime)
			}
		}
	}

	return nil
}
