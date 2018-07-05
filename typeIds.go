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

	callType

	jumpType
	runSyncType

	statusType
	addHeadersType
	setHeadersType
	rmHeadersType
	answerType
)
