package goTezos

import (
	"testing"
)

func TestNewRPCClient(t *testing.T) {

	t.Log("RPC client connecting to a Tezos node over localhost, port 8732")

	gtClient := NewTezosRPCClient("localhost", "8732")
	gt := NewGoTezos()
	gt.AddNewClient(gtClient)

	if !gtClient.Healthcheck() {
		t.Errorf("Unable to query RPC on 'localhost:8732'. Check that a node is accessible.")
	}
}

func TestNewWebClient(t *testing.T) {

	t.Log("Web-based RPC client using https://rpc.tzbeta.net")

	gtClient := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gtClient.IsWebClient(true)

	gt := NewGoTezos()
	gt.AddNewClient(gtClient)

	if !gtClient.Healthcheck() {
		t.Errorf("Unable to query RPC at 'https://rpc.tzbeta.net'.")
	}
}

func TestCreateWalletWithMnemonic(t *testing.T) {

	gt := NewGoTezos()

	t.Log("Create new wallet using Alphanet faucet account")

	mnemonic := "normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout"
	password := "PYh8nXDQLB"
	email := "vksbjweo.qsrgfvbw@tezos.example.org"

	// These values were gathered after manually importing above mnemonic into CLI wallet
	pkh := "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1"
	pk := "edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm"
	sk := "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V"

	// Alphanet 'password' is email & password concatenated together
	myWallet, err := gt.CreateWallet(mnemonic, email+password)
	if err != nil {
		t.Errorf("Unable to create wallet from Mnemonic: %s", err)
	}

	if myWallet.Address != pkh || myWallet.Pk != pk || myWallet.Sk != sk {
		t.Errorf("Created wallet values do not match known answers")
	}
}

func TestImportWalletFullSk(t *testing.T) {

	gt := NewGoTezos()

	t.Log("Import existing wallet using complete secret key")

	pkh := "tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK"
	pk := "edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G"
	sk := "edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR"

	myWallet, err := gt.ImportWallet(pkh, pk, sk)
	if err != nil {
		t.Errorf("%s", err)
	}

	if myWallet.Address != pkh || myWallet.Pk != pk || myWallet.Sk != sk {
		t.Errorf("Created wallet values do not match known answers")
	}
}

func TestImportWalletSeedSk(t *testing.T) {

	gt := NewGoTezos()

	t.Log("Import existing wallet using seed-secret key")

	pkh := "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ"
	pk := "edpkunwa7a3Y5vDr9eoKy4E21pzonuhqvNjscT9XG27aQV4gXq4dNm"
	sks := "edsk362Ypv3qLgbnGvZK7JwqNbwiLGe18XhTMFQY4gUonqnaCPiT6X"
	sk := "edskRjBSseEx9bSRSJJpbypJe5ZXucTtApb6qjechMB1BzEYwcEZyfLooo22Nwk33mPPJ3xZniFoa3o8Js7nNXDdqK9nNjFDi7"

	myWallet, err := gt.ImportWallet(pkh, pk, sks)
	if err != nil {
		t.Errorf("%s", err)
	}

	if myWallet.Address != pkh || myWallet.Pk != pk || myWallet.Sk != sk {
		t.Errorf("Created wallet values do not match known answers")
	}
}

func TestImportEncryptedSecret(t *testing.T) {

	gt := NewGoTezos()

	t.Log("Import wallet using password and encrypted key")

	pw := "password12345##"
	sk := "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2"

	// known answers for testing
	pk := "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5"
	pkh := "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo"

	myWallet, err := gt.ImportEncryptedWallet(pw, sk)
	if err != nil {
		t.Errorf("%s", err)
	}

	if myWallet.Address != pkh || myWallet.Pk != pk {
		t.Errorf("Imported encrypted wallet does not match known answers")
	}
}

func TestGetSnapShot(t *testing.T) {

	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting snapshot 15 from the network")

	snapshot, err := gt.GetSnapShot(15)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(snapshot))
}

func TestGetAllCurrentSnapShots(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all current snapshots")
	snapshots, err := gt.GetAllCurrentSnapShots()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(snapshots))
}

func TestGetChainHead(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all current snapshots")
	head, err := gt.GetChainHead()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(head))
}

func TestGetNetworkConstants(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting network constantss")
	netConts, err := gt.GetNetworkConstants()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(netConts))
}

func TestGetNetworkVersions(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting network versions")
	netVers, err := gt.GetNetworkVersions()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(netVers))
}

func TestGetBranchProtocol(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting branch protocol")
	brProto, err := gt.GetNetworkVersions()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(brProto))
}

func TestGetBranchHash(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting branch hash")
	brHash, err := gt.GetBranchHash()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(brHash))
}

func TestGetBlockLevelHead(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting branch hash")
	levelHead, levelHash, err := gt.GetBlockLevelHead()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(levelHead))
	t.Log(PrettyReport(levelHash))
}

func TestGetBlockAtLevel(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting block at level 100,000")
	block, err := gt.GetBlockAtLevel(100000)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(block))
}

func TestGetBlockByHash(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting block at BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	block, err := gt.GetBlockByHash("BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(block))
}

func TestGetBlockOperationHashesHead(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all operation hashes at head block")
	opHashes, err := gt.GetBlockOperationHashesHead()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(opHashes))
}

func TestGetBlockOperationHashesAtLevel(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all operation hashes at level 100000")
	opHashes, err := gt.GetBlockOperationHashesAtLevel(100000)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(opHashes))
}

func TestGetBlockOperationHashes(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting all operation hashes at hash BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	opHashes, err := gt.GetBlockOperationHashes("BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(opHashes))
}

func TestGetAccountBalanceAtSnapshot(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting account balance for tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc, at cycle 15")
	balance, err := gt.GetAccountBalanceAtSnapshot("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 15)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(balance))
}

func TestGetAccountBalance(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting account balance for tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	balance, err := gt.GetAccountBalance("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(balance))
}

func TestGetDelegateStakingBalance(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting staking balance for tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc at cycle 15")
	balance, err := gt.GetDelegateStakingBalance("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 15)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(balance))
}

func TestGetCurrentCycle(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting the current cycle")
	cycle, err := gt.GetCurrentCycle()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(cycle))
}

func TestGetAccountBalanceAtBlock(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting staking balance for tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc at cycle 15")
	balance, err := gt.GetDelegateStakingBalance("tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc", 15)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(balance))
}

func TestGetChainId(t *testing.T) {
	gt := NewGoTezos()
	client := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gt.AddNewClient(client)

	t.Log("Getting the chain ID for the network")
	chainId, err := gt.GetChainId()
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(chainId))
}
