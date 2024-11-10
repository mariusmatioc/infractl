package utils

import (
	"os/exec"
)

/*
aws eks update-kubeconfig --region <region> --name <cluster_name>
once. Afterwards, the cluster is saved in kubectl. Then, to list all clusters in kubectl, run
kubectl config get-contexts
and to switch between clusters in kubectl, run
kubectl config use-context <NAME>
*/

func RunKubectlCommand(args []string) (output string, err error) {
	out, err := exec.Command("kubectl", args...).Output()
	if err == nil {
		output = string(out)
	}
	return
}
