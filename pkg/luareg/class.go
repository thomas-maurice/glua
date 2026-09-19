// Copyright (c) 2024-2025 Thomas Maurice
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package luareg

import (
	"errors"
	"fmt"
	"reflect"

	glua "github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
)

// AnyClass: non-generic interface stored by Module and Registry. Both Class[T]
// and any future class kinds satisfy it. Used by stub generators (A3) without
// needing to know T.
type AnyClass interface {
	// Name: returns the Lua type name (e.g. "log.Logger").
	Name() string
	// Doc: returns the class documentation string.
	Doc() string
	// MethodsMeta: returns all registered method metadata in registration order.
	MethodsMeta() []*FnMeta
	// GoType: returns the Go reflect.Type for the receiver. Used by stub generators
	// (A3) to match function return types to their registered class name.
	GoType() reflect.Type
	// register: builds the metatable for this class in L. Called by Module.PushTo.
	register(L *lua.LState, classLookup func(reflect.Type) AnyClass)
	// goType: returns the Go reflect.Type for the receiver. Used for auto-wrap
	// lookup when a function returns an instance of this class.
	goType() reflect.Type
}

// Class: a generic Lua class backed by Go type T. T is typically a pointer type
// (e.g. *Logger); the framework supports both value and pointer T.
//
// Example:
//
//	c := luareg.NewClass[*Logger]("log.Logger", "structured logger object")
//	c.Method("info", (*Logger).Info, "log at info level")
//	m.RegisterClass(c)
type Class[T any] struct {
	name    string
	doc     string
	methods []*FnMeta
	byName  map[string]*FnMeta
	typ     reflect.Type // reflect.Type of T
	lenFn   func(T) int  // optional __len backing function; nil means no __len
}

// NewClass: creates a Class with the given Lua type name and summary doc.
// The Go type is captured via the type parameter — no instance is required.
//
// Example:
//
//	c := luareg.NewClass[*Logger]("log.Logger", "structured logger")
func NewClass[T any](luaName, doc string) *Class[T] {
	// reflect.TypeOf((*T)(nil)).Elem() captures the exact type of T at the
	// call site, including whether T is a pointer (e.g. *Logger) or a value.
	typ := reflect.TypeOf((*T)(nil)).Elem()
	return &Class[T]{
		name:   luaName,
		doc:    doc,
		byName: make(map[string]*FnMeta),
		typ:    typ,
	}
}

// Name: returns the class's Lua type name.
func (c *Class[T]) Name() string { return c.name }

// Doc: returns the class documentation string.
func (c *Class[T]) Doc() string { return c.doc }

// MethodsMeta: returns all registered method metadata in registration order.
func (c *Class[T]) MethodsMeta() []*FnMeta { return c.methods }

// GoType: returns the reflect.Type of T. Satisfies AnyClass for stub generators.
func (c *Class[T]) GoType() reflect.Type { return c.typ }

// goType: returns the reflect.Type of T. Used internally for auto-wrap lookup.
func (c *Class[T]) goType() reflect.Type { return c.typ }

// Method: registers a Lua method on this class. methodFn is a Go method
// expression such as (*Logger).Info — a function whose first parameter is the
// receiver T. doc and opts work like Module.Fn.
//
// Panics if luaName was already registered on this class.
//
// Returns the receiver for chaining.
func (c *Class[T]) Method(luaName string, methodFn any, doc string, opts ...FnOpt) *Class[T] {
	if _, exists := c.byName[luaName]; exists {
		panic(fmt.Sprintf("luareg: class %q: duplicate registration of method %q", c.name, luaName))
	}

	validateGoFn(luaName, methodFn)

	// Validate that the first parameter is T (or compatible).
	ft := reflect.TypeOf(methodFn)
	if ft.NumIn() == 0 {
		panic(fmt.Sprintf("luareg: class %q: method %q: methodFn must have at least one parameter (the receiver)", c.name, luaName))
	}
	receiverType := ft.In(0)
	if receiverType != c.typ {
		panic(fmt.Sprintf("luareg: class %q: method %q: first parameter must be %s, got %s", c.name, luaName, c.typ, receiverType))
	}

	meta := &FnMeta{
		LuaName: luaName,
		GoFn:    methodFn,
		Doc:     doc,
	}
	for _, opt := range opts {
		opt(meta)
	}

	c.methods = append(c.methods, meta)
	c.byName[luaName] = meta
	return c
}

// Len: registers an OPT-IN __len metamethod for this class, backed by fn.
// There is no generically sensible definition of "length" for an arbitrary
// wrapped Go type — most classes (loggers, clients, ...) have none at all —
// so __len is only set when a class explicitly opts in via this method,
// rather than bolting an invented meaning (e.g. method count) onto every
// class.
//
// Example:
//
//	c.Len(func(q *Queue) int { return len(q.items) })
//
// Returns the receiver for chaining.
func (c *Class[T]) Len(fn func(T) int) *Class[T] {
	c.lenFn = fn
	return c
}

// Wrap: creates a UserData around an instance of T, bound to this class's
// metatable. L must have the metatable already registered (i.e. PushTo has
// been called on the Module that owns this class).
//
// Used internally by the auto-wrap return-value path, and exposed for code
// that needs to push instances manually.
func (c *Class[T]) Wrap(L *lua.LState, instance T) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = instance
	L.SetMetatable(ud, L.GetTypeMetatable(c.name))
	return ud
}

// Check: extracts an instance of T from a UserData at Lua stack position pos.
// Raises a Lua error via L.ArgError if the value is not the right type.
func (c *Class[T]) Check(L *lua.LState, pos int) T {
	ud := L.CheckUserData(pos)
	if v, ok := ud.Value.(T); ok {
		return v
	}
	L.ArgError(pos, fmt.Sprintf("%s expected", c.name))
	var zero T
	return zero
}

// classTypeMarkerField: hidden field stored on each class's type metatable,
// holding the reflect.Type this Lua class name is bound to. Used at
// registration time by checkClassTypeCollision.
const classTypeMarkerField = "__glua_class_gotype"

// register: builds the metatable for this class in L and registers all methods.
// classLookup is used so method wrappers can auto-wrap return values of class types.
//
// L.NewTypeMetatable is keyed globally in L's registry by Lua name alone —
// uniqueness of c.name is only enforced within a single Module (see
// Module.classByName). If two different modules register a class under the
// same Lua name on the same LState, both would otherwise share the same
// metatable, and the second register() call would silently repoint every
// already-wrapped instance of the first class at the second class's methods.
// checkClassTypeCollision panics loudly instead of allowing that, mirroring
// how validateGoFn already panics on bad registrations at setup time.
func (c *Class[T]) register(L *lua.LState, classLookup func(reflect.Type) AnyClass) {
	mt := L.NewTypeMetatable(c.name)
	checkClassTypeCollision(L, mt, c.name, c.typ)

	methods := L.NewTable()
	for _, meta := range c.methods {
		wrapper := buildMethodWrapper(c, meta, classLookup)
		L.SetField(methods, meta.LuaName, L.NewFunction(wrapper))
	}
	L.SetField(mt, "__index", methods)

	// __metatable hides the real methods table from getmetatable/setmetatable.
	// gopher-lua diverges from stock Lua 5.1 here: baseSetMetatable only
	// type-checks arg 2 (the new metatable), not arg 1, so without this a Lua
	// script could call setmetatable(instance, {}) and strip every method off
	// a live instance. Setting __metatable makes gopher-lua's baselib raise
	// "cannot change a protected metatable" on setmetatable, and makes
	// getmetatable return this string instead of the real methods table.
	L.SetField(mt, "__metatable", lua.LString(c.name))

	// __tostring: "<class name>: <pointer-or-value>". Pointer receivers print
	// their address (stable, short, avoids dumping a potentially huge struct
	// like a k8sclient.Client); value receivers print via %v.
	className := c.name
	L.SetField(mt, "__tostring", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		v := reflect.ValueOf(ud.Value)
		if v.Kind() == reflect.Pointer {
			L.Push(lua.LString(fmt.Sprintf("%s: %p", className, ud.Value)))
		} else {
			L.Push(lua.LString(fmt.Sprintf("%s: %v", className, ud.Value)))
		}
		return 1
	}))

	// __eq: compares the WRAPPED Go values, not the *lua.LUserData wrappers,
	// so two userdata wrapping the same underlying value (e.g. two returns of
	// the same *Logger) compare equal instead of relying on default identity
	// comparison of the wrapper. Lua only invokes __eq when both operands are
	// userdata sharing this exact metamethod, so both values are guaranteed
	// to be T here.
	L.SetField(mt, "__eq", L.NewFunction(func(L *lua.LState) int {
		a := L.CheckUserData(1)
		b := L.CheckUserData(2)
		L.Push(lua.LBool(classValuesEqual(a.Value, b.Value)))
		return 1
	}))

	// __len is opt-in via Class.Len — see there for why it isn't set by default.
	if c.lenFn != nil {
		lenFn := c.lenFn
		L.SetField(mt, "__len", L.NewFunction(func(L *lua.LState) int {
			receiver := c.Check(L, 1)
			L.Push(lua.LNumber(lenFn(receiver)))
			return 1
		}))
	}
}

// classValuesEqual: compares two wrapped Go values for __eq. Pointer types
// compare by pointer identity (matching Go's own == semantics); other
// comparable types (structs of comparable fields, primitives) compare by
// value. Go's == panics on uncomparable dynamic types (slices, maps, funcs,
// or structs embedding them); rather than risk that panic reaching a Lua
// script — or silently falling back to userdata-wrapper identity, which the
// __eq contract here explicitly avoids — an uncomparable Go type is defined
// to never equal anything via __eq, including another wrapper around the
// exact same value. This is a conservative, documented default: it never
// panics and never reports a false positive, at the cost of "==" being
// always false for classes backed by such types.
func classValuesEqual(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	if !ta.Comparable() {
		return false
	}
	return a == b
}

// checkClassTypeCollision: panics if the type metatable for luaName is
// already bound to a Go type other than typ. Because L.NewTypeMetatable is
// keyed globally by name across every module sharing L, the marker recording
// which Go type "owns" a name is stored in a hidden field on the metatable
// itself rather than in a package-level registry keyed by name. A
// package-level map would leak across independent lua.LState instances —
// this codebase's own tests routinely reuse short class names like "Foo" or
// "m.Counter" across many short-lived LStates, which would produce false
// positives once one test's name choice collided with another's. Scoping the
// marker to the metatable means collisions are only detected where they can
// actually occur: two classes sharing one Lua state.
//
// Registering the exact same (name, type) pair again — e.g. Module.PushTo
// called more than once — is legal and a no-op here.
func checkClassTypeCollision(L *lua.LState, mt *lua.LTable, luaName string, typ reflect.Type) {
	existing := L.GetField(mt, classTypeMarkerField)
	if existing == lua.LNil {
		marker := L.NewUserData()
		marker.Value = typ
		L.SetField(mt, classTypeMarkerField, marker)
		return
	}

	ud, ok := existing.(*lua.LUserData)
	if !ok {
		// Unreachable in practice — nothing else writes this field — but fail
		// loud rather than silently ignoring a corrupted marker.
		panic(fmt.Sprintf("luareg: class %q: type marker field was overwritten with a non-userdata value", luaName))
	}
	existingType, ok := ud.Value.(reflect.Type)
	if !ok || existingType != typ {
		panic(fmt.Sprintf(
			"luareg: class name %q is already registered for Go type %v; cannot also register it for %v — Lua class names must be unique across every module pushed to the same Lua state",
			luaName, existingType, typ,
		))
	}
}

// buildMethodWrapper: reflects on meta.GoFn (a method expression whose first
// param is T) and returns an lua.LGFunction that:
//  1. Extracts the receiver T from stack position 1 via c.Check.
//  2. Converts remaining Lua args to the method's non-receiver params.
//  3. Calls the method.
//  4. Pushes return values, auto-wrapping class instances via classLookup.
func buildMethodWrapper[T any](c *Class[T], meta *FnMeta, classLookup func(reflect.Type) AnyClass) lua.LGFunction {
	fn := reflect.ValueOf(meta.GoFn)
	ft := fn.Type()

	// ft.In(0) is T (receiver). ft.In(1) may be *lua.LState (escape hatch).
	// fixedLuaArgCount = number of FIXED (non-variadic, non-LState) Lua args
	// expected AFTER the receiver (self). The receiver is always arg 1 in Lua
	// (via : syntax).
	receiverParamCount := 1
	secondParamIsLState := ft.NumIn() > 1 && ft.In(1) == luaStateType
	isVariadic := ft.IsVariadic()
	fixedLuaArgCount := ft.NumIn() - receiverParamCount
	if secondParamIsLState {
		fixedLuaArgCount-- // *lua.LState does not count as a Lua arg
	}
	if isVariadic {
		fixedLuaArgCount-- // variadic tail does not count as a fixed arg
	}

	// Determine trailing error return.
	numOut := ft.NumOut()
	hasErrorReturn := numOut > 0 && ft.Out(numOut-1).Implements(errorType)
	nonErrReturns := numOut
	if hasErrorReturn {
		nonErrReturns--
	}

	// rawLGFunction: true when the method signature is func(T, *lua.LState) int
	// (or func(T, *lua.LState) with no return). The method manages its own Lua
	// stack pushes and the int return IS the Lua return count. Variadic
	// signatures never qualify.
	rawLGFunction := !isVariadic && secondParamIsLState && numOut == 1 && ft.Out(0).Kind() == reflect.Int

	translator := glua.NewTranslator()

	return func(L *lua.LState) int {
		// Stack: arg1=self(userdata), arg2..N=method args.
		// With *lua.LState escape hatch OR variadic: enforce minimum
		// (self + fixedLuaArgCount), allow extras. Without either: enforce exact.
		if secondParamIsLState || isVariadic {
			if fixedLuaArgCount > 0 && L.GetTop() < 1+fixedLuaArgCount {
				L.RaiseError("expected at least %d argument(s), got %d", fixedLuaArgCount, L.GetTop()-1)
				return 0
			}
		} else {
			expectedTop := 1 + fixedLuaArgCount
			if L.GetTop() != expectedTop {
				L.RaiseError("expected %d argument(s), got %d", fixedLuaArgCount, L.GetTop()-1)
				return 0
			}
		}

		// Extract receiver.
		receiver := c.Check(L, 1)

		// Build Go argument slice: [receiver, (L if escape hatch), arg2, ...].
		args := make([]reflect.Value, ft.NumIn())
		args[0] = reflect.ValueOf(receiver)

		argOffset := 1
		if secondParamIsLState {
			args[1] = reflect.ValueOf(L)
			argOffset = 2
		}

		// Convert fixed Lua args (positions 2..fixedLuaArgCount+1).
		for i := 0; i < fixedLuaArgCount; i++ {
			paramType := ft.In(i + argOffset)
			luaPos := i + 2 // Lua positions start at 1; pos 1 is self
			val, err := luaToGo(L, luaPos, paramType, translator)
			if err != nil {
				L.ArgError(luaPos, err.Error())
				return 0
			}
			args[i+argOffset] = val
		}

		// Collect the variadic tail (if any) into a typed slice.
		// Lua pos 1 is self, fixed args occupy pos 2..1+fixedLuaArgCount, so
		// the variadic tail begins at pos 2+fixedLuaArgCount.
		if isVariadic {
			sliceType := ft.In(ft.NumIn() - 1)
			elemType := sliceType.Elem()
			variadicStart := 2 + fixedLuaArgCount
			top := L.GetTop()
			nVariadic := top - (variadicStart - 1)
			if nVariadic < 0 {
				nVariadic = 0
			}
			variadicSlice := reflect.MakeSlice(sliceType, nVariadic, nVariadic)
			for i := 0; i < nVariadic; i++ {
				luaPos := variadicStart + i
				val, err := luaToGo(L, luaPos, elemType, translator)
				if err != nil {
					L.ArgError(luaPos, err.Error())
					return 0
				}
				variadicSlice.Index(i).Set(val)
			}
			args[ft.NumIn()-1] = variadicSlice
		}

		// Call the method. CallSlice passes args[len-1] as the variadic slice;
		// Call handles the non-variadic case.
		var results []reflect.Value
		if isVariadic {
			results = fn.CallSlice(args)
		} else {
			results = fn.Call(args)
		}

		// Raw lua.LGFunction pattern: the method manages its own stack pushes
		// and the single int return IS the Lua stack count. Return it directly.
		if rawLGFunction {
			return int(results[0].Int())
		}

		// Handle error return.
		if hasErrorReturn {
			errVal := results[len(results)-1]
			if !errVal.IsNil() {
				raiseError(L, errVal.Interface().(error))
				return 0
			}
			results = results[:len(results)-1]
		}

		// Push non-error return values; auto-wrap class instances.
		for _, rv := range results {
			lv, err := goToLuaWithClasses(L, rv, translator, classLookup)
			if err != nil {
				L.RaiseError("return conversion error: %s", err.Error())
				return 0
			}
			L.Push(lv)
		}

		return nonErrReturns
	}
}

// raiseError: raises a Lua error from a Go error. If err is *Error, raises a
// structured table {message=..., kind=...}; otherwise raises a plain string.
func raiseError(L *lua.LState, err error) {
	var luaErr *Error
	if ok := isLuaError(err, &luaErr); ok {
		tbl := L.NewTable()
		L.SetField(tbl, "message", lua.LString(luaErr.Message))
		L.SetField(tbl, "kind", lua.LString(luaErr.Kind))
		L.Error(tbl, 1)
		return
	}
	L.RaiseError("%s", err.Error())
}

// isLuaError: checks if err is (or wraps) a *Error via errors.As.
// Wrapped errors (e.g. fmt.Errorf("context: %w", &Error{...})) are still
// detected and produce a structured Lua error table — the wrapper context
// is lost but the structured kind/message is preserved, which is what most
// callers want for `pcall`-based error classification.
func isLuaError(err error, target **Error) bool {
	return errors.As(err, target)
}

// goToLuaWithClasses: like goToLua but checks classLookup first so that
// return values matching a registered class type are wrapped as UserData
// instead of going through the Translator.
func goToLuaWithClasses(L *lua.LState, v reflect.Value, tr *glua.Translator, classLookup func(reflect.Type) AnyClass) (lua.LValue, error) {
	if classLookup != nil {
		// Check if the value's type matches a registered class.
		// We check both the actual type and, for pointers, the pointer type.
		t := v.Type()
		if cls := classLookup(t); cls != nil {
			// Nil check for pointer types.
			if t.Kind() == reflect.Pointer && v.IsNil() {
				return lua.LNil, nil
			}
			return wrapAnyClass(L, cls, v), nil
		}
	}
	return goToLua(L, v, tr)
}

// wrapAnyClass: wraps a reflect.Value as UserData bound to cls's metatable.
// cls must be the AnyClass whose goType() matches v.Type().
func wrapAnyClass(L *lua.LState, cls AnyClass, v reflect.Value) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = v.Interface()
	L.SetMetatable(ud, L.GetTypeMetatable(cls.Name()))
	return ud
}

// Error: a typed Go error that, when returned from a registered function or
// method, raises a Lua table {message=..., kind=...} via L.Error. Modules can
// return &luareg.Error{Kind:"NotFound", Message:"pods 'foo' not found"} to
// enable pcall callers to distinguish error kinds.
//
// Example:
//
//	return nil, &luareg.Error{Kind: "NotFound", Message: "item not found"}
type Error struct {
	Kind    string
	Message string
}

// Error: implements the error interface.
func (e *Error) Error() string { return e.Message }
