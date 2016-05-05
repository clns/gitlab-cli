package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update this tool to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		rel, err := getLatestRelease()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to retrieve the latest release", err)
			os.Exit(2)
		}
		if *rel.TagName == Version {
			fmt.Fprintln(os.Stdout, "No update available, latest version is", Version)
			os.Exit(0)
		}

		fmt.Printf("New update available: %s. Your current version is %s.\n", *rel.TagName, Version)

		asset, err := getReleaseAsset(rel)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		// Download asset
		fmt.Printf("downloading %d MB...\n", *asset.Size/1024/1024)
		req, err := http.NewRequest("GET", *asset.BrowserDownloadURL, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		req.Header.Add("Accept", "application/octet-stream")
		c := &http.Client{Timeout: 30 * time.Second}
		resp, err := c.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "%s: %s", resp.Status, b)
			os.Exit(2)
		}

		// Update
		if err := update.Apply(resp.Body, update.Options{}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		fmt.Println("gitlab-cli updated successfully to", *rel.TagName)
		saveVersion(*rel.TagName)
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

// Update functions

var lastUpdateCheck = filepath.Join(os.TempDir(), "gitlab-cli-latest-release")

func shouldCheckForUpdate() bool {
	fi, err := os.Stat(lastUpdateCheck)
	return !(err == nil && time.Now().Sub(fi.ModTime()) < 2*time.Minute)
}

func getLatestRelease() (rel *github.RepositoryRelease, err error) {
	gh := github.NewClient(&http.Client{Timeout: 1 * time.Second})
	rel, _, err = gh.Repositories.GetLatestRelease("clns", "gitlab-cli")
	return
}

// getReleaseAsset returns the platform-specific asset of the release.
// Be careful because it can be nil if not found.
func getReleaseAsset(rel *github.RepositoryRelease) (*github.ReleaseAsset, error) {
	var file string
	switch runtime.GOOS {
	case "windows":
		file = "gitlab-cli-Windows-x86_64.exe"
	case "linux":
		file = "gitlab-cli-Linux-x86_64"
	case "darwin":
		file = "gitlab-cli-Darwin-x86_64"
	default:
		return nil, fmt.Errorf("Unsupported platform")
	}
	for _, a := range rel.Assets {
		if *a.Name == file {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("Binary not found for your platform")
}

func CheckUpdate() {
	if !shouldCheckForUpdate() {
		log.Println("update: skip checking")
		b, err := ioutil.ReadFile(lastUpdateCheck)
		if err == nil && len(b) > 0 {
			printUpdateAvl(string(b))
		}
		return
	}
	rel, err := getLatestRelease()
	defer func() {
		ver := ""
		if rel != nil {
			ver = *rel.TagName
		}
		saveVersion(ver)
	}()
	if err != nil {
		log.Println("update:", err)
		return
	}
	printUpdateAvl(*rel.TagName)
}

func printUpdateAvl(latest string) {
	if latest != Version {
		fmt.Printf("New update available: %s. Run 'gitlab-cli update' to update.\n", latest)
	}
}

func saveVersion(ver string) {
	if err := ioutil.WriteFile(lastUpdateCheck, []byte(ver), os.ModePerm); err != nil {
		log.Println("update:", err)
	}
}
