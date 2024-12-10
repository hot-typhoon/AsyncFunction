package handler

import (
	"net/http"

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
		output, err := s.CombinedOutput(bodyParams.Command)
		if err != nil {
			return "", err
		}

		return string(output), nil
	})

	if err != nil {
		util.HttpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpResponse(w, http.StatusOK, output)
}
