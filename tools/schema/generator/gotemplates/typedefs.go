package gotemplates

var typedefsGo = map[string]string{
	// *******************************
	"typedefs.go": `
$#emit goHeader
$#each typedef typedefProxy
`,
	// *******************************
	"typedefProxy": `
$#set mut Immutable
$#if array typedefProxyArray
$#if array typedefProxyAlias
$#if map typedefProxyMap
$#if map typedefProxyAlias
$#set mut Mutable
$#if array typedefProxyArray
$#if array typedefProxyAlias
$#if map typedefProxyMap
$#if map typedefProxyAlias
`,
	// *******************************
	"typedefProxyAlias": `

type $mut$FldName = $proxy
`,
	// *******************************
	"typedefProxyArray": `
$#set proxy ArrayOf$mut$FldType
$#if exist else typedefProxyArrayNew
`,
	// *******************************
	"typedefProxyArrayNew": `

type $proxy struct {
	objID int32
}
$#if mut typedefProxyArrayClear

func (a $proxy) Length() int32 {
	return wasmlib.GetLength(a.objID)
}
$#if basetype typedefProxyArrayNewBaseType typedefProxyArrayNewOtherType
$#set exist $proxy
`,
	// *******************************
	"typedefProxyArrayClear": `

func (a $proxy) Clear() {
	wasmlib.Clear(a.objID)
}
`,
	// *******************************
	"typedefProxyArrayNewBaseType": `

func (a $proxy) Get$FldType(index int32) wasmlib.Sc$mut$FldType {
	return wasmlib.NewSc$mut$FldType(a.objID, wasmlib.Key32(index))
}
`,
	// *******************************
	"typedefProxyArrayNewOtherType": `
$#set OldType $FldType
$#if typedef typedefProxyArrayNewOtherTypeTypeDef typedefProxyArrayNewOtherTypeStruct
`,
	// *******************************
	"typedefProxyArrayNewOtherTypeTypeDef": `
$#set varType wasmlib.TYPE_MAP
$#if array setVarTypeArray

func (a $proxy) Get$OldType(index int32) $mut$OldType {
	subID := wasmlib.GetObjectID(a.objID, wasmlib.Key32(index), $varType)
	return $mut$OldType{objID: subID}
}
`,
	// *******************************
	"typedefProxyArrayNewOtherTypeStruct": `

func (a $proxy) Get$FldType(index int32) $mut$FldType {
	return $mut$FldType{objID: a.objID, keyID: wasmlib.Key32(index)}
}
`,
	// *******************************
	"typedefProxyMap": `
$#set proxy Map$fldMapKey$+To$mut$FldType
$#if exist else typedefProxyMapNew
`,
	// *******************************
	"typedefProxyMapNew": `

type $proxy struct {
	objID int32
}
$#if mut typedefProxyMapClear
$#if basetype typedefProxyMapNewBaseType typedefProxyMapNewOtherType
$#set exist $proxy
`,
	// *******************************
	"typedefProxyMapClear": `

func (m $proxy) Clear() {
	wasmlib.Clear(m.objID)
}
`,
	// *******************************
	"typedefProxyMapNewBaseType": `

func (m $proxy) Get$FldType(key $fldMapKeyLangType) wasmlib.Sc$mut$FldType {
	return wasmlib.NewSc$mut$FldType(m.objID, $fldMapKeyKey.KeyID())
}
`,
	// *******************************
	"typedefProxyMapNewOtherType": `
$#set OldType $FldType
$#set oldMapKeyLangType $fldMapKeyLangType
$#set oldMapKeyKey $fldMapKeyKey
$#if typedef typedefProxyMapNewOtherTypeTypeDef typedefProxyMapNewOtherTypeStruct
`,
	// *******************************
	"typedefProxyMapNewOtherTypeTypeDef": `
$#set varType wasmlib.TYPE_MAP
$#if array setVarTypeArray

func (m $proxy) Get$OldType(key $oldMapKeyLangType) $mut$OldType {
	subID := wasmlib.GetObjectID(m.objID, $oldMapKeyKey.KeyID(), $varType)
	return $mut$OldType{objID: subID}
}
`,
	// *******************************
	"typedefProxyMapNewOtherTypeStruct": `

func (m $proxy) Get$FldType(key $fldMapKeyLangType) $mut$FldType {
	return $mut$FldType{objID: m.objID, keyID: $fldMapKeyKey.KeyID()}
}
`,
	// *******************************
	"setVarTypeArray": `
$#set varType $arrayTypeID|$fldTypeID
`,
}
