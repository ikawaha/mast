package si

import (
	"encoding/binary"
	"fmt"
	"io"
)

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

type configuration struct {
	pc  int
	inp int
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
	for pc := 0; pc < len(vm.prog); {
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

func (vm *FstVM) run(input string) (snap []configuration, accept bool) {
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
		case instAccept:
			snap = append(snap, configuration{pc, hd})
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
	snap = append(snap, configuration{pc, hd})
	return
}

// Search runs a finite state transducer for a given input and returns outputs if accepted otherwise nil.
func (vm *FstVM) Search(input string) []int {
	snap, acc := vm.run(input)
	if !acc || len(snap) == 0 {
		return nil
	}
	c := snap[len(snap)-1]
	pc := c.pc
	sz := int(vm.prog[pc] & valMask)
	pc++
	s := toInt(vm.prog[pc : pc+sz])
	pc += sz
	sz = int(vm.prog[pc])
	pc++
	e := toInt(vm.prog[pc : pc+sz])
	pc += sz
	return vm.data[s:e]
}

// PrefixSearch returns the longest commom prefix keyword and it's length in given input if detected otherwise -1, nil.
func (vm *FstVM) PrefixSearch(input string) (int, []int) {
	snap, _ := vm.run(input)
	if len(snap) == 0 {
		return -1, nil
	}
	c := snap[len(snap)-1]
	pc := c.pc
	sz := int(vm.prog[pc] & valMask)
	pc++
	s := toInt(vm.prog[pc : pc+sz])
	pc += sz
	sz = int(vm.prog[pc])
	pc++
	e := toInt(vm.prog[pc : pc+sz])
	pc += sz
	return c.inp, vm.data[s:e]

}

// CommonPrefixSearch finds keywords sharing common prefix in given input
// and returns it's lengths and outputs. Returns nil, nil if there does not common prefix keywords.
func (vm *FstVM) CommonPrefixSearch(input string) (lens []int, outputs [][]int) {
	snap, _ := vm.run(input)
	if len(snap) == 0 {
		return
	}
	for _, c := range snap {
		pc := c.pc
		sz := int(vm.prog[pc] & valMask)
		pc++
		s := toInt(vm.prog[pc : pc+sz])
		pc += sz
		sz = int(vm.prog[pc])
		pc++
		e := toInt(vm.prog[pc : pc+sz])
		pc += sz
		lens = append(lens, c.inp)
		outputs = append(outputs, vm.data[s:e])
	}
	return

}

// Save FstVM
func (vm FstVM) Save(w io.Writer) (err error) {
	if err = binary.Write(w, binary.LittleEndian, int64(len(vm.prog))); err != nil {
		return
	}
	if _, err = w.Write(vm.prog); err != nil { //TODO compress
		return
	}
	if err = binary.Write(w, binary.LittleEndian, int64(len(vm.data))); err != nil {
		return
	}
	for i := 0; i < len(vm.data); i++ {
		if err = binary.Write(w, binary.LittleEndian, int64(vm.data[i])); err != nil {
			return
		}
	}
	return
}

// Load FstVM
func (vm *FstVM) Load(r io.Reader) (err error) {
	var n int64
	if err = binary.Read(r, binary.LittleEndian, &n); err != nil {
		return
	}
	prog := make([]byte, n, n)
	if _, err = r.Read(prog); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &n); err != nil {
		return
	}
	data := make([]int, n, n)
	for i := 0; i < int(n); i++ {
		var v int64
		if err = binary.Read(r, binary.LittleEndian, &v); err != nil {
			return
		}
		data[i] = int(v)
	}
	vm.prog = prog
	vm.data = data
	return
}
