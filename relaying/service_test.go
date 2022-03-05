package relaying_test

import (
	"github.com/itsnoproblem/pokt-lint/mock"
	"github.com/itsnoproblem/pokt-lint/relaying"
	"testing"
)

func TestNodeChecker_RunRelayTests(t *testing.T) {
	nodeId := "123abc"
	nodeAddress := "https://node-000.mynode.com"
	chains := []string{"0001"}

	client := mock.NewFakeHTTPClient(true)
	svc, err := relaying.NewNodeChecker(nodeId, nodeAddress, chains, client)
	if err != nil {
		t.Fatalf("got error instantiating node checker: %s", err)
	}

	_, err = svc.RunRelayTests()
	if err != nil {
		t.Fatalf("got error running relay tests: %s", err)
	}
}
