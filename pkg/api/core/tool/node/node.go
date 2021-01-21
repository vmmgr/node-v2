package node

import (
	"encoding/json"
	"github.com/vmmgr/controller/pkg/api/core/controller"
	"github.com/vmmgr/node/pkg/api/core/gateway"
	"github.com/vmmgr/node/pkg/api/core/tool/client"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
)

func SendServer(input gateway.Info, code uint, progress uint, data string, error error) {
	for _, srv := range config.Conf.Controller.List {

		var comment string
		var status bool
		if error != nil {
			status = false
			comment = error.Error()
		} else {
			status = true
			comment = data
		}

		sendBody, _ := json.Marshal(controller.Node{
			GroupID:  input.GroupID,
			UUID:     input.UUID,
			Code:     code,
			Progress: progress,
			Status:   status,
			Comment:  comment,
		})
		client.Post(srv.URL, sendBody)
	}
}
