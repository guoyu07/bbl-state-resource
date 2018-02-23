package acceptance_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Test Suite")
}

var (
	checkBinaryPath   string
	inBinaryPath      string
	outBinaryPath     string
	serviceAccountKey string
	projectId         string
)

var _ = BeforeSuite(func() {
	var err error
	checkBinaryPath, err = gexec.Build("github.com/cloudfoundry/bbl-state-resource/cmd/check")
	Expect(err).NotTo(HaveOccurred())
	inBinaryPath, err = gexec.Build("github.com/cloudfoundry/bbl-state-resource/cmd/in")
	Expect(err).NotTo(HaveOccurred())
	outBinaryPath, err = gexec.Build("github.com/cloudfoundry/bbl-state-resource/cmd/out")
	Expect(err).NotTo(HaveOccurred())

	Expect(os.Getenv("BBL_GCP_SERVICE_ACCOUNT_KEY")).NotTo(Equal(""))

	serviceAccountKey, err = getGCPServiceAccountKey(os.Getenv("BBL_GCP_SERVICE_ACCOUNT_KEY"))
	Expect(err).NotTo(HaveOccurred())

	p := struct {
		ProjectId string `json:"project_id"`
	}{}
	if err := json.Unmarshal([]byte(serviceAccountKey), &p); err != nil {
		return Storage{}, fmt.Errorf("Unmarshalling account key for project id: %s", err)
	}
	projectId = p.ProjectId
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func getGCPServiceAccountKey(key string) (string, error) {
	if _, err := os.Stat(key); err != nil {
		return key, nil
	}
	return readGCPServiceAccountKey(key)
}

func readGCPServiceAccountKey(path string) (string, error) {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Reading service account key: %v", err)
	}
	return string(keyBytes), nil
}
