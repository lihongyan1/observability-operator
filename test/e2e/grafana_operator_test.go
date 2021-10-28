package e2e

import (
	"testing"

	"github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	grafanaDeploymentName = "grafana-deployment"
	operatorNamespace     = "monitoring-stack-operator"
)

func TestGrafanaOperatorForResourcesInOwnNamespace(t *testing.T) {
	resources := []client.Object{
		newGrafana(operatorNamespace),
	}
	defer deleteResources(resources...)

	ts := []testCase{
		{
			name: "Operator should create Grafana Operator CRDs",
			scenario: func(t *testing.T) {
				f.AssertResourceEventuallyExists("grafanadashboards.integreatly.org", "", &apiextensionsv1.CustomResourceDefinition{})(t)
				f.AssertResourceEventuallyExists("grafanadatasources.integreatly.org", "", &apiextensionsv1.CustomResourceDefinition{})(t)
				f.AssertResourceEventuallyExists("grafananotificationchannels.integreatly.org", "", &apiextensionsv1.CustomResourceDefinition{})(t)
				f.AssertResourceEventuallyExists("grafanas.integreatly.org", "", &apiextensionsv1.CustomResourceDefinition{})(t)
			},
		},
		{
			name:     "Create grafana operator resources",
			scenario: createResources(resources...),
		},
		{
			name: "Operator should reconcile resources in its own namespace",
			scenario: func(t *testing.T) {
				f.AssertResourceEventuallyExists(grafanaDeploymentName, operatorNamespace, &appsv1.Deployment{})(t)
			},
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, tc.scenario)
	}
}

func TestGrafanaOperatorForResourcesOutsideOfItsOwnNamespace(t *testing.T) {
	resources := []client.Object{
		newGrafana(e2eTestNamespace),
	}
	defer deleteResources(resources...)

	ts := []testCase{
		{
			name:     "Create grafana resource",
			scenario: createResources(resources...),
		},
		{
			name: "Operator should not reconcile resources outside of its own namespace",
			scenario: func(t *testing.T) {
				f.AssertResourceNeverExists(grafanaDeploymentName, e2eTestNamespace, &appsv1.Deployment{})(t)
			},
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, tc.scenario)
	}
}

func newGrafana(namespace string) *v1alpha1.Grafana {
	return &v1alpha1.Grafana{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "integreatly.org/v1alpha1",
			Kind:       "Grafana",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafana",
			Namespace: namespace,
		},
	}
}