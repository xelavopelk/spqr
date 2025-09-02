Feature: Coordinator issues test
  Background:
    #
    # Make host "coordinator" take control
    #
    Given cluster environment is
    """
    ROUTER_CONFIG=/spqr/test/feature/conf/router_cluster.yaml
    """
    Given cluster is up and running
    And host "coordinator2" is stopped
    And host "coordinator2" is started

    When I run SQL on host "coordinator"
    """
    REGISTER ROUTER r1 ADDRESS regress_router::7000
    """
    Then command return code should be "0"

    When I run SQL on host "coordinator"
    """
    CREATE DISTRIBUTION ds1 COLUMN TYPES integer; 
    CREATE KEY RANGE krid2 FROM 100 ROUTE TO sh2 FOR DISTRIBUTION ds1;
    CREATE KEY RANGE krid1 FROM 50 ROUTE TO sh1 FOR DISTRIBUTION ds1;
    ALTER DISTRIBUTION ds1 ATTACH RELATION test DISTRIBUTION KEY id;
    """
    Then command return code should be "0"

    When I run SQL on host "router"
    """
    CREATE TABLE test(id int, name text)
    """
    Then command return code should be "0"


  Scenario: Router synchronization after registration works
    When I run SQL on host "coordinator"
    """
    UNREGISTER ROUTER r1;
    REGISTER ROUTER r1 ADDRESS regress_router::7000
    """
    Then command return code should be "0"
    When I run SQL on host "router-admin"
    """
    SHOW key_ranges
    """
    Then SQL result should match json_exactly
    """
    [{
      "Key range ID":"krid1",
      "Distribution ID":"ds1",
      "Lower bound":"50",
      "Shard ID":"sh1"
    },
    {
      "Key range ID":"krid2",
      "Distribution ID":"ds1",
      "Lower bound":"100",
      "Shard ID":"sh2"
    }]
    """