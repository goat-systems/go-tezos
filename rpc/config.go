package rpc

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

/*
UserActivatedProtocolOverrides represents user activated protocl overrides on the Tezos network.

RPC:
	/config/network/user_activated_protocol_overrides (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-config-network-user-activated-protocol-overrides
*/
type UserActivatedProtocolOverrides struct {
	ReplacedProtocol    string `json:"replaced_protocol"`
	ReplacementProtocol string `json:"replacement_protocol"`
}

/*
UserActivatedProtocolOverrides list of protocols which replace other protocols.

Path:
	/config/network/user_activated_protocol_overrides (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-config-network-user-activated-protocol-overrides
*/
func (c *Client) UserActivatedProtocolOverrides() (*resty.Response, UserActivatedProtocolOverrides, error) {
	resp, err := c.get("/config/network/user_activated_protocol_overrides")
	if err != nil {
		return resp, UserActivatedProtocolOverrides{}, errors.Wrap(err, "failed to get blocks")
	}

	var userActivatedProtocolOverride UserActivatedProtocolOverrides
	err = json.Unmarshal(resp.Body(), &userActivatedProtocolOverride)
	if err != nil {
		return resp, userActivatedProtocolOverride, errors.Wrap(err, "failed to unmarshal blocks")
	}

	return resp, userActivatedProtocolOverride, nil
}
