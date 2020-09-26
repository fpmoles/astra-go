package astra

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNewClusterConfig(t *testing.T) {
	host := "foo"
	port := "8993"
	username := "ronswanson"
	password := "freedom"
	filepath := "/cfg"
	config := NewClusterConfig(host, port, username, password, filepath)
	assert.Equal(t, config.databaseHost, host, "host")
	assert.Equal(t, config.databasePort, port, "port")
	assert.Equal(t, config.databaseUsername, username, "username")
	assert.Equal(t, config.databasePassword, password, "password")
	assert.Equal(t, config.filepath, filepath, "filepath")
}

func TestNewClusterConnection(t *testing.T) {
	err := godotenv.Load("../astra.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	filepath := os.Getenv("CERTPATH")

	config := NewClusterConfig(host, port, username, password, filepath)
	clusterConnection, _ := NewClusterConnection(config)
	iter := clusterConnection.Session.Query("SELECT keyspace_name FROM system_schema.keyspaces;").Iter()
	found := false
	var keyspace_name string
	for iter.Scan(&keyspace_name) {
		if keyspace_name == "system" {
			found = true
		}
	}
	clusterConnection.Session.Close()
	assert.True(t, found)
}
