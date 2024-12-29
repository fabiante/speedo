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

	// Use a proxy for the speedtest. eg: socks://127.0.0.1:7890
	// speedtest.WithUserConfig(&speedtest.UserConfig{Proxy: "socks://127.0.0.1:7890"})(speedtestClient)

	// Select a network card as the data interface.
	// speedtest.WithUserConfig(&speedtest.UserConfig{Source: "192.168.1.101"})(speedtestClient)

	// Get user's network information
	// user, _ := speedtestClient.FetchUserInfo()

	// Get a list of servers near a specified location
	// user.SetLocationByCity("Tokyo")
	// user.SetLocation("Osaka", 34.6952, 135.5006)

	// Search server using serverID.
	// eg: fetch server with ID 28910.
	// speedtest.ErrServerNotFound will be returned if the server cannot be found.
	// server, err := speedtest.FetchServerByID("28910")

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{61387})

	for _, s := range targets {
		var errs []error
		errs = append(errs, s.PingTest(nil))
		errs = append(errs, s.DownloadTest())
		errs = append(errs, s.UploadTest())

		if err := errors.Join(errs...); err != nil {
			logger.Error("Test failed", "err", err.Error(), "errors", errs)
		}

		logger.Info("Test done", "serverID", s.ID, "latency", s.Lat, "downloadMbps", s.DLSpeed.Mbps(), "uploadMbps", s.ULSpeed.Mbps())
		s.Context.Reset() // reset counter
	}
}
