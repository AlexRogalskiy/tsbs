package devops

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hagen1778/tsbs/cmd/tsbs_generate_data/common"
	"github.com/hagen1778/tsbs/cmd/tsbs_generate_data/serialize"
)

var (
	labelNet             = []byte("net") // heap optimization
	labelNetTagInterface = []byte("interface")

	// Reuse NormalDistributions as arguments to other distributions. This is
	// safe to do because the higher-level distribution advances the ND and
	// immediately uses its value and saves the state
	highND = common.ND(50, 1)
	lowND  = common.ND(5, 1)

	NetFields = []labeledDistributionMaker{
		{[]byte("bytes_sent"), func() common.Distribution { return common.MWD(highND, 0) }},
		{[]byte("bytes_recv"), func() common.Distribution { return common.MWD(highND, 0) }},
		{[]byte("packets_sent"), func() common.Distribution { return common.MWD(highND, 0) }},
		{[]byte("packets_recv"), func() common.Distribution { return common.MWD(highND, 0) }},
		{[]byte("err_in"), func() common.Distribution { return common.MWD(lowND, 0) }},
		{[]byte("err_out"), func() common.Distribution { return common.MWD(lowND, 0) }},
		{[]byte("drop_in"), func() common.Distribution { return common.MWD(lowND, 0) }},
		{[]byte("drop_out"), func() common.Distribution { return common.MWD(lowND, 0) }},
	}
)

type NetMeasurement struct {
	*subsystemMeasurement
	interfaceName []byte
}

func NewNetMeasurement(start time.Time) *NetMeasurement {
	sub := newSubsystemMeasurementWithDistributionMakers(start, NetFields)
	interfaceName := []byte(fmt.Sprintf("eth%d", rand.Intn(4)))
	return &NetMeasurement{
		subsystemMeasurement: sub,
		interfaceName:        interfaceName,
	}
}

func (m *NetMeasurement) ToPoint(p *serialize.Point) {
	m.toPointAllInt64(p, labelNet, NetFields)
	p.AppendTag(labelNetTagInterface, m.interfaceName)
}
