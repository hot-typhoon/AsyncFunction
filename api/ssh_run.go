package handler

import (
	"context"
	"net/http"
	"time"

	"asynclab.club/AsyncFunction/pkg/lib/ssh_run"
	"asynclab.club/AsyncFunction/pkg/util"
	"golang.org/x/crypto/ssh"
)

func HandlerSSHRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.HttpResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	bodyParams, err := util.ReadParamsFromBody[ssh_run.BodyParams](r.Body)
	if err != nil {
		util.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := ssh_run.ConsumeSession(bodyParams.Target, bodyParams.Jumpers, func(s *ssh.Session) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		var output []byte
		var err error

		done := make(chan error, 1)
		go func() {
			output, err = s.CombinedOutput(bodyParams.Command)
			done <- err
		}()

		select {
		case <-ctx.Done():
			cancel()
			signalErr := s.Signal(ssh.SIGKILL)
			if signalErr != nil {
				return "", signalErr
			}
			return "", ctx.Err()
		case err := <-done:
			cancel()
			if err != nil {
				return string(output), err
			}
		}

		return string(output), nil
	})

	if err != nil {
		util.HttpResponse(w, http.StatusInternalServerError, "Output: "+output+"; Error: "+err.Error())
		return
	}

	util.HttpResponse(w, http.StatusOK, output)
}
