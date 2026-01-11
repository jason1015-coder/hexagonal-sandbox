#!/usr/local/bin/python3
"""
Run all hooks in a given directory by hook type.
The hook type is the name of the directory in the `BASE_DIR` directory (default: /workspace/.agent_hooks).
The logs are stored in the `LOG_DIR` directory (default: /workspace/.agent_hooks/<hook_type>/_logs).
The hooks should be named like `00_<hook_name>.py`, `01_<hook_name>.py`, etc.
The hooks will be sequentially run in order of the number prefix.
Files with `_` prefixed are ignored.

Usage:
    run_all_hooks.py <hook_type>

    Example:
        run_all_hooks.py shutdown

This will run all hooks in the `shutdown` directory.
"""

import logging
import os
import subprocess
import sys
from pathlib import Path

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

BASE_DIR = Path("/workspace/.agent_hooks")
LOG_DIR_NAME = "_logs"
DEFAULT_HOOK_TIMEOUT_MINUTES = 5


def execute_hook(hook_path: Path, log_dir: Path):
    """Execute a hook

    Args:
        hook_path (Path): The path to the hook executable
    """
    try:
        log_path = log_dir / hook_path.with_suffix(".log").name
        logger.info(f"Running hook {hook_path} with log file {log_path}")
        with open(log_path, "a", buffering=1) as log_file:
            result = subprocess.run(
                [sys.executable, hook_path],
                cwd=hook_path.parent,
                timeout=DEFAULT_HOOK_TIMEOUT_MINUTES * 60,
                text=True,
                stdout=log_file,
                stderr=subprocess.STDOUT,
            )
            # Flush and sync the log file to ensure all content is written to disk
            log_file.flush()
            os.fsync(log_file.fileno())

        logger.info(f"Hook {hook_path} completed with exit code {result.returncode}")
    except subprocess.TimeoutExpired:
        logger.error(f"Hook {hook_path} timed out after {DEFAULT_HOOK_TIMEOUT_MINUTES} minutes")
    except Exception:
        logger.exception(f"Error running hook {hook_path}")


def run_all_hooks(hook_type: str):
    """Run all hooks of a given type

    Args:
        hook_type (str): The type of hook to run
    """
    # Check if the hook directory exists and is a directory
    hook_dir = BASE_DIR / hook_type
    if not hook_dir.exists() or not hook_dir.is_dir():
        raise FileNotFoundError(f"Hooks for type {hook_type} in directory {hook_dir} do not exist")
    # Create log directory if it doesn't exist
    log_dir = hook_dir / LOG_DIR_NAME
    log_dir.mkdir(parents=True, exist_ok=True)
    # Run all hooks in the hook directory
    for hook in sorted(hook_dir.glob("*.py")):
        if hook.name.startswith("_"):
            continue
        execute_hook(hook_path=hook, log_dir=log_dir)

    logger.info(f"Completed running all {hook_type} hooks")


def main():
    if len(sys.argv) != 2:
        # Special case for usage error, print to stderr
        sys.stderr.write("Usage: run_all_hooks.py <hook_type>\n")
        sys.exit(1)
    run_all_hooks(hook_type=sys.argv[1])


if __name__ == "__main__":
    try:
        main()
    except Exception:
        logger.exception(f"Error running hooks")
        sys.exit(1)
