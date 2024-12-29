package main

import (
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
	targets, _ := serverList.FindServer([]int{61387, 5747, 44599, 4087, 8827})

	for _, s := range targets {
		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()
		logger.Info("Test done", "serverID", s.ID, "latency", s.Lat, "downloadMbps", s.DLSpeed.Mbps(), "uploadMbps", s.ULSpeed.Mbps())
		s.Context.Reset() // reset counter
	}
}
