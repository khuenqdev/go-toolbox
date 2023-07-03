package throttler

import "testing"

func TestApplyOperationToList(t *testing.T) {
    inputs := []int{1, 2, 3, 4}
    expOutputs := []int{2, 3, 4, 5}
    outputs, err := ApplyOperationToList[int, int](inputs, 5, 2, func(t int) (int, error) {
        return t + 1, nil
    })
    if err != nil {
        t.Errorf(err.Error())
    }

    if !outputs.Equals(expOutputs) {
        t.Errorf("got %q, wanted %q", outputs, expOutputs)
    }
}
