package glagol

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_IfEmptyBaseUrl(t *testing.T) {
	cli := NewClient("", "", "")
	require.NotEmpty(t, cli.baseUrl)
}

func TestClient_GetDevice(t *testing.T) {
	srv := mockSrv(getNotEmptyMockResponse(), http.StatusOK)
	defer srv.Close()

	cli := NewClient(srv.URL, "id", "token")
	actual, err := cli.GetDevice(context.Background(), "", "")
	expected := NewDevice("id", "", "")

	require.Nil(t, err)
	require.Equal(t, expected.GetId(), actual.GetId())
}

func TestClient_GetDeviceGlagolIsError(t *testing.T) {
	srv := mockSrv(getNotEmptyMockResponse(), http.StatusInternalServerError)
	defer srv.Close()

	cli := NewClient(srv.URL, "unknown", "")
	actual, err := cli.GetDevice(context.Background(), "", "")

	require.NotNil(t, err)
	require.Nil(t, actual)
}

func TestClient_GetDeviceWhenDeviceNotMatched(t *testing.T) {
	srv := mockSrv(getNotEmptyMockResponse(), http.StatusOK)
	defer srv.Close()

	cli := NewClient(srv.URL, "unknown", "")
	actual, err := cli.GetDevice(context.Background(), "", "")

	require.NotNil(t, err)
	require.Nil(t, actual)
}

func TestClient_GetDeviceWhenDeviceNotFound(t *testing.T) {
	srv := mockSrv(DeviceListResponse{}, http.StatusOK)
	defer srv.Close()

	cli := NewClient(srv.URL, "", "")
	actual, err := cli.GetDevice(context.Background(), "", "")

	require.NotNil(t, err)
	require.Nil(t, actual)
}

func mockSrv(mockResp DeviceListResponse, statusCode int) *httptest.Server {
	r := http.NewServeMux()
	r.HandleFunc(fmt.Sprintf("/%s", deviceListAction), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		data, _ := json.Marshal(mockResp)
		_, _ = w.Write(data)
	})
	r.HandleFunc(fmt.Sprintf("/%s", tokenAction), func(w http.ResponseWriter, r *http.Request) {
		body := TokenResponse{Token: "test"}
		data, _ := json.Marshal(body)
		_, _ = w.Write(data)
	})
	return httptest.NewServer(r)
}

func getNotEmptyMockResponse() DeviceListResponse {
	device := DeviceResponse{Id: "id", Platform: "platform"}
	device.Glagol.Security.ServerCertificate = "cert"
	device.Glagol.Security.ServerPrivateKey = "key"

	return DeviceListResponse{
		[]DeviceResponse{device},
	}
}
