# with-env

Written as a counterpart to [pjs](https://github.com/IngCr3at1on/x/tree/master/pjs), provides a godotenv wrapper for resolving a single env file out of a directory of many based on aliases in the file comments.

For example adding the comment `# with-env aliases: sandbox` would cause with-env to select this file with if provided the input `sandbox`.

Supports multiple aliases per file, eq: `# with-env aliases: dev,develop`.

Leverages godotenv for _some_ functionality while mimicking other non-exported functionality.
