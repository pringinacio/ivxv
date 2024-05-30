# IVXV Internet voting framework
"""Backup service management helper."""

import subprocess


def install_backup_crontab():
    """Install crontab for backup automation."""
    subprocess.run(["env", "VISUAL=ivxv-backup-crontab", "crontab", "-e"], check=True)
