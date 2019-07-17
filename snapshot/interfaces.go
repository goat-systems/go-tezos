package snapshot

type TezosSnapshotService interface {
	Get(cycle int) (Snapshot, error)
}
