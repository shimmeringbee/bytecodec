# Shimmering Bee: Byte Codec

[![license](https://img.shields.io/github/license/shimmeringbee/bytecodec.svg)](https://github.com/shimmeringbee/bytecodec/blob/master/LICENSE)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg)](https://github.com/RichardLitt/standard-readme)
[![Actions Status](https://github.com/shimmeringbee/bytecodec/workflows/test/badge.svg)](https://github.com/shimmeringbee/bytecodec/actions)

> Implementation of a byte codec to Marshall/Unmarshall structs to []byte, compatible with Zigbee types, written in Go.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Background

bytecodec is a library to marshall and unmarshall Go structs to a []byte, for transmission on the wire. Compatible with Zigbee types, such as 24 bit integers. Includes features such as:

* implicit or explicit length counts prepending arrays, counting in items or bytes
* bitmasks for enumerated types
* nested structs

## Install

Add an import and most IDEs will `go get` automatically, if it doesn't `go build` will fetch.

```go
import "github.com/shimmeringbee/bytecodec"
```

## Usage

## Maintainers

[@pwood](https://github.com/pwood)

## Contributing

Feel free to dive in! [Open an issue](https://github.com/shimmeringbee/bytecodec/issues/new) or submit PRs.

All Shimmering Bee projects follow the [Contributor Covenant](http://contributor-covenant.org/version/1/3/0/) Code of Conduct.

## License

   Copyright 2019 Peter Wood

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.