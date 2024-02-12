package comb

// The Memo is pre-populated with the first 256 values.
var Memo = []int{6, 10, 15, 20, 21, 35, 28, 56, 70, 36, 84, 126, 45, 120, 210,
	252, 55, 165, 330, 462, 66, 220, 495, 792, 924, 78, 286, 715, 1287, 1716, 91,
	364, 1001, 2002, 3003, 3432, 105, 455, 1365, 3003, 5005, 6435, 120, 560, 1820,
	4368, 8008, 11440, 12870, 136, 680, 2380, 6188, 12376, 19448, 24310, 153, 816,
	3060, 8568, 18564, 31824, 43758, 48620, 171, 969, 3876, 11628, 27132, 50388,
	75582, 92378, 190, 1140, 4845, 15504, 38760, 77520, 125970, 167960, 184756, 210,
	1330, 5985, 20349, 54264, 116280, 203490, 293930, 352716, 231, 1540, 7315,
	26334, 74613, 170544, 319770, 497420, 646646, 705432, 253, 1771, 8855, 33649,
	100947, 245157, 490314, 817190, 1144066, 1352078, 276, 2024, 10626, 42504,
	134596, 346104, 735471, 1307504, 1961256, 2496144, 2704156, 300, 2300, 12650,
	53130, 177100, 480700, 1081575, 2042975, 3268760, 4457400, 5200300, 325, 2600,
	14950, 65780, 230230, 657800, 1562275, 3124550, 5311735, 7726160, 9657700,
	10400600, 351, 2925, 17550, 80730, 296010, 888030, 2220075, 4686825, 8436285,
	13037895, 17383860, 20058300, 378, 3276, 20475, 98280, 376740, 1184040, 3108105,
	6906900, 13123110, 21474180, 30421755, 37442160, 40116600, 406, 3654, 23751,
	118755, 475020, 1560780, 4292145, 10015005, 20030010, 34597290, 51895935,
	67863915, 77558760, 435, 4060, 27405, 142506, 593775, 2035800, 5852925,
	14307150, 30045015, 54627300, 86493225, 119759850, 145422675, 155117520, 465,
	4495, 31465, 169911, 736281, 2629575, 7888725, 20160075, 44352165, 84672315,
	141120525, 206253075, 265182525, 300540195, 496, 4960, 35960, 201376, 906192,
	3365856, 10518300, 28048800, 64512240, 129024480, 225792840, 347373600,
	471435600, 565722720, 601080390, 528, 5456, 40920, 237336, 1107568, 4272048,
	13884156, 38567100, 92561040, 193536720, 354817320, 573166440, 818809200,
	1037158320, 1166803110, 561, 5984, 46376, 278256, 1344904, 5379616, 18156204,
	52451256, 131128140, 286097760, 548354040, 927983760, 1391975640, 1855967520,
	2203961430, 2333606220}

// If we look at the binomial as Pascals triangle, there are two features we
// can take advantage of to compress the memo. The first is symmetry. Which is
// what
//   if i > n/2
// is checking for.
//
// The second thing is that the leftmost diagonal is always one and the second
// diagonal increments. Which could be written
// if i == 0{
//   return 1
// } else if i == 1 {
//   return n
// }
// but it's faster to combine them into
// if i < 2 {
//   return (i * n) + (1 - i)
// }
//
// Also for speed, the exported Binomial function includes checks to make sure
// the values are valid and the memo is large enough. But once these checks have
// been done once, they don't need to be repeated at each step of the recursion.
// So the unexported binomial function handles the actual recursion and skips
// these checks.

// Binomial coefficient of "n choose i".
func Binomial(n, i int) int {
	if n < 0 || i > n || i < 0 {
		return 0
	}
	if i > n/2 {
		i = n - i
	}
	if i < 2 {
		return (i * n) + (1 - i)
	}

	idx := n - 3
	idx = ((idx * idx) / 4) + i - 2

	if idx >= len(Memo) {
		ln := len(Memo) << 1
		for ; ln < idx; ln <<= 1 {
		}
		cp := make([]int, ln)
		copy(cp, Memo)
		Memo = cp
	}

	v := Memo[idx]
	if v == 0 {
		v = (n * binomial(n-1, i-1)) / i
		Memo[idx] = v
	}
	return v
}

func binomial(n, i int) int {
	if i < 2 {
		return (i * n) + (1 - i)
	}

	idx := n - 3
	idx = ((idx * idx) / 4) + i - 2
	v := Memo[idx]
	if v == 0 {
		v = (n * binomial(n-1, i-1)) / i
		Memo[idx] = v
	}
	return v
}
