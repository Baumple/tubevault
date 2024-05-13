package data

import (
	"os"
	"path/filepath"
	"strings"
)

func GetSaveDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	saveDir := filepath.Join(homeDir, ".tubevault")
	err = os.Mkdir(saveDir, 0777)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return "", err
	}

	return saveDir, nil

}

func GetSaveFile() (*os.File, error) {
	dirPath, err := GetSaveDirPath()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(dirPath, "config.json"), os.O_RDWR|os.O_CREATE, 0777)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetPlaylistDir returns a path to where playlists are stored.
// Usually playlists are stored in HOME_DIR/.watchvault/playlists
func GetPlaylistDir() (string, error) {
	saveDir, err := GetSaveDirPath()
	if err != nil {
		return "", nil
	}

	playlistDir := filepath.Join(saveDir, PLAYLIST_DIR)

	err = os.Mkdir(playlistDir, 0777)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return "", err
	}

	return playlistDir, nil
}

type DataRetriever interface {
	GetPlaylists() ([]Playlist, error)
	SavePlaylist(playlist *Playlist) error
	DeletePlaylist(id string) error
	UpdateVideoWatched(playlistId string, id string, watched bool) error
	Close()
}

