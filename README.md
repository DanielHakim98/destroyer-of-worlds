# Destroyer of Worlds

- It's just a normal load tester, I build it because I'm trying this challenge **[coding-challenges: load-tester](https://codingchallenges.fyi/challenges/challenge-load-tester/)**.

## Installation

- For now, you need to build it from source. **Golang 1.22** is required to build the binary. **Python 3.12.0** is used to set up mock server serves static html.

## Usage

```text
Usage:
destroyer-of-worlds [flags]

Flags:
-c, --concurrent int maximum concurent request. Default is 1 (default 1)
-h, --help help for destroyer-of-worlds
-n, --requests int The total requests to be sent. Default is 1 (default 1)
-t, --toggle Help message for toggle
-u, --url string URL to be tested.

```

To setup mock server, make sure current shell is at root project directory, and then run command below

```text
cd www
python -m http.server
```

## Contributing

Nah, I don't think it's worth it anyway. You can just copy and paste, steal it, or even recreate to do something better, robust and more interesting that this. Or just use **[Hey](https://github.com/rakyll/hey)**.
