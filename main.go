package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/showwin/speedtest-go/speedtest"
)

func main() {
	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var speedtestClient = speedtest.New()

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{61387})

	for _, s := range targets {
		var errs []error
		errs = append(errs, s.DownloadTest())
		errs = append(errs, s.UploadTest())

		if err := errors.Join(errs...); err != nil {
			logger.Error("Test failed", "err", err.Error(), "errors", errs)
			continue
		}

		logger.Info("Test done", "serverID", s.ID, "downloadMbps", s.DLSpeed.Mbps(), "uploadMbps", s.ULSpeed.Mbps())
		s.Context.Reset() // reset counter
	}
}
