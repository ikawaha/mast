package ss

import "fmt"

type instOp byte

const (
	valBits        = 5
	instBits       = 8 - valBits
	instShift      = valBits
	valMask   byte = 0xFF >> instBits
	instMask  byte = 0xFF - valMask

	instAccept      instOp = 0x01 << instShift
	instMatch              = 0x02 << instShift
	instBreak              = 0x03 << instShift
	instOutput             = 0x04 << instShift
	instOutputBreak        = 0x05 << instShift
)

var instOpName = [2 << instBits]string{
	"OP0",
	"ACC",
	"MAT",
	"BRK",
	"OUT",
	"OTB",
	"OP6",
	"OP7",
}

// String returns a operation name of a instruction.
func (op instOp) String() string {
	return instOpName[op>>instShift]
}

// FstVM represents a virtual machine of finite state transducers.
type FstVM struct {
	prog []byte
	data string
}

type configuration struct {
	pc   int
	inp  int
	tape int
}

func toInt(b []byte) int {
	var x int
	for i, size := 0, len(b); i < size; i++ {
		x <<= 8
		x += int(b[i])
	}
	return x
}

// String returns a string representation of a program.
func (vm FstVM) String() string {
	ret := ""
	for pc, end := 0, len(vm.prog); pc < end; {
		p := pc
		op := instOp(vm.prog[pc] & instMask)
		sz := int(vm.prog[pc] & valMask)
		pc++
		if op == instAccept {
			if sz == 0 {
				ret += fmt.Sprintf("%3d  %v\n", p, op)
			} else {
				s := toInt(vm.prog[pc : pc+sz])
				pc += sz
				sz = int(vm.prog[pc])
				pc++
				e := toInt(vm.prog[pc : pc+sz])
				var tails []string
				for i := s; i < e; i++ {
					h := i
					for vm.data[i] != 0 {
						i++
					}
					tails = append(tails, vm.data[h:i])
				}
				ret += fmt.Sprintf("%3d  %v %v\n", p, op, tails)
				for j := p + 1; j <= pc; j++ {
					ret += fmt.Sprintf("%3d [TIL addr=%d:%d]\n", j, s, e)
				}
				pc += sz
			}
			continue
		}
		inp := vm.prog[pc]
		pc++
		buf := vm.prog[pc : pc+sz]
		v := toInt(buf)
		ret += fmt.Sprintf("%3d  %v %X %d (sz:%d)\n", p, op, inp, v, sz)
		for j := p + 1; j < pc+sz; j++ {
			ret += fmt.Sprintf("%3d [%v %X %d]\n", j, op, inp, v)
		}
		pc += sz
		if op == instOutput || op == instOutputBreak {
			p := pc
			sz := int(vm.prog[pc])
			pc++
			buf := vm.prog[pc : pc+sz]
			v := toInt(buf)
			pc += sz
			for j := p; j < pc; j++ {
				ret += fmt.Sprintf("%3d [OUT addr=%d]\n", j, v)
			}

		}
	}
	return ret
}

func invert(b []byte) (inv []byte) {
	size := len(b)
	inv = make([]byte, len(b))
	for i := range b {
		inv[i] = b[size-1-i]
	}
	return
}

// Search runs a finite state transducer for a given input and returns outputs if accepted otherwise nil.
func (vm *FstVM) Search(input string) []string {
	tape, snap, acc := vm.run(input)
	if !acc || len(snap) == 0 {
		return nil
	}
	c := snap[len(snap)-1]
	pc := c.pc
	sz := int(vm.prog[pc] & valMask)
	pc++
	if sz == 0 {
		return []string{string(tape[0:c.tape])}
	}
	s := toInt(vm.prog[pc : pc+sz])
	pc += sz
	sz = int(vm.prog[pc])
	pc++
	e := toInt(vm.prog[pc : pc+sz])
	var outs []string
	for i := s; i < e; i++ {
		h := i
		for vm.data[i] != 0 {
			i++
		}
		t := append(tape[0:c.tape], vm.data[h:i]...)
		outs = append(outs, string(t))
	}
	pc += sz
	return outs
}

// PrefixSearch returns the longest commom prefix keyword and it's length in given input if detected otherwise -1, nil.
func (vm *FstVM) PrefixSearch(input string) (int, []string) {
	tape, snap, _ := vm.run(input)
	if len(snap) == 0 {
		return -1, nil
	}
	c := snap[len(snap)-1]
	pc := c.pc
	sz := int(vm.prog[pc] & valMask)
	pc++
	if sz == 0 {
		return c.inp, []string{string(tape[0:c.tape])}
	}
	s := toInt(vm.prog[pc : pc+sz])
	pc += sz
	sz = int(vm.prog[pc])
	pc++
	e := toInt(vm.prog[pc : pc+sz])
	var outs []string
	for i := s; i < e; i++ {
		h := i
		for vm.data[i] != 0 {
			i++
		}
		t := append(tape[0:c.tape], vm.data[h:i]...)
		outs = append(outs, string(t))
	}
	pc += sz
	return c.inp, outs
}

// CommonPrefixSearch finds keywords sharing common prefix in given input
// and returns it's lengths and outputs. Returns nil, nil if there does not common prefix keywords.
func (vm *FstVM) CommonPrefixSearch(input string) (lens []int, outputs [][]string) {
	tape, snap, _ := vm.run(input)
	if len(snap) == 0 {
		return
	}
	for _, c := range snap {
		lens = append(lens, c.inp)
		pc := c.pc
		sz := int(vm.prog[pc] & valMask)
		pc++
		if sz == 0 {
			outputs = append(outputs, []string{string(tape[0:c.tape])})
			continue
		}
		s := toInt(vm.prog[pc : pc+sz])
		pc += sz
		sz = int(vm.prog[pc])
		pc++
		e := toInt(vm.prog[pc : pc+sz])
		var outs []string
		for i := s; i < e; i++ {
			h := i
			for vm.data[i] != 0 {
				i++
			}
			t := make([]byte, 0, c.tape+(i-h))
			t = append(t, tape[0:c.tape]...)
			t = append(t, vm.data[h:i]...)
			outs = append(outs, string(t))
		}
		outputs = append(outputs, outs)
	}
	return
}

func (vm *FstVM) run(input string) (tape []byte, snap []configuration, accept bool) {
	var (
		pc int    // program counter
		op instOp // operation
		sz int    // size
		ch byte   // char
		hd int    // input head
		va int    // value
	)
	for pc < len(vm.prog) && hd < len(input) {
		op = instOp(vm.prog[pc] & instMask)
		sz = int(vm.prog[pc] & valMask)
		//fmt.Printf("pc:%v,op:%v,hd:%v,sz:%v\n", pc, op, hd, sz) //XXX
		switch op {
		case instMatch:
			fallthrough
		case instBreak:
			pc++
			ch = vm.prog[pc]
			pc++
			if ch != input[hd] {
				if op == instBreak {
					return
				}
				if sz > 0 {
					pc += sz
				}
				continue
			}
			if sz > 0 {
				va = toInt(vm.prog[pc : pc+sz])
				pc += sz + va
			} else {
				va = 0
			}
			hd++
			continue
		case instOutput:
			fallthrough
		case instOutputBreak:
			pc++
			ch = vm.prog[pc]
			pc++
			//fmt.Println("ch:", ch, "input[hd]", input[hd]) //XXX
			if ch != input[hd] {
				if op == instOutputBreak {
					return
				}
				if sz > 0 {
					pc += sz
				}
				s := int(vm.prog[pc])
				pc += s + 1
				continue
			}
			if sz > 0 {
				va = toInt(vm.prog[pc : pc+sz])
				pc += sz
			} else {
				va = 0
			}
			hd++
			s := int(vm.prog[pc])
			v := toInt(vm.prog[pc+1 : pc+1+s])
			//fmt.Println("pc:", pc, "s:", s, "v:", v) //XXX
			for p := v; ; {
				if p >= len(vm.data) || vm.data[p] == 0 {
					//fmt.Printf("out>>%s(%v, %v)\n", vm.data[v:p], v, p) //XXX
					tape = append(tape, vm.data[v:p]...)
					break
				}
				p++
			}
			//fmt.Println("pc:", pc, "s:", s, "va:", va)
			pc += s + va + 1
		case instAccept:
			snap = append(snap, configuration{pc, hd, len(tape)})
			pc++
			if sz > 0 {
				pc += sz
				sz = int(vm.prog[pc])
				pc += sz + 1
			}
			continue
		default:
			//fmt.Printf("unknown op:%v\n", op)
			return
		}
	}

	if pc >= len(vm.prog) || hd != len(input) {
		return
	}
	if op = instOp(vm.prog[pc] & instMask); op != instAccept {
		//fmt.Printf("[[FINAL]]pc:%d, op:%s, ch:[%X], sz:%d, v:%d\n", pc, op, ch, sz, va) //XXX
		return

	}
	accept = true
	snap = append(snap, configuration{pc, hd, len(tape)})
	return
}
