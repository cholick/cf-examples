import json
import os
import sys
import traceback

import flask

import auth

app = flask.Flask(__name__)

content_header = {'Content-Type': 'application/json; charset=utf-8'}
drain_url = None


@app.route("/health")
def health():
    return "healthy"


@app.route("/v2/catalog")
@auth.requires_auth
def broker_catalog():
    # Catalog ids were randomly generated guids, per best practices
    # Also note the requires, different from other brokers in that respect
    catalog = {
        "services": [{
            "id": '1fc88b6c-5be0-4fb8-bf7f-d6cf9c527b35',
            "name": 'my-log-drain',
            "description": "Provide a pre-configured drain to applications",
            "bindable": True,
            "requires": ["syslog_drain"],
            "plans": [{
                "id": '3b6b760c-b9e6-4a47-9b48-ae95febde9e1',
                "name": "plan1",
                "description": "Plan 1: the only plan"
            }]
        }]
    }
    return json.dumps(catalog, indent=4)


@app.route("/v2/service_instances/<instance_id>", methods=['PUT'])
@auth.requires_auth
def broker_provision_instance(instance_id):
    # no-op. Real implementation: allocate index? Reserve resources? Something along those lines
    return "{}", 201, content_header


@app.route("/v2/service_instances/<instance_id>/service_bindings/<binding_id>", methods=['PUT'])
@auth.requires_auth
def broker_bind_instance(instance_id, binding_id):
    response_body = json.dumps({"syslog_drain_url": drain_url})
    return response_body, 201, content_header


@app.route("/v2/service_instances/<instance_id>/service_bindings/<binding_id>", methods=['DELETE'])
@auth.requires_auth
def broker_unbind_instance(instance_id, binding_id):
    # no-op
    return "{}", 200, content_header


@app.route("/v2/service_instances/<instance_id>", methods=['DELETE'])
@auth.requires_auth
def broker_deprovision_instance(instance_id):
    # no-op
    return "{}", 200, content_header


@app.errorhandler(500)
def internal_error(error):
    print(error)
    return "Internal server error", 500


if __name__ == "__main__":
    try:
        drain_url = os.getenv("DRAIN_URL")
        if not drain_url:
            print("The environment variable DRAIN_URL is required")
            sys.exit(1)
        if not os.getenv("SECURITY_USER_NAME") or not os.getenv("SECURITY_USER_PASSWORD"):
            print("The environment variables SECURITY_USER_NAME and SECURITY_USER_PASSWORD")
            sys.exit(1)

        app.run(host='0.0.0.0', port=int(os.getenv('PORT', '8080')))
        print("Exited normally")
    except:
        print("* Exited with exception")
        traceback.print_exc()
