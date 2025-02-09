package lmap_test

// 	expected := map[string]int{
// 		"1":  1,
// 		"2":  2,
// 		"4":  4,
// 		"5":  5,
// 		"6":  6,
// 		"8":  8,
// 		"9":  9,
// 		"10": 10,
// 	}

// 	got := lmap.TransformMap(m, tf).Map()
// 	assert.Equal(t, expected, got)

// 	expected := map[string]int{
// 		"1":  1,
// 		"2":  2,
// 		"4":  4,
// 		"5":  5,
// 		"6":  6,
// 		"8":  8,
// 		"9":  9,
// 		"10": 10,
// 	}

// 	got := lmap.TransformMap(m, tf).Map()
// 	assert.Equal(t, expected, got)

// 	buf := lmap.New(map[string]int{
// 		"100": 100,
// 	})
// 	expected["100"] = 100
// 	got = lmap.Transform(lmap.New(m), buf, tf).Map()
// 	assert.Equal(t, expected, got)
// }

// func TestSliceTransform(t *testing.T) {
// 	tf := lmap.NewSliceTransformFunc(func(i int, s string) (string, bool) {
// 		return fmt.Sprintf("%d %s", i, s), s != ""
// 	})

// 	m := map[int]string{
// 		1:  "A",
// 		2:  "B",
// 		3:  "",
// 		4:  "D",
// 		5:  "E",
// 		6:  "F",
// 		7:  "",
// 		8:  "H",
// 		9:  "I",
// 		10: "J",
// 	}

// 	got := tf.TransformMap(m).Sort(slice.LT[string]())
// 	expected := slice.New([]string{"1 A", "10 J", "2 B", "4 D", "5 E", "6 F", "8 H", "9 I"})
// 	assert.Equal(t, expected, got)

// 	got = lmap.SliceTransformMap(m, tf).Sort(slice.LT[string]())
// 	assert.Equal(t, expected, got)

// 	got = lmap.SliceTransform(lmap.New(m), nil, tf).Sort(slice.LT[string]())
// 	assert.Equal(t, expected, got)
// }
