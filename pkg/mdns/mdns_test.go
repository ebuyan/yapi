package mdns

import (
	"github.com/hashicorp/mdns"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestDiscoverDeviceFound(t *testing.T) {
	serv, err := mdns.NewServer(&mdns.Config{Zone: makeServiceWithServiceName(t, "_foobar._tcp")})
	if err != nil {
		require.Nil(t, err)
	}
	defer func() { _ = serv.Shutdown() }()

	actual, err := Discover("123", "_foobar._tcp")

	require.Nil(t, err)
	require.Equal(t, "80", actual.Port)
	require.Equal(t, "192.168.0.42", actual.IpAddr)
}

func TestDiscoverDeviceNotFound(t *testing.T) {
	serv, err := mdns.NewServer(&mdns.Config{Zone: makeServiceWithServiceName(t, "_foobar._tcp")})
	if err != nil {
		require.Nil(t, err)
	}
	defer func() { _ = serv.Shutdown() }()

	actual, err := Discover("124", "_foobar._tcp")

	require.NotNil(t, err)
	require.False(t, actual.Discovered)
}

func makeServiceWithServiceName(t *testing.T, service string) *mdns.MDNSService {
	m, err := mdns.NewMDNSService(
		"hostname",
		service,
		"local.",
		"testhost.",
		80,
		[]net.IP{[]byte{192, 168, 0, 42}, net.ParseIP("2620:0:1000:1900:b0c2:d0b2:c411:18bc")},
		[]string{"deviceId=123"},
	)

	if err != nil {
		require.Nil(t, err)
	}

	return m
}
