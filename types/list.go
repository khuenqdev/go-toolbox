package types

import (
    "crypto/sha256"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "reflect"
)

type List[T any] []T

func (li *List[T]) MarshalJSON() ([]byte, error) {
    return json.Marshal(li.ToSlice())
}

func (li *List[T]) UnmarshalJSON(data []byte) error {
    var result []T
    if err := json.Unmarshal(data, &result); err != nil {
        return err
    }

    if result == nil {
        return errors.New("the unmarshalled result list is nil")
    }

    *li = result

    return nil
}

func (li *List[T]) MarshalGQL(w io.Writer) {
    data, err := li.MarshalJSON()
    if err != nil {
        _, _ = w.Write([]byte("[]"))
        return
    }

    _, _ = w.Write(data)
}

func (li *List[T]) UnmarshalGQL(v interface{}) error {
    data, err := json.Marshal(v)
    if err != nil {
        return err
    }

    if err := li.UnmarshalJSON(data); err != nil {
        return err
    }

    return nil
}

// First get the first element
func (li *List[T]) First() (T, error) {
    if !li.Any() {
        return *new(T), errors.New("the list is empty")
    }

    tmp := *li
    return tmp[0], nil
}

// FirstOrDefault get the first element or return the default if the list is empty
func (li *List[T]) FirstOrDefault(predicate func(t T) bool) (T, error) {
    if predicate == nil {
        return *new(T), errors.New("predicate function must be defined")
    }

    for _, item := range *li {
        if predicate(item) {
            return item, nil
        }
    }

    return *new(T), nil
}

// Select get a list of values extracted by calling the selector function on the list
func (li *List[T]) Select(selector func(t T) (any, error)) (List[any], error) {
    if selector == nil {
        return nil, errors.New("selector function must be defined")
    }

    if !li.Any() {
        return List[any]{}, nil
    }

    var results List[any]
    for _, i := range *li {
        item, err := selector(i)
        if err != nil {
            return List[any]{}, err
        }

        if item == nil {
            continue
        }

        results = append(results, item)
    }

    return results, nil
}

// Any check if the list has any element at all
func (li *List[T]) Any() bool {
    return len(*li) > 0
}

// Find obtain an element that matches a certain condition
func (li *List[T]) Find(match func(t T) bool) (T, error) {
    if match == nil {
        return *new(T), errors.New("matching function must be defined")
    }

    if !li.Any() {
        return *new(T), errors.New(fmt.Sprintf("the list of type %s is empty", *new(T)))
    }

    for _, item := range *li {
        if match(item) {
            return item, nil
        }
    }

    return *new(T), errors.New(fmt.Sprintf("could not find item of type %T in the list", *new(T)))
}

// Last get the last element filtered by the predicate function
func (li *List[T]) Last(predicate func(t T) bool) (T, error) {
    if predicate == nil {
        return *new(T), errors.New("predicate function must be defined")
    }

    if !li.Any() {
        return *new(T), errors.New("the list is empty")
    }

    for _, item := range *li {
        if predicate(item) {
            return item, nil
        }
    }

    return *new(T), errors.New("no element satisfies the condition")
}

// ToSlice convert the list type to a slice
func (li *List[T]) ToSlice() []T {
    var results []T

    for _, item := range *li {
        results = append(results, item)
    }

    return results
}

// Count return the number of elements in the list
func (li *List[T]) Count() int {
    return len(*li)
}

// Contains check if an element exists in the list
func (li *List[T]) Contains(value T) bool {
    if !li.Any() {
        return false
    }

    for _, item := range *li {
        itemBytes, _ := json.Marshal(item)
        valueBytes, _ := json.Marshal(value)
        if itemBytes != nil && valueBytes != nil && string(itemBytes) == string(valueBytes) {
            return true
        }
    }

    return false
}

// Aggregate applies an accumulator function over a sequence. The specified seed value is used as the initial accumulator value.
func (li *List[T]) Aggregate(seed any, aggregator func(accumulator any, t T) any) (any, error) {
    if aggregator == nil {
        return 0, errors.New("aggregator function must be defined")
    }

    if !li.Any() {
        return 0, errors.New("the list is empty")
    }

    for _, item := range *li {
        seed = aggregator(seed, item)
    }

    return seed, nil
}

// Where filters a sequence of values based on a predicate.
func (li *List[T]) Where(predicate func(t T) bool) (List[T], error) {
    if predicate == nil {
        return nil, errors.New("predicate function must be defined")
    }

    if !li.Any() {
        return nil, errors.New("the list is empty")
    }

    var results List[T]

    for _, item := range *li {
        if predicate(item) {
            results = append(results, item)
        }
    }

    return results, nil
}

// Distinct returns distinct elements from a sequence by using the default equality comparer to compare values.
func (li *List[T]) Distinct() List[T] {
    if !li.Any() {
        return *li
    }

    distinctMap := make(map[string]T, 0)
    var results List[T]

    for _, item := range *li {
        itemBytes, _ := json.Marshal(item)
        h := sha256.New()
        h.Write(itemBytes)
        hash := fmt.Sprintf("%x", h.Sum(nil))
        distinctMap[hash] = item
    }

    if distinctMap == nil || len(distinctMap) == 0 {
        return *li
    }

    for _, t := range distinctMap {
        results = append(results, t)
    }

    return results
}

// All determines whether all elements of a sequence satisfy a condition.
func (li *List[T]) All(predicate func(t T) bool) (bool, error) {
    if predicate == nil {
        return false, errors.New("predicate function must be defined")
    }

    if !li.Any() {
        return false, errors.New("the list is empty")
    }

    count := 0
    for _, item := range *li {
        if predicate(item) {
            count++
        }
    }

    return count == len(*li), nil
}

// Equals check whether two lists equal each other
func (li *List[T]) Equals(otherList List[T]) bool {
    if !li.Any() || otherList == nil || !otherList.Any() {
        return false
    }

    return reflect.DeepEqual(*li, otherList)
}
