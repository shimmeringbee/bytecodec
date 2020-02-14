# Shimmering Bee: Byte Codec

[![license](https://img.shields.io/github/license/shimmeringbee/bytecodec.svg)](https://github.com/shimmeringbee/bytecodec/blob/master/LICENSE)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg)](https://github.com/RichardLitt/standard-readme)
[![Actions Status](https://github.com/shimmeringbee/bytecodec/workflows/test/badge.svg)](https://github.com/shimmeringbee/bytecodec/actions)

> Implementation of a byte codec to Marshal/Unmarshal structs to []byte, compatible with Zigbee types, written in Go.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Background

bytecodec is a library to marshal and unmarshal Go structs to a []byte, for transmission on the wire. Compatible with Zigbee types, such as 24 bit integers.

## Install

Add an import and most IDEs will `go get` automatically, if it doesn't `go build` will fetch.

```go
import "github.com/shimmeringbee/bytecodec"
```

## Usage

### Marshalling

Marshalling assumes integers should be expressed as little endian, unless overridden.

Currently supports:
* uint8, uint16, uint32, uint64
* struct
* array/slice
* string (null terminated and length prefixed)

```go
type StructToMarshal struct {
    ByteField              byte
    ArrayOfUint16BigEndian []uint16 `bcendian:"big"`
    Uint16LittleEndian     uint16
}

data := &StructToMarshal{
    ByteField:              0x55,
    ArrayOfUint16BigEndian: []uint16{0x8001, 0x1234},
    Uint16LittleEndian:     0x2233,
}

bytes, err := bytecodec.Marshal(data)

if err != nil {
    // Handle Error
}

// bytes = []byte{0x55,0x80,0x01,0x12,0x34,0x33,0x22}
```

## Maintainers

[@pwood](https://github.com/pwood)

## Contributing

Feel free to dive in! [Open an issue](https://github.com/shimmeringbee/bytecodec/issues/new) or submit PRs.

All Shimmering Bee projects follow the [Contributor Covenant](https://shimmeringbee.io/docs/code_of_conduct/) Code of Conduct.

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