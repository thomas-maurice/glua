---@meta collections

---@class collections
local collections = {}

--- array mode: applies fn to every element, returning a new same-length array of the results
---@param t table the table to iterate, array part only (indices 1..#t)
---@param fn fun(v: any, i: integer): any callback invoked as fn(v, i); its return value becomes the new element at i
---@return any[] result a new array, same length as t
function collections.map(t, fn) end

--- array mode: returns a new, re-indexed array of the elements for which fn is truthy
---@param t table the table to iterate, array part only
---@param fn fun(v: any, i: integer): boolean predicate invoked as fn(v, i); the element is kept when this is truthy
---@return any[] result a new array of the matching elements, in original order
function collections.filter(t, fn) end

--- array mode: folds fn left-to-right over t starting from init, returning the final accumulator
---@param t table the table to iterate, array part only
---@param fn fun(acc: any, v: any, i: integer): any callback invoked as fn(acc, v, i); its return value becomes the next acc
---@param init any the required initial accumulator; pass nil explicitly if that is the desired seed
---@return any result the final accumulator after folding every element
function collections.reduce(t, fn, init) end

--- array mode: returns the first element and index for which fn is truthy, or nil, 0 if none
---@param t table the table to iterate, array part only
---@param fn fun(v: any, i: integer): boolean predicate invoked as fn(v, i)
---@return any value the first matching element, or nil if none matched
---@return number index the 1-based index of value, or 0 if none matched
function collections.find(t, fn) end

--- array mode: reports whether fn is truthy for at least one element (short-circuits)
---@param t table the table to iterate, array part only
---@param fn fun(v: any, i: integer): boolean predicate invoked as fn(v, i)
---@return boolean result true if any element matched; false for an empty table
function collections.any(t, fn) end

--- array mode: reports whether fn is truthy for every element (short-circuits, vacuously true when empty)
---@param t table the table to iterate, array part only
---@param fn fun(v: any, i: integer): boolean predicate invoked as fn(v, i)
---@return boolean result true if every element matched, or the table is empty
function collections.all(t, fn) end

--- array mode: groups elements by the key fn returns, preserving each group's encounter order
---@param t table the table to iterate, array part only
---@param fn fun(v: any, i: integer): any callback invoked as fn(v, i); its return value is the group key
---@return table<any, any[]> groups a table mapping each distinct key to a new array of its elements
function collections.group_by(t, fn) end

--- array mode: returns a new array stably sorted by the key fn returns (all-number or all-string keys only)
---@param t table the table to iterate, array part only
---@param fn fun(v: any): number|string callback invoked as fn(v) (no index); its return value is the sort key, and must be a number or string, uniformly across all elements
---@return any[] result a new array sorted ascending by key; equal keys keep their original relative order
function collections.sort_by(t, fn) end

--- array mode: splits t into a matching array and a rest array based on fn
---@param t table the table to iterate, array part only
---@param fn fun(v: any, i: integer): boolean predicate invoked as fn(v, i)
---@return any[] matching a new array of elements for which fn was truthy
---@return any[] rest a new array of the remaining elements
function collections.partition(t, fn) end

--- array mode: returns a new array with duplicates removed, keeping the first occurrence
---@param t table the table to iterate, array part only
---@return any[] result a new array; primitives dedup by value, tables/functions/userdata by identity
function collections.uniq(t) end

--- array mode, recursive: inlines nested array tables up to depth levels (-1 = fully)
---@param t table the table to flatten, array part only
---@param depth number how many levels of nested tables to inline; -1 flattens fully. Must be >= -1
---@return any[] result a new, flattened array
function collections.flatten(t, depth) end

--- array mode: returns a new array with t's elements in reverse order
---@param t table the table to reverse, array part only
---@return any[] result a new array, elements in reverse order
function collections.reverse(t) end

--- array mode on both tables: returns a new array of {a[i], b[i]} pairs, length min(#a, #b)
---@param a table the first table, array part only
---@param b table the second table, array part only
---@return any[][] result a new array of 2-element arrays; extra elements in the longer input are dropped
function collections.zip(a, b) end

--- array mode: splits t into new arrays of at most size elements each; the last chunk may be shorter
---@param t table the table to chunk, array part only
---@param size number the maximum size of each chunk; must be >= 1
---@return any[][] result a new array of arrays
function collections.chunk(t, size) end

--- map mode: returns a new array of all of t's keys, in unspecified order
---@param t table the table to read, array part and hash part together
---@return any[] keys a new array of t's keys; order is not defined
function collections.keys(t) end

--- map mode: returns a new array of all of t's values, in unspecified order
---@param t table the table to read, array part and hash part together
---@return any[] values a new array of t's values; order is not defined
function collections.values(t) end

--- map mode, shallow: returns a new table with every input table's keys, later tables winning on conflict
---@param t table the base table; never modified
---@param ... table additional tables applied left-to-right after t; never modified
---@return table<any, any> result a new table; a key present in more than one input takes the last input's value
function collections.merge(t, ...) end

--- map mode over names: returns a new table with only the named string-keyed fields of t
---@param t table the table to read from
---@param names string[] the string keys to keep; a name absent from t is silently skipped
---@return table<string, any> result a new table containing only the requested keys
function collections.pick(t, names) end

--- map mode: returns a new table with every key of t except the named string keys
---@param t table the table to read from, array part and hash part together
---@param names string[] the string keys to exclude; non-string keys are never excluded
---@return table<any, any> result a new table with the named keys removed
function collections.omit(t, names) end

--- map mode: walks t along path's dot-separated segments, returning default if any segment is missing or a non-table is indexed
---@param t table the table to walk
---@param path string dot-separated segments, e.g. "spec.containers.1.image"; a segment that parses as an integer is tried as a number key first, then as a string key
---@param default any returned as-is when the path cannot be fully resolved; pass nil explicitly for no default
---@return any value the value found at path, or default
function collections.get_path(t, path, default) end

--- general: reports whether a and b are structurally equal (same keys, recursively equal values); handles cycles
---@param a any the first value
---@param b any the second value
---@return boolean equal true if a and b are structurally equal; metatables are not compared
function collections.deep_equal(a, b) end

--- general: returns a structural copy of v; nested tables are new tables, cycles are preserved
---@param v any the value to copy
---@return any copy a deep copy of v; functions and userdata are copied by reference, metatables are not copied
function collections.deep_copy(v) end

return collections
