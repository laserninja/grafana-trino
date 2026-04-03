package trino

import (
	"net/url"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
)

func TestBuildTransport_Default(t *testing.T) {
	s := &Settings{
		URL:  &url.URL{Scheme: "http", Host: "localhost:8080"},
		Opts: httpclient.Options{},
	}

	transport, err := buildTransport(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.TLSClientConfig.InsecureSkipVerify {
		t.Error("InsecureSkipVerify should be false by default")
	}
	if transport.TLSClientConfig.RootCAs != nil {
		t.Error("RootCAs should be nil by default")
	}
	if len(transport.TLSClientConfig.Certificates) != 0 {
		t.Error("should have no client certificates by default")
	}
}

func TestBuildTransport_InsecureSkipVerify(t *testing.T) {
	s := &Settings{
		URL: &url.URL{Scheme: "https", Host: "localhost:8080"},
		Opts: httpclient.Options{
			TLS: &httpclient.TLSOptions{
				InsecureSkipVerify: true,
			},
		},
	}

	transport, err := buildTransport(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !transport.TLSClientConfig.InsecureSkipVerify {
		t.Error("InsecureSkipVerify should be true")
	}
}

func TestBuildTransport_CACertificate(t *testing.T) {
	// Valid PEM certificate (self-signed test cert)
	caCert := `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2wpSek3WVhMt
acyunmpJkCm+HlWV0FAUQPQE0SHiMX0CIHtxPnwGlLstqEmHLwqnI8+I8WhfIjMN
MzFIM8E/H4hZ
-----END CERTIFICATE-----`

	s := &Settings{
		URL: &url.URL{Scheme: "https", Host: "localhost:8080"},
		Opts: httpclient.Options{
			TLS: &httpclient.TLSOptions{
				CACertificate: caCert,
			},
		},
	}

	transport, err := buildTransport(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.TLSClientConfig.RootCAs == nil {
		t.Error("RootCAs should not be nil when CA cert provided")
	}
}

func TestBuildTransport_InvalidCACert(t *testing.T) {
	s := &Settings{
		URL: &url.URL{Scheme: "https", Host: "localhost:8080"},
		Opts: httpclient.Options{
			TLS: &httpclient.TLSOptions{
				CACertificate: "not-a-valid-pem",
			},
		},
	}

	_, err := buildTransport(s)
	if err == nil {
		t.Fatal("expected error for invalid CA certificate")
	}
}

func TestBuildTransport_ClientCertWithoutKey(t *testing.T) {
	s := &Settings{
		URL: &url.URL{Scheme: "https", Host: "localhost:8080"},
		Opts: httpclient.Options{
			TLS: &httpclient.TLSOptions{
				ClientCertificate: "some-cert",
				ClientKey:         "",
			},
		},
	}

	_, err := buildTransport(s)
	if err == nil {
		t.Fatal("expected error when client cert provided without key")
	}
}
