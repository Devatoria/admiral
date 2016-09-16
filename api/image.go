package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Devatoria/admiral/auth"
	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"
	"github.com/Devatoria/admiral/token"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// getImages returns the images contained in the user namespace
func getImages(c *gin.Context) {
	// Get user
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Get user namespace
	namespace := models.GetNamespaceByName(user.Username)
	if namespace.ID == 0 {
		panic(fmt.Sprintf("User %s has no namespace", user.Username))
	}

	// Get images
	c.JSON(http.StatusOK, namespace.Images)
}

// deleteImage deletes the given image and all its tags from the registry and the database
func deleteImage(c *gin.Context) {
	image, ok := c.Keys["image"].(models.Image)
	if !ok {
		panic("Unable to get image from context")
	}

	// Request registry to remove all tags (and remove from database)
	client := http.Client{}
	registryAddress := fmt.Sprintf("%s:%d", viper.GetString("registry.address"), viper.GetInt("registry.port"))
	for _, tag := range image.Tags {
		// Build token for tag reference
		t := token.NewToken("registry", "admiral", []token.ClaimsAccess{
			token.ClaimsAccess{
				Type:    "repository",
				Name:    image.Name,
				Actions: []string{"pull"},
			},
		})
		tString, err := token.SignToken(t)
		if err != nil {
			panic(err)
		}

		// Prepare request
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/%s/manifests/%s", registryAddress, image.Name, tag.Name), nil)
		if err != nil {
			panic(err)
		}

		// Get tag reference
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tString))
		req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}

		digest := resp.Header.Get("Docker-Content-Digest")
		var m schema2.Manifest
		if err := json.Unmarshal(data, &m); err != nil {
			panic(err)
		}

		// Remove each layer (blob)
		for _, layer := range m.Layers {
			// Build token for delete
			t = token.NewToken("registry", "admiral", []token.ClaimsAccess{
				token.ClaimsAccess{
					Type:    "repository",
					Name:    image.Name,
					Actions: []string{"*"},
				},
			})
			tString, err = token.SignToken(t)
			if err != nil {
				panic(err)
			}

			// Prepare request
			req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/v2/%s/blobs/%s", registryAddress, image.Name, layer.Digest), nil)
			if err != nil {
				panic(err)
			}

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tString))
			resp, err = client.Do(req)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Status is %d for layer %s\n", resp.StatusCode, layer.Digest)
		}

		// Build token for delete
		t = token.NewToken("registry", "admiral", []token.ClaimsAccess{
			token.ClaimsAccess{
				Type:    "repository",
				Name:    image.Name,
				Actions: []string{"*"},
			},
		})
		tString, err = token.SignToken(t)
		if err != nil {
			panic(err)
		}

		// Prepare request
		req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/v2/%s/manifests/%s", registryAddress, image.Name, digest), nil)
		if err != nil {
			panic(err)
		}

		// Do request
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tString))
		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Status is %d for tag %s\n", resp.StatusCode, tag.Name)
		if resp.StatusCode > 299 && resp.StatusCode != 404 {
			continue
		}

		// Delete tag from database
		db.Instance().Delete(&tag)
	}

	// Remove image from database
	db.Instance().Delete(&image)

	c.Status(http.StatusOK)
}
