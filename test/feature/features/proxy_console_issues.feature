Feature: Proxy console issues
    Background:
        #
        # Run routers with coordinators
        # Stop all coordinators
        #
        Given cluster environment is
        """
        ROUTER_CONFIG=/spqr/test/feature/conf/router_with_coordinator.yaml
        ROUTER_COORDINATOR_CONFIG=/spqr/test/feature/conf/router_coordinator.yaml
        ROUTER_2_COORDINATOR_CONFIG=/spqr/test/feature/conf/router_coordinator_2.yaml
        """
        Given cluster is up and running
        And host "coordinator" is stopped
        And host "coordinator2" is stopped

        #
        # Make host "router" take control over coordinator
        #
        Given host "router2" is stopped
        And I wait for coordinator "regress_router" to take control
        And host "router2" is started
        And I wait for host "router2" to respond

        When I run SQL on host "router-admin"
        """
        UNREGISTER ROUTER ALL;
        REGISTER ROUTER r1 ADDRESS regress_router::7000;
        REGISTER ROUTER r2 ADDRESS regress_router_2::7000;
        """
        Then command return code should be "0"

    Scenario: Add key_range is executed in coordinator
        When I run SQL on host "router-admin"
        """
        CREATE DISTRIBUTION ds1 COLUMN TYPES integer;
        CREATE KEY RANGE krid1 FROM 0 ROUTE TO sh1 FOR DISTRIBUTION ds1;
        """
        Then command return code should be "0"

        #
        # Check on first router
        #
        When I run SQL on host "router-admin"
        """
        SHOW key_ranges
        """
        Then SQL result should match json_exactly
        """
        [{
            "Key range ID":"krid1",
            "Distribution ID":"ds1",
            "Lower bound":"0",
            "Shard ID":"sh1"
        }]
        """

        #
        # Check on second router
        #
        When I run SQL on host "router2-admin"
        """
        SHOW key_ranges
        """
        Then SQL result should match json_exactly
        """
        [{
            "Key range ID":"krid1",
            "Distribution ID":"ds1",
            "Lower bound":"0",
            "Shard ID":"sh1"
        }]
        """
