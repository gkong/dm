package gen

// NOTE: THIS FILE WAS PRODUCED BY THE
// ZEBRAPACK CODE GENERATION TOOL (github.com/glycerine/zebrapack)
// DO NOT EDIT

import "github.com/glycerine/zebrapack/msgp"

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *EmailSessData) DecodeMsg(dc *msgp.Reader) (err error) {
	var sawTopNil bool
	if dc.IsNil() {
		sawTopNil = true
		err = dc.ReadNil()
		if err != nil {
			return
		}
		dc.PushAlwaysNil()
	}

	var field []byte
	_ = field
	const maxFields0zxsz = 2

	// -- templateDecodeMsgZid starts here--
	var totalEncodedFields0zxsz uint32
	totalEncodedFields0zxsz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft0zxsz := totalEncodedFields0zxsz
	missingFieldsLeft0zxsz := maxFields0zxsz - totalEncodedFields0zxsz

	var nextMiss0zxsz int = -1
	var found0zxsz [maxFields0zxsz]bool
	var curField0zxsz int

doneWithStruct0zxsz:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft0zxsz > 0 || missingFieldsLeft0zxsz > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft0zxsz, missingFieldsLeft0zxsz, msgp.ShowFound(found0zxsz[:]), decodeMsgFieldOrder0zxsz)
		if encodedFieldsLeft0zxsz > 0 {
			encodedFieldsLeft0zxsz--
			curField0zxsz, err = dc.ReadInt()
			if err != nil {
				return
			}
		} else {
			//missing fields need handling
			if nextMiss0zxsz < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss0zxsz = 0
			}
			for nextMiss0zxsz < maxFields0zxsz && (found0zxsz[nextMiss0zxsz] || decodeMsgFieldSkip0zxsz[nextMiss0zxsz]) {
				nextMiss0zxsz++
			}
			if nextMiss0zxsz == maxFields0zxsz {
				// filled all the empty fields!
				break doneWithStruct0zxsz
			}
			missingFieldsLeft0zxsz--
			curField0zxsz = nextMiss0zxsz
		}
		//fmt.Printf("switching on curField: '%v'\n", curField0zxsz)
		switch curField0zxsz {
		// -- templateDecodeMsgZid ends here --

		case 0:
			// zid 0 for "OldUserKey"
			found0zxsz[0] = true
			z.OldUserKey, err = dc.ReadBytes(z.OldUserKey)
			if err != nil {
				return
			}
		case 1:
			// zid 1 for "NewUserKey"
			found0zxsz[1] = true
			z.NewUserKey, err = dc.ReadBytes(z.NewUserKey)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss0zxsz != -1 {
		dc.PopAlwaysNil()
	}

	if sawTopNil {
		dc.PopAlwaysNil()
	}

	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of EmailSessData
var decodeMsgFieldOrder0zxsz = []string{"OldUserKey", "NewUserKey"}

var decodeMsgFieldSkip0zxsz = []bool{false, false}

// fieldsNotEmpty supports omitempty tags
func (z *EmailSessData) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 2
	}
	var fieldsInUse uint32 = 2
	isempty[0] = (len(z.OldUserKey) == 0) // string, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[1] = (len(z.NewUserKey) == 0) // string, omitempty
	if isempty[1] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z *EmailSessData) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zhji [2]bool
	fieldsInUse_zqpm := z.fieldsNotEmpty(empty_zhji[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zqpm)
	if err != nil {
		return err
	}

	if !empty_zhji[0] {
		// zid 0 for "OldUserKey"
		err = en.Append(0x0)
		if err != nil {
			return err
		}
		err = en.WriteBytes(z.OldUserKey)
		if err != nil {
			return
		}
	}

	if !empty_zhji[1] {
		// zid 1 for "NewUserKey"
		err = en.Append(0x1)
		if err != nil {
			return err
		}
		err = en.WriteBytes(z.NewUserKey)
		if err != nil {
			return
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z *EmailSessData) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [2]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// zid 0 for "OldUserKey"
		o = append(o, 0x0)
		o = msgp.AppendBytes(o, z.OldUserKey)
	}

	if !empty[1] {
		// zid 1 for "NewUserKey"
		o = append(o, 0x1)
		o = msgp.AppendBytes(o, z.NewUserKey)
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EmailSessData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *EmailSessData) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields1zegf = 2

	// -- templateUnmarshalMsgZid starts here--
	var totalEncodedFields1zegf uint32
	if !nbs.AlwaysNil {
		totalEncodedFields1zegf, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft1zegf := totalEncodedFields1zegf
	missingFieldsLeft1zegf := maxFields1zegf - totalEncodedFields1zegf

	var nextMiss1zegf int = -1
	var found1zegf [maxFields1zegf]bool
	var curField1zegf int

doneWithStruct1zegf:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft1zegf > 0 || missingFieldsLeft1zegf > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft1zegf, missingFieldsLeft1zegf, msgp.ShowFound(found1zegf[:]), unmarshalMsgFieldOrder1zegf)
		if encodedFieldsLeft1zegf > 0 {
			encodedFieldsLeft1zegf--
			curField1zegf, bts, err = nbs.ReadIntBytes(bts)
			if err != nil {
				return
			}
		} else {
			//missing fields need handling
			if nextMiss1zegf < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss1zegf = 0
			}
			for nextMiss1zegf < maxFields1zegf && (found1zegf[nextMiss1zegf] || unmarshalMsgFieldSkip1zegf[nextMiss1zegf]) {
				nextMiss1zegf++
			}
			if nextMiss1zegf == maxFields1zegf {
				// filled all the empty fields!
				break doneWithStruct1zegf
			}
			missingFieldsLeft1zegf--
			curField1zegf = nextMiss1zegf
		}
		//fmt.Printf("switching on curField: '%v'\n", curField1zegf)
		switch curField1zegf {
		// -- templateUnmarshalMsgZid ends here --

		case 0:
			// zid 0 for "OldUserKey"
			found1zegf[0] = true
			if nbs.AlwaysNil || msgp.IsNil(bts) {
				if !nbs.AlwaysNil {
					bts = bts[1:]
				}
				z.OldUserKey = z.OldUserKey[:0]
			} else {
				z.OldUserKey, bts, err = nbs.ReadBytesBytes(bts, z.OldUserKey)

				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		case 1:
			// zid 1 for "NewUserKey"
			found1zegf[1] = true
			if nbs.AlwaysNil || msgp.IsNil(bts) {
				if !nbs.AlwaysNil {
					bts = bts[1:]
				}
				z.NewUserKey = z.NewUserKey[:0]
			} else {
				z.NewUserKey, bts, err = nbs.ReadBytesBytes(bts, z.NewUserKey)

				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss1zegf != -1 {
		bts = nbs.PopAlwaysNil()
	}

	if sawTopNil {
		bts = nbs.PopAlwaysNil()
	}
	o = bts
	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of EmailSessData
var unmarshalMsgFieldOrder1zegf = []string{"OldUserKey", "NewUserKey"}

var unmarshalMsgFieldSkip1zegf = []bool{false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *EmailSessData) Msgsize() (s int) {
	s = 1 + 11 + msgp.BytesPrefixSize + len(z.OldUserKey) + 11 + msgp.BytesPrefixSize + len(z.NewUserKey)
	return
}

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *VerifySessData) DecodeMsg(dc *msgp.Reader) (err error) {
	var sawTopNil bool
	if dc.IsNil() {
		sawTopNil = true
		err = dc.ReadNil()
		if err != nil {
			return
		}
		dc.PushAlwaysNil()
	}

	var field []byte
	_ = field
	const maxFields2zrhl = 4

	// -- templateDecodeMsgZid starts here--
	var totalEncodedFields2zrhl uint32
	totalEncodedFields2zrhl, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft2zrhl := totalEncodedFields2zrhl
	missingFieldsLeft2zrhl := maxFields2zrhl - totalEncodedFields2zrhl

	var nextMiss2zrhl int = -1
	var found2zrhl [maxFields2zrhl]bool
	var curField2zrhl int

doneWithStruct2zrhl:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft2zrhl > 0 || missingFieldsLeft2zrhl > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft2zrhl, missingFieldsLeft2zrhl, msgp.ShowFound(found2zrhl[:]), decodeMsgFieldOrder2zrhl)
		if encodedFieldsLeft2zrhl > 0 {
			encodedFieldsLeft2zrhl--
			curField2zrhl, err = dc.ReadInt()
			if err != nil {
				return
			}
		} else {
			//missing fields need handling
			if nextMiss2zrhl < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss2zrhl = 0
			}
			for nextMiss2zrhl < maxFields2zrhl && (found2zrhl[nextMiss2zrhl] || decodeMsgFieldSkip2zrhl[nextMiss2zrhl]) {
				nextMiss2zrhl++
			}
			if nextMiss2zrhl == maxFields2zrhl {
				// filled all the empty fields!
				break doneWithStruct2zrhl
			}
			missingFieldsLeft2zrhl--
			curField2zrhl = nextMiss2zrhl
		}
		//fmt.Printf("switching on curField: '%v'\n", curField2zrhl)
		switch curField2zrhl {
		// -- templateDecodeMsgZid ends here --

		case 0:
			// zid 0 for "UserKey"
			found2zrhl[0] = true
			z.UserKey, err = dc.ReadBytes(z.UserKey)
			if err != nil {
				return
			}
		case 1:
			// zid 1 for "FirstName"
			found2zrhl[1] = true
			z.FirstName, err = dc.ReadString()
			if err != nil {
				return
			}
		case 2:
			// zid 2 for "LastName"
			found2zrhl[2] = true
			z.LastName, err = dc.ReadString()
			if err != nil {
				return
			}
		case 3:
			// zid 3 for "Password"
			found2zrhl[3] = true
			z.Password, err = dc.ReadBytes(z.Password)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss2zrhl != -1 {
		dc.PopAlwaysNil()
	}

	if sawTopNil {
		dc.PopAlwaysNil()
	}

	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of VerifySessData
var decodeMsgFieldOrder2zrhl = []string{"UserKey", "FirstName", "LastName", "Password"}

var decodeMsgFieldSkip2zrhl = []bool{false, false, false, false}

// fieldsNotEmpty supports omitempty tags
func (z *VerifySessData) fieldsNotEmpty(isempty []bool) uint32 {
	if len(isempty) == 0 {
		return 4
	}
	var fieldsInUse uint32 = 4
	isempty[0] = (len(z.UserKey) == 0) // string, omitempty
	if isempty[0] {
		fieldsInUse--
	}
	isempty[1] = (len(z.FirstName) == 0) // string, omitempty
	if isempty[1] {
		fieldsInUse--
	}
	isempty[2] = (len(z.LastName) == 0) // string, omitempty
	if isempty[2] {
		fieldsInUse--
	}
	isempty[3] = (len(z.Password) == 0) // string, omitempty
	if isempty[3] {
		fieldsInUse--
	}

	return fieldsInUse
}

// EncodeMsg implements msgp.Encodable
func (z *VerifySessData) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// honor the omitempty tags
	var empty_zayd [4]bool
	fieldsInUse_zpaz := z.fieldsNotEmpty(empty_zayd[:])

	// map header
	err = en.WriteMapHeader(fieldsInUse_zpaz)
	if err != nil {
		return err
	}

	if !empty_zayd[0] {
		// zid 0 for "UserKey"
		err = en.Append(0x0)
		if err != nil {
			return err
		}
		err = en.WriteBytes(z.UserKey)
		if err != nil {
			return
		}
	}

	if !empty_zayd[1] {
		// zid 1 for "FirstName"
		err = en.Append(0x1)
		if err != nil {
			return err
		}
		err = en.WriteString(z.FirstName)
		if err != nil {
			return
		}
	}

	if !empty_zayd[2] {
		// zid 2 for "LastName"
		err = en.Append(0x2)
		if err != nil {
			return err
		}
		err = en.WriteString(z.LastName)
		if err != nil {
			return
		}
	}

	if !empty_zayd[3] {
		// zid 3 for "Password"
		err = en.Append(0x3)
		if err != nil {
			return err
		}
		err = en.WriteBytes(z.Password)
		if err != nil {
			return
		}
	}

	return
}

// MarshalMsg implements msgp.Marshaler
func (z *VerifySessData) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())

	// honor the omitempty tags
	var empty [4]bool
	fieldsInUse := z.fieldsNotEmpty(empty[:])
	o = msgp.AppendMapHeader(o, fieldsInUse)

	if !empty[0] {
		// zid 0 for "UserKey"
		o = append(o, 0x0)
		o = msgp.AppendBytes(o, z.UserKey)
	}

	if !empty[1] {
		// zid 1 for "FirstName"
		o = append(o, 0x1)
		o = msgp.AppendString(o, z.FirstName)
	}

	if !empty[2] {
		// zid 2 for "LastName"
		o = append(o, 0x2)
		o = msgp.AppendString(o, z.LastName)
	}

	if !empty[3] {
		// zid 3 for "Password"
		o = append(o, 0x3)
		o = msgp.AppendBytes(o, z.Password)
	}

	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VerifySessData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *VerifySessData) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields3zpmf = 4

	// -- templateUnmarshalMsgZid starts here--
	var totalEncodedFields3zpmf uint32
	if !nbs.AlwaysNil {
		totalEncodedFields3zpmf, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft3zpmf := totalEncodedFields3zpmf
	missingFieldsLeft3zpmf := maxFields3zpmf - totalEncodedFields3zpmf

	var nextMiss3zpmf int = -1
	var found3zpmf [maxFields3zpmf]bool
	var curField3zpmf int

doneWithStruct3zpmf:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft3zpmf > 0 || missingFieldsLeft3zpmf > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft3zpmf, missingFieldsLeft3zpmf, msgp.ShowFound(found3zpmf[:]), unmarshalMsgFieldOrder3zpmf)
		if encodedFieldsLeft3zpmf > 0 {
			encodedFieldsLeft3zpmf--
			curField3zpmf, bts, err = nbs.ReadIntBytes(bts)
			if err != nil {
				return
			}
		} else {
			//missing fields need handling
			if nextMiss3zpmf < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss3zpmf = 0
			}
			for nextMiss3zpmf < maxFields3zpmf && (found3zpmf[nextMiss3zpmf] || unmarshalMsgFieldSkip3zpmf[nextMiss3zpmf]) {
				nextMiss3zpmf++
			}
			if nextMiss3zpmf == maxFields3zpmf {
				// filled all the empty fields!
				break doneWithStruct3zpmf
			}
			missingFieldsLeft3zpmf--
			curField3zpmf = nextMiss3zpmf
		}
		//fmt.Printf("switching on curField: '%v'\n", curField3zpmf)
		switch curField3zpmf {
		// -- templateUnmarshalMsgZid ends here --

		case 0:
			// zid 0 for "UserKey"
			found3zpmf[0] = true
			if nbs.AlwaysNil || msgp.IsNil(bts) {
				if !nbs.AlwaysNil {
					bts = bts[1:]
				}
				z.UserKey = z.UserKey[:0]
			} else {
				z.UserKey, bts, err = nbs.ReadBytesBytes(bts, z.UserKey)

				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		case 1:
			// zid 1 for "FirstName"
			found3zpmf[1] = true
			z.FirstName, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case 2:
			// zid 2 for "LastName"
			found3zpmf[2] = true
			z.LastName, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case 3:
			// zid 3 for "Password"
			found3zpmf[3] = true
			if nbs.AlwaysNil || msgp.IsNil(bts) {
				if !nbs.AlwaysNil {
					bts = bts[1:]
				}
				z.Password = z.Password[:0]
			} else {
				z.Password, bts, err = nbs.ReadBytesBytes(bts, z.Password)

				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss3zpmf != -1 {
		bts = nbs.PopAlwaysNil()
	}

	if sawTopNil {
		bts = nbs.PopAlwaysNil()
	}
	o = bts
	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of VerifySessData
var unmarshalMsgFieldOrder3zpmf = []string{"UserKey", "FirstName", "LastName", "Password"}

var unmarshalMsgFieldSkip3zpmf = []bool{false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *VerifySessData) Msgsize() (s int) {
	s = 1 + 8 + msgp.BytesPrefixSize + len(z.UserKey) + 10 + msgp.StringPrefixSize + len(z.FirstName) + 9 + msgp.StringPrefixSize + len(z.LastName) + 9 + msgp.BytesPrefixSize + len(z.Password)
	return
}

// FileSessdata_generated_go holds ZebraPack schema from file 'sessdata.go'
type FileSessdata_generated_go struct{}

// ZebraSchemaInMsgpack2Format provides the ZebraPack Schema in msgpack2 format, length 969 bytes
func (FileSessdata_generated_go) ZebraSchemaInMsgpack2Format() []byte {
	return []byte{
		0x85, 0xaa, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61,
		0x74, 0x68, 0xab, 0x73, 0x65, 0x73, 0x73, 0x64, 0x61, 0x74,
		0x61, 0x2e, 0x67, 0x6f, 0xad, 0x53, 0x6f, 0x75, 0x72, 0x63,
		0x65, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0xa3, 0x67,
		0x65, 0x6e, 0xad, 0x5a, 0x65, 0x62, 0x72, 0x61, 0x53, 0x63,
		0x68, 0x65, 0x6d, 0x61, 0x49, 0x64, 0x00, 0xa7, 0x53, 0x74,
		0x72, 0x75, 0x63, 0x74, 0x73, 0x82, 0xae, 0x56, 0x65, 0x72,
		0x69, 0x66, 0x79, 0x53, 0x65, 0x73, 0x73, 0x44, 0x61, 0x74,
		0x61, 0x82, 0xaa, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x4e,
		0x61, 0x6d, 0x65, 0xae, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
		0x53, 0x65, 0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0xa6, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x73, 0x94, 0x87, 0xa3, 0x5a, 0x69,
		0x64, 0x00, 0xab, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x47, 0x6f,
		0x4e, 0x61, 0x6d, 0x65, 0xa7, 0x55, 0x73, 0x65, 0x72, 0x4b,
		0x65, 0x79, 0xac, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x61,
		0x67, 0x4e, 0x61, 0x6d, 0x65, 0xa7, 0x55, 0x73, 0x65, 0x72,
		0x4b, 0x65, 0x79, 0xac, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54,
		0x79, 0x70, 0x65, 0x53, 0x74, 0x72, 0xa6, 0x5b, 0x5d, 0x62,
		0x79, 0x74, 0x65, 0xad, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x43,
		0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x17, 0xae, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74,
		0x69, 0x76, 0x65, 0x01, 0xad, 0x46, 0x69, 0x65, 0x6c, 0x64,
		0x46, 0x75, 0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x82, 0xa4,
		0x4b, 0x69, 0x6e, 0x64, 0x01, 0xa3, 0x53, 0x74, 0x72, 0xa5,
		0x62, 0x79, 0x74, 0x65, 0x73, 0x87, 0xa3, 0x5a, 0x69, 0x64,
		0x01, 0xab, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x47, 0x6f, 0x4e,
		0x61, 0x6d, 0x65, 0xa9, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4e,
		0x61, 0x6d, 0x65, 0xac, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54,
		0x61, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0xa9, 0x46, 0x69, 0x72,
		0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0xac, 0x46, 0x69, 0x65,
		0x6c, 0x64, 0x54, 0x79, 0x70, 0x65, 0x53, 0x74, 0x72, 0xa6,
		0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0xad, 0x46, 0x69, 0x65,
		0x6c, 0x64, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
		0x17, 0xae, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x50, 0x72, 0x69,
		0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x02, 0xad, 0x46, 0x69,
		0x65, 0x6c, 0x64, 0x46, 0x75, 0x6c, 0x6c, 0x54, 0x79, 0x70,
		0x65, 0x82, 0xa4, 0x4b, 0x69, 0x6e, 0x64, 0x02, 0xa3, 0x53,
		0x74, 0x72, 0xa6, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x87,
		0xa3, 0x5a, 0x69, 0x64, 0x02, 0xab, 0x46, 0x69, 0x65, 0x6c,
		0x64, 0x47, 0x6f, 0x4e, 0x61, 0x6d, 0x65, 0xa8, 0x4c, 0x61,
		0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0xac, 0x46, 0x69, 0x65,
		0x6c, 0x64, 0x54, 0x61, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0xa8,
		0x4c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0xac, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x54, 0x79, 0x70, 0x65, 0x53, 0x74,
		0x72, 0xa6, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0xad, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
		0x72, 0x79, 0x17, 0xae, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x50,
		0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x02, 0xad,
		0x46, 0x69, 0x65, 0x6c, 0x64, 0x46, 0x75, 0x6c, 0x6c, 0x54,
		0x79, 0x70, 0x65, 0x82, 0xa4, 0x4b, 0x69, 0x6e, 0x64, 0x02,
		0xa3, 0x53, 0x74, 0x72, 0xa6, 0x73, 0x74, 0x72, 0x69, 0x6e,
		0x67, 0x87, 0xa3, 0x5a, 0x69, 0x64, 0x03, 0xab, 0x46, 0x69,
		0x65, 0x6c, 0x64, 0x47, 0x6f, 0x4e, 0x61, 0x6d, 0x65, 0xa8,
		0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0xac, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x54, 0x61, 0x67, 0x4e, 0x61, 0x6d,
		0x65, 0xa8, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
		0xac, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x79, 0x70, 0x65,
		0x53, 0x74, 0x72, 0xa6, 0x5b, 0x5d, 0x62, 0x79, 0x74, 0x65,
		0xad, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x43, 0x61, 0x74, 0x65,
		0x67, 0x6f, 0x72, 0x79, 0x17, 0xae, 0x46, 0x69, 0x65, 0x6c,
		0x64, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65,
		0x01, 0xad, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x46, 0x75, 0x6c,
		0x6c, 0x54, 0x79, 0x70, 0x65, 0x82, 0xa4, 0x4b, 0x69, 0x6e,
		0x64, 0x01, 0xa3, 0x53, 0x74, 0x72, 0xa5, 0x62, 0x79, 0x74,
		0x65, 0x73, 0xad, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x65,
		0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0x82, 0xaa, 0x53, 0x74,
		0x72, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0xad, 0x45,
		0x6d, 0x61, 0x69, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x44, 0x61,
		0x74, 0x61, 0xa6, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x92,
		0x87, 0xa3, 0x5a, 0x69, 0x64, 0x00, 0xab, 0x46, 0x69, 0x65,
		0x6c, 0x64, 0x47, 0x6f, 0x4e, 0x61, 0x6d, 0x65, 0xaa, 0x4f,
		0x6c, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4b, 0x65, 0x79, 0xac,
		0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x61, 0x67, 0x4e, 0x61,
		0x6d, 0x65, 0xaa, 0x4f, 0x6c, 0x64, 0x55, 0x73, 0x65, 0x72,
		0x4b, 0x65, 0x79, 0xac, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54,
		0x79, 0x70, 0x65, 0x53, 0x74, 0x72, 0xa6, 0x5b, 0x5d, 0x62,
		0x79, 0x74, 0x65, 0xad, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x43,
		0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x17, 0xae, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74,
		0x69, 0x76, 0x65, 0x01, 0xad, 0x46, 0x69, 0x65, 0x6c, 0x64,
		0x46, 0x75, 0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x82, 0xa4,
		0x4b, 0x69, 0x6e, 0x64, 0x01, 0xa3, 0x53, 0x74, 0x72, 0xa5,
		0x62, 0x79, 0x74, 0x65, 0x73, 0x87, 0xa3, 0x5a, 0x69, 0x64,
		0x01, 0xab, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x47, 0x6f, 0x4e,
		0x61, 0x6d, 0x65, 0xaa, 0x4e, 0x65, 0x77, 0x55, 0x73, 0x65,
		0x72, 0x4b, 0x65, 0x79, 0xac, 0x46, 0x69, 0x65, 0x6c, 0x64,
		0x54, 0x61, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0xaa, 0x4e, 0x65,
		0x77, 0x55, 0x73, 0x65, 0x72, 0x4b, 0x65, 0x79, 0xac, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x54, 0x79, 0x70, 0x65, 0x53, 0x74,
		0x72, 0xa6, 0x5b, 0x5d, 0x62, 0x79, 0x74, 0x65, 0xad, 0x46,
		0x69, 0x65, 0x6c, 0x64, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
		0x72, 0x79, 0x17, 0xae, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x50,
		0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x01, 0xad,
		0x46, 0x69, 0x65, 0x6c, 0x64, 0x46, 0x75, 0x6c, 0x6c, 0x54,
		0x79, 0x70, 0x65, 0x82, 0xa4, 0x4b, 0x69, 0x6e, 0x64, 0x01,
		0xa3, 0x53, 0x74, 0x72, 0xa5, 0x62, 0x79, 0x74, 0x65, 0x73,
		0xa7, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x90,
	}
}

// ZebraSchemaInJsonCompact provides the ZebraPack Schema in compact JSON format, length 1205 bytes
func (FileSessdata_generated_go) ZebraSchemaInJsonCompact() []byte {
	return []byte(`{"SourcePath":"sessdata.go","SourcePackage":"gen","ZebraSchemaId":0,"Structs":{"VerifySessData":{"StructName":"VerifySessData","Fields":[{"Zid":0,"FieldGoName":"UserKey","FieldTagName":"UserKey","FieldTypeStr":"[]byte","FieldCategory":23,"FieldPrimitive":1,"FieldFullType":{"Kind":1,"Str":"bytes"}},{"Zid":1,"FieldGoName":"FirstName","FieldTagName":"FirstName","FieldTypeStr":"string","FieldCategory":23,"FieldPrimitive":2,"FieldFullType":{"Kind":2,"Str":"string"}},{"Zid":2,"FieldGoName":"LastName","FieldTagName":"LastName","FieldTypeStr":"string","FieldCategory":23,"FieldPrimitive":2,"FieldFullType":{"Kind":2,"Str":"string"}},{"Zid":3,"FieldGoName":"Password","FieldTagName":"Password","FieldTypeStr":"[]byte","FieldCategory":23,"FieldPrimitive":1,"FieldFullType":{"Kind":1,"Str":"bytes"}}]},"EmailSessData":{"StructName":"EmailSessData","Fields":[{"Zid":0,"FieldGoName":"OldUserKey","FieldTagName":"OldUserKey","FieldTypeStr":"[]byte","FieldCategory":23,"FieldPrimitive":1,"FieldFullType":{"Kind":1,"Str":"bytes"}},{"Zid":1,"FieldGoName":"NewUserKey","FieldTagName":"NewUserKey","FieldTypeStr":"[]byte","FieldCategory":23,"FieldPrimitive":1,"FieldFullType":{"Kind":1,"Str":"bytes"}}]}},"Imports":[]}`)
}

// ZebraSchemaInJsonPretty provides the ZebraPack Schema in pretty JSON format, length 2927 bytes
func (FileSessdata_generated_go) ZebraSchemaInJsonPretty() []byte {
	return []byte(`{
    "SourcePath": "sessdata.go",
    "SourcePackage": "gen",
    "ZebraSchemaId": 0,
    "Structs": {
        "VerifySessData": {
            "StructName": "VerifySessData",
            "Fields": [
                {
                    "Zid": 0,
                    "FieldGoName": "UserKey",
                    "FieldTagName": "UserKey",
                    "FieldTypeStr": "[]byte",
                    "FieldCategory": 23,
                    "FieldPrimitive": 1,
                    "FieldFullType": {
                        "Kind": 1,
                        "Str": "bytes"
                    }
                },
                {
                    "Zid": 1,
                    "FieldGoName": "FirstName",
                    "FieldTagName": "FirstName",
                    "FieldTypeStr": "string",
                    "FieldCategory": 23,
                    "FieldPrimitive": 2,
                    "FieldFullType": {
                        "Kind": 2,
                        "Str": "string"
                    }
                },
                {
                    "Zid": 2,
                    "FieldGoName": "LastName",
                    "FieldTagName": "LastName",
                    "FieldTypeStr": "string",
                    "FieldCategory": 23,
                    "FieldPrimitive": 2,
                    "FieldFullType": {
                        "Kind": 2,
                        "Str": "string"
                    }
                },
                {
                    "Zid": 3,
                    "FieldGoName": "Password",
                    "FieldTagName": "Password",
                    "FieldTypeStr": "[]byte",
                    "FieldCategory": 23,
                    "FieldPrimitive": 1,
                    "FieldFullType": {
                        "Kind": 1,
                        "Str": "bytes"
                    }
                }
            ]
        },
        "EmailSessData": {
            "StructName": "EmailSessData",
            "Fields": [
                {
                    "Zid": 0,
                    "FieldGoName": "OldUserKey",
                    "FieldTagName": "OldUserKey",
                    "FieldTypeStr": "[]byte",
                    "FieldCategory": 23,
                    "FieldPrimitive": 1,
                    "FieldFullType": {
                        "Kind": 1,
                        "Str": "bytes"
                    }
                },
                {
                    "Zid": 1,
                    "FieldGoName": "NewUserKey",
                    "FieldTagName": "NewUserKey",
                    "FieldTypeStr": "[]byte",
                    "FieldCategory": 23,
                    "FieldPrimitive": 1,
                    "FieldFullType": {
                        "Kind": 1,
                        "Str": "bytes"
                    }
                }
            ]
        }
    },
    "Imports": []
}`)
}
