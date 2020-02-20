package gotezos

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The below variables consist of mock RPC responses which enables proper unit testing.
var (
	mockActiveChainsResp    = []byte(`[{"chain_id":"NetXdQprcVkpaWU"}]`)
	mockBakingRightsResp    = []byte(`[{"level":732756,"delegate":"tz1LmaFsWRkjr7QMCx5PtV6xTUz3AmEpKQiF","priority":0,"estimated_time":"2019-12-12T11:27:11Z"},{"level":732756,"delegate":"tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN","priority":1,"estimated_time":"2019-12-12T11:27:51Z"},{"level":732756,"delegate":"tz1RCFbB9GpALpsZtu6J58sb74dm8qe6XBzv","priority":2,"estimated_time":"2019-12-12T11:28:31Z"},{"level":732756,"delegate":"tz3RB4aoyjov4KEVRbuhvQ1CKJgBJMWhaeB8","priority":3,"estimated_time":"2019-12-12T11:29:11Z"},{"level":732756,"delegate":"tz3e75hU4EhDU3ukyJueh5v6UvEHzGwkg3yC","priority":4,"estimated_time":"2019-12-12T11:29:51Z"},{"level":732756,"delegate":"tz1eEnQhbwf6trb8Q8mPb2RaPkNk2rN7BKi8","priority":5,"estimated_time":"2019-12-12T11:30:31Z"},{"level":732756,"delegate":"tz1WpeqFaBG9Jm73Dmgqamy8eF8NWLz9JCoY","priority":7,"estimated_time":"2019-12-12T11:31:51Z"},{"level":732756,"delegate":"tz1gk3TDbU7cJuiBRMhwQXVvgDnjsxuWhcEA","priority":8,"estimated_time":"2019-12-12T11:32:31Z"},{"level":732756,"delegate":"tz1LLNkQK4UQV6QcFShiXJ2vT2ELw449MzAA","priority":9,"estimated_time":"2019-12-12T11:33:11Z"},{"level":732756,"delegate":"tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs","priority":10,"estimated_time":"2019-12-12T11:33:51Z"},{"level":732756,"delegate":"tz2TSvNTh2epDMhZHrw73nV9piBX7kLZ9K9m","priority":11,"estimated_time":"2019-12-12T11:34:31Z"},{"level":732756,"delegate":"tz3bTdwZinP8U1JmSweNzVKhmwafqWmFWRfk","priority":13,"estimated_time":"2019-12-12T11:35:51Z"},{"level":732756,"delegate":"tz1irJKkXS2DBWkU1NnmFQx1c1L7pbGg4yhk","priority":14,"estimated_time":"2019-12-12T11:36:31Z"},{"level":732756,"delegate":"tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV","priority":15,"estimated_time":"2019-12-12T11:37:11Z"},{"level":732756,"delegate":"tz1Zhv3RkfU2pHrmaiDyxp7kFZpZrUCu1CiF","priority":17,"estimated_time":"2019-12-12T11:38:31Z"},{"level":732756,"delegate":"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q","priority":18,"estimated_time":"2019-12-12T11:39:11Z"},{"level":732756,"delegate":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","priority":19,"estimated_time":"2019-12-12T11:39:51Z"},{"level":732756,"delegate":"tz1SYq214SCBy9naR6cvycQsYcUGpBqQAE8d","priority":20,"estimated_time":"2019-12-12T11:40:31Z"},{"level":732756,"delegate":"tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE","priority":21,"estimated_time":"2019-12-12T11:41:11Z"},{"level":732756,"delegate":"tz1iJ4qgGTzyhaYEzd1RnC6duEkLBd1nzexh","priority":22,"estimated_time":"2019-12-12T11:41:51Z"},{"level":732756,"delegate":"tz1Ldzz6k1BHdhuKvAtMRX7h5kJSMHESMHLC","priority":23,"estimated_time":"2019-12-12T11:42:31Z"},{"level":732756,"delegate":"tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5","priority":25,"estimated_time":"2019-12-12T11:43:51Z"},{"level":732756,"delegate":"tz1c3Wh8gNMMsYwZd67JndQpYxdaaPUV27E7","priority":26,"estimated_time":"2019-12-12T11:44:31Z"},{"level":732756,"delegate":"tz1SohptP53wDPZhzTWzDUFAUcWF6DMBpaJV","priority":27,"estimated_time":"2019-12-12T11:45:11Z"},{"level":732756,"delegate":"tz1MXFrtZoaXckE41bjUCSjAjAap3AFDSr3N","priority":29,"estimated_time":"2019-12-12T11:46:31Z"},{"level":732756,"delegate":"tz1ZNWFe3LmEJYTydctcgD6a5Apemwdtimn4","priority":31,"estimated_time":"2019-12-12T11:47:51Z"},{"level":732756,"delegate":"tz1MecudVJnFZN5FSrriu8ULz2d6dDTR7KaM","priority":34,"estimated_time":"2019-12-12T11:49:51Z"},{"level":732756,"delegate":"tz1TNWtofRofCU11YwCNwTMWNFBodYi6eNqU","priority":41,"estimated_time":"2019-12-12T11:54:31Z"},{"level":732756,"delegate":"tz1NpWrAyDL9k2Lmnyxcgr9xuJakbBxdq7FB","priority":44,"estimated_time":"2019-12-12T11:56:31Z"},{"level":732756,"delegate":"tz1TcH4Nb3aHNDJ7CGZhU7jgAK1BkSP4Lxds","priority":48,"estimated_time":"2019-12-12T11:59:11Z"},{"level":732756,"delegate":"tz1bacP88iSnWHAVUBQShtE4ZnUGYHUpGVBM","priority":53,"estimated_time":"2019-12-12T12:02:31Z"},{"level":732756,"delegate":"tz1f3Re8iw6Pt3KMHAvyccHxDU3NuqL95axD","priority":54,"estimated_time":"2019-12-12T12:03:11Z"},{"level":732756,"delegate":"tz1Nn14BBsDULrPXtkM9UQeXaE4iqJhmqmK5","priority":55,"estimated_time":"2019-12-12T12:03:51Z"},{"level":732756,"delegate":"tz1P2Po7YM526ughEsRbY4oR9zaUPDZjxFrb","priority":56,"estimated_time":"2019-12-12T12:04:31Z"},{"level":732756,"delegate":"tz1VoUzbL6X4RdpxwaccV6eQkeqP6b9sUbYa","priority":59,"estimated_time":"2019-12-12T12:06:31Z"}]`)
	mockBlockResp           = []byte(`{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1","header":{"level":656939,"proto":5,"predecessor":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","timestamp":"2019-10-19T13:44:41Z","validation_pass":4,"operations_hash":"LLoaTVCGSnVwnzcSSLQevLieJJ5zmRLCFzRHsd2ZohF6FT8qdB3PD","fitness":["01","000000000000062b"],"context":"CoVVyoEZA2y25w1C51K6FgdkEKNuyWTmQXfUtb2C4XnzXtVhxVjM","priority":0,"proof_of_work_nonce":"756e6b6feb131fdb","signature":"siggkJPySMihG93d6WYXpsu1Jar1xYQLTuhQ1iVrmFDCBfGc6thnsCxgoyb6ZqqabpZ5jcoYvug6YZwG5Gjo4j5K1TsV3cqp"},"metadata":{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","next_protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","test_chain_status":{"status":"not_running"},"max_operations_ttl":60,"max_operation_data_length":16384,"max_block_header_length":238,"max_operation_list_length":[{"max_size":32768,"max_op":32},{"max_size":32768},{"max_size":135168,"max_op":132},{"max_size":524288}],"baker":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","level":{"level":656939,"level_position":656938,"cycle":160,"cycle_position":1578,"voting_period":20,"voting_period_position":1578,"expected_commitment":false},"voting_period_kind":"proposal","nonce_hash":null,"consumed_gas":"25492","deactivated":[],"balance_updates":[{"kind":"contract","contract":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","change":"-512000000"},{"kind":"freezer","category":"deposits","delegate":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","cycle":160,"change":"512000000"},{"kind":"freezer","category":"rewards","delegate":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","cycle":160,"change":"16000000"}]},"operations":[[{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"opX4UniMfWCHYn7w1buBPxK4yGCpPFssQE3VzBRf4sFCzBHxLuw","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1Kt4P8BCaP93AEV4eA7gmpRryWt5hznjCP","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1Kt4P8BCaP93AEV4eA7gmpRryWt5hznjCP","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1Kt4P8BCaP93AEV4eA7gmpRryWt5hznjCP","cycle":160,"change":"2000000"}],"delegate":"tz1Kt4P8BCaP93AEV4eA7gmpRryWt5hznjCP","slots":[15]}}],"signature":"sigvU29YNjSN8foVQRgqBYWS1wsqSGmWtnE6WrTiq9zfMxAnyrq7zJdPGVQiLUTTmqEzDjjRKRsdRnDbsnXUgc1afnqRApru"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onrJzin1urvfgNAsSZ9t7QsWvuKwDqqamKA6suy6d44CHjxK1KZ","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1P2Po7YM526ughEsRbY4oR9zaUPDZjxFrb","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1P2Po7YM526ughEsRbY4oR9zaUPDZjxFrb","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1P2Po7YM526ughEsRbY4oR9zaUPDZjxFrb","cycle":160,"change":"2000000"}],"delegate":"tz1P2Po7YM526ughEsRbY4oR9zaUPDZjxFrb","slots":[7]}}],"signature":"sigW7M2tPzGMcxLxHsniGJfEKWSrzF3MYzuFcDagwNNa2sYVqoncn6RsgUpKfJzzULFnPB9ux79o9A3UqhvQrMU3Uco7n7K1"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"opGi3dj16BQ5CdKonchZz5c2Pew2b1EGMTbc6XS3SdFgTx6JydP","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1LcuQHNVQEWP2fZjk1QYZGNrfLDwrT3SyZ","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1LcuQHNVQEWP2fZjk1QYZGNrfLDwrT3SyZ","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1LcuQHNVQEWP2fZjk1QYZGNrfLDwrT3SyZ","cycle":160,"change":"2000000"}],"delegate":"tz1LcuQHNVQEWP2fZjk1QYZGNrfLDwrT3SyZ","slots":[30]}}],"signature":"sigWUx3qYo6NZJkt1pRUKUTfmZUArU21xHcSgeu1ttBjgPyqJzhDry3EPYbVmxoKCg5LfB9cwH6BkTsSzT5YovarYSVsVN98"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"opR16Jxv3MLo8FeayHsr12mRjPAh8HepmSrNRH9Asv1mJHyAQXA","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1WBfwbT66FC6BTLexc2BoyCCBM9LG7pnVW","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1WBfwbT66FC6BTLexc2BoyCCBM9LG7pnVW","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1WBfwbT66FC6BTLexc2BoyCCBM9LG7pnVW","cycle":160,"change":"2000000"}],"delegate":"tz1WBfwbT66FC6BTLexc2BoyCCBM9LG7pnVW","slots":[31]}}],"signature":"sigkjQYN8boBH1fqFzEkRpp8DVtdF3q1XiUjXwauWMdaBQpaCQFajqsSDNpNxoJUzL4bXpi98xQkDq6v1FQYSVsb7y7C5d3y"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"op63EzCuyNgEbRYrCyaBYpxe4XDbMUMDy21JK5sFpZX66NCb4hR","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1Xsrfv6hn86fp88YfRs6xcKwt2nTqxVZYM","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1Xsrfv6hn86fp88YfRs6xcKwt2nTqxVZYM","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1Xsrfv6hn86fp88YfRs6xcKwt2nTqxVZYM","cycle":160,"change":"2000000"}],"delegate":"tz1Xsrfv6hn86fp88YfRs6xcKwt2nTqxVZYM","slots":[0]}}],"signature":"sigh3GuwNbR1gfJ5kkm8UUjq75io5w7n3tsBBSFyFk7kMUSrSSMgWLqhuERJdmRuuArS8QL8suwqFMX9fFkdtpakXTvnQoBL"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"op3H6PbKiH7znSasRjvHPYFcUU5N1KKdtmhWDN5j1eqTamBSwCw","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN","change":"-128000000"},{"kind":"freezer","category":"deposits","delegate":"tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN","cycle":160,"change":"128000000"},{"kind":"freezer","category":"rewards","delegate":"tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN","cycle":160,"change":"4000000"}],"delegate":"tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN","slots":[10,5]}}],"signature":"sigo1iC2uiczSDbyDSBAGFvibtzmE3PanroBM8pgytkvCkuJScWNNkgHXJ1YWUbMGkBy1L86TdhfrpAaPiYFGZFKVxp893cg"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onzDm5EDfWbio5M9sg29yQ2PfZvfe1YiuFqjjEDjELf9aMccCfA","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1WctSwL49BZtdWAQ75NY6vr4xsAx6oKdah","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1WctSwL49BZtdWAQ75NY6vr4xsAx6oKdah","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1WctSwL49BZtdWAQ75NY6vr4xsAx6oKdah","cycle":160,"change":"2000000"}],"delegate":"tz1WctSwL49BZtdWAQ75NY6vr4xsAx6oKdah","slots":[18]}}],"signature":"sigqBkdYf12g81rdyu1CEhJGd4rRkKpq72TU5v9GAacoj7s4NoUX85yJzFgcw2HPafuqr7ukVc8fJ7unjLzEChTbREg3Jwpv"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"opQCWzmFvREnBg96e5tdbrNV1znUU8uvUH63dWzLChKnTWSnTmc","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5","change":"-192000000"},{"kind":"freezer","category":"deposits","delegate":"tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5","cycle":160,"change":"192000000"},{"kind":"freezer","category":"rewards","delegate":"tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5","cycle":160,"change":"6000000"}],"delegate":"tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5","slots":[25,9,4]}}],"signature":"sigZnANU2ykmheh9QJKK8jnVACfbVUsA4zYsn1jjfquRUv9vJWRSszUwdEQxaMyRJN2zEGtuuUVXcaeTEnVP1x3zejU5rrrM"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onw3fUAz2nBRgLRchrZuoV4uEr7hVzixSMxMNfVed8DPncQ2F1u","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1NpWrAyDL9k2Lmnyxcgr9xuJakbBxdq7FB","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1NpWrAyDL9k2Lmnyxcgr9xuJakbBxdq7FB","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1NpWrAyDL9k2Lmnyxcgr9xuJakbBxdq7FB","cycle":160,"change":"2000000"}],"delegate":"tz1NpWrAyDL9k2Lmnyxcgr9xuJakbBxdq7FB","slots":[14]}}],"signature":"sigfH8WBcnQouthw9Tigb9Qb1qhqkP41n3E1pnx4fBGy5eMfTzi2wVPvrrqJmaG2WFbf4QPGhBF6D4pK8MJh7VEXHDVNS15H"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"ooHKFhhYfbppHpV9LbQNMW76zLG8PfNF3HVk7q2HSrFPtw9GdbB","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE","cycle":160,"change":"2000000"}],"delegate":"tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE","slots":[20]}}],"signature":"sigmZ3S4yKAyNFyKhRQfsCLnFKf6LuyNqda1cfX999eb2ec2fWbJRUppCLUtMSNNY9cXbt1d7hYhzcHtnMNeaKY6E8gfknA3"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"op5jaXUVbNnBbZvDeffBew7R7Wk64o1z3DDmZsm5eUCUGxqBVLq","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1ZNWFe3LmEJYTydctcgD6a5Apemwdtimn4","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1ZNWFe3LmEJYTydctcgD6a5Apemwdtimn4","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1ZNWFe3LmEJYTydctcgD6a5Apemwdtimn4","cycle":160,"change":"2000000"}],"delegate":"tz1ZNWFe3LmEJYTydctcgD6a5Apemwdtimn4","slots":[12]}}],"signature":"sigdSpxyGzgBzJwjCNMwc7eB7wR2yCg2T4bqFCALxNoHMbwk2BuzZpkcT9oB9uyryFp8D98z1N8Zuiuf2QsDWJYRhQ9UGt4V"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"ooAD9LiA6B2WkPL1TerMKmE9yRS88iLJ6aLu36KYeW7AiMWeFLo","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1cYufsxHXJcvANhvS55h3aY32a9BAFB494","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1cYufsxHXJcvANhvS55h3aY32a9BAFB494","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1cYufsxHXJcvANhvS55h3aY32a9BAFB494","cycle":160,"change":"2000000"}],"delegate":"tz1cYufsxHXJcvANhvS55h3aY32a9BAFB494","slots":[29]}}],"signature":"sigq6vqRXBGtNxAKGXDebkKuAiZ32AuubdFfWYxWVC7rcQZLXdHF77CYh1pWe9a6T9tYJScunPLLcBXPrSa9Ci1g2RDhhquz"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"oo4MbTFvuRbmuuoDk2PhgDSDosnWWh7TTvdMgKFtE9C4VGLXzA1","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1TzaNn7wSQSP5gYPXCnNzBCpyMiidCq1PX","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1TzaNn7wSQSP5gYPXCnNzBCpyMiidCq1PX","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1TzaNn7wSQSP5gYPXCnNzBCpyMiidCq1PX","cycle":160,"change":"2000000"}],"delegate":"tz1TzaNn7wSQSP5gYPXCnNzBCpyMiidCq1PX","slots":[28]}}],"signature":"sigT352d7euTLAd86hmYkoKtiYGiwztnBjGXagJ57VAvkMWrWdLMsr2Xf7XTB9UxiHGL8kAEqi5DsE9GAQybHapeSdk9P71f"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onnM23y3xTP44JERUxGjb7L9NysSAu4F3TUMDxMCLHC3PQW76sD","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q","cycle":160,"change":"2000000"}],"delegate":"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q","slots":[22]}}],"signature":"sigjzk2mwPj3ehguWfo91L9bRNk6GYaXpvGpmamdZVAHemyJNnTSh5akZmfHTSTsBtKhnvo6VmpvpXUTAzA47d3UppsrESuh"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"ooPNNa1EqP4XAoeE4YSXtxkXeR5kwrd3vk9mEGXhrxMBohJBXab","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1MQJPGNMijnXnVoBENFz9rUhaPt3S7rWoz","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1MQJPGNMijnXnVoBENFz9rUhaPt3S7rWoz","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1MQJPGNMijnXnVoBENFz9rUhaPt3S7rWoz","cycle":160,"change":"2000000"}],"delegate":"tz1MQJPGNMijnXnVoBENFz9rUhaPt3S7rWoz","slots":[11]}}],"signature":"sigqATJCs6kGB4nQcFDi3qKyzCVCsgUahUiRgUy1VvW3c6DYso61o5foQgrUs1uimsDKnraf8JaKAdTywb1MwNh9jZDxf1zr"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onm1Vb3EfuG2AV65iy2AM8kuGzwCztSJJyd1nUHT8gA7Xtjc2yq","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ","change":"-128000000"},{"kind":"freezer","category":"deposits","delegate":"tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ","cycle":160,"change":"128000000"},{"kind":"freezer","category":"rewards","delegate":"tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ","cycle":160,"change":"4000000"}],"delegate":"tz1PesW5khQNhy4revu2ETvMtWPtuVyH2XkZ","slots":[8,6]}}],"signature":"sigiigy2hCsCc4mNfMVcCCM6iTVx91e3gSZwLct9umXT78ePLcS3KK7daAWJV7DTt4W3DqjoaDNMGueV6Nrv2uspw4u2f9q1"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"ooBotaGjECaK2zsRvp5XJRQDvuG6VAVdoFMpp3vz7CK8sbPNnzP","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9","change":"-128000000"},{"kind":"freezer","category":"deposits","delegate":"tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9","cycle":160,"change":"128000000"},{"kind":"freezer","category":"rewards","delegate":"tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9","cycle":160,"change":"4000000"}],"delegate":"tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9","slots":[16,1]}}],"signature":"sigu7AH5eBchsQJKy5TEn5dpyqsCtMD1fcAvUKPsHR8J55UpAeAkphkp8ArZ4YsxoTWoxD73LuVwxfbwt2MzJrzVZYvEPAgy"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onn9oJuTL3kDeyJTJLJvwS6tHZCW45pWGTjKS1Z3V8F6FP9LDMG","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz2TSvNTh2epDMhZHrw73nV9piBX7kLZ9K9m","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz2TSvNTh2epDMhZHrw73nV9piBX7kLZ9K9m","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz2TSvNTh2epDMhZHrw73nV9piBX7kLZ9K9m","cycle":160,"change":"2000000"}],"delegate":"tz2TSvNTh2epDMhZHrw73nV9piBX7kLZ9K9m","slots":[17]}}],"signature":"sigXBRHgYUWmiWV3rjUev9ya1KKUwcfcSAZXFsvJGrNt1uyJzLVrBrLN3yArbYEaMvs6GQeG1bQ2SCSfqJxE6XjDMBPHMwc5"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"oot71v2BWHkY4u6YbBSnxKr8yrSC91Timx4dVULCjQ1ozBNoq2q","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV","change":"-192000000"},{"kind":"freezer","category":"deposits","delegate":"tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV","cycle":160,"change":"192000000"},{"kind":"freezer","category":"rewards","delegate":"tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV","cycle":160,"change":"6000000"}],"delegate":"tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV","slots":[27,21,3]}}],"signature":"sigYdXJzgfzi3dvoako4yZRMynMhfHVFuZqgTEGRFUcYH8DNcPexNCWQu68YGaZvTmFUkdZgvULFUzPKJwBr3TBVn8PcjELo"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"ooGN3KRyHG2GN9k2ErGkVXwPEkTbsHCvz1kVAoMTKfiVGnKD3df","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1LLNkQK4UQV6QcFShiXJ2vT2ELw449MzAA","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1LLNkQK4UQV6QcFShiXJ2vT2ELw449MzAA","cycle":160,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1LLNkQK4UQV6QcFShiXJ2vT2ELw449MzAA","cycle":160,"change":"2000000"}],"delegate":"tz1LLNkQK4UQV6QcFShiXJ2vT2ELw449MzAA","slots":[13]}}],"signature":"sigbzN2P5g68V5e5zJYgK9GiNVc3YS3BVk7G7STYKJubXPMbxoSKcQVeDLZVqCUFrhX4bxh6K2uB9N2Cv4T5AoUXUvumuNsr"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onoYsU5KsEG52Kg3D2JDpRBhQkcrr1jj3B9TMb33bf122Dvh7Kp","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs","change":"-128000000"},{"kind":"freezer","category":"deposits","delegate":"tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs","cycle":160,"change":"128000000"},{"kind":"freezer","category":"rewards","delegate":"tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs","cycle":160,"change":"4000000"}],"delegate":"tz1TRqbYbUf2GyrjErf3hBzgBJPzW8y36qEs","slots":[24,23]}}],"signature":"sighrXoa64SrCaF8NEFwY6Y43MymnnhCng3z4tVKEtiBstPeEb5wvpGSSrSLFSqTyWqhJegwrtPfUdUBwmco3fYHFSt1T1pp"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"ooREyyvfhsERYDXJ7j5sAJDSrzejCk43WpfrbS4WFvPYqQ8cuEc","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"endorsement","level":656938,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1eEnQhbwf6trb8Q8mPb2RaPkNk2rN7BKi8","change":"-192000000"},{"kind":"freezer","category":"deposits","delegate":"tz1eEnQhbwf6trb8Q8mPb2RaPkNk2rN7BKi8","cycle":160,"change":"192000000"},{"kind":"freezer","category":"rewards","delegate":"tz1eEnQhbwf6trb8Q8mPb2RaPkNk2rN7BKi8","cycle":160,"change":"6000000"}],"delegate":"tz1eEnQhbwf6trb8Q8mPb2RaPkNk2rN7BKi8","slots":[26,19,2]}}],"signature":"sigPS31652DcdLKvjKq4q2s3YYZjgztBjQDbrL2qpMpm21CP5QRUzbc23vdQPV1QtSNVGPhVM8BrgknrnxKAR5YabQfdqmdA"}],[],[],[{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"ooGypsBLe5Rk3zWVWj67JBYwwCxFTAWhM1YuAvN94K9oGP1Xeyz","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"transaction","source":"tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc","fee":"12700","counter":"4984","gas_limit":"10307","storage_limit":"0","amount":"601000000","destination":"tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH","metadata":{"balance_updates":[{"kind":"contract","contract":"tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc","change":"-12700"},{"kind":"freezer","category":"fees","delegate":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","cycle":160,"change":"12700"}],"operation_result":{"status":"applied","balance_updates":[{"kind":"contract","contract":"tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc","change":"-601000000"},{"kind":"contract","contract":"tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH","change":"601000000"}],"consumed_gas":"10207"}}}],"signature":"sigdRJMJR2PntGv93EStrDmywwTH1y6eyEGnuBqV8CD3mxhar3tijXwTJFRQ23Dow8wDuiVrGJeAB2MCso6sMRU8ASBmGr78"},{"protocol":"PsBabyM1eUXZseaJdmXFApDSBqj8YBfwELoxZHHW77EMcAbbwAS","chain_id":"NetXdQprcVkpaWU","hash":"onyxb5CSqYoosmYtQzhAy2nw5164PTJqGEjMZ7PH4n8yDQmKLrn","branch":"BLtsAJFK9gSbhiUcAeA4aKfDZXADYCQTHoqgcP1nL3JBCypNAE1","contents":[{"kind":"transaction","source":"tz1MnyYPixeTxv339nNotoHzgXhBTG7ERHqE","fee":"5000","counter":"2156617","gas_limit":"18000","storage_limit":"257","amount":"1209000","destination":"KT1MFnMZGpC13cEZP51X2YaBi13Htf4MeRai","metadata":{"balance_updates":[{"kind":"contract","contract":"tz1MnyYPixeTxv339nNotoHzgXhBTG7ERHqE","change":"-5000"},{"kind":"freezer","category":"fees","delegate":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","cycle":160,"change":"5000"}],"operation_result":{"status":"applied","storage":{"bytes":"006c8fa9f196d1f2c7f2956c04d022b20461b15520"},"balance_updates":[{"kind":"contract","contract":"tz1MnyYPixeTxv339nNotoHzgXhBTG7ERHqE","change":"-1209000"},{"kind":"contract","contract":"KT1MFnMZGpC13cEZP51X2YaBi13Htf4MeRai","change":"1209000"}],"consumed_gas":"15285","storage_size":"232"}}}],"signature":"sighbifn5ebqWnMbs7yLC2getEskrDuqqyVa23vEGAYP8cArd6njyt2EVYiJznrgGH3qLr6bwvDgWyb8kFBXne2GFXqbQLm1"}]]}`)
	mockBlocksResp          = []byte(`[["BLUdLeoqJtswBAmboRjokR8bM8aiD22FzfM2LVVp5NR8sxLt15r"]]`)
	mockBootstrapResp       = []byte(`{"block":"BKoarwjfdpFP9W3pAeYLxXDJ3pgc3yTeCV75PceUpbrR1BreDqt","timestamp":"2019-12-12T11:26:11Z"}`)
	mockChainIDResp         = []byte(`"NetXdQprcVkpaWU"`)
	mockCheckpointResp      = []byte(`{"block":{"level":38913,"proto":2,"predecessor":"BMKAi2DP2PrLR6vkmdbaR4LG5UQntuD8iGqZ8FWDDC2SS9DKmi4","timestamp":"2019-12-14T13:28:43Z","validation_pass":4,"operations_hash":"LLoaDd6G4Gre7N6QQDFdJmmnn2Mo4dzKfSrehQoBLkdyPpUijWZui","fitness":["01","0000000000009800"],"context":"CoW4nLbc2FDs2HBRAwoHEFvQurEDondEtG2ToDx2Mp6QZ1rh27my","protocol_data":"0000b1a7b92bfdf70000004e875f142dfd0edbbb8e42bb5d398d3b44e710a4f84a8f8aad121418222785b7f9ba3a45ef24553ef332d052126ad5fdb1e1acbe343aafa22db9e9419f86e109"},"save_point":38913,"caboose":0,"history_mode":"full"}`)
	mockCommitResp          = []byte(`"47e6a0f0134335480f728e245b77461190ca5ac4"`)
	mockConnectionsResp     = []byte(`[{"incoming":false,"peer_id":"idsjDeSrQivFbkSmT24cmFn7zaFcHL","id_point":{"addr":"::ffff:51.158.99.28","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}},{"incoming":false,"peer_id":"idtW4LULbPAaoT5tBzvMJWx94HBkrh","id_point":{"addr":"::ffff:157.230.147.195","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}},{"incoming":false,"peer_id":"idrZgdkNYrdHn5My4uHxUGxFRLpcNt","id_point":{"addr":"::ffff:46.245.179.161","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}},{"incoming":false,"peer_id":"idtFhWy1vk3FgshGoKomYq12M7n7FZ","id_point":{"addr":"::ffff:88.198.23.236","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}},{"incoming":false,"peer_id":"idrAf4BTfTf511fiXWYk1yKHQ2WGzW","id_point":{"addr":"::ffff:34.90.171.21","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}},{"incoming":false,"peer_id":"idtXDzsjXynvPDq6iTSK3VdpcqLHKT","id_point":{"addr":"::ffff:54.214.190.58","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}},{"incoming":false,"peer_id":"idtYrBANfMh8W1WR3KSxzTFXGBh6bS","id_point":{"addr":"::ffff:173.212.230.241","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}},{"incoming":false,"peer_id":"idrMC3YXYWk498kCGan884Hk6FvkKM","id_point":{"addr":"::ffff:34.243.126.77","port":9732},"remote_socket_port":9732,"announced_version":{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0},"private":false,"local_metadata":{"disable_mempool":false,"private_node":false},"remote_metadata":{"disable_mempool":false,"private_node":false}}]`)
	mockConstantsResp       = []byte(`{"proof_of_work_nonce_size":8,"nonce_length":32,"max_revelations_per_block":32,"max_operation_data_length":16384,"max_proposals_per_delegate":20,"preserved_cycles":5,"blocks_per_cycle":4096,"blocks_per_commitment":32,"blocks_per_roll_snapshot":256,"blocks_per_voting_period":32768,"time_between_blocks":["60","40"],"endorsers_per_block":32,"hard_gas_limit_per_operation":"800000","hard_gas_limit_per_block":"8000000","proof_of_work_threshold":"70368744177663","tokens_per_roll":"8000000000","michelson_maximum_type_size":1000,"seed_nonce_revelation_tip":"125000","origination_size":257,"block_security_deposit":"512000000","endorsement_security_deposit":"64000000","block_reward":"16000000","endorsement_reward":"2000000","cost_per_byte":"1000","hard_storage_limit_per_operation":"60000","test_chain_duration":"1966080","quorum_min":2000,"quorum_max":7000,"min_proposal_quorum":500,"initial_endorsers":24,"delay_per_missing_endorsement":"8"}`)
	mockCounterResp         = []byte(`"10"`)
	mockCycleResp           = []byte(`{"last_roll":[],"nonces":[],"random_seed":"04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846","roll_snapshot":4}`)
	mockDelegateResp        = []byte(`{"balance":"172956949254","frozen_balance":"114852428606","frozen_balance_by_cycle":[{"cycle":173,"deposit":"22272000000","fees":"147164","rewards":"673666666"},{"cycle":174,"deposit":"18432000000","fees":"79380","rewards":"570200000"},{"cycle":175,"deposit":"16448000000","fees":"41769","rewards":"509466666"},{"cycle":176,"deposit":"20544000000","fees":"114842","rewards":"629400000"},{"cycle":177,"deposit":"18432000000","fees":"35395","rewards":"570400000"},{"cycle":178,"deposit":"15296000000","fees":"76724","rewards":"474800000"}],"staking_balance":"1216660108948","delegated_contracts":["tz1gxmCTN8BSwuPLghDydtDKTqnAKyD8QTv7","KT1JoAP7MfiigepR332u6xJqza9CG52ycYZ9","KT1JsHBFpoGRVXpcfC763YwvonKtNvaFotpG","KT1Lm4ZSyXSHod7U6znR7z9SGVmexntNQwAp","KT1RbwPHzDwU9oPjnTWZrbCrMGjaFyj8dEtC","KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj","KT1MSFeAGaWk8w7F1gmgUMaarU7mH385ueYC","KT1EbMbqTUS8XnqGVRsdLZVKLhcT7Zc33jR1","tz1VESLfEAEwDEKhyLZJYXVoervFk5ABPUUD","KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy","KT1LfoE9EbpczdzUzowRckGUfikGcd5PyVKg","tz1RomjUZ1j9F2vqE24h2Am8UeGUpcrf6vvJ","tz1PB27kbPL64MWYoNZAfQAEmzCZFi9EvgBw","tz1RoDhaKjJjqcVy9MCN85bVCvbXHEnAFC7j","tz1hMkcTRoxKcWhtShoLTAGxTGsUhsa2g2zJ","KT1NGd6RaRtmvwexYXGibtdvKBnNjjpBNknn","KT1KJ5Qt18yU9DrqN36tgyLtaSvFSZ5r6YL6","tz1LRir5SfRcC4LNfagetqzKRMRjGNBTiHNH","KT19Q8GiYqGpuuUjf9xfXXVu1WY889N8oxRe","KT1QB9UAT1okYfcPQLi4jBmZkYg7LHcepERV","KT1PY2MMiTUkZQv7CPekXy186N1qmu7GikcT","KT1NmVtU3CNqzhNWwLhE5BqAopjkcmHpWzT2","KT19ABG9KxbEz2GrdN6uhGfxLmMY7REikBN8","KT1NxnFWHW7bUxzks1oHVU2jn4heu48KC3eD","KT1GBWviYFdRiNkhwM7LfrKDgHWnpdxpURtx","KT1WQWXvRcMjJB1y6mYZytoS5QsFJyFNDCk5","KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv","KT1JJcydTkinquNqh6kE5HYgFpD2124qHbZp","KT1BjtEUxd25wwdwGH432LoP6PskvUc2bEYV","tz1dfUssfLfTBoYqsWxMu86ycmLUvfF2abng","KT18ni9Yar4UzwZozFbRF7SFUKg2EqyyUPPT","KT1Wp4tXL6GUtABkikB68fT7SaPQY2UuFkuE","KT1VUbpty8fER7npuvsfYDZXf2wVPhAHVqSx","KT1NfMCxyzwev243rKk3Y6SN8GfmdLKwASFQ","KT1CeUNtCrXFNbLmvdGPNnxpcJw2sW5Hcpmc","KT1AThmRzcn51NwMf25NFYTqawjVo62hWiCv","KT1REp3D8dkiVVi37TCSMJNgGeX6UigBtfaL","KT18uqwoNyPRHpHCrg7xBFd7CiAZMbS1Ffne","KT1W3oiS6s9NgSxhZY1nCsazW2QbwkmjkET1","KT1T3dPMBm7D3kKqALKYnW2mViFqMMVCYtmo","KT1RuTPgQ6kdpnE3Adnw7Hr2KFN45uC3BdBy","KT1Aeg9D8kvkbAb6yikUdFcroReXvHtMBaZz","KT1Na4maJ99GE6CGA1vEocWXrKRmxmsVUaTi","KT1HccFB3cn4BR2za9XMuU7Wht64omed2UW8","KT1XiGwpmguFEnZDtBDDGisGxXw6qKJHPjdB","KT1XrBAocuiE3C2vvtgt7PFoazrC1KRi9ZF4","KT1S9VbEnU8nj33ufxrGBYGxBCnqmeoAnKt4","KT1PDBuQmFLVHfiWZjV248QdTrdcmAuSS7Tx","tz1hbXhPVUX1fC8hN7fALyaUpdoC6EMgqM2h","KT1A1sZmBQS9oZnPePRwP3Jyzv41xEppxfbF","KT193c72q6eP1VpaY7hiheE7k1eDZiXeQUUw","KT1KeNNxEM4NyfrmF1CG6TLn3nRSmEGhP7Z2","KT1Dgma8bbDtAbtMbYYS5VmziyCANAZn8M7W","KT1Jw925NVi4FzTVohZk5iLqagnhJGDEQoTS","tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH","KT1CySPLDUSYyJ9vqNCF2dGgit4Rw2yUNEcj","KT18kTf8UujihcF46Zn3rsFdEYFL1ZNFnGY4","KT1Lnh39om2iqr4qb9AarF9T38ayNBLnfAVn","KT1CQiyDJ3mMVDoEqLY8Fz1onFXo5ycp5BDN","KT1VzTs5piA7kYQkkfA9QNApVqGq1h6eMuV4","KT1QLo7DzPZnYK2EhmWpejVUnFjQUuWFKHnc","KT1Cz1jPLuaPR99XamKQDr9PKZY1PTXzTAHH","KT1C28u6DWsBfXk3UMyGrd8zTUVMpsyvjxmp","KT1TDrRrdz6SLYLBw8ZDxLWwJpx7FVpC52bt","KT1XPMJx2wuCbbzKZx5jJyKqLpPJMHv58wni","KT1J2uk1fYSnZjxkJcUhFDkaRDhjCTRBspqv","tz1Lv6nFvAWMvNRbQF7UcX4jobGLrAhKQLNN","KT1TS49jiXxrnwhoJzAvCzGZCXLJs3XV1k6C","KT1MX2TwjSBzPaSsBUeW2k9DKehpiuMGfFcL","KT1A5seo53aLSSyHgJKZFYnh7jTZBtFnNnjz","tz1aX2DF3ioDjqDcTVmrxVuqkxhZh1pLtfHU","KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3","tz1MXhttCeSJYpF3QRmPkMLCfNZaVufEuJmJ","tz1SnvfwMUYfD2uJrHBiaj4XPstW3eUE9RJU","KT1Re5utTU2hrujXgZ3Ux5BgjN8rbru4sns2","KT1J4WFQRV3942phzRrh87WDFWKrNVcDJTP9","KT1JPeGNVarLsPZnSb3hG5xMVmJJmmBnrnpT","KT1AT7N9bGhViSorUrpivuYT6Wxs37hR2p9d","KT1BXmBgMSViAViNyhvkb441e2RBFMiKdnj7","KT1E1MnvNgCDLqnGneStVY9CvmjnyZgcaPaD","KT1EidADxWfYeBgK8L1ZTbf7a9zyjKwCFjfH","KT19Aro5JcjKH7J7RA6sCRihPiBQzQED3oQC","KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC","tz1Z48RMPT1vjqNyUASCexnCEvEEE93J1pwL","KT1MfT8XvQp9ZeGUx4cmCNF3wui55WLNYhq9","KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt","KT1GcSsQaTtMB2HvUKU9b6WRFUnGpGx9JwGk","KT1WBDsJhoRvsvsRCmJirz9AFhSySSvzTWVd","KT1Uh1G9tdq45N63ZBrreDKy7eZF8QVoydm1","KT1UVUasDXH6mg8NCzRRgqvcjMoDUpETYEzH"],"delegated_balance":"1047131093026","deactivated":false,"grace_period":184}`)
	mockDelegatesResp       = []byte(`["tz1Lv6nFvAWMvNRbQF7UcX4jobGLrAhKQLNN","tz1Lv6nFvAWMvNRbQF7UcX4jobGLrAhKQLNN"]`)
	mockDelegationsResp     = []byte(`["tz1fjKL11UAE6DywFh8cN4hmPnQoGzx5WED4","tz1NiLdRoJSubtSinLsE87YpQb2hapjS3Nt8","tz1U4HfT4FPq5BVra6JqGjEzSKnDmECpNB4G","tz1YGLnq1Ls4W3rPanAvCvmcuQ1H5rffnc2V","tz1VG1fQCc59cMnL1g66qVFc5VUACvitvhjk","tz1gwhmUXZ5WU6XRWF5VbKLe1YbT3hcXut4t","tz1SuN5H7coasZuSCHwQxfjoYD897XexGd5Z","tz1Wkq6KDXFr8R1Pfc9qsFKNq1cGZWxsDW8M","tz1YVZtEHtGprtVZtSC342pn5AeskibJ8U5Q","tz1KgZycAC7vvqsFoQAyYufCU6woQDQ7kji3","tz1XN5zXpiKrFae9KE8pLkcSeDhvLxkBSM9d","tz1NSpdiL4LZnsgCwRHhwhY8hfXbTWv1P9m1","tz1PR69DJGUbVkk4bUWJP6mNsQo8EomwfAry","tz1e7axDeFykoUUQUnb64V7xjKrXaCBBPH97","tz1YX4g6wzjKPRz7wm3punsAb5XbTuuq9xV6","tz1LjkdWbWG6YTVsJrdUNFDS6MDksLFzpX84","tz1cwxtKJNuxsb4qBR9gUcBUJUnZc35G6FXe","tz1SMGkowfL5BtWAyr6WrTbMKpcR1SR5Rx3A","tz1N9HduDLV88e881wVZ9QVmvCRTePhTChnP","tz1VxYYPquaFGcDaRUAL4HEbmKVsJCLnb5hZ","tz1NXk3MtEs7dAkHpVPvnfNp4BsJYeXF6VxX","tz1XMPEFYLbwTinPZwVYPy2ctgB8qqrqg1fu","tz1TwuTmmYZkCx8Ws1LRfEGkKUU9ZK2ovdQ7","tz1LeaLL9V2GtRSQy9MeecudKtmKKfrhx2nA","tz1MorwRuLTAFz3c2BDVfAnCCJ4LcaPn3Yhh","tz1bsAUiZR6CBsxDYVqr5vvzRvWAanunwXeC"]`)
	mockEndorsingRightsResp = []byte(`[{"level":822092,"delegate":"tz3bvNMQ95vfAYtG8193ymshqjSvmxiCUuR5","slots":[25],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","slots":[20,17],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN","slots":[15,11],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz3UoffC7FG7zfpmvmjUmUeAaHvzdcUvAj6r","slots":[7],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz3RDC3Jdn4j15J7bBHZd29EUee9gVB1CxD9","slots":[30,26,13,12,8],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz3RB4aoyjov4KEVRbuhvQ1CKJgBJMWhaeB8","slots":[10,5],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz3NExpXn9aPNZPorRE4SdjJ2RGrfbJgMAaV","slots":[19],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz2BzJTyoQp8fNbfhWD4YQgH9JJHDgSGzpdG","slots":[1],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1isXamBXpTUgbByQ6gXgZQg4GWNW7r6rKE","slots":[21],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1irJKkXS2DBWkU1NnmFQx1c1L7pbGg4yhk","slots":[24,14,0],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1gk3TDbU7cJuiBRMhwQXVvgDnjsxuWhcEA","slots":[28],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1gfArv665EUkSg2ojMBzcbfwuPxAvqPvjo","slots":[23],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1eEnQhbwf6trb8Q8mPb2RaPkNk2rN7BKi8","slots":[22],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1Z2jXfEXL7dXhs6bsLmyLFLfmAkXBzA9WE","slots":[4],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1W5VkdB5s7ENMESVBtwyt9kyvLqPcUczRT","slots":[16,9],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1UuQ4HWDu3ALNRgAq94dX9MhqhQhnuY3gC","slots":[6],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1TzaNn7wSQSP5gYPXCnNzBCpyMiidCq1PX","slots":[2],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1Tnjaxk6tbAeC2TmMApPh8UsrEVQvhHvx5","slots":[3],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1S8MNvuFEUsWgjHvi3AxibRBf388NhT1q2","slots":[18],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1LmaFsWRkjr7QMCx5PtV6xTUz3AmEpKQiF","slots":[31],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1Ldzz6k1BHdhuKvAtMRX7h5kJSMHESMHLC","slots":[27],"estimated_time":"2020-02-13T02:46:51Z"},{"level":822092,"delegate":"tz1LBEKXaxQbd5Gtzbc1ATCwc3pppu81aWGc","slots":[29],"estimated_time":"2020-02-13T02:46:51Z"}]`)
	mockFrozenBalanceResp   = []byte(`{"deposits":"15296000000","fees":"76724","rewards":"474800000"}`)
	mockInvalidBlockResp    = []byte(`{"block":"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1","level":10,"errors":[{"kind":"err","error":"message"}]}`)
	mockInvalidBlocksResp   = []byte(`[{"block":"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1","level":10,"errors":[{"kind":"err","error":"message"}]}]`)
	mockOperationHashesResp = []byte(`["BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1","BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1","BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1"]`)
	mockRPCErrorResp        = []byte(`[{"kind":"somekind","Error":"someerror"}]`)
	mockStakingBalanceResp  = []byte(`"1216660108948"`)
	mockVersionResp         = []byte(`{"chain_name":"TEZOS_MAINNET","distributed_db_version":0,"p2p_version":0}`)
)

// The below variables contain mocks that are unmarshaled.
var (
	mockAddressTz1 = "tz1YGLnq1Ls4W3rPanAvCvmcuQ1H5rffnc2V"
	mockBlockHash  = "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1"
)

// Regexes to allow the capture of custom handlers for unit testing.
var (
	regActiveChains       = regexp.MustCompile(`\/monitor\/active_chains`)
	regBakingRights       = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/baking_rights`)
	regBalance            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/balance`)
	regBlock              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+`)
	regBlocks             = regexp.MustCompile(`\/chains\/main\/blocks`)
	regBoostrap           = regexp.MustCompile(`\/monitor\/bootstrapped`)
	regChainID            = regexp.MustCompile(`\/chains\/main\/chain_id`)
	regCheckpoint         = regexp.MustCompile(`\/chains\/main\/checkpoint`)
	regCommit             = regexp.MustCompile(`\/monitor\/commit_hash`)
	regConnections        = regexp.MustCompile(`\/network\/connections`)
	regConstants          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/constants`)
	regCounter            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/counter`)
	regCycle              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/json\/cycle\/[0-9]+`)
	regDelegate           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+`)
	regDelegates          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates`)
	regDelegatedContracts = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/delegated_contracts`)
	regEndorsingRights    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/endorsing_rights`)
	regFrozenBalance      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/json\/contracts\/index\/[A-z0-9]+\/frozen_balance\/[0-9]+`)
	regInjectionBlock     = regexp.MustCompile(`\/injection\/block`)
	regInjectionOperation = regexp.MustCompile(`\/injection\/operation`)
	regInvalidBlocks      = regexp.MustCompile(`\/chains\/main\/invalid_blocks`)
	regOperationHashes    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/operation_hashes`)
	regPreapplyOperations = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/preapply\/operations`)
	regStakingBalance     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/staking_balance`)
	regStorage            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/storage`)
	regVersions           = regexp.MustCompile(`\/network\/version`)
)

// blankHandler handles the end of a http test handler chain
var blankHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

// ----------------------------------------- //
// Mock Handlers
// The below handlers are to simulate the Tezos RPC server for unit testing.
// ----------------------------------------- //

func activeChainsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regActiveChains.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bakingRightsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBakingRights.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func balanceHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBalance.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type blockHandlerMock struct {
	used bool
}

func newBlockMock() *blockHandlerMock {
	return &blockHandlerMock{}
}

func (b *blockHandlerMock) handler(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBlock.MatchString(r.URL.String()) && !b.used {
			w.Write(resp)
			b.used = true
			return
		}

		next.ServeHTTP(w, r)
	})
}

func blocksHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBlocks.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bootstrapHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBoostrap.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func chainIDHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regChainID.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkpointHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCheckpoint.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func commitHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCommit.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func connectionsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regConnections.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type constantsHandlerMock struct {
	used bool
}

func newConstantsMock() *constantsHandlerMock {
	return &constantsHandlerMock{}
}

func (c *constantsHandlerMock) handler(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regConstants.MatchString(r.URL.String()) && !c.used {
			w.Write(resp)
			c.used = true
			return
		}

		next.ServeHTTP(w, r)
	})
}

func counterHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCounter.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cycleHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCycle.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func delegateHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegate.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func delegatesHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegates.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func delegationsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegatedContracts.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func endorsingRightsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regEndorsingRights.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func frozenBalanceHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regFrozenBalance.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func injectionBlockHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regInjectionBlock.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func injectionOperationHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regInjectionOperation.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func invalidBlocksHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regInvalidBlocks.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func operationHashesHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regOperationHashes.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func preapplyOperationsHandlerMock(preapplyResp, blockResp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regPreapplyOperations.MatchString(r.URL.String()) {
			w.Write(preapplyResp)
			return
		}

		if regBlock.MatchString(r.URL.String()) {
			w.Write(blockResp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func stakingBalanceHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regStakingBalance.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func storageHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regStorage.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func versionsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regVersions.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errContains)
	} else {
		assert.Nil(t, err)
	}
}

func testGoTezos(t *testing.T, handler http.Handler) *GoTezos {
	server := httptest.NewServer(handler)
	defer server.Close()

	gt, err := New(server.URL)
	assert.Nil(t, err)
	assert.NotNil(t, gt)

	return gt
}
