package amqputils_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func BenchmarkChannelPool(b *testing.B) {
	ctx := context.Background()
	conn := amqptest.NewConnection(b, testVhost)
	cp := &amqputils.ChannelPool{
		Channel: amqputils.NewChannelGetterConnection(conn),
	}
	chn, err := cp.Get(ctx)
	if err != nil {
		testutils.FatalErr(b, err)
	}
	cp.Put(chn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chn, err = cp.Get(ctx)
		if err != nil {
			testutils.FatalErr(b, err)
		}
		cp.Put(chn)
	}
	err = cp.Close()
	if err != nil {
		testutils.FatalErr(b, err)
	}
}
