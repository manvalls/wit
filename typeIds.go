package wit

const (
	_ = iota

	sliceType
	channelType

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
	textType
	replaceType
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

	jsType
	asyncJSType
	cssType
	asyncCSSType

	callType

	jumpType
	runSyncType
	withKeyType
	clearKeyType
	deferType

	statusType
	redirectType
	addHeadersType
	setHeadersType
	rmHeadersType
	answerType
)
