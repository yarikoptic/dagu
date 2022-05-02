package main

import (
	"errors"
	"log"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yohamta/dagu/internal/config"
	"github.com/yohamta/dagu/internal/controller"
	"github.com/yohamta/dagu/internal/scheduler"
	"github.com/yohamta/dagu/internal/utils"
)

func newStopCommand() *cli.Command {
	cl := &config.Loader{
		HomeDir: utils.MustGetUserHomeDir(),
	}
	return &cli.Command{
		Name:  "stop",
		Usage: "dagu stop <config>",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("config file must be specified.")
			}
			config_file_path := c.Args().Get(0)
			cfg, err := cl.Load(config_file_path, "")
			if err != nil {
				return err
			}
			return stop(cfg)
		},
	}
}

func stop(cfg *config.Config) error {
	status, err := controller.New(cfg).GetStatus()
	if err != nil {
		return err
	}

	if status.Status != scheduler.SchedulerStatus_Running ||
		!status.Pid.IsRunning() {
		log.Printf("DAG is not running.")
		return nil
	}
	syscall.Kill(int(status.Pid), syscall.SIGTERM)
	for {
		time.Sleep(time.Second * 3)
		s, err := controller.New(cfg).GetStatus()
		if err != nil {
			return err
		}
		if s.Pid == status.Pid && s.Status ==
			scheduler.SchedulerStatus_Running {
			continue
		}
		break
	}
	log.Printf("DAG is stopped.")
	return nil
}
