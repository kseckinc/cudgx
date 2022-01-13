package query_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/galaxy-future/cudgx/common/clickhouse"
	"github.com/galaxy-future/cudgx/internal/clients"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Query Suite")
}

var _ = BeforeSuite(func() {
	file, err := os.Open("tests/create_database.sql")
	if err != nil {
		Fail("failed to create test samples" + err.Error())
	}
	cmd := exec.Command("clickhouse-client", "-h", "127.0.0.1", "-mn")
	cmd.Stdin = file
	out, err := cmd.Output()
	if err != nil {
		Fail(fmt.Sprintf("failed to create test samples, err: %s, output: %s ", err.Error(), out))
	}

	file2, err := os.Open("tests/insert_samples.sql")
	if err != nil {
		Fail("failed to create test samples" + err.Error())
	}
	cmd = exec.Command("clickhouse-client", "-h", "127.0.0.1", "-mn")
	cmd.Stdin = file2
	out, err = cmd.Output()
	if err != nil {
		Fail(fmt.Sprintf("failed to create test samples, err: %s, output: %s ", err.Error(), out))
	}

	config := clickhouse.Config{
		Schema:       "http",
		User:         "default",
		Password:     "",
		Database:     "metrics",
		Table:        "metrics_gf_test",
		Hosts:        []string{"127.0.0.1"},
		WriteTimeout: "10s",
		ReadTimeout:  "10s",
	}
	err = clients.InitClickhouseRdCli(&config)
	if err != nil {
		Fail(err.Error())
	}

})

var _ = AfterSuite(func() {
	file, err := os.Open("tests/drop_database.sql")
	if err != nil {
		Fail("failed to drop test samples" + err.Error())
	}

	cmd := exec.Command("clickhouse-client", "-h", "127.0.0.1", "-mn")
	cmd.Stdin = file
	out, err := cmd.Output()
	if err != nil {
		Fail(fmt.Sprintf("failed to drop test samples, err: %s, output: %s ", err.Error(), out))
	}
})
