package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/downloader/internal/pkg/fsys"
	"github.com/dikaeinstein/downloader/pkg/downloader"
	"github.com/dikaeinstein/downloader/pkg/hash"
)

type cli struct {
	cfg           Config
	httpClient    *http.Client
	versionString string
}

type Config struct {
	checksum     string
	checksumURL  string
	checksumFile string
	downloadDir  string
	filename     string
	parallel     bool
	timeout      time.Duration
	url          string
	version      bool
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	if c.cfg.version {
		cmd.Println(c.versionString)
		return nil
	}

	c.httpClient.Timeout = c.cfg.timeout

	hasher, path := c.parseChecksum()
	dl, err := downloader.New(".",
		c.httpClient,
		fsys.OsFS{},
		hasher,
		&downloader.ProgressBar{},
		hash.Verifier{})
	if err != nil {
		return err
	}

	return dl.Download(context.Background(), c.cfg.url, c.cfg.filename, path)
}

func (c *cli) setupConfig(_ *cobra.Command, args []string) error {
	c.cfg.checksum = viper.GetString("checksum")
	c.cfg.checksumURL = viper.GetString("checksum-url")
	c.cfg.checksumFile = viper.GetString("checksum-file")
	c.cfg.parallel = viper.GetBool("parallel")
	c.cfg.timeout = viper.GetDuration("timeout")
	c.cfg.version = viper.GetBool("version")
	c.parseFilename()

	if len(args) == 1 {
		u, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		c.cfg.url = u.String()
	}

	return nil
}

// parseFilename initializes the cli config filename and downloadDir
// with the parsed filename flag .
func (c *cli) parseFilename() {
	filename := viper.GetString("filename")
	if filename == "" {
		c.cfg.filename = ""
		c.cfg.downloadDir = "."
	} else {
		c.cfg.filename = filename
		c.cfg.downloadDir = filepath.Dir(filename)
	}
}

// parseChecksum returns the appropriate Hasher based on input
// checksum flags.
func (c *cli) parseChecksum() (downloader.Hasher, string) {
	if c.cfg.checksum != "" {
		return hash.StringHasher(c.cfg.checksum), c.cfg.checksum
	}
	if c.cfg.checksumFile != "" {
		return hash.LocalHasher{}, c.cfg.checksumFile
	}
	if c.cfg.checksumURL != "" {
		return hash.NewRemoteHasher(c.httpClient), c.cfg.checksumURL
	}

	return nil, ""
}

func setupFlags(cmd *cobra.Command) error {
	cmd.Flags().BoolP("parallel", "p", false,
		"use parallel download")
	cmd.Flags().DurationP("timeout", "t", 300*time.Second,
		"timeout for the download")
	cmd.Flags().StringP("checksum", "c", "",
		"checksum to verify downloaded file")
	cmd.Flags().StringP("checksum-url", "", "",
		"url to download the checksum to verify downloaded file")
	cmd.Flags().StringP("checksum-file", "", "",
		"local file containing the checksum to verify downloaded file")
	cmd.Flags().StringP("filename", "f", "",
		"filename to use")
	cmd.Flags().BoolP("version", "v", false,
		"version of downloaderctl")

	return viper.BindPFlags(cmd.Flags())
}

func validArgLen(cmd *cobra.Command, args []string) error {
	if !viper.GetBool("version") && !viper.GetBool("help") && len(args) == 0 {
		return fmt.Errorf("url is required")
	}

	return nil
}

func newCommand(httpClient *http.Client, vOpt versionOption) *cobra.Command {
	cli := &cli{httpClient: httpClient}
	cli.versionString = fmt.Sprintf(
		"Version: %s\nGo version: %s\nGit hash: %s\nBuilt: %s",
		vOpt.BinaryVersion, vOpt.GoVersion, vOpt.GitHash, vOpt.BuildDate,
	)

	cmd := &cobra.Command{
		Use:     "downloaderctl [flags] url",
		Short:   "downloaderctl is a CLI tool which download files using the given url.",
		Long:    "downloaderctl is a CLI tool which download files using the given url.",
		Args:    validArgLen,
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}

	return cmd
}

type versionOption struct {
	BuildDate     string
	BinaryVersion string
	GitHash       string
	GoVersion     string
}

type Option struct {
	Version versionOption
}

// Run executes the downloaderctl command.
func Run(opt Option) {
	cmd := newCommand(&http.Client{}, opt.Version)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
