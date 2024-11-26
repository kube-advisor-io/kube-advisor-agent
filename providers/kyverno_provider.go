package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/commands/apply"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/processor"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/report"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/store"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/common"
	policyreportv1alpha2 "github.com/kyverno/kyverno/api/policyreport/v1alpha2"

	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

var (
	cmd apply.ApplyCommandConfig
)

type KyvernoPoliciesProvider struct {
	// kyvernoPoliciesResourceProvider *ResourceProvider[kyvernov1.ClusterPolicy]
	dynamicClient                   *dynamic.DynamicClient
	kubeConfig *restclient.Config
	policies *unstructured.UnstructuredList
}

func NewKyvernoPoliciesProvider(dynamicClient *dynamic.DynamicClient, kubeConfig *restclient.Config, config config.Config) *KyvernoPoliciesProvider {
	provider := &KyvernoPoliciesProvider{}
	provider.dynamicClient = dynamicClient
	provider.kubeConfig = kubeConfig
	// provider.kyvernoPoliciesResourceProvider = GetKyvernoResourceProvider(dynamicClient, config.IgnoredNamespaces)
	provider.policies, _ = dynamicClient.Resource(
		schema.GroupVersionResource{Group: "kyverno.io", Resource: "clusterpolicies", Version: "v1"}).
		List(context.Background(), metav1.ListOptions{TimeoutSeconds: &listTimeout})
	// if err != nil{
	// 	log.Error(err)
	// }
	log.Infof("Found cluster policies: %v", provider.policies)
	return provider
}

func (kpp *KyvernoPoliciesProvider) CheckPolicies() []policyreportv1alpha2.ClusterPolicyReport {
	parsedPolicies := []kyvernov1.PolicyInterface{}
	for _, policy := range kpp.policies.Items {
		originalPolicyJson, _ := json.Marshal(policy)
		log.Infof("original policy: %s", originalPolicyJson)
		outPolicy := &kyvernov1.ClusterPolicy{}
		runtime.DefaultUnstructuredConverter.FromUnstructured(policy.Object, &outPolicy)
		parsedPolicyJson, _ := json.Marshal(outPolicy)
		log.Infof("parsed policy: %s", parsedPolicyJson)

		parsedPolicies = append(parsedPolicies, outPolicy)
	}

	dClient, err := kpp.initStoreAndClusterClient()
	if err != nil {
		log.Error(err)
	}
	log.Infof("dClient: %v", dClient)

	resources, err := kpp.loadResources(os.Stdout, parsedPolicies, dClient)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Resources: %v", resources)

	var store store.Store
	store.SetLocal(true)
	store.AllowApiCall(true)
	store.SetRegistryAccess(true)
	var engineResponses []engineapi.EngineResponse
	for _, resource := range resources {
		var rc processor.ResultCounts
		processor := processor.PolicyProcessor{
			Store:                &store,
			Policies:             parsedPolicies,
			Resource:             *resource,
			PolicyExceptions:     nil,
			MutateLogPath:        "",
			MutateLogPathIsDir:   false,
			Variables:            nil,
			UserInfo:             nil,
			PolicyReport:         true,
			NamespaceSelectorMap: nil,
			Stdin:                false,
			Rc:                   &rc,
			PrintPatchResource:   true,
			Client:               dClient,
			AuditWarn:            false,
			Subresources:         nil,
			Out:                  os.Stdout,
		}
		ers, err := processor.ApplyPoliciesOnResource()
		log.Infof("Engine response: %v", ers)
		if err != nil {
			log.Errorf("Error while applying policies %v", err)
		}
		engineResponses = append(engineResponses, ers...)
	}

	clusterReport, _ := report.ComputePolicyReports(false, engineResponses...)
	return clusterReport
}

func (kpp *KyvernoPoliciesProvider) loadResources(out io.Writer, policies []kyvernov1.PolicyInterface, dClient dclient.Interface) ([]*unstructured.Unstructured, error) {
	resources, err := common.GetResources(out, policies, nil, nil, dClient, true, "", true)
	if err != nil {
		return resources, fmt.Errorf("failed to load resources (%w)", err)
	}
	return resources, nil
}

func (kpp *KyvernoPoliciesProvider) initStoreAndClusterClient() (dclient.Interface, error) {
	var err error
	var dClient dclient.Interface
	kubeClient, err := kubernetes.NewForConfig(kpp.kubeConfig)
	if err != nil {
		return nil, err
	}
	dClient, err = dclient.NewClient(context.Background(), kpp.dynamicClient, kubeClient, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	return dClient, err
}
