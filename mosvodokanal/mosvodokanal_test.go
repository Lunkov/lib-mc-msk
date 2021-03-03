// mosvodokanal.go
// Get Data from mosvodokanal.ru
package mosvodokanal

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "flag"
  "time"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-mc"
)

func TestCheckMosvodokanal(t *testing.T) {
  flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()

	glog.Info("Logging configured")
  
  ow := NewWorker()
  
  mc.WorkerRegister(ow)

  go mc.Init("./etc.tests/")
  time.Sleep(2 * time.Second)

  assert.Equal(t, "mosvodokanal.ru", ow.GetAPI())
  assert.Equal(t, true, ow.ClientData.Status.Ok)

  r1 := mc.GetWorkersResults()
  res, _ := (r1["mosvodokanal.ru"][""]).(mc.Result)
  assert.Equal(t, int64(1), res.Status.CntDevices)
  assert.Equal(t, int64(3), res.Status.CntMetrics)
}
