package glagol

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDevice_GetId(t *testing.T) {
	d := getDevice()
	require.Equal(t, "id", d.GetId())
}

func TestDevice_GetCertificate(t *testing.T) {
	d := getDevice()
	require.Equal(t, "cert", d.GetCertificate())
}

func TestDevice_SetHost(t *testing.T) {
	d := getDevice()
	d.SetHost("127.0.0.1", "8080")
	require.Equal(t, "wss://127.0.0.1:8080/", d.GetHost())
}

func TestDevice_RefreshToken(t *testing.T) {
	d := getDevice()
	d.SetRefreshTokenHandler(func(ctx context.Context, deviceId, platform string) (string, error) {
		return "test_token", nil
	})
	err := d.RefreshToken(context.Background())
	require.Nil(t, err)
	require.Equal(t, "test_token", d.GetToken())
}

func TestDevice_SetState(t *testing.T) {
	state := DeviceState{}
	expected, _ := json.Marshal(&state)

	d := getDevice()
	d.SetState(expected)

	require.Equal(t, expected, d.GetState())
}

func getDevice() *Device {
	return NewDevice("id", "platform", "cert")
}
