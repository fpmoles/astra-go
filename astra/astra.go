package astra

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"io/ioutil"
	"path/filepath"
	"time"
)

const filepathStructure = "%s/%s"

type ClusterConnection struct {
	Session *gocql.Session
}

type ClusterConfig struct {
	databaseHost     string
	databasePort     string
	databaseUsername string
	databasePassword string
	filepath         string
	Timeout          *time.Duration
}

//Creates a new database connection to DataStax Astra using GoCQL
func NewClusterConnection(config *ClusterConfig) (*ClusterConnection, error) {
	cluster := gocql.NewCluster(config.databaseHost)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.databaseUsername,
		Password: config.databasePassword,
	}
	cluster.Hosts = []string{config.databaseHost + ":" + config.databasePort}

	certPath, err := filepath.Abs(fmt.Sprintf(filepathStructure, config.filepath, "cert"))
	if err != nil {
		return nil, errors.New("cannot load local cert file from filepath")
	}
	keyPath, err := filepath.Abs(fmt.Sprintf(filepathStructure, config.filepath, "key"))
	if err != nil {
		return nil, errors.New("cannot load local key file from filepath")
	}
	caPath, err := filepath.Abs(fmt.Sprintf(filepathStructure, config.filepath, "ca.crt"))
	if err != nil {
		return nil, errors.New("cannot load local ca file from filepath")
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errors.New("error creating cert from cert and key")
	}
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, errors.New("error reading ca cert from cert path")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	cluster.SslOpts = &gocql.SslOptions{
		Config:                 tlsConfig,
		EnableHostVerification: false,
	}
	cluster.Consistency = gocql.LocalQuorum
	if config.Timeout != nil {
		cluster.Timeout = *config.Timeout
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, errors.New("error creating session")
	}
	return &ClusterConnection{
		Session: session,
	}, nil

}

//Creates a new configuration for GoCQL for DataStax Astra. Filepath is the location of the MTLS certificates from the secure connect bundle
func NewClusterConfig(host string, port string, username string, password string, filepath string) *ClusterConfig {
	return &ClusterConfig{
		databaseHost:     host,
		databasePort:     port,
		databaseUsername: username,
		databasePassword: password,
		filepath:         filepath,
	}
}
