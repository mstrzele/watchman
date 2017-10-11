# Watchman 🕵️

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![CircleCI branch](https://img.shields.io/circleci/project/github/mstrzele/watchman/master.svg)](https://github.com/mstrzele/watchman)

> Watches for new releases on GitHub and sends a message on Slack

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Background

Since GitHub doesn't have a feature for watching new releases (more details in [isaacs/github#410](https://github.com/isaacs/github/issues/410)), I've decided to create my own bot for Slack during the Hack Day 2017 at [@Schibsted-Tech-Polska](https://github.com/Schibsted-Tech-Polska).

There's also a smiliar project, [mystor/gh-release-watch](https://github.com/mystor/gh-release-watch) (written in JavaScript). It sends the notifications by e-mail.

## Install

```
$ go install github.com/mstrzele/watchman
```

## Usage

```
$ watchman -t YOUR_TOKEN_HERE
```

## Maintainers

[@mstrzele](https://github.com/mstrzele)

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2017 Maciej Strzelecki
