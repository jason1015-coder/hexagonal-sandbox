`startup` hooks are run when the sandbox starts.
The sandbox startup process executes the `startup` hooks in the background without a timeout.
`startup` hooks should not block the main runner process, and should put long work on background processes.
