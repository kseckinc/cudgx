package gateway_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		file, err := os.Open("tests/create_test_samples.sql")
		if err != nil {
			Fail("failed to create test samples" + err.Error())
		}
		cmd := exec.Command("mysql", "-uroot")
		cmd.Stdin = file
		out, err := cmd.Output()
		if err != nil {
			Fail(fmt.Sprintf("failed to create test samples, err: %s, output: %s ", err.Error(), out))
		}

	})

	AfterSuite(func() {
		file, err := os.Open("tests/drop_test_samples.sql")
		if err != nil {
			Fail("failed to drop test samples" + err.Error())
		}

		cmd := exec.Command("mysql", "-uroot")
		cmd.Stdin = file
		out, err := cmd.Output()
		if err != nil {
			Fail(fmt.Sprintf("failed to drop test samples, err: %s, output: %s ", err.Error(), out))
		}

	})

	RunSpecs(t, "Gateway Suite")
}
