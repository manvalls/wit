package wit

import "strconv"

const (
	_ = iota

	// Standard

	sliceType

	rootType
	selectorType
	selectorAllType
	parentType
	firstChildType
	lastChildType
	prevSiblingType
	nextSiblingType

	removeType
	clearType

	htmlType
	appendType
	prependType
	insertAfterType
	insertBeforeType

	addAttrType
	setAttrType
	rmAttrType
	addStylesType
	rmStylesType
	addClassType
	rmClassType

	callType

	// Flow control

	channelType
	errorType
	runSyncType
)

var (
	sliceTypeString        = []byte(strconv.Itoa(sliceType))
	rootTypeString         = []byte(strconv.Itoa(rootType))
	selectorTypeString     = []byte(strconv.Itoa(selectorType))
	selectorAllTypeString  = []byte(strconv.Itoa(selectorAllType))
	parentTypeString       = []byte(strconv.Itoa(parentType))
	firstChildTypeString   = []byte(strconv.Itoa(firstChildType))
	lastChildTypeString    = []byte(strconv.Itoa(lastChildType))
	prevSiblingTypeString  = []byte(strconv.Itoa(prevSiblingType))
	nextSiblingTypeString  = []byte(strconv.Itoa(nextSiblingType))
	removeTypeString       = []byte(strconv.Itoa(removeType))
	clearTypeString        = []byte(strconv.Itoa(clearType))
	htmlTypeString         = []byte(strconv.Itoa(htmlType))
	appendTypeString       = []byte(strconv.Itoa(appendType))
	prependTypeString      = []byte(strconv.Itoa(prependType))
	insertAfterTypeString  = []byte(strconv.Itoa(insertAfterType))
	insertBeforeTypeString = []byte(strconv.Itoa(insertBeforeType))
	addAttrTypeString      = []byte(strconv.Itoa(addAttrType))
	setAttrTypeString      = []byte(strconv.Itoa(setAttrType))
	rmAttrTypeString       = []byte(strconv.Itoa(rmAttrType))
	addStylesTypeString    = []byte(strconv.Itoa(addStylesType))
	rmStylesTypeString     = []byte(strconv.Itoa(rmStylesType))
	addClassTypeString     = []byte(strconv.Itoa(addClassType))
	rmClassTypeString      = []byte(strconv.Itoa(rmClassType))
	callTypeString         = []byte(strconv.Itoa(callType))
)
