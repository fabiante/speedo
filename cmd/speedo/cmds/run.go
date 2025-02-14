package cmds

import (
	"context"
	"errors"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/showwin/speedtest-go/speedtest"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func NewRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the speed test",
	}

	// flags := cmd.Flags()

	// output := flags.String("output", "stdout", "How to output speed tests. Possible values: stdout (prints JSON to stdout) influx (sends a data line to InfluxDB)")

	client := influxdb2.NewClient(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"))
	writeAPI := client.WriteAPIBlocking(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))

	cmd.Run = func(_ *cobra.Command, _ []string) {
		// FIXME: Validate output strategy before doing any heavy lifting.

		logger := newLogger()

		var speedtestClient = speedtest.New()

		serverList, _ := speedtestClient.FetchServers()
		targets, _ := serverList.FindServer([]int{61387})

		logger.Debug("Running tests on servers", "servers", targets)

		now := time.Now()

		for _, s := range targets {
			var errs []error
			errs = append(errs, s.DownloadTest())
			errs = append(errs, s.UploadTest())

			if err := errors.Join(errs...); err != nil {
				logger.Error("Test failed", "err", err.Error(), "errors", errs)
				continue
			}

			// Log results
			logger.Info(
				"Test done",
				"serverID",
				s.ID,
				"downloadMegabytesPerSecond",
				s.DLSpeed.Byte(speedtest.UnitTypeDecimalBytes),
				"uploadMegabytesPerSecond",
				s.ULSpeed.Byte(speedtest.UnitTypeDecimalBytes),
			)

			// Send results to database
			tags := map[string]string{
				"host": os.Getenv("HOST"),
			}
			fields := map[string]interface{}{
				"downloadMegabytesPerSecond": float64(s.DLSpeed),
				"uploadMegabytesPerSecond":   float64(s.ULSpeed),
			}
			point := write.NewPoint("speeds", tags, fields, now)
			if err := writeAPI.WritePoint(context.Background(), point); err != nil {
				err = fmt.Errorf("sending value to influx db failed: %w", err)
				panic(err)
			}

			s.Context.Reset() // reset counter
		}
	}

	return cmd
}
