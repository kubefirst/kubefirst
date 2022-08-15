package downloadManager

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/pkg"
)

func DownloadTools(config *configs.Config) error {

	toolsDirPath := fmt.Sprintf("%s/tools", config.K1FolderPath)

	// create folder if it doesn't exist
	if _, err := os.Stat(toolsDirPath); errors.Is(err, fs.ErrNotExist) {
		err = os.Mkdir(toolsDirPath, 0777)
		if err != nil {
			return err
		}
	}

	kubectlDownloadUrl := fmt.Sprintf(
		"https://dl.k8s.io/release/%s/bin/%s/%s/kubectl",
		config.KubectlVersion,
		config.LocalOs,
		config.LocalArchitecture,
	)

	err := downloadFile(config.KubectlClientPath, kubectlDownloadUrl)
	if err != nil {
		return err
	}

	err = os.Chmod(config.KubectlClientPath, 0755)
	if err != nil {
		return err
	}

	// todo: this kubeconfig is not available to us until we have run the terraform in base/
	err = os.Setenv("KUBECONFIG", config.KubeConfigPath)
	if err != nil {
		return err
	}

	log.Println("going to print the kubeconfig env in runtime", os.Getenv("KUBECONFIG"))

	kubectlStdOut, kubectlStdErr, errKubectl := pkg.ExecShellReturnStrings(config.KubectlClientPath, "version", "--client", "--short")
	log.Printf("-> kubectl version:\n\t%s\n\t%s\n", kubectlStdOut, kubectlStdErr)
	if errKubectl != nil {
		log.Panicf("failed to call kubectlVersionCmd.Run(): %v", err)
	}

	// todo: adopt latest helmVersion := "v3.9.0"
	terraformVersion := config.TerraformVersion

	terraformDownloadUrl := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip",
		terraformVersion,
		terraformVersion,
		config.LocalOs,
		config.LocalArchitecture,
	)

	terraformDownloadZipPath := fmt.Sprintf("%s/tools/terraform.zip", config.K1FolderPath)
	err = downloadFile(terraformDownloadZipPath, terraformDownloadUrl)
	if err != nil {
		log.Println("error reading terraform file")
		return err
	}

	unzipDirectory := fmt.Sprintf("%s/tools", config.K1FolderPath)
	unzip(terraformDownloadZipPath, unzipDirectory)

	err = os.Chmod(unzipDirectory, 0777)
	if err != nil {
		return err
	}

	err = os.Chmod(fmt.Sprintf("%s/terraform", unzipDirectory), 0755)
	if err != nil {
		return err
	}
	os.RemoveAll(fmt.Sprintf("%s/terraform.zip", toolsDirPath))

	helmVersion := config.HelmVersion
	helmDownloadUrl := fmt.Sprintf(
		"https://get.helm.sh/helm-%s-%s-%s.tar.gz",
		helmVersion,
		config.LocalOs,
		config.LocalArchitecture,
	)

	helmDownloadTarGzPath := fmt.Sprintf("%s/tools/helm.tar.gz", config.K1FolderPath)
	err = downloadFile(helmDownloadTarGzPath, helmDownloadUrl)
	if err != nil {
		return err
	}

	helmTarDownload, err := os.Open(helmDownloadTarGzPath)
	if err != nil {
		log.Panicf("could not read helm download content")
	}

	extractFileFromTarGz(
		helmTarDownload,
		fmt.Sprintf("%s-%s/helm", config.LocalOs, config.LocalArchitecture),
		config.HelmClientPath,
	)
	err = os.Chmod(config.HelmClientPath, 0755)
	if err != nil {
		return err
	}

	helmStdOut, helmStdErr, errHelm := pkg.ExecShellReturnStrings(
		config.HelmClientPath,
		"version",
		"--client",
		"--short",
	)

	log.Printf("-> kubectl version:\n\t%s\n\t%s\n", helmStdOut, helmStdErr)
	// currently argocd init values is generated by kubefirst ssh
	// todo helm install argocd --create-namespace --wait --values ~/.kubefirst/argocd-init-values.yaml argo/argo-cd
	if errHelm != nil {
		log.Panicf("error executing helm version command: %v", err)
	}

	return nil
}

// downloadFile Downloads a file from the "url" parameter, localFilename is the file destination in the local machine.
func downloadFile(localFilename string, url string) error {
	// create local file
	out, err := os.Create(localFilename)
	if err != nil {
		return err
	}
	defer out.Close()

	// get data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// writer the body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
func extractFileFromTarGz(gzipStream io.Reader, tarAddress string, targetFilePath string) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Panicf("extractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicf("extractTarGz: Next() failed: %s", err.Error())
		}
		log.Println(header.Name)
		if header.Name == tarAddress {
			switch header.Typeflag {
			case tar.TypeReg:
				outFile, err := os.Create(targetFilePath)
				if err != nil {
					log.Panicf("extractTarGz: Create() failed: %s", err.Error())
				}
				if _, err := io.Copy(outFile, tarReader); err != nil {
					log.Panicf("extractTarGz: Copy() failed: %s", err.Error())
				}
				outFile.Close()

			default:
				log.Println(
					"extractTarGz: uknown type: %s in %s",
					header.Typeflag,
					header.Name)
			}

		}
	}
}

func unzip(zipFilepath string, unzipDirectory string) {
	dst := unzipDirectory
	archive, err := zip.OpenReader(zipFilepath)
	if err != nil {
		log.Panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		log.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			log.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			log.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			log.Panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			log.Panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
