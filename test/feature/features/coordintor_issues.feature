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


  Scenario: Coordinator can restart
    #
    # Coordinator is Up
    #
    Given host "coordinator2" is stopped
    When I run SQL on host "coordinator"
    """
    CREATE KEY RANGE krid3 FROM 31 ROUTE TO sh1 FOR DISTRIBUTION ds1
    """
    Then command return code should be "0"

    #
    # Coordinator has been restarted
    #
    Given host "coordinator" is stopped
    And host "coordinator" is started
    When I run SQL on host "coordinator"
    """
    SHOW key_ranges
    """
    Then command return code should be "0"
    And SQL result should match json
    """
    [{
      "Key range ID":"krid3",
      "Distribution ID":"ds1",
      "Lower bound":"31",
      "Shard ID":"sh1"
    }]
    """