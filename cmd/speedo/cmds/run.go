package cmds

import (
	"errors"

	"github.com/showwin/speedtest-go/speedtest"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the speed test",
		Run:   Run,
	}
}

func Run(_ *cobra.Command, _ []string) {
	logger := newLogger()

	var speedtestClient = speedtest.New()

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{61387})

	logger.Debug("Running tests on servers", "servers", targets)

	for _, s := range targets {
		var errs []error
		errs = append(errs, s.DownloadTest())
		errs = append(errs, s.UploadTest())

		if err := errors.Join(errs...); err != nil {
			logger.Error("Test failed", "err", err.Error(), "errors", errs)
			continue
		}

		logger.Info(
			"Test done",
			"serverID",
			s.ID,
			"downloadMegabytesPerSecond",
			s.DLSpeed.Byte(speedtest.UnitTypeDecimalBytes),
			"uploadMegabytesPerSecond",
			s.ULSpeed.Byte(speedtest.UnitTypeDecimalBytes),
		)
		s.Context.Reset() // reset counter
	}
}
