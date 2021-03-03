// mosvodokanal.go
// Get Data from mosvodokanal.ru
package mosvodokanal

import (
  "time"
  "strings"
  "strconv"
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "github.com/golang/glog"
  "github.com/Lunkov/lib-mc"
)


type WorkerInfo struct {
  mc.WorkerInfo
}

func NewWorker() *WorkerInfo {
  w := new(WorkerInfo)
  w.API = "mosvodokanal.ru"
  return w
}

type ParamsInfo struct {
  Id int `json:"id"`
  Value string `json:"value"`
  InRange int `json:"inrange"`
}

type ResultInfo struct {
  DT_from string `json:"dt_from"`
  DT_to string `json:"dt_to"`
  Params []ParamsInfo `json:"params"`
}

type Info struct {
  Message string `json:"message"`
  Code int `json:"code"`
  Result ResultInfo `json:"result"`
}

func (w *WorkerInfo) httpGet() Info {
  var data Info
  resp, err := http.Get(w.ClientData.Url)
  if err != nil {
    w.ClientData.Status.Ok = false
    w.ClientData.Status.LastError = fmt.Sprintf("Mosvodokanal: %s", err)
    glog.Errorf("ERR: %s\n", w.ClientData.Status.LastError)
    return data
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if glog.V(9) {
    glog.Infof("DBG: Mosvodokanal BODY(%s)", body)
  }
  if err := json.Unmarshal(body, &data); err != nil {
    w.ClientData.Status.Ok = false
    w.ClientData.Status.LastError = fmt.Sprintf("Mosvodokanal: %s", err)
    glog.Errorf("ERR: %s\n", w.ClientData.Status.LastError)
    return data
  }
  w.ClientData.Status.Ok = true
  return data
}

func (w *WorkerInfo) GetData() {
  if glog.V(2) {
    glog.Infof("LOG: Mosvodokanal started")
  }
  data := w.httpGet()
  if !w.ClientData.Status.Ok {
    glog.Infof("ERR: Mosvodokanal Error")
    return
  }
  var dm []mc.DeviceMetric
  dt := time.Now()
  w.ClientData.Status.CntDevices = 1
  if glog.V(9) {
    glog.Infof("DBG: Mosvodokanal DATA(%v)", data)
  }
  r := strings.NewReplacer(">", "",
                           "<", "",
                           "=", "")
  for _, v := range data.Result.Params {
    if mCODE, ok := w.ClientData.ParamsCode[strconv.Itoa(v.Id)]; ok {
      v.Value = r.Replace(v.Value)
      if glog.V(9) {
        glog.Infof("DBG: Mosvodokanal (%s) Value=%s", mCODE, v.Value)
      }
      if ff, err := strconv.ParseFloat(v.Value, 32); err == nil {
        if glog.V(9) {
          glog.Infof("DBG: Mosvodokanal (%s) Value=%f\n", mCODE, ff)
        }
        dm = append(dm, mc.DeviceMetric{Device_CODE: w.ClientData.Device_CODE, Metric_CODE: mCODE, DT: dt, Value: ff})
        w.ClientData.Status.CntMetrics ++
      }
    }
  }
  w.ClientData.Status.Ok = true
  w.SendMetrics(&dm)
  if glog.V(2) {
    glog.Infof("LOG: Mosvodokanal finished: Ok=%v", w.ClientData.Status.Ok)
  }
}
