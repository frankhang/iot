// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/frankhang/util/sys/linux"
	"github.com/frankhang/util/tcp"

	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/log"

	"github.com/frankhang/util/metrics"
	"github.com/frankhang/util/signal"
	"github.com/frankhang/util/systimemon"
	"github.com/opentracing/opentracing-go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/struCoder/pidusage"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

// Flag Names
const (
	nmVersion          = "V"
	nmConfig           = "config"
	nmConfigCheck      = "config-check"
	nmConfigStrict     = "config-strict"
	nmStore            = "store"
	nmStorePath        = "path"
	nmHost             = "host"
	nmAdvertiseAddress = "advertise-address"
	nmPort             = "P"
	nmCors             = "cors"
	nmSocket           = "socket"
	nmEnableBinlog     = "enable-binlog"
	nmRunDDL           = "run-ddl"
	nmLogLevel         = "L"
	nmLogFile          = "log-file"
	nmLogSlowQuery     = "log-slow-query"
	nmReportStatus     = "report-status"
	nmStatusHost       = "status-host"
	nmStatusPort       = "status"
	nmMetricsAddr      = "metrics-addr"
	nmMetricsInterval  = "metrics-interval"
	nmDdlLease         = "lease"
	nmTokenLimit       = "token-limit"
	nmPluginDir        = "plugin-dir"
	nmPluginLoad       = "plugin-load"

	nmProxyProtocolNetworks      = "proxy-protocol-networks"
	nmProxyProtocolHeaderTimeout = "proxy-protocol-header-timeout"
	nmAffinityCPU                = "affinity-cpus"
)

var (
	version      = flagBoolean(nmVersion, false, "print version information and exit")
	configPath   = flag.String(nmConfig, "", "config file path")
	configCheck  = flagBoolean(nmConfigCheck, false, "check config file validity and exit")
	configStrict = flagBoolean(nmConfigStrict, false, "enforce config file validity")

	// Base
	store            = flag.String(nmStore, "mocktikv", "registered store name, [tikv, mocktikv]")
	storePath        = flag.String(nmStorePath, "/tmp/tidb", "tidb storage path")
	host             = flag.String(nmHost, "0.0.0.0", "tidb server host")
	advertiseAddress = flag.String(nmAdvertiseAddress, "", "tidb server advertise IP")
	port             = flag.String(nmPort, "4000", "tidb server port")
	cors             = flag.String(nmCors, "", "tidb server allow cors origin")
	socket           = flag.String(nmSocket, "", "The socket file to use for connection.")
	enableBinlog     = flagBoolean(nmEnableBinlog, false, "enable generate binlog")
	runDDL           = flagBoolean(nmRunDDL, true, "run ddl worker on this tidb-server")
	ddlLease         = flag.String(nmDdlLease, "45s", "schema lease duration, very dangerous to change only if you know what you do")
	tokenLimit       = flag.Int(nmTokenLimit, 1000, "the limit of concurrent executed sessions")
	pluginDir        = flag.String(nmPluginDir, "/data/deploy/plugin", "the folder that hold plugin")
	pluginLoad       = flag.String(nmPluginLoad, "", "wait load plugin name(separated by comma)")
	affinityCPU      = flag.String(nmAffinityCPU, "", "affinity cpu (cpu-no. separated by comma, e.g. 1,2,3)")

	// Log
	logLevel     = flag.String(nmLogLevel, "info", "log level: info, debug, warn, error, fatal")
	logFile      = flag.String(nmLogFile, "", "log file path")
	logSlowQuery = flag.String(nmLogSlowQuery, "", "slow query file path")

	// Status
	reportStatus    = flagBoolean(nmReportStatus, true, "If enable status report HTTP service.")
	statusHost      = flag.String(nmStatusHost, "0.0.0.0", "tidb server status host")
	statusPort      = flag.String(nmStatusPort, "10080", "tidb server status port")
	metricsAddr     = flag.String(nmMetricsAddr, "", "prometheus pushgateway address, leaves it empty will disable prometheus push.")
	metricsInterval = flag.Uint(nmMetricsInterval, 15, "prometheus client push interval in second, set \"0\" to disable prometheus push.")

	// PROXY Protocol
	proxyProtocolNetworks      = flag.String(nmProxyProtocolNetworks, "", "proxy protocol networks allowed IP or *, empty mean disable proxy protocol support")
	proxyProtocolHeaderTimeout = flag.Uint(nmProxyProtocolHeaderTimeout, 5, "proxy protocol header read timeout, unit is second.")
)

var (
	cfg      *tcp.Config
	svr      *tcp.Server
	graceful bool
)

func main() {
	flag.Parse()
	if *version {
		//fmt.Println(printer.GetTiDBInfo())
		os.Exit(0)
	}

	registerMetrics()
	configWarning := loadConfig()

	if err := cfg.Valid(); err != nil {
		fmt.Fprintln(os.Stderr, "invalid config", err)
		os.Exit(1)
	}
	if *configCheck {
		fmt.Println("config check successful")
		os.Exit(0)
	}
	setGlobalVars()
	setCPUAffinity()
	setupLog()
	// If configStrict had been specified, and there had been an error, the server would already
	// have exited by now. If configWarning is not an empty string, write it to the log now that
	// it's been properly set up.
	if configWarning != "" {
		log.Warn(configWarning)
	}
	setupTracing() // Should before createServer and after setup config.
	printInfo()
	setupMetrics()
	createServer()
	signal.SetupSignalHandler(serverShutdown)
	runServer()
	cleanup()
	syncLog()
}

func exit() {
	syncLog()
	os.Exit(0)
}

func syncLog() {
	if err := log.Sync(); err != nil {
		fmt.Fprintln(os.Stderr, "sync log err:", err)
		os.Exit(1)
	}
}

func setCPUAffinity() {
	if affinityCPU == nil || len(*affinityCPU) == 0 {
		return
	}
	var cpu []int
	for _, af := range strings.Split(*affinityCPU, ",") {
		af = strings.TrimSpace(af)
		if len(af) > 0 {
			c, err := strconv.Atoi(af)
			if err != nil {
				fmt.Fprintf(os.Stderr, "wrong affinity cpu config: %s", *affinityCPU)
				exit()
			}
			cpu = append(cpu, c)
		}
	}
	err := linux.SetAffinity(cpu)
	if err != nil {
		fmt.Fprintf(os.Stderr, "set cpu affinity failure: %v", err)
		exit()
	}
	runtime.GOMAXPROCS(len(cpu))
}

func registerMetrics() {
	metrics.RegisterMetrics()
}

// Prometheus push.
const zeroDuration = time.Duration(0)

// pushMetric pushes metrics in background.
func pushMetric(addr string, interval time.Duration) {
	if interval == zeroDuration || len(addr) == 0 {
		log.Info("disable Prometheus push client")
		return
	}
	log.Info("start prometheus push client", zap.String("server addr", addr), zap.String("interval", interval.String()))
	go prometheusPushClient(addr, interval)
}

// prometheusPushClient pushes metrics to Prometheus Pushgateway.
func prometheusPushClient(addr string, interval time.Duration) {
	// TODO: TiDB do not have uniq name, so we use host+port to compose a name.
	job := "tidb"
	pusher := push.New(addr, job)
	pusher = pusher.Gatherer(prometheus.DefaultGatherer)
	pusher = pusher.Grouping("instance", instanceName())
	for {
		err := pusher.Push()
		if err != nil {
			log.Error("could not push metrics to prometheus pushgateway", zap.String("err", err.Error()))
		}
		time.Sleep(interval)
	}
}

func instanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return fmt.Sprintf("%s_%d", hostname, cfg.Port)
}

// parseDuration parses lease argument string.
func parseDuration(lease string) time.Duration {
	dur, err := time.ParseDuration(lease)
	if err != nil {
		dur, err = time.ParseDuration(lease + "s")
	}
	if err != nil || dur < 0 {
		log.Fatal("invalid lease duration", zap.String("lease", lease))
	}
	return dur
}

func flagBoolean(name string, defaultVal bool, usage string) *bool {
	if !defaultVal {
		// Fix #4125, golang do not print default false value in usage, so we append it.
		usage = fmt.Sprintf("%s (default false)", usage)
		return flag.Bool(name, defaultVal, usage)
	}
	return flag.Bool(name, defaultVal, usage)
}

//var deprecatedConfig = map[string]struct{}{
//	"pessimistic-txn.ttl": {},
//	"log.rotate":          {},
//}

func loadConfig() string {

	//errors.MustNil(err)

	return ""
}

func setGlobalVars() {

}

func setupLog() {
	//err := logutil.InitZapLogger(cfg.Log.ToLogConfig())
	//errors.MustNil(err)
	//
	//err = logutil.InitLogger(cfg.Log.ToLogConfig())
	//errors.MustNil(err)
}

func printInfo() {
	// Make sure the info is always printed.
	level := log.GetLevel()
	log.SetLevel(zap.InfoLevel)
	//printer.PrintTiDBInfo()
	log.SetLevel(level)
}

func createServer() {
	tierDriver := NewTireDriver(cfg, &DefaultTierPacketIO, &DefaultTierHanlder)
	var err error
	svr, err = tcp.NewServer(cfg, tierDriver)
	errors.MustNil(err)

}

func serverShutdown(isgraceful bool) {
	if isgraceful {
		graceful = true
	}
	svr.Close()
}

func setupMetrics() {
	// Enable the mutex profile, 1/10 of mutex blocking event sampling.
	runtime.SetMutexProfileFraction(10)
	systimeErrHandler := func() {
		metrics.TimeJumpBackCounter.Inc()
	}
	callBackCount := 0
	sucessCallBack := func() {
		callBackCount++
		// It is callback by monitor per second, we increase metrics.KeepAliveCounter per 5s.
		if callBackCount >= 5 {
			callBackCount = 0
			metrics.KeepAliveCounter.Inc()
			updateCPUUsageMetrics()
		}
	}
	go systimemon.StartMonitor(time.Now, systimeErrHandler, sucessCallBack)

	pushMetric(cfg.Status.MetricsAddr, time.Duration(cfg.Status.MetricsInterval)*time.Second)
}

func updateCPUUsageMetrics() {
	sysInfo, err := pidusage.GetStat(os.Getpid())
	if err != nil {
		return
	}
	metrics.CPUUsagePercentageGauge.Set(sysInfo.CPU)
}

func setupTracing() {
	tracingCfg := cfg.OpenTracing.ToTracingConfig()
	tracer, _, err := tracingCfg.New("TiDB")
	if err != nil {
		log.Fatal("setup jaeger tracer failed", zap.String("error message", err.Error()))
	}
	opentracing.SetGlobalTracer(tracer)
}

func runServer() {
	err := svr.Run()
	errors.MustNil(err)
}

func cleanup() {
	if graceful {
		svr.GracefulDown(context.Background(), nil)
	} else {
		svr.TryGracefulDown()
	}

}
