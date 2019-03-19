Feature: payments

  Scenario: create new payment
    Given payment record
    When I create new payment via HTTP
    And I get id of created payment
    Then I can find payment
    # clean after self (should be improved)
    And I delete payment

  Scenario: delete payment
    Given payment record
    When I create new payment
    And I delete payment
    Then I cannot find payment