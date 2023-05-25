# base16-builder-go

<!--toc:start-->

- [base16-builder-go](#base16-builder-go)
  - [Building](#building)
  - [Commands](#commands)
  - [Notes](#notes) - [Additional variables](#additional-variables)
  <!--toc:end-->

A simple builder for base16 templates and schemes.

This currently implements version 0.10.0 of the [base16_spec].

## Building

Currently version 1.18 or higher of the Go compiler is needed.

Because the schemes are stored in a separate repo, the schemes repo needs to
be cloned before building.

The following command will clone the schemes directory:

```text
git clone https://github.com/tinted-theming/base16-schemes.git schemes
```

Now that the repo is cloned, you can use `go build` to create a binary. You may
wish to update the schemes directory to get new included schemes.

## Commands

By default, this builder will build the template in the current directory using
the compiled-in schemes. If you want to update schemes independently, you can
use the `-schemes-dir` flag to point to another directory.

```text
Usage of base16-builder-go:
  -schemes-dir string
    Target directory for scheme data.
    The default value uses internal schemes. (default "-")
  -template-dir string
    Target template directory to build. (default ".")
  -verbose
    Log all debug messages
```

## Notes

I'm open to making a few template-specific tweaks as long as they'll be useful
to other templates. Below is a listing of the additions to the base16 spec which
this builder supports.

### Additional variables

`scheme-slug-underscored` - A version of the scheme slug where dashes have
been replaced with underscores.

## References

[base16_spec]: https://github.com/tinted-theming/home
