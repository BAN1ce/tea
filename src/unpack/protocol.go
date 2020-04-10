package unpack

import "bufio"

type Protocol func() bufio.SplitFunc
