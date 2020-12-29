package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

/*
BallotList represents a list of casted ballots in a block.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-ballot-list
*/
type BallotList []struct {
	PublicKeyHash string `json:"pkh"`
	Ballot        string `json:"ballot"`
}

/*
BallotList returns ballots casted so far during a voting period.

Path:
	../<block_id>/votes/ballot_list (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-ballot-list
*/
func (c *Client) BallotList(blockID BlockID) (*resty.Response, BallotList, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/ballot_list", c.chain, blockID.ID()))
	if err != nil {
		return nil, BallotList{}, errors.Wrapf(err, "failed to get ballot list")
	}

	var ballotList BallotList
	err = json.Unmarshal(resp.Body(), &ballotList)
	if err != nil {
		return resp, BallotList{}, errors.Wrapf(err, "failed to get ballot list: failed to parse json")
	}

	return resp, ballotList, nil
}

/*
Ballots represents a ballot total.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-ballots
*/
type Ballots struct {
	Yay  int `json:"yay"`
	Nay  int `json:"nay"`
	Pass int `json:"pass"`
}

/*
Ballots returns sum of ballots casted so far during a voting period.

Path:
	../<block_id>/votes/ballots (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-ballots
*/
func (c *Client) Ballots(blockID BlockID) (*resty.Response, Ballots, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/ballots", c.chain, blockID.ID()))
	if err != nil {
		return nil, Ballots{}, errors.Wrapf(err, "failed to get ballots")
	}

	var ballots Ballots
	err = json.Unmarshal(resp.Body(), &ballots)
	if err != nil {
		return resp, Ballots{}, errors.Wrapf(err, "failed to get ballots: failed to parse json")
	}

	return resp, ballots, nil
}

/*
VotingPeriod is the the voting period (index, kind, starting position) and related information (position, remaining) of the interrogated block.

Path:
	../<block_id>/votes/current_period (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-current-period
*/
type VotingPeriod struct {
	VotingPeriod struct {
		Index         int    `json:"index"`
		Kind          string `json:"kind"`
		StartPosition int    `json:"start_position"`
	} `json:"voting_period"`
	Position  int `json:"position"`
	Remaining int `json:"remaining"`
}

/*
CurrentPeriod returns the voting period (index, kind, starting position) and related information (position, remaining) of the interrogated block.

Path:
	../<block_id>/votes/current_period (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-current-period
*/
func (c *Client) CurrentPeriod(blockID BlockID) (*resty.Response, VotingPeriod, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/current_period", c.chain, blockID.ID()))
	if err != nil {
		return resp, VotingPeriod{}, errors.Wrapf(err, "failed to get current period")
	}

	var currentPeriod VotingPeriod
	err = json.Unmarshal(resp.Body(), &currentPeriod)
	if err != nil {
		return resp, VotingPeriod{}, errors.Wrapf(err, "failed to get current period: failed to parse json")
	}

	return resp, currentPeriod, nil
}

/*
CurrentPeriodKind returns the current period kind.

Path:
	../<block_id>/votes/current_period_kind (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-current-period-kind
*/
func (c *Client) CurrentPeriodKind(blockID BlockID) (*resty.Response, string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/current_period_kind", c.chain, blockID.ID()))
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get current period kind")
	}

	var currentPeriodKind string
	err = json.Unmarshal(resp.Body(), &currentPeriodKind)
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get current period kind: failed to parse json")
	}

	return resp, currentPeriodKind, nil
}

/*
CurrentProposal returns the current proposal under evaluation.

Path:
	../<block_id>/votes/current_proposal (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-current-proposal
*/
func (c *Client) CurrentProposal(blockID BlockID) (*resty.Response, string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/current_proposal", c.chain, blockID.ID()))
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get current proposal")
	}

	var currentProposal string
	err = json.Unmarshal(resp.Body(), &currentProposal)
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get current proposal: failed to parse json")
	}

	return resp, currentProposal, nil
}

/*
CurrentQuorum returns the current expected quorum.

Path:
	../<block_id>/votes/current_proposal (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-current-quorum
*/
func (c *Client) CurrentQuorum(blockID BlockID) (*resty.Response, int, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/current_quorum", c.chain, blockID.ID()))
	if err != nil {
		return resp, 0, errors.Wrapf(err, "failed to get current quorum")
	}

	var currentQuorum int
	err = json.Unmarshal(resp.Body(), &currentQuorum)
	if err != nil {
		return resp, 0, errors.Wrapf(err, "failed to get current quorum: failed to parse json")
	}

	return resp, currentQuorum, nil
}

/*
Listings represents a list of delegates with their voting weight, in number of rolls.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-listings
*/
type Listings []struct {
	PublicKeyHash string `json:"pkh"`
	Rolls         int    `json:"rolls"`
}

/*
Listings returns a list of delegates with their voting weight, in number of rolls.

Path:
	../<block_id>/votes/listings (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-listings
*/
func (c *Client) Listings(blockID BlockID) (*resty.Response, Listings, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/listings", c.chain, blockID.ID()))
	if err != nil {
		return resp, Listings{}, errors.Wrapf(err, "failed to get listings")
	}

	var listings Listings
	err = json.Unmarshal(resp.Body(), &listings)
	if err != nil {
		return resp, Listings{}, errors.Wrapf(err, "failed to get listings: failed to parse json")
	}

	return resp, listings, nil
}

/*
Proposals is the list of proposals with number of supporters.

RPC:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-proposals
*/
type Proposals []struct {
	Hash       string
	Supporters int
}

// UnmarshalJSON satisfies the json.Marshaler
func (p *Proposals) UnmarshalJSON(b []byte) error {
	var out [][]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return err
	}

	var proposals Proposals
	for _, x := range out {
		if len(x) != 2 {
			return errors.New("unexpected bytes")
		}

		hash := fmt.Sprintf("%v", x[0])
		supportersStr := fmt.Sprintf("%v", x[1])
		supporters, err := strconv.Atoi(supportersStr)
		if err != nil {
			return errors.New("unexpected bytes")
		}

		proposals = append(proposals, struct {
			Hash       string
			Supporters int
		}{
			Hash:       hash,
			Supporters: supporters,
		})
	}

	p = &proposals
	return nil
}

/*
Proposals returns a list of proposals with number of supporters.

Path:
	../<block_id>/votes/proposals (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-proposals
*/
func (c *Client) Proposals(blockID BlockID) (*resty.Response, Proposals, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/proposals", c.chain, blockID.ID()))
	if err != nil {
		return resp, Proposals{}, errors.Wrapf(err, "failed to get proposals")
	}

	var proposals Proposals
	err = json.Unmarshal(resp.Body(), &proposals)
	if err != nil {
		return resp, Proposals{}, errors.Wrapf(err, "failed to get proposals: failed to parse json")
	}

	return resp, proposals, nil
}

/*
SuccessorPeriod returns the voting period (index, kind, starting position) and related information (position, remaining) of the next block.

Path:
	../<block_id>/votes/successor_period (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-successor-period
*/
func (c *Client) SuccessorPeriod(blockID BlockID) (*resty.Response, VotingPeriod, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/proposals", c.chain, blockID.ID()))
	if err != nil {
		return resp, VotingPeriod{}, errors.Wrapf(err, "failed to get successor period")
	}

	var votingPeriod VotingPeriod
	err = json.Unmarshal(resp.Body(), &votingPeriod)
	if err != nil {
		return resp, VotingPeriod{}, errors.Wrapf(err, "failed to get successor period: failed to parse json")
	}

	return resp, votingPeriod, nil
}

/*
TotalVotingPower returns the total number of rolls for the delegates in the voting listings.

Path:
	../<block_id>/votes/current_proposal (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-votes-total-voting-power
*/
func (c *Client) TotalVotingPower(blockID BlockID) (*resty.Response, int, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/votes/total_voting_power", c.chain, blockID.ID()))
	if err != nil {
		return resp, 0, errors.Wrapf(err, "failed to get total voting power")
	}

	var votingPower int
	err = json.Unmarshal(resp.Body(), &votingPower)
	if err != nil {
		return resp, 0, errors.Wrapf(err, "failed to get total voting power: failed to parse json")
	}

	return resp, votingPower, nil
}
