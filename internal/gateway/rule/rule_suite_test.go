package rule_test

import (
	"github.com/galaxy-future/cudgx/internal/gateway/rule"
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var manager *rule.Manager

func TestRule(t *testing.T) {
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

		validOption := &rule.MysqlOption{
			Dsn:            "root:@tcp(localhost:3306)/cudgx_test?charset=utf8mb4&parseTime=True&loc=Local",
			RefreshSeconds: 5,
		}
		manager, err = rule.NewRuleManager(validOption)
		if err != nil {
			Fail(fmt.Sprintf("failed to create RuleManager, err: %s, output: %s ", err.Error(), out))
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

	RunSpecs(t, "Rule Suite")

}
