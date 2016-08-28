package jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/spf13/viper"
)

type Catalog struct {
	Repositories []string `json:"repositories"`
}

func SynchronizeCatalog() error {
	// Request catalog from registry
	registryAddress := viper.GetString("registry.address")
	registryPort := viper.GetInt("registry.port")
	log.Printf("Querying catalog from %s:%d...\n", registryAddress, registryPort)
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s:%d/v2/_catalog", registryAddress, registryPort), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Read response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse catalog
	var catalog Catalog
	err = json.Unmarshal(data, &catalog)
	if err != nil {
		return err
	}

	log.Printf("Found %d entries in catalog\n", len(catalog.Repositories))

	// Get all namespaces from database
	var namespaces []models.Namespace
	db.Instance().Find(&namespaces)

	// Prepare a map to check existing repositories
	existingNamespaces := make(map[string]struct{})
	for _, namespace := range namespaces {
		existingNamespaces[namespace.Name] = struct{}{}
	}

	// Parse namespace and insert if doesn't exist
	for _, repository := range catalog.Repositories {
		repSplit := strings.SplitN(repository, "/", 2)
		name := repSplit[0]
		if _, ok := existingNamespaces[name]; !ok {
			log.Printf("Creating %s namespace\n", name)
			db.Instance().Create(&models.Namespace{Name: name})
		}
	}

	log.Println("All done!")

	return nil
}
