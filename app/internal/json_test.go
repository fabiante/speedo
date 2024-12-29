package internal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

const logStr = `{"time":"2024-12-29T11:33:09.180353921Z","level":"INFO","msg":"Test done","serverID":"61387","downloadMbps":42.2179619498218,"uploadMbps":11.840630020360306}
{"time":"2024-12-29T11:33:31.737024763Z","level":"INFO","msg":"Test done","serverID":"5747","downloadMbps":67.06493106824415,"uploadMbps":11.48459446046503}
{"time":"2024-12-29T11:33:56.114286512Z","level":"INFO","msg":"Test done","serverID":"44599","downloadMbps":63.85029234769865,"uploadMbps":10.855524410402655}
{"time":"2024-12-29T11:34:18.916153583Z","level":"INFO","msg":"Test done","serverID":"4087","downloadMbps":69.43697006014746,"uploadMbps":11.173423881436058}
{"time":"2024-12-29T11:34:41.689635342Z","level":"INFO","msg":"Test done","serverID":"8827","downloadMbps":59.297743734555475,"uploadMbps":12.228545856200613}
{"time":"2024-12-29T11:40:09.365673899Z","level":"INFO","msg":"Test done","serverID":"61387","downloadMbps":36.06291909814911,"uploadMbps":12.724933754728996}
{"time":"2024-12-29T11:40:37.166599832Z","level":"INFO","msg":"Test done","serverID":"5747","downloadMbps":15.47387624581305,"uploadMbps":10.84535869156537}
{"time":"2024-12-29T11:40:59.888744991Z","level":"INFO","msg":"Test done","serverID":"44599","downloadMbps":64.0008008626761,"uploadMbps":10.606752259511548}
{"time":"2024-12-29T11:41:22.478660288Z","level":"INFO","msg":"Test done","serverID":"4087","downloadMbps":69.29717190418434,"uploadMbps":11.01642886113242}
{"time":"2024-12-29T11:41:45.18130534Z","level":"INFO","msg":"Test done","serverID":"8827","downloadMbps":61.32759494845209,"uploadMbps":12.166745045819727}
{"time":"2024-12-29T12:40:30.555009436Z","level":"INFO","msg":"Test done","serverID":"61387","downloadMbps":42.397848314562964,"uploadMbps":13.231294560757313}`

func TestDecodeLog(t *testing.T) {
	lines, err := DecodeLog(bytes.NewReader([]byte(logStr)))
	require.NoError(t, err)
	require.Equal(t, 11, len(lines))
}
