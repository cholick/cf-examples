import json
import os
import traceback
import html

import flask
import requests
import requests.auth

app = flask.Flask(__name__)


def debug_request():
    if app.debug_request:
        print("Request for:", flask.request.url)
        content = flask.request.get_json()
        if content is not None:
            print("Body:", json.dumps(content, indent=4), "\n")

        querystring = flask.request.query_string
        if len(querystring) > 0:
            print("Query string: ", querystring, "\n")


def path_to_url(url: str) -> str:
    return "http://{}{}".format(app.host, url)


@app.route("/env")
def env():
    return flask.jsonify(dict(os.environ))


@app.route("/dashboard_login")
def dashboard_login():
    permissions = "+".join([
        "openid", "cloud_controller_service_permissions.read",
    ])
    redirect_uri = "https://uaa.{}/oauth/authorize?response_type=code&client_id={}&scope={}".format(
        app.system_domain, app.name + "-client-id", permissions
    )
    return flask.redirect(redirect_uri)


@app.route("/dashboard")
def dashboard():
    debug_request()

    error = flask.request.args.get("error", None)
    if error:
        return "failure: " + error

    code = flask.request.args.get("code", None)
    if not code:
        return "failure: no code found"

    token_url = "https://uaa.{}/oauth/token".format(
        app.system_domain,
    )

    client_id = app.name + "-client-id"
    uaa_response = requests.post(token_url, data={
        "client_id": client_id,
        "client_secret": app.shared_secret,
        "grant_type": "authorization_code",
        "code": code,
        "redirect_uri": path_to_url("/dashboard")
    }, verify=False)

    if uaa_response.status_code != 200:
        return "failure:" + uaa_response.text

    try:
        uaa_response_json = uaa_response.json()
    except:
        return "failure: <pre>{}</pre>".format(html.escape(uaa_response.text))

    print("---------------")
    print(json.dumps(uaa_response_json, indent=4))
    print("---------------")

    # using this token, check back against apis's
    # v2/service_instances/<instance guid>/permissions
    # to get read/manage permission json. For example:
    #    {"manage": true, "read": true}
    # see https://docs.cloudfoundry.org/services/dashboard-sso.html#checking-user-permissions

    return "success"


@app.route("/v2/catalog")
def broker_catalog():
    return flask.jsonify({
        "services": [{
            "id": app.name + "-id",
            "name": app.name,
            "description": "Test/debugging service",
            "bindable": True,
            "dashboard_client": {
                "id": app.name + "-client-id",
                "secret": app.shared_secret,
                "redirect_uri": path_to_url("/dashboard"),
            },
            "plans": [{
                "id": app.name + "-plan1-id",
                "name": "plan1",
                "description": "Default plan"
            }]
        }]
    })


@app.route("/v2/service_instances/<instance_id>", methods=["PUT"])
def broker_provision_instance(instance_id):
    debug_request()

    return flask.jsonify({
        "dashboard_url": path_to_url("/dashboard_login")
    }), 200


@app.route("/v2/service_instances/<instance_id>", methods=["DELETE"])
def broker_deprovision_instance(instance_id):
    debug_request()

    return flask.jsonify({}), 200


@app.route("/v2/service_instances/<instance_id>/service_bindings/<binding_id>", methods=["PUT"])
def broker_bind_instance(instance_id, binding_id):
    debug_request()

    return flask.jsonify({
        "credentials": {
            "user": "root", "password": "monkey123"
        }
    }), 201


@app.route("/v2/service_instances/<instance_id>/service_bindings/<binding_id>", methods=["DELETE"])
def broker_unbind_instance(instance_id, binding_id):
    debug_request()

    return flask.jsonify({}), 200


@app.errorhandler(500)
def internal_error(error):
    print(error)
    return "Internal server error", 500


if __name__ == "__main__":
    try:
        vcap_application = json.loads(
            os.getenv("VCAP_APPLICATION", json.dumps({
                "application_uris": ["http://localhost:8080"],
                "application_name": "broker-sso"
            })))
        app.host = vcap_application["application_uris"][0]
        app.name = vcap_application["application_name"]
        app.debug_request = True

        app.shared_secret = os.getenv("SSO_SHARED_SECRET", "note-secure")
        if "apps" in app.host:
            app.system_domain = "sys." + app.host.split(".apps.")[1]
        else:
            app.system_domain = "localhost:8080"
            print("SSO will not work, looks like a local environment")

        app.run(host="0.0.0.0", port=int(os.getenv("PORT", "8080")))
    except:
        print("Exited with exception")
        traceback.print_exc()
