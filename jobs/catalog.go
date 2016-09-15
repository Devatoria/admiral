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
	"github.com/Devatoria/admiral/token"

	"github.com/spf13/viper"
)

// Catalog represents the Docker Registry catalog JSON format
type Catalog struct {
	Repositories []string `json:"repositories"`
}

// Tags represents the Docker Registry tags list JSON format
type Tags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// SynchronizeCatalog parses the registry catalog to get namespaces, images and associated tags and inserts it database if needed
func SynchronizeCatalog(args []string) error {
	// Prepare bearer token
	t := token.NewToken("registry", "admiral", []token.ClaimsAccess{
		token.ClaimsAccess{
			Type:    "registry",
			Name:    "catalog",
			Actions: []string{"*"},
		},
	})
	tString, err := token.SignToken(t)
	if err != nil {
		panic(err)
	}

	// Request catalog from registry
	registryAddress := viper.GetString("registry.address")
	registryPort := viper.GetInt("registry.port")
	log.Printf("Querying catalog from %s:%d...\n", registryAddress, registryPort)
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s:%d/v2/_catalog", registryAddress, registryPort), nil)
	if err != nil {
		return err
	}

	req.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", tString)}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Read response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
		var image models.Image
		if len(repSplit) == 1 {
			if _, ok := existingImages[repository]; !ok {
				image = models.Image{Name: repository}
				log.Printf("Creating public image %s\n", image.Name)
				db.Instance().Create(&image)
				existingImages[image.Name] = image.ID
			}
		} else {
			// Create image if namespace already exists but not image
			if _, ok := existingNamespaces[repSplit[0]]; ok {
				if _, ok := existingImages[repository]; !ok {
					image = models.Image{Name: repository, NamespaceID: existingNamespaces[repSplit[0]]}
					log.Printf("Creating image %s\n", image.Name)
					db.Instance().Create(&image)
					existingImages[image.Name] = image.ID
				}
			}
		}

		// If the image has not been created (unexisting namespace, for example), we skip the tags
		if _, ok := existingImages[repository]; !ok {
			fmt.Printf("Skipping %s tags because namespace does not exist\n", repository)
			continue
		}

		// Prepare token
		tTag := token.NewToken("registry", "admiral", []token.ClaimsAccess{
			token.ClaimsAccess{
				Type:    "repository",
				Name:    repository,
				Actions: []string{"pull"},
			},
		})
		tTagString, err := token.SignToken(tTag)
		if err != nil {
			panic(err)
		}

		// Retrieve tags for current image
		reqTags, err := http.NewRequest("GET", fmt.Sprintf("%s:%d/v2/%s/tags/list", registryAddress, registryPort, repository), nil)
		if err != nil {
			log.Printf("Unable to create HTTP request: %v", err)
			continue
		}

		reqTags.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tTagString))
		respTags, err := client.Do(reqTags)
		if err != nil {
			log.Printf("Unable to do HTTP request: %v", err)
			continue
		}
		defer respTags.Body.Close()

		dataTags, err := ioutil.ReadAll(respTags.Body)
		if err != nil {
			log.Printf("Unable to read response body: %v", err)
			continue
		}

		// Parse tags
		var tags Tags
		err = json.Unmarshal(dataTags, &tags)
		if err != nil {
			log.Printf("Unable to parse response: %v", err)
			continue
		}

		// Create entities
		fmt.Printf("Found %d tags for image %s\n", len(tags.Tags), repository)
		for _, tag := range tags.Tags {
			tagEntity := models.Tag{
				Name:    tag,
				ImageID: existingImages[repository],
			}

			// Check if tag exists
			db.Instance().Where(&tagEntity).Find(&tagEntity)
			if tagEntity.ID != 0 {
				continue
			}

			fmt.Printf("Creating tag %s for image %s\n", tag, repository)
			db.Instance().Create(&tagEntity)
		}
	}

	log.Println("All done!")

	return nil
}
