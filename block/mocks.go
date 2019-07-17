package block

var (
	goldenBlock = []byte(`{
		"chain_id": "NetXdQprcVkpaWU",
		"hash": "BLTGSUUjDpaHe7BYZa1zsrccJ7skurNiHZ1mpCz3cak9GnDfRoT",
		"header": {
			"context": "CoVpxj7gxbSQhX4J21HLaqFJJAavinyq3VmZwdgV3TtTrJYwxMXd",
			"fitness": [
				"00",
				"0000000000fb7afb"
			],
			"level": 524067,
			"operations_hash": "LLoaz5CZ4VDwJKUwd3SwcNzVz91NKy6jCNJG4jJqAvLxESjmvxoyc",
			"predecessor": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
			"priority": 0,
			"proof_of_work_nonce": "00000003ec49cde9",
			"proto": 4,
			"signature": "sigUgGKVNfJGCePA1jCkcfvTiAo4z9aVELXkWxEECPVjLJorNCeSmVYKwFeqTDRFFkHvdH6QV6PoksJYVk7Sy6rgR1UEBmUZ",
			"timestamp": "2019-07-16T14:59:56Z",
			"validation_pass": 4
		},
		"metadata": {
			"baker": "tz1VQnqCCqX4K5sP3FNkVSNKTdCAMJDd3E1n",
			"balance_updates": [
				{
					"change": "-512000000",
					"contract": "tz1VQnqCCqX4K5sP3FNkVSNKTdCAMJDd3E1n",
					"kind": "contract"
				},
				{
					"category": "deposits",
					"change": "512000000",
					"cycle": 127,
					"delegate": "tz1VQnqCCqX4K5sP3FNkVSNKTdCAMJDd3E1n",
					"kind": "freezer"
				},
				{
					"category": "rewards",
					"change": "16000000",
					"cycle": 127,
					"delegate": "tz1VQnqCCqX4K5sP3FNkVSNKTdCAMJDd3E1n",
					"kind": "freezer"
				}
			],
			"consumed_gas": "0",
			"deactivated": [],
			"level": {
				"cycle": 127,
				"cycle_position": 3874,
				"expected_commitment": false,
				"level": 524067,
				"level_position": 524066,
				"voting_period": 15,
				"voting_period_position": 32546
			},
			"max_block_header_length": 238,
			"max_operation_data_length": 16384,
			"max_operation_list_length": [
				{
					"max_op": 32,
					"max_size": 32768
				},
				{
					"max_size": 32768
				},
				{
					"max_op": 132,
					"max_size": 135168
				},
				{
					"max_size": 524288
				}
			],
			"max_operations_ttl": 60,
			"next_protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
			"nonce_hash": null,
			"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
			"test_chain_status": {
				"status": "not_running"
			},
			"voting_period_kind": "testing_vote"
		},
		"operations": [
			[
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1hx8hMmmeyDBi6WJgpKwK4n5S2qAEpavx2",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1hx8hMmmeyDBi6WJgpKwK4n5S2qAEpavx2",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1hx8hMmmeyDBi6WJgpKwK4n5S2qAEpavx2",
										"kind": "freezer"
									}
								],
								"delegate": "tz1hx8hMmmeyDBi6WJgpKwK4n5S2qAEpavx2",
								"slots": [
									28
								]
							}
						}
					],
					"hash": "onzhYXdxBcpNmMBSDoDnpeToNMTaw6G9Btdm8bKsrdYqeCsR1Tk",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigPEkymbx5vMX6CSsEh1vvDS4uUqESRSMUWHRLqvZ97FsoqvT1hPjXubaJypnnehmzUnUHZR3aHLLY4Byng1fWnXinewn1J"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ",
										"kind": "freezer"
									}
								],
								"delegate": "tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ",
								"slots": [
									13
								]
							}
						}
					],
					"hash": "op5NAznMhAibhpzbArakyxczZXiQqbf4jNuhjdDwPL1xnfenYpE",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigP2ku4WzVRGbjQAH4rp5s8b25aVAUJcaL7L1DLUijUkhYqDNR96wvBWVGCgzNNvvztBGVpq1XNVdUVfxuMW3uGv3AHrvJD"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-192000000",
										"contract": "tz3UoffC7FG7zfpmvmjUmUeAaHvzdcUvAj6r",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "192000000",
										"cycle": 127,
										"delegate": "tz3UoffC7FG7zfpmvmjUmUeAaHvzdcUvAj6r",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "6000000",
										"cycle": 127,
										"delegate": "tz3UoffC7FG7zfpmvmjUmUeAaHvzdcUvAj6r",
										"kind": "freezer"
									}
								],
								"delegate": "tz3UoffC7FG7zfpmvmjUmUeAaHvzdcUvAj6r",
								"slots": [
									31,
									26,
									23
								]
							}
						}
					],
					"hash": "oosrsi9kuJuUxq2FCX1Rp4J5ZK1mQXM1JycdyQVjo7AqrpJoeBX",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigw27ZZFQhyU7TPmeWDhXbdmv5zXp5N2QDwMm2nTv45jhSjTsHL8TuZoYpF2bJqkQtZFhHzLzHCr4rNFSygweFR3eoR8U8w"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-128000000",
										"contract": "tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "128000000",
										"cycle": 127,
										"delegate": "tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "4000000",
										"cycle": 127,
										"delegate": "tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5",
										"kind": "freezer"
									}
								],
								"delegate": "tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5",
								"slots": [
									29,
									22
								]
							}
						}
					],
					"hash": "oo684swNgfoybeqEzSkumSy8K44o3TpQ97Jy9bWNW7U3xPNsqdo",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigf5oAkWBeKq5ZrYDT9YoqKKAkLxem2WRWt4wKB2HrYMrW29uGfttNvTBAk55s4RWXA4YBeooJ6nwyJGZ1fsYwpd45an3MK"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-192000000",
										"contract": "tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "192000000",
										"cycle": 127,
										"delegate": "tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "6000000",
										"cycle": 127,
										"delegate": "tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE",
										"kind": "freezer"
									}
								],
								"delegate": "tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE",
								"slots": [
									19,
									15,
									12
								]
							}
						}
					],
					"hash": "opEaHx9ryDGfXSHGeDDEjVPdnC776SWwFPiwS5aMq2fP8mYS2ct",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigcDz9osLYxxCP6vU7BQhwPryhcKtnnJLdZDJmjQpFEGKaQ7uSz1wWGCoteR76M5goBG8mhchvbFicvn4xdHtkqYzPTKWDb"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
										"kind": "freezer"
									}
								],
								"delegate": "tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
								"slots": [
									8
								]
							}
						}
					],
					"hash": "onj2BfBPG9uYyrDZ4FmGXsfrwGhw1QiE1YWF1aBzwdRbNxvJbtS",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sige6TKCzDhGapiZwpXGz53iByxBgUTHD6XWgjmkp8M1Zh7fX2aXFUPcshW1fEyXAUk7P2XEG3wEjHAofZ4zPWZMUAefLE2P"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-128000000",
										"contract": "tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "128000000",
										"cycle": 127,
										"delegate": "tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "4000000",
										"cycle": 127,
										"delegate": "tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9",
										"kind": "freezer"
									}
								],
								"delegate": "tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9",
								"slots": [
									9,
									1
								]
							}
						}
					],
					"hash": "opNtvju5XXUsBTeBrRasyVenM8hx5QFMFrkKcxhA66UsAXfEQCt",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigpHU7cdZ3QF2My6rfjbBNDMC84H8sAq29ZDi3x4FEez8TLXT2CZjurUxTmE5LDVhr1AceRsidQqkeY5gzzoiG9Jt1VHCzw"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-192000000",
										"contract": "tz1Ldzz6k1BHdhuKvAtMRX7h5kJSMHESMHLC",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "192000000",
										"cycle": 127,
										"delegate": "tz1Ldzz6k1BHdhuKvAtMRX7h5kJSMHESMHLC",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "6000000",
										"cycle": 127,
										"delegate": "tz1Ldzz6k1BHdhuKvAtMRX7h5kJSMHESMHLC",
										"kind": "freezer"
									}
								],
								"delegate": "tz1Ldzz6k1BHdhuKvAtMRX7h5kJSMHESMHLC",
								"slots": [
									20,
									14,
									10
								]
							}
						}
					],
					"hash": "ooMi6cMvKkR3cmSn2XmxsAKc5ahvPidZarWsKji3zrbEf4vf5eQ",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigbAAJzTaKGTz8RmjLn5fXMLmtKzUEyStGdyuHUUAcPNLNCzaUcW1N8ubbNpnoKFMzKj7ypQiVcPuZqU772nyRiQiLaMsg3"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-128000000",
										"contract": "tz1TNWtofRofCU11YwCNwTMWNFBodYi6eNqU",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "128000000",
										"cycle": 127,
										"delegate": "tz1TNWtofRofCU11YwCNwTMWNFBodYi6eNqU",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "4000000",
										"cycle": 127,
										"delegate": "tz1TNWtofRofCU11YwCNwTMWNFBodYi6eNqU",
										"kind": "freezer"
									}
								],
								"delegate": "tz1TNWtofRofCU11YwCNwTMWNFBodYi6eNqU",
								"slots": [
									30,
									18
								]
							}
						}
					],
					"hash": "oo8UEx9zf6TeBq2qEgLE5ye3sDRHcUFjNcAYRGvcSD3VNGtSDoz",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigkWB2RnMTSw5eVGggw5PnBQMaAyfPTrBbPAJydFMenVFbUxydWw32KpB8kVQ3UUbts2LX5iZrz7PsYFkJBr62EaLAXvjqD"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1KtvGSYU5hdKD288a1koTBURWYuADJGrLE",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1KtvGSYU5hdKD288a1koTBURWYuADJGrLE",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1KtvGSYU5hdKD288a1koTBURWYuADJGrLE",
										"kind": "freezer"
									}
								],
								"delegate": "tz1KtvGSYU5hdKD288a1koTBURWYuADJGrLE",
								"slots": [
									7
								]
							}
						}
					],
					"hash": "ooD8414eFLGiqVWrteEdD7djF63RSRVe82J75wNth5heTN71FP2",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigsvszfQKvgUsUeSXti1JTzfj7nLc6hnudCmtrroVUjLrCdG7Uac1YHt1SZqAHGpvAp1WSWuhovHshi3e245CJ6cmRiWREQ"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz3bTdwZinP8U1JmSweNzVKhmwafqWmFWRfk",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz3bTdwZinP8U1JmSweNzVKhmwafqWmFWRfk",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz3bTdwZinP8U1JmSweNzVKhmwafqWmFWRfk",
										"kind": "freezer"
									}
								],
								"delegate": "tz3bTdwZinP8U1JmSweNzVKhmwafqWmFWRfk",
								"slots": [
									3
								]
							}
						}
					],
					"hash": "onsuHbfyDkYEZoKCu71CsaRb6f46aNm5R56ngUt35Nj2rwrFKzX",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigQfGGzByTDRb3hFRB8bUTbupTtUsZo12CuXya281LYGvKU2dSc24rZuwXCshNiaERfwaGqt9R95nt4d8mVTiqk6u9ceDoq"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1TDSmoZXwVevLTEvKCTHWpomG76oC9S2fJ",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1TDSmoZXwVevLTEvKCTHWpomG76oC9S2fJ",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1TDSmoZXwVevLTEvKCTHWpomG76oC9S2fJ",
										"kind": "freezer"
									}
								],
								"delegate": "tz1TDSmoZXwVevLTEvKCTHWpomG76oC9S2fJ",
								"slots": [
									11
								]
							}
						}
					],
					"hash": "oopuuvdS8rHkt8guwqFYLgBvC8ryEPBYZs8HSJr1cQqCUasSG6T",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigd7WstkuLL6Mm24PryDL83Hs7jpMSaGDsR7zhx5pw81P2aVZU6cCryTJ62MYCuKkWUbBi2SVemDcCFHmSXmSw4C8SuGsoY"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-128000000",
										"contract": "tz1Tnjaxk6tbAeC2TmMApPh8UsrEVQvhHvx5",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "128000000",
										"cycle": 127,
										"delegate": "tz1Tnjaxk6tbAeC2TmMApPh8UsrEVQvhHvx5",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "4000000",
										"cycle": 127,
										"delegate": "tz1Tnjaxk6tbAeC2TmMApPh8UsrEVQvhHvx5",
										"kind": "freezer"
									}
								],
								"delegate": "tz1Tnjaxk6tbAeC2TmMApPh8UsrEVQvhHvx5",
								"slots": [
									2,
									0
								]
							}
						}
					],
					"hash": "ooiS6KdfVTPDb1B9poJ7nF1SKruvVMduMGeFomEPkXCbsyyjQs6",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigSeQy8siunUFFdFCVZpswMRGiCMVfkY8tgWJd7DZ9XH2ovpYbgBcdpGXuyFEi3CYxLq5PT4y3iBZUQqXQ2xkkC74eqJ6pS"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1bHzftcTKZMTZgLLtnrXydCm6UEqf4ivca",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1bHzftcTKZMTZgLLtnrXydCm6UEqf4ivca",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1bHzftcTKZMTZgLLtnrXydCm6UEqf4ivca",
										"kind": "freezer"
									}
								],
								"delegate": "tz1bHzftcTKZMTZgLLtnrXydCm6UEqf4ivca",
								"slots": [
									5
								]
							}
						}
					],
					"hash": "opDb9Yxq4wTS6TTbLe6kc86moc3q6PTeewVKYyqBkNKdxBXwe7A",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigfDa7g9qyoNU1TKqFDXYu2TkCRxCqA6xkKczhWrvnkoRdFyTKqRroyKxG2QrARZUCWUVcdpBf2Vmo2jsnqYutRHvMe3RdG"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1LmaFsWRkjr7QMCx5PtV6xTUz3AmEpKQiF",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1LmaFsWRkjr7QMCx5PtV6xTUz3AmEpKQiF",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1LmaFsWRkjr7QMCx5PtV6xTUz3AmEpKQiF",
										"kind": "freezer"
									}
								],
								"delegate": "tz1LmaFsWRkjr7QMCx5PtV6xTUz3AmEpKQiF",
								"slots": [
									27
								]
							}
						}
					],
					"hash": "onjp1wwPLQuJd8vMWVtYjU6bBDWDXP8mevwKbke9aXR5NKJVyz7",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigpBZJKNEepQg2wcu9X5na54Nu4ZU4YiH6fpsTs9H9eYjHFP45knvXvrCn9xLMHczUgckzDtMxybPWqZ1QBPVxR7LdKC1So"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs",
										"kind": "freezer"
									}
								],
								"delegate": "tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs",
								"slots": [
									16
								]
							}
						}
					],
					"hash": "ooQKVMB7yRCHWjnDykZkzuBs7KXhgf1aXRPycP47R56nwcUiKU9",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigcd8WP9bEG9tnfgJHzgRRUyUEnRV4hzAUcRKkidV4t93yHvZjTvCBVSM5V2EsjDh4T9o2VpkFD25ES1XSAcWeb9kdFUn9U"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV",
										"kind": "freezer"
									}
								],
								"delegate": "tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV",
								"slots": [
									6
								]
							}
						}
					],
					"hash": "opZJ7raAX3AXxp4J5frMmxBQGrz4y89Exazsp4JpR7KBN7ND64J",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigS8NmjWS83id1Y6cmgkcjscwGBTPcNwrPtNFS4h4eseK1WrYiVGimBH5x3yseM843WtkPfckBNB2NgTW6M6gDCfk3UgH3o"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-128000000",
										"contract": "tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "128000000",
										"cycle": 127,
										"delegate": "tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "4000000",
										"cycle": 127,
										"delegate": "tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K",
										"kind": "freezer"
									}
								],
								"delegate": "tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K",
								"slots": [
									25,
									17
								]
							}
						}
					],
					"hash": "oowyPgfzShPHxYFhC8YhmDXriFx2R1Ldez1z62r8QDcHHKNkzq9",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigmsNz5xcJ3JtAazRB9zkznBjAQf13S8tK9gbipCELaf8rEyx1K7Q2UghimabEjABAWDFuwAaN8k7pJrtwrM9ofbBNoDE3L"
				},
				{
					"branch": "BMdw66rEAHYSu1WRwpVehpWUrB2tdt8RmGRYEt5YT6vs63zuWPU",
					"chain_id": "NetXdQprcVkpaWU",
					"contents": [
						{
							"kind": "endorsement",
							"level": 524066,
							"metadata": {
								"balance_updates": [
									{
										"change": "-64000000",
										"contract": "tz1UHxJUMWHY4FxK3RxgbSdwMXAhEzmoLVWA",
										"kind": "contract"
									},
									{
										"category": "deposits",
										"change": "64000000",
										"cycle": 127,
										"delegate": "tz1UHxJUMWHY4FxK3RxgbSdwMXAhEzmoLVWA",
										"kind": "freezer"
									},
									{
										"category": "rewards",
										"change": "2000000",
										"cycle": 127,
										"delegate": "tz1UHxJUMWHY4FxK3RxgbSdwMXAhEzmoLVWA",
										"kind": "freezer"
									}
								],
								"delegate": "tz1UHxJUMWHY4FxK3RxgbSdwMXAhEzmoLVWA",
								"slots": [
									21
								]
							}
						}
					],
					"hash": "opHpcLKt5KgGRhzxc2RbzHk1mfSjp4A4wTUuQvbd5e5LmxSb3gv",
					"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
					"signature": "sigXGkVuguJJYuj1rMwJu3aLZqU13zFesedkAwM4TshyeMgjPU1aadifCdo3g76G4oQbchRQxfBSVfht1JV6rEMJvgEEsaTR"
				}
			],
			[],
			[],
			[]
		],
		"protocol": "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd"
	}`)
)

type client struct {
	ReturnBody []byte
}

func (c *client) Post(path, args string) ([]byte, error) {
	return c.ReturnBody, nil
}

func (c *client) Get(path string, params map[string]string) ([]byte, error) {
	return c.ReturnBody, nil
}
