# AWS Credentials Package

The `awscred` package is an easy to use CLI for checking what [AWS](https://aws.amazon.com/) credentials are available in your `.aws/credentials` file and which of these credentials are valid.

Easily apply and run a command within the context of the desired profile.


## Features

- View available credentials.
- Remove expired or invalid credentials from both `.aws/credentials` and `.aws/config`.
- Apply a profile and use it to run a command with that profile set.

# Installation

```sh
go install github.com/mauriceLC92/awscred
```

## Usage

```sh
awscred print
```
![](https://drive.google.com/file/d/1wN4jnYLMplJT6ftiYI1g4paVod3_PKmC/view?usp=sharing)


# License

This package is licensed under the MIT License - see the LICENSE file for details.