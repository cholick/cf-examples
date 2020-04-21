import json
import os
import traceback

import flask

app = flask.Flask(__name__)


def debug_request():
    if app.debug_request:
        content = flask.request.get_json()
        print("\n", json.dumps(content, indent=4), "\n")


@app.route("/env")
def env():
    return flask.jsonify(dict(os.environ))


@app.route("/v2/catalog")
def broker_catalog():
    return flask.jsonify({
        "services": [{
            "id": app.name + '-id',
            "name": app.name + '-service',
            "description": "Test/debugging service",
            "bindable": True,
            "plans": [{
                "id": app.name + '-plan-1-id',
                "name": "plan-1",
                "description": "Default plan"
            }]
        }]
    })


@app.route("/v2/service_instances/<instance_id>", methods=['PUT'])
def broker_provision_instance(instance_id):
    debug_request()

    return json.dumps({"operation": "working_on_it"}, indent=4), 202


@app.route("/v2/service_instances/<instance_id>/last_operation", methods=['GET'])
def broker_last_operation(instance_id):
    return flask.jsonify({
        "state": "in progress",
        "description": "user facing (state is enumeration)"
    }), 200


@app.route("/v2/service_instances/<instance_id>", methods=['DELETE'])
def broker_deprovision_instance(instance_id):
    debug_request()

    return flask.jsonify({}), 200


@app.route("/v2/service_instances/<instance_id>/service_bindings/<binding_id>", methods=['PUT'])
def broker_bind_instance(instance_id, binding_id):
    debug_request()

    return flask.jsonify({
        'credentials': {
            'user': 'root', 'password': 'monkey123'
        }
    }), 201


@app.route("/v2/service_instances/<instance_id>/service_bindings/<binding_id>", methods=['DELETE'])
def broker_unbind_instance(instance_id, binding_id):
    debug_request()

    return json.dumps({}, indent=4), 200


@app.errorhandler(500)
def internal_error(error):
    print(error)
    return "Internal server error", 500


if __name__ == "__main__":
    try:
        vcap_application = json.loads(
            os.getenv('VCAP_APPLICATION', '{ "name": "none", "application_uris": [ "http://localhost:8080" ] }')
        )
        app.host = vcap_application['application_uris'][0]
        app.name = "cholick-async"
        app.debug_request = True

        app.run(host='0.0.0.0', port=int(os.getenv('PORT', '8080')))
    except:
        print("Exited with exception")
        traceback.print_exc()
