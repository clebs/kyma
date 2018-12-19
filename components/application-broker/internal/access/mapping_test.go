package access

import (
	"context"
	"testing"
	"time"

	"github.com/kyma-project/kyma/components/application-broker/pkg/apis/applicationconnector/v1alpha1"
	"github.com/kyma-project/kyma/components/application-broker/pkg/client/clientset/versioned/fake"
	informers "github.com/kyma-project/kyma/components/application-broker/pkg/client/informers/externalversions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestApplicationMappingService_IsApplicationEnabled(t *testing.T) {
	for tn, tc := range map[string]struct {
		givenMappings  []runtime.Object
		namespace      string
		name           string
		expectedResult bool
	}{
		"ApplicationMapping exists": {
			givenMappings: []runtime.Object{
				fixEnvMapping("prod", "ec"),
			},
			namespace:      "prod",
			name:           "ec",
			expectedResult: true,
		},
		"ApplicationMapping does not exists": {
			givenMappings: []runtime.Object{
				fixEnvMapping("prod", "ec"),
				fixEnvMapping("stage", "marketing"),
			},
			namespace:      "prod",
			name:           "marketing",
			expectedResult: false,
		},
	} {
		t.Run(tn, func(t *testing.T) {
			// GIVEN
			cs := fake.NewSimpleClientset(tc.givenMappings...)
			informerFactory := informers.NewSharedInformerFactory(cs, time.Hour)

			svc := NewApplicationMappingService(informerFactory.Applicationconnector().V1alpha1().ApplicationMappings().Lister())

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			informerFactory.Start(ctx.Done())
			informerFactory.WaitForCacheSync(ctx.Done())

			// WHEN
			result, err := svc.IsApplicationEnabled(tc.namespace, tc.name)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)

		})
	}
}

func fixEnvMapping(namespace, name string) *v1alpha1.ApplicationMapping {
	return &v1alpha1.ApplicationMapping{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
