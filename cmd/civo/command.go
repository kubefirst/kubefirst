package civo

import "github.com/spf13/cobra"

var (
	adminEmailFlag             string
	cloudRegionFlag            string
	clusterNameFlag            string
	clusterTypeFlag            string
	githubOwnerFlag            string
	gitopsTemplateURLFlag      string
	gitopsTemplateBranchFlag   string
	metaphorTemplateBranchFlag string
	metaphorTemplateURLFlag    string
	domainNameFlag             string
	kbotPasswordFlag           string
	silentModeFlag             bool
	useTelemetryFlag           bool
)

func NewCommand() *cobra.Command {

	civoCmd := &cobra.Command{
		Use:     "civo",
		Short:   "kubefirst civo installation",
		Long:    "kubefirst civo",
		PreRunE: validateCivo, // todo what should this function be called?
		RunE:    runCivo,
		// PostRunE: runPostCivo,
	}

	// todo review defaults and update descriptions
	civoCmd.Flags().StringVar(&adminEmailFlag, "admin-email", "", "email address for let's encrypt certificate notifications")
	civoCmd.Flags().StringVar(&cloudRegionFlag, "cloud-region", "NYC1", "the civo region to provision infrastructure in")
	civoCmd.Flags().StringVar(&clusterNameFlag, "cluster-name", "kubefirst", "the name of the cluster to create")
	civoCmd.Flags().StringVar(&clusterTypeFlag, "cluster-type", "mgmt", "the type of cluster to create (i.e. mgmt|workload)")
	civoCmd.Flags().StringVar(&domainNameFlag, "domain-name", "", "the Civo DNS Name to use for DNS records (i.e. your-domain.com|subdomain.your-domain.com)")
	civoCmd.Flags().StringVar(&githubOwnerFlag, "github-owner", "", "the GitHub owner of the new gitops and metaphor repositories")
	civoCmd.MarkFlagRequired("admin-email")
	civoCmd.MarkFlagRequired("domain-name")
	civoCmd.MarkFlagRequired("github-owner")
	civoCmd.Flags().StringVar(&gitopsTemplateBranchFlag, "gitops-template-branch", "main", "the branch to clone for the gitops-template repository")
	civoCmd.Flags().StringVar(&gitopsTemplateURLFlag, "gitops-template-url", "https://github.com/kubefirst/gitops-template.git", "the fully qualified url to the gitops-template repository to clone")
	civoCmd.Flags().StringVar(&kbotPasswordFlag, "kbot-password", "", "the default password to use for the kbot user")
	civoCmd.Flags().StringVar(&metaphorTemplateBranchFlag, "metaphor-template-branch", "main", "the branch to clone for the metaphor-template repository")
	civoCmd.Flags().StringVar(&metaphorTemplateURLFlag, "metaphor-template-url", "https://github.com/kubefirst/metaphor-frontend-template.git", "the fully qualified url to the metaphor-template repository to clone")

	civoCmd.Flags().BoolVar(&silentModeFlag, "silent-mode", false, "suppress output to the terminal")
	civoCmd.Flags().BoolVar(&useTelemetryFlag, "use-telemetry", true, "whether to emit telemetry")

	// on error, doesnt show helper/usage
	civoCmd.SilenceUsage = true

	// wire up new commands
	civoCmd.AddCommand(Destroy())

	return civoCmd
}

func Destroy() *cobra.Command {
	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy civo cloud",
		Long:  "todo",
		RunE:  destroyCivo,
	}

	return destroyCmd
}
