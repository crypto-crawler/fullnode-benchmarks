const apps = [
  {
    "name": "bloxroute_pair_reserve_cloud",
    "script": "bloxroute_pair_reserve",
    "args": "-cert external_gateway_cert.pem -key external_gateway_key.pem",
    "exec_interpreter": "none",
    "exec_mode": "fork_mode",
    "instances": 1,
    "restart_delay": 1000
  },
  {
    "name": "bloxroute_pair_reserve_gateway",
    "script": "bloxroute_pair_reserve",
    "args": "-header YOUR_HEADER -gateway ws://localhost:28334 -output bloxroute-pair-reserve-gateway.json",
    "exec_interpreter": "none",
    "exec_mode": "fork_mode",
    "instances": 1,
    "restart_delay": 1000
  },
  {
    "name": "fullnode_pair_reserve",
    "script": "fullnode_pair_reserve",
    "exec_interpreter": "none",
    "exec_mode": "fork_mode",
    "instances": 1,
    "restart_delay": 1000
  },
  {
    "name": "fullnode_pair_reserve_bulk",
    "script": "fullnode_pair_reserve_bulk",
    "exec_interpreter": "none",
    "exec_mode": "fork_mode",
    "instances": 1,
    "restart_delay": 1000
  },
  {
    "name": "fullnode_pair_reserve_bulk_head",
    "script": "fullnode_pair_reserve_bulk_head",
    "exec_interpreter": "none",
    "exec_mode": "fork_mode",
    "instances": 1,
    "restart_delay": 1000
  }
]

module.exports = {
  apps,
};
