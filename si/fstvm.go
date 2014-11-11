package si

import "fmt"

type instOp byte

const (
	valBits        = 6
	instBits       = 8 - valBits
	instShift      = valBits
	valMask   byte = 0xFF >> instBits
	instMask  byte = 0xFF - valMask

	instAccept instOp = 0x01 << instShift
	instMatch         = 0x02 << instShift
	instBreak         = 0x03 << instShift
)

var instOpName = [2 << instBits]string{
	"OP0",
	"ACC",
	"MAT",
	"BRK",
}

// String returns a operation name of a instruction.
func (op instOp) String() string {
	return instOpName[op>>instShift]
}

// FstVM represents a virtual machine of finite state transducers.
type FstVM struct {
	prog []byte
	data []int
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
				ret += fmt.Sprintf("%3d  %v %v\n", p, op, vm.data[s:e])
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
func (vm *FstVM) Search(input string) []int {
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
					return nil
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
		case instAccept:
			pc++
			if sz > 0 {
				pc += sz
				sz = int(vm.prog[pc])
				pc += sz + 1
			}
			continue
		default:
			//fmt.Printf("unknown op:%v\n", op)
			return nil
		}
	}

	if pc >= len(vm.prog) || hd != len(input) {
		return nil
	}
	if op = instOp(vm.prog[pc] & instMask); op != instAccept {
		//fmt.Printf("[[FINAL]]pc:%d, op:%s, ch:[%X], sz:%d, v:%d\n", pc, op, ch, sz, va) //XXX
		return nil

	}
	sz = int(vm.prog[pc] & valMask)
	pc++
	s := toInt(vm.prog[pc : pc+sz])
	pc += sz
	sz = int(vm.prog[pc])
	pc++
	e := toInt(vm.prog[pc : pc+sz])
	pc += sz
	return vm.data[s:e]
}
