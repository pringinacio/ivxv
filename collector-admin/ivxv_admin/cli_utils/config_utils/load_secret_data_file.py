# IVXV Internet voting framework
"""CLI utilites for sercet loading."""

import hashlib

import OpenSSL

from ... import (COLLECTOR_STATE_CONFIGURED, COLLECTOR_STATE_FAILURE,
                 COLLECTOR_STATE_INSTALLED, COLLECTOR_STATE_PARTIAL_FAILURE,
                 SERVICE_SECRET_TYPES, SERVICE_STATE_CONFIGURED,
                 SERVICE_STATE_FAILURE, SERVICE_STATE_INSTALLED,
                 SERVICE_TYPE_PARAMS, lib)
from ...event_log import register_service_event
from ...lib import IvxvError
from ...service.service import Service
from .. import init_cli_util, log


def main():
    """Load secret data to IVXV services."""
    # validate CLI arguments
    args = init_cli_util("""
    Load secret data to IVXV services.

    This utility loads file that contains secret data to services.

    Supported secret types are:

        tls-cert - TLS certificate for service.

            Certificate (and key) is used for securing
            communication between services and service instances.

        tls-key - TLS key for service.

            Key is used together with service certificate.

        tsp-regkey - PKIX TSP registration key for voting services.

            Key is used for signing Time Stamp Protocol requests.

            Key file must be in PEM format and must be not password protected.

        mid-token-key - Mobile-ID/Smart-ID/Web eID identity token for
                        choices, mobile-id and voting services.

            Key file must be 32 bytes long.

    Usage: ivxv-secret-load [--service=<service-id>] <secret-type> <keyfile>
    """)
    secret_type = args['<secret-type>'].lower()
    filepath = args['<keyfile>']
    if secret_type not in SERVICE_SECRET_TYPES:
        log.error("Invalid secret type %r", secret_type)
        log.info('Supported secret types are: %s',
                 ', '.join(SERVICE_SECRET_TYPES))
        return 1
    secret_descr = SERVICE_SECRET_TYPES[secret_type]['description']

    # load file
    log.debug('Loading %s file', secret_descr)
    try:
        with open(filepath, 'rb') as fp:
            file_content = fp.read()
    except (FileNotFoundError, PermissionError) as err:
        log.error("Unable to load file %r: %s", filepath, err.strerror)
        return 1

    # validate file
    try:
        validate_secret_file(secret_type, file_content, args["--service"])
    except IvxvError as err:
        log.error('Error while validating %s: %s', secret_descr, err)
        return 1

    # calculate file checksum
    file_checksum = hashlib.sha256(file_content).hexdigest()

    # generate list of services that are in required state
    key_param = {
        'tls-cert': 'require_tls',
        'tls-key': 'require_tls',
        'tsp-regkey': 'tspreg',
        'mid-token-key': 'mobile_id',
    }[secret_type]
    service_types_affected = sorted(
        set(
            service_type
            for service_type, service_params in SERVICE_TYPE_PARAMS.items()
            if service_params[key_param]
        )
    )

    services = lib.get_services(
        include_types=service_types_affected,
        require_collector_state=[
            COLLECTOR_STATE_INSTALLED, COLLECTOR_STATE_CONFIGURED,
            COLLECTOR_STATE_FAILURE, COLLECTOR_STATE_PARTIAL_FAILURE
        ],
        service_state=[
            SERVICE_STATE_INSTALLED, SERVICE_STATE_CONFIGURED,
            SERVICE_STATE_FAILURE
        ])
    if not services:
        return 1

    # copy key to service hosts
    # FIXME avoid multiple copying of shared secret to single host
    services_updated = []
    for service_id, service_data in sorted(services.items()):
        if args['--service'] and service_id != args['--service']:
            continue
        if (service_data.get(SERVICE_SECRET_TYPES[secret_type]['db-key']) ==
                file_checksum):
            log.info('Service %s already contains specified %s',
                     service_id, secret_descr)
            continue

        service = Service(service_id, service_data)
        if not service.load_secret_file(secret_type, filepath, file_checksum):
            return 1
        register_service_event(
            'SECRET_INSTALL',
            service=service_id,
            params={
                'secret_descr': secret_descr,
            })
        services_updated.append(service_id)

    if args['--service'] and args['--service'] not in services_updated:
        log.error('%s was not loaded to service %s',
                  secret_descr, args['--service'])
        return 1

    if args['--service']:
        log.info('%s is loaded to service %s', secret_descr, args['--service'])
    else:
        log.info('%s is loaded to services', secret_descr)

    return 0


def validate_secret_file(secret_type, file_content, service_id):
    """Validate secret data file.

    :raises IvxvError:
    """
    if secret_type in ('tls-cert', 'tls-key') and not service_id:
        raise IvxvError("Service ID is not specified")
    if secret_type == 'mid-token-key' and len(file_content) != 32:
        raise IvxvError(
            f'File is not 32 bytes long '
            f'(actual size: {len(file_content)} bytes)')
    if secret_type == 'tsp-regkey':
        try:
            privkey = OpenSSL.crypto.load_privatekey(
                OpenSSL.crypto.FILETYPE_PEM, file_content
            )
        except OpenSSL.crypto.Error as err:
            err_lib, err_func, err_reason = err.args[0][0]
            raise IvxvError(
                f'Error in {err_lib} library {err_func} '
                f'function: {err_reason}')
        with open("/var/lib/ivxv/service/tspreg-pubkey.pem", "wb") as fd:
            fd.write(
                OpenSSL.crypto.dump_publickey(OpenSSL.crypto.FILETYPE_PEM, privkey)
            )
