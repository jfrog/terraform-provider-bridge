# How to Contribute

We'd love to accept your patches and contributions to this project. There are just a few guidelines you need to follow.

## Contributor License Agreement

Contributions to this project must be accompanied by a Contributor License Agreement. You (or your employer) retain the copyright to your contribution; this simply gives us permission to use and redistribute your contributions as part of the project. Head over to <https://cla.jfrog.com/> to see your current agreements on file or to sign a new one.

You generally only need to submit a CLA once, so if you've already submitted one (even if it was for a different project), you probably don't need to do it again.

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/jfrog/terraform-provider-bridge`

```sh
mkdir -p $GOPATH/src/github.com/jfrog
cd $GOPATH/src/github.com/jfrog
git clone git@github.com:jfrog/terraform-provider-bridge
```

Enter the provider directory and build the provider:

```sh
cd $GOPATH/src/github.com/jfrog/terraform-provider-bridge
go mod tidy
make build
```

Install locally for testing:

```sh
make install
```

## Testing

### Environment Setup

Set the following environment variables for testing:

```sh
export JFROG_URL="https://your-instance.jfrog.io"
export JFROG_ACCESS_TOKEN="your-admin-token"
```

### Unit Tests

```sh
make test
```

### Acceptance Tests

To run the full suite of Acceptance tests, run `make acceptance`.

*Note:* Acceptance tests create real resources. You should expect that the full acceptance test suite may take time to run.

```sh
make acceptance
```

## Generating Documentation

To generate documentation, run:

```sh
make doc
```

## Code Reviews

All submissions, including submissions by project members, require review. We use GitHub pull requests for this purpose. Consult [GitHub Help](https://help.github.com/articles/about-pull-requests/) for more information on using pull requests.

## Resources

The provider currently supports the following resource:

- `bridge` - Manages JFrog Bridge connections between Bridge Server and Bridge Client

For more details, see the [README](README.md) and [documentation](docs/index.md).