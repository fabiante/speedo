package cmds

import (
	"encoding/csv"
	"os"

	"github.com/fabiante/speedo/app"
	"github.com/spf13/cobra"
)

func NewCSV() *cobra.Command {
	return &cobra.Command{
		Use:   "csv",
		Short: "Converts the log file to a CSV file",
		Args:  cobra.ExactArgs(1),
		Run:   CSV,
	}
}

func CSV(_ *cobra.Command, args []string) {
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

	writer := csv.NewWriter(os.Stdout)
	writer.Comma = ';'
	writer.UseCRLF = true
	defer writer.Flush()

	writer.Write([]string{"time", "downloadMbps", "uploadMbps"})

	for _, line := range lines {
		writer.Write([]string{
			line.Time.Format("01.02.2006 15:04"),
			line.DownloadMegabytesPerSecond,
			line.UploadMegabytesPerSecond,
		})
	}
}
