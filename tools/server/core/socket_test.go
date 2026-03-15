package core_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/tools/server/core"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// selfSignedCert writes a self-signed cert/key pair for "localhost" into a
// temp dir, for testing ListenAndServe's TLS branch.
func selfSignedCert(t *testing.T) (certPath, keyPath string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	assert.NoError(t, err)

	dir := t.TempDir()
	certPath = dir + "/cert.pem"
	keyPath = dir + "/key.pem"

	certOut, err := os.Create(certPath)
	assert.NoError(t, err)
	assert.NoError(t, pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}))
	assert.NoError(t, certOut.Close())

	keyOut, err := os.Create(keyPath)
	assert.NoError(t, err)
	assert.NoError(t, pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}))
	assert.NoError(t, keyOut.Close())
	return certPath, keyPath
}

func TestListenAndServeTLS(t *testing.T) {
	certPath, keyPath := selfSignedCert(t)
	cfg := core.Config{
		Addr: ":53455",
		SSL:  core.SSL{Cert: certPath, Key: keyPath},
	}
	srv := cfg.NewServer()
	expected := []byte("tls ok")
	srv.Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	})

	closed := make(chan error, 1)
	go func() { closed <- srv.ListenAndServe() }()

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	var resp *http.Response
	err := timeout.After(500, func() {
		for {
			var gerr error
			resp, gerr = client.Get("https://localhost" + cfg.Addr + "/test")
			if gerr == nil {
				return
			}
			time.Sleep(time.Millisecond)
		}
	})
	assert.NoError(t, err)
	got, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	assert.NoError(t, srv.Close())
	err = <-closed
	assert.Equal(t, http.ErrServerClosed, err)
}

// syncBuffer is a bytes.Buffer safe to read from one goroutine while a CLI
// loop writes from another.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func (s *syncBuffer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buf.Reset()
}


func TestServerSocketNotConfigured(t *testing.T) {
	cfg := core.Config{Addr: ":53458"}
	srv := cfg.NewServer()
	assert.False(t, srv.AwaitSocket())
}


func TestRunStdIO(t *testing.T) {
	cfg := core.Config{Addr: ":53459"}
	srv := cfg.NewServer()
	srv.CliHandler = cmdrFactory(srv)

	oldIn, oldOut := cli.StdIn, cli.StdOut
	defer func() { cli.StdIn, cli.StdOut = oldIn, oldOut }()
	pr, pw := io.Pipe()
	cli.StdIn = pr
	buf := &syncBuffer{}
	cli.StdOut = buf

	// RunStdIO blocks running the CLI over StdIn/StdOut, same as a real
	// admin session; it isn't meant to return, so this leaves the
	// goroutine running for the rest of the test binary's life.
	go srv.RunStdIO()

	err := timeout.After(200, func() {
		for !strings.Contains(buf.String(), "> ") {
			time.Sleep(time.Millisecond)
		}
	})
	assert.NoError(t, err)

	pw.Write([]byte("help\n"))
	err = timeout.After(200, func() {
		for !strings.Contains(buf.String(), "h, help") {
			time.Sleep(time.Millisecond)
		}
	})
	assert.NoError(t, err)
}

func TestServerSocket(t *testing.T) {
	dir := t.TempDir()
	cfg := core.Config{
		Addr:   ":53457",
		Socket: dir + "/admin.sock",
	}
	srv := cfg.NewServer()
	srv.CliHandler = cmdrFactory(srv)

	closed := make(chan bool)
	go func() {
		srv.Run()
		closed <- true
	}()

	assert.True(t, srv.AwaitSocket())

	assert.NoError(t, srv.Close())
	err := timeout.After(100, closed)
	assert.NoError(t, err)
}
