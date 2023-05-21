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
![print](https://drive.google.com/uc?export=view&id=1wN4jnYLMplJT6ftiYI1g4paVod3_PKmC)

```sh
awscred check
```
![check](https://drive.google.com/uc?export=view&id=1zCOWEY5RPJLKWP_Id_gtIR0Q5DhpJeYf)

```sh
awscred apply workshop-admin
```
![apply](https://drive.google.com/uc?export=view&id=1PZnsBh-3hSiLaVyWWgAhyB5eYSVe913o)

```sh
awscred help
```
![help](https://drive.google.com/uc?export=view&id=1UmBDAT13hC1r20okrRVtTWbgCExkPPY5)

# License

This package is licensed under the MIT License - see the LICENSE file for details.