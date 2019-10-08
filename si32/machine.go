//  Copyright (c) 2015 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package si32

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"unsafe"
)

type operation byte

const (
	opAccept      operation = 1
	opAcceptBreak operation = 2
	opMatch       operation = 3
	opMatchBreak  operation = 4
	opOutput      operation = 5
	opOutputBreak operation = 6
)

func (o operation) String() string {
	opName := []string{
		"UNDEF0",
		"ACCEPT",
		"ACCEPTB",
		"MATCH",
		"MATCHB",
		"OUTPUT",
		"OUTPUTB",
		"UNDEF7",
	}
	if int(o) >= len(opName) {
		return fmt.Sprintf("NA[%d]", o)
	}
	return opName[o]
}

type Program []Instruction
type Instruction [4]byte

// FST represents a finite State transducer VM.
type FST struct {
	Program []Instruction
	Data    []int32
}

// Configuration represents a FST Configuration.
type Configuration struct {
	PC      int     // program counter
	Head    int     // input head
	Outputs []int32 // outputs
}

const maxUint16 = 1<<16 - 1

type byteSlice []byte

func (p byteSlice) Len() int           { return len(p) }
func (p byteSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p byteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p Program) Reverse() {
	size := len(p)
	for i := 0; i < size/2; i++ {
		p[i], p[size-1-i] = p[size-1-i], p[i]
	}
}

func (m MAST) BuildFST() (t FST, err error) {
	var (
		prog Program
		data []int32
		code Instruction // tmp Instruction
	)
	var edges []byte
	addrMap := make(map[int]int)
	for _, s := range m.States {
		edges = edges[:0]
		for ch := range s.Trans {
			edges = append(edges, ch)
		}
		if len(edges) > 0 && !sort.IsSorted(byteSlice(edges)) {
			sort.Sort(byteSlice(edges))
		}
		for i, size := 0, len(edges); i < size; i++ {
			ch := edges[size-1-i]
			next := s.Trans[ch]
			addr, ok := addrMap[next.ID]
			if !ok && !next.IsFinal {
				err = fmt.Errorf("next addr is undefined: State(%v), input(%X)", s.ID, ch)
				return
			}
			jump := len(prog) - addr + 1

			var op operation
			out, ok := s.Output[ch]
			if !ok {
				if i == 0 {
					op = opMatchBreak
				} else {
					op = opMatch
				}
			} else {
				if i == 0 {
					op = opOutputBreak
				} else {
					op = opOutput
				}
			}

			if jump > maxUint16 {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(jump)
				prog = append(prog, code)
				jump = 0
			}
			if ok {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(out)
				prog = append(prog, code)
			}

			code[0] = byte(op)
			code[1] = ch
			p := unsafe.Pointer(&code[2])
			(*(*uint16)(p)) = uint16(jump)
			prog = append(prog, code)
		}
		if s.IsFinal {
			if len(s.Tail) > 0 {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(len(data))
				prog = append(prog, code)
				var tmp int32Slice
				for t := range s.Tail {
					tmp = append(tmp, t)
				}
				sort.Sort(tmp)
				data = append(data, tmp...)
				p = unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(len(data))
				prog = append(prog, code)
			}
			if len(s.Trans) == 0 {
				code[0] = byte(opAcceptBreak)
			} else {
				code[0] = byte(opAccept)
			}
			code[1], code[2], code[3] = 0, 0, 0 // clear
			if len(s.Tail) > 0 {
				code[1] = 1
			}

			prog = append(prog, code)
		}
		addrMap[s.ID] = len(prog)
	}

	prog.Reverse()
	t = FST{Program: prog, Data: data}
	return
}

// String returns debug codes of a fst virtual machine.
func (t FST) String() string {
	var (
		pc   int
		code Instruction
		op   operation
		ch   byte
		v16  uint16
		v32  int32
	)
	ret := ""
	for pc = 0; pc < len(t.Program); pc++ {
		code = t.Program[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*uint16)(unsafe.Pointer(&code[2])))
		switch operation(op) {
		case opAccept:
			fallthrough
		case opAcceptBreak:
			//fmt.Printf("%3d %v\t%X %d\n", PC, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%d %d\n", pc, op, ch, v16)
			if ch == 0 {
				break
			}
			pc++
			code = t.Program[pc]
			to := (*(*int32)(unsafe.Pointer(&code[0])))
			ret += fmt.Sprintf("%3d [%d]\n", pc, to)
			pc++
			code = t.Program[pc]
			from := (*(*int32)(unsafe.Pointer(&code[0])))
			ret += fmt.Sprintf("%3d [%d] %v\n", pc, from, t.Data[from:to]) //FIXME
		case opMatch:
			fallthrough
		case opMatchBreak:
			//fmt.Printf("%3d %v\t%02X %d\n", PC, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%02X(%c) %d\n", pc, op, ch, ch, v16)
			if v16 == 0 {
				pc++
				code = t.Program[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("%3d [%d]\n", PC, v32) //XXX
				ret += fmt.Sprintf("%3d jmp[%d]\n", pc, v32)
				//break
			}
		case opOutput:
			fallthrough
		case opOutputBreak:
			//fmt.Printf("%3d %v\t%02X %d\n", PC, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%02X(%c) %d\n", pc, op, ch, ch, v16)
			if v16 == 0 {
				pc++
				code = t.Program[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("%3d [%d]\n", PC, v32) //XXX
				ret += fmt.Sprintf("%3d jmp[%d]\n", pc, v32)
				//break
			}
			pc++
			code = t.Program[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			//fmt.Printf("%3d [%d]\n", PC, v32) //XXX
			ret += fmt.Sprintf("%3d [%d]\n", pc, v32)
		default:
			//fmt.Printf("%3d UNDEF %v\n", PC, code)
			ret += fmt.Sprintf("%3d UNDEF %v\n", pc, code)
		}
	}
	return ret
}

func (t *FST) run(input string) (snap []Configuration, accept bool) {
	var (
		pc  int       // program counter
		op  operation // operation
		ch  byte      // char
		v16 uint16    // 16bit register
		v32 int32     // 32bit register
		hd  int       // input head
		out int32     // output

		code Instruction // tmp Instruction
	)
	for pc < len(t.Program) && hd <= len(input) {
		code = t.Program[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*uint16)(unsafe.Pointer(&code[2])))
		//fmt.Printf("PC:%v,op:%v,Head:%v,v16:%v,Outputs:%v\n", PC, op, Head, v16, Outputs) //XXX
		switch op {
		case opMatch:
			fallthrough
		case opMatchBreak:
			if hd == len(input) {
				goto L_END
			}
			if ch != input[hd] {
				if op == opMatchBreak {
					return
				}
				if v16 == 0 {
					pc++
				}
				pc++
				continue
			}
			if v16 > 0 {
				pc += int(v16)
			} else {
				pc++
				code = t.Program[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("ex jump:%d\n", v32) //XXX
				pc += int(v32)
			}
			hd++
			continue
		case opOutput:
			fallthrough
		case opOutputBreak:
			if hd == len(input) {
				goto L_END
			}
			if ch != input[hd] {
				if op == opOutputBreak {
					return
				}
				if v16 == 0 {
					pc++
				}
				pc++
				pc++
				continue
			}
			pc++
			code = t.Program[pc]
			out = (*(*int32)(unsafe.Pointer(&code[0])))
			if v16 > 0 {
				pc += int(v16)
			} else {
				pc++
				code = t.Program[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("ex jump:%d\n", v32) //XXX
				pc += int(v32)
			}
			hd++
			continue
		case opAccept:
			fallthrough
		case opAcceptBreak:
			c := Configuration{PC: pc, Head: hd}
			pc++
			if ch == 0 {
				c.Outputs = []int32{out}
			} else {
				code = t.Program[pc]
				to := (*(*int32)(unsafe.Pointer(&code[0])))
				pc++
				code = t.Program[pc]
				from := (*(*int32)(unsafe.Pointer(&code[0])))
				c.Outputs = t.Data[from:to]
				pc++
			}
			//fmt.Printf("conf: %+v\n", c) //XXX
			snap = append(snap, c)
			if hd == len(input) {
				goto L_END
			}
			if op == opAcceptBreak {
				goto L_END
			}
			continue
		default:
			//fmt.Printf("unknown op:%v\n", op) //XXX
			return
		}
	}
L_END:
	//fmt.Printf("[[L_END]]PC:%d, op:%s, ch:[%X]\n", PC, op, ch) //XXX
	if hd != len(input) {
		return
	}
	if op != opAccept && op != opAcceptBreak {
		//fmt.Printf("[[NOT ACCEPT]]PC:%d, op:%s, ch:[%X]\n", PC, op, ch) //XXX
		return

	}
	accept = true
	//fmt.Printf("[[ACCEPT]]PC:%d, op:%s, ch:[%X]\n", PC, op, ch) //XXX
	return
}

// Search runs a finite State transducer for a given input and returns outputs if accepted otherwise nil.
func (t FST) Search(input string) []int32 {
	snap, acc := t.run(input)
	if !acc || len(snap) == 0 {
		return nil
	}
	c := snap[len(snap)-1]
	return c.Outputs
}

// PrefixSearch returns the longest commom prefix keyword and it's length in given input
// if detected otherwise -1, nil.
func (t FST) PrefixSearch(input string) (length int, output []int32) {
	snap, _ := t.run(input)
	if len(snap) == 0 {
		return -1, nil
	}
	c := snap[len(snap)-1]
	return c.Head, c.Outputs

}

// CommonPrefixSearch finds keywords sharing common prefix in given input
// and returns it's lengths and outputs. Returns nil, nil if there does not common prefix keywords.
func (t FST) CommonPrefixSearch(input string) (lens []int, outputs [][]int32) {
	snap, _ := t.run(input)
	if len(snap) == 0 {
		return
	}
	for _, c := range snap {
		lens = append(lens, c.Head)
		outputs = append(outputs, c.Outputs)
	}
	return

}

// WriteTo saves a program of finite State transducer.
func (t FST) WriteTo(w io.Writer) (n int64, err error) {
	var (
		pc   int
		code Instruction
		op   operation
		ch   byte
		v16  uint16
		v32  int32
	)
	dataLen := int64(len(t.Data))
	//fmt.Println("Data len:", dataLen)
	if err = binary.Write(w, binary.LittleEndian, dataLen); err != nil {
		return
	}
	n += int64(binary.Size(dataLen))
	//fmt.Println("Data len:", dataLen) //XXX
	for _, v := range t.Data {
		if err = binary.Write(w, binary.LittleEndian, v); err != nil {
			return
		}
		n += int64(binary.Size(v))
	}

	progLen := int64(len(t.Program))
	if err = binary.Write(w, binary.LittleEndian, progLen); err != nil {
		return
	}
	n += int64(binary.Size(progLen))

	//fmt.Println("Program len:", progLen) //XXX
	for pc = 0; pc < len(t.Program); pc++ {
		code = t.Program[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*uint16)(unsafe.Pointer(&code[2])))

		// write op and ch
		var tmp int
		tmp, err = w.Write(code[0:2])
		if err != nil {
			return
		}
		n += int64(tmp)
		//fmt.Printf("%3d %v\t%X %d\n", PC, op, ch, v16) //XXX
		switch operation(op) {
		case opAccept:
			fallthrough
		case opAcceptBreak:
			if ch == 0 {
				break
			}
			pc++
			code = t.Program[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0]))) //to addr
			if err = binary.Write(w, binary.LittleEndian, v32); err != nil {
				return
			}
			n += int64(binary.Size(v32))
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
			pc++
			code = t.Program[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0]))) //from addr
			if err = binary.Write(w, binary.LittleEndian, v32); err != nil {
				return
			}
			n += int64(binary.Size(v32))
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
		case opMatch:
			fallthrough
		case opMatchBreak:
			if err = binary.Write(w, binary.LittleEndian, v16); err != nil {
				return
			}
			n += int64(binary.Size(v16))
			if v16 != 0 {
				break
			}
			pc++
			code = t.Program[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			if err = binary.Write(w, binary.LittleEndian, v32); err != nil {
				return
			}
			n += int64(binary.Size(v32))
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
		case opOutput:
			fallthrough
		case opOutputBreak:
			if err = binary.Write(w, binary.LittleEndian, v16); err != nil {
				return
			}
			n += int64(binary.Size(v16))
			pc++
			code = t.Program[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			if err = binary.Write(w, binary.LittleEndian, v32); err != nil {
				return
			}
			n += int64(binary.Size(v32))
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX

			if v16 != 0 {
				break
			}
			pc++
			code = t.Program[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			if err = binary.Write(w, binary.LittleEndian, v32); err != nil {
				return
			}
			n += int64(binary.Size(v32))
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
		default:
			return n, fmt.Errorf("undefined operation error")
		}
	}
	return
}

// Read loads a program of finite State transducer.
func Read(r io.Reader) (t FST, e error) {
	var (
		code Instruction
		op   byte
		ch   byte
		v16  uint16
		v32  int32
		p    unsafe.Pointer
		//PC   int //XXX
	)

	rd := bufio.NewReader(r)

	var dataLen int64
	if e = binary.Read(rd, binary.LittleEndian, &dataLen); e != nil {
		return
	}
	//fmt.Println("Data len:", dataLen) //XXX
	t.Data = make([]int32, 0, dataLen)
	for i := 0; i < int(dataLen); i++ {
		if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
			return
		}
		t.Data = append(t.Data, v32)
	}

	var progLen int64
	if e = binary.Read(rd, binary.LittleEndian, &progLen); e != nil {
		return
	}
	//fmt.Println("Program len:", progLen) //XXX
	t.Program = make([]Instruction, 0, progLen)

	for e == nil {
		if op, e = rd.ReadByte(); e != nil {
			break
		}
		if ch, e = rd.ReadByte(); e != nil {
			break
		}
		switch operation(op) {
		case opAccept:
			fallthrough
		case opAcceptBreak:
			code[0], code[1], code[2], code[3] = op, ch, 0, 0
			t.Program = append(t.Program, code)
			//fmt.Printf("%3d %v\t%X %d\n", PC, operation(op), ch, 0) //XXX
			//PC++                                                    //XXX
			if ch == 0 {
				break
			}
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
			//PC++                                //XXX
			t.Program = append(t.Program, code)

			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
			//PC++                                //XXX
			t.Program = append(t.Program, code)
		case opMatch:
			fallthrough
		case opMatchBreak:
			code[0], code[1] = op, ch
			if e = binary.Read(rd, binary.LittleEndian, &v16); e != nil {
				break
			}
			p = unsafe.Pointer(&code[2])
			(*(*uint16)(p)) = uint16(v16)
			//fmt.Printf("%3d %v\t%X %d\n", PC, operation(op), ch, v16) //XXX
			//PC++                                                      //XXX
			t.Program = append(t.Program, code)

			if v16 != 0 {
				break
			}
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
			//PC++                                //XXX
			t.Program = append(t.Program, code)
		case opOutput:
			fallthrough
		case opOutputBreak:
			code[0], code[1] = op, ch
			if e = binary.Read(rd, binary.LittleEndian, &v16); e != nil {
				break
			}
			p = unsafe.Pointer(&code[2])
			(*(*uint16)(p)) = uint16(v16)
			//fmt.Printf("%3d %v\t%X %d\n", PC, operation(op), ch, v16) //XXX
			//PC++                                                      //XXX
			t.Program = append(t.Program, code)
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
			//PC++                                //XXX
			t.Program = append(t.Program, code)

			if v16 != 0 {
				break
			}
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", PC, v32) //XXX
			//PC++                                //XXX
			t.Program = append(t.Program, code)
		default:
			e = fmt.Errorf("invalid format: undefined operation error")
			break
		}
	}
	if e == io.EOF {
		e = nil
	}
	return
}
