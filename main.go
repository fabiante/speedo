package main

import (
	csv2 "encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fabiante/speedo/app"
	"github.com/showwin/speedtest-go/speedtest"
)

func main() {
	args := os.Args[1:]

	var cmd string

	if len(args) == 0 {
		cmd = "run"
	} else {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "run":
		run()
	case "csv":
		csv(args)
	default:
		fmt.Println("Unknown command", "cmd", cmd)
		os.Exit(1)
	}
}

func csv(args []string) {
	filePath := args[0]

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	lines, err := app.DecodeLog(file)
	if err != nil {
		panic(err)
	}

	writer := csv2.NewWriter(os.Stdout)
	writer.Comma = ';'
	writer.UseCRLF = true
	defer writer.Flush()

	writer.Write([]string{"time", "downloadMbps", "uploadMbps"})

	for _, line := range lines {
		writer.Write([]string{
			line.Time.Format("01.02.2006 15:04"),
			fmt.Sprintf("%f", line.DownloadMbps),
			fmt.Sprintf("%f", line.UploadMbps),
		})
	}
}

func run() {
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

		logger.Info("Test done", "serverID", s.ID, "downloadMbps", s.DLSpeed.Mbps(), "uploadMbps", s.ULSpeed.Mbps())
		s.Context.Reset() // reset counter
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func newLogger() *slog.Logger {
	levels := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	levelStr := strings.ToLower(getEnv("LOG_LEVEL", "info"))
	level := levels[levelStr]

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
