package goTezos
/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: This file contains structures used for the goTezos lib
License: MIT
*/

import "time"

/*
Description: A structure to hold a snapshot query in.
Cycle (int): The cycle the snapshot was taken in
Number (int): Snap shot number decided. Empty if decided == false
Decided (bool): true or false if snap shot is decided
AssociatedBlock: The block number the snapshot reference, only available if snapshot is decided.
*/
type SnapShot struct {
    Cycle int
    Number int
    AssociatedBlock int
}

/*
Description: A way to repesent each Delegated Contact, and their share for each cycle
Address: A string value representing the delegated contracts address
Commitments: An array of a structure that holds the amount commited for a cycle, and the percentage share
Delegator: Is this contract the delegator?
*/
type DelegatedContract struct {
    Address string //Public Key Hash
    Contracts []Contract //Percentage of total delegation for profit share for each cycle participated
    Delegate bool //If this client is yourself or not.
    TimeStamp time.Time
    TotalPayout float64
}

/*
Description: A representation of the amount commited in a cycle, and the percentage share for that amount.
Cycle: The cycle number
Amount: XTZ value of the amount commited in the cycle
SharePercentage: The percentage value of the amount to all commitments made in that cycle
Payout: Amount of rewards to be paid out for the commitment
Timestamp: A timestamp to show when the commitment was made
*/
type Contract struct {
  Cycle int
  Amount float64
  RollInclusion float64
  SharePercentage float64
  GrossPayout float64
  NetPayout float64
  Fee float64
}

/*
Description: A structure to represent a known address.
Address: a string value representing the address
Alias: the alias assigned to the known address
Sk: The protection around the Sk, unencrypted, legder, etc
*/
type KnownAddress struct {
  Address string
  Alias string
  Sk string
}
