package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/downloader/internal/pkg/fsys"
	"github.com/dikaeinstein/downloader/pkg/downloader"
	"github.com/dikaeinstein/downloader/pkg/hash"
)

type cli struct {
	cfg        Config
	httpClient *http.Client
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
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	c.httpClient.Timeout = c.cfg.timeout

	hasher, path := c.parseChecksum()
	dl, err := downloader.New(".",
		c.httpClient,
		fsys.OsFS{},
		hasher,
		downloader.DefaultProgress{},
		hash.Verifier{})
	if err != nil {
		return err
	}

	return dl.Download(context.Background(), c.cfg.url, c.cfg.filename, path)
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	c.cfg.checksum = viper.GetString("checksum")
	c.cfg.checksumURL = viper.GetString("checksum-url")
	c.cfg.checksumFile = viper.GetString("checksum-file")
	c.cfg.parallel = viper.GetBool("parallel")
	c.cfg.timeout = viper.GetDuration("timeout")
	c.parseFilename()
	c.cfg.url = args[0]

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
	cmd.Flags().BoolP("parallel", "p", false, "Use parallel download.")
	cmd.Flags().DurationP("timeout", "t", 10*time.Second, "Timeout for the download.")
	cmd.Flags().StringP("checksum", "c", "", "Checksum of the file.")
	cmd.Flags().StringP("checksum-url", "", "", "Url to download the checksum from.")
	cmd.Flags().StringP("checksum-file", "", "", "Local file containing the checksum.")
	cmd.Flags().StringP("filename", "f", "", "Filename to use.")

	return viper.BindPFlags(cmd.Flags())
}

func newCommand(httpClient *http.Client, vOpt VersionOption) *cobra.Command {
	cli := &cli{httpClient: httpClient}

	cmd := &cobra.Command{
		Use:     "downloaderctl",
		Short:   "downloaderctl is a CLI tool which download files using the given url.",
		Long:    "downloaderctl is a CLI tool which download files using the given url.",
		Example: "downloaderctl [<options>] [url]",
		Args:    cobra.ExactArgs(1),
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}

	cmd.SetVersionTemplate(
		fmt.Sprintf("Version: %s\nGo version: %s\nGit hash: %s\nBuilt: %s\n",
			vOpt.BinaryVersion, vOpt.GoVersion, vOpt.GitHash, vOpt.BuildDate),
	)

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}

	return cmd
}

type VersionOption struct {
	BuildDate     string
	BinaryVersion string
	GitHash       string
	GoVersion     string
}

type Option struct {
	Version VersionOption
}

func Run(opt Option) {
	cmd := newCommand(&http.Client{}, opt.Version)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}