package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rogpeppe/go-internal/semver"

	"github.com/hetznercloud/cli/internal/state/config"
)

const cliRepo = "hetznercloud/cli"

func GetLatestReleaseVersion() (ver string, err error) {
	url := "https://api.github.com/repos/" + cliRepo + "/releases/latest"
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(err, res.Body.Close())
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func CheckForUpdate(cfg config.Config) {
	lastCheck := cfg.LastUpdateCheck()
	if time.Since(lastCheck) < 24*time.Hour {
		return
	}
	cfg.SetLastUpdateCheck(time.Now())
	_ = cfg.Write()

	current := "v" + version
	latest, err := GetLatestReleaseVersion()
	if err != nil {
		return
	}

	if semver.Compare(current, latest) < 0 {
		_, _ = fmt.Fprintf(os.Stderr, "You are using an outdated version of the Hetzner Cloud CLI.\nYour version: %s\nCurrent version: %s\n\n", current, latest)
	}
}
