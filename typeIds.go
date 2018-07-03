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
	cssType

	callType

	jumpType
	runSyncType
	deferType

	statusType
	redirectType
	addHeadersType
	setHeadersType
	rmHeadersType
	answerType
)
