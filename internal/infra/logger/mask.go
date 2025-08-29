// package redact (ou onde preferir)

package logger

import (
	"reflect"
	"strings"
)

const Mask = "********"

// RedactStruct faz uma cópia de v e mascara campos conforme paths.
// Paths podem ser:
//   - nomes simples: "password", "password_confirm" (match por nome do json tag, em qualquer nível)
//   - dot-notation:  "password", "addresses.0.document", "body.password"
func RedactStruct[T any](v T, paths ...string) T {
	if len(paths) == 0 {
		return v
	}
	// separa full paths e nomes simples
	full := make(map[string]struct{})
	names := make(map[string]struct{})
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, ".") {
			full[p] = struct{}{}
		} else {
			names[p] = struct{}{}
		}
	}
	rv := reflect.ValueOf(v)
	out := deepRedact(rv, full, names, "").Interface()
	return out.(T)
}

func deepRedact(rv reflect.Value, full map[string]struct{}, names map[string]struct{}, prefix string) reflect.Value {
	if !rv.IsValid() {
		return rv
	}
	// desreferencia ponteiros/interfaces
	for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return rv
		}
		elem := rv.Elem()
		red := deepRedact(elem, full, names, prefix)
		// re-empacota
		if rv.Kind() == reflect.Pointer {
			ptr := reflect.New(elem.Type())
			ptr.Elem().Set(red)
			return ptr
		}
		return red
	}

	switch rv.Kind() {
	case reflect.Struct:
		out := reflect.New(rv.Type()).Elem()
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			fv := rv.Field(i)
			ft := rt.Field(i)
			// ignora campos não exportados
			if ft.PkgPath != "" {
				out.Field(i).Set(fv)
				continue
			}
			jsonName := jsonTagName(ft)
			if jsonName == "-" {
				out.Field(i).Set(fv)
				continue
			}
			// se não tem tag, usa o nome do campo
			if jsonName == "" {
				jsonName = ft.Name
			}
			path := jsonName
			if prefix != "" {
				path = prefix + "." + jsonName
			}
			// match por full path
			if _, ok := full[path]; ok {
				out.Field(i).Set(reflect.ValueOf(Mask).Convert(ft.Type))
				continue
			}
			// match por nome simples (qualquer nível)
			if _, ok := names[jsonName]; ok {
				out.Field(i).Set(reflect.ValueOf(Mask).Convert(ft.Type))
				continue
			}
			// recursão
			out.Field(i).Set(deepRedact(fv, full, names, path))
		}
		return out

	case reflect.Map:
		// só suporta map[string]any / map[string]T
		if rv.Type().Key().Kind() != reflect.String {
			return rv
		}
		out := reflect.MakeMapWithSize(rv.Type(), rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key()
			kstr := k.String()
			path := kstr
			if prefix != "" {
				path = prefix + "." + kstr
			}
			if _, ok := full[path]; ok {
				out.SetMapIndex(k, reflect.ValueOf(Mask).Convert(rv.Type().Elem()))
				continue
			}
			if _, ok := names[kstr]; ok {
				out.SetMapIndex(k, reflect.ValueOf(Mask).Convert(rv.Type().Elem()))
				continue
			}
			out.SetMapIndex(k, deepRedact(iter.Value(), full, names, path))
		}
		return out

	case reflect.Slice, reflect.Array:
		n := rv.Len()
		out := reflect.MakeSlice(rv.Type(), n, n)
		for i := 0; i < n; i++ {
			idxPath := prefix + "." + itoa(i)
			out.Index(i).Set(deepRedact(rv.Index(i), full, names, idxPath))
		}
		return out

	default:
		// tipos primitivos: se o prefixo completo bater, mascara
		if _, ok := full[prefix]; ok && prefix != "" {
			return reflect.ValueOf(Mask).Convert(rv.Type())
		}
		return rv
	}
}

func jsonTagName(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" {
		return ""
	}
	// pega parte antes da vírgula (ex.: "password,omitempty")
	if i := strings.Index(tag, ","); i >= 0 {
		return tag[:i]
	}
	return tag
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	n := i
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
