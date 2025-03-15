# Ignite App: CsClient

Generaetes a C# client library for your cosmos blockchain.

## Installation

### Within Project Directory

To use the "cs-client" app within your project, execute the following command inside the project directory:

```bash
ignite app install github.com/ignite/apps/cs-client
```

The app will be available only when running `ignite` inside the project directory.

### Globally

To use the "cs-client" app globally, execute the following command:

```bash
ignite app install -g github.com/ignite/apps/cs-client
```

This command will compile the app and make it immediately available to the `ignite` command lists.

## Requirements

- Go (version 1.23.6 or higher)
- Ignite CLI (version 28.1.1 or higher)

## Usage

```bash
ignite generate cs-client
```

### Configuration
```
Flags:
  -h, --help         help for cs-client
  -o, --out string   csharp output directory
  -y, --yes          answers interactive yes/no questions with yes
```

Also the output directory can be set via the `config.yaml`, like so:
```yaml
version: 1
client:
  cs-client:
    path: path/to/where/ur/lib/should/be
```
