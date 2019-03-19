from behave import *
import requests
import json

@given('payment record')
def step_impl(context):
    payment = {
        "type": "Instant",
        "attributes": {
            "amount": "",
            "currency": "",
            "beneficiary_party": {
                "account_name": "",
                "account_number": "",
            },
            "debtor_party": {
                "account_name": "",
                "account_number": "",
            },
            "payment_id": "",
            "payment_type": "",
            "processing_date": "",
            "reference": "",
        }
    }
    context.payment = payment
    pass


@when('I create new payment via HTTP')
def step_impl(context):
    payment = context.payment
    payment_json = json.dumps(payment)
    resp = requests.post("http://localhost:8080/", payment_json)
    context.response = resp
    pass


@then('I get id of created payment')
@when('I get id of created payment')
def step_impl(context):
    response = context.response
    if response.status_code != 201:
        # normally we should have more business oriented language
        fail("not created; status code not 201")

    location = response.headers["Location"]
    context.id = location
    if location == "":
        # normally we should have more business oriented language
        fail("empty location id")
    pass

@then('I can find payment')
def step_impl(context):
    resp = requests.get("http://localhost:8081/")
    if resp.status_code != 200:
        # normally we should have more business oriented language
        fail("status code not 200")
pass

@when('I create new payment')
def step_impl(context):
    context.execute_steps("""
        when I create new payment via HTTP
        then I get id of created payment
    """)
    pass


@when('I delete payment')
@then('I delete payment')
def step_impl(context):
    print(f'id {context.id}')
    resp = requests.delete(f'http://localhost:8080/{context.id}')
    context.response = resp
    pass


@then('I cannot find payment')
def step_impl(context):
    response = context.response
    if response.status_code != 200:
        # normally we should have more business oriented language
        fail("status code not 200")

    resp = requests.get(f'http://localhost:8081/{context.id}')
    if resp.status_code != 404:
        # normally we should have more business oriented language
        fail("status code not 404")
    pass
