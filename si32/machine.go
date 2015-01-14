package si32

import (
	"fmt"
	"sort"
	"unsafe"
)

type operation byte

const (
	opAccept      operation = 1
	opMatch       operation = 2
	opBreak       operation = 3
	opOutput      operation = 4
	opOutputBreak operation = 5
)

func (o operation) String() string {
	opName := []string{"OP0", "ACC", "MTC", "BRK", "OUT", "OUB", "OP6", "OP7"}
	if int(o) >= len(opName) {
		return fmt.Sprintf("NA%d", o)
	}
	return opName[o]
}

type instruction [4]byte

// FST represents a finite state transducer (virtual machine).
type FST struct {
	prog []instruction
	data []int32
}

// Configuration represents a FST (virtual machine) configuration.
type configuration struct {
	pc  int     // program counter
	hd  int     // input head
	out []int32 // outputs
}

const maxUint16 = 1<<16 - 1

type int32Slice []int32

func (p int32Slice) Len() int           { return len(p) }
func (p int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type byteSlice []byte

func (p byteSlice) Len() int           { return len(p) }
func (p byteSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p byteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func invert(prog []instruction) []instruction {
	size := len(prog)
	inv := make([]instruction, size)
	for i := range prog {
		inv[i] = prog[size-1-i]
	}
	return inv
}

func buildFST(m mast) (t FST, err error) {
	var (
		prog []instruction
		data []int32
		code instruction // temporary
	)
	var edges []byte
	addrMap := make(map[int]int)
	for _, s := range m.states {
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
			out := s.Output[ch]
			addr, ok := addrMap[next.ID]
			if !ok && !next.IsFinal {
				err = fmt.Errorf("next addr is undefined: state(%v), input(%X)", s.ID, ch)
				return
			}
			jump := len(prog) - addr + 1
			var op operation
			if out != 0 {
				if i == 0 {
					op = opOutputBreak
				} else {
					op = opOutput
				}
			} else if i == 0 {
				op = opBreak
			} else {
				op = opMatch
			}

			if out != 0 {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(out)
				prog = append(prog, code)
			}
			if jump > maxUint16 {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(jump)
				prog = append(prog, code)
				jump = 0
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
			code[0] = byte(opAccept)
			code[1], code[2], code[3] = 0, 0, 0 // clear
			if len(s.Tail) > 0 {
				code[1] = 1
			}

			prog = append(prog, code)
		}
		addrMap[s.ID] = len(prog)
	}
	t = FST{prog: invert(prog), data: data}
	return
}

func (t FST) String() string {
	var (
		pc   int
		code instruction
		op   operation
		ch   byte
		v16  int16
		v32  int32
	)
	ret := ""
	for pc = 0; pc < len(t.prog); pc++ {
		code = t.prog[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*int16)(unsafe.Pointer(&code[2])))
		switch operation(op) {
		case opAccept:
			//fmt.Printf("%3d %v\t%X %d\n", pc, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%X %d\n", pc, op, ch, v16)
			if ch == 0 {
				break
			}
			pc++
			code = t.prog[pc]
			to := (*(*int32)(unsafe.Pointer(&code[0])))
			ret += fmt.Sprintf("%3d [%d]\n", pc, to)
			pc++
			code = t.prog[pc]
			from := (*(*int32)(unsafe.Pointer(&code[0])))
			ret += fmt.Sprintf("%3d [%d] %v\n", pc, from, t.data[from:to]) //FIXME
		case opMatch:
			fallthrough
		case opBreak:
			//fmt.Printf("%3d %v\t%02X %d\n", pc, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%02X %d\n", pc, op, ch, v16)
			if v16 != 0 {
				break
			}
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			//fmt.Printf("%3d [%d]\n", pc, v32) //XXX
			ret += fmt.Sprintf("%3d [%d]\n", pc, v32)
		case opOutput:
			fallthrough
		case opOutputBreak:
			//fmt.Printf("%3d %v\t%02X %d\n", pc, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%02X %d\n", pc, op, ch, v16)
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			//fmt.Printf("%3d [%d]\n", pc, v32) //XXX
			ret += fmt.Sprintf("%3d [%d]\n", pc, v32)
			//fmt.Println("pc: ", pc) //XXX
			if v16 != 0 {
				break
			}
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			//fmt.Printf("%3d [%d]\n", pc, v32) //XXX
			ret += fmt.Sprintf("%3d [%d]\n", pc, v32)
		default:
			//fmt.Printf("%3d UNDEF %v\n", pc, code)
			ret += fmt.Sprintf("%3d UNDEF %v\n", pc, code)
		}
	}
	return ret
}

func (t *FST) run(input string) (snap []configuration, accept bool) {
	var (
		pc  int       // program counter
		op  operation // operation
		ch  byte      // char
		v16 int16     // 16bit register
		v32 int32     // 32bit register
		hd  int       // input head
		out int32     // output

		code instruction // tmp
	)
	for pc < len(t.prog) && hd <= len(input) {
		code = t.prog[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*int16)(unsafe.Pointer(&code[2])))
		//fmt.Printf("pc:%v,op:%v,hd:%v,v16:%v,out:%v\n", pc, op, hd, v16, out) //XXX
		switch op {
		case opMatch:
			fallthrough
		case opBreak:
			if hd == len(input) {
				goto L_END
			}
			if ch != input[hd] {
				if op == opBreak {
					return
				}
				pc++
				continue
			}
			if v16 > 0 {
				pc += int(v16)
			} else {
				pc++
				code = t.prog[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
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
				pc++
				pc++
				continue
			}
			pc++
			code = t.prog[pc]
			out = (*(*int32)(unsafe.Pointer(&code[0])))
			if v16 > 0 {
				pc += int(v16)
			} else {
				pc++
				code = t.prog[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				pc += int(v32)
			}
			hd++
			continue
		case opAccept:
			c := configuration{pc: pc, hd: hd}
			pc++
			if ch == 0 {
				c.out = []int32{out}
			} else {
				code = t.prog[pc]
				to := (*(*int32)(unsafe.Pointer(&code[0])))
				pc++
				code = t.prog[pc]
				from := (*(*int32)(unsafe.Pointer(&code[0])))
				c.out = t.data[from:to]
				pc++
			}
			snap = append(snap, c)
			if hd == len(input) {
				goto L_END
			}
			continue
		default:
			//fmt.Printf("unknown op:%v\n", op)
			return
		}
	}
L_END:
	if hd != len(input) {
		return
	}
	if op != opAccept {
		//fmt.Printf("[[FINAL]]pc:%d, op:%s, ch:[%X], sz:%d, v:%d\n", pc, op, ch, sz, va) //XXX
		return

	}
	accept = true
	return
}

// Search runs a finite state transducer for a given input and returns outputs if accepted otherwise nil.
func (t FST) Search(input string) []int32 {
	snap, acc := t.run(input)
	if !acc || len(snap) == 0 {
		return nil
	}
	c := snap[len(snap)-1]
	return c.out
}

// PrefixSearch returns the longest commom prefix keyword and it's length in given input
// if detected otherwise -1, nil.
func (t FST) PrefixSearch(input string) (length int, output []int32) {
	snap, _ := t.run(input)
	if len(snap) == 0 {
		return -1, nil
	}
	c := snap[len(snap)-1]
	return c.hd, c.out

}

// CommonPrefixSearch finds keywords sharing common prefix in given input
// and returns it's lengths and outputs. Returns nil, nil if there does not common prefix keywords.
func (t FST) CommonPrefixSearch(input string) (lens []int, outputs [][]int32) {
	snap, _ := t.run(input)
	if len(snap) == 0 {
		return
	}
	for _, c := range snap {
		lens = append(lens, c.hd)
		outputs = append(outputs, c.out)
	}
	return

}
