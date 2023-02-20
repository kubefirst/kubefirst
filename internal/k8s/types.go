package k8s

// podSessionOptions provides a struct to assign parameters to an exec session
type PodSessionOptions struct {
	Command    []string
	Namespace  string
	PodName    string
	TtyEnabled bool
}
