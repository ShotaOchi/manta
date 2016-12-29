package manta

import (
	"strconv"
	"strings"
)

var huffTree = newHuffmanTree()

type fieldPath struct {
	path []int
	last int
	done bool
}

type fieldPathOp struct {
	name   string
	weight int
	fn     func(r *reader, fp *fieldPath)
}

var fieldPathTable = []fieldPathOp{
	fieldPathOp{"PlusOne", 36271, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 1
	}},
	fieldPathOp{"PlusTwo", 10334, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 2
	}},
	fieldPathOp{"PlusThree", 1375, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 3
	}},
	fieldPathOp{"PlusFour", 646, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 4
	}},
	fieldPathOp{"PlusN", 4128, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += r.readUBitVarFieldPath() + 5
	}},
	fieldPathOp{"PushOneLeftDeltaZeroRightZero", 35, func(r *reader, fp *fieldPath) {
		fp.last++
		fp.path[fp.last] = 0
	}},
	fieldPathOp{"PushOneLeftDeltaZeroRightNonZero", 3, func(r *reader, fp *fieldPath) {
		fp.last++
		fp.path[fp.last] = r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushOneLeftDeltaOneRightZero", 521, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 1
		fp.last++
		fp.path[fp.last] = 0
	}},
	fieldPathOp{"PushOneLeftDeltaOneRightNonZero", 2942, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 1
		fp.last++
		fp.path[fp.last] = r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushOneLeftDeltaNRightZero", 560, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] = 0
	}},
	fieldPathOp{"PushOneLeftDeltaNRightNonZero", 471, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += r.readUBitVarFieldPath() + 2
		fp.last++
		fp.path[fp.last] = r.readUBitVarFieldPath() + 1
	}},
	fieldPathOp{"PushOneLeftDeltaNRightNonZeroPack6Bits", 10530, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += int(r.readBits(3)) + 2
		fp.last++
		fp.path[fp.last] = int(r.readBits(3)) + 1
	}},
	fieldPathOp{"PushOneLeftDeltaNRightNonZeroPack8Bits", 251, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += int(r.readBits(4)) + 2
		fp.last++
		fp.path[fp.last] = int(r.readBits(4)) + 1
	}},
	fieldPathOp{"PushTwoLeftDeltaZero", 0, func(r *reader, fp *fieldPath) {
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushTwoPack5LeftDeltaZero", 0, func(r *reader, fp *fieldPath) {
		fp.last++
		fp.path[fp.last] = int(r.readBits(5))
		fp.last++
		fp.path[fp.last] = int(r.readBits(5))
	}},
	fieldPathOp{"PushThreeLeftDeltaZero", 0, func(r *reader, fp *fieldPath) {
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushThreePack5LeftDeltaZero", 0, func(r *reader, fp *fieldPath) {
		fp.last++
		fp.path[fp.last] = int(r.readBits(5))
		fp.last++
		fp.path[fp.last] = int(r.readBits(5))
		fp.last++
		fp.path[fp.last] = int(r.readBits(5))
	}},
	fieldPathOp{"PushTwoLeftDeltaOne", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 1
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushTwoPack5LeftDeltaOne", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 1
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
	}},
	fieldPathOp{"PushThreeLeftDeltaOne", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 1
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushThreePack5LeftDeltaOne", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += 1
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
	}},
	fieldPathOp{"PushTwoLeftDeltaN", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += int(r.readUBitVar()) + 2
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushTwoPack5LeftDeltaN", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += int(r.readUBitVar()) + 2
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
	}},
	fieldPathOp{"PushThreeLeftDeltaN", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += int(r.readUBitVar()) + 2
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
		fp.last++
		fp.path[fp.last] += r.readUBitVarFieldPath()
	}},
	fieldPathOp{"PushThreePack5LeftDeltaN", 0, func(r *reader, fp *fieldPath) {
		fp.path[fp.last] += int(r.readUBitVar()) + 2
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
		fp.last++
		fp.path[fp.last] += int(r.readBits(5))
	}},
	fieldPathOp{"PushN", 0, func(r *reader, fp *fieldPath) {
		n := int(r.readUBitVar())
		fp.path[fp.last] += int(r.readUBitVar())
		for i := 0; i < n; i++ {
			fp.last++
			fp.path[fp.last] += r.readUBitVarFieldPath()
		}

	}},
	fieldPathOp{"PushNAndNonTopological", 310, func(r *reader, fp *fieldPath) {
		for i := 0; i <= fp.last; i++ {
			if r.readBoolean() {
				fp.path[i] += int(r.readVarInt32()) + 1
			}
		}
		count := int(r.readUBitVar())
		for i := 0; i < count; i++ {
			fp.last++
			fp.path[fp.last] = r.readUBitVarFieldPath()
		}
	}},
	fieldPathOp{"PopOnePlusOne", 2, func(r *reader, fp *fieldPath) {
		fp.last--
		fp.path[fp.last] += 1
	}},
	fieldPathOp{"PopOnePlusN", 0, func(r *reader, fp *fieldPath) {
		fp.last--
		fp.path[fp.last] += r.readUBitVarFieldPath() + 1
	}},
	fieldPathOp{"PopAllButOnePlusOne", 1837, func(r *reader, fp *fieldPath) {
		fp.last = 0
		fp.path[0] += 1

	}},
	fieldPathOp{"PopAllButOnePlusN", 149, func(r *reader, fp *fieldPath) {
		fp.last = 0
		fp.path[0] += r.readUBitVarFieldPath() + 1
	}},
	fieldPathOp{"PopAllButOnePlusNPack3Bits", 300, func(r *reader, fp *fieldPath) {
		fp.last = 0
		fp.path[0] += int(r.readBits(3)) + 1
	}},
	fieldPathOp{"PopAllButOnePlusNPack6Bits", 634, func(r *reader, fp *fieldPath) {
		fp.last = 0
		fp.path[0] += int(r.readBits(6)) + 1
	}},
	fieldPathOp{"PopNPlusOne", 0, func(r *reader, fp *fieldPath) {
		fp.last -= r.readUBitVarFieldPath()
		fp.path[fp.last] += 1
	}},
	fieldPathOp{"PopNPlusN", 0, func(r *reader, fp *fieldPath) {
		fp.last -= r.readUBitVarFieldPath()
		fp.path[fp.last] += int(r.readVarInt32())
	}},
	fieldPathOp{"PopNAndNonTopographical", 1, func(r *reader, fp *fieldPath) {
		fp.last -= r.readUBitVarFieldPath()
		for i := 0; i <= fp.last; i++ {
			if r.readBoolean() {
				fp.path[i] += int(r.readVarInt32())
			}
		}
	}},
	fieldPathOp{"NonTopoComplex", 76, func(r *reader, fp *fieldPath) {
		for i := 0; i <= fp.last; i++ {
			if r.readBoolean() {
				fp.path[i] += int(r.readVarInt32())
			}
		}
	}},
	fieldPathOp{"NonTopoPenultimatePlusOne", 271, func(r *reader, fp *fieldPath) {
		fp.path[fp.last-1] += 1
	}},
	fieldPathOp{"NonTopoComplexPack4Bits", 99, func(r *reader, fp *fieldPath) {
		for i := 0; i <= fp.last; i++ {
			if r.readBoolean() {
				fp.path[i] += int(r.readBits(4)) - 7
			}
		}
	}},
	fieldPathOp{"FieldPathEncodeFinish", 25474, func(r *reader, fp *fieldPath) {
		fp.done = true
	}},
}

func (fp *fieldPath) copy() *fieldPath {
	x := &fieldPath{
		path: make([]int, fp.last+1),
		last: fp.last,
	}
	copy(x.path, fp.path[:fp.last+1])
	x.done = fp.done
	return x
}

func (fp *fieldPath) String() string {
	ss := make([]string, len(fp.path))
	for i, n := range fp.path {
		ss[i] = strconv.Itoa(n)
	}
	return strings.Join(ss, "/")
}

func readFieldPaths(r *reader) []*fieldPath {
	fp := &fieldPath{
		path: []int{-1, 0, 0, 0, 0, 0},
		last: 0,
		done: false,
	}

	node, next := huffTree, huffTree

	paths := []*fieldPath{}

	for !fp.done {
		if r.readBits(1) == 1 {
			next = node.Right()
		} else {
			next = node.Left()
		}

		if next.IsLeaf() {
			node = huffTree
			op := fieldPathTable[next.Value()]

			if waldnew {
				_printf("NEW FP BEFORE: %s (%s) pos=%s", fp.copy().String(), op.name, r.position())
			}

			op.fn(r, fp)

			if waldnew {
				_printf("NEW FP AFTER: %s (%s) pos=%s", fp.copy().String(), op.name, r.position())
			}

			if !fp.done {
				paths = append(paths, fp.copy())
			}
		} else {
			node = next
		}
	}

	return paths
}

func newHuffmanTree() huffmanTree {
	freqs := make([]int, len(fieldpathLookup))
	for i, op := range fieldPathTable {
		freqs[i] = op.weight
	}
	return buildHuffmanTree(freqs)
}