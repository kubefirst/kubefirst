package k3d

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubefirst/kubefirst/configs"
)

// DetokenizeCivoGithubGitops - Translate tokens by values on a given path
func DetokenizeK3dGithubGitops(path string, tokens *K3dTokenValues) error {

	err := filepath.Walk(path, detokenizeK3dGithubdGitops(path, tokens))
	if err != nil {
		return err
	}

	return nil
}

func detokenizeK3dGithubdMetaphor(path string, tokens *K3dTokenValues) filepath.WalkFunc {
	// todo implement
	return nil
}

func detokenizeK3dGithubdGitops(path string, tokens *K3dTokenValues) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !!fi.IsDir() {
			return nil
		}

		// var matched bool
		matched, err := filepath.Match("*", fi.Name())
		if matched {
			read, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// todo reduce to terraform tokens by moving to helm chart?
			newContents := string(read)
			newContents = strings.Replace(newContents, "<ALERTS_EMAIL>", "your@email.com", -1) //
			newContents = strings.Replace(newContents, "<ARGO_CD_INGRESS_URL>", tokens.ArgocdIngressURL, -1)
			newContents = strings.Replace(newContents, "<ARGO_WORKFLOWS_INGRESS_URL>", tokens.ArgoWorkflowsIngressURL, -1)
			newContents = strings.Replace(newContents, "<ATLANTIS_ALLOW_LIST>", tokens.AtlantisAllowList, -1)
			newContents = strings.Replace(newContents, "<ATLANTIS_INGRESS_URL>", tokens.AtlantisIngressURL, -1)
			newContents = strings.Replace(newContents, "<CLUSTER_NAME>", tokens.ClusterName, -1)
			newContents = strings.Replace(newContents, "<DOMAIN_NAME>", DomainName, -1)
			newContents = strings.Replace(newContents, "<KUBEFIRST_VERSION>", configs.K1Version, -1)
			newContents = strings.Replace(newContents, "<METAPHOR_DEVELPOMENT_INGRESS_URL>", tokens.MetaphorDevelopmentIngressURL, -1)
			newContents = strings.Replace(newContents, "<METAPHOR_STAGING_INGRESS_URL>", tokens.MetaphorStagingIngressURL, -1)
			newContents = strings.Replace(newContents, "<METAPHOR_PRODUCTION_INGRESS_URL>", tokens.MetaphorProductionIngressURL, -1)
			newContents = strings.Replace(newContents, "<GITHUB_HOST>", tokens.GithubHost, -1)
			newContents = strings.Replace(newContents, "<GITHUB_OWNER>", tokens.GithubOwner, -1)
			newContents = strings.Replace(newContents, "<GITHUB_USER>", tokens.GithubUser, -1)
			newContents = strings.Replace(newContents, "<GITOPS_REPO_GIT_URL>", tokens.GitopsRepoGitURL, -1)
			newContents = strings.Replace(newContents, "<NGROK_HOST>", tokens.NgrokHost, -1)
			newContents = strings.Replace(newContents, "<VAULT_INGRESS_URL>", tokens.VaultIngressURL, -1)

			err = ioutil.WriteFile(path, []byte(newContents), 0)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
