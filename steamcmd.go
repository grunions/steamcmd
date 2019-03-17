package steamcmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

var (
	// ErrNotImplemented gets returned if something doesnt work yet
	ErrNotImplemented = errors.New("Not implemented")
)

// SteamCmd is a wrapper around the Steam CMD
type SteamCmd struct {
	sync.Mutex // mutex for operations

	SteamCmdDir string
	AppBasePath string

	LoginUser string
	LoginPass string

	Debug bool
}

// New creates a new steamcmd instance, if path is empty, a temporary
// path will be created. Otherwise an existing instance will be reused.
func New(user, pass, path string) *SteamCmd {
	var err error

	// if path was left empty, we create a temporary dir for steam
	if path == "" {
		path, err = ioutil.TempDir("", "steamcmd")
		if err != nil {
			err = errors.Wrap(err, "Could not create steamcmd temp dir")
			panic(err)
		}
	}

	gamesPath, err := ioutil.TempDir("", "games")
	if err != nil {
		err = errors.Wrap(err, "Could not create apps temp dir")
		panic(err)
	}

	scmd := &SteamCmd{
		SteamCmdDir: path,
		AppBasePath: gamesPath,
		LoginUser:   user,
		LoginPass:   pass,
	}

	return scmd
}

// EnsureInstalled checks if the SteamCmd is executable, and bootstraps it otherwise.
// Remember, steam needs curl, bzip2, tar and lib32gcc1
func (scmd SteamCmd) EnsureInstalled() error {
	const SteamCmdURL = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"

	// check if SteamCmd needs to be bootstrapped
	fi, err := os.Stat(filepath.Join(scmd.SteamCmdDir, "steamcmd.sh"))

	// check whether it does not exist or is not executable
	if err != nil || (fi.Mode()&0111 == 0) {
		// we should bootstrap it again
		if err := os.MkdirAll(scmd.SteamCmdDir, 0755); err != nil {
			return errors.Wrap(err, "Could not create steam path")
		}

		task := exec.Command("bash", "-c", "curl -sqL \""+SteamCmdURL+"\" | tar xzvf -")

		task.Dir = scmd.SteamCmdDir

		if scmd.Debug {
			task.Stdout = os.Stdout
			task.Stderr = os.Stderr
		}

		err := task.Run()
		if err != nil {
			return errors.Wrap(err, "SteamCmd download failed")
		}

	}
	return nil
}

// GetAppPath returns the path where an app would be installed
func (scmd *SteamCmd) GetAppPath(id int) string {
	return filepath.Join(scmd.AppBasePath, strconv.Itoa(id))
}

// InstallUpdateApp installs and updates a given app
func (scmd *SteamCmd) InstallUpdateApp(id int) error {
	return scmd.run("+login", "anonymous",
		"+force_install_dir", scmd.GetAppPath(id),
		"+app_update", strconv.Itoa(id), "validate")
}

// AppInstalledVersion returns the Build ID of an Steam App
func (scmd *SteamCmd) AppInstalledVersion(id int) (int, error) {
	return 0, ErrNotImplemented
}

// AppAvailableVersion returns the latest Build ID for the Public branch
// of an Steam App
func (scmd *SteamCmd) AppAvailableVersion(id int) (int, error) {
	return 0, ErrNotImplemented
}

// DownloadWorkshopMod tries to download a mod from the workshop
func (scmd *SteamCmd) DownloadWorkshopMod(appid, id int) error {
	return ErrNotImplemented
}

// run helper
// * exit status 8 - no subscription
func (scmd *SteamCmd) run(params ...string) error {
	params = append(params, "+quit")

	task := exec.Command("./steamcmd.sh", params...)
	task.Dir = scmd.SteamCmdDir
	if scmd.Debug {
		task.Stdout = os.Stdout
		task.Stderr = os.Stderr
	}
	if err := task.Run(); err != nil {
		return errors.Wrap(err, "raw command failed")
	}

	return nil
}
