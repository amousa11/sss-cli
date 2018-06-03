# sss-cli
Shamir Secret Sharing Command Line Interface

A simple command line interface for generation and recovery of a shared secret. 

## Usage

This tool has 2 main functions:
- Generation of Shares of a secret
- Recovery of a secret, using those shares

### Generating a Secret and Shares

To generate shares, run the following command:

`./sss-cli generate minNumShares numShares`

minNumShares is the minmum number of shares required to recover the secret

numShares is the total number of shares generated

*minNumShares must never be greater than the total number of shares*

After running this command, the output will be printed to the console and will look like this:

```
>./sss-cli generate 4 8
FIELD MODULUS = fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f
Generating  8 shares with threshold 4  for recovery:
Secret :  5b89a30d37b59e9a849c7f114d0215e415db6d082703898e0813ec4c9bdd9a3
Shares :
1        16741fc130afe8e324135418f1bbb88137676baa635c4ae5da583489077256f3
2        fa6717b6f1ee01bc4e94da1a6363c5ca9c47228b4acdc6d270551335a855a484
3        573cfe4082b489aa9224bec5c450a1d38c3dbd6fae5c7bc3c9c559e35f7e4919
4        d2a14f8c4e8265e3591966eb6f0aa535238c1e5403a0391f0df687a5e002b860
5        123f87c8c0d67b9c0dc9375bbe1a28887e732734c030ce4964361b95dcf9791c
6        bbc32324452fb00a1a8a94e70c078466b933ba0e59a60aa7f3d194c70978fefb
7        74d79dcd470ce862e9b3e45db35b1168f00eb8dd4597bd9fe41672521897d0c0
8        e32873f231ed09dbe59b8a900e9d28283f45059df99db6965c52334bbd6c65ea
```

Also, a folder named `shares` will be generated in the same folder as the executable. It will contain text files named by numbers. Each file corresponds to a single share which a user would distribute amongst individuals sharing this secret

```
the-top-folder
    sss-cli
    shares
        1
        2
        3
        4
        5
        ...
```

### Recovering a Secret

To recover a secret, first create a folder named `shares` in the same directory as the executable so your directory looks like this:

```
the-top-folder
    sss-cli
    shares
        1
        2
        3
        4
        5
        ...
```

Put share files in the `shares` folder like in the example above.

Now run the following command:

`./sss-cli recover`

The output should look like this:

```
Secret successfully recovered from 8 shares: 5b89a30d37b59e9a849c7f114d0215e415db6d082703898e0813ec4c9bdd9a3
```

## Help

The following is the help text for the cli

```
NAME:
   Generation of Shares and Recovery of Secrets - Shamir Secret Sharing - Generate Shares and Recover Secrets using a threshold of those shares

USAGE:
   sss-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     generate, g, gen  generate [minNumShares] [numShares] - generates numShares shares of a secret which can only be recovered by numMinShares shares
     recover, c        recover [pathToSharesFolder] - recovers a secret given a path to a folder with the generated shares
     help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```