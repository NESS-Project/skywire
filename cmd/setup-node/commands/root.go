package commands

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	logrussyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/skycoin/dmsg/discord"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/spf13/cobra"

	"github.com/skycoin/skywire/pkg/metrics"
	"github.com/skycoin/skywire/pkg/setup"
	"github.com/skycoin/skywire/pkg/util/buildinfo"
)

var (
	metricsAddr  string
	syslogAddr   string
	tag          string
	cfgFromStdin bool
)

var rootCmd = &cobra.Command{
	Use:   "setup-node [config.json]",
	Short: "Route Setup Node for skywire",
	Run: func(_ *cobra.Command, args []string) {
		if _, err := buildinfo.Get().WriteTo(log.Writer()); err != nil {
			log.Printf("Failed to output build info: %v", err)
		}

		logger := logging.MustGetLogger(tag)
		if syslogAddr != "" {
			hook, err := logrussyslog.NewSyslogHook("udp", syslogAddr, syslog.LOG_INFO, tag)
			if err != nil {
				logger.Fatalf("Unable to connect to syslog daemon on %v", syslogAddr)
			}
			logging.AddHook(hook)
		}

		if discordWebhookURL := discord.GetWebhookURLFromEnv(); discordWebhookURL != "" {
			hook := discord.NewHook(tag, discordWebhookURL)
			logging.AddHook(hook)
		}

		var rdr io.Reader
		var err error

		if !cfgFromStdin {
			configFile := "config.json"

			if len(args) > 0 {
				configFile = args[0]
			}
			rdr, err = os.Open(configFile)
			if err != nil {
				log.Fatalf("Failed to open config: %s", err)
			}
		} else {
			logger.Info("Reading config from STDIN")
			rdr = bufio.NewReader(os.Stdin)
		}

		conf := &setup.Config{}

		raw, err := ioutil.ReadAll(rdr)
		if err != nil {
			logger.Fatalf("Failed to read config: %v", err)
		}

		if err := json.Unmarshal(raw, &conf); err != nil {
			logger.WithField("raw", string(raw)).Fatalf("Failed to decode config: %s", err)
		}

		logger.Infof("Config: %#v", conf)

		sn, err := setup.NewNode(conf, metrics.NewPrometheus("setupnode"))
		if err != nil {
			logger.Fatal("Failed to create a setup node: ", err)
		}

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(metricsAddr, nil); err != nil {
				logger.Println("Failed to start metrics API:", err)
			}
		}()

		logger.Fatal(sn.Serve())
	},
}

func init() {
	rootCmd.Flags().StringVarP(&metricsAddr, "metrics", "m", ":2121", "address to bind metrics API to")
	rootCmd.Flags().StringVar(&syslogAddr, "syslog", "", "syslog server address. E.g. localhost:514")
	rootCmd.Flags().StringVar(&tag, "tag", "setup-node", "logging tag")
	rootCmd.Flags().BoolVarP(&cfgFromStdin, "stdin", "i", false, "read config from STDIN")
}

// Execute executes root CLI command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
