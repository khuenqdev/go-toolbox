package json

import "encoding/json"

// Reserialize Take an object as an input then convert it to another object via JSON serialization
func Reserialize[T any](obj any) T {
    serialized, err := json.Marshal(obj)
    if err != nil {
        return *new(T)
    }

    result := new(T)
    err = json.Unmarshal(serialized, &result)
    if err != nil {
        return *new(T)
    }

    return *result
}
