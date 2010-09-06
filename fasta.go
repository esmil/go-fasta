/*
 * This file is part of go-fasta.
 *
 * go-fasta is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of
 * the License, or(at your option) any later version.
 *
 * go-fasta is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with go-fasta.  If not, see <http://www.gnu.org/licenses/>.
 */
package fasta

import (
	"os"
	"io"
)

const (
	A = iota
	B
	C
	D
	E
	F
	G
	H
	I
	K
	L
	M
	N
	O
	P
	Q
	R
	S
	T
	U
	V
	W
	X
	Y
	Z
	GAP
	INVALID_CHARACTER
)

const _SHIFT = 5

const (
	cCH = iota << _SHIFT
	cST
	cCM
	cET
	cNL
	_CLASSES = iota
)

const (
	sN1 = iota
	sN2
	sCM
	sTX
	sWS
	sSQ
	_STATES
	xTX
	xDN
)

var character_class [256]byte = [256]byte{
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cNL, cET, cET, cNL, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,

	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cST, cET, cET, GAP, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cCM, cET, cET, cCM, cET,

	cET, A,   B,   C,   D,   E,   F,   G,
	H,   I,   cET, K,   L,   M,   N,   O,
	P,   Q,   R,   S,   T,   U,   V,   W,
	X,   Y,   Z,   cET, cET, cET, cET, cET,

	cET, A,   B,   C,   D,   E,   F,   G,
	H,   I,   cET, K,   L,   M,   N,   O,
	P,   Q,   R,   S,   T,   U,   V,   W,
	X,   Y,   Z,   cET, cET, cET, cET, cET,

	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,

	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,

	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,

	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
	cET, cET, cET, cET, cET, cET, cET, cET,
}

var transition [_STATES][_CLASSES]byte = [_STATES][_CLASSES]byte{
	//                       cCH  cST  cCM  cET  cNL
	/* sN1 */ [_CLASSES]byte{xTX, xTX, sCM, xTX, sN1},
	/* sN2 */ [_CLASSES]byte{sSQ, xDN, xDN, sWS, sN2},
	/* sCM */ [_CLASSES]byte{sTX, sTX, sTX, sTX, sN1},
	/* sTX */ [_CLASSES]byte{sTX, sTX, sTX, sTX, sN1},
	/* sWS */ [_CLASSES]byte{sSQ, xDN, sWS, sWS, sN2},
	/* sSQ */ [_CLASSES]byte{sSQ, xDN, sWS, sWS, sN2},
}

func CharToSymbol(char byte) byte {
	sym := character_class[char]
	if sym >= INVALID_CHARACTER {
		return INVALID_CHARACTER
	}
	return sym
}

func SymbolToChar(sym byte) byte {
	switch {
	case sym < K:
		return sym + 'A'
	case sym < GAP:
		return sym + 'A' + 1
	case sym < INVALID_CHARACTER:
		return '-'
	}

	return ' '
}

type FASTA struct {
	Text string
	Data []byte
}

type Parser struct {
	state byte
	i     int
	buf   []byte
	text  string
}

func (p *Parser) growBuffer(i int) []byte {
	newBuf := make([]byte, 2*(i+128))
	copy(newBuf, p.buf)
	p.buf = newBuf
	return newBuf
}

func (p *Parser) Feed(data []byte) (done bool, leftover []byte) {
	state, buf, i := p.state, p.buf, p.i

out:
	for k := range data {
		char := data[k]
		class := character_class[char]

	again:
		state = transition[state][class>>_SHIFT]
		//print(k, ": ", class, " -> ", state, "\n")
		switch state {
		case sCM:
			if i >= len(buf) {
				buf = p.growBuffer(i)
			}
			buf[i] = '\n'
			i++

		case sTX:
			if i >= len(buf) {
				buf = p.growBuffer(i)
			}
			buf[i] = char
			i++

		case xTX:
			if i == 0 {
				i = 1
			}
			p.text = string(buf[1:i])
			i = 0
			state = sN2
			goto again

		case sSQ:
			if i >= len(buf) {
				buf = p.growBuffer(i)
			}
			buf[i] = class
			i++

		case xDN:
			done = true
			leftover = data[k+1:]
			break out
		}
	}

	p.state, p.buf, p.i = state, buf, i
	return
}

func (p *Parser) Result(f *FASTA) {
	text, buf, i := p.text, p.buf, p.i

	switch p.state {
	case sN1: fallthrough
	case sCM: fallthrough
	case sTX:
		if i == 0 {
			i = 1
		}
		text = string(buf[1:i])
		i = 0
	}

	res := make([]byte, i)
	copy(res, buf[0:i])

	f.Text, f.Data = text, res
	p.state, p.i, p.text = sN1, 0, ""
}

func ParseOne(input io.Reader, f *FASTA) os.Error {
	var (
		buf [1024]byte
		p   Parser
	)

	for {
		n, err := input.Read(buf[0:])
		if err != nil {
			if err == os.EOF {
				break
			}
			return err
		}

		if done, _ := p.Feed(buf[0:n]); done {
			break
		}
	}

	p.Result(f)

	return nil
}

func ParseAll(input io.Reader) ([]*FASTA, os.Error) {
	var (
		buf [1024]byte
		p   Parser
	)

	ret, i := make([]*FASTA, 32), 0
	for {
		n, err := input.Read(buf[0:])
		if err != nil {
			if err == os.EOF {
				break
			}
			return ret[0:i], err
		}

		for slice := buf[0:n]; slice != nil; {
			var done bool

			done, slice = p.Feed(slice)
			if done {
				f := &FASTA{}
				p.Result(f)
				if i >= len(ret) {
					newRet := make([]*FASTA, 2*i)
					copy(newRet, ret)
					ret = newRet
				}
				ret[i] = f
				i++
			}
		}
	}

	f := &FASTA{}
	p.Result(f)
	if i >= len(ret) {
		newRet := make([]*FASTA, i+1)
		copy(newRet, ret)
		ret = newRet
	}
	ret[i] = f
	i++

	return ret[0:i], nil
}

func (f *FASTA) String() string {
	text, data := f.Text, f.Data

	i := len(f.Data)
	i += (i / 70) + len(text) + 5
	buf := make([]byte, i)

	i = 0
	buf[i] = '>'; i++
	buf[i] = ' '; i++
	for k := 0; k < len(text); k++ {
		buf[i] = text[k]
		i++
	}
	for k := 0; k < len(data); k++ {
		if k%70 == 0 {
			buf[i] = '\n'
			i++
		}
		buf[i] = SymbolToChar(data[k])
		i++
	}
	buf[i] = '*'; i++
	buf[i] = '\n'

	return string(buf)
}
