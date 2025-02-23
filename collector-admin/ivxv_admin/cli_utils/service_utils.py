# IVXV Internet voting framework
"""CLI utilities for service management."""

import datetime
import os
import shutil
import json
import logging
import random
import subprocess
import sys
import tempfile
import time
import zipfile
from collections import OrderedDict
from jinja2 import Environment, PackageLoader

import fasteners
import yaml
import requests

from .. import (
    COLLECTOR_STATE_CONFIGURED,
    COLLECTOR_STATE_FAILURE,
    COLLECTOR_STATE_INSTALLED,
    COLLECTOR_STATE_NOT_INSTALLED,
    COLLECTOR_STATE_PARTIAL_FAILURE,
    SERVICE_STATE_CONFIGURED,
    SERVICE_STATE_FAILURE,
    SERVICE_STATE_INSTALLED,
    SERVICE_STATE_REMOVED,
    SERVICE_TYPE_PARAMS,
    __version__,
    command_file,
)
from ..agent_daemon import get_collector_data, ping_service
from ..config import cfg_path
from ..db import IVXVManagerDb
from ..lib import IvxvError, get_services
from ..service import logging as service_logging
from ..service.service import Service, exec_remote_cmd
from . import init_cli_util, log

#: JSON formatting options
JSON_DUMP_ARGS = dict(indent=2, sort_keys=True)


def export_votes_util():
    """Export collected votes."""
    args = init_cli_util("""
    Export collected votes.

    This utility copies current ballot box from voting service to backup
    service and outputs ballot box content.

    Usage: ivxv-export-votes [--consolidate] <output-file>

    Options:
        --consolidate   Consolidate all collected votes
    """)

    try:
        export_votes(args["--consolidate"], args["<output-file>"])
    except IvxvError as err:
        log.error(err)
        return 1

    return 0


def export_votes(must_consolidate, output_filename):
    """Export votes."""
    # fail if output file already exist
    if os.path.exists(output_filename):
        raise IvxvError(f"Output file {output_filename!r} already exist")

    # create backup of current ballot box
    log.info('Creating backup copy from current ballot box')
    try:
        subprocess.run(['ivxv-backup', 'ballot-box'], check=True)
    except OSError as err:
        raise IvxvError(f"Creating ballot box backup failed: {err.strerror}")
    except subprocess.CalledProcessError as err:
        raise IvxvError(f"Creating ballot box backup failed: {err}")

    # create handler for backup service
    services = get_services(include_types=['backup'])
    if not services:
        raise IvxvError('Backup service is not defined')
    assert len(services) == 1

    backup_service = Service(*list(services.items())[0])
    log.debug('Backup service: %s', backup_service.service_id)
    ballot_box_filepath = datetime.datetime.now().strftime(
        '/var/lib/ivxv/ballot-box-consolidated-%Y%m%d_%H%M.zip')

    # run consolidation in backup service
    if must_consolidate:
        proc = backup_service.ssh([
            'ivxv-voteunion', ballot_box_filepath,
            '/var/backups/ivxv/ballot-box/ballot-box-????????_????.zip'
        ])
        if proc.returncode:
            raise IvxvError('Consolidation command failed in backup service')
    else:  # export without consolidation
        # detect last backup filename
        cmd = [
            'ls', '/var/backups/ivxv/ballot-box/ballot-box-????????_????.zip'
        ]
        proc = backup_service.ssh(cmd, stdout=subprocess.PIPE, check=True)
        ballot_box_filepath = proc.stdout.decode('UTF-8').strip().split(
            '\n')[-1]

    # copy consolidated ballot box from backup
    log.info('Copying ballot box to management service')
    if not backup_service.scp(
            output_filename,
            ballot_box_filepath,
            "consolidated ballot box" if must_consolidate else "ballot box",
            to_remote=False):
        raise IvxvError('Failed to copy ballot box to management service')
    if must_consolidate:
        log.info('Removing consolidated ballot box from backup service')
        backup_service.ssh(['rm', '-v', ballot_box_filepath])

    log.info("Collected votes archive is written to %r", output_filename)


def generate_processor_input_util():
    """Generate input for processor application."""
    args = init_cli_util(
        """
    Generate input for processor application.

    This utility generates ZIP container with data files
    for processor application to validate ballot box:

        1. District list;
        2. Voter lists;
        3. Validation key for vote registration requests;
        4. Configuration for processor application.

    Usage: ivxv-generate-processor-input <output-file>
    """
    )

    try:
        election_id, cfg = generate_processor_check_cfg()
        generate_processor_input(election_id, cfg, args["<output-file>"])
    except IvxvError as err:
        log.error(err)
        return 1

    return 0


def generate_processor_check_cfg():
    """Generate config file for processor application."""
    log.info("Generating processor application config")

    processor_check_cfg = OrderedDict(
        districts=None,
        registrationlist="registrationlist.zip",
        registrationlist_checksum="registrationlist.zip.sha256sum.asice",
        vlkey=None,
        tskey="ts.key",
        voterlists=[],
        election_start=None,
        voterforeignehak=None,
        out=None,
    )

    # read collector config from management database
    with IVXVManagerDb() as db:
        election_id = db.get_value("election/election-id")
        processor_check_cfg["election_start"] = db.get_value("election/electionstart")
        districts_version = db.get_value("list/districts")
        voter_list_states = []
        for changeset_no in range(10_000):
            try:
                voter_list_state = db.get_value(f"list/voters{changeset_no:04d}-state")
            except KeyError:
                break
            if voter_list_state not in ["APPLIED", "SKIPPED"]:
                raise IvxvError(
                    f"Voter list changeset #{changeset_no} "
                    f"has unexpected state {voter_list_state!r}"
                )
            voter_list_states.append(voter_list_state)
    if not election_id:
        raise IvxvError("Election is not configured")
    if not districts_version:
        raise IvxvError("District list is not loaded")
    if not voter_list_states:
        raise IvxvError("Initial voter list is not loaded")

    processor_check_cfg.update(
        {
            "districts": f"{election_id}.districts.json.bdoc",
            "out": f"{election_id}-out-1",
        }
    )

    # read election config from command file
    try:
        election_cfg = command_file.load_cfg_file_content(
            "election", f"{election_id}.election.yaml", "/etc/ivxv/election.bdoc"
        )
    except (IvxvError, OSError) as err:
        raise IvxvError(f"Error while loading config file: {err}")
    except UnicodeDecodeError as err:
        raise IvxvError(f"Error while decoding config file: {err}")

    # copy election config values
    if "voterforeignehak" in election_cfg:
        processor_check_cfg["voterforeignehak"] = election_cfg["voterforeignehak"]
    else:
        processor_check_cfg.pop("voterforeignehak")
    processor_check_cfg["vlkey"] = election_cfg["voterlist"]["key"]

    # generate voter list filenames
    for changeset_no, voter_list_state in enumerate(voter_list_states):
        voter_list_record = dict(
            path=f"{changeset_no:02d}.{election_id}.voters.utf",
            signature=f"{changeset_no:02d}.{election_id}.voters.sig",
        )
        if voter_list_state == "SKIPPED":
            voter_list_record[
                "skip_cmd"
            ] = f"{changeset_no:02d}.{election_id}.voters-skip.yaml.bdoc"
        processor_check_cfg["voterlists"].append(voter_list_record)

    return election_id, processor_check_cfg


def generate_processor_input(election_id, processor_check_cfg, output_filename):
    """Generate input file for processor application."""
    # fail if output file already exist
    if os.path.exists(output_filename):
        raise IvxvError(f"Output file {output_filename!r} already exist")

    log.info("Creating input file for processor application")

    with tempfile.TemporaryDirectory() as output_dir:
        log.info("Preparing container structure in directory %r", output_dir)

        # district list
        filename = processor_check_cfg["districts"]
        log.info("Copying district list %r", filename)
        shutil.copy(
            cfg_path("active_config_files_path", "districts.bdoc"),
            f"{output_dir}/{filename}",
        )

        # key to verify voter list signer
        filename = "voterfile.pub.key"
        log.info("Copying voter list signing key %r", filename)
        with open(f"{output_dir}/{filename}", "w") as fd:
            fd.write(processor_check_cfg["vlkey"])
        processor_check_cfg["vlkey"] = filename

        # voter lists
        for changeset_no, voter_list_files in enumerate(
            processor_check_cfg["voterlists"]
        ):
            input_filepath = cfg_path(
                "active_config_files_path",
                f"voters{changeset_no:04d}.zip" if changeset_no else "voters0000.bdoc",
            )
            with zipfile.ZipFile(input_filepath) as input_zip:
                for description, key in [
                    ["content", "path"],
                    ["signature", "signature"],
                    ["skipping command", "skip_cmd"],
                ]:
                    try:
                        filename = voter_list_files[key]
                    except KeyError:
                        continue
                    file_ext = os.path.splitext(filename)[-1]
                    log.info(
                        "Copying voter list #%d %s %r",
                        changeset_no,
                        description,
                        filename,
                    )
                    if key == "skip_cmd":
                        shutil.copy(
                            cfg_path(
                                "active_config_files_path",
                                f"voters{changeset_no:04d}.bdoc",
                            ),
                            f"{output_dir}/"
                            f"{changeset_no:02d}.{election_id}.voters-skip.yaml.bdoc",
                        )
                    else:
                        with input_zip.open(
                            f"{election_id}-voters-{changeset_no}{file_ext}"
                        ) as zip_fd:
                            with open(f"{output_dir}/{filename}", "wb") as fd:
                                fd.write(zip_fd.read())

        # key to verify registration requests
        filename = processor_check_cfg["tskey"]
        log.info("Copying registration requests verification key %r", filename)
        shutil.copy(
            "/var/lib/ivxv/service/tspreg-pubkey.pem",
            f"{output_dir}/{filename}",
        )

        # setup YAML module to write OrderedDict
        def represent_dictionary_order(self, dict_data):
            return self.represent_mapping("tag:yaml.org,2002:map", dict_data.items())

        def setup_yaml():
            yaml.add_representer(OrderedDict, represent_dictionary_order)

        setup_yaml()

        # config file for processor application
        filename = f"{election_id}.processor.yaml"
        log.info("Writing processor application config %r", filename)
        with open(f"{output_dir}/{filename}", "w") as fd:
            fd.write(
                "# Automaatselt genereeritud fail "
                f"{datetime.datetime.now():%d.%m.%Y %R}\n\n"
                f"# Generaator: {__file__}\n\n"
            )
            yaml.dump(dict(check=processor_check_cfg), fd, default_flow_style=False)
            fd.write('\n')
            fd.write('  # BALLOT BOX\n')
            fd.write('  #ballotbox: votes.zip\n')
            fd.write('  #ballotbox_checksum: votes.zip.sha256sum\n')

        # generate ZIP container
        filenames = sorted(os.listdir(output_dir))
        output_filename_tmp = f"{output_filename}.tmp"
        log.info(
            "Generating ZIP container %r with %d files",
            output_filename_tmp,
            len(filenames),
        )
        with zipfile.ZipFile(output_filename_tmp, "w", allowZip64=True) as zip_file:
            for filename in filenames:
                log.info("Adding %r to ZIP container", filename)
                zip_file.write(f"{output_dir}/{filename}", filename)

    shutil.move(output_filename_tmp, output_filename)

    log.info("Processor input is written to %r", output_filename)


def copy_logs_to_logmon_util():
    """Copy IVXV log files from service hosts to Log Monitor."""
    # validate CLI arguments
    args = init_cli_util("""
    Copy IVXV log files from service hosts to Log Monitor.

    This utility transports collected IVXV log files from IVXV services
    (including Log Collector Service) to Log Monitor.

    Usage: ivxv-copy-log-to-logmon [--log-level=<level>] [<hostname> ...]

    Options:
        <hostname>              Service host name.
        --log-level=<level>     Logging level [Default: INFO].
    """)

    # check collector state
    with IVXVManagerDb() as db:
        collector_state = db.get_value('collector/state')
        logmonitor_address = db.get_value('logmonitor/address')
        services = db.get_all_values('service')
    if collector_state == COLLECTOR_STATE_NOT_INSTALLED:
        if not sys.stdout.isatty():  # suppress warning if executed from crontab
            return 0
        log.warning("Collector is not installed")
        return 1

    # create service host list
    hostnames = args['<hostname>']
    if not hostnames:
        hostnames = []
        for service in services.values():
            is_main_service = (
                SERVICE_TYPE_PARAMS[service['service-type']]['main_service'])
            if ((is_main_service
                    or service['service-type'] == 'log')
                    and service['state'] != SERVICE_STATE_REMOVED):
                hostnames.append(service['ip-address'].split(':')[0])
        hostnames = list(set(hostnames))

    # checking access to Log Monitor account
    if not logmonitor_address:
        log.error('Log monitor is not defined')
        return 1
    log.info("Using address %r for log monitor", logmonitor_address)
    logmon_account = f'logmon@{logmonitor_address}'
    log.info('Checking SSH access to Log Monitor account %s', logmon_account)
    proc = exec_remote_cmd(['ssh', logmon_account, 'true'])
    if proc.returncode:
        log.error('Cannot access to Log Monitor (%s)', logmon_account)
        return 1

    # copy log file from service hosts to Log Monitor
    exit_code = 0
    for hostname in hostnames:
        remote_account = f'ivxv-admin@{hostname}'

        # check if service host have Log Monitor host key
        log.info("Checking if %r have Log Monitor host key", remote_account)
        proc = exec_remote_cmd(
            ['ssh', hostname, 'ssh-keygen', '-F', logmonitor_address],
            stdout=subprocess.PIPE)
        if proc.returncode:
            log.error(
                'Failed to check Log Monitor host key for %s', remote_account)
            log.info('Installing Log Monitor host key for %s', remote_account)
            proc = subprocess.run(
                ['sh', '-c',
                 f'ssh-keygen -F {logmonitor_address} | '
                 'grep -v ^# | '
                 f'ssh {hostname} tee --append .ssh/known_hosts'],
                check=False,
            )
            if proc.returncode:
                log.error('Failed to install Log Monitor host key for %s',
                          remote_account)
                continue

        if not copy_logs_from_host_to_logmon(hostname, logmon_account):
            exit_code = 1

    return exit_code


def copy_logs_from_host_to_logmon(hostname, logmon_account):
    """Copy log files from service host to Log Monitor."""
    # acquire process lock
    lockfile_path = cfg_path(
        'ivxv_admin_data_path',
        f'service/copy-log-from-{hostname}-to-logmon.lock')
    lock = fasteners.InterProcessLock(lockfile_path)
    if not lock.acquire(blocking=False):
        log.warning('Lock exists for process "copy logs from host %s '
                    'to Log Monitor"', hostname)
        return False

    # execute rsync command in host to copy log file to Log Monitor
    proc = None
    try:
        log.info(
            "Copying IVXV service log files from host %r to Log Monitor",
            hostname)
        transfer_cmd = [
            'ivxv-admin-helper', 'copy-logs-to-logmon', hostname,
            logmon_account
        ]
        cmd = ['ssh-agent', 'ssh', '-A', hostname] + transfer_cmd

        proc = subprocess.run(cmd, check=False)
    finally:
        lock.release()

    if not proc.returncode:
        return True

    log.error("Failed to copy log file from host %r to Log Monitor",
              hostname)
    return False


def update_software_pkg_util():
    """Update service packages in service hosts."""
    # validate CLI arguments
    args = init_cli_util("""
    Update service packages in IVXV service hosts.

    This utility checks versions of software packages in service hosts
    and installs new versions if required.

    Usage: ivxv-update-packages [--force]

    Options:
        --force     Update even package version does not require update
    """)

    # generate list of voting services that are in required state
    services = get_services(
        require_collector_state=[
            COLLECTOR_STATE_INSTALLED,
            COLLECTOR_STATE_CONFIGURED,
            COLLECTOR_STATE_FAILURE,
            COLLECTOR_STATE_PARTIAL_FAILURE,
        ],
        service_state=[
            SERVICE_STATE_INSTALLED,
            SERVICE_STATE_CONFIGURED,
            SERVICE_STATE_FAILURE,
        ])
    if not services:
        return 1

    host_versions = dict(
        [services[service_id]['ip-address'].split(':')[0], '']
        for service_id in sorted(services))
    update_res = dict(failure=[], install=[], skip=[])

    def get_installed_pkg_ver():
        """Get version string if installed package in service host."""
        service.log.info('Detect %s package version', pkg_name)
        proc = service.ssh(
            f'dpkg --status {pkg_name} | grep ^Version: | cut -d: -f2',
            stdout=subprocess.PIPE,
            account='ivxv-admin')
        return proc.stdout.decode('UTF-8').strip()

    for service_id, service_data in sorted(services.items()):
        service = Service(service_id, service_data)

        # check ivxv-common version
        host_version = host_versions.get(service.hostname)
        pkg_name = 'ivxv-common'
        install_pkg = bool(args['--force'])
        if not install_pkg and host_version != __version__:
            host_versions[service.hostname] = get_installed_pkg_ver()
            install_pkg = host_versions[service.hostname] != __version__

        # install ivxv-common if required
        if install_pkg:
            if not service.update_ivxv_common_pkg():
                update_res['failure'].append([service_id, pkg_name])
                continue
            update_res['install'].append([service_id, pkg_name])
            host_versions[service.hostname] = __version__
        else:
            update_res['skip'].append([service_id, pkg_name])

        # check service package version and upgrade if required
        pkg_name = service.deb_pkg_name
        install_pkg = args['--force'] or get_installed_pkg_ver() != __version__

        # install package if required
        if install_pkg:
            if not service.install_service_pkg(is_update=True):
                update_res['failure'].append([service_id, pkg_name])
                continue
            update_res['install'].append([service_id, pkg_name])
        else:
            update_res['skip'].append([service_id, pkg_name])

    # output result
    for service_id, pkg_name in update_res['install']:
        log.info('Successfully installed service %s package %s',
                 service_id, pkg_name)
    for service_id, pkg_name in update_res['failure']:
        log.error('Failed to install service %s package %s',
                  service_id, pkg_name)
    log.info('Service update stats: %d packages installed, '
             '%d package installations failed, %s packages skipped',
             len(update_res['install']),
             len(update_res['failure']),
             len(update_res['skip']))

    return 1 if update_res['failure'] else None


def manage_service():
    """Manage IVXV services."""
    # validate CLI arguments
    args = init_cli_util("""
    Manage IVXV services.

    Usage: ivxv-service <action> <service-id> ...

    Options:
        <action>    Management action: start, stop, restart, ping
    """)
    action = args['<action>']
    if action not in ['start', 'stop', 'restart', 'ping']:
        log.error('Invalid action: %s', action)
        return 1

    services = get_collector_data()
    if services is None:
        log.error('Election data is not loaded')
        return 1

    exit_code = 0
    for service_id in args['<service-id>']:
        if service_id not in services:
            log.error('Unknown service %s', service_id)
            exit_code = 1
            continue
        if action == 'ping':
            log.info('Pinging service %s', service_id)
            if ping_service(service_id, services[service_id]):
                log.info('Service %s is alive', service_id)
            else:
                log.error('Failed to query service %s status', service_id)
                exit_code = 1
            continue

        service = Service(service_id, services[service_id])

        if action == 'stop':
            log.info('Stopping service %s', service_id)
            if service.stop_service():
                log.info('Service %s stopped', service_id)
            else:
                log.error('Failed to stop service %s', service_id)
                exit_code = 1
            continue

        # start, restart
        log.info('%sing service %s', action.capitalize(), service_id)
        service = Service(service_id, services[service_id])
        if service.restart_service():
            log.info('Service %s %sed', service_id, action)
        else:
            log.error('Failed to %s service %s', action, service_id)
            exit_code = 1

    return exit_code


def voterstats_util():
    """Import voter stats from voting service and export to VIS."""
    args = init_cli_util(
        """
        Import voter stats from voting service and export common stats to VIS.

        Usage: ivxv-voterstats <TYPE> [--action=<action>] [--file=<file>]
                    [--service-id=<service_id>] [--log-level=<level>]

        Options:
            <TYPE>                      Stats type "common" or "detail".
            --action=<action>           Limit actions for "common" stats type.
                                        Possible values are "import" and "export".
                                        [Default: all]
            --file=<file>               Path to stats file.
            --service-id=<service_id>   Voting service ID [Default: random].
            --log-level=<level>         Logging level [Default: INFO].
        """
    )
    stats_type = args["<TYPE>"]
    action = args["--action"]
    filepath = args["--file"]
    if not filepath:
        filepath = f"/var/lib/ivxv/admin-ui-data/voterstats-{stats_type}.json"
    service_id = args["--service-id"]

    # validate CLI args
    if stats_type not in ["common", "detail"]:
        log.error("Unexpected stats type %r", stats_type)
        return 1

    # check collector config state
    try:
        with IVXVManagerDb() as db:
            if not db.get_value("list/districts"):
                raise IvxvError("District list is not loaded")
            services = get_services(
                db=db, service_state=["CONFIGURED"], include_types=["voting"]
            )
            if not services:
                raise IvxvError("No configured voting service found")
    except IvxvError as err:
        if not sys.stdout.isatty():  # suppress warning if executed from crontab
            return 0
        log.warning(err)
        return 1

    # choose voting service
    if service_id == "random":
        service_id = random.choice(list(services.keys()))
    if service_id not in services:
        services = get_services(include_types=["voting"])
        if service_id in services:
            log.error(
                "Voting service %r is not configured (current state: %s)",
                service_id,
                services[service_id]["state"],
            )
        else:
            log.error("No voting service %r found", service_id)
        return 1
    service_logging.log.setLevel(args["--log-level"])
    service = Service(service_id, services[service_id])

    # import stats
    try:
        if action in ["all", "import"]:
            import_voterstats(stats_type, service, filepath, args["--log-level"])
    except IvxvError:
        log.error("Failed to import %s stats from service %r", stats_type, service_id)
        return 1

    # export stats to VIS only if requested
    command_file.log.setLevel(args["--log-level"])
    try:
        if action in ["all", "export"]:
            export_voterstats(stats_type, filepath)
    except IvxvError:
        log.error("Failed to export %s stats to VIS", stats_type)
        return 1

    return 0


def import_voterstats(stats_type, service, filepath, log_level):
    """Import voter stats from voting service."""
    log.info("Generating %s stats in voting service %r", stats_type, service.service_id)
    cmd = ["ivxv-voterstats", "-instance", service.service_id]
    if stats_type == "detail":
        cmd.append("-detailed")
    if logging.getLevelName(log_level) >= 30:
        cmd.append("-q")
    remote_filepath = (
        f"/var/lib/ivxv/user/ivxv-voting/ivxv-voterstats-{os.getpid()}.json"
    )
    cmd.append(remote_filepath)
    proc = service.ssh(cmd)
    if proc.returncode:
        raise IvxvError(
            f"Failed to generate {stats_type} stats "
            f"in voting service {service.service_id}"
        )

    log.info(
        "Importing %s stats from voting service %r", stats_type, service.service_id
    )
    if not service.scp(
        filepath, remote_filepath, f"{stats_type} stats", to_remote=False
    ):
        raise IvxvError(f"Failed to {stats_type} stats to management service")


def export_voterstats(stats_type, filepath):
    """Export voter stats to VIS."""
    log.info("Exporting %s stats to VIS", stats_type)
    with open(filepath) as fd:
        stats = json.load(fd)

    # get params from election config
    cfg = command_file.load_collector_cmd_file("election", "/etc/ivxv/election.bdoc")
    if not cfg:
        raise IvxvError("Election config not found")
    ca_certs_filepath = None
    if cfg["vis"].get("ca"):
        ca_certs_filepath = cfg_path("vis_path", "ca.pem")
        with open(ca_certs_filepath, "w") as fd:
            fd.write("\n".join(cfg["vis"]["ca"]))

    # upload voter stats
    url = f"{cfg['vis']['url']}voters-by-county"
    log.info("Upload voter stats to %r", url)
    resp = requests.post(
        url,
        verify=ca_certs_filepath,
        cert=(
            "/etc/ssl/certs/ivxv-admin-client.crt",
            "/etc/ssl/private/ivxv-admin-client.key",
        ),
        json=stats,
    )
    log.info(
        "VIS responded to voter stats upload: %r %s", resp.status_code, resp.reason
    )

    if resp.status_code != 204:
        raise IvxvError(
            f"Server responded with status {resp.status_code} - {resp.reason}"
        )


def voting_sessions_util():
    """Utility to import voting sessions from Log Monitor."""
    args = init_cli_util(
        """
        Import list of voting sessions from Log Monitor.

        Session data is in CSV format.

        Usage: ivxv-voting-sessions (vote | verify) <output_file> [--anonymize] [--uniq]
                    [--log-level=<level>]

        Options:
            <output_file>               Write sessions to file.
            --anonymize                 Anonymize session data
                                        (IP addresses and ID codes).
            --uniq                      Consolidate session data.
            --log-level=<level>         Logging level [Default: INFO].
        """
    )

    # check collector state
    with IVXVManagerDb() as db:
        collector_state = db.get_value("collector/state")
        logmonitor_address = db.get_value("logmonitor/address")
    if collector_state == COLLECTOR_STATE_NOT_INSTALLED:
        # emit warning only if attached to console
        # (suppress if executed from crontab)
        if sys.stdout.isatty():
            log.warning("Collector is not installed")
            return 1
        return 0

    # checking access to Log Monitor account
    if not logmonitor_address:
        log.error("Log monitor is not defined")
        return 1
    log.info("Using address %r for log monitor", logmonitor_address)
    logmon_account = f"logmon@{logmonitor_address}"
    log.info("Checking SSH access to Log Monitor account %s", logmon_account)
    proc = exec_remote_cmd(["ssh", logmon_account, "true"])
    if proc.returncode:
        log.error("Cannot access to Log Monitor (%s)", logmon_account)
        return 1

    try:
        import_voting_sessions(
            logmon_account,
            session_type="vote" if args["vote"] else "verify",
            anonymize=args["--anonymize"],
            uniq=args["--uniq"],
            output_filepath=args["<output_file>"],
            log_level=args["--log-level"],
        )
    except IvxvError as err:
        log.error(err)
        return 1

    return 0


def import_voting_sessions(
    logmon_account, session_type, anonymize, uniq, output_filepath, log_level
):
    """Import voting sessions from Log Monitor."""
    # generate CSV
    log.info("Generating voting sessions file in Log Monitor")
    remote_outfile = f"~logmon/voting-sessions-{session_type}-{os.getpid()}.json"
    remote_cmd = ["ivxv-export-voting-sessions", f"--log-level={log_level}"]
    if anonymize:
        remote_cmd.append("--anonymize")
    if uniq:
        remote_cmd.append("--uniq")
    remote_cmd += [session_type, remote_outfile]
    cmd = ["ssh", logmon_account] + remote_cmd
    try:
        subprocess.run(cmd, stdout=subprocess.PIPE, check=True)
    except subprocess.CalledProcessError as err:
        raise IvxvError("Failed to generate voting sessions") from err
    log.info("Voting sessions file successfully generated in Log Monitor")

    # import CSV
    log.info("Importing voting sessions file from Log Monitor")
    cmd = ["scp", f"{logmon_account}:{remote_outfile}", output_filepath]
    try:
        subprocess.run(cmd, stdout=subprocess.PIPE, check=True)
    except subprocess.CalledProcessError as err:
        raise IvxvError("Failed to import voting sessions") from err
    log.info("Voting sessions file successfully imported from Log Monitor")


def remove_ivxv_admin_crontab():
    """Remove ivxv-admin user crontab if exists."""
    log.info('Removing ivxv-admin user crontab (if exists)')
    proc = subprocess.run(["crontab", "-r"], check=False)
    assert proc.returncode in [0, 1], 'Unexpected exit code for crontab command'


def generate_detail_stats_crontab(cfg: dict):
    """Generate crontab for detailed statistics exchange with VIS.
    Generated file will be stored to db.

    :param cfg: content of configuration file
    :type cfg: dict
    """

    # load Jinja2 crontab template for detailed statistics
    template_dir = Environment(loader=PackageLoader('ivxv_admin', 'templates'))
    template = template_dir.get_template('ivxv_detail_stats_crontab.jinja')

    # read `stats:` section from a configuration file,
    # use default values if configuration file doesn't have `stats:` section
    cron_cfg = cfg.get('stats', {}).get('detail', {}).get('scheduler', {})
    crontab_params = {
        'minute': cron_cfg.get('cron', {}).get('min', '*/15') or '*/15',
        'hour': cron_cfg.get('cron', {}).get('hour', '*') or '*',
        'day': cron_cfg.get('cron', {}).get('day', '*') or '*',
        'month': cron_cfg.get('cron', {}).get('month', '*') or '*',
        'weekday': cron_cfg.get('cron', {}).get('weekday', '*') or '*',
    }

    # render Jinja2 template
    rendered_template = template.render(
        time_generated=datetime.datetime.now().strftime('%d.%M.%Y %H:%M:%S'),
        **crontab_params,
    )

    # override crontab rendered template in a db
    with IVXVManagerDb(for_update=True) as db:
        db.set_value('stats/detail/scheduler/cron', rendered_template)


def detail_stats_crontab_editor() -> int:
    args = init_cli_util("""
    Generate crontab for IVXV detail statistics export to VIS automation.

    This utility must be called as editor by crontab utility:

        $ env VISUAL=ivxv-detail-stats-crontab crontab -e

    Usage: ivxv-detail-stats-crontab <filename>
    """)
    filepath = args['<filename>']

    # when you run `crontab -e`, cron generates temporary file at '/tmp/XYZ/crontab',
    # and that temporary file is passed here as args['<filename>'], i.e 'filepath'
    crontab_tmp_file_content: str

    # read crontab temporary file
    with open(filepath, 'r') as fp:
        try:
            crontab_tmp_file_content = fp.read()
        except Exception as err:
            msg = "Can't read crontab temporary file %r: %s"
            log.error(msg, filepath, err.__str__())
            return 1

    # get rendered Jinja2 crontab template for detail stats from a db
    with IVXVManagerDb() as db:
        # get from db
        detail_stats_crontab = db.get_value('stats/detail/scheduler/cron')
        # remove old block:
        # ### block ivxv_detail_stats_crontab ###
        # ...
        # ### endblock ivxv_detail_stats_crontab ###
        # from a crontab temporary file if any exists
        try:
            without_detail_stats_crontab = remove_detail_stats_crontab(
                data=crontab_tmp_file_content)
        except ValueError as err:
            msg = "Can't remove ivxv_detail_stats block from a temporary file %r: %s"
            log.error(msg, filepath, err.__str__())
            return 1
        # add new block to a crontab temporary file
        crontab = insert_detail_stats_crontab(
            data=without_detail_stats_crontab, crontab=detail_stats_crontab)

    # Pause for 1 second. Crontab checks mtime to detect file modifications. It
    # seems that crontab can't detect mtime change if changes happens too
    # quickly (tested in Ubuntu Xenial).
    time.sleep(1)

    # override crontab temporary file with a new content
    with open(filepath, "w") as fp:
        try:
            fp.write(crontab)
        except Exception as err:
            msg = "Can't write to a crontab temporary file %r: %s"
            log.error(msg, filepath, err.__str__())
            return 1


def remove_detail_stats_crontab(data: str) -> str:
    """Remove ivxv_detail_stats_crontab block from a data.

    :param data: any data
    :type data: str
    :return: data without ivxv_detail_stats_crontab block
    :rtype: str
    """
    block, found = __get_detail_stats_crontab_block(data=data)
    if not found:
        return data
    return data.replace(block, '').strip()


def insert_detail_stats_crontab(data: str, crontab: str) -> str:
    """Wrap crontab in an ivxv_detail_stats_crontab block and insert it to a data.

    :param data: any data
    :type data: str
    :param crontab: rendered ivxv_detail_stats Jinja2 crontab template file
    :type crontab: str
    :return: data with an ivxv_detail_stats_crontab block
    :rtype: str
    """
    header = '### block ivxv_detail_stats_crontab ###'
    tail = '### endblock ivxv_detail_stats_crontab ###'
    return data + "\n\n" + header + "\n" + crontab + "\n" + tail + "\n"


def __get_detail_stats_crontab_block(data: str) -> (str, bool):
    """Get ivxv_detail_stats_crontab block from a data.

    :param data: any data
    :type data: str
    :return: ivxv_detail_stats_crontab block and True on success,
    otherwise data and False
    :rtype: str, bool
    :raise ValueError: if ivxv_detail_stats_crontab block's header/tail is malformed,
    however though, if both header and tail are malformed then function assumes that
    ivxv_detail_stats_crontab block doesn't exist in a data
    """
    header = '### block ivxv_detail_stats_crontab ###'
    tail = '### endblock ivxv_detail_stats_crontab ###'

    # remove leading and trailing whitespaces/newlines
    data_stripped = data.strip()

    # get start index of a header
    header_start_index = data_stripped.find(header)

    # get start index of a tail
    tail_start_index = data_stripped.find(tail)

    # both header and tail aren't present in a data_stripped
    if header_start_index < 0 and tail_start_index < 0:
        return data, False
    # both header and tail present in a data_stripped
    elif header_start_index >= 0 and tail_start_index >= 0:
        # get end index of a tail
        tail_end_index = tail_start_index + len(tail)
        # extract ivxv_detail_stats_crontab block from a data_stripped
        return data_stripped[header_start_index:tail_end_index], True
    else:
        if header_start_index < 0:
            raise ValueError("malformed ### block ivxv_detail_stats_crontab ###")
        else:
            raise ValueError("malformed ### endblock ivxv_detail_stats_crontab ###")


def install_detail_stats_crontab():
    """Install crontab for detail stats export with VIS automation."""
    subprocess.run(
        ["env", "VISUAL=ivxv-detail-stats-crontab", "crontab", "-e"], check=True)


def generate_voting_facts_crontab(cfg: dict):
    """Generate crontab for ivxv-storageorder automation.
    Generated file will be stored to db.

    :param cfg: content of configuration file
    :type cfg: dict
    """

    # load Jinja2 crontab template for voting facts
    template_dir = Environment(loader=PackageLoader('ivxv_admin', 'templates'))
    template = template_dir.get_template('ivxv_voting_facts_crontab.jinja')

    # read `stats:` section from a configuration file,
    # use default values if configuration file doesn't have `stats:` section
    cron_cfg = cfg.get('stats', {}).get('voting_facts', {}).get('scheduler', {})
    crontab_params = {
        'minute': cron_cfg.get('cron', {}).get('min', '*/15') or '*/15',
        'hour': cron_cfg.get('cron', {}).get('hour', '*') or '*',
        'day': cron_cfg.get('cron', {}).get('day', '*') or '*',
        'month': cron_cfg.get('cron', {}).get('month', '*') or '*',
        'weekday': cron_cfg.get('cron', {}).get('weekday', '*') or '*',
    }

    # render Jinja2 template
    rendered_template = template.render(
        time_generated=datetime.datetime.now().strftime('%d.%M.%Y %H:%M:%S'),
        **crontab_params,
    )

    # override crontab rendered template in a db
    with IVXVManagerDb(for_update=True) as db:
        db.set_value('stats/voting_facts/scheduler/cron', rendered_template)


def voting_facts_crontab_editor() -> int:
    args = init_cli_util("""
    Generate crontab for ivxv-storageorder automation.

    This utility must be called as editor by crontab utility:

        $ env VISUAL=ivxv-voting-facts-crontab crontab -e

    Usage: ivxv-voting-facts-crontab <filename>
    """)
    filepath = args['<filename>']

    # when you run `crontab -e`, cron generates temporary file at '/tmp/XYZ/crontab',
    # and that temporary file is passed here as args['<filename>'], i.e 'filepath'
    crontab_tmp_file_content: str

    # read crontab temporary file
    with open(filepath, 'r') as fp:
        try:
            crontab_tmp_file_content = fp.read()
        except Exception as err:
            msg = "Can't read crontab temporary file %r: %s"
            log.error(msg, filepath, err.__str__())
            return 1

    # get rendered Jinja2 crontab template for voting facts from a db
    with IVXVManagerDb() as db:
        # get from db
        voting_facts_crontab = db.get_value('stats/voting_facts/scheduler/cron')
        # remove old block:
        # ### block ivxv_voting_facts_crontab ###
        # ...
        # ### endblock ivxv_voting_facts_crontab ###
        # from a crontab temporary file if any exists
        try:
            without_voting_facts_crontab = remove_voting_facts_crontab(
                data=crontab_tmp_file_content)
        except ValueError as err:
            msg = "Can't remove ivxv_voting_facts_crontab block from a tmp file %r: %s"
            log.error(msg, filepath, err.__str__())
            return 1
        # add new block to a crontab temporary file
        crontab = insert_voting_facts_crontab(
            data=without_voting_facts_crontab, crontab=voting_facts_crontab)

    # Pause for 1 second. Crontab checks mtime to detect file modifications. It
    # seems that crontab can't detect mtime change if changes happens too
    # quickly (tested in Ubuntu Xenial).
    time.sleep(1)

    # override crontab temporary file with a new content
    with open(filepath, "w") as fp:
        try:
            fp.write(crontab)
        except Exception as err:
            msg = "Can't write to a crontab temporary file %r: %s"
            log.error(msg, filepath, err.__str__())
            return 1


def remove_voting_facts_crontab(data: str) -> str:
    """Remove ivxv_voting_facts_crontab block from a data.

    :param data: any data
    :type data: str
    :return: data without ivxv_voting_facts_crontab block
    :rtype: str
    """
    block, found = __get_voting_facts_crontab_block(data=data)
    if not found:
        return data
    return data.replace(block, '').strip()


def insert_voting_facts_crontab(data: str, crontab: str) -> str:
    """Wrap crontab in an ivxv_voting_facts_crontab block and insert it to a data.

    :param data: any data
    :type data: str
    :param crontab: rendered ivxv_voting_facts_crontab Jinja2 crontab template file
    :type crontab: str
    :return: data with an ivxv_voting_facts_crontab block
    :rtype: str
    """
    header = '### block ivxv_voting_facts_crontab ###'
    tail = '### endblock ivxv_voting_facts_crontab ###'
    return data + "\n\n" + header + "\n" + crontab + "\n" + tail + "\n"


def __get_voting_facts_crontab_block(data: str) -> (str, bool):
    """Get ivxv_voting_facts_crontab block from a data.

    :param data: any data
    :type data: str
    :return: ivxv_voting_facts_crontab block and True on success,
    otherwise data and False
    :rtype: str, bool
    :raise ValueError: if ivxv_voting_facts_crontab block's header/tail is malformed,
    however though, if both header and tail are malformed then function assumes that
    ivxv_voting_facts_crontab block doesn't exist in a data
    """
    header = '### block ivxv_voting_facts_crontab ###'
    tail = '### endblock ivxv_voting_facts_crontab ###'

    # remove leading and trailing whitespaces/newlines
    data_stripped = data.strip()

    # get start index of a header
    header_start_index = data_stripped.find(header)

    # get start index of a tail
    tail_start_index = data_stripped.find(tail)

    # both header and tail aren't present in a data_stripped
    if header_start_index < 0 and tail_start_index < 0:
        return data, False
    # both header and tail present in a data_stripped
    elif header_start_index >= 0 and tail_start_index >= 0:
        # get end index of a tail
        tail_end_index = tail_start_index + len(tail)
        # extract ivxv_voting_facts_crontab block from a data_stripped
        return data_stripped[header_start_index:tail_end_index], True
    else:
        if header_start_index < 0:
            raise ValueError("malformed ### block ivxv_voting_facts_crontab ###")
        else:
            raise ValueError("malformed ### endblock ivxv_voting_facts_crontab ###")


def install_voting_facts_crontab():
    """Install crontab for ivxv-storageorder automation."""
    subprocess.run(
        ["env", "VISUAL=ivxv-voting-facts-crontab", "crontab", "-e"], check=True)


def voting_facts_util():
    """Launch log analyzer on IVXV Logmonitor to search for missing voting facts, then
    export these facts as .csv file and run 'ivxv-storageorder' on a first Collector
    host, where operation succeeds, in case of unsuccessful operation it will choose
    another Collector host, until succeeds, or returns an error."""
    init_cli_util(
        """
        Launch log analyzer on IVXV Logmonitor to search for missing voting facts, then
        export these facts as .csv file and then run 'ivxv-storageorder' on a first
        Collector host, where operation succeeds, in case of unsuccessful operation it
        will choose another Collector host, until succeeds, or returns an error.

        Usage: ivxv-voting-facts
        """
    )

    filename = f"votesorder-{time.time_ns()}.csv"

    with IVXVManagerDb() as db:
        collector_state = db.get_value("collector/state")
        services = db.get_all_values("service")
        logmon_host = db.get_value('logmonitor/address')
        if not logmon_host:
            log.error("Log monitor is not defined")
            return 1
    if collector_state != COLLECTOR_STATE_CONFIGURED:
        if not sys.stdout.isatty():  # suppress warning if executed from crontab
            return 0
        log.warning("Collector is not configured")
        return 1

    votesorder_services = {}

    # Find all votesorder service host addresses
    for service_id, service in services.items():
        service_type = service["service-type"]
        is_votesorder_service = service_type == "votesorder"
        if is_votesorder_service and service["state"] != SERVICE_STATE_REMOVED:
            votesorder_services[service_id] = service["ip-address"].split(":")[0]

    if len(votesorder_services.keys()) < 1:
        log.error("No votesorder services found")
        return 1

    # Prepare missing voting facts .csv file on first succeeded Logmonitor host
    logmon_account = f"logmon@{logmon_host}"
    proc = exec_remote_cmd(['ssh', logmon_account,
                            "ivxv-storageorder-prepare-csv.sh",
                            "--filename", filename])
    if proc.returncode:
        logmon_account = ""
        if proc.returncode == 2:
            log.info("No missing voting facts found")
            return 0
        log.error(
            f"Error while preparing .csv: {proc}")

    if logmon_account == "":
        log.error(f"Logmonitor host is not active {proc}")
        return 1

    # Copy .csv file from succeeded Log host to first active Collector host
    collector_account = ""
    for service_host in votesorder_services.values():
        collector_account = f"ivxv-votesorder@{service_host}"
        proc = exec_remote_cmd(
            ["scp", f"{logmon_account}:{filename}", f"{collector_account}:{filename}"])
        if proc.returncode:
            collector_account = ""
            log.error(
                f"Cannot copy .csv file from {logmon_account} to {collector_account} "
                f"due to an error: {proc}")
            continue

    if collector_account == "":
        log.error(f"No active Collector hosts among {votesorder_services.values()}")
        return 1

    # Add missing voting facts using ivxv-storageorder
    for service_id, service_host in votesorder_services.items():
        collector_account = f"ivxv-votesorder@{service_host}"
        proc = exec_remote_cmd(
            ["ssh", f"{collector_account}",
             "ivxv-storageorder", "-file", f"{filename}",
             "-instance", f"{service_id}"])
        if proc.returncode:
            collector_account = ""
            log.error(f"Running ivxv-storageorder failed with error: {proc}")
            continue

    if collector_account == "":
        log.error(f"ivxv-storageorder failed on {votesorder_services.values()} hosts")
        return 1

    return 0
