package si32

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Operation represents the instruction code.
type Operation byte

const (
	// Accept is an operation code that accepts if character matches, otherwise it increments the program counter.
	Accept Operation = 1
	// AcceptBreak is an operation code that jumps if character matches, otherwise it breaks the program.
	AcceptBreak Operation = 2
	// Match is an operation code that jumps if character matches, otherwise it increments the program counter.
	Match Operation = 3
	// MatchBreak is an operation code that jumps if character matches, otherwise it breaks the program.
	MatchBreak Operation = 4
	// Output is an operation code that outputs.
	Output Operation = 5
	// OutputBreak is an operation code that outputs and breaks program.
	OutputBreak Operation = 6
)

// String returns the name of the operation.
func (o Operation) String() string {
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

// Program represents program that executes an FST.
type Program []Instruction

// Instruction represents an instruction of program.
type Instruction uint32

// FST represents a finite state transducer.
type FST struct {
	Program []Instruction
	Data    []int32
}

// Configuration represents a configuration of an FST.
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

type int32Slice []int32

func (p int32Slice) Len() int           { return len(p) }
func (p int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Reverse rearrange instructions in reverse order.
func (p Program) Reverse() {
	size := len(p)
	for i := 0; i < size/2; i++ {
		p[i], p[size-1-i] = p[size-1-i], p[i]
	}
}

// BuildFST generates virtual machine code of an FST from a minimal acyclic subsequential transducer
func (m MAST) BuildFST() (*FST, error) {
	var (
		prog  Program
		data  []int32
		edges []byte
	)
	addrMap := map[int]int{}
	for _, s := range m.States {
		edges = edges[:0]
		for ch := range s.Trans {
			edges = append(edges, ch)
		}
		if len(edges) > 0 {
			sort.Sort(byteSlice(edges))
		}
		for i, size := 0, len(edges); i < size; i++ {
			ch := edges[size-1-i]
			next := s.Trans[ch]
			addr, ok := addrMap[next.ID]
			if !ok && !next.IsFinal {
				return nil, fmt.Errorf("next addr is undefined: State(%v), input(%X)", s.ID, ch)
			}

			var op Operation
			out, ok := s.Output[ch]
			if !ok {
				if i == 0 {
					op = MatchBreak
				} else {
					op = Match
				}
			} else if i == 0 {
				op = OutputBreak
			} else {
				op = Output
			}

			jump := len(prog) - addr + 1
			if jump > maxUint16 {
				prog = append(prog, Instruction(jump))
				jump = 0
			}
			if ok {
				prog = append(prog, Instruction(out))
			}
			prog = append(prog, Instruction((int(op)<<24)+(int(ch)<<16)+jump))
		}
		if s.IsFinal {
			if len(s.Tail) > 0 {
				prog = append(prog, Instruction(len(data))) // from
				tmp := make(int32Slice, 0, len(s.Tail))
				for t := range s.Tail {
					tmp = append(tmp, t)
				}
				sort.Sort(tmp)
				data = append(data, tmp...)
				prog = append(prog, Instruction(len(data))) // to
			}
			var inst Instruction
			if len(s.Trans) == 0 {
				inst = Instruction(int(AcceptBreak) << 24)
			} else {
				inst = Instruction(int(Accept) << 24)
			}
			if len(s.Tail) > 0 {
				inst += Instruction(1 << 16)
			}
			prog = append(prog, inst)
		}
		addrMap[s.ID] = len(prog)
	}
	prog.Reverse()
	return &FST{Program: prog, Data: data}, nil
}

// String returns virtual machine code of the FST.
func (t FST) String() string {
	var (
		pc   int
		code Instruction
		op   Operation
		ch   byte
		v16  uint16
		v32  int32
	)
	var b strings.Builder
	for pc = 0; pc < len(t.Program); pc++ {
		code = t.Program[pc]
		op = Operation((code & 0xFF000000) >> 24)
		ch = byte((code & 0x00FF0000) >> 16)
		v16 = uint16(code & 0x0000FFFF)
		switch Operation(op) {
		case Accept:
			fallthrough
		case AcceptBreak:
			fmt.Fprintf(&b, "%3d %v\t%d %d\n", pc, op, ch, v16)
			if ch == 0 {
				break
			}
			pc++
			code = t.Program[pc]
			to := code
			fmt.Fprintf(&b, "%3d [%d]\n", pc, to)
			pc++
			code = t.Program[pc]
			from := code
			fmt.Fprintf(&b, "%3d [%d] %v\n", pc, from, t.Data[from:to])
		case Match:
			fallthrough
		case MatchBreak:
			fmt.Fprintf(&b, "%3d %v\t%02X(%c) %d\n", pc, op, ch, ch, v16)
			if v16 == 0 {
				pc++
				code = t.Program[pc]
				v32 = int32(code)
				fmt.Fprintf(&b, "%3d jmp[%d]\n", pc, v32)
			}
		case Output:
			fallthrough
		case OutputBreak:
			fmt.Fprintf(&b, "%3d %v\t%02X(%c) %d\n", pc, op, ch, ch, v16)
			if v16 == 0 {
				pc++
				code = t.Program[pc]
				v32 = int32(code)
				fmt.Fprintf(&b, "%3d jmp[%d]\n", pc, v32)
			}
			pc++
			code = t.Program[pc]
			v32 = int32(code)
			fmt.Fprintf(&b, "%3d [%d]\n", pc, v32)
		default:
			fmt.Fprintf(&b, "%3d UNDEF %v\n", pc, code)
		}
	}
	return b.String()
}

// Run runs virtual machine code of the FST.
func (t *FST) Run(input string) (snap []Configuration, accept bool) {
	var (
		pc   int       // program counter
		op   Operation // operation
		ch   byte      // char
		v16  uint16    // 16bit register
		v32  int32     // 32bit register
		head int       // input head
		out  int32     // output

		inst Instruction // tmp instruction
	)
	for pc < len(t.Program) && head <= len(input) {
		inst = t.Program[pc]
		op = Operation((inst & 0xFF000000) >> 24)
		ch = byte((inst & 0x00FF0000) >> 16)
		v16 = uint16(inst & 0x0000FFFF)
		//fmt.Printf("PC:%v,op:%v,Head:%v,v16:%v,Outputs:%v\n", PC, op, Head, v16, Outputs) //XXX
		switch op {
		case Match, MatchBreak:
			if head == len(input) {
				goto L_END
			}
			if ch != input[head] {
				if op == MatchBreak {
					return snap, false
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
				inst = t.Program[pc]
				v32 = int32(inst)
				//fmt.Printf("ex jump:%d\n", v32) //XXX
				pc += int(v32)
			}
			head++
			continue
		case Output, OutputBreak:
			if head == len(input) {
				goto L_END
			}
			if ch != input[head] {
				if op == OutputBreak {
					return snap, false
				}
				if v16 == 0 {
					pc++
				}
				pc++
				pc++
				continue
			}
			pc++
			inst = t.Program[pc]
			out = int32(inst)
			if v16 > 0 {
				pc += int(v16)
			} else {
				pc++
				inst = t.Program[pc]
				v32 = int32(inst)
				//fmt.Printf("ex jump:%d\n", v32) //XXX
				pc += int(v32)
			}
			head++
			continue
		case Accept, AcceptBreak:
			c := Configuration{PC: pc, Head: head}
			pc++
			if ch == 0 {
				c.Outputs = []int32{out}
			} else {
				inst = t.Program[pc]
				to := inst
				pc++
				inst = t.Program[pc]
				from := inst
				c.Outputs = t.Data[from:to]
				pc++
			}
			//fmt.Printf("conf: %+v\n", c) //XXX
			snap = append(snap, c)
			if head == len(input) {
				goto L_END
			}
			if op == AcceptBreak {
				goto L_END
			}
			continue
		default:
			//fmt.Printf("unknown op:%v\n", op) //XXX
			return snap, false
		}
	}
L_END:
	//fmt.Printf("[[L_END]]PC:%d, op:%s, ch:[%X]\n", PC, op, ch) //XXX
	if head != len(input) {
		return snap, false
	}
	if op != Accept && op != AcceptBreak {
		//fmt.Printf("[[NOT ACCEPT]]PC:%d, op:%s, ch:[%X]\n", PC, op, ch) //XXX
		return snap, false

	}
	//fmt.Printf("[[ACCEPT]]PC:%d, op:%s, ch:[%X]\n", PC, op, ch) //XXX
	return snap, true
}

// Search runs the FST for the given input and it returns outputs if accepted otherwise nil.
func (t FST) Search(input string) []int32 {
	snap, acc := t.Run(input)
	if !acc || len(snap) == 0 {
		return nil
	}
	c := snap[len(snap)-1]
	return c.Outputs
}

// PrefixSearch returns the longest common prefix keyword and its length.
// If there is no common prefix keyword, it returns (-1, nil).
func (t FST) PrefixSearch(input string) (length int, output []int32) {
	snap, _ := t.Run(input)
	if len(snap) == 0 {
		return -1, nil
	}
	c := snap[len(snap)-1]
	return c.Head, c.Outputs
}

// CommonPrefixSearch finds keywords sharing common prefix and it returns its lengths and outputs.
// If there are no common prefix keywords, it returns (nil, nil).
func (t FST) CommonPrefixSearch(input string) (lens []int, outputs [][]int32) {
	snap, _ := t.Run(input)
	if len(snap) == 0 {
		return lens, outputs
	}
	for _, c := range snap {
		lens = append(lens, c.Head)
		outputs = append(outputs, c.Outputs)
	}
	return lens, outputs

}

// WriteTo saves program of the FST.
func (t FST) WriteTo(w io.Writer) (n int64, err error) {
	dataLen := int64(len(t.Data))
	if err = binary.Write(w, binary.LittleEndian, dataLen); err != nil {
		return n, err
	}
	n += int64(binary.Size(dataLen))
	for _, v := range t.Data {
		if err = binary.Write(w, binary.LittleEndian, v); err != nil {
			return
		}
		n += int64(binary.Size(v))
	}

	progLen := int64(len(t.Program))
	if err = binary.Write(w, binary.LittleEndian, progLen); err != nil {
		return n, err
	}
	n += int64(binary.Size(progLen))
	for _, code := range t.Program {
		if err := binary.Write(w, binary.LittleEndian, code); err != nil {
			return n, err
		}
		n += int64(binary.Size(code))
	}
	return n, nil
}

// Read loads program of the FST.
func Read(r io.Reader) (*FST, error) {
	rd := bufio.NewReader(r)
	var dataLen int64
	if err := binary.Read(rd, binary.LittleEndian, &dataLen); err != nil {
		return nil, err
	}
	//fmt.Println("Data len:", dataLen)
	data := make([]int32, 0, dataLen)
	for i := 0; i < int(dataLen); i++ {
		var v32 int32
		if err := binary.Read(rd, binary.LittleEndian, &v32); err != nil {
			return nil, err
		}
		data = append(data, v32)
	}

	var progLen int64
	if err := binary.Read(rd, binary.LittleEndian, &progLen); err != nil {
		return nil, err
	}
	//fmt.Println("Program len:", progLen) //XXX
	program := make([]Instruction, 0, progLen)
	for i := 0; i < int(progLen); i++ {
		var v32 Instruction
		if err := binary.Read(rd, binary.LittleEndian, &v32); err != nil {
			return nil, err
		}
		program = append(program, v32)
	}
	return &FST{
		Program: program,
		Data:    data,
	}, nil
}
