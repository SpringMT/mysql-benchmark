package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/Songmu/wrapcommander"
	"github.com/jessevdk/go-flags"
)

type mysqlBench struct {
	Host                 string `long:"host" default:"localhost" description:"DBのhost名"`
	Port                 string `long:"port" default:"3306" description:"DBのPort番号"`
	Concurrent           int    `short:"c" default:"50" description:"並列接続数"`
	RequestNum           int    `short:"n" default:"100" description:"リクエスト総数"`
	Sleep                int    `short:"s" default:"1000" description:"スリープ ミリ秒"`
	DbName               string `short:"d" default:"test_seq" description:"データベース名"`
	Table                string `short:"t" default:"sequence" description:"テーブル名"`
	User                 string `short:"u" default:"root" description:"DBのユーザー"`
	Password             string `short:"p" default:"" description:"DBのpassword"`
	Column               string `long:"clumn" default:"id" description:"対象のカラム"`
	Verbose              []bool `short:"v" long:"verbose" description:"verbose output. it can be stacked like -vv for more detailed log"`
	outStream, errStream io.Writer
}

func (mb *mysqlBench) run() error {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	mb.log(debug, "main start")

	var wg sync.WaitGroup
	before := MemConsumed()
	concurrentStream := make(chan interface{}, mb.Concurrent)
	heartBeatStream := make(chan heartBeat, mb.RequestNum)
	for i := 0; i < mb.RequestNum; i++ {
		wg.Add(1)
		concurrentStream <- true
		go func() {
			mb.log(debug, "mysql incr start")
			client := MysqlNewClient(mb.Host, mb.User, mb.Password, mb.Port, mb.DbName)
			beginTime := time.Now()
			var res int64
			var err error
			res, err = client.increment(mb.Table, mb.Column)
			if err == nil {
				heartBeatStream <- heartBeat{
					Time:     now(),
					Status:   Success,
					Incr:     res,
					Duration: time.Since(beginTime),
				}
				mb.logf(info, "%d", res)
			} else {
				heartBeatStream <- heartBeat{
					Time:   now(),
					Status: Failed,
				}
				mb.logf(warn, "error %s", err)
			}
			time.Sleep(time.Duration(mb.Sleep+rand.Intn(mb.Sleep)) * time.Millisecond)
			defer func() {
				wg.Done()
				client.close()
				<-concurrentStream
			}()
		}()
	}
	wg.Wait()
	after := MemConsumed()
	mb.logf(info, "Memory %.3f kb", float64(after-before)/1000)
	// channelはcloseしないとメモリリークの原因になる
	close(concurrentStream)
	close(heartBeatStream)
	results := NewHeartBeatResult()
	for hb := range heartBeatStream {
		results.add(hb)
	}
	results.show()
	return nil
}

func now() *time.Time {
	now := time.Now()
	return &now
}

func parseArgs(args []string) (*flags.Parser, *mysqlBench, error) {
	mb := &mysqlBench{}
	p := flags.NewParser(mb, flags.Default)
	p.Usage = fmt.Sprintf(`--host localhost [...]

Version: %s (rev: %s/%s)`, version, revision, runtime.Version())
	_, err := p.ParseArgs(args)
	mb.outStream = os.Stdout
	mb.errStream = os.Stderr
	return p, mb, err
}

// Run benchmark
func Run(args []string) int {
	p, mb, err := parseArgs(args)
	if err != nil {
		if ferr, ok := err.(*flags.Error); !ok || ferr.Type != flags.ErrHelp {
			p.WriteHelp(mb.errStream)
		}
		return 2
	}
	runErr := mb.run()
	if runErr != nil {
		return wrapcommander.ResolveExitCode(runErr)
	}

	return 0
}

func main() {
	rand.Seed(time.Now().UnixNano())
	os.Exit(Run(os.Args[1:]))
}
