# CircleCI Context Secret Manager

A tool for managing your CircleCI context secrets (environment variables) as idempotent configuration.

## Features

- Extension - you can base one contexts secrets on another (or many others) with overrides.
- Know exactly what all your secrets are at any time. The tool uses an idempotent methadology to ensure that the secrets in a context exactly match what the tool knows about.

## Recommendations

- Use source control for your secrets. Secrets are a form of configuration, store them somewhere that is versioned with a good change history.
- Using git for secrets is great but make sure that your secrets are encrypted. Check out [git-crypt](https://github.com/AGWA/git-crypt).

## Usage

Using the tool should be easy, check out our [examples](/examples) configuration files.

The CLI is self documenting, if something isn't obvious please raise an issue and we will address it.

```sh
ccsm --help
```

### Installing

#### GoLang

```sh
go get github.com/armakuni/circleci-context-secret-manager/cmd/ccsm
```

#### Manual

1. Download the [release](https://github.com/armakuni/circleci-context-secret-manager/releases) for your OS - if we don't provide one then you can compile it yourself or let us know and we will look at adding a binary for you.
2. Rename the binary to `ccsm` or `ccsm.exe` on Windows.
3. Make the binary executable.
4. Put the binary somewhere on your `PATH`

### Configuration reference

All configuration files should be stored in a single folder, this defaults to `contexts` from wherever the tool is run, but can be overriden with the `--contexts` flag or the `CONTEXTS_DIR` environment variable.

The configuration for each yaml file should follow:

| Key           | Type                | Required                         | Default | Description                                                                                                                                                                                          |
|---------------|---------------------|----------------------------------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `context_id`  | `string`            | Yes (unless skip_deploy is true) | N/A     | The context ID your configuration is for. **Note**: The file name must also match the context name (context `main` would need a file of `main.yml`)                                                  |
| `skip_deploy` | `bool`              | No                               | false   | If true the context file can be used by `extends` but will not update a context, useful if you want to share a set of secrets between two contexts and do not need the secrets in a specific context |
| `extends`     | `[]string`          | No                               | N/A     | An array of context file names to extend, it will load all `secrets` from extended files in order and then override with any `secrets` defined locally                                               |
| `secrets`     | `map[string]string` | No                               | N/A     | The secrets you wish to configure for your context, this will override anything imported via `extends`. **Note**: Leaving `secrets` blank will delete all secrets on an `apply`                      |

#### Projects

While this tools main focus remains on managing environment variables on a context you might find times where you need some configuration to be managed on projects, you can find a guide to using the tool with projects [here](/GUIDE.md#Projects)

### Getting started

1. Create some contexts in CircleCI
2. Get the IDs of your contexts
    1. Browse to your orgs contexts, something like <https://ui.circleci.com/settings/organization/github/armakuni/contexts>
    2. Click on the context you want the ID for.
    3. The URL bar will include the ID after `contexts`: `https://ui.circleci.com/settings/organization/github/armakuni/contexts/3e614332-9bbd-4b02-ab30-b5d18b11ae01`
3. Set the ID in your yaml file `context_id`

    > **Note**: Your yaml file name has to match the context name, this is for a safety check to ensure you have got the ID for the correct context.

    ```yaml
    extends:
    - main.yml

    context_id: <your_context_id>

    secrets:
      PASSWORD: fwibble
    ```

4. Check your yaml file will plan `ccsm plan --context dev.yml`
5. Run a `dry-run` to see what will change on the remote `ccsm -t <api_token> dry-run`

    > **Note**: Dry run will only show what secrets are being `added/deleted/updated` even if nothing has changed all the secrets will show up as `updated` as the CircleCI APIs return a truncated value so it is not possible to compare.
