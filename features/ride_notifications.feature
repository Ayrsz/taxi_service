Feature: Ride Notifications for Driver
  As a taxi driver
  I want to receive relevant notifications about my rides
  So that I can effectively manage my duties

  Scenario: Driver receives destination arrival notification
    Given I am logged in as driver "João" with ID "1"
    And I have an active ride with destination at (40.7128, -74.0060) and origin at (40.7000, -74.0100)
    When my current location updates to (40.7127, -74.0061)
    Then I should receive a notification "Você chegou ao destino"

  Scenario: System displays cancellation confirmation to driver
    Given I am logged in as driver "João" with ID "1"
    And I have an accepted ride with ID "101" with estimated distance 8 km and estimated value R$ 20.00
    And I am viewing the estimated arrival time of 10 minutes
    When I select the "Cancel ride" option for ride "101"
    Then the system should display a notification with message "Tem certeza que deseja cancelar a corrida? Cancelamentos frequentes podem impactar sua avaliação."
    And the notification should include options "Sim, quero cancelar" and "Não, continuar com a corrida"

  Scenario: Successfully completed ride is recorded in history
    Given I am logged in as driver "João" with ID "1"
    And I have completed a ride with ID "102", actual distance 15.2 km and actual value R$ 32.50
    When I access my ride history
    Then I should see the completed ride with ID "102" with current date and time
    And the total value "32.50" and distance traveled "15.2" km are displayed for ride "102"

  Scenario: System displays cancelled ride with details in driver's history
    Given I am logged in as driver "João" with ID "1"
    And I have cancelled a ride with ID "103" that had estimated destination 12 km and estimated value R$ 28.00
    When I access my ride history
    Then the system should display an entry in the history for ride "103" with status "cancelled"
    And this entry for ride "103" should show the cancellation date and time, estimated distance "12" km, and estimated value R$ "28.00" of the ride

  Scenario: Estimated Time of Arrival notification is displayed upon accepting a ride
    Given I am logged in as driver "João" with ID "1"
    And an available ride with ID "104" from (40.7000, -74.0100) to (40.7200, -74.0200) exists
    When I accept the ride with ID "104"
    Then I should see the estimated time of arrival to the pickup location displayed