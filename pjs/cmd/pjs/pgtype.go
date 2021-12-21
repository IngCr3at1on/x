package main

import (
	"github.com/ingcr3at1on/x/pjs"
	"github.com/jackc/pgtype"
)

func addToOutput(rec pjs.Receiver) interface{} {
	switch rec.DataTypeOID {
	case pgtype.ByteaOID:
		val := rec.Val.(*pgtype.ByteaArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.CIDRArrayOID:
		val := rec.Val.(*pgtype.CIDRArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.BoolArrayOID:
		val := rec.Val.(*pgtype.BoolArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.Int2ArrayOID:
		val := rec.Val.(*pgtype.Int2Array)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.Int4ArrayOID:
		val := rec.Val.(*pgtype.Int4Array)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.TextArrayOID:
		val := rec.Val.(*pgtype.TextArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.ByteaArrayOID:
		val := rec.Val.(*pgtype.ByteaArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.BPCharArrayOID:
		val := rec.Val.(*pgtype.BPCharArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.VarcharArrayOID:
		val := rec.Val.(*pgtype.VarcharArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.Int8ArrayOID:
		val := rec.Val.(*pgtype.Int8Array)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.Float4ArrayOID:
		val := rec.Val.(*pgtype.Float4Array)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.Float8ArrayOID:
		val := rec.Val.(*pgtype.Float8Array)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.ACLItemArrayOID:
		val := rec.Val.(*pgtype.ACLItemArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.InetArrayOID:
		val := rec.Val.(*pgtype.InetArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.TimestampArrayOID:
		val := rec.Val.(*pgtype.TimestampArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.DateArrayOID:
		val := rec.Val.(*pgtype.DateArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.TimestamptzArrayOID:
		val := rec.Val.(*pgtype.TimestamptzArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.NumericArrayOID:
		val := rec.Val.(*pgtype.NumericArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.UUIDArrayOID:
		val := rec.Val.(*pgtype.UUIDArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.JSONBArrayOID:
		val := rec.Val.(*pgtype.JSONBArray)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return []interface{}{}
		default:
			return val.Elements
		}
	case pgtype.NameOID:
		val := rec.Val.(*pgtype.Name)
		switch val.Status {
		case pgtype.Null:
			return nil
		case pgtype.Undefined:
			return ``
		default:
			return val.String
		}
	default:
		return rec.Val
	}
}
