package certs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// CertificateProviderType is the enum for certificate providers.
type CertificateProviderType int

const (
	// ServerProvider uses server.pem/server-key.pem.
	ServerProvider CertificateProviderType = iota
	// ClientProvider uses client.pem/client-key.pem.
	ClientProvider
	// CertProvider uses cert.pem/key.pem.
	CertProvider
)

// String returns the string representation of CertificateProviderType.
func (c CertificateProviderType) String() string {
	switch c {
	case ServerProvider:
		return "ServerProvider"
	case ClientProvider:
		return "ClientProvider"
	case CertProvider:
		return "CertProvider"
	}

	return ""
}

// CertProviderFromString returns a CertificateProviderType from a supplied string.
func CertProviderFromString(in string) CertificateProviderType {
	switch strings.ToLower(in) {
	case "serverprovider", "server":
		return ServerProvider
	case "clientprovider", "client":
		return ClientProvider
	case "certprovider", "cert":
		return CertProvider
	}

	return ClientProvider
}

// ErrNoValidCertificates returned when no valid certificates are found in ca.pem.
var ErrNoValidCertificates = errors.New("no valid certificates present")

// CertificateProvider is an interface to a provider for certificates used with gRPC server and clients.
type CertificateProvider interface {
	IdentityCert() tls.Certificate
	CAPool() *x509.CertPool
	ServerOption() grpc.ServerOption
	DialOption(serverName string) grpc.DialOption
}

// FileCertificateProvider uses files for the source of certificates and keys.
type FileCertificateProvider struct {
	certDir      string
	identityCert tls.Certificate
	caPool       *x509.CertPool
}

// NewFileCertificateProvider returns a new FileCertificateProvider using certs from the specified directory
// optionally also can be used for gRPC clients by setting server to false.
func NewFileCertificateProvider(certDir string, providerType CertificateProviderType) (c CertificateProvider, err error) {
	f := &FileCertificateProvider{
		certDir: os.ExpandEnv(certDir),
	}

	switch providerType {
	case ServerProvider:
		f.identityCert, err = tls.LoadX509KeyPair(f.serverCertPath(), f.serverKeyPath())
	case ClientProvider:
		f.identityCert, err = tls.LoadX509KeyPair(f.clientCertPath(), f.clientKeyPath())
	case CertProvider:
		f.identityCert, err = tls.LoadX509KeyPair(f.certCertPath(), f.certKeyPath())
	}

	if err != nil {
		return nil, fmt.Errorf("failed loading certificates: %w", err)
	}

	// populate certs.IdentityCert.Leaf. This has already been parsed, but
	// intentionally discarded by LoadX509KeyPair, for some reason.
	f.identityCert.Leaf, err = x509.ParseCertificate(f.identityCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed loading certificates: %w", err)
	}

	ca, err := ioutil.ReadFile(f.caCertPath())
	if err != nil {
		return nil, fmt.Errorf("failed loading certificates: %w", err)
	}

	f.caPool = x509.NewCertPool()
	switch providerType {
	case ClientProvider, CertProvider:
		f.caPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed loading certificates: %w", err)
		}
	}

	if ok := f.caPool.AppendCertsFromPEM(ca); !ok {
		return nil, ErrNoValidCertificates
	}

	return f, nil
}

func (c *FileCertificateProvider) certCertPath() string {
	return path.Join(c.certDir, "cert.pem")
}

func (c *FileCertificateProvider) certKeyPath() string {
	return path.Join(c.certDir, "key.pem")
}

func (c *FileCertificateProvider) serverCertPath() string {
	return path.Join(c.certDir, "server.pem")
}

func (c *FileCertificateProvider) serverKeyPath() string {
	return path.Join(c.certDir, "server-key.pem")
}

func (c *FileCertificateProvider) clientCertPath() string {
	return path.Join(c.certDir, "client.pem")
}

func (c *FileCertificateProvider) clientKeyPath() string {
	return path.Join(c.certDir, "client-key.pem")
}

func (c *FileCertificateProvider) caCertPath() string {
	return path.Join(c.certDir, "ca.pem")
}

// IdentityCert returns the Identity Certificate used for the connection.
func (c *FileCertificateProvider) IdentityCert() tls.Certificate {
	return c.identityCert
}

// CAPool returns the CA Pool for the connection.
func (c *FileCertificateProvider) CAPool() *x509.CertPool {
	return c.caPool
}

// ServerOption returns the grpc.ServerOption for use with a new gRPC server.
func (c *FileCertificateProvider) ServerOption() grpc.ServerOption {
	creds := credentials.NewTLS(&tls.Config{
		ClientCAs:    c.CAPool(),
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{c.IdentityCert()},
		MinVersion:   tls.VersionTLS13,
	})

	return grpc.Creds(creds)
}

// DialOption returns the grpc.DialOption used with a gRPC client.
func (c *FileCertificateProvider) DialOption(serverName string) grpc.DialOption {
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   serverName,
		RootCAs:      c.CAPool(),
		Certificates: []tls.Certificate{c.IdentityCert()},
		MinVersion:   tls.VersionTLS13,
	})

	return grpc.WithTransportCredentials(creds)
}
