#!/usr/local/bin/python3
"""
Track processes running on ports on the sandbox.
Filters out processes that are managed by the sandbox itself. (e.g. TTY terminal, VSCode server, etc.)
Generates a supervisor config file to track the processes in `SUPERVISORD_CONF_FILE` (default: /etc/supervisor/conf.d/_superninja_startup.conf).

Usage:
    python3 00_track_processes_on_ports.py
"""
import logging
import subprocess
import sys
from dataclasses import dataclass

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

PROCESS_NAMES_TO_IGNORE = {
    "daytona",  # When running Daytona some processes are started
    "sshd",  # SSH server for TTY terminal
}
PROCESS_PORTS_TO_IGNORE = {
    2222,  # TTY terminal internal
    3222,  # TTY terminal external for Nginx
    4000,  # VSCode server internal
    5000,  # VSCode server external for Nginx
    5900,  # x11VNC
    5901,  # x11VNC
    6080,  # noVNC
    8000,  # API server
    8002,  # Browser API server
    8080,  # HTTP server
    9222,  # Chrome remote debugging
}

# Supervisor config file
# This is the file that will be used to track the processes
SUPERVISORD_CONF_FILE = "/etc/supervisor/conf.d/_superninja_startup.conf"

# Supervisor config template
# This is the template that will be used to generate each process entry in the supervisor config file
SUPERVISORD_CONFIG_TEMPLATE = """[program:{SUPERVISOR_NAME}]
command=/bin/bash -c "source ~/.bashrc 2>/dev/null; exec {COMMAND}"
directory={CWD}
user={PROC_USER}
autostart=true
autorestart=true
startsecs=3
startretries=3
stderr_logfile=/var/log/supervisor/{SUPERVISOR_NAME}.err.log
stdout_logfile=/var/log/supervisor/{SUPERVISOR_NAME}.out.log
stdout_logfile_maxbytes=5MB
stdout_logfile_backups=2
stderr_logfile_maxbytes=5MB
stderr_logfile_backups=2
"""


@dataclass
class Process:
    port: int
    pid: int


@dataclass
class ProcessInfo:
    port: int
    pid: int
    command: str
    cwd: str
    proc_name: str
    proc_user: str


def _run_command(command: str):
    """Runs a command and returns the stdout

    Args:
        command (str): The command to run

    Returns:
        str: The stdout of the command
    """
    result = subprocess.run(
        command.split(),
        text=True,
        check=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
    )
    return result.stdout.strip()


def get_processes_listening_on_ports() -> list[Process]:
    """Get processes listening on ports using `lsof`

    Example `lsof` output:
    COMMAND PID     USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
    nginx   591     root    7u  IPv4    556      0t0  TCP *:5000 (LISTEN)

    Returns:
        list[Process]: A list of processes listening on ports with their port and PID
    """
    logger.info("Getting processes listening on ports by running lsof...")
    result = _run_command("lsof -nP -iTCP -sTCP:LISTEN")
    lines = result.split("\n")[1:]  # Skip header line
    processes: list[Process] = []
    for line in lines:
        command, pid, user, fd, type, device, size_off, node, name, _ = line.split()
        port = int(name.split(":")[-1])
        if command in PROCESS_NAMES_TO_IGNORE or port in PROCESS_PORTS_TO_IGNORE:
            continue
        processes.append(Process(port=port, pid=int(pid)))

    logger.info(f"Found {len(processes)} processes listening on ports")
    if len(processes) > 0:
        logger.info("\n".join([f"  {process.pid}: {process.port}" for process in processes]))
    return processes


def deduplicate_processes_by_port(processes: list[Process]) -> list[Process]:
    """Removes multiple processes listening on the same port, keeping the process with the lowest PID.

    Args:
        processes (list[Process]): A list of processes

    Returns:
        list[Process]: A list of processes, deduplicated by port
    """
    if len(processes) == 0:
        logger.info("No processes to deduplicate by port")
        return processes
    processes_by_port: dict[int, Process] = {}
    for process in processes:
        if (
            process.port not in processes_by_port
            or process.pid < processes_by_port[process.port].pid
        ):
            processes_by_port[process.port] = process
    deduplicated_processes = list(processes_by_port.values())
    logger.info(
        f"Deduplicated {len(processes)} processes by port, keeping {len(deduplicated_processes)} processes"
    )
    logger.info(
        "\n".join([f"  {process.pid}: {process.port}" for process in deduplicated_processes])
    )
    return deduplicated_processes


def deduplicate_processes_by_command(processes: list[ProcessInfo]) -> list[ProcessInfo]:
    """Removes multiple processes with the same command, keeping the first one (order doesn't matter).

    Args:
        processes (list[ProcessInfo]): A list of processes information

    Returns:
        list[ProcessInfo]: A list of processes information, deduplicated by command
    """
    if len(processes) == 0:
        logger.info("No processes to deduplicate by command")
        return processes
    processes_by_command: dict[str, ProcessInfo] = {}
    for process_info in processes:
        if process_info.command not in processes_by_command:
            processes_by_command[process_info.command] = process_info
    deduplicated_processes = list(processes_by_command.values())
    logger.info(
        f"Deduplicated {len(processes)} processes by command, keeping {len(processes_by_command)} processes"
    )
    logger.info(
        "\n".join(
            [
                f"  {process_info.command} ({process_info.pid}): {process_info.port}"
                for process_info in deduplicated_processes
            ]
        )
    )
    return deduplicated_processes


def get_process_info(process: Process) -> ProcessInfo | None:
    """Get additional process information using `ps` and `pwdx`

    Args:
        process (Process): A process

    Returns:
        ProcessInfo | None: The process information or None if the process information cannot be retrieved
    """
    try:
        logger.info(f"Getting process info for PID {process.pid}...")
        return ProcessInfo(
            port=process.port,
            pid=process.pid,
            command=_run_command(f"ps -p {process.pid} -o command="),
            # pwdx output format: "PID: /path"
            cwd=_run_command(f"pwdx {process.pid}").split(":")[-1],
            proc_name=_run_command(f"ps -p {process.pid} -o comm="),
            proc_user=_run_command(f"ps -p {process.pid} -o user="),
        )
    except Exception:
        logger.exception(f"Failed to get process info for PID {process.pid}")
        return None


def main():
    # Only track the lowest PID per port to avoid duplicates
    try:
        processes: list[Process] = get_processes_listening_on_ports()
        processes = deduplicate_processes_by_port(processes)
    except Exception:
        logger.exception(f"Error getting processes running on ports")
        sys.exit(1)

    # Track each unique command only once
    try:
        processes_info: list[ProcessInfo] = [
            process_info
            for process in processes
            if (process_info := get_process_info(process)) is not None
        ]
        processes_info = deduplicate_processes_by_command(processes_info)
    except Exception:
        logger.exception(f"Error getting process info")
        sys.exit(1)

    supervisor_config_entries: list[str] = []
    for process_info in processes_info:
        # Generate config entry for the process
        supervisor_name = f"{process_info.port}_{process_info.proc_name}"
        supervisor_config_entries.append(
            SUPERVISORD_CONFIG_TEMPLATE.format(
                SUPERVISOR_NAME=supervisor_name,
                COMMAND=process_info.command,
                CWD=process_info.cwd,
                PROC_USER=process_info.proc_user,
            )
        )

    # Write new supervisor config file
    try:
        with open(SUPERVISORD_CONF_FILE, "w") as f:
            f.write("\n".join(supervisor_config_entries))
    except Exception:
        logger.exception(f"Error writing supervisor config file")
        sys.exit(1)

    logger.info(
        f"Script executed successfully, tracked {len(supervisor_config_entries)} unique processes"
    )


if __name__ == "__main__":
    try:
        main()
    except Exception:
        logger.exception(f"Error tracking processes on ports")
        sys.exit(1)
