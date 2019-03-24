package gotezos

import (
	"encoding/json"
	"testing"
)

func TestNewRPCClient(t *testing.T) {
	var cases = []struct {
		host string
		port string
	}{
		{"localhost", "8732"},
		{"mainnet-node.tzscan.io", "80"},
	}

	for _, c := range cases {
		gtClient := NewTezosRPCClient(c.host, c.port)
		gt := NewGoTezos()
		gt.AddNewClient(gtClient)

		if !gtClient.Healthcheck() {
			t.Logf("Unable to query RPC on '%s:%s'. Check that a node is accessible.", c.host, c.port)
		}
	}
}

func TestCreateWalletWithMnemonic(t *testing.T) {
	t.Log("Create new wallet using Alphanet faucet account")

	var cases = []struct {
		mnemonic string
		password string
		email    string
		result   Wallet
	}{
		{"normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout",
			"PYh8nXDQLB",
			"vksbjweo.qsrgfvbw@tezos.example.org",
			Wallet{Sk: "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
				Pk:      "edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm",
				Address: "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
			},
		},
	}

	gt := NewGoTezos()

	for _, c := range cases {
		myWallet, err := gt.CreateWallet(c.mnemonic, c.email+c.password)
		if err != nil {
			t.Errorf("Unable to create wallet from Mnemonic: %s", err)
		}

		if myWallet.Address != c.result.Address || myWallet.Pk != c.result.Pk || myWallet.Sk != c.result.Sk {
			t.Errorf("Created wallet values do not match known answers")
		}
	}
}

func TestImportWalletFullSk(t *testing.T) {

	t.Log("Import existing wallet using complete secret key")

	var cases = []struct {
		pkh string
		pk  string
		sk  string
	}{
		{
			"tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
			"edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
			"edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
		},
	}

	gt := NewGoTezos()

	for _, c := range cases {
		myWallet, err := gt.ImportWallet(c.pkh, c.pk, c.sk)
		if err != nil {
			t.Errorf("%s", err)
		}

		if myWallet.Address != c.pkh || myWallet.Pk != c.pk || myWallet.Sk != c.sk {
			t.Errorf("Created wallet values do not match known answers")
		}
	}
}

func TestImportWalletSeedSk(t *testing.T) {

	t.Log("Import existing wallet using seed-secret key")

	var cases = []struct {
		pkh string
		pk  string
		sk  string
		sks string
	}{
		{
			"tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
			"edpkunwa7a3Y5vDr9eoKy4E21pzonuhqvNjscT9XG27aQV4gXq4dNm",
			"edskRjBSseEx9bSRSJJpbypJe5ZXucTtApb6qjechMB1BzEYwcEZyfLooo22Nwk33mPPJ3xZniFoa3o8Js7nNXDdqK9nNjFDi7",
			"edsk362Ypv3qLgbnGvZK7JwqNbwiLGe18XhTMFQY4gUonqnaCPiT6X",
		},
	}

	gt := NewGoTezos()

	for _, c := range cases {
		myWallet, err := gt.ImportWallet(c.pkh, c.pk, c.sks)
		if err != nil {
			t.Errorf("%s", err)
		}

		if myWallet.Address != c.pkh || myWallet.Pk != c.pk || myWallet.Sk != c.sk {
			t.Errorf("Created wallet values do not match known answers")
		}
	}
}

func TestImportEncryptedSecret(t *testing.T) {

	t.Log("Import wallet using password and encrypted key")

	var cases = []struct {
		pw  string
		sk  string
		pk  string
		pkh string
	}{
		{
			"password12345##",
			"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
			"edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5",
			"tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
		},
	}

	gt := NewGoTezos()

	for _, c := range cases {

		myWallet, err := gt.ImportEncryptedWallet(c.pw, c.sk)
		if err != nil {
			t.Errorf("%s", err)
		}

		if myWallet.Address != c.pkh || myWallet.Pk != c.pk {
			t.Errorf("Imported encrypted wallet does not match known answers")
		}
	}
}
func TestGetSnapShot(t *testing.T) {
	var cases = []struct {
		in  int
		out SnapShot
	}{
		{15, SnapShot{Cycle: 15, AssociatedBlock: 34048, AssociatedHash: "BLQ8HALSaSMaDPASYAG4tCBXrTinpfzhMZ7uD2JJ8k2zvxDrzEQ"}},
		{20, SnapShot{Cycle: 20, AssociatedBlock: 54272, AssociatedHash: "BM7AshJjzA9vDNDMbwDiGYfHPdWawy9nZCdabCttuCWZDA72SqL"}},
		{80, SnapShot{Cycle: 80, AssociatedBlock: 301568, AssociatedHash: "BMbeaq7EGPQtjS4Pr4kS6F8s4sSWZn9YQ9DJtzE6BCespEkmL7H"}},
	}

	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	for _, c := range cases {
		snapshot, err := gt.GetSnapShot(c.in)
		if err != nil {
			t.Error(err)
		}
		if c.out.AssociatedBlock != snapshot.AssociatedBlock || c.out.AssociatedHash != snapshot.AssociatedHash || c.out.Cycle != snapshot.Cycle {
			t.Errorf("Snap Shot %v, does not match the snapshot queryied: %v", c.out, snapshot)
		}
	}
}

func TestGetChainHead(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetChainHead()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetNetworkConstants(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetNetworkConstants()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetNetworkVersions(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetNetworkVersions()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBranchProtocol(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "443")
	gt.AddNewClient(client)

	_, err := gt.GetNetworkVersions()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBranchHash(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetBranchHash()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockLevelHead(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, _, err := gt.GetBlockLevelHead()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockAtLevel(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetBlockAtLevel(100000)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockByHash(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetBlockByHash("BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockOperationHashesHead(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetBlockOperationHashesHead()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockOperationHashesAtLevel(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetBlockOperationHashesAtLevel(100000)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockOperationHashes(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetBlockOperationHashes("BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetAccountBalanceAtSnapshot(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetAccountBalanceAtSnapshot("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 15)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetAccountBalance(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetAccountBalance("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetDelegateStakingBalance(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetDelegateStakingBalance("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 15)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetCurrentCycle(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetCurrentCycle()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetAccountBalanceAtBlock(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetDelegateStakingBalance("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 15)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetChainID(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)
	_, err := gt.GetChainID()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetDelegationsForDelegate(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	_, err := gt.GetDelegationsForDelegate("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetDelegationsForDelegateByCycle(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("mainnet-node.tzscan.io", "80")
	gt.AddNewClient(client)

	t.Log("Getting delegations for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycle 60")
	delegations, err := gt.GetDelegationsForDelegateByCycle("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(delegations)
}

func TestGetRewardsForDelegateForCycles(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting rewards for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycles 60-64")
	rewards, err := gt.GetRewardsForDelegateForCycles("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60, 64)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(rewards))
}

func TestGetRewardsForDelegateCycle(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting rewards for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycle 60")
	rewards, err := gt.GetRewardsForDelegateCycle("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(rewards))
}

func TestGetCycleRewards(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting rewards for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycles 60")
	rewards, err := gt.GetCycleRewards("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(rewards))
}

func TestGetDelegateRewardsForCycle(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting the rewards earned by delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for a cycle 60.")
	rewards, err := gt.GetDelegateRewardsForCycle("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(rewards))
}

func TestGetContractRewardsForDelegate(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	rewards, err := gt.GetDelegateRewardsForCycle("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log("Getting gross rewards and share for all delegations for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	contractRewards, err := gt.getContractRewardsForDelegate("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", rewards, 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(contractRewards))
}

func TestGetShareOfContract(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting the the share for delegation KT1EidADxWfYeBgK8L1ZTbf7a9zyjKwCFjfH on tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycle 60")
	share, _, err := gt.GetShareOfContract("KT1EidADxWfYeBgK8L1ZTbf7a9zyjKwCFjfH", "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(share))
}

func TestGetDelegate(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting info for delegate on tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc head block")
	delegate, err := gt.GetDelegate("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(delegate))
}

func TestGetStakingBalanceAtCycle(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting the staking balance for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc at cycle 60")
	stakingBalance, err := gt.GetStakingBalanceAtCycle("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(stakingBalance))
}

func TestGetBakingRights(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all baking rights for cycle 60")
	bakingRights, err := gt.GetBakingRights(60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(bakingRights))
}

func TestGetBakingRightsForDelegate(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting baking rights with priotrity 2 for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycle 60")
	bakingRights, err := gt.GetBakingRightsForDelegate(60, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 2)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(bakingRights))
}

func TestGetBakingRightsForDelegateForCycles(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting baking rights with priotrity 2 for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycles 60-64")
	bakingRights, err := gt.GetBakingRightsForDelegateForCycles(60, 64, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 2)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(bakingRights))
}

func TestGetEndorsingRightsForDelegate(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all endorsing rights for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycles 60")
	endorsingRights, err := gt.GetEndorsingRightsForDelegate(60, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(endorsingRights))
}

func TestGetEndorsingRightsForDelegateForCycles(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all endorsing rights for delegate tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc for cycles 60-64")
	endorsingRights, err := gt.GetEndorsingRightsForDelegateForCycles(60, 64, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(endorsingRights))
}

func TestGetEndorsingRights(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all endorsing rights for cycles 60")
	endorsingRights, err := gt.GetEndorsingRights(60)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(endorsingRights))
}

func TestGetAllDelegatesByHash(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all delegates at block hash BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	delegates, err := gt.GetAllDelegatesByHash("BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(delegates))
}

func TestGetAllDelegates(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all delegates at head")
	delegates, err := gt.GetAllDelegates()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(delegates))
}

//Takes an interface v and returns a pretty json string.
func PrettyReport(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return ""
}
