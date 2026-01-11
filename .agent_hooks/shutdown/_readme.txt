`shutdown` hooks are run when the sandbox stops.
The sandbox shutdown process waits on the execution of all hooks before continuing.
All `shutdown` hooks are collectively given a timeout of 60 seconds.
If the `shutdown` hooks time out, the process is forcibly killed and shutdown continues.
