# oas2md

oas2md is a simple program to create markdown documentation of an OpenAPI defined API.  It uses [libopenapi](https://github.com/pb33f/libopenapi) to parse the specification and go templates for rendering the pages.

The program was designed to create pages for a [Hugo](https://gohugo.io/) website but it can also create pages more suitable for display on Github by passing the `-g` option.


## Who this project is for
This project is intended for developers who want to create static documentation from their OpenAPI specifications.


## Project dependencies
Before using oas2md, ensure you have:
* An OpenAPI specification.  This can be a file or a URL to a spec.
* golang 1.18+ if developing

When using the binary the templates are embedded.

## Instructions for using oas2md
Get started with oas2md by installing the program.


### Installing
1. Download the latest binary from [Github](https://github.com/tjdavis3/oas2md/releases)

### Running

Oas2md will create markdown files from your OpenAPI specification.  All configuration is done via command-line switches.  By default it will create Hugo-compatible files with the main file being `_index.md`.

**Command-line switches:**

`-d`
: The directory in which to create the file(s)

`-s`
: The location of the specification.  Can be a filepath or URL

`-g`
: makes Github-compatible files.  The main page is README.md

`-1`
: creates a single file for the entire API rather than a file per path



### Troubleshooting

The program will panic if there are any issues creating files, rendering templates, etc.  Most errors can be traced back to issues with the specification.  The program attempts to validate the spec before running, but some issues can still get through.  It is recommended to run the spec through [Vacuum](https://quobix.com/vacuum/) and clean up any issues before attempting to generate the documentation.

<!--
Other troubleshooting supports:
* {Link to FAQs}
* {Link to runbooks}
* {Link to other relevant support information}
-->

## Contributing guidelines
If you wish to contribute to oas2md fork the project and submit a pull request.  

<!--
## Additional documentation
{Include links and brief descriptions to additional documentation. Examples provided in README template guide.}

For more information:
* Reference link 1
* Reference link 2
* Reference link 3...
-->

## How to get help
Requests for help can be made by starting a [discussion](https://github.com/tjdavis3/oas2md/discussions) on Github.  Problems can be reported by creating an [issue](https://github.com/tjdavis3/oas2md/issues/new/choose). 

## Terms of use
oas2md is licensed under the [MIT License](LICENSE).

---

