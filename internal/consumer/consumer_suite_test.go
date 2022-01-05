package consumer_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConsumer(t *testing.T) {
	RegisterFailHandler(Fail)
	BeforeSuite(func() {
		file, err := os.Open("tests/create_database.sql")
		if err != nil {
			Fail("failed to create test samples" + err.Error())
		}
		cmd := exec.Command("clickhouse-client", "-h", "localhost", "-u", "default", "--password", "test")
		cmd.Stdin = file
		out, err := cmd.Output()
		if err != nil {
			Fail(fmt.Sprintf("failed to create test samples, err: %s, output: %s ", err.Error(), out))
		}

	})

	AfterSuite(func() {
		file, err := os.Open("tests/drop_database.sql")
		if err != nil {
			Fail("failed to drop test samples" + err.Error())
		}

		cmd := exec.Command("clickhouse-client", "-h", "localhost", "-u", "default", "--password", "test")
		cmd.Stdin = file
		out, err := cmd.Output()
		if err != nil {
			Fail(fmt.Sprintf("failed to drop test samples, err: %s, output: %s ", err.Error(), out))
		}
	})
	RunSpecs(t, "Consumer Suite")
}
