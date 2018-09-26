package gracehttp

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"filelog"

	"github.com/facebookgo/grace/gracenet"
)

const (
	defaultShutdownTimeout = time.Minute
)

var (
	isChild   = os.Getenv("LISTEN_FDS") != "" // gracenet 中定义的环境变量。
	ppid      = os.Getppid()
	pid       = os.Getpid()
	serverLog = &filelog.ReqLogger{ReqId: "server", Logger: filelog.ServerLogger}
)

type Server struct {
	*http.Server
	ln       net.Listener
	CertFile string
	KeyFile  string

	Name            string
	ShutdownTimeout time.Duration
	WrappedListener func(net.Listener) net.Listener
	BeforeShutdown  func()
}

type App struct {
	servers []*Server
	net     *gracenet.Net
}

func NewApp(servers []*Server) *App {
	app := &App{
		servers: servers,
		net:     new(gracenet.Net),
	}
	return app
}

func (a *App) Run() error {
	if serverLog.Logger == nil {
		serverLog.Logger = filelog.ServerLogger
	}
	// create or take over listeners.
	for _, server := range a.servers {
		ln, err := a.net.Listen("tcp", server.Addr)
		if err != nil {
			return err
		}
		if server.WrappedListener != nil {
			ln = server.WrappedListener(ln)
		}
		server.ln = ln
	}

	// print infomation of listeners.
	msg := fmt.Sprintf("PID(%d) is listening on [%s]", pid, prettyAddr(a.servers))
	if isChild {
		msg += fmt.Sprintf(" which are taken over from ppid(%d)", ppid)
	}
	serverLog.Info(msg)

	// serve http.
	wg := new(sync.WaitGroup)
	wg.Add(len(a.servers))
	for _, server := range a.servers {
		go func(server *Server) {
			defer wg.Done()
			serverLog.Infof("PID(%d) %s is running on %s", pid, server.Name, server.ln.Addr().String())
			var err error
			if server.CertFile != "" && server.KeyFile != "" {
				err = server.ServeTLS(server.ln, server.CertFile, server.KeyFile)
			} else {
				err = server.Serve(server.ln)
			}
			if err != nil {
				serverLog.Warnf("%s exits, err: %v", server.Name, err)
			}
		}(server)
	}

	// since servers are up, now close the parent process.
	if isChild && ppid != 1 {
		PProcess, err := os.FindProcess(ppid)
		if err != nil {
			serverLog.Errorf("Failed to find parent process PID(%d), err: %v", ppid, err)
			return err
		}
		err = PProcess.Kill()
		if err != nil {
			serverLog.Errorf("Failed to find parent process PID(%d), err: %v", ppid, err)
			return err
		}
	}

	// signal control.
	go a.signalHandler(wg)

	// waiting server to exit and gracefully shutdown.
	wg.Wait()

	return nil
}

// SIGINT & SIGTERM 优雅退出信号
// SIGHUP 热重启信号
func (a *App) signalHandler(wg *sync.WaitGroup) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		sig := <-ch
		serverLog.Infof("Receiving signal: %v", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:

			// 当收到优雅关闭信号后，不再特殊处理后续的任何信号，并退出信号处理器。
			signal.Stop(ch)
			a.shutdown(wg)
			return
		case syscall.SIGHUP:
			childPID, err := a.net.StartProcess()
			if err != nil {
				serverLog.Errorf("Failed to start child process, err: %v", err)
				continue
			}

			serverLog.Infof("Success to start child process, pid: %d", childPID)

			go func() {
				// 子进程启动后可能会在向父进程发送 SIGTERM 前异常退出，为了避免其成为僵尸进程，这里需要回收子进程。
				childProcess, err := os.FindProcess(childPID)
				if err == nil {
					state, err := childProcess.Wait()
					serverLog.Warnf("Child process exits, state: %v, err: %v", state, err)
				}
			}()

			// 启动子进程后不退出，由子进程给自己发送 SIGTERM 后再退出。
		}
	}
}

func (a *App) shutdown(wg *sync.WaitGroup) {
	wg.Add(len(a.servers))
	for _, server := range a.servers {
		timeout := server.ShutdownTimeout
		if timeout == 0 {
			timeout = defaultShutdownTimeout
		}
		go func(server *Server) {
			defer wg.Done()

			if server.BeforeShutdown != nil {
				server.BeforeShutdown()
			}
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			serverLog.Infof("PID(%d) %s is shutting down...", pid, server.Name)
			err := server.Shutdown(ctx)
			serverLog.Infof("PID(%d) %s finish shutting down, cost: %v, err: %v", pid, server.Name, time.Since(start), err)
		}(server)
	}
}

func prettyAddr(ss []*Server) string {
	var buf bytes.Buffer
	for i, s := range ss {
		if i != 0 {
			fmt.Fprint(&buf, ", ")
		}
		fmt.Fprint(&buf, s.ln.Addr())
	}
	return buf.String()
}
