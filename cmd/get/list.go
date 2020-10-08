package get

import (
	"path"

	"github.com/CorentinB/Zeno/config"
	"github.com/CorentinB/Zeno/internal/pkg/crawl"
	"github.com/CorentinB/Zeno/internal/pkg/frontier"
	"github.com/google/uuid"
	"github.com/remeh/sizedwaitgroup"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewGetListCmd() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "Start crawling with a seed list",
		Action:    CmdGetList,
		Flags:     []cli.Flag{},
		UsageText: "<FILE> [ARGUMENTS]",
	}
}

func CmdGetList(c *cli.Context) error {
	err := initLogging(c)
	if err != nil {
		log.Error("Unable to parse arguments")
		return err
	}

	// Initialize Crawl
	crawl, err := crawl.Create()
	if err != nil {
		log.Error("Unable to initialize crawl job")
		return err
	}

	// If the job name isn't specified, we generate a random name
	if len(config.App.Flags.Job) == 0 {
		UUID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		config.App.Flags.Job = UUID.String()
	}

	crawl.WARC = config.App.Flags.WARC
	crawl.WARCRetry = config.App.Flags.WARCRetry
	crawl.WARCPrefix = config.App.Flags.WARCPrefix
	crawl.WARCOperator = config.App.Flags.WARCOperator
	crawl.API = config.App.Flags.API
	crawl.DomainsCrawl = config.App.Flags.DomainsCrawl
	crawl.Job = config.App.Flags.Job
	crawl.JobPath = path.Join("jobs", config.App.Flags.Job)
	crawl.UserAgent = config.App.Flags.UserAgent
	crawl.Headless = config.App.Flags.Headless
	crawl.Workers = config.App.Flags.Workers
	crawl.WorkerPool = sizedwaitgroup.New(crawl.Workers)
	crawl.Seencheck = config.App.Flags.Seencheck
	crawl.Proxy = config.App.Flags.Proxy
	crawl.MaxHops = uint8(config.App.Flags.MaxHops)

	// Initialize client
	crawl.InitHTTPClient()

	// Initialize initial seed list
	crawl.SeedList, err = frontier.IsSeedList(c.Args().Get(0))
	if err != nil || len(crawl.SeedList) <= 0 {
		logrus.WithFields(logrus.Fields{
			"input": c.Args().Get(0),
			"error": err.Error(),
		}).Error("This is not a valid input")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"input":      c.Args().Get(0),
		"seedsCount": len(crawl.SeedList),
	}).Print("Seed list loaded")

	// Start crawl
	err = crawl.Start()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"crawl": crawl,
			"error": err,
		}).Error("Crawl exited due to error")
		return err
	}

	return nil
}
