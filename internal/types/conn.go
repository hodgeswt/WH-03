package types

import (
	"fmt"
	"strconv"

	"github.com/hodgeswt/WH-03/internal/util"
)

type IConn interface {
	GetName() string
	SetName(name string)

	GetBitwidth() int
	SetBitwidth(bitwidth int)

	GetData() int
	GetDataStr() string

	SetData(data int) error
	SetDataStr(data string) error
}

type Conn struct {
	name     string
	bitwidth int
	data     int
}

func (it *Conn) GetName() string {
    return it.name
}

func (it *Conn) SetName(name string) {
    it.name = name
}

func (it *Conn) GetBitwidth() int {
    return it.bitwidth
}

func (it *Conn) SetBitwidth(bitwidth int) {
    it.bitwidth = bitwidth
}

func (it *Conn) GetData() int {
    return it.data
}

func (it *Conn) GetDataStr() string {
    return fmt.Sprintf("%08b", it.data)
}

func (it *Conn) maxData() int {
    return (2 ^ it.bitwidth) - 1
}

func (it *Conn) SetData(data int) error {
    max := it.maxData()

    if data > max {
        return ErrDataTooLarge
    }

    it.data = data
    return nil
}

func (it *Conn) SetDataStr(data string) error {
    if !util.ContainsOnly(data, bitchars) {
        return ErrInvalidChars
    }

    num64, err := strconv.ParseInt(data, 2, 64)
    if err != nil {
        return err
    }

    max := it.maxData()
    num := int(num64)

    if num > max {
        return ErrDataTooLarge
    }

    it.data = num
    return nil
}

var bitchars = map[rune]bool{'0': true, '1': true}
