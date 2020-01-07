package t

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const InnerQps int = 10

type respon struct {
	err           error
	statusCode    int
	duration      float64
	contentLength int64
	body          string
}

type TaoR struct {
	Method, Url, Form string
}

type Tao struct {
	client    *http.Client
	duration  int
	qps       int
	okCount   int
	reqCount  int
	durations float64
	reqStat   reqStatus
	bodyShow  bool
	req       *TaoR
}

type reqStatus struct {
	t50      int
	t50_99   int
	t100_199 int
	t200_299 int
	t300_399 int
	t400_499 int
	t500     int
}

var (
	cliOnce sync.Once
	client  *http.Client
)

func GetClient(disKA, disCompression bool, timeout int) *http.Client {
	cliOnce.Do(func() {
		httpClient := new(http.Client)
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableKeepAlives:   disKA,
			DisableCompression:  disCompression,
			TLSHandshakeTimeout: time.Duration(timeout) * time.Millisecond,
		}
		client = httpClient
	})
	return client
}

func SetupTao(c *http.Client, duration, qps int, r *TaoR, bodyShow bool) *Tao {
	t := new(Tao)
	t.client = c
	t.req = r
	t.duration = duration
	t.qps = qps
	t.okCount = 0
	t.reqCount = duration * qps
	t.bodyShow = bodyShow
	return t
}

func (t *Tao) newRequest() *http.Request {
	form := t.req.Form
	req, err := http.NewRequest(t.req.Method, t.req.Url, strings.NewReader(form))
	if err != nil {
		println(err.Error())
		panic("new request error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func (t *Tao) Run() {
	num_chanel := t.qps / InnerQps
	now := time.Now()
	var wg sync.WaitGroup
	wg.Add(num_chanel)
	for i := 0; i < num_chanel; i++ {
		go func() {
			hits := 0
			for c := range time.NewTicker(time.Duration(1e3/InnerQps) * time.Millisecond).C {
				var size int64
				var code int
				if resp, err := t.client.Do(t.newRequest()); resp != nil {
					if err == nil && resp.StatusCode == 200 {
						size = resp.ContentLength
						code = resp.StatusCode
						bs, _ := ioutil.ReadAll(resp.Body)
						resp.Body.Close()
						t.okCount += 1
						hits += 1
						duration := time.Now().Sub(c).Seconds() * 1000
						printResult(respon{
							statusCode:    code,
							duration:      duration,
							contentLength: size,
							body:          string(bs),
						}, t.bodyShow)
						t.durations += duration
						switch {
						case duration < float64(50):
							t.reqStat.t50 += 1
						case duration < float64(100):
							t.reqStat.t50_99 += 1
						case duration < float64(200):
							t.reqStat.t100_199 += 1
						case duration < float64(300):
							t.reqStat.t200_299 += 1
						case duration < float64(400):
							t.reqStat.t300_399 += 1
						case duration < float64(500):
							t.reqStat.t400_499 += 1
						case duration >= float64(500):
							t.reqStat.t500 += 1
						}
					} else {
						printErrorResult(respon{
							err:           err,
							statusCode:    resp.StatusCode,
							duration:      time.Now().Sub(c).Seconds() * 1000,
							contentLength: resp.ContentLength,
						})
					}
					if hits >= t.duration*InnerQps {
						wg.Done()
						break
					}
				} else {
					println(err.Error())
					panic("Client Do error")
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Done! request count-[%d] | ok count-[%d] | success rate-[%1.1f] | average request time-[%4.2f]ms | request status %s | cost time-[%4.3f]seconds", t.reqCount, t.okCount, float64(t.okCount)/float64(t.reqCount)*100, t.durations/float64(t.reqCount), printStatMap(t.reqStat), time.Now().Sub(now).Seconds())
}
