package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type FileSystemPlayerStore struct {
	Database *json.Encoder
	league   League
}

func initializePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("Problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initializePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("Problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("Problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		Database: json.NewEncoder(&Tape{file}),
		league:   league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.Database.Encode(f.league)
}

func (f *FileSystemPlayerStore) GetPLayerScore(name string) int {
	player := f.league.Find(name)

	if player != nil {

		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (wins int) {
	for _, player := range f.GetLeague() {
		if player.Name == name {
			wins = player.Wins
			break
		}
	}

	return
}

func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
  db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

  if err != nil {
    return nil, nil, fmt.Errorf("Problem opening %s %v", path, err)
  }

  closeFunc := func() {
    db.Close()
  }

  store, err := NewFileSystemPlayerStore(db)

  if err != nil {
    return nil, nil, fmt.Errorf("problem creating file system player store, %v", err)
  }

  return store, closeFunc, nil
}
