package cftools_test

import (
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/cloudnativego/cf-tools"
)

func TestGetVCAPServiceProperty(t *testing.T) {
	fakeVCAP := []string{
		`VCAP_APPLICATION={}`,
		`VCAP_SERVICES={"user-provided": [{"credentials": {"target-url": "http://beer.com"}, "label": "user-provided", "name": "beer", "syslog_drain_url": "", "tags": []}]}`,
	}

	testEnv := cfenv.Env(fakeVCAP)
	cfenv, err := cfenv.New(testEnv)

	if cfenv == nil || err != nil {
		t.Errorf("cfenv is nil: %v", err.Error())
	}

	propertyValue, err := GetVCAPServiceProperty("beer", "target-url", cfenv)
	if err != nil {
		t.Errorf("Could not find service property: %v", err.Error())
	}

	expectedValue := "http://beer.com"
	if propertyValue != expectedValue {
		t.Errorf("Expected property to equal %s, got %s instead.", expectedValue, propertyValue)
	}
}
