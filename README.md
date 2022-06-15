# SKM

> A ssh private key manager.

# Download

Please go to [release](https://github.com/lovexy-fun/skm/releases)

# Usage

Use `skm` to show the current key.

input

```shell
skm
```

output

```shell
Effective key: testkey
```

Use `-h` for help

input

```shell
skm -h
```
output

```shell
SSH key manager

Usage:
  skm [flags]
  skm [command]

Available Commands:
  add         Add a key to manager
  del         Delete a key from manager
  help        Help about any command
  ls          List all keys
  sel         Choose a key to make it effective

Flags:
  -h, --help   help for skm

Use "skm [command] --help" for more information about a command.
```

## add

Add a ssh key to manager.

For example:

```shell
skm add -f ./id_rsa -n testkey
```

## del

Delete a ssh key from manager.

For example:

```shell
skm del
```

This command use promptui, you need use ↑/↓ to select key.

## ls

List all keys.

For example:
```shell
skm ls
```

## sel

Select a key to make it effective.

For example:

```shell
skm sel
```

This command use promptui, you need use ↑/↓ to select key.

# Other

When I finished it, I found a better [skm](https://github.com/TimothyYe/skm).