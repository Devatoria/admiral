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
	existingNamespaces := make(map[string]uint)
	for _, namespace := range namespaces {
		existingNamespaces[namespace.Name] = namespace.ID
	}

	// Get all public images from database
	var images []models.Image
	db.Instance().Find(&images)
	existingImages := make(map[string]uint)
	for _, image := range images {
		existingImages[image.Name] = image.ID
	}

	// Parse namespace and insert if doesn't exist
	for _, repository := range catalog.Repositories {
		repSplit := strings.SplitN(repository, "/", 2)

		// If public image (no namespace), just create image with null namespace
		// Else, ensure namespace exists (or create it), and then create image
		if len(repSplit) == 1 {
			if _, ok := existingImages[repSplit[0]]; !ok {
				image := models.Image{Name: repSplit[0]}
				log.Printf("Creating public image %s\n", image.Name)
				db.Instance().Create(&image)
				existingImages[image.Name] = image.ID
			}
		} else {
			if _, ok := existingNamespaces[repSplit[0]]; !ok {
				namespace := models.Namespace{Name: repSplit[0]}
				log.Printf("Creating namespace %s\n", namespace.Name)
				db.Instance().Create(&namespace)
				existingNamespaces[namespace.Name] = namespace.ID
			}

			if _, ok := existingImages[repository]; !ok {
				image := models.Image{Name: repository, NamespaceID: existingNamespaces[repSplit[0]]}
				log.Printf("Creating image %s\n", image.Name)
				db.Instance().Create(&image)
				existingImages[image.Name] = image.ID
			}
		}
	}

	log.Println("All done!")

	return nil
}
