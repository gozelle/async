package chunk

import (
	"fmt"
	"github.com/gozelle/async"
)

// Int64s 将 [a,b] 范围按指定步长分段
// 示例：[0 9 3] => [{0 2} {3 5} {6 8} {9 9}]
// 示例：[1 5 2] => [{1 2} {3 4} {5 5}]
// 示例：[0 3 5] => [{0 3}]
// 遍历示例：每小段都是左闭、右闭区间
// r, _ := Int64s(0,9,3)
// for _,v := range r {
//     for i:=v.Begin;i<= v.End; i++{
//          // ..
//      }
// }
func Int64s(a, b, step int64) (ranges []async.Int64Range, err error) {
	if step <= 0 {
		err = fmt.Errorf("expect step > 0,got: step=%d", step)
		return
	}
	
	if b < a {
		err = fmt.Errorf("expect b >= a,got: a=%d, b=%d", a, b)
		return
	}
	
	if a == b {
		ranges = append(ranges, async.Int64Range{Begin: a, End: b})
		return
	}
	
	i := int64(0)
	a -= 1
	for {
		aa := a + i*step + 1
		bb := a + (i+1)*step
		if bb < b && aa < b {
			ranges = append(ranges, async.Int64Range{Begin: aa, End: bb})
		} else {
			ranges = append(ranges, async.Int64Range{Begin: aa, End: b})
			if bb >= b {
				break
			}
		}
		i++
	}
	return
}
