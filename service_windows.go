// +build windows
package main

import (
	"fmt"
	"os"

	ptp "github.com/subutai-io/p2p/lib"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

type P2PService struct{}

func (m *P2PService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	go ExecDaemon(52523, TargetURL, "", "", "", DefaultLog, ptp.DefaultMTU, ptp.UsePMTU)
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	//	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				//				changes <- svc.Status{State: svc.StopPending}
				changes <- svc.Status{State: svc.Stopped, Accepts: cmdsAccepted}
				break loop
			case svc.Pause:
				//				changes <- svc.Status{State: svc.PausePending}
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				//				changes <- svc.Status{State: svc.ContinuePending}
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				ptp.Log(ptp.Error, "Unexpected control request #%d", c)
			}
		}
	}
	return
}

func ExecService() error {
	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		ptp.Log(ptp.Error, "Failed to determine if we are running in an interactive session: %v", err)
		os.Exit(106)
		return nil
	}
	if isIntSess {
		ptp.Log(ptp.Info, "Running in an interactive session")
		elog := debug.New("Subutai P2P")
		defer elog.Close()
		elog.Info(1, fmt.Sprintf("Debug mode ON"))
		run := debug.Run
		err = run("Subutai P2P", &P2PService{})
		if err != nil {
			elog.Info(1, fmt.Sprintf("Failed to run service: %s", err))
			return nil
		}
	} else {
		elog, err := eventlog.Open("Subutai P2P")
		if err != nil {
			ptp.Log(ptp.Error, "Failed to get access to event logger")
			return nil
		}
		defer elog.Close()
		elog.Info(1, fmt.Sprintf("Running in a non-interactive mode"))
		run := svc.Run
		err = run("Subutai P2P", &P2PService{})
		if err != nil {
			elog.Info(1, fmt.Sprintf("Failed to run service: %s", err))
			return nil
		}
	}
	return nil
}
