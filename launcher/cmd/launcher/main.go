package main

import (
	"bufio"
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
)

func main() {

	argLen := len(os.Args) - 1
	if argLen%2 > 0 {
		log.Fatalln("Args should be by env,cmd pairs")
	}

	signalCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	errGr, errCtx := errgroup.WithContext(signalCtx)
	wg := sync.WaitGroup{}

	for i := 1; i < argLen; i += 2 {
		wg.Add(1)

		envPath := os.Args[i]
		cmdPath := os.Args[i+1]

		errGr.Go(func() error {
			defer func() {
				wg.Done()
			}()

			if err := fork(errCtx, envPath, cmdPath); err != nil {
				return err
			}

			return nil
		})
	}

	errExitCh := make(chan struct{})
	go func() {
		err := errGr.Wait()
		if err != nil {
			log.Println(err)
		}
		errExitCh <- struct{}{}
	}()

	select {
	case <-signalCtx.Done():
		break
	case <-errExitCh:
		break
	}

	wg.Wait()
}

func fork(ctx context.Context, envPath string, cmdPath string) error {
	file, err := os.Open(envPath)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	var env []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		env = append(env, scanner.Text())
	}

	cmd := exec.Command(cmdPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	exitCh := make(chan error, 1)
	defer close(exitCh)

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		select {
		case <-exitCh:
			break
		case <-ctx.Done():
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				log.Println(err)
			}
			break
		}
	}()

	err = cmd.Wait()
	exitCh <- err

	if err != nil {
		return err
	}

	return nil
}
