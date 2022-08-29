package gitlab_test

import (
	"fmt"
	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/pkg"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

// this is called when GitLab should be up and running
func TestApplyGitlabTerraform(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	config := configs.ReadConfig()
	err := pkg.SetupViper(config)
	if err != nil {
		t.Error(err)
	}

	argoURL := fmt.Sprintf("https://gitlab.%s", viper.GetString("aws.hostedzonename"))

	req, err := http.NewRequest(http.MethodGet, argoURL, nil)
	if err != nil {
		t.Error(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted http status code 200, got %d", res.StatusCode)
	}
}
