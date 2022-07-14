package cmd

import (
	"fmt"
	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/internal/progressPrinter"
	"github.com/kubefirst/kubefirst/internal/telemetry"
	"github.com/kubefirst/kubefirst/pkg"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

// todo: move it to internals/ArgoCD
func setArgocdCreds(dryRun bool) {
	if dryRun {
		log.Printf("[#99] Dry-run mode, setArgocdCreds skipped.")
		viper.Set("argocd.admin.password", "dry-run-not-real-pwd")
		viper.Set("argocd.admin.username", "dry-run-not-admin")
		viper.WriteConfig()
		return
	}

	cfg := configs.ReadConfig()
	config, err := clientcmd.BuildConfigFromFlags("", cfg.KubeConfigPath)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	argocdSecretClient = clientset.CoreV1().Secrets("argocd")

	argocdPassword := getSecretValue(argocdSecretClient, "argocd-initial-admin-secret", "password")

	viper.Set("argocd.admin.password", argocdPassword)
	viper.Set("argocd.admin.username", "admin")
	viper.WriteConfig()
}

func sendStartedInstallTelemetry(dryRun bool) {
	metricName := "kubefirst.mgmt_cluster_install.started"
	if !dryRun {
		telemetry.SendTelemetry(viper.GetString("aws.hostedzonename"), metricName)
	} else {
		log.Printf("[#99] Dry-run mode, telemetry skipped:  %s", metricName)
	}
}

func sendCompleteInstallTelemetry(dryRun bool) {
	metricName := "kubefirst.mgmt_cluster_install.completed"
	if !dryRun {
		telemetry.SendTelemetry(viper.GetString("aws.hostedzonename"), metricName)
	} else {
		log.Printf("[#99] Dry-run mode, telemetry skipped:  %s", metricName)
	}
}

func waitArgoCDToBeReady(dryRun bool) {
	if dryRun {
		log.Printf("[#99] Dry-run mode, waitArgoCDToBeReady skipped.")
		return
	}
	config := configs.ReadConfig()
	x := 50
	for i := 0; i < x; i++ {
		_, _, err := pkg.ExecShellReturnStrings(config.KubectlClientPath, "--kubeconfig", config.KubeConfigPath, "get", "namespace/argocd")
		if err != nil {
			log.Println("Waiting argocd to be born")
			time.Sleep(10 * time.Second)
		} else {
			log.Println("argocd namespace found, continuing")
			time.Sleep(5 * time.Second)
			break
		}
	}
	for i := 0; i < x; i++ {
		_, _, err := pkg.ExecShellReturnStrings(config.KubectlClientPath, "--kubeconfig", config.KubeConfigPath, "get", "pods", "-l", "app.kubernetes.io/name=argocd-server")
		if err != nil {
			log.Println("Waiting for argocd pods to create, checking in 10 seconds")
			time.Sleep(10 * time.Second)
		} else {
			log.Println("argocd pods found, continuing")
			time.Sleep(15 * time.Second)
			break
		}
	}
}

func waitVaultToBeInitialized(dryRun bool) {
	if dryRun {
		log.Printf("[#99] Dry-run mode, waitVaultToBeInitialized skipped.")
		return
	}
	config := configs.ReadConfig()
	x := 50
	for i := 0; i < x; i++ {
		_, _, err := pkg.ExecShellReturnStrings(config.KubectlClientPath, "--kubeconfig", config.KubeConfigPath, "get", "namespace/vault")
		if err != nil {
			log.Println("Waiting vault to be born")
			time.Sleep(10 * time.Second)
		} else {
			log.Println("vault namespace found, continuing")
			time.Sleep(25 * time.Second)
			break
		}
	}
	x = 50
	for i := 0; i < x; i++ {
		_, _, err := pkg.ExecShellReturnStrings(config.KubectlClientPath, "--kubeconfig", config.KubeConfigPath, "-n", "vault", "get", "pods", "-l", "vault-initialized=true")
		if err != nil {
			log.Println("Waiting vault pods to create")
			time.Sleep(10 * time.Second)
		} else {
			log.Println("vault pods found, continuing")
			time.Sleep(15 * time.Second)
			break
		}
	}
}

func waitGitlabToBeReady(dryRun bool) {
	if dryRun {
		log.Printf("[#99] Dry-run mode, waitVaultToBeInitialized skipped.")
		return
	}
	config := configs.ReadConfig()
	x := 50
	for i := 0; i < x; i++ {
		_, _, err := pkg.ExecShellReturnStrings(config.KubectlClientPath, "--kubeconfig", config.KubeConfigPath, "get", "namespace/gitlab")
		if err != nil {
			log.Println("Waiting gitlab namespace to be born")
			time.Sleep(10 * time.Second)
		} else {
			log.Println("gitlab namespace found, continuing")
			time.Sleep(5 * time.Second)
			break
		}
	}
	x = 50
	for i := 0; i < x; i++ {
		_, _, err := pkg.ExecShellReturnStrings(config.KubectlClientPath, "--kubeconfig", config.KubeConfigPath, "-n", "gitlab", "get", "pods", "-l", "app=webservice")
		if err != nil {
			log.Println("Waiting gitlab pods to be born")
			time.Sleep(10 * time.Second)
		} else {
			log.Println("gitlab pods found, continuing")
			time.Sleep(15 * time.Second)
			break
		}
	}

}

//Notify user in the STOUT and also logfile
func informUser(message string) {
	log.Println(message)
	progressPrinter.LogMessage(fmt.Sprintf("- %s", message))
}
