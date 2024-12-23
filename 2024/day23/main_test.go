package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA1 string = `kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn`

func Test(t *testing.T) {
	t.Run("TestComputerNetworks", func(t *testing.T) {

		network := ParseInputData(EXAMPLEDATA1)

		assert.Equal(t, 7, CountTConnections(network))
		assert.Equal(t, "co,de,ka,ta", FindMaxCluster(network))
	})
}
