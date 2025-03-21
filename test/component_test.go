package test

import (
	// "context"
	"testing"
	// "fmt"
	"strings"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	// awsHelper "github.com/cloudposse/test-helpers/pkg/aws"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	// "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/gruntwork-io/terratest/modules/random"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "eks/alb-controller-ingress-group/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	defer s.DestroyAtmosComponent(s.T(), component, stack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, nil)
	assert.NotNil(s.T(), options)

	expectedAnnotations := map[string]string{
		"alb.ingress.kubernetes.io/group.name": "alb-controller-ingress-group",
		"alb.ingress.kubernetes.io/listen-ports": "[{\"HTTP\": 80}, {\"HTTPS\": 443}]",
		"alb.ingress.kubernetes.io/scheme": "internet-facing",
		"alb.ingress.kubernetes.io/ssl-policy": "ELBSecurityPolicy-TLS13-1-2-2021-06",
		"alb.ingress.kubernetes.io/ssl-redirect": "443",
		"alb.ingress.kubernetes.io/target-type": "ip",
		"kubernetes.io/ingress.class": "default",
	}
	annotations := atmos.OutputMap(s.T(), options, "annotations")
	for key, expectedValue := range expectedAnnotations {
		assert.Equal(s.T(), expectedValue, annotations[key], "Annotation %s does not match", key)
	}

	groupName := atmos.Output(s.T(), options, "group_name")
	assert.Equal(s.T(), groupName, "alb-controller-ingress-group")

	loadBalancerName := atmos.Output(s.T(), options, "load_balancer_name")
	assert.True(s.T(), strings.HasPrefix(loadBalancerName, "eg-default-ue2-test-alb-co-"))

	host := atmos.Output(s.T(), options, "host")
	assert.NotEmpty(s.T(), host)

	messageBodyLength := atmos.Output(s.T(), options, "message_body_length")
	assert.NotNil(s.T(), messageBodyLength)
	assert.EqualValues(s.T(), messageBodyLength, "1007")

	loadBalancerScheme := atmos.Output(s.T(), options, "load_balancer_scheme")
	assert.NotNil(s.T(), loadBalancerScheme)
	assert.Equal(s.T(), loadBalancerScheme, "internet-facing")

	ingressClass := atmos.Output(s.T(), options, "ingress_class")
	assert.NotNil(s.T(), ingressClass)
	assert.Equal(s.T(), ingressClass, "default")

	s.DriftTest(component, stack, nil)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "eks/alb-controller-ingress-group/disabled"
	const stack = "default-test"
	s.VerifyEnabledFlag(component, stack, nil)
}

func (s *ComponentSuite) SetupSuite() {
	s.TestSuite.InitConfig()
	s.TestSuite.Config.ComponentDestDir = "components/terraform/eks/alb-controller-ingress-group"
	s.TestSuite.SetupSuite()
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)
	suite.AddDependency(t, "vpc", "default-test", nil)
	suite.AddDependency(t, "eks/cluster", "default-test", nil)
	suite.AddDependency(t, "eks/alb-controller", "default-test", nil)

	subdomain := strings.ToLower(random.UniqueId())
	inputs := map[string]interface{}{
		"zone_config": []map[string]interface{}{
			{
				"subdomain": subdomain,
				"zone_name": "components.cptest.test-automation.app",
			},
		},
	}
	suite.AddDependency(t, "dns-delegated", "default-test", &inputs)
	helper.Run(t, suite)
}
