// Copyright 2016 The go-modernizingpark Authors
// This file is part of the go-modernizingpark library.
//
// The go-modernizingpark library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-modernizingpark library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-modernizingpark library. If not, see <http://www.gnu.org/licenses/>.

package debug

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/modernizingpark/go-modernizingpark/log"
	"github.com/modernizingpark/go-modernizingpark/metrics"
	"github.com/modernizingpark/go-modernizingpark/metrics/exp"
	"github.com/fjl/memsize/memsizeui"
	colorable "github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"gopkg.in/urfave/cli.v1"
)

var Memsize memsizeui.Handler
var ID string

var (
	verbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail",
		Value: 3,
	}
	logPathFlag = cli.StringFlag{
		Name:  "logpath",
		Usage: "File path for log files",
		Value: "",
	}

	metricLogFlag = cli.BoolFlag{
		Name:  "metriclog",
		Usage: "Write metric info to log files",
	}

	vmoduleFlag = cli.StringFlag{
		Name:  "vmodule",
		Usage: "Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. eth/*=5,p2p=4)",
		Value: "",
	}
	backtraceAtFlag = cli.StringFlag{
		Name:  "backtrace",
		Usage: "Request a stack trace at a specific logging statement (e.g. \"block.go:271\")",
		Value: "",
	}
	debugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "Prepends log messages with call-site location (file and line number)",
	}
	pprofFlag = cli.BoolFlag{
		Name:  "pprof",
		Usage: "Enable the pprof HTTP server",
	}
	pprofPortFlag = cli.IntFlag{
		Name:  "pprof.port",
		Usage: "pprof HTTP server listening port",
		Value: 6060,
	}
	pprofAddrFlag = cli.StringFlag{
		Name:  "pprof.addr",
		Usage: "pprof HTTP server listening interface",
		Value: "127.0.0.1",
	}
	memprofilerateFlag = cli.IntFlag{
		Name:  "pprof.memprofilerate",
		Usage: "Turn on memory profiling with the given rate",
		Value: runtime.MemProfileRate,
	}
	blockprofilerateFlag = cli.IntFlag{
		Name:  "pprof.blockprofilerate",
		Usage: "Turn on block profiling with the given rate",
	}
	cpuprofileFlag = cli.StringFlag{
		Name:  "pprof.cpuprofile",
		Usage: "Write CPU profile to the given file",
	}
	traceFlag = cli.StringFlag{
		Name:  "trace",
		Usage: "Write execution trace to the given file",
	}
	// (Deprecated April 2020)
	legacyPprofPortFlag = cli.IntFlag{
		Name:  "pprofport",
		Usage: "pprof HTTP server listening port (deprecated, use --pprof.port)",
		Value: 6060,
	}
	legacyPprofAddrFlag = cli.StringFlag{
		Name:  "pprofaddr",
		Usage: "pprof HTTP server listening interface (deprecated, use --pprof.addr)",
		Value: "127.0.0.1",
	}
	legacyMemprofilerateFlag = cli.IntFlag{
		Name:  "memprofilerate",
		Usage: "Turn on memory profiling with the given rate (deprecated, use --pprof.memprofilerate)",
		Value: runtime.MemProfileRate,
	}
	legacyBlockprofilerateFlag = cli.IntFlag{
		Name:  "blockprofilerate",
		Usage: "Turn on block profiling with the given rate (deprecated, use --pprof.blockprofilerate)",
	}
	legacyCpuprofileFlag = cli.StringFlag{
		Name:  "cpuprofile",
		Usage: "Write CPU profile to the given file (deprecated, use --pprof.cpuprofile)",
	}
)

// Flags holds all command-line flags required for debugging.
var Flags = []cli.Flag{
	verbosityFlag, logPathFlag, metricLogFlag, vmoduleFlag, backtraceAtFlag, debugFlag,
	pprofFlag, pprofAddrFlag, pprofPortFlag, memprofilerateFlag,
	blockprofilerateFlag, cpuprofileFlag, traceFlag,
}

var DeprecatedFlags = []cli.Flag{
	legacyPprofPortFlag, legacyPprofAddrFlag, legacyMemprofilerateFlag,
	legacyBlockprofilerateFlag, legacyCpuprofileFlag,
}

var (
	glogger *log.GlogHandler
)

const (
	metricLogFile = "metric.log"
	metricKey     = "metric"
)

func setupLogHandler(ctx *cli.Context) (handler log.Handler) {
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"

	if ctx.GlobalString(logPathFlag.Name) == "" {
		output := io.Writer(os.Stderr)
		if usecolor {
			output = colorable.NewColorableStderr()
		}
		handler = log.StreamHandler(output, log.TerminalFormat(usecolor))
		return
	}

	rConfig := log.NewRotateConfig()
	rConfig.LogDir = ctx.GlobalString(logPathFlag.Name)
	handler1 := log.NewFileRotateHandler(rConfig, log.TerminalFormat(usecolor))
	if !ctx.GlobalBool(metricLogFlag.Name) {
		handler = handler1
		return
	}

	mConfig := log.NewRotateConfig()
	mConfig.LogDir = ctx.GlobalString(logPathFlag.Name)
	mConfig.Filename = metricLogFile
	handler2 := log.NewFileRotateHandler(mConfig, log.JSONFormat())

	handler = log.FuncHandler(func(r *log.Record) error {
		if r.Msg == metricKey {
			r.Ctx = append(r.Ctx, "id", ID)
			return handler2.Log(r)
		} else {
			return handler1.Log(r)
		}
	})

	return
}

// Setup initializes profiling and logging based on the CLI flags.
// It should be called as early as possible in the program.
func Setup(ctx *cli.Context) error {
	// logging
	handler := setupLogHandler(ctx)
	glogger = log.NewGlogHandler(handler)

	log.PrintOrigins(ctx.GlobalBool(debugFlag.Name))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(verbosityFlag.Name)))
	glogger.Vmodule(ctx.GlobalString(vmoduleFlag.Name))
	glogger.BacktraceAt(ctx.GlobalString(backtraceAtFlag.Name))
	log.Root().SetHandler(glogger)

	// profiling, tracing
	if ctx.GlobalIsSet(legacyMemprofilerateFlag.Name) {
		runtime.MemProfileRate = ctx.GlobalInt(legacyMemprofilerateFlag.Name)
		log.Warn("The flag --memprofilerate is deprecated and will be removed in the future, please use --pprof.memprofilerate")
	}
	runtime.MemProfileRate = ctx.GlobalInt(memprofilerateFlag.Name)

	if ctx.GlobalIsSet(legacyBlockprofilerateFlag.Name) {
		Handler.SetBlockProfileRate(ctx.GlobalInt(legacyBlockprofilerateFlag.Name))
		log.Warn("The flag --blockprofilerate is deprecated and will be removed in the future, please use --pprof.blockprofilerate")
	}
	Handler.SetBlockProfileRate(ctx.GlobalInt(blockprofilerateFlag.Name))

	if traceFile := ctx.GlobalString(traceFlag.Name); traceFile != "" {
		if err := Handler.StartGoTrace(traceFile); err != nil {
			return err
		}
	}

	if cpuFile := ctx.GlobalString(cpuprofileFlag.Name); cpuFile != "" {
		if err := Handler.StartCPUProfile(cpuFile); err != nil {
			return err
		}
	}
	if cpuFile := ctx.GlobalString(legacyCpuprofileFlag.Name); cpuFile != "" {
		log.Warn("The flag --cpuprofile is deprecated and will be removed in the future, please use --pprof.cpuprofile")
		if err := Handler.StartCPUProfile(cpuFile); err != nil {
			return err
		}
	}

	// pprof server
	if ctx.GlobalBool(pprofFlag.Name) {
		listenHost := ctx.GlobalString(pprofAddrFlag.Name)
		if ctx.GlobalIsSet(legacyPprofAddrFlag.Name) && !ctx.GlobalIsSet(pprofAddrFlag.Name) {
			listenHost = ctx.GlobalString(legacyPprofAddrFlag.Name)
			log.Warn("The flag --pprofaddr is deprecated and will be removed in the future, please use --pprof.addr")
		}

		port := ctx.GlobalInt(pprofPortFlag.Name)
		if ctx.GlobalIsSet(legacyPprofPortFlag.Name) && !ctx.GlobalIsSet(pprofPortFlag.Name) {
			port = ctx.GlobalInt(legacyPprofPortFlag.Name)
			log.Warn("The flag --pprofport is deprecated and will be removed in the future, please use --pprof.port")
		}

		address := fmt.Sprintf("%s:%d", listenHost, port)
		// This context value ("metrics.addr") represents the utils.MetricsHTTPFlag.Name.
		// It cannot be imported because it will cause a cyclical dependency.
		StartPProf(address, !ctx.GlobalIsSet("metrics.addr"))
	}
	return nil
}

func StartPProf(address string, withMetrics bool) {
	// Hook go-metrics into expvar on any /debug/metrics request, load all vars
	// from the registry into expvar, and execute regular expvar handler.
	if withMetrics {
		exp.Exp(metrics.DefaultRegistry)
	}
	http.Handle("/memsize/", http.StripPrefix("/memsize", &Memsize))
	log.Info("Starting pprof server", "addr", fmt.Sprintf("http://%s/debug/pprof", address))
	go func() {
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Error("Failure in running pprof server", "err", err)
		}
	}()
}

// Exit stops all running profiles, flushing their output to the
// respective file.
func Exit() {
	Handler.StopCPUProfile()
	Handler.StopGoTrace()
}
