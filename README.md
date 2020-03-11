# CircleCI Context Secret Manager

A tool for managing your CircleCI context secrets (environment variables) as idempotent configuration.

## Features

- Extension - you can base one contexts secrets on another (or many others) with overrides.
- Know exactly what all your secrets are at any time. The tool uses an idempotent methadology to ensure that the secrets in a context exactly match what the tool knows about.

## Recommendations

- Use source control for your secrets. Secrets are a form of configuration, store them somewhere that is versioned with a good change history.
- Using git for secrets is great but make sure that your secrets are encrypted. Check out [git-crypt](https://github.com/AGWA/git-crypt).

## Usage

Using the tool should be easy, check out our [example](/example) configuration files.

The CLI is self documenting, if something isn't obvious please raise an issue and we will address it.

```sh
ccsm --help
```

### Getting started

1. Create some contexts in CircleCI
2. Get the IDs of your contexts
    1. Something like <https://ui.circleci.com/settings/organization/github/armakuni/contexts>
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
